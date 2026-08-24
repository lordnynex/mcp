package tools

import (
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/transitions"
	"github.com/lordnynex/mcp/obs/session"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerTransitionsTools(s *mcp.Server, h *session.Host) {
	addRequest(s, h, "GetTransitionKindList", "Gets an array of all available transition kinds.",
		func(c *goobs.Client, in *transitions.GetTransitionKindListParams) (*transitions.GetTransitionKindListResponse, error) {
			return c.Transitions.GetTransitionKindList(in)
		})
	addRequest(s, h, "GetSceneTransitionList", "Gets an array of all scene transitions in OBS.",
		func(c *goobs.Client, in *transitions.GetSceneTransitionListParams) (*transitions.GetSceneTransitionListResponse, error) {
			return c.Transitions.GetSceneTransitionList(in)
		})
	addRequest(s, h, "GetCurrentSceneTransition", "Gets information about the current scene transition.",
		func(c *goobs.Client, in *transitions.GetCurrentSceneTransitionParams) (*transitions.GetCurrentSceneTransitionResponse, error) {
			return c.Transitions.GetCurrentSceneTransition(in)
		})
	addRequest(s, h, "SetCurrentSceneTransition", "Sets the current scene transition.",
		func(c *goobs.Client, in *transitions.SetCurrentSceneTransitionParams) (*transitions.SetCurrentSceneTransitionResponse, error) {
			return c.Transitions.SetCurrentSceneTransition(in)
		})
	addRequest(s, h, "SetCurrentSceneTransitionDuration", "Sets the duration of the current scene transition, if it is not fixed.",
		func(c *goobs.Client, in *transitions.SetCurrentSceneTransitionDurationParams) (*transitions.SetCurrentSceneTransitionDurationResponse, error) {
			return c.Transitions.SetCurrentSceneTransitionDuration(in)
		})
	addRequest(s, h, "SetCurrentSceneTransitionSettings", "Sets the settings of the current scene transition.",
		func(c *goobs.Client, in *transitions.SetCurrentSceneTransitionSettingsParams) (*transitions.SetCurrentSceneTransitionSettingsResponse, error) {
			return c.Transitions.SetCurrentSceneTransitionSettings(in)
		})
	addRequest(s, h, "GetCurrentSceneTransitionCursor", "Gets the cursor position of the current scene transition.",
		func(c *goobs.Client, in *transitions.GetCurrentSceneTransitionCursorParams) (*transitions.GetCurrentSceneTransitionCursorResponse, error) {
			return c.Transitions.GetCurrentSceneTransitionCursor(in)
		})
	addRequest(s, h, "TriggerStudioModeTransition", "Triggers the current scene transition. Same functionality as the Transition button in studio mode.",
		func(c *goobs.Client, in *transitions.TriggerStudioModeTransitionParams) (*transitions.TriggerStudioModeTransitionResponse, error) {
			return c.Transitions.TriggerStudioModeTransition(in)
		})
	addRequest(s, h, "SetTBarPosition", "Sets the position of the TBar.",
		func(c *goobs.Client, in *transitions.SetTBarPositionParams) (*transitions.SetTBarPositionResponse, error) {
			return c.Transitions.SetTBarPosition(in)
		})
}
