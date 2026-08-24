package obs

import (
	"context"
	"sync"
	"testing"

	"github.com/lordnynex/gest"
	"github.com/lordnynex/mcp/obs/complete"
	"github.com/lordnynex/mcp/obs/prompts"
	"github.com/lordnynex/mcp/obs/session"
	"github.com/lordnynex/mcp/obs/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestPrompts(t *testing.T) {
	gest.Run(t, "Prompts", func(s *gest.S) {
		s.It("lists always-on prompts before connect", func(t *gest.T) {
			srv := New()
			sess := connectClient(t, srv, nil)
			res, err := sess.ListPrompts(t.Context(), nil)
			t.Require(err).To(gest.BeNil())
			names := promptNames(res)
			for _, n := range prompts.AlwaysNames {
				t.Expect(names).To(gest.Contain(n))
			}
			t.Expect(names).NotTo(gest.Contain("obs-current-program"))
		})

		s.It("adds obs-current-program after RegisterConnected", func(t *gest.T) {
			srv := New()
			prompts.RegisterConnected(srv, session.Default())
			sess := connectClient(t, srv, nil)
			res, err := sess.ListPrompts(t.Context(), nil)
			t.Require(err).To(gest.BeNil())
			t.Expect(promptNames(res)).To(gest.Contain("obs-current-program"))
		})
	})
}

func TestCompletions(t *testing.T) {
	gest.Run(t, "Completions", func(s *gest.S) {
		s.It("completes event subscription categories", func(t *gest.T) {
			res, err := complete.HandleHost(t.Context(), session.New(), &mcp.CompleteRequest{
				Params: &mcp.CompleteParams{
					Ref:      &mcp.CompleteReference{Type: "ref/prompt", Name: "obs-subscribe-events"},
					Argument: mcp.CompleteParamsArgument{Name: "category", Value: "Sc"},
				},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(res.Completion.Values).To(gest.Contain("Scenes"))
			t.Expect(res.Completion.Values).NotTo(gest.Contain("Inputs"))
		})

		s.It("completes event resource type names", func(t *gest.T) {
			res, err := complete.HandleHost(t.Context(), session.New(), &mcp.CompleteRequest{
				Params: &mcp.CompleteParams{
					Ref:      &mcp.CompleteReference{Type: "ref/resource", URI: "obs://events/{eventType}"},
					Argument: mcp.CompleteParamsArgument{Name: "eventType", Value: "CurrentProgram"},
				},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(res.Completion.Values).To(gest.Contain("CurrentProgramSceneChanged"))
		})

		s.It("returns no scene names when OBS is disconnected", func(t *gest.T) {
			res, err := complete.HandleHost(t.Context(), session.New(), &mcp.CompleteRequest{
				Params: &mcp.CompleteParams{
					Ref:      &mcp.CompleteReference{Type: "ref/prompt", Name: "obs-switch-scene"},
					Argument: mcp.CompleteParamsArgument{Name: "sceneName", Value: ""},
				},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(res.Completion.Values).To(gest.Equal([]string{}))
		})
	})
}

func TestPagination(t *testing.T) {
	gest.Run(t, "Pagination", func(s *gest.S) {
		s.It("returns a next cursor when connected tools exceed page size", func(t *gest.T) {
			SetMCPDefaults(MCPDefaults{PageSize: 50})
			srv := New()
			tools.RegisterConnected(srv, session.Default())
			sess := connectClient(t, srv, nil)
			res, err := sess.ListTools(t.Context(), nil)
			t.Require(err).To(gest.BeNil())
			t.Expect(res.NextCursor).NotTo(gest.Equal(""))
			t.Expect(len(res.Tools) <= 50).To(gest.Equal(true))
		})
	})
}

func TestElicitation(t *testing.T) {
	gest.Run(t, "Elicitation", func(s *gest.S) {
		s.It("elicits StartStream and does not elicit StartRecord", func(t *gest.T) {
			var messages []string
			srv := mcp.NewServer(&mcp.Implementation{Name: "obs-test", Version: "0"}, &mcp.ServerOptions{PageSize: 1000})
			h := session.New()
			tools.RegisterConnected(srv, h)
			sess := connectClient(t, srv, &mcp.ClientOptions{
				ElicitationHandler: func(_ context.Context, req *mcp.ElicitRequest) (*mcp.ElicitResult, error) {
					if req != nil && req.Params != nil {
						messages = append(messages, req.Params.Message)
					}
					return &mcp.ElicitResult{Action: "decline"}, nil
				},
			})

			streamRes, err := sess.CallTool(t.Context(), &mcp.CallToolParams{Name: "StartStream"})
			t.Require(err).To(gest.BeNil())
			t.Expect(streamRes.IsError).To(gest.Equal(true))
			t.Expect(messages).To(gest.Contain("Start streaming on this OBS instance?"))

			n := len(messages)
			recordRes, err := sess.CallTool(t.Context(), &mcp.CallToolParams{Name: "StartRecord"})
			t.Require(err).To(gest.BeNil())
			t.Expect(recordRes.IsError).To(gest.Equal(true))
			t.Expect(len(messages)).To(gest.Equal(n))
		})
	})
}

func TestProgress(t *testing.T) {
	gest.Run(t, "Progress", func(s *gest.S) {
		s.It("notifies progress during Connect", func(t *gest.T) {
			var mu sync.Mutex
			var messages []string
			srv := New()
			sess := connectClient(t, srv, &mcp.ClientOptions{
				ProgressNotificationHandler: func(_ context.Context, req *mcp.ProgressNotificationClientRequest) {
					if req == nil || req.Params == nil {
						return
					}
					mu.Lock()
					messages = append(messages, req.Params.Message)
					mu.Unlock()
				},
			})
			params := &mcp.CallToolParams{
				Name: "Connect",
				Arguments: map[string]any{
					"host": "127.0.0.1:1",
				},
			}
			params.SetProgressToken("connect-progress")
			_, err := sess.CallTool(t.Context(), params)
			t.Expect(err).To(gest.BeNil())
			mu.Lock()
			defer mu.Unlock()
			t.Expect(messages).To(gest.Contain("identifying with OBS"))
		})
	})
}

func promptNames(res *mcp.ListPromptsResult) []string {
	if res == nil {
		return nil
	}
	names := make([]string, 0, len(res.Prompts))
	for _, p := range res.Prompts {
		names = append(names, p.Name)
	}
	return names
}
