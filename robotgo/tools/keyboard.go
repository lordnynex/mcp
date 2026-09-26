package tools

import (
	"context"

	"github.com/lordnynex/mcp/robotgo/desktop"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type typeInput struct {
	Text    string `json:"text" jsonschema:"UTF-8 string to type"`
	DelayMs int    `json:"delayMs,omitempty" jsonschema:"Milliseconds between runes"`
	PID     int    `json:"pid,omitempty" jsonschema:"Optional target process id"`
}

type keyTapInput struct {
	Key       string   `json:"key" jsonschema:"robotgo key name or a single letter/digit"`
	Modifiers []string `json:"modifiers,omitempty" jsonschema:"cmd, alt, ctrl, shift (and left/right variants)"`
	PID       int      `json:"pid,omitempty" jsonschema:"Optional target process id"`
}

type keyToggleInput struct {
	Key       string   `json:"key" jsonschema:"robotgo key name or a single letter/digit"`
	State     string   `json:"state,omitempty" jsonschema:"down or up (default down)"`
	Modifiers []string `json:"modifiers,omitempty" jsonschema:"Optional modifiers"`
	PID       int      `json:"pid,omitempty" jsonschema:"Optional target process id"`
}

type keyPressInput struct {
	Key       string   `json:"key" jsonschema:"robotgo key name or a single letter/digit"`
	Modifiers []string `json:"modifiers,omitempty" jsonschema:"Optional modifiers"`
	PID       int      `json:"pid,omitempty" jsonschema:"Optional target process id"`
}

type okOutput struct {
	OK bool `json:"ok"`
}

func registerKeyboard(s *mcp.Server, d desktop.Driver) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "Type",
		Title:       "Type",
		Description: "Type a UTF-8 string via Unicode events. Prefer Type for text and KeyTap for shortcuts. Cmd/Ctrl is available from ObserveDesktop.cmdCtrl.",
		Annotations: annotations("Type"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, in typeInput) (*mcp.CallToolResult, *okOutput, error) {
		if err := doType(d, in); err != nil {
			return nil, nil, err
		}
		return nil, &okOutput{OK: true}, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "KeyTap",
		Title:       humanTitle("KeyTap"),
		Description: "Tap a named key with optional modifiers. Use ObserveDesktop.cmdCtrl for the platform accelerator (cmd on macOS, ctrl elsewhere). Do not invent key names.",
		Annotations: annotations("KeyTap"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, in keyTapInput) (*mcp.CallToolResult, *okOutput, error) {
		if err := doKeyTap(d, in); err != nil {
			return nil, nil, err
		}
		return nil, &okOutput{OK: true}, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "KeyToggle",
		Title:       humanTitle("KeyToggle"),
		Description: "Hold or release a key (down/up) with optional modifiers.",
		Annotations: annotations("KeyToggle"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, in keyToggleInput) (*mcp.CallToolResult, *okOutput, error) {
		if err := doKeyToggle(d, in); err != nil {
			return nil, nil, err
		}
		return nil, &okOutput{OK: true}, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "KeyPress",
		Title:       humanTitle("KeyPress"),
		Description: "Press a key down, wait a short delay, then release (robotgo KeyPress).",
		Annotations: annotations("KeyPress"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, in keyPressInput) (*mcp.CallToolResult, *okOutput, error) {
		if err := doKeyPress(d, in); err != nil {
			return nil, nil, err
		}
		return nil, &okOutput{OK: true}, nil
	})
}
