package goastedit

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// Span is a half-open byte range in the original source.
type Span struct {
	Start int
	End   int
}

// Decl is one function or method.
type Decl struct {
	Addr string // Func, (*T).Method, or T.Method
	Name string
	Doc  Span // leading doc comment, or zero
	Decl Span // declaration without the doc comment
	Sig  Span // from func through the end of the signature
	Body Span // block including braces, or zero when there is no body
	Node *ast.FuncDecl
}

// Index is a parsed Go file and its function declarations.
type Index struct {
	Src   []byte
	Fset  *token.FileSet
	AST   *ast.File
	Decls []Decl
}

// Parse parses src as a single Go file. Comments stay attached so edits can
// keep or drop doc comments with the declaration they belong to.
func Parse(src []byte) (*Index, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "sample.go", src, parser.ParseComments|parser.SkipObjectResolution)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	decls, err := indexFuncs(fset, file)
	if err != nil {
		return nil, err
	}
	return &Index{Src: src, Fset: fset, AST: file, Decls: decls}, nil
}

func indexFuncs(fset *token.FileSet, file *ast.File) ([]Decl, error) {
	var decls []Decl
	seen := make(map[string]struct{}, len(file.Decls))
	for _, d := range file.Decls {
		fn, ok := d.(*ast.FuncDecl)
		if !ok || fn.Name == nil {
			continue
		}
		decl := declFrom(fset, fn)
		if _, ok := seen[decl.Addr]; ok {
			return nil, fmt.Errorf("duplicate declaration %s", decl.Addr)
		}
		seen[decl.Addr] = struct{}{}
		decls = append(decls, decl)
	}
	return decls, nil
}

func declFrom(fset *token.FileSet, fn *ast.FuncDecl) Decl {
	d := Decl{
		Addr: addrOf(fn),
		Name: fn.Name.Name,
		Decl: spanOf(fset, fn.Pos(), fn.End()),
		Node: fn,
	}
	if fn.Doc != nil {
		d.Doc = spanOf(fset, fn.Doc.Pos(), fn.Doc.End())
	}
	sigEnd := fn.End()
	if fn.Body != nil {
		sigEnd = fn.Body.Lbrace
		d.Body = spanOf(fset, fn.Body.Lbrace, fn.Body.End())
	}
	d.Sig = spanOf(fset, fn.Pos(), sigEnd)
	return d
}

func addrOf(fn *ast.FuncDecl) string {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return fn.Name.Name
	}
	switch recv := fn.Recv.List[0].Type.(type) {
	case *ast.StarExpr:
		return "(*" + exprName(recv.X) + ")." + fn.Name.Name
	default:
		return exprName(recv) + "." + fn.Name.Name
	}
}

func exprName(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + exprName(t.X)
	default:
		return "expr"
	}
}

func spanOf(fset *token.FileSet, start, end token.Pos) Span {
	return Span{Start: offset(fset, start), End: offset(fset, end)}
}

func offset(fset *token.FileSet, p token.Pos) int {
	if !p.IsValid() {
		return 0
	}
	return fset.Position(p).Offset
}

// Text returns the original source for s.
func (x *Index) Text(s Span) string {
	if s.End < s.Start || s.Start < 0 || s.End > len(x.Src) {
		return ""
	}
	return string(x.Src[s.Start:s.End])
}

// Lookup resolves a full address or a short name. A short name that matches
// more than one declaration returns an error naming each address.
func (x *Index) Lookup(name string) (Decl, error) {
	var exact []Decl
	var short []Decl
	for _, d := range x.Decls {
		if d.Addr == name {
			exact = append(exact, d)
		}
		if d.Name == name {
			short = append(short, d)
		}
	}
	hits := exact
	if len(hits) == 0 {
		hits = short
	}
	switch len(hits) {
	case 0:
		return Decl{}, fmt.Errorf("declaration %q not found", name)
	case 1:
		return hits[0], nil
	default:
		names := make([]string, len(hits))
		for i, h := range hits {
			names[i] = h.Addr
		}
		return Decl{}, fmt.Errorf("ambiguous name %q: %s", name, strings.Join(names, ", "))
	}
}

// Full returns the declaration span including its doc comment.
func (d Decl) Full() Span {
	start := d.Decl.Start
	if d.Doc.End > d.Doc.Start && d.Doc.Start < start {
		start = d.Doc.Start
	}
	return Span{Start: start, End: d.Decl.End}
}
