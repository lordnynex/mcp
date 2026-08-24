package protocol

import (
	"os"
	"testing"

	"github.com/andreykaipov/goobs/api/events/subscriptions"
	"github.com/lordnynex/gest"
)

func TestProtocolCoverage(t *testing.T) {
	gest.Run(t, "ProtocolCoverage", func(s *gest.S) {
		s.It("matches testdata/protocol.md request TOC", func(t *gest.T) {
			md := readProtocol(t)
			names := ParseTOCNames(md, "## Requests Table of Contents")
			t.Expect(names).To(gest.Equal(RequestNames))
			t.Expect(RequestNames).To(gest.HaveLen(147))
		})

		s.It("matches testdata/protocol.md event TOC", func(t *gest.T) {
			md := readProtocol(t)
			names := ParseTOCNames(md, "## Events Table of Contents")
			t.Expect(names).To(gest.Equal(EventNames))
			t.Expect(EventNames).To(gest.HaveLen(60))
		})
	})
}

func TestExpandEventTypes(t *testing.T) {
	gest.Run(t, "ExpandEventTypes", func(s *gest.S) {
		s.It("rejects unknown event names", func(t *gest.T) {
			_, err := ExpandEventTypes([]string{"NotAnOfficialEvent"})
			t.Expect(err).NotTo(gest.BeNil())
		})

		s.It("expands category aliases without high-volume events", func(t *gest.T) {
			types, err := ExpandEventTypes([]string{"Scenes"})
			t.Require(err).To(gest.BeNil())
			t.Expect(types).To(gest.Contain("CurrentProgramSceneChanged"))
			t.Expect(types).NotTo(gest.Contain("SceneItemTransformChanged"))
			t.Expect(types).NotTo(gest.Contain("InputVolumeMeters"))
		})

		s.It("accepts every official event name", func(t *gest.T) {
			types, err := ExpandEventTypes(EventNames)
			t.Require(err).To(gest.BeNil())
			t.Expect(types).To(gest.Equal(EventNames))
			for _, n := range types {
				t.Expect(EventResourceURI(n)).To(gest.Equal("obs://events/" + n))
			}
		})

		s.It("maps high-volume names to official bits", func(t *gest.T) {
			t.Expect(HighVolumeEvents["InputVolumeMeters"]).To(gest.Equal(subscriptions.InputVolumeMeters))
			t.Expect(HighVolumeEvents["SceneItemTransformChanged"]).To(gest.Equal(subscriptions.SceneItemTransformChanged))
		})
	})
}

func readProtocol(t *gest.T) string {
	t.Helper()
	b, err := os.ReadFile("../testdata/protocol.md")
	t.Require(err).To(gest.BeNil())
	return string(b)
}
