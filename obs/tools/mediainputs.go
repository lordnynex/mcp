package tools

import (
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/mediainputs"
	"github.com/lordnynex/mcp/obs/session"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerMediaInputsTools(s *mcp.Server, h *session.Host) {
	addRequest(s, h, "GetMediaInputStatus", "Gets the status of a media input.",
		func(c *goobs.Client, in *mediainputs.GetMediaInputStatusParams) (*mediainputs.GetMediaInputStatusResponse, error) {
			return c.MediaInputs.GetMediaInputStatus(in)
		})
	addRequest(s, h, "SetMediaInputCursor", "Sets the cursor position of a media input.",
		func(c *goobs.Client, in *mediainputs.SetMediaInputCursorParams) (*mediainputs.SetMediaInputCursorResponse, error) {
			return c.MediaInputs.SetMediaInputCursor(in)
		})
	addRequest(s, h, "OffsetMediaInputCursor", "Offsets the current cursor position of a media input by the specified value.",
		func(c *goobs.Client, in *mediainputs.OffsetMediaInputCursorParams) (*mediainputs.OffsetMediaInputCursorResponse, error) {
			return c.MediaInputs.OffsetMediaInputCursor(in)
		})
	addRequest(s, h, "TriggerMediaInputAction", "Triggers an action on a media input.",
		func(c *goobs.Client, in *mediainputs.TriggerMediaInputActionParams) (*mediainputs.TriggerMediaInputActionResponse, error) {
			return c.MediaInputs.TriggerMediaInputAction(in)
		})
}
