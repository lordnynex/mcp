package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/stream"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/lordnynex/mcp/obs/session"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerStreamTools(s *mcp.Server, h *session.Host) {
	addRequest(s, h, "GetStreamStatus", "Gets the status of the stream output.",
		func(c *goobs.Client, in *stream.GetStreamStatusParams) (*stream.GetStreamStatusResponse, error) {
			return c.Stream.GetStreamStatus(in)
		})
	addRequest(s, h, "ToggleStream", "Toggles the status of the stream output.",
		func(c *goobs.Client, in *stream.ToggleStreamParams) (*stream.ToggleStreamResponse, error) {
			return c.Stream.ToggleStream(in)
		})
	registerStartStream(s, h)
	addRequest(s, h, "StopStream", "Stops the stream output.",
		func(c *goobs.Client, in *stream.StopStreamParams) (*stream.StopStreamResponse, error) {
			return c.Stream.StopStream(in)
		})
	addRequest(s, h, "SendStreamCaption", "Sends CEA-608 caption text over the stream output.",
		func(c *goobs.Client, in *stream.SendStreamCaptionParams) (*stream.SendStreamCaptionResponse, error) {
			return c.Stream.SendStreamCaption(in)
		})
}

func registerStartStream(s *mcp.Server, h *session.Host) {
	registerDispatcher("StartStream", func(c *goobs.Client, raw json.RawMessage) (any, error) {
		var in stream.StartStreamParams
		if err := unmarshalParams(raw, &in); err != nil {
			return nil, err
		}
		return c.Stream.StartStream(&in)
	})
	mcp.AddTool(s, &mcp.Tool{
		Name:        "StartStream",
		Title:       humanTitle("StartStream"),
		Description: "Starts the stream output. Asks for confirmation via elicitation before starting.",
		Annotations: annotationsFor("StartStream"),
	}, func(ctx context.Context, req *mcp.CallToolRequest, in stream.StartStreamParams) (*mcp.CallToolResult, *stream.StartStreamResponse, error) {
		if req == nil || req.Params == nil || len(req.Params.InputResponses) == 0 {
			return startStreamElicit(), nil, nil
		}
		if err := acceptedElicit(req, "confirm"); err != nil {
			return nil, nil, err
		}
		c, err := h.RequireClient()
		if err != nil {
			return nil, nil, err
		}
		out, err := c.Stream.StartStream(&in)
		return nil, out, err
	})
}

func startStreamElicit() *mcp.CallToolResult {
	return &mcp.CallToolResult{
		InputRequests: mcp.InputRequestMap{
			"confirm": &mcp.ElicitParams{
				Message: "Start streaming on this OBS instance?",
				RequestedSchema: &jsonschema.Schema{
					Type: "object",
					Properties: map[string]*jsonschema.Schema{
						"confirm": {Type: "boolean"},
					},
				},
			},
		},
		RequestState: "elicit-start-stream",
	}
}

func acceptedElicit(req *mcp.CallToolRequest, key string) error {
	raw, ok := req.Params.InputResponses[key]
	if !ok {
		return fmt.Errorf("StartStream cancelled: missing elicitation response")
	}
	res, ok := raw.(*mcp.ElicitResult)
	if !ok {
		return fmt.Errorf("StartStream cancelled: unexpected elicitation response")
	}
	if res.Action != "accept" {
		return fmt.Errorf("StartStream cancelled")
	}
	return nil
}
