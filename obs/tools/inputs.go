package tools

import (
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/inputs"
	"github.com/lordnynex/mcp/obs/session"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerInputsTools(s *mcp.Server, h *session.Host) {
	addRequest(s, h, "GetInputList", "Gets an array of all inputs in OBS.",
		func(c *goobs.Client, in *inputs.GetInputListParams) (*inputs.GetInputListResponse, error) {
			return c.Inputs.GetInputList(in)
		})
	addRequest(s, h, "GetInputKindList", "Gets an array of all available input kinds in OBS.",
		func(c *goobs.Client, in *inputs.GetInputKindListParams) (*inputs.GetInputKindListResponse, error) {
			return c.Inputs.GetInputKindList(in)
		})
	addRequest(s, h, "GetSpecialInputs", "Gets the names of all special inputs.",
		func(c *goobs.Client, in *inputs.GetSpecialInputsParams) (*inputs.GetSpecialInputsResponse, error) {
			return c.Inputs.GetSpecialInputs(in)
		})
	addRequest(s, h, "CreateInput", "Creates a new input, adding it as a scene item to the specified scene.",
		func(c *goobs.Client, in *inputs.CreateInputParams) (*inputs.CreateInputResponse, error) {
			return c.Inputs.CreateInput(in)
		})
	addRequest(s, h, "RemoveInput", "Removes an existing input.",
		func(c *goobs.Client, in *inputs.RemoveInputParams) (*inputs.RemoveInputResponse, error) {
			return c.Inputs.RemoveInput(in)
		})
	addRequest(s, h, "SetInputName", "Sets the name of an input (rename).",
		func(c *goobs.Client, in *inputs.SetInputNameParams) (*inputs.SetInputNameResponse, error) {
			return c.Inputs.SetInputName(in)
		})
	addRequest(s, h, "GetInputDefaultSettings", "Gets the default settings for an input kind.",
		func(c *goobs.Client, in *inputs.GetInputDefaultSettingsParams) (*inputs.GetInputDefaultSettingsResponse, error) {
			return c.Inputs.GetInputDefaultSettings(in)
		})
	addRequest(s, h, "GetInputSettings", "Gets the settings of an input.",
		func(c *goobs.Client, in *inputs.GetInputSettingsParams) (*inputs.GetInputSettingsResponse, error) {
			return c.Inputs.GetInputSettings(in)
		})
	addRequest(s, h, "SetInputSettings", "Sets the settings of an input.",
		func(c *goobs.Client, in *inputs.SetInputSettingsParams) (*inputs.SetInputSettingsResponse, error) {
			return c.Inputs.SetInputSettings(in)
		})
	addRequest(s, h, "GetInputMute", "Gets the audio mute state of an input.",
		func(c *goobs.Client, in *inputs.GetInputMuteParams) (*inputs.GetInputMuteResponse, error) {
			return c.Inputs.GetInputMute(in)
		})
	addRequest(s, h, "SetInputMute", "Sets the audio mute state of an input.",
		func(c *goobs.Client, in *inputs.SetInputMuteParams) (*inputs.SetInputMuteResponse, error) {
			return c.Inputs.SetInputMute(in)
		})
	addRequest(s, h, "ToggleInputMute", "Toggles the audio mute state of an input.",
		func(c *goobs.Client, in *inputs.ToggleInputMuteParams) (*inputs.ToggleInputMuteResponse, error) {
			return c.Inputs.ToggleInputMute(in)
		})
	addRequest(s, h, "GetInputVolume", "Gets the current volume setting of an input.",
		func(c *goobs.Client, in *inputs.GetInputVolumeParams) (*inputs.GetInputVolumeResponse, error) {
			return c.Inputs.GetInputVolume(in)
		})
	addRequest(s, h, "SetInputVolume", "Sets the volume setting of an input.",
		func(c *goobs.Client, in *inputs.SetInputVolumeParams) (*inputs.SetInputVolumeResponse, error) {
			return c.Inputs.SetInputVolume(in)
		})
	addRequest(s, h, "GetInputAudioBalance", "Gets the audio balance of an input.",
		func(c *goobs.Client, in *inputs.GetInputAudioBalanceParams) (*inputs.GetInputAudioBalanceResponse, error) {
			return c.Inputs.GetInputAudioBalance(in)
		})
	addRequest(s, h, "SetInputAudioBalance", "Sets the audio balance of an input.",
		func(c *goobs.Client, in *inputs.SetInputAudioBalanceParams) (*inputs.SetInputAudioBalanceResponse, error) {
			return c.Inputs.SetInputAudioBalance(in)
		})
	addRequest(s, h, "GetInputAudioSyncOffset", "Gets the audio sync offset of an input.",
		func(c *goobs.Client, in *inputs.GetInputAudioSyncOffsetParams) (*inputs.GetInputAudioSyncOffsetResponse, error) {
			return c.Inputs.GetInputAudioSyncOffset(in)
		})
	addRequest(s, h, "SetInputAudioSyncOffset", "Sets the audio sync offset of an input.",
		func(c *goobs.Client, in *inputs.SetInputAudioSyncOffsetParams) (*inputs.SetInputAudioSyncOffsetResponse, error) {
			return c.Inputs.SetInputAudioSyncOffset(in)
		})
	addRequest(s, h, "GetInputAudioMonitorType", "Gets the audio monitor type of an input.",
		func(c *goobs.Client, in *inputs.GetInputAudioMonitorTypeParams) (*inputs.GetInputAudioMonitorTypeResponse, error) {
			return c.Inputs.GetInputAudioMonitorType(in)
		})
	addRequest(s, h, "SetInputAudioMonitorType", "Sets the audio monitor type of an input.",
		func(c *goobs.Client, in *inputs.SetInputAudioMonitorTypeParams) (*inputs.SetInputAudioMonitorTypeResponse, error) {
			return c.Inputs.SetInputAudioMonitorType(in)
		})
	addRequest(s, h, "GetInputAudioTracks", "Gets the enable state of all audio tracks of an input.",
		func(c *goobs.Client, in *inputs.GetInputAudioTracksParams) (*inputs.GetInputAudioTracksResponse, error) {
			return c.Inputs.GetInputAudioTracks(in)
		})
	addRequest(s, h, "SetInputAudioTracks", "Sets the enable state of audio tracks of an input.",
		func(c *goobs.Client, in *inputs.SetInputAudioTracksParams) (*inputs.SetInputAudioTracksResponse, error) {
			return c.Inputs.SetInputAudioTracks(in)
		})
	addRequest(s, h, "GetInputDeinterlaceMode", "Gets the deinterlace mode of an input.",
		func(c *goobs.Client, in *inputs.GetInputDeinterlaceModeParams) (*inputs.GetInputDeinterlaceModeResponse, error) {
			return c.Inputs.GetInputDeinterlaceMode(in)
		})
	addRequest(s, h, "SetInputDeinterlaceMode", "Sets the deinterlace mode of an input.",
		func(c *goobs.Client, in *inputs.SetInputDeinterlaceModeParams) (*inputs.SetInputDeinterlaceModeResponse, error) {
			return c.Inputs.SetInputDeinterlaceMode(in)
		})
	addRequest(s, h, "GetInputDeinterlaceFieldOrder", "Gets the deinterlace field order of an input.",
		func(c *goobs.Client, in *inputs.GetInputDeinterlaceFieldOrderParams) (*inputs.GetInputDeinterlaceFieldOrderResponse, error) {
			return c.Inputs.GetInputDeinterlaceFieldOrder(in)
		})
	addRequest(s, h, "SetInputDeinterlaceFieldOrder", "Sets the deinterlace field order of an input.",
		func(c *goobs.Client, in *inputs.SetInputDeinterlaceFieldOrderParams) (*inputs.SetInputDeinterlaceFieldOrderResponse, error) {
			return c.Inputs.SetInputDeinterlaceFieldOrder(in)
		})
	addRequest(s, h, "GetInputPropertiesListPropertyItems", "Gets the items of a list property from an input's properties.",
		func(c *goobs.Client, in *inputs.GetInputPropertiesListPropertyItemsParams) (*inputs.GetInputPropertiesListPropertyItemsResponse, error) {
			return c.Inputs.GetInputPropertiesListPropertyItems(in)
		})
	addRequest(s, h, "PressInputPropertiesButton", "Presses a button in the properties of an input.",
		func(c *goobs.Client, in *inputs.PressInputPropertiesButtonParams) (*inputs.PressInputPropertiesButtonResponse, error) {
			return c.Inputs.PressInputPropertiesButton(in)
		})
}
