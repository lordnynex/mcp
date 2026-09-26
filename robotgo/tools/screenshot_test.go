package tools

import (
	"testing"

	"github.com/lordnynex/gest"
	"github.com/lordnynex/mcp/robotgo/desktop"
)

func TestMouseRegion(t *testing.T) {
	gest.Run(t, "MouseRegion", func(s *gest.S) {
		s.It("uses the main display for a full capture", func(t *gest.T) {
			fake := desktop.NewFake()
			m := mouseRegion(fake, desktop.CaptureSpec{DisplayID: -1})
			t.Expect(m.OriginX).To(gest.Equal(0))
			t.Expect(m.OriginY).To(gest.Equal(0))
			t.Expect(m.MouseWidth).To(gest.Equal(1920))
			t.Expect(m.MouseHeight).To(gest.Equal(1080))
		})

		s.It("uses the crop origin and size", func(t *gest.T) {
			fake := desktop.NewFake()
			m := mouseRegion(fake, desktop.CaptureSpec{
				HasRegion: true, X: 10, Y: 20, W: 30, H: 40,
			})
			t.Expect(m.OriginX).To(gest.Equal(10))
			t.Expect(m.OriginY).To(gest.Equal(20))
			t.Expect(m.MouseWidth).To(gest.Equal(30))
			t.Expect(m.MouseHeight).To(gest.Equal(40))
		})
	})
}
