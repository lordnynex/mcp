package robotgo

import (
	"slices"
	"testing"

	"github.com/lordnynex/gest"
	"github.com/lordnynex/mcp/robotgo/desktop"
	"github.com/lordnynex/mcp/robotgo/prompts"
	"github.com/lordnynex/mcp/robotgo/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestServer(t *testing.T) {
	gest.Run(t, "Server", func(s *gest.S) {
		s.Describe("New", func(s *gest.S) {
			s.It("returns a server", func(t *gest.T) {
				t.Expect(New()).NotTo(gest.BeNil())
			})
		})
	})
}

func TestRegister(t *testing.T) {
	gest.Run(t, "Register", func(s *gest.S) {
		s.It("can be called more than once on the same server", func(t *gest.T) {
			srv := mcp.NewServer(&mcp.Implementation{Name: "aggregate", Version: "0.0.0"}, nil)
			t.Expect(func() {
				Register(srv)
				Register(srv)
			}).NotTo(gest.Panic())
		})
	})
}

func TestTools(t *testing.T) {
	gest.Run(t, "Tools", func(s *gest.S) {
		s.It("registers every planned tool", func(t *gest.T) {
			srv := New()
			names := listAllToolNames(t, srv)
			for _, want := range tools.ToolNames {
				t.Expect(names).To(gest.Contain(want))
			}
			t.Expect(len(names)).To(gest.Equal(len(tools.ToolNames)))
		})

		s.It("does not register invented tool names", func(t *gest.T) {
			srv := New()
			names := listAllToolNames(t, srv)
			official := make(map[string]struct{}, len(tools.ToolNames))
			for _, n := range tools.ToolNames {
				official[n] = struct{}{}
			}
			for _, n := range names {
				_, ok := official[n]
				t.Expect(ok).To(gest.Equal(true))
			}
		})
	})
}

