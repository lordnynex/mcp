package prompts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/lordnynex/mcp/obs/session"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const currentProgramName = "obs-current-program"

// AlwaysNames are registered before Connect.
var AlwaysNames = []string{
	"obs-connect",
	"obs-switch-scene",
	"obs-studio-status",
	"obs-subscribe-events",
	"obs-start-stream",
	"obs-start-record",
	"obs-record-clip",
	"obs-create-browser-source",
	"obs-create-input",
	"obs-request-batch",
}

// RegisterAlways adds prompts that do not require an OBS connection.
func RegisterAlways(s *mcp.Server, h *session.Host) {
	add(s, &mcp.Prompt{
		Name:        "obs-connect",
		Title:       "Connect to OBS",
		Description: "Identify with OBS Studio so request tools become available.",
		Arguments: []*mcp.PromptArgument{
			{Name: "host", Title: "Host", Description: "OBS websocket host:port"},
			{Name: "password", Title: "Password", Description: "OBS websocket password"},
		},
	}, func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		host := arg(req, "host")
		if host == "" {
			host = session.DefaultHost
		}
		pw := arg(req, "password")
		pwNote := "omit password if the CLI already has OBS_PASSWORD"
		if pw != "" {
			pwNote = "include the provided password argument"
		}
		text := fmt.Sprintf("Call the Connect tool with host %q and %s. After success, wait for notifications/tools/list_changed before using OBS request tools. ConnectionStatus reports whether Identify succeeded.", host, pwNote)
		return userPrompt("Connect to OBS", text), nil
	})

	add(s, &mcp.Prompt{
		Name:        "obs-switch-scene",
		Title:       "Switch program scene",
		Description: "Switch the OBS program scene by official scene name.",
		Arguments: []*mcp.PromptArgument{
			{Name: "sceneName", Title: "Scene name", Description: "Official OBS scene name", Required: true},
		},
	}, func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		name := arg(req, "sceneName")
		if name == "" {
			return nil, fmt.Errorf("sceneName is required")
		}
		return userPrompt("Switch program scene", fmt.Sprintf("Call SetCurrentProgramScene with sceneName %q. If the tool is missing, call Connect first.", name)), nil
	})

	add(s, &mcp.Prompt{
		Name:        "obs-studio-status",
		Title:       "OBS studio status",
		Description: "Read connection, version, stream, and record status.",
	}, func(_ context.Context, _ *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return userPrompt("OBS studio status", "Call ConnectionStatus. If connected, then call GetVersion, GetStreamStatus, and GetRecordStatus. Summarize whether OBS is identified, streaming, and recording."), nil
	})

	add(s, &mcp.Prompt{
		Name:        "obs-subscribe-events",
		Title:       "Subscribe to OBS events",
		Description: "Subscribe to an official EventSubscription category.",
		Arguments: []*mcp.PromptArgument{
			{Name: "category", Title: "Category", Description: "EventSubscription category such as Scenes or Inputs", Required: true},
		},
	}, func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		cat := arg(req, "category")
		if cat == "" {
			return nil, fmt.Errorf("category is required")
		}
		text := fmt.Sprintf("Call SubscribeEvents with eventTypes [%q]. Then call resources/subscribe or subscriptions/listen on each returned obs://events/{eventType} URI so notifications/resources/updated is delivered.", cat)
		return userPrompt("Subscribe to OBS events", text), nil
	})

	add(s, &mcp.Prompt{
		Name:        "obs-start-stream",
		Title:       "Start streaming",
		Description: "Checklist before StartStream. StartStream elicits confirmation.",
	}, func(_ context.Context, _ *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return userPrompt("Start streaming", "Call GetStreamStatus. If not already streaming, call StartStream and accept the elicitation confirmation. Do not call StartRecord unless the user asked to record."), nil
	})

	add(s, &mcp.Prompt{
		Name:        "obs-start-record",
		Title:       "Start recording",
		Description: "Start recording immediately. Recording tools never elicit.",
	}, func(_ context.Context, _ *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return userPrompt("Start recording", "Call GetRecordStatus. If not recording, call StartRecord immediately. Do not wait for elicitation; recording tools never ask for confirmation so unattended video capture can finish."), nil
	})

	add(s, &mcp.Prompt{
		Name:        "obs-record-clip",
		Title:       "Record a timed clip",
		Description: "Record a finite clip, stop, restore the program scene, and return the file path.",
		Arguments: []*mcp.PromptArgument{
			{Name: "durationMs", Title: "Duration (ms)", Description: "Minimum recording duration in milliseconds", Required: true},
			{Name: "sceneName", Title: "Scene name", Description: "Optional program scene to record; restore the previous scene after StopRecord"},
		},
	}, func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		ms, err := parseDurationMs(arg(req, "durationMs"))
		if err != nil {
			return nil, err
		}
		scene := arg(req, "sceneName")
		text := fmt.Sprintf("Call GetRecordStatus. If outputActive is true, call StopRecord first so this clip is a new file. Call GetCurrentProgramScene and remember currentProgramSceneName. %sCall StartRecord immediately (recording tools never elicit). Poll GetRecordStatus.outputDuration until it is at least %d milliseconds. Call StopRecord and report outputPath. Restore the remembered program scene with SetCurrentProgramScene. Do not leave the record output running. If the tools are missing, call Connect first.", sceneSwitchSentence(scene), ms)
		return userPrompt("Record a timed clip", text), nil
	})

	add(s, &mcp.Prompt{
		Name:        "obs-create-browser-source",
		Title:       "Create a browser source",
		Description: "Create a transparent browser_source and verify it painted before recording.",
		Arguments: []*mcp.PromptArgument{
			{Name: "sceneName", Title: "Scene name", Description: "Scene that receives the new input", Required: true},
			{Name: "inputName", Title: "Input name", Description: "Name of the new browser source", Required: true},
			{Name: "url", Title: "URL", Description: "Page URL; omit when localFile is set"},
			{Name: "localFile", Title: "Local file", Description: "Absolute HTML path; sets is_local_file"},
			{Name: "width", Title: "Width", Description: "Browser width; default GetVideoSettings.baseWidth"},
			{Name: "height", Title: "Height", Description: "Browser height; default GetVideoSettings.baseHeight"},
		},
	}, func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		scene := arg(req, "sceneName")
		input := arg(req, "inputName")
		if scene == "" {
			return nil, fmt.Errorf("sceneName is required")
		}
		if input == "" {
			return nil, fmt.Errorf("inputName is required")
		}
		page, err := browserPageSentence(arg(req, "url"), arg(req, "localFile"))
		if err != nil {
			return nil, err
		}
		size := browserSizeSentence(arg(req, "width"), arg(req, "height"))
		text := fmt.Sprintf("Call GetInputDefaultSettings with inputKind \"browser_source\". %sCall CreateInput with sceneName %q, inputName %q, inputKind \"browser_source\", and inputSettings using %s, the resolved width/height, and the official default transparent CSS (body { background-color: rgba(0, 0, 0, 0); margin: 0px auto; overflow: hidden; }). There is no inject-JavaScript request; HTML and JS live in url or local_file. Prefer http(s) or is_local_file; file:// and data: URLs often stay blank on macOS CEF. After create, call GetSourceActive and GetSourceScreenshot or SaveSourceScreenshot on the new input before claiming success. If the page did not paint, do not StartRecord. If the tools are missing, call Connect first.", size, scene, input, page)
		return userPrompt("Create a browser source", text), nil
	})

	add(s, &mcp.Prompt{
		Name:        "obs-create-input",
		Title:       "Create an input",
		Description: "Create any input kind from official defaults and verify it before use.",
		Arguments: []*mcp.PromptArgument{
			{Name: "sceneName", Title: "Scene name", Description: "Scene that receives the new input", Required: true},
			{Name: "inputName", Title: "Input name", Description: "Name of the new input", Required: true},
			{Name: "inputKind", Title: "Input kind", Description: "Official input kind from GetInputKindList", Required: true},
		},
	}, func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		scene := arg(req, "sceneName")
		input := arg(req, "inputName")
		kind := arg(req, "inputKind")
		if scene == "" {
			return nil, fmt.Errorf("sceneName is required")
		}
		if input == "" {
			return nil, fmt.Errorf("inputName is required")
		}
		if kind == "" {
			return nil, fmt.Errorf("inputKind is required")
		}
		text := fmt.Sprintf("Call GetInputKindList if inputKind is uncertain; do not invent kinds. Call GetInputDefaultSettings with inputKind %q. Call CreateInput with sceneName %q, inputName %q, inputKind %q, and inputSettings overlaid on those defaults. If you later call SetSceneItemTransform, boundsWidth and boundsHeight must be at least 1 even when boundsType is OBS_BOUNDS_NONE. After create, call GetSourceActive and GetSourceScreenshot or SaveSourceScreenshot on the new input before claiming success. If the tools are missing, call Connect first.", kind, scene, input, kind)
		return userPrompt("Create an input", text), nil
	})

	add(s, &mcp.Prompt{
		Name:        "obs-request-batch",
		Title:       "OBS request batch",
		Description: "How to run several official requests in one call.",
	}, func(_ context.Context, _ *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return userPrompt("OBS request batch", "Use RequestBatch with official requestType names. executionType 0 is SerialRealtime; 1 SerialFrame is unsupported; 2 is Parallel. Sleep is only valid with SerialRealtime. Do not invent request names."), nil
	})

	_ = h
}

