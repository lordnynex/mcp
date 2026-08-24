package obs

import (
	"sync"
	"time"

	"github.com/lordnynex/mcp/obs/session"
)

const (
	defaultPageSize                  = 50
	defaultKeepAlive                 = 60 * time.Second
	defaultKeepAliveFailureThreshold = 2
)

// ConnectDefaults are used when Connect tool arguments are omitted.
type ConnectDefaults = session.ConnectDefaults

// MCPDefaults are tunables for ServerOptions. Zero values fall back to defaults.
type MCPDefaults struct {
	PageSize                  int
	KeepAlive                 time.Duration
	KeepAliveFailureThreshold int
}

var (
	mcpMu       sync.RWMutex
	mcpDefaults MCPDefaults
)

// SetConnectDefaults sets host/password/timeout used when Connect arguments are empty.
func SetConnectDefaults(d ConnectDefaults) {
	session.SetDefaults(d)
}

// SetMCPDefaults sets page size and keepalive used by New.
func SetMCPDefaults(d MCPDefaults) {
	mcpMu.Lock()
	defer mcpMu.Unlock()
	mcpDefaults = d
}

func currentMCP() MCPDefaults {
	mcpMu.RLock()
	d := mcpDefaults
	mcpMu.RUnlock()
	if d.PageSize <= 0 {
		d.PageSize = defaultPageSize
	}
	if d.KeepAlive <= 0 {
		d.KeepAlive = defaultKeepAlive
	}
	if d.KeepAliveFailureThreshold <= 0 {
		d.KeepAliveFailureThreshold = defaultKeepAliveFailureThreshold
	}
	return d
}
