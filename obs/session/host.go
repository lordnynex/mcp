package session

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"time"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api"
	"github.com/andreykaipov/goobs/api/events/subscriptions"
	"github.com/lordnynex/mcp/obs/protocol"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ErrNotConnected is returned when a tool runs before Connect succeeds.
var ErrNotConnected = errors.New("not connected to OBS; call Connect first")

// ConnectInput is the Connect tool argument set.
type ConnectInput struct {
	Host               string `json:"host,omitempty" jsonschema:"OBS websocket host:port (default localhost:4455 or CLI --host)"`
	Password           string `json:"password,omitempty" jsonschema:"OBS websocket password"`
	EventSubscriptions *int   `json:"eventSubscriptions,omitempty" jsonschema:"EventSubscription bitmask; default All non-high-volume events"`
}

// ConnectResult is returned by a successful Connect.
type ConnectResult struct {
	Connected           bool    `json:"connected"`
	Host                string  `json:"host"`
	ObsVersion          string  `json:"obsVersion,omitempty"`
	ObsWebSocketVersion string  `json:"obsWebSocketVersion,omitempty"`
	RpcVersion          float64 `json:"rpcVersion,omitempty"`
}

// Status is the ConnectionStatus payload.
type Status struct {
	Connected  bool   `json:"connected"`
	Host       string `json:"host,omitempty"`
	ObsVersion string `json:"obsVersion,omitempty"`
}

// Event is one OBS websocket event stored for resource reads.
type Event struct {
	EventType string `json:"eventType"`
	EventData any    `json:"eventData"`
}

// Host owns the process-wide OBS client and MCP feature gating.
type Host struct {
	mu sync.Mutex

	client     *goobs.Client
	host       string
	eventSubs  int
	obsVersion string

	servers []*mcp.Server
	toolsOn map[*mcp.Server]bool

	listenCancel context.CancelFunc

	eventMu  sync.Mutex
	latest   map[string]Event
	buffer   []Event
	interest map[*mcp.ServerSession]map[string]struct{}

	// OnConnected registers protocol tools/resources/prompts for one server.
	OnConnected func(s *mcp.Server)
	// OnDisconnected removes connection-only prompts (tools/resources are removed by Host).
	OnDisconnected func(s *mcp.Server)
}

var defaultHost = New()

// Default returns the process-wide Host.
func Default() *Host {
	return defaultHost
}

// New constructs an empty Host.
func New() *Host {
	return &Host{
		toolsOn:  make(map[*mcp.Server]bool),
		latest:   make(map[string]Event),
		interest: make(map[*mcp.ServerSession]map[string]struct{}),
	}
}

// AddServer records an MCP server that should receive protocol tools after Connect.
func (h *Host) AddServer(s *mcp.Server) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if slices.Contains(h.servers, s) {
		return
	}
	h.servers = append(h.servers, s)
}

// RequireClient returns the connected goobs client or ErrNotConnected.
func (h *Host) RequireClient() (*goobs.Client, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.client == nil {
		return nil, ErrNotConnected
	}
	return h.client, nil
}

// Status reports whether OBS is connected.
func (h *Host) Status() Status {
	h.mu.Lock()
	defer h.mu.Unlock()
	return Status{
		Connected:  h.client != nil,
		Host:       h.host,
		ObsVersion: h.obsVersion,
	}
}

