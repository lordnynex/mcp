package desktop

import (
	"fmt"
	"image"
	"image/color"
	"sync"
	"time"
)

// Call is one recorded Fake method invocation.
type Call struct {
	Name string
	Args []any
}

// Fake is an in-memory Driver for tests. It never touches the real desktop.
type Fake struct {
	mu sync.Mutex

	Calls []Call

	Image      image.Image
	CursorX    int
	CursorY    int
	Screen     []Display
	Title      string
	Procs      []Process
	Clipboard  string
	CmdCtrlKey string
	KeySleep   int
	MouseSleep int
	CaptureErr error
	KillErr    error
	Killed     []int
}

// NewFake returns a Fake with a 2x2 PNG-friendly image and one display.
func NewFake() *Fake {
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	for y := range 4 {
		for x := range 4 {
			img.Set(x, y, color.RGBA{R: 0x11, G: 0x22, B: 0x33, A: 0xff})
		}
	}
	return &Fake{
		Image:      img,
		CmdCtrlKey: "cmd",
		Screen: []Display{{
			ID: 0, X: 0, Y: 0, W: 1920, H: 1080,
			ScaleWidth: 1920, ScaleHeight: 1080, Main: true,
		}},
		Title: "Fake Window",
		Procs: []Process{
			{PID: 1, Name: "init"},
			{PID: 42, Name: "FakeApp"},
		},
	}
}

func (f *Fake) record(name string, args ...any) {
	f.Calls = append(f.Calls, Call{Name: name, Args: args})
}

func (f *Fake) Move(x, y, displayID int, relative, smooth bool, low, high float64, delay int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.record("Move", x, y, displayID, relative, smooth, low, high, delay)
	if relative {
		f.CursorX += x
		f.CursorY += y
		return nil
	}
	f.CursorX = x
	f.CursorY = y
	return nil
}

func (f *Fake) Click(button string, double bool, count int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.record("Click", button, double, count)
	return nil
}

func (f *Fake) Toggle(button, state string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.record("Toggle", button, state)
	return nil
}

func (f *Fake) Scroll(x, y int, dir string, smooth bool, steps, sleepMs int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.record("Scroll", x, y, dir, smooth, steps, sleepMs)
	return nil
}

func (f *Fake) Drag(x, y int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.record("Drag", x, y)
	f.CursorX = x
	f.CursorY = y
	return nil
}

func (f *Fake) Location() (int, int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.CursorX, f.CursorY
}

func (f *Fake) Type(text string, delayMs, pid int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.record("Type", text, delayMs, pid)
	return nil
}

func (f *Fake) KeyTap(key string, modifiers []string, pid int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.record("KeyTap", key, modifiers, pid)
	return nil
}

func (f *Fake) KeyToggle(key, state string, modifiers []string, pid int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.record("KeyToggle", key, state, modifiers, pid)
	return nil
}

func (f *Fake) KeyPress(key string, modifiers []string, pid int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.record("KeyPress", key, modifiers, pid)
	return nil
}

func (f *Fake) CmdCtrl() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.CmdCtrlKey == "" {
		return "ctrl"
	}
	return f.CmdCtrlKey
}

func (f *Fake) Capture(spec CaptureSpec) (image.Image, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.record("Capture", spec)
	if f.CaptureErr != nil {
		return nil, f.CaptureErr
	}
	return f.Image, nil
}

func (f *Fake) Displays() []Display {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]Display, len(f.Screen))
	copy(out, f.Screen)
	return out
}

func (f *Fake) ScreenSize() (int, int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.Screen) == 0 {
		return 0, 0
	}
	return f.Screen[0].W, f.Screen[0].H
}

func (f *Fake) PixelColor(x, y, displayID int, atCursor bool) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.record("PixelColor", x, y, displayID, atCursor)
	return "112233", nil
}

func (f *Fake) ClipboardRead() (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.record("ClipboardRead")
	return f.Clipboard, nil
}

func (f *Fake) ClipboardWrite(text string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.record("ClipboardWrite", text)
	f.Clipboard = text
	return nil
}

func (f *Fake) ClipboardPaste(text string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.record("ClipboardPaste", text)
	f.Clipboard = text
	return nil
}

func (f *Fake) WindowTitle() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.Title
}

func (f *Fake) FocusWindow(name string, pid int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.record("FocusWindow", name, pid)
	if name != "" {
		f.Title = name
	}
	return nil
}

func (f *Fake) WindowBounds(pid int) (int, int, int, int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.record("WindowBounds", pid)
	return 0, 0, 800, 600, nil
}

func (f *Fake) MinWindow(pid int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.record("MinWindow", pid)
	return nil
}

func (f *Fake) MaxWindow(pid int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.record("MaxWindow", pid)
	return nil
}

func (f *Fake) CloseWindow(pid int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.record("CloseWindow", pid)
	return nil
}

func (f *Fake) Processes() ([]Process, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]Process, len(f.Procs))
	copy(out, f.Procs)
	return out, nil
}

func (f *Fake) FindProcess(name string) ([]Process, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.record("FindProcess", name)
	var out []Process
	for _, p := range f.Procs {
		if p.Name == name {
			out = append(out, p)
		}
	}
	return out, nil
}

func (f *Fake) Kill(pid int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.record("Kill", pid)
	if f.KillErr != nil {
		return f.KillErr
	}
	f.Killed = append(f.Killed, pid)
	return nil
}

func (f *Fake) SetDelay(keySleep, mouseSleep int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.record("SetDelay", keySleep, mouseSleep)
	f.KeySleep = keySleep
	f.MouseSleep = mouseSleep
}

func (f *Fake) Sleep(ms int) {
	if ms > 0 {
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
}

// Called reports whether a method was recorded.
func (f *Fake) Called(name string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, c := range f.Calls {
		if c.Name == name {
			return true
		}
	}
	return false
}

// Last returns the most recent call with name, or an error if none.
func (f *Fake) Last(name string) (Call, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for i := len(f.Calls) - 1; i >= 0; i-- {
		if f.Calls[i].Name == name {
			return f.Calls[i], nil
		}
	}
	return Call{}, fmt.Errorf("no call named %q", name)
}
