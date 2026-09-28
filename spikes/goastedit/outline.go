package goastedit

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"
)

// Outline lists signatures and type names without bodies.
func Outline(src []byte) (string, error) {
	idx, err := Parse(src)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	for _, d := range idx.AST.Decls {
		switch n := d.(type) {
		case *ast.FuncDecl:
			decl, err := idx.Lookup(addrOf(n))
			if err != nil {
				return "", err
			}
			sig := strings.Join(strings.Fields(idx.Text(decl.Sig)), " ")
			b.WriteString(sig)
			b.WriteByte('\n')
		case *ast.GenDecl:
			for _, spec := range n.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					fmt.Fprintf(&b, "type %s %s\n", s.Name.Name, typeKind(idx.Fset, s.Type))
				case *ast.ValueSpec:
					kind := "var"
					if n.Tok == token.CONST {
						kind = "const"
					}
					for _, name := range s.Names {
						fmt.Fprintf(&b, "%s %s\n", kind, name.Name)
					}
				}
			}
		}
	}
	return b.String(), nil
}

func typeKind(fset *token.FileSet, expr ast.Expr) string {
	switch expr.(type) {
	case *ast.StructType:
		return "struct"
	case *ast.InterfaceType:
		return "interface"
	default:
		var buf bytes.Buffer
		if err := printer.Fprint(&buf, fset, expr); err != nil {
			return ""
		}
		return buf.String()
	}
}
