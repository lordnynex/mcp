package tools

import (
	"context"
	"fmt"

	"github.com/lordnynex/mcp/robotgo/desktop"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type mousePathInput struct {
	Points    []pathPoint `json:"points" jsonschema:"Absolute or relative vertices in ObserveDesktop space"`
	Smooth    bool        `json:"smooth,omitempty" jsonschema:"MoveSmooth for every segment"`
	Relative  bool        `json:"relative,omitempty" jsonschema:"Treat each point as a delta"`
	Hold      bool        `json:"hold,omitempty" jsonschema:"Left-button down after the first point, up after the last"`
	DisplayID *int        `json:"displayId,omitempty" jsonschema:"Target display for absolute moves"`
	Low       float64     `json:"low,omitempty" jsonschema:"Smooth speed low"`
	High      float64     `json:"high,omitempty" jsonschema:"Smooth speed high"`
	Delay     int         `json:"delay,omitempty" jsonschema:"Smooth mouse delay in ms"`
}

type pathPoint struct {
	X int `json:"x" jsonschema:"X"`
	Y int `json:"y" jsonschema:"Y"`
}

type mousePathOutput struct {
	Points []pointOutput `json:"points"`
}

func registerPath(s *mcp.Server, d desktop.Driver) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "MousePath",
		Title:       humanTitle("MousePath"),
		Description: "Move through an array of points in one call. Shared smooth/relative/displayId apply to every segment. Use smooth false for geometric traces (star, zigzag). hold presses the left button after arriving at the first point and releases after the last (drag-select or draw). Hover-only by default. Also valid as a DesktopOps op.",
		Annotations: annotations("MousePath"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, in mousePathInput) (*mcp.CallToolResult, *mousePathOutput, error) {
		out, err := doPath(d, in)
		if err != nil {
			return nil, nil, err
		}
		return nil, &out, nil
	})
}

func doPath(d desktop.Driver, in mousePathInput) (mousePathOutput, error) {
	if len(in.Points) == 0 {
		return mousePathOutput{}, fmt.Errorf("points must not be empty")
	}
	if len(in.Points) > maxDesktopOps {
		return mousePathOutput{}, fmt.Errorf("points exceeds cap of %d", maxDesktopOps)
	}
	displayID := -1
	if in.DisplayID != nil {
		displayID = *in.DisplayID
	}
	out := mousePathOutput{Points: make([]pointOutput, 0, len(in.Points))}
	held := false
	if in.Hold {
		defer func() {
			if held {
				_ = d.Toggle("left", "up")
			}
		}()
	}
	for i, p := range in.Points {
		if err := d.Move(p.X, p.Y, displayID, in.Relative, in.Smooth, in.Low, in.High, in.Delay); err != nil {
			return out, err
		}
		out.Points = append(out.Points, locationOf(d))
		if in.Hold && i == 0 {
			if err := d.Toggle("left", "down"); err != nil {
				return out, err
			}
			held = true
		}
	}
	return out, nil
}
