package tools

import (
	"context"

	"github.com/lordnynex/mcp/obs/session"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerAlwaysTools(s *mcp.Server, h *session.Host) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "Connect",
		Title:       "Connect",
		Description: "Connect to an OBS Studio websocket server. On success, OBS request tools become available via notifications/tools/list_changed.",
		Annotations: annotationsFor("Connect"),
	}, func(ctx context.Context, req *mcp.CallToolRequest, in session.ConnectInput) (*mcp.CallToolResult, session.ConnectResult, error) {
		notifyProgress(ctx, req, 0, 2, "identifying with OBS")
		out, err := h.Connect(in)
		if err != nil {
			return nil, session.ConnectResult{}, err
		}
		notifyProgress(ctx, req, 2, 2, "connected")
		return nil, out, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "ConnectionStatus",
		Title:       "Connection status",
		Description: "Report whether this MCP server is connected to OBS.",
		Annotations: annotationsFor("GetConnectionStatus"),
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, session.Status, error) {
		return nil, h.Status(), nil
	})
}

func registerDisconnectTool(s *mcp.Server, h *session.Host) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "Disconnect",
		Title:       "Disconnect",
		Description: "Disconnect from OBS and remove OBS request tools from the tool list.",
		Annotations: annotationsFor("Disconnect"),
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, session.Status, error) {
		if err := h.Disconnect(); err != nil {
			return nil, session.Status{}, err
		}
		return nil, h.Status(), nil
	})
}
