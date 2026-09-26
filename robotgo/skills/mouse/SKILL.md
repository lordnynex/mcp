---
name: mouse
description: robotgo mouse tools, buttons, scrolling, bulk paths, and multi-monitor coordinates.
---

# Mouse

Call `ObserveDesktop` before planning a sequence. Map image pixels with `originX`/`mouseWidth`/`width` (see `robotgo://skills/observe-then-act`). Secondary monitors can start at negative coordinates.

Aim `x`,`y` at the **geometric center** of the control, not the first glyph or an edge. See `robotgo://skills/aim-center`.

## Tools

- `MouseMove` — absolute by default; `relative` for a delta; `smooth` for `MoveSmooth`.
- `MouseClick` — `left` (default), `center`, `right`, `wheelUp`, `wheelDown`, `wheelLeft`, `wheelRight`. Optional `x`/`y` moves first (center of the hit target). `double` or `count`.
- `MouseToggle` — `down` / `up` for click-and-hold.
- `MouseScroll` — `dir` (`up`/`down`/`left`/`right`) or `x`/`y` ticks; `smooth` for `ScrollSmooth`.
- `MouseDrag` — left-button `DragSmooth` to `x`,`y`. Cursor must already be at the start.
- `MouseSelect` — drag-select from `x`,`y` to `toX`,`toY` (button down, move, up). Use this to highlight text.
- `MousePath` — array of points with shared `smooth` / `relative`. Also a `DesktopOps` op. `hold` keeps the left button down after the first point.
- `MouseLocation` — current cursor.
- `DesktopOps` — mix mouse, keyboard, clipboard, sleep, window ops, and `MousePath` in one serial call.

Prefer one `DesktopOps` when the path or click-then-type sequence is already known. Do not issue one `MouseMove` per vertex.

## Smooth

Geometric traces (star, zigzag, a known polyline) use `smooth: false`. Use `smooth: true` only for the last approach to a UI control.

## Select text

- Range: `MouseSelect` (or `MouseToggle` down → `MouseMove` → up). Not a click on the first character.
- Polyline: `MousePath` with `hold`.
- Word / line: `MouseClick` `double` or `count` 2 / 3 at the text center.
- Shift-click: `DesktopOps` with `KeyToggle` shift down, `MouseClick` at the other end, `KeyToggle` up.

## After a planned sequence

Call `ObserveDesktop` once or set `DesktopOps.observeAfter` (crop with `x,y,w,h` when verifying a field). If the cursor is not where you intended, inset 40–80px toward the box interior instead of repeating the same edge or starting a crop loop.
