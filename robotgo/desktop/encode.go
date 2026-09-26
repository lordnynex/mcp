package desktop

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"sync"
)

const (
	DefaultMaxWidth = 1280
	DefaultFormat   = "png"
)

// EncodeOptions control how screenshots are serialized for MCP image content.
type EncodeOptions struct {
	MaxWidth int
	Format   string
}

var (
	encodeMu sync.RWMutex
	encode   = EncodeOptions{MaxWidth: DefaultMaxWidth, Format: DefaultFormat}
)

// SetEncodeOptions sets process-wide screenshot encoding defaults.
func SetEncodeOptions(o EncodeOptions) {
	encodeMu.Lock()
	defer encodeMu.Unlock()
	if o.MaxWidth <= 0 {
		o.MaxWidth = DefaultMaxWidth
	}
	if o.Format == "" {
		o.Format = DefaultFormat
	}
	encode = o
}

// CurrentEncode returns screenshot encoding defaults.
func CurrentEncode() EncodeOptions {
	encodeMu.RLock()
	defer encodeMu.RUnlock()
	return encode
}

// EncodeImage resizes img to maxWidth and encodes it as png or jpeg.
func EncodeImage(img image.Image, maxWidth int, format string) (data []byte, mime string, width, height int, err error) {
	if img == nil {
		return nil, "", 0, 0, fmt.Errorf("no image")
	}
	format, err = NormalizeFormat(format)
	if err != nil {
		return nil, "", 0, 0, err
	}
	scaled := scaleToMaxWidth(img, maxWidth)
	b := scaled.Bounds()
	width, height = b.Dx(), b.Dy()
	var buf bytes.Buffer
	switch format {
	case "jpeg":
		if err := jpeg.Encode(&buf, scaled, &jpeg.Options{Quality: 80}); err != nil {
			return nil, "", 0, 0, err
		}
		return buf.Bytes(), "image/jpeg", width, height, nil
	default:
		if err := png.Encode(&buf, scaled); err != nil {
			return nil, "", 0, 0, err
		}
		return buf.Bytes(), "image/png", width, height, nil
	}
}

func scaleToMaxWidth(img image.Image, maxW int) image.Image {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	if maxW <= 0 || w <= maxW {
		return img
	}
	newW := maxW
	newH := h * maxW / w
	if newH < 1 {
		newH = 1
	}
	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	for y := range newH {
		for x := range newW {
			srcX := b.Min.X + x*w/newW
			srcY := b.Min.Y + y*h/newH
			dst.Set(x, y, img.At(srcX, srcY))
		}
	}
	return dst
}
