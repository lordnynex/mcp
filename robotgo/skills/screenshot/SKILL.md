---
name: screenshot
description: How to capture the robotgo desktop for vision: crop, max width, PNG vs JPEG.
---

# Screenshots

`ObserveDesktop` is the default. It returns `mcp.ImageContent` plus cursor, displays, window title, and image-to-mouse mapping.

`CaptureScreen` is image-only when metadata is not needed (it still returns the mapping fields).

## Mapping

`width` / `height` are encoded image pixels. `mouseWidth` / `mouseHeight` and `originX` / `originY` are the `MouseClick` space.

`mouseX = originX + ix * mouseWidth / width`

`mouseY = originY + iy * mouseHeight / height`

HiDPI captures are often 2× logical pixels, then scaled down to `maxWidth`. Do not treat image pixels as mouse coordinates.

## Size

Images are scaled to `--screenshot-max-width` (default 1280) so vision tokens stay usable. Override per call with `maxWidth`. Crop with `x`,`y`,`w`,`h` (all four required) when the target is a small region.

`format` is `png` (default, lossless) or `jpeg` (smaller). JPEG is fine for dense UIs; use PNG when you need crisp text.

## Failure

If capture returns not supported, the binary was built with Linux `libei` (GNOME/KDE input-only). Screenshots need `-tags "purego,x11"` or `-tags "purego,wayland"` on a wlroots compositor. macOS needs Screen Recording permission.

## Workflow

1. Full-display `ObserveDesktop` to orient and read the mapping fields.
2. One `DesktopOps` for the planned sequence (`MousePath` may be an op; `smooth: false` for traces).
3. One verify: `observeAfter` with an optional crop, or one `ObserveDesktop`. Do not capture between every move. On a miss, inset 40–80px; do not crop-loop.
