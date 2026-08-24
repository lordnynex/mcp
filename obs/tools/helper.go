package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"unicode"

	"github.com/andreykaipov/goobs"
	"github.com/lordnynex/mcp/obs/session"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type requestCall func(c *goobs.Client, raw json.RawMessage) (any, error)

var (
	dispatchMu sync.RWMutex
	dispatch   = map[string]requestCall{}
)

func registerDispatcher(name string, call requestCall) {
	dispatchMu.Lock()
	defer dispatchMu.Unlock()
	dispatch[name] = call
}

func lookupDispatcher(name string) (requestCall, bool) {
	dispatchMu.RLock()
	defer dispatchMu.RUnlock()
	fn, ok := dispatch[name]
	return fn, ok
}

func addRequest[In, Out any](s *mcp.Server, h *session.Host, name, desc string, call func(*goobs.Client, *In) (*Out, error)) {
	registerDispatcher(name, func(c *goobs.Client, raw json.RawMessage) (any, error) {
		var in In
		if err := unmarshalParams(raw, &in); err != nil {
			return nil, err
		}
		return call(c, &in)
	})
	mcp.AddTool(s, &mcp.Tool{
		Name:        name,
		Description: desc,
		Title:       humanTitle(name),
		Annotations: annotationsFor(name),
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in In) (*mcp.CallToolResult, *Out, error) {
		c, err := h.RequireClient()
		if err != nil {
			return nil, nil, err
		}
		out, err := call(c, &in)
		return nil, out, err
	})
}

func unmarshalParams(raw json.RawMessage, dest any) error {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}
	if err := json.Unmarshal(raw, dest); err != nil {
		return fmt.Errorf("decode requestData: %w", err)
	}
	return nil
}

func annotationsFor(name string) *mcp.ToolAnnotations {
	open := true
	a := &mcp.ToolAnnotations{
		Title:         humanTitle(name),
		OpenWorldHint: &open,
	}
	switch {
	case strings.HasPrefix(name, "Get"):
		a.ReadOnlyHint = true
	case strings.HasPrefix(name, "Remove"),
		strings.HasPrefix(name, "Stop"),
		name == "StartStream",
		name == "StartRecord",
		name == "SetCurrentProfile",
		name == "SetCurrentSceneCollection":
		d := true
		a.DestructiveHint = &d
	case strings.HasPrefix(name, "Set"):
		a.IdempotentHint = true
	}
	return a
}

func humanTitle(name string) string {
	var b strings.Builder
	for i, r := range name {
		if i > 0 && unicode.IsUpper(r) {
			b.WriteByte(' ')
			b.WriteRune(unicode.ToLower(r))
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func notifyProgress(ctx context.Context, req *mcp.CallToolRequest, progress, total float64, message string) {
	if req == nil || req.Params == nil || req.Session == nil {
		return
	}
	token := req.Params.GetProgressToken()
	if token == nil {
		return
	}
	_ = req.Session.NotifyProgress(ctx, &mcp.ProgressNotificationParams{
		ProgressToken: token,
		Progress:      progress,
		Total:         total,
		Message:       message,
	})
}

// RegisterAlways registers Connect and ConnectionStatus.
func RegisterAlways(s *mcp.Server, h *session.Host) {
	registerAlwaysTools(s, h)
}

// RegisterConnected registers protocol tools available after Identify.
func RegisterConnected(s *mcp.Server, h *session.Host) {
	registerDisconnectTool(s, h)
	registerBatchTool(s, h)
	registerEventTools(s, h)
	registerGeneralTools(s, h)
	registerConfigTools(s, h)
	registerSourcesTools(s, h)
	registerCanvasesTools(s, h)
	registerScenesTools(s, h)
	registerInputsTools(s, h)
	registerTransitionsTools(s, h)
	registerFiltersTools(s, h)
	registerSceneItemsTools(s, h)
	registerOutputsTools(s, h)
	registerStreamTools(s, h)
	registerRecordTools(s, h)
	registerMediaInputsTools(s, h)
	registerUiTools(s, h)
}
