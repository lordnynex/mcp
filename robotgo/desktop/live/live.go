//go:build purego || mac || x11 || wayland || libei || win || robotgocgo

package live

import (
	"fmt"
	"image"
	"strings"
	"sync"

	rg "github.com/go-vgo/robotgo"
	"github.com/lordnynex/mcp/robotgo/desktop"
)

const captureHint = "screen capture is not supported by this backend (libei is input-only); rebuild with -tags \"purego,x11\" or -tags \"purego,wayland\" on a compositor that supports zwlr_screencopy_v1"

// New returns a Driver that calls go-vgo/robotgo. All calls are serialized
// because robotgo keeps DisplayID, KeySleep, and MouseSleep as process globals.
func New() desktop.Driver {
	return &driver{}
}

type driver struct {
	mu sync.Mutex
}

func (d *driver) Move(x, y, displayID int, relative, smooth bool, low, high float64, delay int) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if relative && smooth {
		args := smoothArgs(low, high, delay)
		if len(args) > 0 {
			rg.MoveSmoothRelative(x, y, args...)
		} else {
			rg.MoveSmoothRelative(x, y)
		}
		return nil
	}
	if relative {
		rg.MoveRelative(x, y)
		return nil
	}
	if smooth {
		args := smoothArgs(low, high, delay)
		if displayID >= 0 {
			// MoveSmooth has no displayId; set the global for backends that read it.
			prev := rg.DisplayID
			rg.DisplayID = displayID
			defer func() { rg.DisplayID = prev }()
		}
		if len(args) > 0 {
			rg.MoveSmooth(x, y, args...)
		} else {
			rg.MoveSmooth(x, y)
		}
		return nil
	}
	if displayID >= 0 {
		rg.Move(x, y, displayID)
		return nil
	}
	rg.Move(x, y)
	return nil
}

func smoothArgs(low, high float64, delay int) []any {
	if low == 0 && high == 0 && delay == 0 {
		return nil
	}
	out := []any{}
	if low != 0 || high != 0 {
		out = append(out, low, high)
	}
	if delay > 0 {
		out = append(out, delay)
	}
	return out
}

func (d *driver) Click(button string, double bool, count int) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if button == "" {
		button = "left"
	}
	if double {
		return rg.Click(button, true)
	}
	n := count
	if n < 1 {
		n = 1
	}
	for range n {
		if err := rg.Click(button); err != nil {
			return err
		}
	}
	return nil
}

func (d *driver) Toggle(button, state string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if button == "" {
		button = "left"
	}
	if state == "" || state == "down" {
		return rg.Toggle(button)
	}
	return rg.Toggle(button, "up")
}

func (d *driver) Scroll(x, y int, dir string, smooth bool, steps, sleepMs int) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if dir != "" {
		if smooth {
			args := []int{}
			if steps > 0 {
				args = append(args, steps)
			}
			if sleepMs > 0 {
				args = append(args, sleepMs)
			}
			amount := x
			if amount == 0 {
				amount = y
			}
			if len(args) > 0 {
				rg.ScrollSmooth(amount, args...)
			} else {
				rg.ScrollSmooth(amount)
			}
			return nil
		}
		rg.ScrollDir(max(x, y), dir)
		return nil
	}
	rg.Scroll(x, y)
	return nil
}

func (d *driver) Drag(x, y int) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	rg.DragSmooth(x, y)
	return nil
}

func (d *driver) Location() (int, int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return rg.Location()
}

func (d *driver) Type(text string, delayMs, pid int) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	switch {
	case pid > 0 && delayMs > 0:
		rg.Type(text, pid, delayMs)
	case pid > 0:
		rg.Type(text, pid)
	case delayMs > 0:
		rg.Type(text, 0, delayMs)
	default:
		rg.Type(text)
	}
	return nil
}

func (d *driver) KeyTap(key string, modifiers []string, pid int) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return rg.KeyTap(key, keyArgs(modifiers, pid)...)
}

func (d *driver) KeyToggle(key, state string, modifiers []string, pid int) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	args := make([]any, 0, len(modifiers)+2)
	if state != "" && state != "down" {
		args = append(args, state)
	}
	for _, m := range modifiers {
		args = append(args, m)
	}
	if pid > 0 {
		args = append(args, pid)
	}
	if len(args) == 0 {
		return rg.KeyToggle(key)
	}
	return rg.KeyToggle(key, args...)
}

func (d *driver) KeyPress(key string, modifiers []string, pid int) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return rg.KeyPress(key, keyArgs(modifiers, pid)...)
}

func (d *driver) CmdCtrl() string {
	return rg.CmdCtrl()
}

func keyArgs(modifiers []string, pid int) []any {
	if pid > 0 && len(modifiers) > 0 {
		return append([]any{pid}, toAny(modifiers)...)
	}
	if pid > 0 {
		return []any{pid}
	}
	if len(modifiers) == 1 {
		return []any{modifiers[0]}
	}
	if len(modifiers) > 1 {
		return []any{modifiers}
	}
	return nil
}

