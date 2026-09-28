package goastedit

import jsonv2 "encoding/json/v2"

func marshal(v any) (string, error) {
	b, err := jsonv2.Marshal(v, jsonv2.Deterministic(true))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

type jsonBody struct {
	Op      string `json:"op"`
	Target  string `json:"target,omitempty"`
	Content string `json:"content,omitempty"`
	Path    string `json:"path,omitempty"`
	Node    any    `json:"node,omitempty"`
}

// Arm is one way of writing an edit. NA is set when that way cannot express it.
type Arm struct {
	Name    string
	Payload string
	NA      bool
}

// Payloads renders the model output for each arm of req.
func Payloads(src []byte, req Request) ([]Arm, error) {
	chg, err := plan(src, req)
	if err != nil {
		return nil, err
	}
	old := string(src[chg.quoteStart:chg.quoteEnd])
	search := "SEARCH\n" + old + "\nREPLACE\n" + chg.quoteNew
	minimal := string(spliceBytes(src, chg.start, chg.end, chg.repl))
	diff := unifiedDiff(string(src), minimal)

	splice, err := marshal(jsonBody{
		Op:      string(req.Op),
		Target:  req.Target,
		Content: spliceContent(req),
		Path:    req.Import,
	})
	if err != nil {
		return nil, err
	}
	node, err := parseSnippet(req.Op, req.Content, req.Import)
	if err != nil {
		return nil, err
	}
	var encoded any
	if node != nil {
		encoded, err = Encode(node)
		if err != nil {
			return nil, err
		}
	}
	astJSON, err := marshal(jsonBody{
		Op:     string(req.Op),
		Target: req.Target,
		Path:   req.Import,
		Node:   encoded,
	})
	if err != nil {
		return nil, err
	}
	intent, err := intentPayload(req)
	if err != nil {
		return nil, err
	}
	return []Arm{
		{Name: "search/replace", Payload: search},
		{Name: "unified diff", Payload: diff},
		{Name: "symbol splice", Payload: splice},
		{Name: "ast json", Payload: astJSON},
		intent,
	}, nil
}

func spliceContent(req Request) string {
	if req.Op == OpDelete || req.Op == OpAddImport {
		return ""
	}
	return req.Content
}

func intentPayload(req Request) (Arm, error) {
	arm := Arm{Name: "intent", NA: true}
	switch req.Op {
	case OpDelete:
		s, err := marshal(jsonBody{Op: string(req.Op), Target: req.Target})
		if err != nil {
			return Arm{}, err
		}
		return Arm{Name: "intent", Payload: s}, nil
	case OpAddImport:
		s, err := marshal(jsonBody{Op: string(req.Op), Path: req.Import})
		if err != nil {
			return Arm{}, err
		}
		return Arm{Name: "intent", Payload: s}, nil
	default:
		return arm, nil
	}
}

// Result is what a tool would send back after the edit, independent of which
// arm the model used to ask for it.
type Result struct {
	Summary string
	Diff    string
	File    string
	AST     string
}

// ToolResult applies req and measures the four result shapes.
func ToolResult(src []byte, req Request) (Result, error) {
	out, err := Apply(src, req)
	if err != nil {
		return Result{}, err
	}
	idx, err := Parse(out)
	if err != nil {
		return Result{}, err
	}
	enc, err := Encode(idx.AST)
	if err != nil {
		return Result{}, err
	}
	astJSON, err := marshal(enc)
	if err != nil {
		return Result{}, err
	}
	return Result{
		Summary: summary(req),
		Diff:    unifiedDiff(string(src), string(out)),
		File:    string(out),
		AST:     astJSON,
	}, nil
}

func summary(req Request) string {
	switch req.Op {
	case OpReplaceBody:
		return "replaced body of " + req.Target
	case OpReplaceDecl:
		return "replaced declaration " + req.Target
	case OpInsertBefore:
		return "inserted before " + req.Target
	case OpInsertAfter:
		return "inserted after " + req.Target
	case OpPrependStmt:
		return "prepended statement in " + req.Target
	case OpDelete:
		return "deleted " + req.Target
	case OpAddImport:
		return "added import " + req.Import
	default:
		return string(req.Op)
	}
}
