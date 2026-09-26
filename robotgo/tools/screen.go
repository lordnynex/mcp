package tools

import (
	"context"

	"github.com/lordnynex/mcp/robotgo/desktop"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type captureInput struct {
	DisplayID *int   `json:"displayId,omitempty" jsonschema:"Display index; omit for the default display"`
	X         *int   `json:"x,omitempty" jsonschema:"Region left; set x,y,w,h together"`
	Y         *int   `json:"y,omitempty" jsonschema:"Region top"`
	W         *int   `json:"w,omitempty" jsonschema:"Region width"`
	H         *int   `json:"h,omitempty" jsonschema:"Region height"`
	MaxWidth  *int   `json:"maxWidth,omitempty" jsonschema:"Max encoded width in pixels"`
	Format    string `json:"format,omitempty" jsonschema:"png or jpeg"`
}

type captureOutput struct {
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	OriginX     int    `json:"originX"`
	OriginY     int    `json:"originY"`
	MouseWidth  int    `json:"mouseWidth"`
	MouseHeight int    `json:"mouseHeight"`
	Format      string `json:"format"`
}

type screenInfoOutput struct {
	Displays     []desktop.Display `json:"displays"`
	ScreenWidth  int               `json:"screenWidth"`
	ScreenHeight int               `json:"screenHeight"`
	CmdCtrl      string            `json:"cmdCtrl"`
}

type pixelInput struct {
	X         int  `json:"x,omitempty" jsonschema:"Pixel X (ignored when atCursor is true)"`
	Y         int  `json:"y,omitempty" jsonschema:"Pixel Y"`
	DisplayID *int `json:"displayId,omitempty" jsonschema:"Optional display index"`
	AtCursor  bool `json:"atCursor,omitempty" jsonschema:"Read the pixel under the cursor"`
}

type pixelOutput struct {
	Hex string `json:"hex"`
	X   int    `json:"x"`
	Y   int    `json:"y"`
}

func registerScreen(s *mcp.Server, d desktop.Driver) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "CaptureScreen",
		Title:       humanTitle("CaptureScreen"),
		Description: "Capture the screen (or a region) and return an image. Prefer ObserveDesktop when you also need cursor and window metadata. " + imageToMouseNote,
		Annotations: annotations("CaptureScreen"),
	}, func(ctx context.Context, req *mcp.CallToolRequest, in captureInput) (*mcp.CallToolResult, *captureOutput, error) {
		notifyProgress(ctx, req, 1, 2, "capturing screen")
		spec, err := captureSpec(in.DisplayID, in.X, in.Y, in.W, in.H)
		if err != nil {
			return nil, nil, err
		}
		data, mime, w, h, err := encodeCapture(d, spec, in.MaxWidth, in.Format)
		if err != nil {
			return nil, nil, err
		}
		region := mouseRegion(d, spec)
		notifyProgress(ctx, req, 2, 2, "screenshot ready")
		return imageResult(data, mime), &captureOutput{
			Width: w, Height: h, Format: formatFromMIME(mime),
			OriginX: region.OriginX, OriginY: region.OriginY,
			MouseWidth: region.MouseWidth, MouseHeight: region.MouseHeight,
		}, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "GetScreenInfo",
		Title:       humanTitle("GetScreenInfo"),
		Description: "List displays, bounds, scale sizes, and the platform accelerator key (cmd vs ctrl).",
		Annotations: annotations("GetScreenInfo"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, *screenInfoOutput, error) {
		sw, sh := d.ScreenSize()
		return nil, &screenInfoOutput{
			Displays:     d.Displays(),
			ScreenWidth:  sw,
			ScreenHeight: sh,
			CmdCtrl:      d.CmdCtrl(),
		}, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "GetPixelColor",
		Title:       humanTitle("GetPixelColor"),
		Description: "Read the hex color at x,y or under the cursor.",
		Annotations: annotations("GetPixelColor"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, in pixelInput) (*mcp.CallToolResult, *pixelOutput, error) {
		displayID := -1
		if in.DisplayID != nil {
			displayID = *in.DisplayID
		}
		hex, err := d.PixelColor(in.X, in.Y, displayID, in.AtCursor)
		if err != nil {
			return nil, nil, err
		}
		x, y := in.X, in.Y
		if in.AtCursor {
			x, y = d.Location()
		}
		return nil, &pixelOutput{Hex: hex, X: x, Y: y}, nil
	})
}
