package robotgo

import (
	"context"
	"crypto/rand"
	"log/slog"

	"github.com/lordnynex/mcp/robotgo/complete"
	"github.com/lordnynex/mcp/robotgo/desktop"
	"github.com/lordnynex/mcp/robotgo/prompts"
	"github.com/lordnynex/mcp/robotgo/resources"
	"github.com/lordnynex/mcp/robotgo/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	Name    = "robotgo"
	Version = "0.1.0"
	Title   = "robotgo desktop"
)

const instructions = `This server drives the local desktop through robotgo (keyboard, mouse, clipboard, screenshots). Call ObserveDesktop once to plan. Map image pixels with mouseX = originX + ix * mouseWidth / width (same for y). Then one DesktopOps for the known sequence (MousePath is a valid op; use smooth false for geometric traces), then observeAfter or ObserveDesktop to verify. Aim clicks at the geometric center of the control, not the first glyph or an edge. On a miss, inset 40-80px; do not crop-loop. Select text with MouseSelect. Do not invent key names; use KeyTap modifiers and ObserveDesktop.cmdCtrl (cmd on macOS, ctrl elsewhere). Linux binaries must be built with an explicit backend tag (purego,x11 / purego,wayland / purego,libei). libei cannot screenshot. Prompts: robotgo-observe-then-act, robotgo-click-target, robotgo-type-into-field, robotgo-shortcut, robotgo-scroll-read, robotgo-copy-paste, robotgo-focus-app, robotgo-drag-drop, robotgo-desktop-ops, robotgo-backend-limits. Skills: robotgo://skills/{name}. KillProcess elicits confirmation.`

var schemaCache = mcp.NewSchemaCache()

// New returns a standalone robotgo MCP server.
func New() *mcp.Server {
	opts := currentMCP()
	s := mcp.NewServer(&mcp.Implementation{
		Name:        Name,
		Title:       Title,
		Description: "MCP server for desktop automation via go-vgo/robotgo",
		Version:     Version,
		WebsiteURL:  "https://github.com/lordnynex/mcp/blob/main/robotgo/README.md",
	}, &mcp.ServerOptions{
		Instructions: instructions,
		Logger:       slog.Default(),
		InitializedHandler: func(ctx context.Context, req *mcp.InitializedRequest) {
			logInitialized(req)
		},
		PageSize:                    opts.PageSize,
		KeepAlive:                   opts.KeepAlive,
		KeepAliveFailureThreshold:   opts.KeepAliveFailureThreshold,
		ProgressNotificationHandler: handleClientProgress,
		CompletionHandler:           complete.Handle,
		SubscribeHandler:            resources.Subscribe,
		UnsubscribeHandler:          resources.Unsubscribe,
		SchemaCache:                 schemaCache,
		GetSessionID:                rand.Text,
		Capabilities: &mcp.ServerCapabilities{
			Tools:       &mcp.ToolCapabilities{ListChanged: true},
			Resources:   &mcp.ResourceCapabilities{ListChanged: true, Subscribe: true},
			Prompts:     &mcp.PromptCapabilities{ListChanged: true},
			Completions: &mcp.CompletionCapabilities{},
		},
	})
	s.AddReceivingMiddleware(logMethods)
	s.AddSendingMiddleware(listCacheTTL)
	Register(s)
	return s
}

// Register mounts robotgo tools, prompts, and skill resources onto an existing MCP server.
func Register(s *mcp.Server) {
	tools.Register(s, desktop.Current())
	prompts.Register(s)
	resources.Register(s)
}

func logInitialized(req *mcp.InitializedRequest) {
	if req == nil || req.Session == nil {
		slog.Info("MCP client initialized")
		return
	}
	p := req.Session.InitializeParams()
	if p == nil {
		slog.Info("MCP client initialized")
		return
	}
	name, ver := "", ""
	if p.ClientInfo != nil {
		name = p.ClientInfo.Name
		ver = p.ClientInfo.Version
	}
	slog.Info("MCP client initialized",
		"client", name,
		"version", ver,
		"protocol", p.ProtocolVersion,
	)
}

func handleClientProgress(_ context.Context, req *mcp.ProgressNotificationServerRequest) {
	if req == nil || req.Params == nil {
		return
	}
	slog.Info("client progress",
		"progress", req.Params.Progress,
		"total", req.Params.Total,
		"message", req.Params.Message,
	)
}

func logMethods(next mcp.MethodHandler) mcp.MethodHandler {
	return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
		slog.Debug("MCP method", "method", method)
		return next(ctx, method, req)
	}
}

func listCacheTTL(next mcp.MethodHandler) mcp.MethodHandler {
	return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
		res, err := next(ctx, method, req)
		if err != nil {
			return res, err
		}
		switch r := res.(type) {
		case *mcp.ListToolsResult:
			r.TTLMs = 5_000
		case *mcp.ListResourcesResult:
			r.TTLMs = 30_000
		case *mcp.ListPromptsResult:
			r.TTLMs = 30_000
		}
		return res, err
	}
}
