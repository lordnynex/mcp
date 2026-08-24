package session

import (
	"sync"
	"time"
)

const DefaultHost = "localhost:4455"

// ConnectDefaults are used when Connect tool arguments are omitted.
type ConnectDefaults struct {
	Host            string
	Password        string
	ResponseTimeout time.Duration
}

var (
	defaultMu       sync.RWMutex
	connectDefaults ConnectDefaults
)

// SetDefaults sets host/password/timeout used when Connect arguments are empty.
func SetDefaults(d ConnectDefaults) {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	connectDefaults = d
}

func currentDefaults() ConnectDefaults {
	defaultMu.RLock()
	defer defaultMu.RUnlock()
	return connectDefaults
}
