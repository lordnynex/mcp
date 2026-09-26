package tools

import (
	"context"
	"strings"
	"unicode"

	"github.com/lordnynex/mcp/robotgo/desktop"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ToolNames is the complete registered tool list.
var ToolNames = []string{
	"ObserveDesktop",
	"MouseMove", "MouseClick", "MouseToggle", "MouseScroll", "MouseDrag", "MouseSelect", "MousePath", "MouseLocation",
	"Type", "KeyTap", "KeyToggle", "KeyPress",
	"CaptureScreen", "GetScreenInfo", "GetPixelColor",
	"ClipboardRead", "ClipboardWrite", "ClipboardPaste",
	"GetWindowTitle", "FocusWindow", "GetWindowBounds", "MinWindow", "MaxWindow", "CloseWindow",
	"ListProcesses", "FindProcess", "KillProcess",
	"SetDelay", "Sleep",
	"DesktopOps",
}

// Register mounts every robotgo tool on s.
func Register(s *mcp.Server, d desktop.Driver) {
	if d == nil {
		d = desktop.Current()
	}
	registerObserve(s, d)
	registerMouse(s, d)
	registerPath(s, d)
	registerKeyboard(s, d)
	registerScreen(s, d)
	registerClipboard(s, d)
	registerWindow(s, d)
	registerProcess(s, d)
	registerTiming(s, d)
	registerBatch(s, d)
}

func annotations(name string) *mcp.ToolAnnotations {
	open := true
	a := &mcp.ToolAnnotations{
		Title:         humanTitle(name),
		OpenWorldHint: &open,
	}
	switch name {
	case "ObserveDesktop", "MouseLocation", "CaptureScreen", "GetScreenInfo",
		"GetPixelColor", "ClipboardRead", "GetWindowTitle", "GetWindowBounds",
		"ListProcesses", "FindProcess":
		a.ReadOnlyHint = true
	case "KillProcess", "CloseWindow":
		d := true
		a.DestructiveHint = &d
	case "SetDelay":
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

func ptrInt(p *int) (int, bool) {
	if p == nil {
		return 0, false
	}
	return *p, true
}
