package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/general"
	"github.com/lordnynex/mcp/obs/session"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerGeneralTools(s *mcp.Server, h *session.Host) {
	addRequest(s, h, "GetVersion", "Gets data about the current plugin and RPC version.",
		func(c *goobs.Client, in *general.GetVersionParams) (*general.GetVersionResponse, error) {
			return c.General.GetVersion(in)
		})
	addRequest(s, h, "GetStats", "Gets statistics about OBS, obs-websocket, and the current session.",
		func(c *goobs.Client, in *general.GetStatsParams) (*general.GetStatsResponse, error) {
			return c.General.GetStats(in)
		})
	addRequest(s, h, "BroadcastCustomEvent", "Broadcasts a CustomEvent to all WebSocket clients. Receivers are clients which are identified and subscribed.",
		func(c *goobs.Client, in *general.BroadcastCustomEventParams) (*general.BroadcastCustomEventResponse, error) {
			return c.General.BroadcastCustomEvent(in)
		})
	addRequest(s, h, "CallVendorRequest", "Call a request registered to a vendor.",
		func(c *goobs.Client, in *general.CallVendorRequestParams) (*general.CallVendorRequestResponse, error) {
			return c.General.CallVendorRequest(in)
		})
	addRequest(s, h, "GetHotkeyList", "Gets an array of all hotkey names in OBS.",
		func(c *goobs.Client, in *general.GetHotkeyListParams) (*general.GetHotkeyListResponse, error) {
			return c.General.GetHotkeyList(in)
		})
	addRequest(s, h, "TriggerHotkeyByName", "Triggers a hotkey using its name. See GetHotkeyList.",
		func(c *goobs.Client, in *general.TriggerHotkeyByNameParams) (*general.TriggerHotkeyByNameResponse, error) {
			return c.General.TriggerHotkeyByName(in)
		})
	addRequest(s, h, "TriggerHotkeyByKeySequence", "Triggers a hotkey using a sequence of keys.",
		func(c *goobs.Client, in *general.TriggerHotkeyByKeySequenceParams) (*general.TriggerHotkeyByKeySequenceResponse, error) {
			return c.General.TriggerHotkeyByKeySequence(in)
		})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "Sleep",
		Title:       "Sleep",
		Description: "Sleeps for a time duration or number of frames. Only available in request batches with types SERIAL_REALTIME or SERIAL_FRAME.",
		Annotations: annotationsFor("Sleep"),
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ general.SleepParams) (*mcp.CallToolResult, *general.SleepResponse, error) {
		return nil, nil, fmt.Errorf("Sleep is only valid inside RequestBatch with SerialRealtime (goobs cannot send OpCode 8 SerialFrame)")
	})
	registerDispatcher("Sleep", func(c *goobs.Client, raw json.RawMessage) (any, error) {
		var in general.SleepParams
		if err := unmarshalParams(raw, &in); err != nil {
			return nil, err
		}
		return c.General.Sleep(&in)
	})
}
