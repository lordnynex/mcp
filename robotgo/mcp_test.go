package robotgo

import (
	"context"
	"testing"

	"github.com/lordnynex/gest"
	"github.com/lordnynex/mcp/robotgo/complete"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestPrompts(t *testing.T) {
	gest.Run(t, "Prompts", func(s *gest.S) {
		s.It("robotgo-click-target requires target and names the tools", func(t *gest.T) {
			sess := connectClient(t, New(), nil)
			_, err := sess.GetPrompt(t.Context(), &mcp.GetPromptParams{Name: "robotgo-click-target"})
			t.Expect(err).NotTo(gest.BeNil())

			res, err := sess.GetPrompt(t.Context(), &mcp.GetPromptParams{
				Name:      "robotgo-click-target",
				Arguments: map[string]string{"target": "Save"},
			})
			t.Require(err).To(gest.BeNil())
			text := promptText(res)
			t.Expect(text).To(gest.Contain("Save"))
			t.Expect(text).To(gest.Contain("ObserveDesktop"))
			t.Expect(text).To(gest.Contain("MouseClick"))
		})

		s.It("robotgo-type-into-field requires text", func(t *gest.T) {
			sess := connectClient(t, New(), nil)
			_, err := sess.GetPrompt(t.Context(), &mcp.GetPromptParams{Name: "robotgo-type-into-field"})
			t.Expect(err).NotTo(gest.BeNil())

			res, err := sess.GetPrompt(t.Context(), &mcp.GetPromptParams{
				Name: "robotgo-type-into-field",
				Arguments: map[string]string{
					"text":  "hello",
					"field": "Search",
				},
			})
			t.Require(err).To(gest.BeNil())
			text := promptText(res)
			t.Expect(text).To(gest.Contain("hello"))
			t.Expect(text).To(gest.Contain("Search"))
			t.Expect(text).To(gest.Contain("Type"))
		})

		s.It("robotgo-desktop-ops lists allowed ops", func(t *gest.T) {
			sess := connectClient(t, New(), nil)
			res, err := sess.GetPrompt(t.Context(), &mcp.GetPromptParams{Name: "robotgo-desktop-ops"})
			t.Require(err).To(gest.BeNil())
			text := promptText(res)
			t.Expect(text).To(gest.Contain("DesktopOps"))
			t.Expect(text).To(gest.Contain("MouseSelect"))
			t.Expect(text).To(gest.Contain("FocusWindow"))
			t.Expect(text).To(gest.Contain("KillProcess is not allowed"))
			t.Expect(text).To(gest.Contain("MousePath"))
			t.Expect(text).To(gest.Contain("smooth false"))
		})

		s.It("robotgo-backend-limits mentions libei and tags", func(t *gest.T) {
			sess := connectClient(t, New(), nil)
			res, err := sess.GetPrompt(t.Context(), &mcp.GetPromptParams{Name: "robotgo-backend-limits"})
			t.Require(err).To(gest.BeNil())
			text := promptText(res)
			t.Expect(text).To(gest.Contain("libei"))
			t.Expect(text).To(gest.Contain("purego,x11"))
		})
	})
}

func TestCompletions(t *testing.T) {
	gest.Run(t, "Completions", func(s *gest.S) {
		s.It("completes key names", func(t *gest.T) {
			res, err := complete.Handle(t.Context(), &mcp.CompleteRequest{
				Params: &mcp.CompleteParams{
					Ref:      &mcp.CompleteReference{Type: "ref/prompt", Name: "robotgo-shortcut"},
					Argument: mcp.CompleteParamsArgument{Name: "key", Value: "ent"},
				},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(res.Completion.Values).To(gest.Contain("enter"))
			t.Expect(res.Completion.Values).NotTo(gest.Contain("escape"))
		})

		s.It("completes scroll directions", func(t *gest.T) {
			res, err := complete.Handle(t.Context(), &mcp.CompleteRequest{
				Params: &mcp.CompleteParams{
					Ref:      &mcp.CompleteReference{Type: "ref/prompt", Name: "robotgo-scroll-read"},
					Argument: mcp.CompleteParamsArgument{Name: "dir", Value: "d"},
				},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(res.Completion.Values).To(gest.Contain("down"))
			t.Expect(res.Completion.Values).NotTo(gest.Contain("up"))
		})

		s.It("completes DesktopOps op names", func(t *gest.T) {
			res, err := complete.Handle(t.Context(), &mcp.CompleteRequest{
				Params: &mcp.CompleteParams{
					Ref:      &mcp.CompleteReference{Type: "ref/prompt", Name: "robotgo-desktop-ops"},
					Argument: mcp.CompleteParamsArgument{Name: "op", Value: "MouseS"},
				},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(res.Completion.Values).To(gest.Contain("MouseSelect"))
			t.Expect(res.Completion.Values).To(gest.Contain("MouseScroll"))
			t.Expect(res.Completion.Values).NotTo(gest.Contain("MouseMove"))
		})

		s.It("completes skill resource names", func(t *gest.T) {
			res, err := complete.Handle(t.Context(), &mcp.CompleteRequest{
				Params: &mcp.CompleteParams{
					Ref:      &mcp.CompleteReference{Type: "ref/resource", URI: "robotgo://skills/{name}"},
					Argument: mcp.CompleteParamsArgument{Name: "name", Value: "key"},
				},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(res.Completion.Values).To(gest.Contain("keyboard"))
			t.Expect(res.Completion.Values).NotTo(gest.Contain("mouse"))
		})
	})
}

func TestElicitation(t *testing.T) {
	gest.Run(t, "Elicitation", func(s *gest.S) {
		s.It("elicits KillProcess and does not kill when declined", func(t *gest.T) {
			fake := withFake(t)
			var messages []string
			sess := connectClient(t, New(), &mcp.ClientOptions{
				ElicitationHandler: func(_ context.Context, req *mcp.ElicitRequest) (*mcp.ElicitResult, error) {
					if req != nil && req.Params != nil {
						messages = append(messages, req.Params.Message)
					}
					return &mcp.ElicitResult{Action: "decline"}, nil
				},
			})
			res, err := sess.CallTool(t.Context(), &mcp.CallToolParams{
				Name:      "KillProcess",
				Arguments: map[string]any{"pid": 42},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(res.IsError).To(gest.Equal(true))
			t.Expect(messages).To(gest.Contain("Kill process 42 on this machine?"))
			t.Expect(fake.Called("Kill")).To(gest.Equal(false))
		})

		s.It("kills after elicitation is accepted", func(t *gest.T) {
			fake := withFake(t)
			sess := connectClient(t, New(), &mcp.ClientOptions{
				ElicitationHandler: func(_ context.Context, _ *mcp.ElicitRequest) (*mcp.ElicitResult, error) {
					return &mcp.ElicitResult{Action: "accept"}, nil
				},
			})
			res, err := sess.CallTool(t.Context(), &mcp.CallToolParams{
				Name:      "KillProcess",
				Arguments: map[string]any{"pid": 42},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(res.IsError).To(gest.Equal(false))
			t.Expect(fake.Killed).To(gest.Equal([]int{42}))
		})
	})
}

func TestSkillResource(t *testing.T) {
	gest.Run(t, "SkillResource", func(s *gest.S) {
		s.It("reads observe-then-act markdown", func(t *gest.T) {
			sess := connectClient(t, New(), nil)
			res, err := sess.ReadResource(t.Context(), &mcp.ReadResourceParams{URI: "robotgo://skills/observe-then-act"})
			t.Require(err).To(gest.BeNil())
			t.Require(len(res.Contents) > 0).To(gest.Equal(true))
			t.Expect(res.Contents[0].Text).To(gest.Contain("ObserveDesktop"))
			t.Expect(res.Contents[0].Text).To(gest.Contain("DesktopOps"))
			t.Expect(res.Contents[0].Text).To(gest.Contain("mouseWidth"))
		})

		s.It("reads aim-center markdown", func(t *gest.T) {
			sess := connectClient(t, New(), nil)
			res, err := sess.ReadResource(t.Context(), &mcp.ReadResourceParams{URI: "robotgo://skills/aim-center"})
			t.Require(err).To(gest.BeNil())
			t.Require(len(res.Contents) > 0).To(gest.Equal(true))
			t.Expect(res.Contents[0].Text).To(gest.Contain("geometric center"))
			t.Expect(res.Contents[0].Text).To(gest.Contain("MouseSelect"))
			t.Expect(res.Contents[0].Text).To(gest.Contain("mouseWidth"))
		})
	})
}

func promptText(res *mcp.GetPromptResult) string {
	if res == nil || len(res.Messages) == 0 || res.Messages[0] == nil {
		return ""
	}
	text, _ := res.Messages[0].Content.(*mcp.TextContent)
	if text == nil {
		return ""
	}
	return text.Text
}
