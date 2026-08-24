package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lordnynex/mcp/obs/protocol"
	"github.com/lordnynex/mcp/obs/session"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Register adds the obs://events resource and obs://events/{eventType} template.
func Register(s *mcp.Server, h *session.Host) {
	s.AddResource(&mcp.Resource{
		URI:         protocol.EventsResourceURI,
		Name:        "OBS events",
		Title:       "OBS events",
		Description: "Recent OBS websocket events (ring buffer).",
		MIMEType:    "application/json",
	}, handle(h))

	s.AddResourceTemplate(&mcp.ResourceTemplate{
		URITemplate: protocol.EventURITemplate,
		Name:        "OBS event",
		Title:       "OBS event",
		Description: "Latest payload for an official OBS websocket event type.",
		MIMEType:    "application/json",
	}, handle(h))
}

func handle(h *session.Host) mcp.ResourceHandler {
	return func(_ context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		uri := ""
		if req != nil && req.Params != nil {
			uri = req.Params.URI
		}
		var body []byte
		var err error
		switch {
		case uri == protocol.EventsResourceURI:
			body, err = json.Marshal(h.EventBuffer())
		case strings.HasPrefix(uri, protocol.EventURIPrefix):
			name := strings.TrimPrefix(uri, protocol.EventURIPrefix)
			if !protocol.IsKnownEvent(name) {
				return nil, mcp.ResourceNotFoundError(uri)
			}
			ev, ok := h.LatestEvent(name)
			if !ok {
				body = []byte("null")
				break
			}
			body, err = json.Marshal(ev)
		default:
			return nil, mcp.ResourceNotFoundError(uri)
		}
		if err != nil {
			return nil, err
		}
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{{
				URI:      uri,
				MIMEType: "application/json",
				Text:     string(body),
			}},
		}, nil
	}
}

// Subscribe is the MCP resources/subscribe handler.
func Subscribe(_ context.Context, req *mcp.SubscribeRequest) error {
	if req == nil || req.Params == nil {
		return fmt.Errorf("missing subscribe params")
	}
	return validateEventURI(req.Params.URI)
}

// Unsubscribe is the MCP resources/unsubscribe handler.
func Unsubscribe(_ context.Context, req *mcp.UnsubscribeRequest) error {
	if req == nil || req.Params == nil {
		return fmt.Errorf("missing unsubscribe params")
	}
	return validateEventURI(req.Params.URI)
}

func validateEventURI(uri string) error {
	if uri == protocol.EventsResourceURI {
		return nil
	}
	if strings.HasPrefix(uri, protocol.EventURIPrefix) {
		name := strings.TrimPrefix(uri, protocol.EventURIPrefix)
		if protocol.IsKnownEvent(name) {
			return nil
		}
		return fmt.Errorf("unknown OBS event type %q", name)
	}
	return fmt.Errorf("unsupported resource URI %q", uri)
}
