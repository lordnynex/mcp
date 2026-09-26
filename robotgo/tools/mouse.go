package tools

import (
	"context"

	"github.com/lordnynex/mcp/robotgo/desktop"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type mouseMoveInput struct {
	X         int     `json:"x" jsonschema:"X coordinate (absolute) or delta (relative)"`
	Y         int     `json:"y" jsonschema:"Y coordinate (absolute) or delta (relative)"`
	Relative  bool    `json:"relative,omitempty" jsonschema:"Move relative to the current cursor"`
	Smooth    bool    `json:"smooth,omitempty" jsonschema:"Human-like MoveSmooth"`
	DisplayID *int    `json:"displayId,omitempty" jsonschema:"Target display for absolute moves"`
	Low       float64 `json:"low,omitempty" jsonschema:"Smooth speed low (robotgo MoveSmooth)"`
	High      float64 `json:"high,omitempty" jsonschema:"Smooth speed high"`
	Delay     int     `json:"delay,omitempty" jsonschema:"Smooth mouse delay in ms"`
}

type mouseClickInput struct {
	Button string `json:"button,omitempty" jsonschema:"left, center, right, wheelUp, wheelDown, wheelLeft, wheelRight"`
	Double bool   `json:"double,omitempty" jsonschema:"Double-click"`
	Count  int    `json:"count,omitempty" jsonschema:"Click count (default 1)"`
	X      *int   `json:"x,omitempty" jsonschema:"Optional move-to X before clicking"`
	Y      *int   `json:"y,omitempty" jsonschema:"Optional move-to Y before clicking"`
}

type mouseToggleInput struct {
	Button string `json:"button,omitempty" jsonschema:"left, center, right, or a wheel button"`
	State  string `json:"state,omitempty" jsonschema:"down or up (default down)"`
}

type mouseScrollInput struct {
	X      int    `json:"x,omitempty" jsonschema:"Horizontal scroll ticks"`
	Y      int    `json:"y,omitempty" jsonschema:"Vertical scroll ticks"`
	Dir    string `json:"dir,omitempty" jsonschema:"up, down, left, or right (uses ScrollDir)"`
	Smooth bool   `json:"smooth,omitempty" jsonschema:"Use ScrollSmooth"`
	Steps  int    `json:"steps,omitempty" jsonschema:"Smooth step count"`
	Sleep  int    `json:"sleepMs,omitempty" jsonschema:"Smooth sleep between steps"`
}

type mouseDragInput struct {
	X int `json:"x" jsonschema:"Drag destination X"`
	Y int `json:"y" jsonschema:"Drag destination Y"`
}

type mouseSelectInput struct {
	X         int     `json:"x" jsonschema:"Selection start X (same space as ObserveDesktop)"`
	Y         int     `json:"y" jsonschema:"Selection start Y"`
	ToX       int     `json:"toX" jsonschema:"Selection end X"`
	ToY       int     `json:"toY" jsonschema:"Selection end Y"`
	Smooth    bool    `json:"smooth,omitempty" jsonschema:"Human-like MoveSmooth for both ends"`
	DisplayID *int    `json:"displayId,omitempty" jsonschema:"Target display"`
	Low       float64 `json:"low,omitempty" jsonschema:"Smooth speed low"`
	High      float64 `json:"high,omitempty" jsonschema:"Smooth speed high"`
	Delay     int     `json:"delay,omitempty" jsonschema:"Smooth mouse delay in ms"`
}

type pointOutput struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func registerMouse(s *mcp.Server, d desktop.Driver) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "MouseMove",
		Title:       humanTitle("MouseMove"),
		Description: "Move the cursor. Absolute by default (same pixel space as ObserveDesktop). Set relative for a delta. Set smooth for robotgo MoveSmooth.",
		Annotations: annotations("MouseMove"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, in mouseMoveInput) (*mcp.CallToolResult, *pointOutput, error) {
		out, err := doMove(d, in)
		if err != nil {
			return nil, nil, err
		}
		return nil, &out, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "MouseClick",
		Title:       humanTitle("MouseClick"),
		Description: "Click a mouse button. Optionally move to x,y first — use the control's geometric center, not the first glyph or an edge. Buttons: left, center, right, wheelUp, wheelDown, wheelLeft, wheelRight.",
		Annotations: annotations("MouseClick"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, in mouseClickInput) (*mcp.CallToolResult, *pointOutput, error) {
		out, err := doClick(d, in)
		if err != nil {
			return nil, nil, err
		}
		return nil, &out, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "MouseToggle",
		Title:       humanTitle("MouseToggle"),
		Description: "Press or release a mouse button (down/up). Use MouseDrag for a click-and-drag.",
		Annotations: annotations("MouseToggle"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, in mouseToggleInput) (*mcp.CallToolResult, *struct{}, error) {
		if err := doToggle(d, in); err != nil {
			return nil, nil, err
		}
		return nil, &struct{}{}, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "MouseScroll",
		Title:       humanTitle("MouseScroll"),
		Description: "Scroll the mouse wheel. Use x/y ticks or dir (up/down/left/right). Set smooth for ScrollSmooth.",
		Annotations: annotations("MouseScroll"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, in mouseScrollInput) (*mcp.CallToolResult, *struct{}, error) {
		if err := doScroll(d, in); err != nil {
			return nil, nil, err
		}
		return nil, &struct{}{}, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "MouseDrag",
		Title:       humanTitle("MouseDrag"),
		Description: "Drag the left button smoothly to x,y (robotgo DragSmooth). Cursor must already be at the start; prefer MouseSelect for from-to text selection.",
		Annotations: annotations("MouseDrag"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, in mouseDragInput) (*mcp.CallToolResult, *pointOutput, error) {
		out, err := doDrag(d, in)
		if err != nil {
			return nil, nil, err
		}
		return nil, &out, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "MouseSelect",
		Title:       humanTitle("MouseSelect"),
		Description: "Select by dragging from x,y to toX,toY (left button down, move, up). Same pixel space as ObserveDesktop. Prefer this over a click when highlighting text.",
		Annotations: annotations("MouseSelect"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, in mouseSelectInput) (*mcp.CallToolResult, *pointOutput, error) {
		out, err := doSelect(d, in)
		if err != nil {
			return nil, nil, err
		}
		return nil, &out, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "MouseLocation",
		Title:       humanTitle("MouseLocation"),
		Description: "Return the current cursor position in the same coordinate space as ObserveDesktop.",
		Annotations: annotations("MouseLocation"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, *pointOutput, error) {
		out := locationOf(d)
		return nil, &out, nil
	})
}
