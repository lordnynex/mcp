package prompts

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Names is the complete prompt list.
var Names = []string{
	"robotgo-observe-then-act",
	"robotgo-click-target",
	"robotgo-type-into-field",
	"robotgo-shortcut",
	"robotgo-scroll-read",
	"robotgo-copy-paste",
	"robotgo-focus-app",
	"robotgo-drag-drop",
	"robotgo-desktop-ops",
	"robotgo-backend-limits",
}

// Register adds computer-use workflow prompts.
func Register(s *mcp.Server) {
	add(s, &mcp.Prompt{
		Name:        "robotgo-observe-then-act",
		Title:       "Observe then act",
		Description: "Always screenshot before mouse or keyboard input, then verify.",
	}, func(_ context.Context, _ *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return userPrompt("Observe then act", "Call ObserveDesktop once to plan. Map image pixels: mouseX = originX + ix * mouseWidth / width (same for y). Aim clicks at the geometric center of the control, not the first glyph or an edge. If the next actions do not need a new screenshot, send them as one DesktopOps call (MousePath is a valid op; use smooth false for geometric traces). Verify with observeAfter or one ObserveDesktop. On a miss, inset 40-80px; do not crop-loop. Do not invent coordinates. If capture fails with a backend error, stop and report robotgo-backend-limits."), nil
	})

	add(s, &mcp.Prompt{
		Name:        "robotgo-click-target",
		Title:       "Click a visible target",
		Description: "Find on-screen coordinates from a screenshot and click.",
		Arguments: []*mcp.PromptArgument{
			{Name: "target", Title: "Target", Description: "What to click, in the user's words", Required: true},
		},
	}, func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		target := arg(req, "target")
		if target == "" {
			return nil, fmt.Errorf("target is required")
		}
		text := fmt.Sprintf("Call ObserveDesktop. Map image pixels with originX/mouseWidth/width. Locate %q and estimate its axis-aligned box. Click the geometric center (x+w/2, y+h/2), inset from borders, placeholders, and chrome — not the first glyph. Prefer DesktopOps with MouseClick at that mapped center. Verify once. On a miss, inset 40-80px; do not crop-loop. If the target is not visible, MouseScroll or switch windows, then observe again. Do not guess coordinates.", target)
		return userPrompt("Click a visible target", text), nil
	})

	add(s, &mcp.Prompt{
		Name:        "robotgo-type-into-field",
		Title:       "Type into a field",
		Description: "Focus a text field, type, and verify.",
		Arguments: []*mcp.PromptArgument{
			{Name: "text", Title: "Text", Description: "String to type", Required: true},
			{Name: "field", Title: "Field", Description: "Visible field to click first"},
		},
	}, func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		text := arg(req, "text")
		if text == "" {
			return nil, fmt.Errorf("text is required")
		}
		field := arg(req, "field")
		click := "Call ObserveDesktop. If a text field is already focused, skip the click. "
		if field != "" {
			click = fmt.Sprintf("Call ObserveDesktop. Click the empty interior center of the field described as %q (not the placeholder's first letter). ", field)
		}
		body := click + fmt.Sprintf("Map image pixels with originX/mouseWidth/width. Send one DesktopOps with MouseClick at that mapped center then Type with text %q. Use KeyTap for Enter/Tab/Escape, not Type. Verify once (observeAfter crop is fine). On a miss, inset 40-80px; do not crop-loop. If a shortcut is needed, use KeyTap with ObserveDesktop.cmdCtrl instead of typing the modifier glyph.", text)
		return userPrompt("Type into a field", body), nil
	})

	add(s, &mcp.Prompt{
		Name:        "robotgo-shortcut",
		Title:       "Keyboard shortcut",
		Description: "Tap a named key with platform modifiers.",
		Arguments: []*mcp.PromptArgument{
			{Name: "key", Title: "Key", Description: "robotgo key name", Required: true},
			{Name: "modifiers", Title: "Modifiers", Description: "Comma-separated modifiers; use cmdctrl for the platform accelerator"},
		},
	}, func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		key := arg(req, "key")
		if key == "" {
			return nil, fmt.Errorf("key is required")
		}
		mods := arg(req, "modifiers")
		modNote := "Call ObserveDesktop and read cmdCtrl. "
		if mods != "" {
			modNote += fmt.Sprintf("Use modifiers %q; replace cmdctrl with the cmdCtrl value. ", mods)
		} else {
			modNote += "Use cmdCtrl when the shortcut is the platform accelerator. "
		}
		body := modNote + fmt.Sprintf("Call KeyTap with key %q and the resolved modifiers array. Do not invent key names. Call ObserveDesktop after the tap.", key)
		return userPrompt("Keyboard shortcut", body), nil
	})

	add(s, &mcp.Prompt{
		Name:        "robotgo-scroll-read",
		Title:       "Scroll and read",
		Description: "Scroll, then screenshot again.",
		Arguments: []*mcp.PromptArgument{
			{Name: "dir", Title: "Direction", Description: "up, down, left, or right"},
		},
	}, func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		dir := arg(req, "dir")
		if dir == "" {
			dir = "down"
		}
		text := fmt.Sprintf("Call ObserveDesktop. Call MouseScroll with dir %q (or x/y ticks). Call ObserveDesktop again and read the new contents. Repeat until the target is visible or the page stops changing.", dir)
		return userPrompt("Scroll and read", text), nil
	})

	add(s, &mcp.Prompt{
		Name:        "robotgo-copy-paste",
		Title:       "Copy or paste",
		Description: "Use the clipboard instead of retyping long text.",
		Arguments: []*mcp.PromptArgument{
			{Name: "text", Title: "Text", Description: "Text to paste; omit to read the clipboard"},
		},
	}, func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		text := arg(req, "text")
		if text == "" {
			return userPrompt("Copy or paste", "Call ClipboardRead and return the text. To copy the current selection, call ObserveDesktop, then KeyTap with key c and modifiers [cmdCtrl from ObserveDesktop]."), nil
		}
		body := fmt.Sprintf("Call ObserveDesktop and click the destination field if needed. Call ClipboardPaste with text %q (writes the clipboard and sends Cmd+V / Ctrl+V). Prefer this over Type for long or formatted text. Call ObserveDesktop to verify.", text)
		return userPrompt("Copy or paste", body), nil
	})

	add(s, &mcp.Prompt{
		Name:        "robotgo-focus-app",
		Title:       "Focus an application",
		Description: "Find a process and focus its window.",
		Arguments: []*mcp.PromptArgument{
			{Name: "name", Title: "Name", Description: "App or window name substring", Required: true},
		},
	}, func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		name := arg(req, "name")
		if name == "" {
			return nil, fmt.Errorf("name is required")
		}
		text := fmt.Sprintf("Call FindProcess with name %q (or ListProcesses if the name is uncertain). Call FocusWindow with name %q. Call ObserveDesktop and confirm windowTitle or the screenshot. On macOS -tags purego and Linux libei, FocusWindow returns not supported; tell the user to focus the app manually or rebuild with a backend that implements window management.", name, name)
		return userPrompt("Focus an application", text), nil
	})

	add(s, &mcp.Prompt{
		Name:        "robotgo-drag-drop",
		Title:       "Drag and drop",
		Description: "Drag from one point to another after observing.",
		Arguments: []*mcp.PromptArgument{
			{Name: "from", Title: "From", Description: "Description of the drag source", Required: true},
			{Name: "to", Title: "To", Description: "Description of the drop target", Required: true},
		},
	}, func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		from := arg(req, "from")
		to := arg(req, "to")
		if from == "" || to == "" {
			return nil, fmt.Errorf("from and to are required")
		}
		text := fmt.Sprintf("Call ObserveDesktop. Read coordinates for %q, MouseMove there. Read coordinates for %q. Call MouseDrag with the destination x,y (DragSmooth). Call ObserveDesktop and confirm the drop. Do not guess coordinates.", from, to)
		return userPrompt("Drag and drop", text), nil
	})

	add(s, &mcp.Prompt{
		Name:        "robotgo-desktop-ops",
		Title:       "Desktop ops batch",
		Description: "When to use DesktopOps / MousePath and which op names are valid.",
		Arguments: []*mcp.PromptArgument{
			{Name: "op", Title: "Op", Description: "Optional DesktopOps op name to mention"},
		},
	}, func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		body := "After one ObserveDesktop, map image pixels (mouseX = originX + ix * mouseWidth / width) and send a known sequence as one DesktopOps (serial). MousePath is a valid op; use smooth false for geometric traces. haltOnFailure defaults true. observeAfter attaches one verify screenshot (crop with x,y,w,h). observeOnError also captures after a halt. Allowed op names: MouseMove, MouseClick, MouseToggle, MouseScroll, MouseDrag, MouseSelect, MousePath, Type, KeyTap, KeyToggle, KeyPress, ClipboardWrite, ClipboardPaste, Sleep, SetDelay, GetWindowTitle, FocusWindow, GetWindowBounds, MinWindow, MaxWindow, CloseWindow. KillProcess is not allowed. Select text with MouseSelect, not a click on the first glyph. Aim clicks at control centers. On a miss, inset 40-80px; do not crop-loop. Observe again only when the next decision needs new pixels."
		if op := arg(req, "op"); op != "" {
			body += fmt.Sprintf(" The caller asked about op %q.", op)
		}
		return userPrompt("Desktop ops batch", body), nil
	})

	add(s, &mcp.Prompt{
		Name:        "robotgo-backend-limits",
		Title:       "Backend limits",
		Description: "Which build tags provide screenshots and window APIs.",
	}, func(_ context.Context, _ *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		text := strings.Join([]string{
			"This server is CGO-free and the backend is selected at build time.",
			"macOS: -tags purego (or mac) has keyboard, mouse, and screenshots; window APIs return not supported.",
			"Linux X11/XWayland: -tags \"purego,x11\" has input, screenshots, and window APIs.",
			"Linux wlroots (Sway, Hyprland, Wayfire): -tags \"purego,wayland\" has input, screenshots, and window APIs.",
			"Linux GNOME/KDE: -tags \"purego,libei\" has input only. ObserveDesktop and CaptureScreen fail. Look-before-act is not possible; tell the user to rebuild with x11 or wayland.",
			"Bare -tags purego on Linux selects wayland, which fails on GNOME/KDE.",
			"Linux clipboard needs xclip or xsel. macOS needs Accessibility and Screen Recording permission.",
		}, " ")
		return userPrompt("Backend limits", text), nil
	})
}

func add(s *mcp.Server, p *mcp.Prompt, h mcp.PromptHandler) {
	s.AddPrompt(p, h)
}

func arg(req *mcp.GetPromptRequest, name string) string {
	if req == nil || req.Params == nil || req.Params.Arguments == nil {
		return ""
	}
	return strings.TrimSpace(req.Params.Arguments[name])
}

func userPrompt(desc, text string) *mcp.GetPromptResult {
	return &mcp.GetPromptResult{
		Description: desc,
		Messages: []*mcp.PromptMessage{{
			Role:    "user",
			Content: &mcp.TextContent{Text: text},
		}},
	}
}
