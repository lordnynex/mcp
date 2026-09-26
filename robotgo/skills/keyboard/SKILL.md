---
name: keyboard
description: robotgo key names, modifiers, and when to Type versus KeyTap.
---

# Keyboard

## Type vs tap

- `Type` sends UTF-8 text (letters, spaces, Unicode). Use it for field content.
- `KeyTap` sends a named key plus modifiers. Use it for Enter, Tab, Escape, arrows, function keys, and shortcuts.
- `KeyToggle` holds or releases a key (`down` / `up`).
- `KeyPress` is down, a short delay, then up.

Never type modifier glyphs. Get the platform accelerator from `ObserveDesktop.cmdCtrl` (`cmd` on macOS, `ctrl` elsewhere).

## Named keys

Letters `a`–`z` and digits `0`–`9` are valid keys. Also: `enter`, `tab`, `esc`/`escape`, `space`, `backspace`, `delete`, `up`, `down`, `left`, `right`, `home`, `end`, `pageup`, `pagedown`, `f1`–`f24`, media keys (`audio_play`, …), numpad (`num0`–`num9`, `num_enter`).

Modifiers: `cmd`, `cmdl`, `cmdr`, `command`, `alt`, `altl`, `altr`, `ctrl`, `ctrll`, `ctrlr`, `control`, `shift`, `shiftl`, `shiftr`, `right_shift`.

Do not invent names. Completions list official keys.

## Examples

- Paste shortcut: `KeyTap` key `v`, modifiers `[cmdCtrl]`.
- New tab: `KeyTap` key `t`, modifiers `[cmdCtrl]`.
- Confirm a dialog: `KeyTap` key `enter`.
