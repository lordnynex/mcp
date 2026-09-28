package goastedit

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
)

// Encode walks n and returns a JSON-ready tree. Node kinds are recorded, and
// positions and comments are left out so the tree is not charged for byte offsets.
func Encode(n ast.Node) (any, error) {
	if n == nil {
		return nil, fmt.Errorf("encode nil node")
	}
	v := encodeValue(reflect.ValueOf(n))
	if v == nil {
		return nil, fmt.Errorf("encode %T produced nothing", n)
	}
	return v, nil
}

var (
	posType   = reflect.TypeFor[token.Pos]()
	tokenType = reflect.TypeFor[token.Token]()
)

func encodeValue(rv reflect.Value) any {
	if !rv.IsValid() {
		return nil
	}
	if rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return nil
		}
		elem := rv.Elem()
		et := elem.Type()
		if et.PkgPath() == "go/ast" && et.Kind() == reflect.Struct {
			switch et.Name() {
			case "CommentGroup", "Comment", "Scope", "Object":
				return nil
			}
			return encodeStruct(et, elem)
		}
	}
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		if (rv.Kind() == reflect.Slice && rv.IsNil()) || rv.Len() == 0 {
			return nil
		}
		out := make([]any, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			ev := encodeValue(rv.Index(i))
			if ev == nil {
				continue
			}
			out = append(out, ev)
		}
		if len(out) == 0 {
			return nil
		}
		return out
	case reflect.String:
		if rv.String() == "" {
			return nil
		}
		return rv.String()
	case reflect.Bool:
		if !rv.Bool() {
			return nil
		}
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if rv.Type() == tokenType {
			tok := token.Token(rv.Int())
			if tok == token.ILLEGAL {
				return nil
			}
			return tok.String()
		}
		if rv.Type() == posType || rv.Int() == 0 {
			return nil
		}
		return rv.Int()
	default:
		return nil
	}
}

func encodeStruct(t reflect.Type, rv reflect.Value) map[string]any {
	m := map[string]any{"kind": t.Name()}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() || f.Type == posType {
			continue
		}
		// File.Imports repeats the import declaration already present in Decls.
		if t.Name() == "File" && (f.Name == "Imports" || f.Name == "Unresolved" || f.Name == "Comments") {
			continue
		}
		key := lowerFirst(f.Name)
		if key == "kind" {
			key = "lit"
		}
		ev := encodeValue(rv.Field(i))
		if ev == nil {
			continue
		}
		m[key] = ev
	}
	return m
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	return string(s[0]+('a'-'A')) + s[1:]
}

func parseSnippet(op Op, content, importPath string) (ast.Node, error) {
	switch op {
	case OpReplaceBody:
		return parseWrapped("package p\nfunc _() "+content, func(f *ast.File) ast.Node {
			return f.Decls[0].(*ast.FuncDecl).Body
		})
	case OpReplaceDecl, OpInsertBefore, OpInsertAfter:
		return parseWrapped("package p\n"+content, func(f *ast.File) ast.Node {
			return f.Decls[0]
		})
	case OpPrependStmt:
		return parseWrapped("package p\nfunc _() {\n"+content+"\n}", func(f *ast.File) ast.Node {
			list := f.Decls[0].(*ast.FuncDecl).Body.List
			if len(list) == 1 {
				return list[0]
			}
			return &ast.BlockStmt{List: list}
		})
	case OpAddImport:
		return &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(importPath)}}, nil
	case OpDelete:
		return nil, nil
	default:
		return nil, fmt.Errorf("no syntax for op %q", op)
	}
}

func parseWrapped(src string, pick func(*ast.File) ast.Node) (ast.Node, error) {
	f, err := parser.ParseFile(token.NewFileSet(), "", src, parser.SkipObjectResolution)
	if err != nil {
		return nil, fmt.Errorf("parse snippet: %w", err)
	}
	if len(f.Decls) == 0 {
		return nil, fmt.Errorf("parse snippet: no declarations")
	}
	return pick(f), nil
}
