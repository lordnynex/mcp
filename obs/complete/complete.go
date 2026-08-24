package complete

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/lordnynex/mcp/obs/protocol"
	"github.com/lordnynex/mcp/obs/session"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Handle answers completion/complete for prompt arguments and event resource URIs.
func Handle(ctx context.Context, req *mcp.CompleteRequest) (*mcp.CompleteResult, error) {
	return HandleHost(ctx, session.Default(), req)
}

// HandleHost is Handle with an explicit Host (for tests).
func HandleHost(_ context.Context, h *session.Host, req *mcp.CompleteRequest) (*mcp.CompleteResult, error) {
	if req == nil || req.Params == nil || req.Params.Ref == nil {
		return nil, fmt.Errorf("missing completion reference")
	}
	prefix := ""
	if req.Params.Argument.Value != "" {
		prefix = req.Params.Argument.Value
	}
	switch req.Params.Ref.Type {
	case "ref/prompt":
		return completePrompt(h, req.Params.Ref.Name, req.Params.Argument.Name, prefix)
	case "ref/resource":
		return completeResource(req.Params.Ref.URI, prefix)
	default:
		return nil, fmt.Errorf("unrecognized reference type %q", req.Params.Ref.Type)
	}
}

func completePrompt(h *session.Host, prompt, arg, prefix string) (*mcp.CompleteResult, error) {
	var values []string
	switch {
	case sceneNamePrompt(prompt) && arg == "sceneName":
		values = sceneNames(h, prefix)
	case prompt == "obs-subscribe-events" && arg == "category":
		values = prefixMatch(slices.Sorted(maps.Keys(protocol.EventCategories)), prefix)
	case prompt == "obs-connect" && arg == "host":
		values = prefixMatch([]string{session.DefaultHost}, prefix)
	default:
		values = nil
	}
	return completeResult(values), nil
}

func completeResource(uri, prefix string) (*mcp.CompleteResult, error) {
	if !strings.HasPrefix(uri, protocol.EventURIPrefix) && uri != protocol.EventURITemplate && uri != protocol.EventsResourceURI {
		return completeResult(nil), nil
	}
	return completeResult(prefixMatch(protocol.EventNames, prefix)), nil
}

func sceneNamePrompt(name string) bool {
	switch name {
	case "obs-switch-scene", "obs-record-clip", "obs-create-browser-source", "obs-create-input":
		return true
	default:
		return false
	}
}

func sceneNames(h *session.Host, prefix string) []string {
	if h == nil {
		return nil
	}
	c, err := h.RequireClient()
	if err != nil {
		return nil
	}
	list, err := c.Scenes.GetSceneList(&scenes.GetSceneListParams{})
	if err != nil || list == nil {
		return nil
	}
	names := make([]string, 0, len(list.Scenes))
	for _, sc := range list.Scenes {
		if sc != nil && sc.SceneName != "" {
			names = append(names, sc.SceneName)
		}
	}
	slices.Sort(names)
	return prefixMatch(names, prefix)
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
