package tools

import (
	"context"
	"fmt"

	"github.com/lordnynex/mcp/robotgo/desktop"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const maxDesktopOps = 256

const desktopOpsDesc = `Run a serial batch of desktop operations in one call. Prefer this over one MouseMove/MouseClick/Type per turn when the sequence is already known from ObserveDesktop. Each ops[] item has op plus the same fields as that tool. Allowed op: MouseMove, MouseClick, MouseToggle, MouseScroll, MouseDrag, MouseSelect, MousePath, Type, KeyTap, KeyToggle, KeyPress, ClipboardWrite, ClipboardPaste, Sleep, SetDelay, GetWindowTitle, FocusWindow, GetWindowBounds, MinWindow, MaxWindow, CloseWindow. haltOnFailure defaults true. observeAfter attaches one screenshot after the last successful batch; set x,y,w,h to crop that verify shot. observeOnError also captures after a halt. KillProcess is not allowed. Aim MouseClick x,y at the control center. Geometric MousePath traces should use smooth false.`

// DesktopOpNames is the allowed DesktopOps op list (existing tool names).
var DesktopOpNames = []string{
	"MouseMove", "MouseClick", "MouseToggle", "MouseScroll", "MouseDrag", "MouseSelect", "MousePath",
	"Type", "KeyTap", "KeyToggle", "KeyPress",
	"ClipboardWrite", "ClipboardPaste",
	"Sleep", "SetDelay",
	"GetWindowTitle", "FocusWindow", "GetWindowBounds", "MinWindow", "MaxWindow", "CloseWindow",
}

var desktopOpSet = func() map[string]struct{} {
	m := make(map[string]struct{}, len(DesktopOpNames))
	for _, n := range DesktopOpNames {
		m[n] = struct{}{}
	}
	return m
}()

type desktopOpsInput struct {
	HaltOnFailure  *bool       `json:"haltOnFailure,omitempty" jsonschema:"Stop the batch on first failure (default true)"`
	ObserveAfter   bool        `json:"observeAfter,omitempty" jsonschema:"Attach one screenshot after the last op"`
	ObserveOnError bool        `json:"observeOnError,omitempty" jsonschema:"Also screenshot when the batch halted (default false)"`
	DisplayID      *int        `json:"displayId,omitempty" jsonschema:"observeAfter display index"`
	X              *int        `json:"x,omitempty" jsonschema:"observeAfter crop left; set x,y,w,h together"`
	Y              *int        `json:"y,omitempty" jsonschema:"observeAfter crop top"`
	W              *int        `json:"w,omitempty" jsonschema:"observeAfter crop width"`
	H              *int        `json:"h,omitempty" jsonschema:"observeAfter crop height"`
	MaxWidth       *int        `json:"maxWidth,omitempty" jsonschema:"observeAfter max encoded width"`
	Format         string      `json:"format,omitempty" jsonschema:"observeAfter png or jpeg"`
	Ops            []desktopOp `json:"ops" jsonschema:"Serial operations; each item has op plus that tool's fields"`
}

type desktopOp struct {
	Op string `json:"op" jsonschema:"Tool name (MouseMove, MouseClick, Type, FocusWindow, …)"`

	X         *int    `json:"x,omitempty" jsonschema:"X / selection start X / scroll ticks"`
	Y         *int    `json:"y,omitempty" jsonschema:"Y / selection start Y / scroll ticks"`
	ToX       *int    `json:"toX,omitempty" jsonschema:"MouseSelect end X"`
	ToY       *int    `json:"toY,omitempty" jsonschema:"MouseSelect end Y"`
	Relative  bool    `json:"relative,omitempty" jsonschema:"MouseMove relative delta"`
	Smooth    bool    `json:"smooth,omitempty" jsonschema:"MoveSmooth / ScrollSmooth"`
	DisplayID *int    `json:"displayId,omitempty" jsonschema:"Target display"`
	Low       float64 `json:"low,omitempty" jsonschema:"Smooth speed low"`
	High      float64 `json:"high,omitempty" jsonschema:"Smooth speed high"`
	Delay     int     `json:"delay,omitempty" jsonschema:"Smooth mouse delay in ms"`

	Button string `json:"button,omitempty" jsonschema:"Mouse button"`
	Double bool   `json:"double,omitempty" jsonschema:"Double-click"`
	Count  int    `json:"count,omitempty" jsonschema:"Click count"`
	State  string `json:"state,omitempty" jsonschema:"down or up"`
	Dir    string `json:"dir,omitempty" jsonschema:"Scroll direction"`
	Steps  int    `json:"steps,omitempty" jsonschema:"Smooth scroll steps"`
	Sleep  int    `json:"sleepMs,omitempty" jsonschema:"Smooth scroll sleep"`

	Text      string   `json:"text,omitempty" jsonschema:"Type / clipboard text"`
	DelayMs   int      `json:"delayMs,omitempty" jsonschema:"Milliseconds between Type runes"`
	Key       string   `json:"key,omitempty" jsonschema:"KeyTap / KeyToggle / KeyPress key"`
	Modifiers []string `json:"modifiers,omitempty" jsonschema:"Key modifiers"`
	PID       int      `json:"pid,omitempty" jsonschema:"Process or window pid"`
	Name      string   `json:"name,omitempty" jsonschema:"FocusWindow name"`

	KeySleep   *int `json:"keySleep,omitempty" jsonschema:"SetDelay key sleep"`
	MouseSleep *int `json:"mouseSleep,omitempty" jsonschema:"SetDelay mouse sleep"`
	Ms         int  `json:"ms,omitempty" jsonschema:"Sleep milliseconds"`

	Points []pathPoint `json:"points,omitempty" jsonschema:"MousePath vertices"`
	Hold   bool        `json:"hold,omitempty" jsonschema:"MousePath left-button hold after the first point"`
}

type desktopOpsOutput struct {
	Results      []desktopOpResult `json:"results"`
	CursorX      int               `json:"cursorX"`
	CursorY      int               `json:"cursorY"`
	Halted       bool              `json:"halted,omitempty"`
	Width        int               `json:"width,omitempty"`
	Height       int               `json:"height,omitempty"`
	OriginX      int               `json:"originX,omitempty"`
	OriginY      int               `json:"originY,omitempty"`
	MouseWidth   int               `json:"mouseWidth,omitempty"`
	MouseHeight  int               `json:"mouseHeight,omitempty"`
	Format       string            `json:"format,omitempty"`
	WindowTitle  string            `json:"windowTitle,omitempty"`
	Displays     []desktop.Display `json:"displays,omitempty"`
	CmdCtrl      string            `json:"cmdCtrl,omitempty"`
	ObserveError string            `json:"observeError,omitempty"`
}

type desktopOpResult struct {
	Op       string `json:"op"`
	OK       bool   `json:"ok"`
	Error    string `json:"error,omitempty"`
	X        *int   `json:"x,omitempty"`
	Y        *int   `json:"y,omitempty"`
	Response any    `json:"response,omitempty"`
}

func registerBatch(s *mcp.Server, d desktop.Driver) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "DesktopOps",
		Title:       humanTitle("DesktopOps"),
		Description: desktopOpsDesc,
		Annotations: annotations("DesktopOps"),
	}, func(ctx context.Context, req *mcp.CallToolRequest, in desktopOpsInput) (*mcp.CallToolResult, *desktopOpsOutput, error) {
		out, img, err := runDesktopOps(ctx, req, d, in)
		if err != nil {
			return nil, nil, err
		}
		return img, &out, nil
	})
}

