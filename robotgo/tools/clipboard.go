package tools

import (
	"context"

	"github.com/lordnynex/mcp/robotgo/desktop"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type clipboardText struct {
	Text string `json:"text" jsonschema:"Clipboard text"`
}

func registerClipboard(s *mcp.Server, d desktop.Driver) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "ClipboardRead",
		Title:       humanTitle("ClipboardRead"),
		Description: "Read the system clipboard text. Linux needs xclip or xsel on PATH.",
		Annotations: annotations("ClipboardRead"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, *clipboardText, error) {
		text, err := d.ClipboardRead()
		if err != nil {
			return nil, nil, err
		}
		return nil, &clipboardText{Text: text}, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "ClipboardWrite",
		Title:       humanTitle("ClipboardWrite"),
		Description: "Write text to the system clipboard without pasting.",
		Annotations: annotations("ClipboardWrite"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, in clipboardText) (*mcp.CallToolResult, *okOutput, error) {
		if err := doClipboardWrite(d, in); err != nil {
			return nil, nil, err
		}
		return nil, &okOutput{OK: true}, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "ClipboardPaste",
		Title:       humanTitle("ClipboardPaste"),
		Description: "Write text to the clipboard and send the paste shortcut (Cmd+V or Ctrl+V).",
		Annotations: annotations("ClipboardPaste"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, in clipboardText) (*mcp.CallToolResult, *okOutput, error) {
		if err := doClipboardPaste(d, in); err != nil {
			return nil, nil, err
		}
		return nil, &okOutput{OK: true}, nil
	})
}
