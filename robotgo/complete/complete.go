package complete

import (
	"context"
	"fmt"
	"strings"

	"github.com/lordnynex/mcp/robotgo/desktop"
	"github.com/lordnynex/mcp/robotgo/skills"
	"github.com/lordnynex/mcp/robotgo/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Handle answers completion/complete for prompt arguments and skill resources.
func Handle(_ context.Context, req *mcp.CompleteRequest) (*mcp.CompleteResult, error) {
	if req == nil || req.Params == nil || req.Params.Ref == nil {
		return nil, fmt.Errorf("missing completion reference")
	}
	prefix := ""
	if req.Params.Argument.Value != "" {
		prefix = req.Params.Argument.Value
	}
	switch req.Params.Ref.Type {
	case "ref/prompt":
		return completePrompt(req.Params.Ref.Name, req.Params.Argument.Name, prefix), nil
	case "ref/resource":
		return completeResource(req.Params.Ref.URI, prefix), nil
	default:
		return nil, fmt.Errorf("unrecognized reference type %q", req.Params.Ref.Type)
	}
}

func completePrompt(prompt, arg, prefix string) *mcp.CompleteResult {
	var values []string
	switch {
	case prompt == "robotgo-shortcut" && arg == "key":
		values = prefixMatch(desktop.Keys, prefix)
	case prompt == "robotgo-shortcut" && arg == "modifiers":
		values = prefixMatch(desktop.Modifiers, prefix)
	case prompt == "robotgo-scroll-read" && arg == "dir":
		values = prefixMatch(desktop.ScrollDirs, prefix)
	case prompt == "robotgo-observe-then-act" && arg == "format":
		values = prefixMatch([]string{"png", "jpeg"}, prefix)
	case prompt == "robotgo-desktop-ops" && arg == "op":
		values = prefixMatch(tools.DesktopOpNames, prefix)
	default:
		values = []string{}
	}
	return completeResult(values)
}

func completeResource(uri, prefix string) *mcp.CompleteResult {
	if !strings.HasPrefix(uri, skills.URIPrefix) && uri != skills.URITemplate && uri != skills.IndexURI {
		return completeResult(nil)
	}
	return completeResult(prefixMatch(skills.Names, prefix))
}

func prefixMatch(values []string, prefix string) []string {
	if prefix == "" {
		return values
	}
	out := make([]string, 0, len(values))
	for _, v := range values {
		if strings.HasPrefix(strings.ToLower(v), strings.ToLower(prefix)) {
			out = append(out, v)
		}
	}
	return out
}

func completeResult(values []string) *mcp.CompleteResult {
	if values == nil {
		values = []string{}
	}
	return &mcp.CompleteResult{
		Completion: mcp.CompletionResultDetails{
			Values:  values,
			Total:   len(values),
			HasMore: false,
		},
	}
}