func TestObserveDesktop(t *testing.T) {
	gest.Run(t, "ObserveDesktop", func(s *gest.S) {
		s.It("returns image content and cursor metadata", func(t *gest.T) {
			fake := withFake(t)
			fake.CursorX = 12
			fake.CursorY = 34
			sess := connectClient(t, New(), nil)
			res, err := sess.CallTool(t.Context(), &mcp.CallToolParams{Name: "ObserveDesktop"})
			t.Require(err).To(gest.BeNil())
			t.Expect(res.IsError).To(gest.Equal(false))
			t.Require(len(res.Content) > 0).To(gest.Equal(true))
			img, ok := res.Content[0].(*mcp.ImageContent)
			t.Require(ok).To(gest.Equal(true))
			t.Expect(img.MIMEType).To(gest.Equal("image/png"))
			t.Expect(len(img.Data) > 0).To(gest.Equal(true))
			t.Expect(fake.Called("Capture")).To(gest.Equal(true))
		})

		s.It("returns image-to-mouse mapping fields", func(t *gest.T) {
			withFake(t)
			sess := connectClient(t, New(), nil)
			res, err := sess.CallTool(t.Context(), &mcp.CallToolParams{Name: "ObserveDesktop"})
			t.Require(err).To(gest.BeNil())
			t.Expect(res.IsError).To(gest.Equal(false))
			m, ok := res.StructuredContent.(map[string]any)
			t.Require(ok).To(gest.Equal(true))
			t.Expect(m["mouseWidth"]).To(gest.Equal(float64(1920)))
			t.Expect(m["mouseHeight"]).To(gest.Equal(float64(1080)))
			t.Expect(m["originX"]).To(gest.Equal(float64(0)))
			t.Expect(m["width"]).To(gest.Equal(float64(4)))
			t.Expect(m["height"]).To(gest.Equal(float64(4)))
		})

		s.It("rejects a partial region", func(t *gest.T) {
			withFake(t)
			sess := connectClient(t, New(), nil)
			res, err := sess.CallTool(t.Context(), &mcp.CallToolParams{
				Name: "ObserveDesktop",
				Arguments: map[string]any{
					"x": 1,
					"y": 2,
				},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(res.IsError).To(gest.Equal(true))
		})
	})
}

func TestMouseAndKeys(t *testing.T) {
	gest.Run(t, "MouseAndKeys", func(s *gest.S) {
		s.It("records MouseMove and KeyTap arguments", func(t *gest.T) {
			fake := withFake(t)
			sess := connectClient(t, New(), nil)
			_, err := sess.CallTool(t.Context(), &mcp.CallToolParams{
				Name: "MouseMove",
				Arguments: map[string]any{
					"x": 100,
					"y": 200,
				},
			})
			t.Require(err).To(gest.BeNil())
			call, err := fake.Last("Move")
			t.Require(err).To(gest.BeNil())
			t.Expect(call.Args[0]).To(gest.Equal(100))
			t.Expect(call.Args[1]).To(gest.Equal(200))

			res, err := sess.CallTool(t.Context(), &mcp.CallToolParams{
				Name: "KeyTap",
				Arguments: map[string]any{
					"key":       "enter",
					"modifiers": []string{"cmd"},
				},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(res.IsError).To(gest.Equal(false))
			t.Expect(fake.Called("KeyTap")).To(gest.Equal(true))
		})

		s.It("rejects an unknown key", func(t *gest.T) {
			withFake(t)
			sess := connectClient(t, New(), nil)
			res, err := sess.CallTool(t.Context(), &mcp.CallToolParams{
				Name:      "KeyTap",
				Arguments: map[string]any{"key": "not-a-key"},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(res.IsError).To(gest.Equal(true))
		})

		s.It("rejects empty Type text", func(t *gest.T) {
			withFake(t)
			sess := connectClient(t, New(), nil)
			res, err := sess.CallTool(t.Context(), &mcp.CallToolParams{
				Name:      "Type",
				Arguments: map[string]any{"text": ""},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(res.IsError).To(gest.Equal(true))
		})

		s.It("rejects an unknown mouse button", func(t *gest.T) {
			withFake(t)
			sess := connectClient(t, New(), nil)
			res, err := sess.CallTool(t.Context(), &mcp.CallToolParams{
				Name:      "MouseClick",
				Arguments: map[string]any{"button": "middle-finger"},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(res.IsError).To(gest.Equal(true))
		})
	})
}

func TestDesktopOpsTool(t *testing.T) {
	gest.Run(t, "DesktopOpsTool", func(s *gest.S) {
		s.It("runs a mixed batch through MCP", func(t *gest.T) {
			fake := withFake(t)
			sess := connectClient(t, New(), nil)
			res, err := sess.CallTool(t.Context(), &mcp.CallToolParams{
				Name: "DesktopOps",
				Arguments: map[string]any{
					"ops": []any{
						map[string]any{"op": "MouseMove", "x": 10, "y": 20},
						map[string]any{"op": "MouseClick"},
						map[string]any{"op": "Type", "text": "hi"},
					},
				},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(res.IsError).To(gest.Equal(false))
			t.Expect(fake.Called("Move")).To(gest.Equal(true))
			t.Expect(fake.Called("Click")).To(gest.Equal(true))
			t.Expect(fake.Called("Type")).To(gest.Equal(true))
		})

		s.It("rejects KillProcess in a batch", func(t *gest.T) {
			fake := withFake(t)
			sess := connectClient(t, New(), nil)
			res, err := sess.CallTool(t.Context(), &mcp.CallToolParams{
				Name: "DesktopOps",
				Arguments: map[string]any{
					"ops": []any{map[string]any{"op": "KillProcess", "pid": 42}},
				},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(res.IsError).To(gest.Equal(true))
			t.Expect(fake.Called("Kill")).To(gest.Equal(false))
		})
	})
}

func TestMousePathTool(t *testing.T) {
	gest.Run(t, "MousePathTool", func(s *gest.S) {
		s.It("moves through each point", func(t *gest.T) {
			fake := withFake(t)
			sess := connectClient(t, New(), nil)
			res, err := sess.CallTool(t.Context(), &mcp.CallToolParams{
				Name: "MousePath",
				Arguments: map[string]any{
					"points": []any{
						map[string]any{"x": 1, "y": 2},
						map[string]any{"x": 3, "y": 4},
					},
				},
			})
			t.Require(err).To(gest.BeNil())
			t.Expect(res.IsError).To(gest.Equal(false))
			n := 0
			for _, c := range fake.Calls {
				if c.Name == "Move" {
					n++
				}
			}
			t.Expect(n).To(gest.Equal(2))
		})
	})
}

func TestPromptsListed(t *testing.T) {
	gest.Run(t, "PromptsListed", func(s *gest.S) {
		s.It("lists every planned prompt", func(t *gest.T) {
			sess := connectClient(t, New(), nil)
			res, err := sess.ListPrompts(t.Context(), nil)
			t.Require(err).To(gest.BeNil())
			names := promptNames(res)
			for _, want := range prompts.Names {
				t.Expect(names).To(gest.Contain(want))
			}
		})
	})
}

func withFake(t *gest.T) *desktop.Fake {
	t.Helper()
	prev := desktop.Current()
	fake := desktop.NewFake()
	desktop.SetDefault(fake)
	t.Cleanup(func() { desktop.SetDefault(prev) })
	return fake
}

func listAllToolNames(t *gest.T, srv *mcp.Server) []string {
	t.Helper()
	sess := connectClient(t, srv, nil)
	var names []string
	var cursor string
	for {
		res, err := sess.ListTools(t.Context(), &mcp.ListToolsParams{Cursor: cursor})
		t.Require(err).To(gest.BeNil())
		for _, tool := range res.Tools {
			names = append(names, tool.Name)
		}
		if res.NextCursor == "" {
			break
		}
		cursor = res.NextCursor
	}
	slices.Sort(names)
	return names
}

func connectClient(t *gest.T, srv *mcp.Server, opts *mcp.ClientOptions) *mcp.ClientSession {
	t.Helper()
	ctx := t.Context()
	st, ct := mcp.NewInMemoryTransports()
	_, err := srv.Connect(ctx, st, nil)
	t.Require(err).To(gest.BeNil())
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "0"}, opts)
	sess, err := client.Connect(ctx, ct, nil)
	t.Require(err).To(gest.BeNil())
	return sess
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
