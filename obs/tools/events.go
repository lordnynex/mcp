package tools

import (
	"context"

	"github.com/lordnynex/mcp/obs/protocol"
	"github.com/lordnynex/mcp/obs/session"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type subscribeEventsInput struct {
	EventTypes []string `json:"eventTypes" jsonschema:"Official obs-websocket event names and/or EventSubscription category names (General, Config, Scenes, Inputs, Transitions, Filters, Outputs, SceneItems, MediaInputs, Vendors, Ui, Canvases)"`
}

type subscribeEventsResult struct {
	EventTypes []string `json:"eventTypes"`
	URIs       []string `json:"uris"`
	Note       string   `json:"note"`
}

type unsubscribeEventsInput struct {
	EventTypes []string `json:"eventTypes,omitempty" jsonschema:"Event types to drop; omit to drop all subscriptions for this session"`
}

type unsubscribeEventsResult struct {
	EventTypes []string `json:"eventTypes"`
}

func registerEventTools(s *mcp.Server, h *session.Host) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "SubscribeEvents",
		Title:       "Subscribe events",
		Description: "Select official OBS websocket event types to receive. Also call resources/subscribe (or subscriptions/listen) on the returned URIs so notifications/resources/updated is delivered on the MCP event-source stream.",
		Annotations: annotationsFor("SubscribeEvents"),
	}, func(ctx context.Context, req *mcp.CallToolRequest, in subscribeEventsInput) (*mcp.CallToolResult, subscribeEventsResult, error) {
		types, err := protocol.ExpandEventTypes(in.EventTypes)
		if err != nil {
			return nil, subscribeEventsResult{}, err
		}
		if err := h.EnsureHighVolume(types); err != nil {
			return nil, subscribeEventsResult{}, err
		}
		if req != nil && req.Session != nil {
			h.SetInterest(req.Session, types, true)
		}
		uris := make([]string, len(types))
		for i, n := range types {
			uris[i] = protocol.EventResourceURI(n)
		}
		return nil, subscribeEventsResult{
			EventTypes: types,
			URIs:       uris,
			Note:       "Call resources/subscribe (legacy) or subscriptions/listen with these URIs so notifications/resources/updated is delivered.",
		}, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "UnsubscribeEvents",
		Title:       "Unsubscribe events",
		Description: "Stop selecting OBS websocket event types for this MCP session.",
		Annotations: annotationsFor("UnsubscribeEvents"),
	}, func(ctx context.Context, req *mcp.CallToolRequest, in unsubscribeEventsInput) (*mcp.CallToolResult, unsubscribeEventsResult, error) {
		var types []string
		if len(in.EventTypes) == 0 {
			if req != nil && req.Session != nil {
				types = h.ClearInterest(req.Session)
			}
		} else {
			var err error
			types, err = protocol.ExpandEventTypes(in.EventTypes)
			if err != nil {
				return nil, unsubscribeEventsResult{}, err
			}
			if req != nil && req.Session != nil {
				h.SetInterest(req.Session, types, false)
			}
		}
		return nil, unsubscribeEventsResult{EventTypes: types}, nil
	})
}
