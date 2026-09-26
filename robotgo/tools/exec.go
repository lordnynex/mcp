package tools

import (
	"fmt"

	"github.com/lordnynex/mcp/robotgo/desktop"
)

func doMove(d desktop.Driver, in mouseMoveInput) (pointOutput, error) {
	displayID := -1
	if in.DisplayID != nil {
		displayID = *in.DisplayID
	}
	if err := d.Move(in.X, in.Y, displayID, in.Relative, in.Smooth, in.Low, in.High, in.Delay); err != nil {
		return pointOutput{}, err
	}
	x, y := d.Location()
	return pointOutput{X: x, Y: y}, nil
}

func doClick(d desktop.Driver, in mouseClickInput) (pointOutput, error) {
	if err := desktop.CheckButton(in.Button); err != nil {
		return pointOutput{}, err
	}
	if (in.X != nil) != (in.Y != nil) {
		return pointOutput{}, fmt.Errorf("x and y must be set together")
	}
	if in.X != nil {
		if err := d.Move(*in.X, *in.Y, -1, false, false, 0, 0, 0); err != nil {
			return pointOutput{}, err
		}
	}
	if err := d.Click(in.Button, in.Double, in.Count); err != nil {
		return pointOutput{}, err
	}
	x, y := d.Location()
	return pointOutput{X: x, Y: y}, nil
}

func doToggle(d desktop.Driver, in mouseToggleInput) error {
	if err := desktop.CheckButton(in.Button); err != nil {
		return err
	}
	state := in.State
	if state == "" {
		state = "down"
	}
	if state != "down" && state != "up" {
		return fmt.Errorf("state must be down or up")
	}
	return d.Toggle(in.Button, state)
}

func doScroll(d desktop.Driver, in mouseScrollInput) error {
	if err := desktop.CheckScrollDir(in.Dir); err != nil {
		return err
	}
	return d.Scroll(in.X, in.Y, in.Dir, in.Smooth, in.Steps, in.Sleep)
}

func doDrag(d desktop.Driver, in mouseDragInput) (pointOutput, error) {
	if err := d.Drag(in.X, in.Y); err != nil {
		return pointOutput{}, err
	}
	x, y := d.Location()
	return pointOutput{X: x, Y: y}, nil
}

func doSelect(d desktop.Driver, in mouseSelectInput) (pointOutput, error) {
	displayID := -1
	if in.DisplayID != nil {
		displayID = *in.DisplayID
	}
	if err := d.Move(in.X, in.Y, displayID, false, in.Smooth, in.Low, in.High, in.Delay); err != nil {
		return pointOutput{}, err
	}
	if err := d.Toggle("left", "down"); err != nil {
		return pointOutput{}, err
	}
	if err := d.Move(in.ToX, in.ToY, displayID, false, in.Smooth, in.Low, in.High, in.Delay); err != nil {
		_ = d.Toggle("left", "up")
		return pointOutput{}, err
	}
	if err := d.Toggle("left", "up"); err != nil {
		return pointOutput{}, err
	}
	x, y := d.Location()
	return pointOutput{X: x, Y: y}, nil
}

func doType(d desktop.Driver, in typeInput) error {
	if in.Text == "" {
		return fmt.Errorf("text is required")
	}
	return d.Type(in.Text, in.DelayMs, in.PID)
}

func doKeyTap(d desktop.Driver, in keyTapInput) error {
	if err := desktop.CheckKey(in.Key); err != nil {
		return err
	}
	if err := desktop.CheckModifiers(in.Modifiers); err != nil {
		return err
	}
	return d.KeyTap(in.Key, in.Modifiers, in.PID)
}

func doKeyToggle(d desktop.Driver, in keyToggleInput) error {
	if err := desktop.CheckKey(in.Key); err != nil {
		return err
	}
	if err := desktop.CheckModifiers(in.Modifiers); err != nil {
		return err
	}
	state := in.State
	if state == "" {
		state = "down"
	}
	if state != "down" && state != "up" {
		return fmt.Errorf("state must be down or up")
	}
	return d.KeyToggle(in.Key, state, in.Modifiers, in.PID)
}

func doKeyPress(d desktop.Driver, in keyPressInput) error {
	if err := desktop.CheckKey(in.Key); err != nil {
		return err
	}
	if err := desktop.CheckModifiers(in.Modifiers); err != nil {
		return err
	}
	return d.KeyPress(in.Key, in.Modifiers, in.PID)
}

func doClipboardWrite(d desktop.Driver, in clipboardText) error {
	return d.ClipboardWrite(in.Text)
}

func doClipboardPaste(d desktop.Driver, in clipboardText) error {
	if in.Text == "" {
		return fmt.Errorf("text is required")
	}
	return d.ClipboardPaste(in.Text)
}

func doSetDelay(d desktop.Driver, in delayInput) delayOutput {
	key, mouse := 0, 0
	if in.KeySleep != nil {
		key = *in.KeySleep
	}
	if in.MouseSleep != nil {
		mouse = *in.MouseSleep
	}
	d.SetDelay(key, mouse)
	return delayOutput{KeySleep: key, MouseSleep: mouse}
}

func doSleep(d desktop.Driver, in sleepInput) {
	if in.Ms < 0 {
		in.Ms = 0
	}
	d.Sleep(in.Ms)
}

func doFocusWindow(d desktop.Driver, in focusInput) (titleOutput, error) {
	if in.Name == "" && in.PID <= 0 {
		return titleOutput{}, fmt.Errorf("name or pid is required")
	}
	if err := d.FocusWindow(in.Name, in.PID); err != nil {
		return titleOutput{}, err
	}
	return titleOutput{Title: d.WindowTitle()}, nil
}

func doWindowBounds(d desktop.Driver, in boundsInput) (boundsOutput, error) {
	if in.PID <= 0 {
		return boundsOutput{}, fmt.Errorf("pid is required")
	}
	x, y, w, h, err := d.WindowBounds(in.PID)
	if err != nil {
		return boundsOutput{}, err
	}
	return boundsOutput{X: x, Y: y, W: w, H: h}, nil
}

func doMinWindow(d desktop.Driver, in pidInput) error {
	return d.MinWindow(in.PID)
}

func doMaxWindow(d desktop.Driver, in pidInput) error {
	return d.MaxWindow(in.PID)
}

func doCloseWindow(d desktop.Driver, in pidInput) error {
	return d.CloseWindow(in.PID)
}

func locationOf(d desktop.Driver) pointOutput {
	x, y := d.Location()
	return pointOutput{X: x, Y: y}
}
