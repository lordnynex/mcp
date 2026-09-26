package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/lordnynex/mcp/robotgo/skills"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Register adds robotgo://skills and robotgo://skills/{name}.
func Register(s *mcp.Server) {
	s.AddResource(&mcp.Resource{
		URI:         skills.IndexURI,
		Name:        "robotgo skills",
		Title:       "robotgo skills",
		Description: "Index of computer-use skills for calling agents.",
		MIMEType:    "text/plain",
	}, handle)

	s.AddResourceTemplate(&mcp.ResourceTemplate{
		URITemplate: skills.URITemplate,
		Name:        "robotgo skill",
		Title:       "robotgo skill",
		Description: "SKILL.md for a named computer-use workflow.",
		MIMEType:    "text/markdown",
	}, handle)
}

func handle(_ context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	uri := ""
	if req != nil && req.Params != nil {
		uri = req.Params.URI
	}
	var (
		body string
		mime string
		err  error
	)
	switch {
	case uri == skills.IndexURI:
		body = strings.Join(skills.Names, "\n") + "\n"
		mime = "text/plain"
	case strings.HasPrefix(uri, skills.URIPrefix):
		name := strings.TrimPrefix(uri, skills.URIPrefix)
		body, err = skills.Read(name)
		if err != nil {
			return nil, mcp.ResourceNotFoundError(uri)
		}
		mime = "text/markdown"
	default:
		return nil, mcp.ResourceNotFoundError(uri)
	}
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{{
			URI:      uri,
			MIMEType: mime,
			Text:     body,
		}},
	}, nil
}

// Subscribe accepts skill resource URIs.
func Subscribe(_ context.Context, req *mcp.SubscribeRequest) error {
	if req == nil || req.Params == nil {
		return fmt.Errorf("missing subscribe params")
	}
	return validateSkillURI(req.Params.URI)
}

// Unsubscribe accepts skill resource URIs.
func Unsubscribe(_ context.Context, req *mcp.UnsubscribeRequest) error {
	if req == nil || req.Params == nil {
		return fmt.Errorf("missing unsubscribe params")
	}
	return validateSkillURI(req.Params.URI)
}

func validateSkillURI(uri string) error {
	if uri == skills.IndexURI {
		return nil
	}
	if strings.HasPrefix(uri, skills.URIPrefix) {
		name := strings.TrimPrefix(uri, skills.URIPrefix)
		if _, err := skills.Read(name); err == nil {
			return nil
		}
	}
	return fmt.Errorf("unknown skill URI %q", uri)
}