func runDesktopOps(ctx context.Context, req *mcp.CallToolRequest, d desktop.Driver, in desktopOpsInput) (desktopOpsOutput, *mcp.CallToolResult, error) {
	if len(in.Ops) == 0 {
		return desktopOpsOutput{}, nil, fmt.Errorf("ops must not be empty")
	}
	if len(in.Ops) > maxDesktopOps {
		return desktopOpsOutput{}, nil, fmt.Errorf("ops exceeds cap of %d", maxDesktopOps)
	}
	for _, op := range in.Ops {
		if err := checkDesktopOp(op.Op); err != nil {
			return desktopOpsOutput{}, nil, err
		}
	}
	halt := true
	if in.HaltOnFailure != nil {
		halt = *in.HaltOnFailure
	}
	total := float64(len(in.Ops))
	out := desktopOpsOutput{Results: make([]desktopOpResult, 0, len(in.Ops))}
	for i, op := range in.Ops {
		if ctx.Err() != nil {
			return out, nil, ctx.Err()
		}
		notifyProgress(ctx, req, float64(i), total, op.Op)
		item := runOneDesktopOp(d, op)
		out.Results = append(out.Results, item)
		if halt && !item.OK {
			out.Halted = true
			break
		}
	}
	notifyProgress(ctx, req, float64(len(out.Results)), total, "batch complete")
	loc := locationOf(d)
	out.CursorX, out.CursorY = loc.X, loc.Y

	var img *mcp.CallToolResult
	wantObserve := in.ObserveAfter && (!out.Halted || in.ObserveOnError)
	if wantObserve {
		spec, err := captureSpec(in.DisplayID, in.X, in.Y, in.W, in.H)
		if err != nil {
			out.ObserveError = err.Error()
		} else {
			data, mime, w, h, err := encodeCapture(d, spec, in.MaxWidth, in.Format)
			if err != nil {
				out.ObserveError = err.Error()
			} else {
				region := mouseRegion(d, spec)
				out.Width = w
				out.Height = h
				out.OriginX = region.OriginX
				out.OriginY = region.OriginY
				out.MouseWidth = region.MouseWidth
				out.MouseHeight = region.MouseHeight
				out.Format = formatFromMIME(mime)
				out.WindowTitle = d.WindowTitle()
				out.Displays = d.Displays()
				out.CmdCtrl = d.CmdCtrl()
				cx, cy := d.Location()
				out.CursorX, out.CursorY = cx, cy
				img = imageResult(data, mime)
			}
		}
	}
	return out, img, nil
}

