package logger

import (
	"bytes"
	"log/slog"
	"os"
	"testing"

	"github.com/lordnynex/gest"
)

func TestLogger(t *testing.T) {
	gest.Run(t, "Logger", func(s *gest.S) {
		s.Describe("color", func(s *gest.S) {
			s.It("emits ANSI when color is forced on", func(t *gest.T) {
				out := logLine(t, Options{Color: ptr(true), Level: slog.LevelInfo}, func(l *slog.Logger) {
					l.Error("boom")
				})
				t.Expect(out).To(gest.Contain("\x1b["))
				t.Expect(out).To(gest.Contain(ansiRed))
				t.Expect(out).To(gest.Contain("ERROR"))
				t.Expect(out).To(gest.Contain("boom"))
			})

			s.It("omits ANSI when color is forced off", func(t *gest.T) {
				out := logLine(t, Options{Color: ptr(false)}, func(l *slog.Logger) {
					l.Info("hello")
				})
				t.Expect(out).NotTo(gest.Contain("\x1b["))
				t.Expect(out).To(gest.Contain("INFO"))
				t.Expect(out).To(gest.Contain("hello"))
			})
		})

		s.Describe("attrs", func(s *gest.S) {
			s.It("writes key-value attrs", func(t *gest.T) {
				out := logLine(t, Options{Color: ptr(false)}, func(l *slog.Logger) {
					l.Info("connected", "user", "ada")
				})
				t.Expect(out).To(gest.Contain("user=ada"))
			})

			s.It("prefixes grouped attrs", func(t *gest.T) {
				out := logLine(t, Options{Color: ptr(false)}, func(l *slog.Logger) {
					l.WithGroup("obs").Info("ready", "host", "localhost")
				})
				t.Expect(out).To(gest.Contain("obs.host=localhost"))
			})
		})

		s.Describe("defaults", func(s *gest.S) {
			s.It("writes to stderr when Writer is nil", func(t *gest.T) {
				t.Expect(writerOrDefault(nil) == os.Stderr).To(gest.BeTrue())
			})

			s.It("SetDefault does not panic", func(t *gest.T) {
				t.Expect(func() {
					SetDefault(Options{Writer: ioDiscard{}, Color: ptr(false)})
				}).NotTo(gest.Panic())
			})
		})
	})
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) { return len(p), nil }

func logLine(t *gest.T, opts Options, fn func(*slog.Logger)) string {
	t.Helper()
	var buf bytes.Buffer
	opts.Writer = &buf
	fn(New(opts))
	return buf.String()
}

func ptr[T any](v T) *T {
	return &v
}
