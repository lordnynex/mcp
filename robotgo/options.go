package robotgo

import (
	"sync"
	"time"

	"github.com/lordnynex/mcp/robotgo/desktop"
)

const (
	defaultPageSize                  = 50
	defaultKeepAlive                 = 60 * time.Second
	defaultKeepAliveFailureThreshold = 2
)

// MCPDefaults are tunables for ServerOptions. Zero values fall back to defaults.
type MCPDefaults struct {
	PageSize                  int
	KeepAlive                 time.Duration
	KeepAliveFailureThreshold int
}

// DesktopDefaults are tunables applied to captures and robotgo delays.
type DesktopDefaults struct {
	KeySleep           int
	MouseSleep         int
	ScreenshotMaxWidth int
	ScreenshotFormat   string
}

var (
	mcpMu       sync.RWMutex
	mcpDefaults MCPDefaults
)

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

// SetDesktopDefaults sets screenshot encoding and delay defaults used by tools.
func SetDesktopDefaults(d DesktopDefaults) {
	maxW := d.ScreenshotMaxWidth
	if maxW <= 0 {
		maxW = desktop.DefaultMaxWidth
	}
	format := d.ScreenshotFormat
	if format == "" {
		format = desktop.DefaultFormat
	}
	desktop.SetEncodeOptions(desktop.EncodeOptions{MaxWidth: maxW, Format: format})
	desktop.Current().SetDelay(d.KeySleep, d.MouseSleep)
}