func toAny(ss []string) []any {
	out := make([]any, len(ss))
	for i, s := range ss {
		out[i] = s
	}
	return out
}

func (d *driver) Capture(spec desktop.CaptureSpec) (image.Image, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	prev := rg.DisplayID
	if spec.DisplayID >= 0 {
		rg.DisplayID = spec.DisplayID
	}
	defer func() { rg.DisplayID = prev }()

	var (
		img image.Image
		err error
	)
	if spec.HasRegion {
		img, err = rg.CaptureImg(spec.X, spec.Y, spec.W, spec.H)
	} else {
		img, err = rg.CaptureImg()
	}
	if err != nil {
		if isNotSupported(err) {
			return nil, fmt.Errorf("%w: %s", err, captureHint)
		}
		return nil, err
	}
	return img, nil
}

func (d *driver) Displays() []desktop.Display {
	d.mu.Lock()
	defer d.mu.Unlock()
	n := rg.DisplaysNum()
	main := rg.GetMainId()
	out := make([]desktop.Display, 0, n)
	for i := range n {
		x, y, w, h := rg.GetDisplayBounds(i)
		sw, sh := rg.GetScaleSize(i)
		out = append(out, desktop.Display{
			ID:          i,
			X:           x,
			Y:           y,
			W:           w,
			H:           h,
			ScaleWidth:  sw,
			ScaleHeight: sh,
			Main:        i == main,
		})
	}
	return out
}

func (d *driver) ScreenSize() (int, int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return rg.GetScreenSize()
}

func (d *driver) PixelColor(x, y, displayID int, atCursor bool) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if atCursor {
		x, y = rg.Location()
	}
	if displayID >= 0 {
		return rg.GetPixelColor(x, y, displayID), nil
	}
	return rg.GetPixelColor(x, y), nil
}

func (d *driver) ClipboardRead() (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return rg.ReadAll()
}

func (d *driver) ClipboardWrite(text string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return rg.WriteAll(text)
}

func (d *driver) ClipboardPaste(text string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return rg.Paste(text)
}

func (d *driver) WindowTitle() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	return rg.GetTitle()
}

func (d *driver) FocusWindow(name string, pid int) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if name != "" {
		if err := rg.ActiveName(name); err != nil {
			return windowErr(err)
		}
		return nil
	}
	if pid > 0 {
		found, err := rg.FindName(pid)
		if err != nil {
			return err
		}
		if found == "" {
			return fmt.Errorf("no process name for pid %d", pid)
		}
		if err := rg.ActiveName(found); err != nil {
			return windowErr(fmt.Errorf("focus by pid uses ActiveName on the pure-Go backend: %w", err))
		}
		return nil
	}
	return fmt.Errorf("name or pid is required")
}

func (d *driver) WindowBounds(pid int) (int, int, int, int, error) {
	return 0, 0, 0, 0, fmt.Errorf("window bounds are not available on the pure-Go robotgo backend")
}

func (d *driver) MinWindow(pid int) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	rg.MinWindow(pid)
	return nil
}

func (d *driver) MaxWindow(pid int) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	rg.MaxWindow(pid)
	return nil
}

func (d *driver) CloseWindow(pid int) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if pid > 0 {
		rg.CloseWindow(pid)
		return nil
	}
	rg.CloseWindow()
	return nil
}

func (d *driver) Processes() ([]desktop.Process, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	list, err := rg.Process()
	if err != nil {
		return nil, err
	}
	out := make([]desktop.Process, 0, len(list))
	for _, p := range list {
		out = append(out, desktop.Process{PID: p.Pid, Name: p.Name})
	}
	return out, nil
}

func (d *driver) FindProcess(name string) ([]desktop.Process, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	ids, err := rg.FindIds(name)
	if err != nil {
		return nil, err
	}
	out := make([]desktop.Process, 0, len(ids))
	for _, id := range ids {
		n, nerr := rg.FindName(id)
		if nerr != nil {
			n = name
		}
		out = append(out, desktop.Process{PID: id, Name: n})
	}
	return out, nil
}

func (d *driver) Kill(pid int) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return rg.Kill(pid)
}

func (d *driver) SetDelay(keySleep, mouseSleep int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	rg.KeySleep = keySleep
	rg.MouseSleep = mouseSleep
}

func (d *driver) Sleep(ms int) {
	rg.MilliSleep(ms)
}

func windowErr(err error) error {
	if err == nil {
		return nil
	}
	if isNotSupported(err) {
		return fmt.Errorf("%w: window management is not implemented on this backend (macOS -tags purego and Linux libei)", err)
	}
	return err
}

func isNotSupported(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "not supported")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
