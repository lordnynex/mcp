package tools

import (
	"fmt"

	"github.com/lordnynex/mcp/robotgo/desktop"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func captureSpec(displayID, x, y, w, h *int) (desktop.CaptureSpec, error) {
	spec := desktop.CaptureSpec{DisplayID: -1}
	if id, ok := ptrInt(displayID); ok {
		spec.DisplayID = id
	}
	has := x != nil || y != nil || w != nil || h != nil
	if !has {
		return spec, nil
	}
	if x == nil || y == nil || w == nil || h == nil {
		return spec, fmt.Errorf("x, y, w, and h must be set together")
	}
	if *w <= 0 || *h <= 0 {
		return spec, fmt.Errorf("w and h must be positive")
	}
	spec.X, spec.Y, spec.W, spec.H = *x, *y, *w, *h
	spec.HasRegion = true
	return spec, nil
}

// captureMap is the mouse-space region that an encoded screenshot covers.
type captureMap struct {
	OriginX     int `json:"originX"`
	OriginY     int `json:"originY"`
	MouseWidth  int `json:"mouseWidth"`
	MouseHeight int `json:"mouseHeight"`
}

const imageToMouseNote = `width/height are encoded image pixels. Mouse coordinates: mouseX = originX + ix * mouseWidth / width; mouseY = originY + iy * mouseHeight / height.`

func mouseRegion(d desktop.Driver, spec desktop.CaptureSpec) captureMap {
	if spec.HasRegion {
		return captureMap{OriginX: spec.X, OriginY: spec.Y, MouseWidth: spec.W, MouseHeight: spec.H}
	}
	id := spec.DisplayID
	for _, disp := range d.Displays() {
		if id >= 0 && disp.ID != id {
			continue
		}
		if id >= 0 || disp.Main {
			return captureMap{OriginX: disp.X, OriginY: disp.Y, MouseWidth: disp.W, MouseHeight: disp.H}
		}
	}
	if displays := d.Displays(); len(displays) > 0 {
		disp := displays[0]
		return captureMap{OriginX: disp.X, OriginY: disp.Y, MouseWidth: disp.W, MouseHeight: disp.H}
	}
	w, h := d.ScreenSize()
	return captureMap{MouseWidth: w, MouseHeight: h}
}

func formatFromMIME(mime string) string {
	if mime == "image/jpeg" {
		return "jpeg"
	}
	return "png"
}

func encodeCapture(d desktop.Driver, spec desktop.CaptureSpec, maxWidth *int, format string) ([]byte, string, int, int, error) {
	img, err := d.Capture(spec)
	if err != nil {
		return nil, "", 0, 0, err
	}
	opts := desktop.CurrentEncode()
	mw := opts.MaxWidth
	if maxWidth != nil && *maxWidth > 0 {
		mw = *maxWidth
	}
	if format == "" {
		format = opts.Format
	}
	return desktop.EncodeImage(img, mw, format)
}

func imageResult(data []byte, mime string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.ImageContent{Data: data, MIMEType: mime},
		},
	}
}
