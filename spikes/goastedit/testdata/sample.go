// Package sample is a stable fixture for token measurements.
// Edit scenarios target these declarations. Keep the shape steady
// so later runs stay comparable.
package sample

import (
	"fmt"
	"strings"
)

// Box stores string values by key and counts lookups.
type Box struct {
	items map[string]string
	hits  int
	miss  int
}

// Get returns the value stored for key.
// A missing key returns an empty string and increments the miss counter.
func (b *Box) Get(key string) string {
	key = strings.TrimSpace(key)
	if b.items == nil {
		b.miss++
		return ""
	}
	val, ok := b.items[key]
	if !ok {
		b.miss++
		return ""
	}
	b.hits++
	val = strings.TrimSpace(val)
	if val == "" {
		return ""
	}
	return val
}

// Put stores val at key, allocating the map on first use.
func (b *Box) Put(key, val string) {
	key = strings.TrimSpace(key)
	if key == "" {
		return
	}
	if b.items == nil {
		b.items = make(map[string]string)
	}
	b.items[key] = strings.TrimSpace(val)
}

// Delete removes key and reports whether it was present.
func (b *Box) Delete(key string) bool {
	key = strings.TrimSpace(key)
	if b.items == nil {
		return false
	}
	_, ok := b.items[key]
	if !ok {
		return false
	}
	delete(b.items, key)
	return true
}

// Len reports how many keys are stored.
func (b *Box) Len() int {
	return len(b.items)
}

// Other is a second type so the short name Get is ambiguous.
type Other struct {
	n int
}

// Get returns the stored count.
func (o *Other) Get() int {
	return o.n
}

// Label formats the count. The receiver is a value, so the address is Other.Label.
func (o Other) Label() string {
	return fmt.Sprintf("other=%d", o.n)
}

// Helper formats a label.
func Helper(label string) string {
	label = strings.TrimSpace(label)
	if label == "" {
		return "empty"
	}
	parts := strings.Split(label, " ")
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	return fmt.Sprintf("label=%s", strings.Join(parts, "-"))
}

// This comment sits between Helper and Neighbor and must survive deletion.

// Neighbor returns a fixed label.
func Neighbor() string {
	return "stay"
}
