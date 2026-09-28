package goastedit

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Scenario is one logical edit measured on the fixture.
type Scenario struct {
	Name string
	Req  Request
}

// Scenarios are the edits compared by the report. New Go source is written once
// here and reused by every arm.
func Scenarios() []Scenario {
	return []Scenario{
		{
			Name: "replace body",
			Req: Request{
				Op:     OpReplaceBody,
				Target: "(*Box).Get",
				Content: `{
	if key == "" {
		return "missing"
	}
	val := b.items[key]
	return strings.TrimSpace(val)
}`,
			},
		},
		{
			Name: "add method",
			Req: Request{
				Op:     OpInsertAfter,
				Target: "(*Box).Put",
				Content: `// Has reports whether key is present.
func (b *Box) Has(key string) bool {
	_, ok := b.items[key]
	return ok
}`,
			},
		},
		{
			Name: "change signature",
			Req: Request{
				Op:     OpReplaceDecl,
				Target: "Helper",
				Content: `// Helper formats a label with an optional prefix.
func Helper(prefix, label string) string {
	label = strings.TrimSpace(label)
	if label == "" {
		return prefix + "empty"
	}
	return fmt.Sprintf("%s=%s", prefix, label)
}`,
			},
		},
		{
			Name: "delete helper",
			Req:  Request{Op: OpDelete, Target: "Helper"},
		},
		{
			Name: "add import",
			Req:  Request{Op: OpAddImport, Import: "strconv"},
		},
		{
			Name: "prepend statement",
			Req: Request{
				Op:     OpPrependStmt,
				Target: "(*Box).Put",
				Content: `if strings.Contains(key, " ") {
	return
}`,
			},
		},
	}
}

// Report renders the comparison for fixture and an additional read-only file.
func Report(fixture []byte, largestLabel string, largest []byte) (string, error) {
	var b strings.Builder
	b.WriteString("## Schema (per session)\n\n")
	b.WriteString("| Arm | Bytes | Tokens | Role |\n| --- | ---: | ---: | --- |\n")
	schemas := []struct {
		name, text, role string
	}{
		{"search/replace", SchemaSearch, "baseline already paid"},
		{"unified diff", SchemaDiff, "baseline already paid"},
		{"symbol splice", SchemaSplice, "incremental"},
		{"ast json", SchemaAST, "incremental"},
		{"intent", SchemaIntent, "incremental"},
	}
	for _, s := range schemas {
		if err := writeMeasure(&b, s.name, s.text, s.role); err != nil {
			return "", err
		}
	}

	b.WriteString("\n## Read\n\n")
	b.WriteString("| File | View | Bytes | Tokens |\n| --- | --- | ---: | ---: |\n")
	if err := writeRead(&b, "sample.go", fixture); err != nil {
		return "", err
	}
	if err := writeRead(&b, largestLabel, largest); err != nil {
		return "", err
	}

	b.WriteString("\n## Model output\n\n")
	b.WriteString("| Scenario | Arm | Bytes | Tokens |\n| --- | --- | ---: | ---: |\n")
	var results strings.Builder
	results.WriteString("\n## Tool result\n\n")
	results.WriteString("| Scenario | Result | Bytes | Tokens |\n| --- | --- | ---: | ---: |\n")
	for _, sc := range Scenarios() {
		arms, err := Payloads(fixture, sc.Req)
		if err != nil {
			return "", fmt.Errorf("%s: %w", sc.Name, err)
		}
		for _, arm := range arms {
			if arm.NA {
				fmt.Fprintf(&b, "| %s | %s | n/a | n/a |\n", sc.Name, arm.Name)
				continue
			}
			u, err := Measure(arm.Payload)
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&b, "| %s | %s | %d | %d |\n", sc.Name, arm.Name, u.Bytes, u.Tokens)
		}
		res, err := ToolResult(fixture, sc.Req)
		if err != nil {
			return "", fmt.Errorf("%s result: %w", sc.Name, err)
		}
		for _, row := range []struct{ name, text string }{
			{"summary", res.Summary},
			{"diff", res.Diff},
			{"file", res.File},
			{"ast json", res.AST},
		} {
			u, err := Measure(row.text)
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&results, "| %s | %s | %d | %d |\n", sc.Name, row.name, u.Bytes, u.Tokens)
		}
	}
	b.WriteString(results.String())
	return b.String(), nil
}

func writeMeasure(b *strings.Builder, name, text, role string) error {
	u, err := Measure(text)
	if err != nil {
		return err
	}
	fmt.Fprintf(b, "| %s | %d | %d | %s |\n", name, u.Bytes, u.Tokens, role)
	return nil
}

func writeRead(b *strings.Builder, label string, src []byte) error {
	outline, err := Outline(src)
	if err != nil {
		return fmt.Errorf("outline %s: %w", label, err)
	}
	idx, err := Parse(src)
	if err != nil {
		return fmt.Errorf("parse %s: %w", label, err)
	}
	enc, err := Encode(idx.AST)
	if err != nil {
		return fmt.Errorf("encode %s: %w", label, err)
	}
	astJSON, err := marshal(enc)
	if err != nil {
		return err
	}
	for _, row := range []struct{ view, text string }{
		{"source", string(src)},
		{"outline", outline},
		{"ast json", astJSON},
	} {
		u, err := Measure(row.text)
		if err != nil {
			return err
		}
		fmt.Fprintf(b, "| %s | %s | %d | %d |\n", label, row.view, u.Bytes, u.Tokens)
	}
	return nil
}

// LargestGoFile returns the largest parseable .go file under root. Hidden
// directories and vendor are skipped so local scratch trees are not measured.
// The path is relative to root.
func LargestGoFile(root string) (string, []byte, error) {
	var bestPath string
	var best []byte
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if name == "vendor" || (name != "." && strings.HasPrefix(name, ".")) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		size := int(info.Size())
		if len(best) > 0 && (size < len(best) || (size == len(best) && path >= bestPath)) {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if _, err := Parse(src); err != nil {
			return nil
		}
		best = src
		bestPath = path
		return nil
	})
	if err != nil {
		return "", nil, err
	}
	if bestPath == "" {
		return "", nil, fmt.Errorf("no parseable Go file under %s", root)
	}
	rel, err := filepath.Rel(root, bestPath)
	if err != nil {
		return "", nil, err
	}
	return rel, best, nil
}
