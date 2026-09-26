package desktop

import (
	"testing"

	"github.com/lordnynex/gest"
)

func TestKeys(t *testing.T) {
	gest.Run(t, "Keys", func(s *gest.S) {
		s.It("accepts named keys and single letters", func(t *gest.T) {
			t.Expect(ValidKey("enter")).To(gest.Equal(true))
			t.Expect(ValidKey("a")).To(gest.Equal(true))
			t.Expect(ValidKey("7")).To(gest.Equal(true))
			t.Expect(CheckKey("not-a-key")).NotTo(gest.BeNil())
		})

		s.It("validates buttons and scroll dirs", func(t *gest.T) {
			t.Expect(CheckButton("left")).To(gest.BeNil())
			t.Expect(CheckButton("middle-finger")).NotTo(gest.BeNil())
			t.Expect(CheckScrollDir("down")).To(gest.BeNil())
			t.Expect(CheckScrollDir("sideways")).NotTo(gest.BeNil())
		})

		s.It("normalizes image formats", func(t *gest.T) {
			f, err := NormalizeFormat("JPG")
			t.Require(err).To(gest.BeNil())
			t.Expect(f).To(gest.Equal("jpeg"))
			_, err = NormalizeFormat("gif")
			t.Expect(err).NotTo(gest.BeNil())
		})
	})
}

func TestEncodeImage(t *testing.T) {
	gest.Run(t, "EncodeImage", func(s *gest.S) {
		s.It("encodes a fake screenshot as png", func(t *gest.T) {
			img := NewFake().Image
			data, mime, w, h, err := EncodeImage(img, 1280, "png")
			t.Require(err).To(gest.BeNil())
			t.Expect(mime).To(gest.Equal("image/png"))
			t.Expect(len(data) > 0).To(gest.Equal(true))
			t.Expect(w).To(gest.Equal(4))
			t.Expect(h).To(gest.Equal(4))
		})
	})
}
