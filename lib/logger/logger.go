package logger

import (
	"io"
	"log/slog"
	"os"
)

// Options configure a pretty slog logger for MCP server processes.
type Options struct {
	Level  slog.Level
	Writer io.Writer // default os.Stderr (MCP stdio uses stdout)
	Color  *bool     // nil = auto-detect TTY on Writer
}

// New returns a pretty slog logger. Output defaults to stderr.
func New(opts Options) *slog.Logger {
	return slog.New(newHandler(opts))
}

// SetDefault installs a pretty logger as the process-wide slog default.
func SetDefault(opts Options) {
	slog.SetDefault(New(opts))
}

func writerOrDefault(w io.Writer) io.Writer {
	if w == nil {
		return os.Stderr
	}
	return w
}

func colorEnabled(opts Options) bool {
	if opts.Color != nil {
		return *opts.Color
	}
	return isTerminal(writerOrDefault(opts.Writer))
}
