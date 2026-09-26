---
name: observe-then-act
description: Computer-use loop for the robotgo MCP server. Use when driving a GUI with screenshots, mouse, and keyboard.
---

# Observe then act

This MCP server controls the machine it runs on. Treat every click and keypress as acting on a live desktop.

## Loop

1. Call `ObserveDesktop` once. Read the image plus `originX`, `originY`, `mouseWidth`, `mouseHeight`, `width`, `height`, `cursorX`, `cursorY`, `displays`, `windowTitle`, and `cmdCtrl`.
2. Map image pixels to mouse space before acting:

   `mouseX = originX + ix * mouseWidth / width`

   `mouseY = originY + iy * mouseHeight / height`

3. Decide the next **sequence**. Do not invent coordinates. Aim at control centers (`robotgo://skills/aim-center`).
4. If the next actions do not need a new screenshot (path, click, type), send them as **one** `DesktopOps` call. `MousePath` is a valid op (`smooth: false` for geometric traces).
5. Verify with `observeAfter` (optionally crop with `x,y,w,h`) or one `ObserveDesktop`. Confirm the UI changed before the next decision.

On a missed click, inset 40–80px into the box and click again. Do not crop-loop.

Observe again only when the next decision needs new pixels (menu opened, scroll, unknown dialog).

If capture fails, stop. Read the `robotgo-backend-limits` prompt. Linux libei cannot screenshot.

## Coordinate space

`width`/`height` are encoded image pixels (Retina capture may be 2×, then scaled to `maxWidth`). `mouseWidth`/`mouseHeight` and `originX`/`originY` are the `MouseMove` / `MouseClick` space. Multi-monitor layouts may use negative `originX`/`originY`. After a move, trust `cursorX`/`cursorY` from the next observe.

## Do not

- Type passwords into logs.
- Call `KillProcess` unless the user asked to kill that pid (it elicits). `KillProcess` is not valid inside `DesktopOps`.
- Skip the post-batch screenshot.
