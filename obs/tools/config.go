package tools

import (
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/config"
	"github.com/lordnynex/mcp/obs/session"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerConfigTools(s *mcp.Server, h *session.Host) {
	addRequest(s, h, "GetPersistentData", "Gets the value of a \\\"slot\\\" from the selected persistent data realm.",
		func(c *goobs.Client, in *config.GetPersistentDataParams) (*config.GetPersistentDataResponse, error) {
			return c.Config.GetPersistentData(in)
		})
	addRequest(s, h, "SetPersistentData", "Sets the value of a \\\"slot\\\" from the selected persistent data realm.",
		func(c *goobs.Client, in *config.SetPersistentDataParams) (*config.SetPersistentDataResponse, error) {
			return c.Config.SetPersistentData(in)
		})
	addRequest(s, h, "GetSceneCollectionList", "Gets an array of all scene collections",
		func(c *goobs.Client, in *config.GetSceneCollectionListParams) (*config.GetSceneCollectionListResponse, error) {
			return c.Config.GetSceneCollectionList(in)
		})
	addRequest(s, h, "SetCurrentSceneCollection", "Switches to a scene collection.",
		func(c *goobs.Client, in *config.SetCurrentSceneCollectionParams) (*config.SetCurrentSceneCollectionResponse, error) {
			return c.Config.SetCurrentSceneCollection(in)
		})
	addRequest(s, h, "CreateSceneCollection", "Creates a new scene collection, switching to it in the process.",
		func(c *goobs.Client, in *config.CreateSceneCollectionParams) (*config.CreateSceneCollectionResponse, error) {
			return c.Config.CreateSceneCollection(in)
		})
	addRequest(s, h, "GetProfileList", "Gets an array of all profiles",
		func(c *goobs.Client, in *config.GetProfileListParams) (*config.GetProfileListResponse, error) {
			return c.Config.GetProfileList(in)
		})
	addRequest(s, h, "SetCurrentProfile", "Switches to a profile.",
		func(c *goobs.Client, in *config.SetCurrentProfileParams) (*config.SetCurrentProfileResponse, error) {
			return c.Config.SetCurrentProfile(in)
		})
	addRequest(s, h, "CreateProfile", "Creates a new profile, switching to it in the process",
		func(c *goobs.Client, in *config.CreateProfileParams) (*config.CreateProfileResponse, error) {
			return c.Config.CreateProfile(in)
		})
	addRequest(s, h, "RemoveProfile", "Removes a profile. If the current profile is chosen, it will change to a different profile first.",
		func(c *goobs.Client, in *config.RemoveProfileParams) (*config.RemoveProfileResponse, error) {
			return c.Config.RemoveProfile(in)
		})
	addRequest(s, h, "GetProfileParameter", "Gets a parameter from the current profile's configuration.",
		func(c *goobs.Client, in *config.GetProfileParameterParams) (*config.GetProfileParameterResponse, error) {
			return c.Config.GetProfileParameter(in)
		})
	addRequest(s, h, "SetProfileParameter", "Sets the value of a parameter in the current profile's configuration.",
		func(c *goobs.Client, in *config.SetProfileParameterParams) (*config.SetProfileParameterResponse, error) {
			return c.Config.SetProfileParameter(in)
		})
	addRequest(s, h, "GetVideoSettings", "Gets the current video settings.",
		func(c *goobs.Client, in *config.GetVideoSettingsParams) (*config.GetVideoSettingsResponse, error) {
			return c.Config.GetVideoSettings(in)
		})
	addRequest(s, h, "SetVideoSettings", "Sets the current video settings.",
		func(c *goobs.Client, in *config.SetVideoSettingsParams) (*config.SetVideoSettingsResponse, error) {
			return c.Config.SetVideoSettings(in)
		})
	addRequest(s, h, "GetStreamServiceSettings", "Gets the current stream service settings (stream destination).",
		func(c *goobs.Client, in *config.GetStreamServiceSettingsParams) (*config.GetStreamServiceSettingsResponse, error) {
			return c.Config.GetStreamServiceSettings(in)
		})
	addRequest(s, h, "SetStreamServiceSettings", "Sets the current stream service settings (stream destination).",
		func(c *goobs.Client, in *config.SetStreamServiceSettingsParams) (*config.SetStreamServiceSettingsResponse, error) {
			return c.Config.SetStreamServiceSettings(in)
		})
	addRequest(s, h, "GetRecordDirectory", "Gets the current directory that the record output is set to.",
		func(c *goobs.Client, in *config.GetRecordDirectoryParams) (*config.GetRecordDirectoryResponse, error) {
			return c.Config.GetRecordDirectory(in)
		})
	addRequest(s, h, "SetRecordDirectory", "Sets the current directory that the record output writes files to.",
		func(c *goobs.Client, in *config.SetRecordDirectoryParams) (*config.SetRecordDirectoryResponse, error) {
			return c.Config.SetRecordDirectory(in)
		})
}
