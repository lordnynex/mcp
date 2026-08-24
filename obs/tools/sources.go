package tools

import (
	"context"
	"encoding/json"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/sources"
	"github.com/lordnynex/mcp/obs/session"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerSourcesTools(s *mcp.Server, h *session.Host) {
	addRequest(s, h, "GetSourceActive", "Gets the active and show state of a source.",
		func(c *goobs.Client, in *sources.GetSourceActiveParams) (*sources.GetSourceActiveResponse, error) {
			return c.Sources.GetSourceActive(in)
		})
	addRequest(s, h, "SaveSourceScreenshot", "Saves a screenshot of a source to the filesystem.",
		func(c *goobs.Client, in *sources.SaveSourceScreenshotParams) (*sources.SaveSourceScreenshotResponse, error) {
			return c.Sources.SaveSourceScreenshot(in)
		})

	registerDispatcher("GetSourceScreenshot", func(c *goobs.Client, raw json.RawMessage) (any, error) {
		var in sources.GetSourceScreenshotParams
		if err := unmarshalParams(raw, &in); err != nil {
			return nil, err
		}
		return c.Sources.GetSourceScreenshot(&in)
	})
	mcp.AddTool(s, &mcp.Tool{
		Name:        "GetSourceScreenshot",
		Title:       humanTitle("GetSourceScreenshot"),
		Description: "Gets a Base64-encoded screenshot of a source. The imageWidth and imageHeight parameters are treated as \"scale to inner\", meaning the smallest ratio will be used and the aspect ratio of the original resolution is kept.",
		Annotations: annotationsFor("GetSourceScreenshot"),
	}, func(ctx context.Context, req *mcp.CallToolRequest, in sources.GetSourceScreenshotParams) (*mcp.CallToolResult, *sources.GetSourceScreenshotResponse, error) {
		c, err := h.RequireClient()
		if err != nil {
			return nil, nil, err
		}
		notifyProgress(ctx, req, 0, 2, "capturing source")
		out, err := c.Sources.GetSourceScreenshot(&in)
		if err != nil {
			return nil, nil, err
		}
		notifyProgress(ctx, req, 1, 2, "encoding image")
		res := &mcp.CallToolResult{}
		if out != nil && out.ImageData != "" {
			format := ""
			if in.ImageFormat != nil {
				format = *in.ImageFormat
			}
			res.Content = []mcp.Content{screenshotImage(out.ImageData, format)}
		}
		notifyProgress(ctx, req, 2, 2, "screenshot ready")
		return res, out, nil
	})
}
