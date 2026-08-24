package obs

import (
	"context"
	"crypto/rand"
	"log/slog"

	"github.com/lordnynex/mcp/obs/complete"
	"github.com/lordnynex/mcp/obs/prompts"
	"github.com/lordnynex/mcp/obs/resources"
	"github.com/lordnynex/mcp/obs/session"
	"github.com/lordnynex/mcp/obs/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	Name    = "obs"
	Version = "0.1.0"
	Title   = "OBS Studio"
)

const instructions = `This server controls OBS Studio over obs-websocket. Call Connect before any OBS request tools appear (notifications/tools/list_changed). Use SubscribeEvents then resources/subscribe or subscriptions/listen on the returned obs://events/{eventType} URIs. Prompts: obs-connect, obs-switch-scene, obs-studio-status, obs-subscribe-events, obs-start-stream, obs-start-record, obs-record-clip, obs-create-browser-source, obs-create-input, obs-request-batch. StartStream asks for confirmation via elicitation. Recording tools never elicit. Sleep is only valid inside RequestBatch with SerialRealtime.`

var schemaCache = mcp.NewSchemaCache()

// New returns a standalone OBS MCP server with session tools registered.
func New() *mcp.Server {
	opts := currentMCP()
	s := mcp.NewServer(&mcp.Implementation{
		Name:        Name,
		Title:       Title,
		Description: "MCP server for OBS Studio via obs-websocket",
		Version:     Version,
		WebsiteURL:  "https://github.com/lordnynex/mcp/blob/main/obs/README.md",
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

// Register mounts OBS session tools onto an existing MCP server so this package
// can run standalone or be composed into an aggregate server.
func Register(s *mcp.Server) {
	h := session.Default()
	h.OnConnected = func(srv *mcp.Server) {
		tools.RegisterConnected(srv, h)
		resources.Register(srv, h)
		prompts.RegisterConnected(srv, h)
	}
	h.OnDisconnected = func(srv *mcp.Server) {
		prompts.UnregisterConnected(srv)
	}
	h.AddServer(s)
	tools.RegisterAlways(s, h)
	prompts.RegisterAlways(s, h)
}

// ResourceSubscribe is the MCP resources/subscribe handler. Aggregate servers
// that call Register should pass this as ServerOptions.SubscribeHandler.
func ResourceSubscribe(ctx context.Context, req *mcp.SubscribeRequest) error {
	return resources.Subscribe(ctx, req)
}

// ResourceUnsubscribe is the MCP resources/unsubscribe handler.
func ResourceUnsubscribe(ctx context.Context, req *mcp.UnsubscribeRequest) error {
	return resources.Unsubscribe(ctx, req)
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
			r.TTLMs = 5_000
		case *mcp.ListPromptsResult:
			r.TTLMs = 30_000
		}
		return res, err
	}
}
