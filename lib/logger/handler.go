package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"strings"
	"sync"
	"time"
)

const timeLayout = "15:04:05.000"

type handler struct {
	opts   Options
	color  bool
	attrs  []slog.Attr
	groups []string
	mu     *sync.Mutex
}

func newHandler(opts Options) *handler {
	opts.Writer = writerOrDefault(opts.Writer)
	return &handler{
		opts:  opts,
		color: colorEnabled(opts),
		mu:    &sync.Mutex{},
	}
}

func (h *handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level
}

func (h *handler) Handle(_ context.Context, r slog.Record) error {
	var b strings.Builder

	ts := r.Time
	if ts.IsZero() {
		ts = time.Now()
	}
	writeColored(&b, h.color, ansiGray, ts.Format(timeLayout))
	b.WriteByte(' ')
	writeColored(&b, h.color, levelColor(r.Level), formatLevel(r.Level))
	b.WriteByte(' ')
	b.WriteString(r.Message)

	write := func(a slog.Attr) {
		writeAttr(&b, h.color, h.groups, a)
	}
	for _, a := range h.attrs {
		write(a)
	}
	r.Attrs(func(a slog.Attr) bool {
		write(a)
		return true
	})
	b.WriteByte('\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := io.WriteString(h.opts.Writer, b.String())
	return err
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	c := *h
	c.attrs = append(slices.Clip(h.attrs), attrs...)
	return &c
}

func (h *handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	c := *h
	c.groups = append(slices.Clip(h.groups), name)
	return &c
}

func writeAttr(b *strings.Builder, color bool, groups []string, a slog.Attr) {
	a.Value = a.Value.Resolve()
	if a.Equal(slog.Attr{}) {
		return
	}
	if a.Value.Kind() == slog.KindGroup {
		name := a.Key
		next := groups
		if name != "" {
			next = append(slices.Clip(groups), name)
		}
		for _, child := range a.Value.Group() {
			writeAttr(b, color, next, child)
		}
		return
	}

	key := a.Key
	if len(groups) > 0 {
		key = strings.Join(groups, ".") + "." + key
	}
	b.WriteByte(' ')
	writeColored(b, color, ansiCyan, key)
	b.WriteByte('=')
	b.WriteString(formatValue(a.Value))
}

func formatValue(v slog.Value) string {
	switch v.Kind() {
	case slog.KindString:
		return v.String()
	case slog.KindTime:
		return v.Time().Format(time.RFC3339)
	default:
		return fmt.Sprint(v.Any())
	}
}

func formatLevel(level slog.Level) string {
	switch {
	case level < slog.LevelInfo:
		return "DEBUG"
	case level < slog.LevelWarn:
		return "INFO"
	case level < slog.LevelError:
		return "WARN"
	default:
		return "ERROR"
	}
}
