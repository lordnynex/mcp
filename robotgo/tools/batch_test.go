package tools

import (
	"context"
	"testing"

	"github.com/lordnynex/gest"
	"github.com/lordnynex/mcp/robotgo/desktop"
)

func TestDesktopOps(t *testing.T) {
	gest.Run(t, "DesktopOps", func(s *gest.S) {
		s.It("runs move click type in order", func(t *gest.T) {
			fake := desktop.NewFake()
			out, img, err := runDesktopOps(context.Background(), nil, fake, desktopOpsInput{
				Ops: []desktopOp{
					{Op: "MouseMove", X: ip(10), Y: ip(20)},
					{Op: "MouseClick"},
					{Op: "Type", Text: "hi"},
				},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(img).To(gest.BeNil())
			t.Expect(len(out.Results)).To(gest.Equal(3))
			t.Expect(out.Results[0].OK).To(gest.Equal(true))
			t.Expect(out.Results[1].OK).To(gest.Equal(true))
			t.Expect(out.Results[2].OK).To(gest.Equal(true))
			t.Expect(out.CursorX).To(gest.Equal(10))
			t.Expect(out.CursorY).To(gest.Equal(20))
			t.Require(len(fake.Calls) >= 3).To(gest.Equal(true))
			t.Expect(fake.Calls[0].Name).To(gest.Equal("Move"))
			t.Expect(fake.Calls[1].Name).To(gest.Equal("Click"))
			t.Expect(fake.Calls[2].Name).To(gest.Equal("Type"))
		})

		s.It("halts on failure by default", func(t *gest.T) {
			fake := desktop.NewFake()
			out, _, err := runDesktopOps(context.Background(), nil, fake, desktopOpsInput{
				Ops: []desktopOp{
					{Op: "MouseMove", X: ip(1), Y: ip(2)},
					{Op: "KeyTap", Key: "not-a-key"},
					{Op: "Type", Text: "nope"},
				},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(out.Halted).To(gest.Equal(true))
			t.Expect(len(out.Results)).To(gest.Equal(2))
			t.Expect(out.Results[1].OK).To(gest.Equal(false))
			t.Expect(fake.Called("Type")).To(gest.Equal(false))
		})

		s.It("continues when haltOnFailure is false", func(t *gest.T) {
			fake := desktop.NewFake()
			halt := false
			out, _, err := runDesktopOps(context.Background(), nil, fake, desktopOpsInput{
				HaltOnFailure: &halt,
				Ops: []desktopOp{
					{Op: "KeyTap", Key: "not-a-key"},
					{Op: "Type", Text: "ok"},
				},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(out.Halted).To(gest.Equal(false))
			t.Expect(len(out.Results)).To(gest.Equal(2))
			t.Expect(out.Results[1].OK).To(gest.Equal(true))
			t.Expect(fake.Called("Type")).To(gest.Equal(true))
		})

		s.It("rejects empty ops", func(t *gest.T) {
			_, _, err := runDesktopOps(context.Background(), nil, desktop.NewFake(), desktopOpsInput{})
			t.Expect(err).NotTo(gest.BeNil())
		})

		s.It("rejects an unknown op", func(t *gest.T) {
			_, _, err := runDesktopOps(context.Background(), nil, desktop.NewFake(), desktopOpsInput{
				Ops: []desktopOp{{Op: "InventedOp"}},
			})
			t.Expect(err).NotTo(gest.BeNil())
			t.Expect(err.Error()).To(gest.MatchRegexp("unknown op"))
		})

		s.It("rejects KillProcess", func(t *gest.T) {
			fake := desktop.NewFake()
			_, _, err := runDesktopOps(context.Background(), nil, fake, desktopOpsInput{
				Ops: []desktopOp{{Op: "KillProcess", PID: 42}},
			})
			t.Expect(err).NotTo(gest.BeNil())
			t.Expect(err.Error()).To(gest.MatchRegexp("KillProcess"))
			t.Expect(fake.Called("Kill")).To(gest.Equal(false))
		})

		s.It("rejects over the ops cap", func(t *gest.T) {
			ops := make([]desktopOp, maxDesktopOps+1)
			for i := range ops {
				ops[i] = desktopOp{Op: "MouseLocation"}
			}
			_, _, err := runDesktopOps(context.Background(), nil, desktop.NewFake(), desktopOpsInput{Ops: ops})
			t.Expect(err).NotTo(gest.BeNil())
		})

		s.It("runs window ops and returns payloads", func(t *gest.T) {
			fake := desktop.NewFake()
			out, _, err := runDesktopOps(context.Background(), nil, fake, desktopOpsInput{
				Ops: []desktopOp{
					{Op: "FocusWindow", Name: "FakeApp"},
					{Op: "GetWindowBounds", PID: 42},
					{Op: "CloseWindow", PID: 42},
					{Op: "GetWindowTitle"},
				},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(len(out.Results)).To(gest.Equal(4))
			t.Expect(fake.Called("FocusWindow")).To(gest.Equal(true))
			t.Expect(fake.Called("WindowBounds")).To(gest.Equal(true))
			t.Expect(fake.Called("CloseWindow")).To(gest.Equal(true))
			bounds, ok := out.Results[1].Response.(boundsOutput)
			t.Require(ok).To(gest.Equal(true))
			t.Expect(bounds.W).To(gest.Equal(800))
			title, ok := out.Results[3].Response.(titleOutput)
			t.Require(ok).To(gest.Equal(true))
			t.Expect(title.Title).To(gest.Equal("FakeApp"))
		})

		s.It("MouseSelect toggles down then up", func(t *gest.T) {
			fake := desktop.NewFake()
			out, _, err := runDesktopOps(context.Background(), nil, fake, desktopOpsInput{
				Ops: []desktopOp{{
					Op: "MouseSelect", X: ip(400), Y: ip(300), ToX: ip(720), ToY: ip(300),
				}},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(out.Results[0].OK).To(gest.Equal(true))
			t.Require(len(fake.Calls)).To(gest.Equal(4))
			t.Expect(fake.Calls[0].Name).To(gest.Equal("Move"))
			t.Expect(fake.Calls[1].Name).To(gest.Equal("Toggle"))
			t.Expect(fake.Calls[1].Args[1]).To(gest.Equal("down"))
			t.Expect(fake.Calls[2].Name).To(gest.Equal("Move"))
			t.Expect(fake.Calls[3].Name).To(gest.Equal("Toggle"))
			t.Expect(fake.Calls[3].Args[1]).To(gest.Equal("up"))
		})

		s.It("observeAfter captures once", func(t *gest.T) {
			fake := desktop.NewFake()
			out, img, err := runDesktopOps(context.Background(), nil, fake, desktopOpsInput{
				ObserveAfter: true,
				Ops:          []desktopOp{{Op: "MouseMove", X: ip(3), Y: ip(4)}},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(fake.Called("Capture")).To(gest.Equal(true))
			t.Expect(img).NotTo(gest.BeNil())
			t.Expect(out.Width > 0).To(gest.Equal(true))
			t.Expect(out.Format).To(gest.Equal("png"))
			t.Expect(out.MouseWidth).To(gest.Equal(1920))
		})

		s.It("observeAfter crop records the region", func(t *gest.T) {
			fake := desktop.NewFake()
			out, img, err := runDesktopOps(context.Background(), nil, fake, desktopOpsInput{
				ObserveAfter: true,
				X:            ip(1260),
				Y:            ip(790),
				W:            ip(240),
				H:            ip(160),
				Ops:          []desktopOp{{Op: "MouseMove", X: ip(3), Y: ip(4)}},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(img).NotTo(gest.BeNil())
			t.Expect(out.OriginX).To(gest.Equal(1260))
			t.Expect(out.OriginY).To(gest.Equal(790))
			t.Expect(out.MouseWidth).To(gest.Equal(240))
			t.Expect(out.MouseHeight).To(gest.Equal(160))
			call, err := fake.Last("Capture")
			t.Require(err).To(gest.BeNil())
			spec, ok := call.Args[0].(desktop.CaptureSpec)
			t.Require(ok).To(gest.Equal(true))
			t.Expect(spec.HasRegion).To(gest.Equal(true))
			t.Expect(spec.X).To(gest.Equal(1260))
			t.Expect(spec.W).To(gest.Equal(240))
		})

		s.It("skips observeAfter when the batch halted", func(t *gest.T) {
			fake := desktop.NewFake()
			out, img, err := runDesktopOps(context.Background(), nil, fake, desktopOpsInput{
				ObserveAfter: true,
				Ops:          []desktopOp{{Op: "KeyTap", Key: "not-a-key"}},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(out.Halted).To(gest.Equal(true))
			t.Expect(img).To(gest.BeNil())
			t.Expect(fake.Called("Capture")).To(gest.Equal(false))
		})

		s.It("observeOnError captures after a halt", func(t *gest.T) {
			fake := desktop.NewFake()
			out, img, err := runDesktopOps(context.Background(), nil, fake, desktopOpsInput{
				ObserveAfter:   true,
				ObserveOnError: true,
				Ops:            []desktopOp{{Op: "KeyTap", Key: "not-a-key"}},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(out.Halted).To(gest.Equal(true))
			t.Expect(img).NotTo(gest.BeNil())
			t.Expect(fake.Called("Capture")).To(gest.Equal(true))
		})

		s.It("runs MousePath then Type", func(t *gest.T) {
			fake := desktop.NewFake()
			out, _, err := runDesktopOps(context.Background(), nil, fake, desktopOpsInput{
				Ops: []desktopOp{
					{Op: "MousePath", Points: []pathPoint{{X: 1, Y: 2}, {X: 3, Y: 4}}},
					{Op: "Type", Text: "hi"},
				},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(len(out.Results)).To(gest.Equal(2))
			t.Expect(out.Results[0].OK).To(gest.Equal(true))
			t.Expect(out.Results[1].OK).To(gest.Equal(true))
			n := 0
			for _, c := range fake.Calls {
				if c.Name == "Move" {
					n++
				}
			}
			t.Expect(n).To(gest.Equal(2))
			t.Expect(fake.Called("Type")).To(gest.Equal(true))
		})
	})
}

func TestMousePath(t *testing.T) {
	gest.Run(t, "MousePath", func(s *gest.S) {
		s.It("visits every point", func(t *gest.T) {
			fake := desktop.NewFake()
			out, err := doPath(fake, mousePathInput{
				Points: []pathPoint{{X: 1, Y: 2}, {X: 3, Y: 4}, {X: 5, Y: 6}},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(len(out.Points)).To(gest.Equal(3))
			t.Expect(out.Points[2]).To(gest.Equal(pointOutput{X: 5, Y: 6}))
			t.Expect(fake.Called("Toggle")).To(gest.Equal(false))
			n := 0
			for _, c := range fake.Calls {
				if c.Name == "Move" {
					n++
				}
			}
			t.Expect(n).To(gest.Equal(3))
		})

		s.It("hold toggles the button after the first point", func(t *gest.T) {
			fake := desktop.NewFake()
			_, err := doPath(fake, mousePathInput{
				Points: []pathPoint{{X: 10, Y: 10}, {X: 20, Y: 20}},
				Hold:   true,
			})
			t.Require(err).To(gest.BeNil())
			t.Require(len(fake.Calls)).To(gest.Equal(4))
			t.Expect(fake.Calls[0].Name).To(gest.Equal("Move"))
			t.Expect(fake.Calls[1].Name).To(gest.Equal("Toggle"))
			t.Expect(fake.Calls[1].Args[1]).To(gest.Equal("down"))
			t.Expect(fake.Calls[2].Name).To(gest.Equal("Move"))
			t.Expect(fake.Calls[3].Name).To(gest.Equal("Toggle"))
			t.Expect(fake.Calls[3].Args[1]).To(gest.Equal("up"))
		})

		s.It("rejects empty points", func(t *gest.T) {
			_, err := doPath(desktop.NewFake(), mousePathInput{})
			t.Expect(err).NotTo(gest.BeNil())
		})
	})
}

func ip(v int) *int { return &v }
