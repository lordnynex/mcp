package tools

import (
	"context"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/lordnynex/mcp/robotgo/desktop"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type processListOutput struct {
	Processes []desktop.Process `json:"processes"`
}

type findProcessInput struct {
	Name string `json:"name" jsonschema:"Substring matched by robotgo FindIds (case insensitive)"`
}

type killInput struct {
	PID int `json:"pid" jsonschema:"Process ID to kill"`
}

type killOutput struct {
	PID int `json:"pid"`
}

func registerProcess(s *mcp.Server, d desktop.Driver) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "ListProcesses",
		Title:       humanTitle("ListProcesses"),
		Description: "List running processes (pid and name).",
		Annotations: annotations("ListProcesses"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, *processListOutput, error) {
		list, err := d.Processes()
		if err != nil {
			return nil, nil, err
		}
		return nil, &processListOutput{Processes: list}, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "FindProcess",
		Title:       humanTitle("FindProcess"),
		Description: "Find processes whose names contain the given substring.",
		Annotations: annotations("FindProcess"),
	}, func(_ context.Context, _ *mcp.CallToolRequest, in findProcessInput) (*mcp.CallToolResult, *processListOutput, error) {
		if in.Name == "" {
			return nil, nil, fmt.Errorf("name is required")
		}
		list, err := d.FindProcess(in.Name)
		if err != nil {
			return nil, nil, err
		}
		return nil, &processListOutput{Processes: list}, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "KillProcess",
		Title:       humanTitle("KillProcess"),
		Description: "Kill a process by pid. Asks for confirmation via elicitation before sending the signal.",
		Annotations: annotations("KillProcess"),
	}, func(_ context.Context, req *mcp.CallToolRequest, in killInput) (*mcp.CallToolResult, *killOutput, error) {
		if in.PID <= 0 {
			return nil, nil, fmt.Errorf("pid must be a positive integer")
		}
		if req == nil || req.Params == nil || len(req.Params.InputResponses) == 0 {
			return killElicit(in.PID), nil, nil
		}
		if err := acceptedElicit(req, "confirm"); err != nil {
			return nil, nil, err
		}
		if err := d.Kill(in.PID); err != nil {
			return nil, nil, err
		}
		return nil, &killOutput{PID: in.PID}, nil
	})
}

func killElicit(pid int) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		InputRequests: mcp.InputRequestMap{
			"confirm": &mcp.ElicitParams{
				Message: fmt.Sprintf("Kill process %d on this machine?", pid),
				RequestedSchema: &jsonschema.Schema{
					Type: "object",
					Properties: map[string]*jsonschema.Schema{
						"confirm": {Type: "boolean"},
					},
				},
			},
		},
		RequestState: "elicit-kill-process",
	}
}

func acceptedElicit(req *mcp.CallToolRequest, key string) error {
	raw, ok := req.Params.InputResponses[key]
	if !ok {
		return fmt.Errorf("KillProcess cancelled: missing elicitation response")
	}
	res, ok := raw.(*mcp.ElicitResult)
	if !ok {
		return fmt.Errorf("KillProcess cancelled: unexpected elicitation response")
	}
	if res.Action != "accept" {
		return fmt.Errorf("KillProcess cancelled")
	}
	return nil
}
