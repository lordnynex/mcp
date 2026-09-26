package tools

import (
	"context"

	"github.com/lordnynex/mcp/robotgo/desktop"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type titleOutput struct {
	Title string `json:"title"`
}

type focusInput struct {
	Name string `json:"name,omitempty" jsonschema:"Window or app name passed to ActiveName"`
	PID  int    `json:"pid,omitempty" jsonschema:"Process id; on pure-Go backends this resolves a name then calls ActiveName"`
}

type boundsInput struct {
	PID int `json:"pid" jsonschema:"Process id"`
}

type boundsOutput struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type pidInput struct {
	PID int `json:"pid,omitempty" jsonschema:"Process id; omit to target the front window where supported"`
}

func registerWindow(s *mcp.Server, d desktop.Driver) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "GetWindowTitle",
		Title:       humanTitle("GetWindowTitle"),
		Description: "Return the focused window title. Empty on macOS -tags purego and Linux libei (ErrNotSupported backends).",
		Annotations: annotations("GetWindowTitle"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, *titleOutput, error) {
		return nil, &titleOutput{Title: d.WindowTitle()}, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "FocusWindow",
		Title:       humanTitle("FocusWindow"),
		Description: "Focus a window by name (ActiveName) or pid. Unsupported on macOS -tags purego and Linux libei.",
		Annotations: annotations("FocusWindow"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, in focusInput) (*mcp.CallToolResult, *titleOutput, error) {
		out, err := doFocusWindow(d, in)
		if err != nil {
			return nil, nil, err
		}
		return nil, &out, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "GetWindowBounds",
		Title:       humanTitle("GetWindowBounds"),
		Description: "Window bounds by pid. Not available on the pure-Go robotgo backends.",
		Annotations: annotations("GetWindowBounds"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, in boundsInput) (*mcp.CallToolResult, *boundsOutput, error) {
		out, err := doWindowBounds(d, in)
		if err != nil {
			return nil, nil, err
		}
		return nil, &out, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "MinWindow",
		Title:       humanTitle("MinWindow"),
		Description: "Minimize a window by pid. May no-op on macOS -tags purego and Linux libei.",
		Annotations: annotations("MinWindow"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, in pidInput) (*mcp.CallToolResult, *okOutput, error) {
		if err := doMinWindow(d, in); err != nil {
			return nil, nil, err
		}
		return nil, &okOutput{OK: true}, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "MaxWindow",
		Title:       humanTitle("MaxWindow"),
		Description: "Maximize a window by pid. May no-op on macOS -tags purego and Linux libei.",
		Annotations: annotations("MaxWindow"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, in pidInput) (*mcp.CallToolResult, *okOutput, error) {
		if err := doMaxWindow(d, in); err != nil {
			return nil, nil, err
		}
		return nil, &okOutput{OK: true}, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "CloseWindow",
		Title:       humanTitle("CloseWindow"),
		Description: "Close a window by pid (or the front window if pid is omitted). Destructive. May no-op on unsupported backends.",
		Annotations: annotations("CloseWindow"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, in pidInput) (*mcp.CallToolResult, *okOutput, error) {
		if err := doCloseWindow(d, in); err != nil {
			return nil, nil, err
		}
		return nil, &okOutput{OK: true}, nil
	})
}
