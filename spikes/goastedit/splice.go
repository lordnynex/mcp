package goastedit

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"strconv"
	"strings"
)

// Op is one edit the spike can apply.
type Op string

const (
	OpReplaceBody  Op = "replace_body"
	OpReplaceDecl  Op = "replace_decl"
	OpInsertBefore Op = "insert_before"
	OpInsertAfter  Op = "insert_after"
	OpPrependStmt  Op = "prepend_stmt"
	OpDelete       Op = "delete"
	OpAddImport    Op = "add_import"
)

// Request is a logical edit. Content is Go source. Import is a package path
// without quotes.
type Request struct {
	Op      Op
	Target  string
	Content string
	Import  string
}

// Apply splices req into src and gofmts the file. The splice uses original
// bytes; it does not rebuild the file from syntax nodes.
func Apply(src []byte, req Request) ([]byte, error) {
	chg, err := plan(src, req)
	if err != nil {
		return nil, err
	}
	raw := spliceBytes(src, chg.start, chg.end, chg.repl)
	out, err := format.Source(raw)
	if err != nil {
		return nil, fmt.Errorf("format edited source: %w", err)
	}
	return out, nil
}

// change is the minimal byte edit. quote is the region a search/replace edit
// must repeat: insertions quote the enclosing declaration or body, because an
// empty old string is not a usable text match.
type change struct {
	start, end           int
	repl                 string
	quoteStart, quoteEnd int
	quoteNew             string
}

func plan(src []byte, req Request) (change, error) {
	if req.Op == OpAddImport {
		return planImport(src, req.Import)
	}
	idx, err := Parse(src)
	if err != nil {
		return change{}, err
	}
	d, err := idx.Lookup(req.Target)
	if err != nil {
		return change{}, err
	}
	switch req.Op {
	case OpReplaceBody:
		if err := needContent(req); err != nil {
			return change{}, err
		}
		if d.Body.End <= d.Body.Start {
			return change{}, fmt.Errorf("%s has no body", d.Addr)
		}
		return change{
			start: d.Body.Start, end: d.Body.End, repl: req.Content,
			quoteStart: d.Body.Start, quoteEnd: d.Body.End, quoteNew: req.Content,
		}, nil
	case OpReplaceDecl:
		if err := needContent(req); err != nil {
			return change{}, err
		}
		full := d.Full()
		return change{
			start: full.Start, end: full.End, repl: req.Content,
			quoteStart: full.Start, quoteEnd: full.End, quoteNew: req.Content,
		}, nil
	case OpInsertBefore:
		if err := needContent(req); err != nil {
			return change{}, err
		}
		full := d.Full()
		repl := strings.TrimRight(req.Content, "\n") + "\n\n"
		old := idx.Text(full)
		return change{
			start: full.Start, end: full.Start, repl: repl,
			quoteStart: full.Start, quoteEnd: full.End, quoteNew: repl + old,
		}, nil
	case OpInsertAfter:
		if err := needContent(req); err != nil {
			return change{}, err
		}
		full := d.Full()
		repl := "\n\n" + strings.TrimRight(req.Content, "\n") + "\n"
		old := idx.Text(full)
		return change{
			start: full.End, end: full.End, repl: repl,
			quoteStart: full.Start, quoteEnd: full.End, quoteNew: old + repl,
		}, nil
	case OpPrependStmt:
		if err := needContent(req); err != nil {
			return change{}, err
		}
		if d.Body.End <= d.Body.Start {
			return change{}, fmt.Errorf("%s has no body", d.Addr)
		}
		body := idx.Text(d.Body)
		stmt := strings.TrimRight(req.Content, "\n")
		repl := "\n" + stmt + "\n"
		quoteNew := body[:1] + repl + body[1:]
		return change{
			start: d.Body.Start + 1, end: d.Body.Start + 1, repl: repl,
			quoteStart: d.Body.Start, quoteEnd: d.Body.End, quoteNew: quoteNew,
		}, nil
	case OpDelete:
		full := d.Full()
		return change{
			start: full.Start, end: full.End, repl: "",
			quoteStart: full.Start, quoteEnd: full.End, quoteNew: "",
		}, nil
	default:
		return change{}, fmt.Errorf("unknown op %q", req.Op)
	}
}

func needContent(req Request) error {
	if strings.TrimSpace(req.Content) == "" {
		return fmt.Errorf("%s requires content", req.Op)
	}
	return nil
}

func planImport(src []byte, path string) (change, error) {
	if path == "" || strings.ContainsAny(path, "\"\n\r") {
		return change{}, fmt.Errorf("invalid import path %q", path)
	}
	idx, err := Parse(src)
	if err != nil {
		return change{}, err
	}
	gen := importGen(idx.AST)
	if gen == nil {
		return change{}, fmt.Errorf("file has no import declaration")
	}
	quoted := strconv.Quote(path)
	for _, spec := range gen.Specs {
		im, ok := spec.(*ast.ImportSpec)
		if !ok || im.Path == nil {
			continue
		}
		got, err := strconv.Unquote(im.Path.Value)
		if err != nil {
			continue
		}
		if got == path {
			full := spanOf(idx.Fset, gen.Pos(), gen.End())
			text := idx.Text(full)
			return change{
				start: full.End, end: full.End, repl: "",
				quoteStart: full.Start, quoteEnd: full.End, quoteNew: text,
			}, nil
		}
	}
	if gen.Rparen.IsValid() {
		at := offset(idx.Fset, gen.Rparen)
		full := spanOf(idx.Fset, gen.Pos(), gen.End())
		old := idx.Text(full)
		line := "\t" + quoted + "\n"
		rel := at - full.Start
		return change{
			start: at, end: at, repl: line,
			quoteStart: full.Start, quoteEnd: full.End, quoteNew: old[:rel] + line + old[rel:],
		}, nil
	}
	// A single import spec becomes a group.
	if len(gen.Specs) != 1 {
		return change{}, fmt.Errorf("unparenthesized import has %d specs", len(gen.Specs))
	}
	im, ok := gen.Specs[0].(*ast.ImportSpec)
	if !ok || im.Path == nil {
		return change{}, fmt.Errorf("import spec is not a path")
	}
	full := spanOf(idx.Fset, gen.Pos(), gen.End())
	repl := "import (\n\t" + im.Path.Value + "\n\t" + quoted + "\n)"
	return change{
		start: full.Start, end: full.End, repl: repl,
		quoteStart: full.Start, quoteEnd: full.End, quoteNew: repl,
	}, nil
}

func importGen(file *ast.File) *ast.GenDecl {
	for _, d := range file.Decls {
		g, ok := d.(*ast.GenDecl)
		if ok && g.Tok == token.IMPORT {
			return g
		}
	}
	return nil
}

func spliceBytes(src []byte, start, end int, repl string) []byte {
	out := make([]byte, 0, len(src)-(end-start)+len(repl))
	out = append(out, src[:start]...)
	out = append(out, repl...)
	out = append(out, src[end:]...)
	return out
}