// RegisterConnected adds prompts that need a live OBS connection.
func RegisterConnected(s *mcp.Server, h *session.Host) {
	add(s, &mcp.Prompt{
		Name:        currentProgramName,
		Title:       "Current program scene",
		Description: "Describe the current OBS program scene.",
	}, func(_ context.Context, _ *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		c, err := h.RequireClient()
		if err != nil {
			return nil, err
		}
		cur, err := c.Scenes.GetCurrentProgramScene(&scenes.GetCurrentProgramSceneParams{})
		if err != nil {
			return nil, err
		}
		name := cur.SceneName
		if name == "" {
			name = cur.CurrentProgramSceneName
		}
		return userPrompt("Current program scene", fmt.Sprintf("The current OBS program scene is %q. Use GetCurrentProgramScene or SetCurrentProgramScene to read or change it.", name)), nil
	})
}

// UnregisterConnected removes connection-only prompts.
func UnregisterConnected(s *mcp.Server) {
	s.RemovePrompts(currentProgramName)
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

func parseDurationMs(s string) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("durationMs is required")
	}
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("durationMs must be a positive integer")
	}
	return n, nil
}

func sceneSwitchSentence(scene string) string {
	if scene == "" {
		return "Leave the current program scene unchanged unless the user asked to record a different scene. "
	}
	return fmt.Sprintf("Call SetCurrentProgramScene with sceneName %q. ", scene)
}

func browserPageSentence(url, localFile string) (string, error) {
	switch {
	case localFile != "":
		return fmt.Sprintf("is_local_file true and local_file %q", localFile), nil
	case url != "":
		return fmt.Sprintf("url %q", url), nil
	default:
		return "", fmt.Errorf("url or localFile is required")
	}
}

func browserSizeSentence(width, height string) string {
	if width == "" && height == "" {
		return "If width or height were omitted, call GetVideoSettings and use baseWidth/baseHeight. "
	}
	parts := make([]string, 0, 2)
	if width != "" {
		parts = append(parts, "width "+width)
	}
	if height != "" {
		parts = append(parts, "height "+height)
	}
	return "Use " + strings.Join(parts, " and ") + " from the prompt arguments; call GetVideoSettings only for a missing dimension. "
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
