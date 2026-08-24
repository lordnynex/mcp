package session

import (
	"testing"

	"github.com/andreykaipov/goobs/api/events"
	"github.com/lordnynex/gest"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestHost(t *testing.T) {
	gest.Run(t, "Host", func(s *gest.S) {
		s.It("fails Connect to an unreachable host", func(t *gest.T) {
			h := New()
			_, err := h.Connect(ConnectInput{Host: "127.0.0.1:1"})
			t.Expect(err).NotTo(gest.BeNil())
			t.Expect(h.Status().Connected).To(gest.Equal(false))
		})

		s.It("Disconnect without a connection returns ErrNotConnected", func(t *gest.T) {
			h := New()
			err := h.Disconnect()
			t.Expect(err).To(gest.MatchErrorIs(ErrNotConnected))
		})

		s.It("RequireClient returns ErrNotConnected", func(t *gest.T) {
			h := New()
			_, err := h.RequireClient()
			t.Expect(err).To(gest.MatchErrorIs(ErrNotConnected))
		})

		s.It("requires high-volume bits that were not identified", func(t *gest.T) {
			h := New()
			err := h.EnsureHighVolume([]string{"InputVolumeMeters"})
			t.Expect(err).NotTo(gest.BeNil())
		})

		s.It("records events and fans out only when a session subscribed", func(t *gest.T) {
			h := New()
			h.ClearEvents()
			ev := &events.CurrentProgramSceneChanged{SceneName: "Cam"}
			h.OnOBSEvent(ev)
			got, ok := h.LatestEvent("CurrentProgramSceneChanged")
			t.Expect(ok).To(gest.Equal(true))
			t.Expect(got.EventType).To(gest.Equal("CurrentProgramSceneChanged"))
			t.Expect(h.AnyInterest("CurrentProgramSceneChanged")).To(gest.Equal(false))

			sess := &mcp.ServerSession{}
			h.SetInterest(sess, []string{"CurrentProgramSceneChanged"}, true)
			t.Expect(h.AnyInterest("CurrentProgramSceneChanged")).To(gest.Equal(true))
			h.OnOBSEvent(ev)
			buf := h.EventBuffer()
			t.Expect(len(buf) >= 2).To(gest.Equal(true))
			h.ClearInterest(sess)
			t.Expect(h.AnyInterest("CurrentProgramSceneChanged")).To(gest.Equal(false))
		})
	})
}
