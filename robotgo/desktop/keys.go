package desktop

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Buttons are robotgo mouse button names.
var Buttons = []string{
	"left", "center", "right",
	"wheelDown", "wheelUp", "wheelLeft", "wheelRight",
}

// ScrollDirs are robotgo ScrollDir values.
var ScrollDirs = []string{"up", "down", "left", "right"}

// Modifiers are valid KeyTap / KeyToggle modifier names.
var Modifiers = []string{
	"cmd", "cmdl", "cmdr", "command",
	"alt", "altl", "altr",
	"ctrl", "ctrll", "ctrlr", "control",
	"shift", "shiftl", "shiftr", "right_shift",
}

// Keys are named robotgo keys (not including single letters/digits).
var Keys = []string{
	"backspace", "delete", "enter", "tab", "esc", "escape",
	"up", "down", "left", "right",
	"home", "end", "pageup", "pagedown",
	"insert", "space", "capslock", "scroll_lock", "pause_break",
	"f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8", "f9", "f10",
	"f11", "f12", "f13", "f14", "f15", "f16", "f17", "f18", "f19", "f20",
	"f21", "f22", "f23", "f24",
	"cmd", "cmdl", "cmdr", "command",
	"alt", "altl", "altr",
	"ctrl", "ctrll", "ctrlr", "control",
	"shift", "shiftl", "shiftr", "right_shift",
	"print", "printscreen", "menu",
	"audio_mute", "audio_vol_down", "audio_vol_up",
	"audio_play", "audio_stop", "audio_pause", "audio_prev", "audio_next",
	"audio_rewind", "audio_forward", "audio_repeat", "audio_random",
	"num0", "num1", "num2", "num3", "num4", "num5", "num6", "num7", "num8", "num9",
	"num_lock", "num.", "num+", "num-", "num*", "num/", "num_clear", "num_enter", "num_equal",
	"lights_mon_up", "lights_mon_down", "lights_kbd_toggle", "lights_kbd_up", "lights_kbd_down",
}

var (
	buttonSet   = setOf(Buttons)
	scrollSet   = setOf(ScrollDirs)
	modifierSet = setOf(Modifiers)
	keySet      = setOf(Keys)
)

func setOf(vals []string) map[string]struct{} {
	m := make(map[string]struct{}, len(vals))
	for _, v := range vals {
		m[v] = struct{}{}
	}
	return m
}

// ValidButton reports whether name is a robotgo mouse button.
func ValidButton(name string) bool {
	_, ok := buttonSet[name]
	return ok
}

// ValidScrollDir reports whether name is a robotgo scroll direction.
func ValidScrollDir(name string) bool {
	_, ok := scrollSet[name]
	return ok
}

// ValidModifier reports whether name is a robotgo modifier.
func ValidModifier(name string) bool {
	_, ok := modifierSet[name]
	return ok
}

// ValidKey reports whether name is a robotgo key or a single letter/digit.
func ValidKey(name string) bool {
	if name == "" {
		return false
	}
	if _, ok := keySet[name]; ok {
		return true
	}
	r, size := utf8.DecodeRuneInString(name)
	if size != len(name) || r == utf8.RuneError {
		return false
	}
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

// CheckKey returns an error if name is not a valid robotgo key.
func CheckKey(name string) error {
	if ValidKey(name) {
		return nil
	}
	return fmt.Errorf("unknown key %q", name)
}

// CheckButton returns an error if name is not a valid mouse button.
func CheckButton(name string) error {
	if name == "" || ValidButton(name) {
		return nil
	}
	return fmt.Errorf("unknown mouse button %q", name)
}

// CheckModifiers returns an error if any modifier is unknown.
func CheckModifiers(mods []string) error {
	for _, m := range mods {
		if !ValidModifier(m) {
			return fmt.Errorf("unknown modifier %q", m)
		}
	}
	return nil
}

// CheckScrollDir returns an error if name is not a valid scroll direction.
func CheckScrollDir(name string) error {
	if name == "" || ValidScrollDir(name) {
		return nil
	}
	return fmt.Errorf("unknown scroll direction %q", name)
}

// NormalizeFormat returns png or jpeg.
func NormalizeFormat(format string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "", "png":
		return "png", nil
	case "jpeg", "jpg":
		return "jpeg", nil
	default:
		return "", fmt.Errorf("unknown image format %q (png or jpeg)", format)
	}
}
