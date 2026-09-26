package tools

import (
	"context"

	"github.com/lordnynex/mcp/robotgo/desktop"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type observeInput struct {
	DisplayID *int   `json:"displayId,omitempty" jsonschema:"Display index; omit for the default display"`
	X         *int   `json:"x,omitempty" jsonschema:"Region left; set x,y,w,h together"`
	Y         *int   `json:"y,omitempty" jsonschema:"Region top"`
	W         *int   `json:"w,omitempty" jsonschema:"Region width"`
	H         *int   `json:"h,omitempty" jsonschema:"Region height"`
	MaxWidth  *int   `json:"maxWidth,omitempty" jsonschema:"Max encoded width in pixels (default 1280)"`
	Format    string `json:"format,omitempty" jsonschema:"png or jpeg"`
}

type observeOutput struct {
	CursorX     int               `json:"cursorX"`
	CursorY     int               `json:"cursorY"`
	Displays    []desktop.Display `json:"displays"`
	WindowTitle string            `json:"windowTitle,omitempty"`
	Width       int               `json:"width"`
	Height      int               `json:"height"`
	OriginX     int               `json:"originX"`
	OriginY     int               `json:"originY"`
	MouseWidth  int               `json:"mouseWidth"`
	MouseHeight int               `json:"mouseHeight"`
	Format      string            `json:"format"`
	CmdCtrl     string            `json:"cmdCtrl"`
}

func registerObserve(s *mcp.Server, d desktop.Driver) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "ObserveDesktop",
		Title:       humanTitle("ObserveDesktop"),
		Description: "Capture the screen and report cursor position, display bounds, and the focused window title. Call once to plan, then one DesktopOps (MousePath may be an op), then verify. " + imageToMouseNote,
		Annotations: annotations("ObserveDesktop"),
	}, func(ctx context.Context, req *mcp.CallToolRequest, in observeInput) (*mcp.CallToolResult, *observeOutput, error) {
		notifyProgress(ctx, req, 1, 2, "capturing screen")
		spec, err := captureSpec(in.DisplayID, in.X, in.Y, in.W, in.H)
		if err != nil {
			return nil, nil, err
		}
		data, mime, w, h, err := encodeCapture(d, spec, in.MaxWidth, in.Format)
		if err != nil {
			return nil, nil, err
		}
		cx, cy := d.Location()
		region := mouseRegion(d, spec)
		out := &observeOutput{
			CursorX:     cx,
			CursorY:     cy,
			Displays:    d.Displays(),
			WindowTitle: d.WindowTitle(),
			Width:       w,
			Height:      h,
			OriginX:     region.OriginX,
			OriginY:     region.OriginY,
			MouseWidth:  region.MouseWidth,
			MouseHeight: region.MouseHeight,
			Format:      formatFromMIME(mime),
			CmdCtrl:     d.CmdCtrl(),
		}
		notifyProgress(ctx, req, 2, 2, "screenshot ready")
		return imageResult(data, mime), out, nil
	})
}
