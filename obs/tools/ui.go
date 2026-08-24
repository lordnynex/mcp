package tools

import (
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/ui"
	"github.com/lordnynex/mcp/obs/session"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerUiTools(s *mcp.Server, h *session.Host) {
	addRequest(s, h, "GetStudioModeEnabled", "Gets whether studio is enabled.",
		func(c *goobs.Client, in *ui.GetStudioModeEnabledParams) (*ui.GetStudioModeEnabledResponse, error) {
			return c.Ui.GetStudioModeEnabled(in)
		})
	addRequest(s, h, "SetStudioModeEnabled", "Enables or disables studio mode",
		func(c *goobs.Client, in *ui.SetStudioModeEnabledParams) (*ui.SetStudioModeEnabledResponse, error) {
			return c.Ui.SetStudioModeEnabled(in)
		})
	addRequest(s, h, "OpenInputPropertiesDialog", "Opens the properties dialog of an input.",
		func(c *goobs.Client, in *ui.OpenInputPropertiesDialogParams) (*ui.OpenInputPropertiesDialogResponse, error) {
			return c.Ui.OpenInputPropertiesDialog(in)
		})
	addRequest(s, h, "OpenInputFiltersDialog", "Opens the filters dialog of an input.",
		func(c *goobs.Client, in *ui.OpenInputFiltersDialogParams) (*ui.OpenInputFiltersDialogResponse, error) {
			return c.Ui.OpenInputFiltersDialog(in)
		})
	addRequest(s, h, "OpenInputInteractDialog", "Opens the interact dialog of an input.",
		func(c *goobs.Client, in *ui.OpenInputInteractDialogParams) (*ui.OpenInputInteractDialogResponse, error) {
			return c.Ui.OpenInputInteractDialog(in)
		})
	addRequest(s, h, "GetMonitorList", "Gets a list of connected monitors and information about them.",
		func(c *goobs.Client, in *ui.GetMonitorListParams) (*ui.GetMonitorListResponse, error) {
			return c.Ui.GetMonitorList(in)
		})
	addRequest(s, h, "OpenVideoMixProjector", "Opens a projector for a specific output video mix.",
		func(c *goobs.Client, in *ui.OpenVideoMixProjectorParams) (*ui.OpenVideoMixProjectorResponse, error) {
			return c.Ui.OpenVideoMixProjector(in)
		})
	addRequest(s, h, "OpenSourceProjector", "Opens a projector for a source.",
		func(c *goobs.Client, in *ui.OpenSourceProjectorParams) (*ui.OpenSourceProjectorResponse, error) {
			return c.Ui.OpenSourceProjector(in)
		})
}