func checkDesktopOp(name string) error {
	if name == "" {
		return fmt.Errorf("op is required")
	}
	if name == "KillProcess" {
		return fmt.Errorf("KillProcess is not allowed in DesktopOps")
	}
	if _, ok := desktopOpSet[name]; !ok {
		return fmt.Errorf("unknown op %q", name)
	}
	return nil
}

func runOneDesktopOp(d desktop.Driver, op desktopOp) desktopOpResult {
	item := desktopOpResult{Op: op.Op}
	resp, err := executeDesktopOp(d, op)
	if err != nil {
		item.Error = err.Error()
		return item
	}
	item.OK = true
	item.Response = resp
	if pt, ok := resp.(pointOutput); ok {
		item.X, item.Y = &pt.X, &pt.Y
	}
	return item
}

func executeDesktopOp(d desktop.Driver, op desktopOp) (any, error) {
	switch op.Op {
	case "MouseMove":
		x, y, err := requireXY(op.X, op.Y)
		if err != nil {
			return nil, err
		}
		return doMove(d, mouseMoveInput{
			X: x, Y: y, Relative: op.Relative, Smooth: op.Smooth,
			DisplayID: op.DisplayID, Low: op.Low, High: op.High, Delay: op.Delay,
		})
	case "MouseClick":
		return doClick(d, mouseClickInput{
			Button: op.Button, Double: op.Double, Count: op.Count, X: op.X, Y: op.Y,
		})
	case "MouseToggle":
		if err := doToggle(d, mouseToggleInput{Button: op.Button, State: op.State}); err != nil {
			return nil, err
		}
		return okOutput{OK: true}, nil
	case "MouseScroll":
		sx, sy := 0, 0
		if op.X != nil {
			sx = *op.X
		}
		if op.Y != nil {
			sy = *op.Y
		}
		if err := doScroll(d, mouseScrollInput{
			X: sx, Y: sy, Dir: op.Dir, Smooth: op.Smooth, Steps: op.Steps, Sleep: op.Sleep,
		}); err != nil {
			return nil, err
		}
		return okOutput{OK: true}, nil
	case "MouseDrag":
		x, y, err := requireXY(op.X, op.Y)
		if err != nil {
			return nil, err
		}
		return doDrag(d, mouseDragInput{X: x, Y: y})
	case "MouseSelect":
		x, y, err := requireXY(op.X, op.Y)
		if err != nil {
			return nil, err
		}
		toX, toY, err := requirePair("toX", "toY", op.ToX, op.ToY)
		if err != nil {
			return nil, err
		}
		return doSelect(d, mouseSelectInput{
			X: x, Y: y, ToX: toX, ToY: toY, Smooth: op.Smooth,
			DisplayID: op.DisplayID, Low: op.Low, High: op.High, Delay: op.Delay,
		})
	case "MousePath":
		return doPath(d, mousePathInput{
			Points: op.Points, Smooth: op.Smooth, Relative: op.Relative, Hold: op.Hold,
			DisplayID: op.DisplayID, Low: op.Low, High: op.High, Delay: op.Delay,
		})
	case "Type":
		if err := doType(d, typeInput{Text: op.Text, DelayMs: op.DelayMs, PID: op.PID}); err != nil {
			return nil, err
		}
		return okOutput{OK: true}, nil
	case "KeyTap":
		if err := doKeyTap(d, keyTapInput{Key: op.Key, Modifiers: op.Modifiers, PID: op.PID}); err != nil {
			return nil, err
		}
		return okOutput{OK: true}, nil
	case "KeyToggle":
		if err := doKeyToggle(d, keyToggleInput{Key: op.Key, State: op.State, Modifiers: op.Modifiers, PID: op.PID}); err != nil {
			return nil, err
		}
		return okOutput{OK: true}, nil
	case "KeyPress":
		if err := doKeyPress(d, keyPressInput{Key: op.Key, Modifiers: op.Modifiers, PID: op.PID}); err != nil {
			return nil, err
		}
		return okOutput{OK: true}, nil
	case "ClipboardWrite":
		if err := doClipboardWrite(d, clipboardText{Text: op.Text}); err != nil {
			return nil, err
		}
		return okOutput{OK: true}, nil
	case "ClipboardPaste":
		if err := doClipboardPaste(d, clipboardText{Text: op.Text}); err != nil {
			return nil, err
		}
		return okOutput{OK: true}, nil
	case "Sleep":
		doSleep(d, sleepInput{Ms: op.Ms})
		return okOutput{OK: true}, nil
	case "SetDelay":
		return doSetDelay(d, delayInput{KeySleep: op.KeySleep, MouseSleep: op.MouseSleep}), nil
	case "GetWindowTitle":
		return titleOutput{Title: d.WindowTitle()}, nil
	case "FocusWindow":
		return doFocusWindow(d, focusInput{Name: op.Name, PID: op.PID})
	case "GetWindowBounds":
		return doWindowBounds(d, boundsInput{PID: op.PID})
	case "MinWindow":
		if err := doMinWindow(d, pidInput{PID: op.PID}); err != nil {
			return nil, err
		}
		return okOutput{OK: true}, nil
	case "MaxWindow":
		if err := doMaxWindow(d, pidInput{PID: op.PID}); err != nil {
			return nil, err
		}
		return okOutput{OK: true}, nil
	case "CloseWindow":
		if err := doCloseWindow(d, pidInput{PID: op.PID}); err != nil {
			return nil, err
		}
		return okOutput{OK: true}, nil
	default:
		return nil, fmt.Errorf("unknown op %q", op.Op)
	}
}

func requireXY(x, y *int) (int, int, error) {
	return requirePair("x", "y", x, y)
}

func requirePair(xn, yn string, x, y *int) (int, int, error) {
	if x == nil || y == nil {
		return 0, 0, fmt.Errorf("%s and %s are required", xn, yn)
	}
	return *x, *y, nil
}
