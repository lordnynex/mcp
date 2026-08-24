package prompts

import (
	"context"
	"fmt"
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

func userPrompt(desc, text string) *mcp.GetPromptResult {
	return &mcp.GetPromptResult{
		Description: desc,
		Messages: []*mcp.PromptMessage{{
			Role:    "user",
			Content: &mcp.TextContent{Text: text},
		}},
	}
}