// Connect identifies with OBS and registers protocol features.
func (h *Host) Connect(in ConnectInput) (ConnectResult, error) {
	d := currentDefaults()
	host := firstNonEmpty(in.Host, d.Host, DefaultHost)
	password := in.Password
	if password == "" {
		password = d.Password
	}
	subs := subscriptions.All
	if in.EventSubscriptions != nil {
		subs = *in.EventSubscriptions
	}
	timeout := d.ResponseTimeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	h.mu.Lock()
	wasConnected := h.client != nil
	if wasConnected {
		h.teardownClientLocked()
	}
	servers := slices.Clone(h.servers)
	h.mu.Unlock()

	opts := []goobs.Option{
		goobs.WithPassword(password),
		goobs.WithEventSubscriptions(subs),
		goobs.WithLogger(slogPrinter{}),
		goobs.WithResponseTimeoutDuration(timeout),
	}
	client, err := goobs.New(host, opts...)
	if err != nil {
		if wasConnected {
			h.unregisterProtocolFeatures(servers)
		}
		return ConnectResult{}, fmt.Errorf("connect to OBS at %s: %w", host, err)
	}

	ver, err := client.General.GetVersion()
	if err != nil {
		_ = client.Disconnect()
		if wasConnected {
			h.unregisterProtocolFeatures(servers)
		}
		return ConnectResult{}, fmt.Errorf("GetVersion after connect: %w", err)
	}

	h.mu.Lock()
	h.client = client
	h.host = host
	h.eventSubs = subs
	h.obsVersion = ver.ObsVersion
	h.mu.Unlock()

	if !wasConnected {
		h.registerProtocolFeatures(servers)
	}

	h.startListen(client)
	slog.Info("connected to OBS; protocol tools registered",
		"host", host,
		"obsVersion", ver.ObsVersion,
		"obsWebSocketVersion", ver.ObsWebSocketVersion,
	)

	return ConnectResult{
		Connected:           true,
		Host:                host,
		ObsVersion:          ver.ObsVersion,
		ObsWebSocketVersion: ver.ObsWebSocketVersion,
		RpcVersion:          ver.RpcVersion,
	}, nil
}

// Disconnect closes the OBS client and removes protocol tools.
func (h *Host) Disconnect() error {
	h.mu.Lock()
	if h.client == nil {
		h.mu.Unlock()
		return ErrNotConnected
	}
	servers := slices.Clone(h.servers)
	h.teardownClientLocked()
	h.mu.Unlock()
	h.unregisterProtocolFeatures(servers)
	h.ClearEvents()
	slog.Info("disconnected from OBS; protocol tools removed")
	return nil
}

func (h *Host) teardownClientLocked() {
	if h.listenCancel != nil {
		h.listenCancel()
		h.listenCancel = nil
	}
	if h.client != nil {
		_ = h.client.Disconnect()
		h.client = nil
	}
	h.host = ""
	h.obsVersion = ""
	h.eventSubs = 0
}

func (h *Host) startListen(client *goobs.Client) {
	h.mu.Lock()
	if h.listenCancel != nil {
		h.listenCancel()
	}
	ctx, cancel := context.WithCancel(context.Background())
	h.listenCancel = cancel
	h.mu.Unlock()

	go func() {
		client.Listen(func(event any) {
			if ctx.Err() != nil {
				return
			}
			h.OnOBSEvent(event)
		})
		select {
		case <-ctx.Done():
			return
		default:
		}
		h.mu.Lock()
		if h.client == client {
			servers := slices.Clone(h.servers)
			h.client = nil
			h.host = ""
			h.obsVersion = ""
			h.eventSubs = 0
			h.listenCancel = nil
			h.mu.Unlock()
			h.unregisterProtocolFeatures(servers)
			h.ClearEvents()
			slog.Info("OBS connection closed; protocol tools removed")
			return
		}
		h.mu.Unlock()
	}()
}

func (h *Host) registerProtocolFeatures(servers []*mcp.Server) {
	for _, s := range servers {
		h.mu.Lock()
		already := h.toolsOn[s]
		h.mu.Unlock()
		if already {
			continue
		}
		if h.OnConnected != nil {
			h.OnConnected(s)
		}
		h.mu.Lock()
		h.toolsOn[s] = true
		h.mu.Unlock()
	}
}

func (h *Host) unregisterProtocolFeatures(servers []*mcp.Server) {
	names := protocol.ConnectedToolNames()
	for _, s := range servers {
		s.RemoveTools(names...)
		s.RemoveResources(protocol.EventsResourceURI)
		s.RemoveResourceTemplates(protocol.EventURITemplate)
		if h.OnDisconnected != nil {
			h.OnDisconnected(s)
		}
		h.mu.Lock()
		delete(h.toolsOn, s)
		h.mu.Unlock()
	}
}

// IdentifiedSubs returns the EventSubscription bitmask from Identify.
func (h *Host) IdentifiedSubs() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.eventSubs
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

type slogPrinter struct{}

func (slogPrinter) Printf(format string, args ...any) {
	slog.Debug(fmt.Sprintf(format, args...))
}

var _ api.Logger = slogPrinter{}
