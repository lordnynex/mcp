package tools

import (
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/outputs"
	"github.com/lordnynex/mcp/obs/session"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerOutputsTools(s *mcp.Server, h *session.Host) {
	addRequest(s, h, "GetVirtualCamStatus", "Gets the status of the virtualcam output.",
		func(c *goobs.Client, in *outputs.GetVirtualCamStatusParams) (*outputs.GetVirtualCamStatusResponse, error) {
			return c.Outputs.GetVirtualCamStatus(in)
		})
	addRequest(s, h, "ToggleVirtualCam", "Toggles the state of the virtualcam output.",
		func(c *goobs.Client, in *outputs.ToggleVirtualCamParams) (*outputs.ToggleVirtualCamResponse, error) {
			return c.Outputs.ToggleVirtualCam(in)
		})
	addRequest(s, h, "StartVirtualCam", "Starts the virtualcam output.",
		func(c *goobs.Client, in *outputs.StartVirtualCamParams) (*outputs.StartVirtualCamResponse, error) {
			return c.Outputs.StartVirtualCam(in)
		})
	addRequest(s, h, "StopVirtualCam", "Stops the virtualcam output.",
		func(c *goobs.Client, in *outputs.StopVirtualCamParams) (*outputs.StopVirtualCamResponse, error) {
			return c.Outputs.StopVirtualCam(in)
		})
	addRequest(s, h, "GetReplayBufferStatus", "Gets the status of the replay buffer output.",
		func(c *goobs.Client, in *outputs.GetReplayBufferStatusParams) (*outputs.GetReplayBufferStatusResponse, error) {
			return c.Outputs.GetReplayBufferStatus(in)
		})
	addRequest(s, h, "ToggleReplayBuffer", "Toggles the state of the replay buffer output.",
		func(c *goobs.Client, in *outputs.ToggleReplayBufferParams) (*outputs.ToggleReplayBufferResponse, error) {
			return c.Outputs.ToggleReplayBuffer(in)
		})
	addRequest(s, h, "StartReplayBuffer", "Starts the replay buffer output.",
		func(c *goobs.Client, in *outputs.StartReplayBufferParams) (*outputs.StartReplayBufferResponse, error) {
			return c.Outputs.StartReplayBuffer(in)
		})
	addRequest(s, h, "StopReplayBuffer", "Stops the replay buffer output.",
		func(c *goobs.Client, in *outputs.StopReplayBufferParams) (*outputs.StopReplayBufferResponse, error) {
			return c.Outputs.StopReplayBuffer(in)
		})
	addRequest(s, h, "SaveReplayBuffer", "Saves the contents of the replay buffer output.",
		func(c *goobs.Client, in *outputs.SaveReplayBufferParams) (*outputs.SaveReplayBufferResponse, error) {
			return c.Outputs.SaveReplayBuffer(in)
		})
	addRequest(s, h, "GetLastReplayBufferReplay", "Gets the filename of the last replay buffer save file.",
		func(c *goobs.Client, in *outputs.GetLastReplayBufferReplayParams) (*outputs.GetLastReplayBufferReplayResponse, error) {
			return c.Outputs.GetLastReplayBufferReplay(in)
		})
	addRequest(s, h, "GetOutputList", "Gets the list of available outputs.",
		func(c *goobs.Client, in *outputs.GetOutputListParams) (*outputs.GetOutputListResponse, error) {
			return c.Outputs.GetOutputList(in)
		})
	addRequest(s, h, "GetOutputStatus", "Gets the status of an output.",
		func(c *goobs.Client, in *outputs.GetOutputStatusParams) (*outputs.GetOutputStatusResponse, error) {
			return c.Outputs.GetOutputStatus(in)
		})
	addRequest(s, h, "ToggleOutput", "Toggles the status of an output.",
		func(c *goobs.Client, in *outputs.ToggleOutputParams) (*outputs.ToggleOutputResponse, error) {
			return c.Outputs.ToggleOutput(in)
		})
	addRequest(s, h, "StartOutput", "Starts an output.",
		func(c *goobs.Client, in *outputs.StartOutputParams) (*outputs.StartOutputResponse, error) {
			return c.Outputs.StartOutput(in)
		})
	addRequest(s, h, "StopOutput", "Stops an output.",
		func(c *goobs.Client, in *outputs.StopOutputParams) (*outputs.StopOutputResponse, error) {
			return c.Outputs.StopOutput(in)
		})
	addRequest(s, h, "GetOutputSettings", "Gets the settings of an output.",
		func(c *goobs.Client, in *outputs.GetOutputSettingsParams) (*outputs.GetOutputSettingsResponse, error) {
			return c.Outputs.GetOutputSettings(in)
		})
	addRequest(s, h, "SetOutputSettings", "Sets the settings of an output.",
		func(c *goobs.Client, in *outputs.SetOutputSettingsParams) (*outputs.SetOutputSettingsResponse, error) {
			return c.Outputs.SetOutputSettings(in)
		})
}
