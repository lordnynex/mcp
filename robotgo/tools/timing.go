package tools

import (
	"context"

	"github.com/lordnynex/mcp/robotgo/desktop"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type delayInput struct {
	KeySleep   *int `json:"keySleep,omitempty" jsonschema:"Milliseconds after key events (robotgo KeySleep)"`
	MouseSleep *int `json:"mouseSleep,omitempty" jsonschema:"Milliseconds after mouse events (robotgo MouseSleep)"`
}

type delayOutput struct {
	KeySleep   int `json:"keySleep"`
	MouseSleep int `json:"mouseSleep"`
}

type sleepInput struct {
	Ms int `json:"ms" jsonschema:"Milliseconds to sleep"`
}

func registerTiming(s *mcp.Server, d desktop.Driver) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "SetDelay",
		Title:       humanTitle("SetDelay"),
		Description: "Set robotgo KeySleep and/or MouseSleep in milliseconds. Omitted fields keep their current values when both cannot be read; pass the pair you want.",
		Annotations: annotations("SetDelay"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, in delayInput) (*mcp.CallToolResult, *delayOutput, error) {
		out := doSetDelay(d, in)
		return nil, &out, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "Sleep",
		Title:       "Sleep",
		Description: "Sleep the server for ms milliseconds. Prefer this over spinning; do not sleep inside every click.",
		Annotations: annotations("Sleep"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, in sleepInput) (*mcp.CallToolResult, *okOutput, error) {
		doSleep(d, in)
		return nil, &okOutput{OK: true}, nil
	})
}
