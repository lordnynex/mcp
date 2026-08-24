package session

import (
	"context"
	"fmt"
	"reflect"

	"github.com/andreykaipov/goobs/api/events"
	"github.com/lordnynex/mcp/obs/protocol"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// EnsureHighVolume checks that high-volume events were enabled at Identify.
func (h *Host) EnsureHighVolume(types []string) error {
	subs := h.IdentifiedSubs()
	for _, n := range types {
		bit, ok := protocol.HighVolumeEvents[n]
		if !ok {
			continue
		}
		if subs&bit == 0 {
			return fmt.Errorf("event %s is high-volume and was not enabled at Connect; reconnect with eventSubscriptions including bit %d", n, bit)
		}
	}
	return nil
}

// SetInterest adds or removes event types for a session.
func (h *Host) SetInterest(sess *mcp.ServerSession, types []string, add bool) {
	h.eventMu.Lock()
	defer h.eventMu.Unlock()
	cur := h.interest[sess]
	if add {
		if cur == nil {
			cur = make(map[string]struct{})
			h.interest[sess] = cur
		}
		for _, n := range types {
			cur[n] = struct{}{}
		}
		return
	}
	if cur == nil {
		return
	}
	for _, n := range types {
		delete(cur, n)
	}
	if len(cur) == 0 {
		delete(h.interest, sess)
	}
}

// ClearInterest drops all event types for a session and returns the previous set.
func (h *Host) ClearInterest(sess *mcp.ServerSession) []string {
	h.eventMu.Lock()
	defer h.eventMu.Unlock()
	cur := h.interest[sess]
	out := make([]string, 0, len(cur))
	for n := range cur {
		out = append(out, n)
	}
	delete(h.interest, sess)
	return out
}

// AnyInterest reports whether any session selected eventType.
func (h *Host) AnyInterest(eventType string) bool {
	h.eventMu.Lock()
	defer h.eventMu.Unlock()
	for _, set := range h.interest {
		if _, ok := set[eventType]; ok {
			return true
		}
	}
	return false
}

func (h *Host) recordEvent(ev Event) {
	h.eventMu.Lock()
	defer h.eventMu.Unlock()
	h.latest[ev.EventType] = ev
	h.buffer = append(h.buffer, ev)
	if len(h.buffer) > protocol.EventBufferSize {
		h.buffer = h.buffer[len(h.buffer)-protocol.EventBufferSize:]
	}
}

// ClearEvents resets the ring buffer and interest map.
func (h *Host) ClearEvents() {
	h.eventMu.Lock()
	defer h.eventMu.Unlock()
	h.latest = make(map[string]Event)
	h.buffer = nil
	h.interest = make(map[*mcp.ServerSession]map[string]struct{})
}

// EventBuffer returns a copy of recent events.
func (h *Host) EventBuffer() []Event {
	h.eventMu.Lock()
	defer h.eventMu.Unlock()
	out := make([]Event, len(h.buffer))
	copy(out, h.buffer)
	return out
}

// LatestEvent returns the most recent event of the given type.
func (h *Host) LatestEvent(eventType string) (Event, bool) {
	h.eventMu.Lock()
	defer h.eventMu.Unlock()
	ev, ok := h.latest[eventType]
	return ev, ok
}

// OnOBSEvent records an OBS event and notifies subscribed MCP sessions.
func (h *Host) OnOBSEvent(event any) {
	if event == nil {
		return
	}
	name := eventTypeName(event)
	if !protocol.IsKnownEvent(name) {
		return
	}
	ev := Event{EventType: name, EventData: event}
	h.recordEvent(ev)
	if _, ok := event.(*events.ExitStarted); ok {
		// Connection teardown is handled when Listen returns.
	}
	if !h.AnyInterest(name) {
		return
	}
	uri := protocol.EventResourceURI(name)
	params := &mcp.ResourceUpdatedNotificationParams{URI: uri}
	params.Meta = mcp.Meta{"obsEvent": ev}

	h.mu.Lock()
	servers := append([]*mcp.Server{}, h.servers...)
	h.mu.Unlock()
	for _, s := range servers {
		_ = s.ResourceUpdated(context.Background(), params)
	}
}

func eventTypeName(event any) string {
	t := reflect.TypeOf(event)
	if t == nil {
		return ""
	}
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t.Name()
}
