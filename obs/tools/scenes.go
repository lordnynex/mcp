package tools

import (
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/lordnynex/mcp/obs/session"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerScenesTools(s *mcp.Server, h *session.Host) {
	addRequest(s, h, "GetSceneList", "Gets an array of scenes in OBS.",
		func(c *goobs.Client, in *scenes.GetSceneListParams) (*scenes.GetSceneListResponse, error) {
			return c.Scenes.GetSceneList(in)
		})
	addRequest(s, h, "GetGroupList", "Gets an array of all groups in OBS.",
		func(c *goobs.Client, in *scenes.GetGroupListParams) (*scenes.GetGroupListResponse, error) {
			return c.Scenes.GetGroupList(in)
		})
	addRequest(s, h, "GetCurrentProgramScene", "Gets the current program scene.",
		func(c *goobs.Client, in *scenes.GetCurrentProgramSceneParams) (*scenes.GetCurrentProgramSceneResponse, error) {
			return c.Scenes.GetCurrentProgramScene(in)
		})
	addRequest(s, h, "SetCurrentProgramScene", "Sets the current program scene.",
		func(c *goobs.Client, in *scenes.SetCurrentProgramSceneParams) (*scenes.SetCurrentProgramSceneResponse, error) {
			return c.Scenes.SetCurrentProgramScene(in)
		})
	addRequest(s, h, "GetCurrentPreviewScene", "Gets the current preview scene.",
		func(c *goobs.Client, in *scenes.GetCurrentPreviewSceneParams) (*scenes.GetCurrentPreviewSceneResponse, error) {
			return c.Scenes.GetCurrentPreviewScene(in)
		})
	addRequest(s, h, "SetCurrentPreviewScene", "Sets the current preview scene.",
		func(c *goobs.Client, in *scenes.SetCurrentPreviewSceneParams) (*scenes.SetCurrentPreviewSceneResponse, error) {
			return c.Scenes.SetCurrentPreviewScene(in)
		})
	addRequest(s, h, "CreateScene", "Creates a new scene in OBS.",
		func(c *goobs.Client, in *scenes.CreateSceneParams) (*scenes.CreateSceneResponse, error) {
			return c.Scenes.CreateScene(in)
		})
	addRequest(s, h, "RemoveScene", "Removes a scene from OBS.",
		func(c *goobs.Client, in *scenes.RemoveSceneParams) (*scenes.RemoveSceneResponse, error) {
			return c.Scenes.RemoveScene(in)
		})
	addRequest(s, h, "SetSceneName", "Sets the name of a scene (rename).",
		func(c *goobs.Client, in *scenes.SetSceneNameParams) (*scenes.SetSceneNameResponse, error) {
			return c.Scenes.SetSceneName(in)
		})
	addRequest(s, h, "GetSceneSceneTransitionOverride", "Gets the scene transition overridden for a scene.",
		func(c *goobs.Client, in *scenes.GetSceneSceneTransitionOverrideParams) (*scenes.GetSceneSceneTransitionOverrideResponse, error) {
			return c.Scenes.GetSceneSceneTransitionOverride(in)
		})
	addRequest(s, h, "SetSceneSceneTransitionOverride", "Sets the scene transition overridden for a scene.",
		func(c *goobs.Client, in *scenes.SetSceneSceneTransitionOverrideParams) (*scenes.SetSceneSceneTransitionOverrideResponse, error) {
			return c.Scenes.SetSceneSceneTransitionOverride(in)
		})
}
