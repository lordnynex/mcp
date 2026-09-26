package desktop

import (
	"fmt"
	"image"
	"sync"
)

// RebuildTags is the error text for binaries built without a live backend.
const RebuildTags = `this binary was built without a robotgo backend; rebuild with -tags purego (macOS) or -tags "purego,x11" / "purego,wayland" / "purego,libei" (Linux)`

// Display is one monitor's bounds and sizes.
type Display struct {
	ID          int  `json:"id"`
	X           int  `json:"x"`
	Y           int  `json:"y"`
	W           int  `json:"w"`
	H           int  `json:"h"`
	ScaleWidth  int  `json:"scaleWidth,omitempty"`
	ScaleHeight int  `json:"scaleHeight,omitempty"`
	Main        bool `json:"main,omitempty"`
}

// Process is a running OS process.
type Process struct {
	PID  int    `json:"pid"`
	Name string `json:"name"`
}

// CaptureSpec selects a screen region.
type CaptureSpec struct {
	DisplayID  int
	X, Y, W, H int
	HasRegion  bool
}

// Driver is the desktop automation surface used by MCP tools.
type Driver interface {
	Move(x, y int, displayID int, relative, smooth bool, low, high float64, delay int) error
	Click(button string, double bool, count int) error
	Toggle(button, state string) error
	Scroll(x, y int, dir string, smooth bool, steps, sleepMs int) error
	Drag(x, y int) error
	Location() (int, int)

	Type(text string, delayMs, pid int) error
	KeyTap(key string, modifiers []string, pid int) error
	KeyToggle(key, state string, modifiers []string, pid int) error
	KeyPress(key string, modifiers []string, pid int) error
	CmdCtrl() string

	Capture(spec CaptureSpec) (image.Image, error)
	Displays() []Display
	ScreenSize() (int, int)
	PixelColor(x, y, displayID int, atCursor bool) (string, error)

	ClipboardRead() (string, error)
	ClipboardWrite(text string) error
	ClipboardPaste(text string) error

	WindowTitle() string
	FocusWindow(name string, pid int) error
	WindowBounds(pid int) (x, y, w, h int, err error)
	MinWindow(pid int) error
	MaxWindow(pid int) error
	CloseWindow(pid int) error

	Processes() ([]Process, error)
	FindProcess(name string) ([]Process, error)
	Kill(pid int) error

	SetDelay(keySleep, mouseSleep int)
	Sleep(ms int)
}

var (
	driverMu sync.RWMutex
	current  Driver = Missing()
)

// SetDefault installs the process-wide driver used by New and Register.
func SetDefault(d Driver) {
	driverMu.Lock()
	defer driverMu.Unlock()
	if d == nil {
		current = Missing()
		return
	}
	current = d
}

// Current returns the process-wide driver.
func Current() Driver {
	driverMu.RLock()
	defer driverMu.RUnlock()
	return current
}

// Missing returns a driver that errors on every operation with RebuildTags.
func Missing() Driver {
	return missingDriver{}
}

type missingDriver struct{}

func (missingDriver) err() error { return fmt.Errorf("%s", RebuildTags) }

func (missingDriver) Move(int, int, int, bool, bool, float64, float64, int) error {
	return missingDriver{}.err()
}
func (missingDriver) Click(string, bool, int) error                 { return missingDriver{}.err() }
func (missingDriver) Toggle(string, string) error                   { return missingDriver{}.err() }
func (missingDriver) Scroll(int, int, string, bool, int, int) error { return missingDriver{}.err() }
func (missingDriver) Drag(int, int) error                           { return missingDriver{}.err() }
func (missingDriver) Location() (int, int)                          { return 0, 0 }
func (missingDriver) Type(string, int, int) error                   { return missingDriver{}.err() }
func (missingDriver) KeyTap(string, []string, int) error            { return missingDriver{}.err() }
func (missingDriver) KeyToggle(string, string, []string, int) error { return missingDriver{}.err() }
func (missingDriver) KeyPress(string, []string, int) error          { return missingDriver{}.err() }
func (missingDriver) CmdCtrl() string                               { return "ctrl" }
func (missingDriver) Capture(CaptureSpec) (image.Image, error)      { return nil, missingDriver{}.err() }
func (missingDriver) Displays() []Display                           { return nil }
func (missingDriver) ScreenSize() (int, int)                        { return 0, 0 }
func (missingDriver) PixelColor(int, int, int, bool) (string, error) {
	return "", missingDriver{}.err()
}
func (missingDriver) ClipboardRead() (string, error) { return "", missingDriver{}.err() }
func (missingDriver) ClipboardWrite(string) error    { return missingDriver{}.err() }
func (missingDriver) ClipboardPaste(string) error    { return missingDriver{}.err() }
func (missingDriver) WindowTitle() string            { return "" }
func (missingDriver) FocusWindow(string, int) error  { return missingDriver{}.err() }
func (missingDriver) WindowBounds(int) (int, int, int, int, error) {
	return 0, 0, 0, 0, missingDriver{}.err()
}
func (missingDriver) MinWindow(int) error           { return missingDriver{}.err() }
func (missingDriver) MaxWindow(int) error           { return missingDriver{}.err() }
func (missingDriver) CloseWindow(int) error         { return missingDriver{}.err() }
func (missingDriver) Processes() ([]Process, error) { return nil, missingDriver{}.err() }
func (missingDriver) FindProcess(string) ([]Process, error) {
	return nil, missingDriver{}.err()
}
func (missingDriver) Kill(int) error    { return missingDriver{}.err() }
func (missingDriver) SetDelay(int, int) {}
func (missingDriver) Sleep(int)         {}
