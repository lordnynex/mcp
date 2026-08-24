package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"golang.org/x/term"
)

const (
	ansiReset  = "\x1b[0m"
	ansiGray   = "\x1b[90m"
	ansiCyan   = "\x1b[36m"
	ansiGreen  = "\x1b[32m"
	ansiYellow = "\x1b[33m"
	ansiRed    = "\x1b[31m"
)

func levelColor(level slog.Level) string {
	switch {
	case level < slog.LevelInfo:
		return ansiGray
	case level < slog.LevelWarn:
		return ansiGreen
	case level < slog.LevelError:
		return ansiYellow
	default:
		return ansiRed
	}
}

func writeColored(b *strings.Builder, color bool, code, text string) {
	if color {
		b.WriteString(code)
		b.WriteString(text)
		b.WriteString(ansiReset)
		return
	}
	b.WriteString(text)
}

func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(int(f.Fd()))
}
