package tools

import (
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/canvases"
	"github.com/lordnynex/mcp/obs/session"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerCanvasesTools(s *mcp.Server, h *session.Host) {
	addRequest(s, h, "GetCanvasList", "Gets an array of canvases in OBS.",
		func(c *goobs.Client, in *canvases.GetCanvasListParams) (*canvases.GetCanvasListResponse, error) {
			return c.Canvases.GetCanvasList(in)
		})
}
