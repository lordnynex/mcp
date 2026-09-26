# robotgo

Desktop automation MCP server. It drives the local keyboard, mouse, clipboard, and screen through [go-vgo/robotgo](https://github.com/go-vgo/robotgo) (the `v2.0.0-beta2` commit, pinned as a v1-compatible pseudo-version because the v2 tag is not a valid Go module path) so a calling LLM can look at the display before it moves or types.

There is no Connect gate: this process *is* the desktop. Call `ObserveDesktop` once to plan, run a known sequence with `DesktopOps` or `MousePath`, then observe again to verify.

The process logs to stderr via [`lib/logger`](../lib/logger). Startup lines describe server identity, MCP tunables, and screenshot defaults.

## Run

The live robotgo backend is **CGO-free** and selected with build tags. Run the package directory, not a single file.

```bash
# macOS (Quartz via purego)
go run -tags purego ./robotgo/cmd/robotgo
go run -tags purego ./robotgo/cmd/robotgo --help

# Linux X11 / XWayland
go run -tags "purego,x11" ./robotgo/cmd/robotgo

# Linux wlroots Wayland (Sway, Hyprland, Wayfire)
go run -tags "purego,wayland" ./robotgo/cmd/robotgo

# Linux GNOME / KDE (input only — no screenshots)
go run -tags "purego,libei" ./robotgo/cmd/robotgo
```

`go run robotgo/cmd/robotgo/main.go` only compiles `main.go` and fails with `undefined: rootCmd`.

The server speaks MCP on stdout. Logs go to stderr.

A binary built **without** a backend tag installs a stub driver that errors with rebuild instructions. `go test ./...` uses that default so the suite stays CGO-free.

## Backends

| Target | Build | Input | Screenshot | Window APIs |
| --- | --- | --- | --- | --- |
| macOS | `-tags purego` (or `-tags mac`) | yes | yes | no (`ErrNotSupported`) |
| Linux X11 / XWayland | `-tags "purego,x11"` | yes | yes | yes |
| Linux wlroots | `-tags "purego,wayland"` | yes | yes | yes |
| Linux GNOME / KDE | `-tags "purego,libei"` | yes | **no** | **no** |

Bare `-tags purego` on Linux selects **wayland**, which fails on GNOME/KDE. Linux builds must name `x11`, `wayland`, or `libei`.

`libei` cannot meet the look-before-act criterion. `ObserveDesktop` / `CaptureScreen` return a clear error. Rebuild with `x11` or `wayland`, or use a compositor that supports `zwlr_screencopy_v1`.

Optional CGO (no tags, `CGO_ENABLED=1`) is an escape hatch for macOS window management. It is not the default.

## System dependencies

**All**

- Go 1.27
- Runtime permission to control input (and, where required, to capture the screen)

**macOS (`-tags purego`)**

- System Settings → Privacy & Security → **Accessibility**
- System Settings → Privacy & Security → **Screen & System Audio Recording**
- Clipboard: `pbcopy` / `pbpaste` (built-in)
- Xcode is **not** required for the pure-Go backend

**Linux X11 (`-tags "purego,x11"`)**

- An X11 (or XWayland) session and `DISPLAY`
- Clipboard: `xclip` or `xsel` on `PATH`

```bash
# Debian / Ubuntu
sudo apt install xclip   # or xsel
```

**Linux Wayland / wlroots (`-tags "purego,wayland"`)**

- A compositor that supports:
  - `zwlr_virtual_pointer_v1`
  - `zwp_virtual_keyboard_v1`
  - `zwlr_screencopy_v1`
  - `zwlr_foreign_toplevel_management_v1`
- Works on Sway, Hyprland, Wayfire. Not GNOME or KDE.
- Clipboard: `xclip` or `xsel` if you use an Xwayland clipboard; otherwise the compositor's clipboard tools

**Linux libei / GNOME / KDE (`-tags "purego,libei"`)**

- `xdg-desktop-portal`
- `xdg-desktop-portal-gnome` or `xdg-desktop-portal-kde`
- Input only. No screenshots, no window APIs.

**Optional CGO backend**

- macOS: Xcode Command Line Tools (`xcode-select --install`) plus the same Accessibility / Screen Recording permissions
- Linux: `gcc`, `libc6-dev`, `libx11-dev`, `libxtst-dev`, `xclip`/`xsel`

```bash
# Debian / Ubuntu (CGO)
sudo apt install gcc libc6-dev libx11-dev xorg-dev libxtst-dev xclip
```

## Configuration

Runtime values come from Viper: flags, `ROBOTGO_` environment variables, and an optional config file.

| Key | Flag | Env | Default |
| --- | --- | --- | --- |
| `log-level` | `--log-level` | `ROBOTGO_LOG_LEVEL` | `info` |
| `page-size` | `--page-size` | `ROBOTGO_PAGE_SIZE` | 50 |
| `keep-alive` | `--keep-alive` | `ROBOTGO_KEEP_ALIVE` | 60s |
| `keep-alive-failure-threshold` | `--keep-alive-failure-threshold` | `ROBOTGO_KEEP_ALIVE_FAILURE_THRESHOLD` | 2 |
| `http-addr` | `--http-addr` | `ROBOTGO_HTTP_ADDR` | empty (HTTP off) |
| `http-path` | `--http-path` | `ROBOTGO_HTTP_PATH` | `/mcp` |
| `http-stateless` | `--http-stateless` | `ROBOTGO_HTTP_STATELESS` | false |
| `http-json` | `--http-json` | `ROBOTGO_HTTP_JSON` | false (SSE) |
| `http-session-timeout` | `--http-session-timeout` | `ROBOTGO_HTTP_SESSION_TIMEOUT` | 30m |
| `http-only` | `--http-only` | `ROBOTGO_HTTP_ONLY` | false |
| `http-bearer-token` | `--http-bearer-token` | `ROBOTGO_HTTP_BEARER_TOKEN` | empty |
| `key-sleep` | `--key-sleep` | `ROBOTGO_KEY_SLEEP` | 10 |
| `mouse-sleep` | `--mouse-sleep` | `ROBOTGO_MOUSE_SLEEP` | 0 |
| `screenshot-max-width` | `--screenshot-max-width` | `ROBOTGO_SCREENSHOT_MAX_WIDTH` | 1280 |
| `screenshot-format` | `--screenshot-format` | `ROBOTGO_SCREENSHOT_FORMAT` | `png` |

`--config` selects a file. If omitted, the CLI looks for `robotgo.yaml` in the current directory and `$HOME/.config/mcp`. A missing file is ignored.

```bash
go run -tags purego ./robotgo/cmd/robotgo --screenshot-max-width 1024 --screenshot-format jpeg
ROBOTGO_LOG_LEVEL=debug go run -tags purego ./robotgo/cmd/robotgo
```

## Observe-first workflow

1. Call `ObserveDesktop` once (image + cursor + displays + mapping fields).
2. Map image pixels: `mouseX = originX + ix * mouseWidth / width` (same for `y`). `width`/`height` are the encoded image; `mouseWidth` is `MouseClick` space.
3. If the next actions do not need a new screenshot, send them as one `DesktopOps` call (`MousePath` may be an op; `smooth: false` for geometric traces). Aim clicks at the geometric center of the control.
4. Verify with `observeAfter` (optional crop) or one `ObserveDesktop`. On a miss, inset 40–80px; do not crop-loop.

Prompts and skills encode this loop. Do not invent coordinates or key names. `ObserveDesktop.cmdCtrl` is `cmd` on macOS and `ctrl` elsewhere.

## Tools

| Tool | Purpose |
| --- | --- |
| `ObserveDesktop` | Screenshot plus cursor, displays, focused title, and image-to-mouse mapping (`originX`, `mouseWidth`, …). |
| `MouseMove` | Absolute / relative / smooth move. |
| `MouseClick` | Click; optional move-to; `double` / `count`. |
| `MouseToggle` | Button down/up. |
| `MouseScroll` | Ticks or `dir`; optional smooth. |
| `MouseDrag` | `DragSmooth` to x,y (cursor must already be at the start). |
| `MouseSelect` | Drag-select from `x,y` to `toX,toY` (left down, move, up). |
| `MousePath` | Array of points; shared `smooth` / `relative`; `hold` to drag. Also a `DesktopOps` op. Use `smooth: false` for traces. |
| `MouseLocation` | Current cursor. |
| `Type` | UTF-8 string. |
| `KeyTap` | Named key + modifiers. |
| `KeyToggle` | Key down/up. |
| `KeyPress` | Down, short delay, up. |
| `CaptureScreen` | Image plus the same mapping fields as `ObserveDesktop`. |
| `GetScreenInfo` | Displays, sizes, `cmdCtrl`. |
| `GetPixelColor` | Hex at x,y or cursor. |
| `ClipboardRead` / `ClipboardWrite` / `ClipboardPaste` | Clipboard; paste is write + Cmd/Ctrl+V. |
| `GetWindowTitle` / `FocusWindow` / `GetWindowBounds` / `MinWindow` / `MaxWindow` / `CloseWindow` | Window APIs (unsupported on mac purego / libei). |
| `ListProcesses` / `FindProcess` | Process list. |
| `KillProcess` | Kill by pid. **Elicits confirmation.** |
| `SetDelay` | `KeySleep` / `MouseSleep`. |
| `Sleep` | Server-side milliseconds. |
| `DesktopOps` | Serial batch including `MousePath`. `haltOnFailure` defaults true. `observeAfter` can crop; `observeOnError` captures after a halt. `KillProcess` is not allowed. |

## Prompts and skills

**MCP prompts:** `robotgo-observe-then-act`, `robotgo-click-target`, `robotgo-type-into-field`, `robotgo-shortcut`, `robotgo-scroll-read`, `robotgo-copy-paste`, `robotgo-focus-app`, `robotgo-drag-drop`, `robotgo-desktop-ops`, `robotgo-backend-limits`.

**Skills** (markdown for calling agents, also `robotgo://skills/{name}`):

| Skill | Path |
| --- | --- |
| aim-center | [skills/aim-center/SKILL.md](skills/aim-center/SKILL.md) |
| observe-then-act | [skills/observe-then-act/SKILL.md](skills/observe-then-act/SKILL.md) |
| keyboard | [skills/keyboard/SKILL.md](skills/keyboard/SKILL.md) |
| mouse | [skills/mouse/SKILL.md](skills/mouse/SKILL.md) |
| screenshot | [skills/screenshot/SKILL.md](skills/screenshot/SKILL.md) |

Completions cover key names, modifiers, scroll directions, and skill resource names.

## HTTP

Stdio is the default. Set `--http-addr` (for example `127.0.0.1:8080`) to serve streamable HTTP. If stdin is a terminal, stdio is skipped so keepalive pings are not written to the console. Use `--http-only` to skip stdio even when stdin is a pipe.

- **Stateful (default).** `Mcp-Session-Id` after initialize. `--http-stateless` for protocol `2026-07-28` HTTP clients (elicitation cannot complete on that path).
- **`--http-json`.** `application/json` instead of SSE.
- **`--http-bearer-token`.** Optional static bearer. The token is never logged (`tokenSet` only).

## MCP features

- **Instructions** — observe-then-batch loop, center-aim, key names, Linux tags, prompt names.
- **Prompts and resources** — workflows plus `robotgo://skills/{name}`.
- **Completions** — keys, modifiers, scroll dirs, DesktopOps op names, skill names.
- **Progress** — `ObserveDesktop`, `CaptureScreen`, and `DesktopOps` when the client sends a progress token.
- **Elicitation** — `KillProcess` only.
- **Pagination** — `--page-size`.
- **Keepalive** — `--keep-alive` / `--keep-alive-failure-threshold`.
- **List cache hints** — `ttlMs` on tools/resources/prompts.

Not implemented: CBitmap APIs, OpenCV/`gcv`, `gohook` event hooks, OCR, ADB, native `Alert` dialogs, deprecated robotgo aliases (`TypeStr`, `Drag`, `GetMousePos`). TLS and full OAuth are not in this server.

## Tests

```bash
go test ./...
```

Unit tests use a fake desktop driver and never move the real pointer. The live adapter (`desktop/live`) is compiled only with a backend tag.
