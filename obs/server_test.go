package obs

import (
	"slices"
	"testing"

	"github.com/lordnynex/gest"
	"github.com/lordnynex/mcp/obs/protocol"
	"github.com/lordnynex/mcp/obs/session"
	"github.com/lordnynex/mcp/obs/tools"
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

func TestSessionTools(t *testing.T) {
	gest.Run(t, "SessionTools", func(s *gest.S) {
		s.It("lists only Connect and ConnectionStatus before OBS connect", func(t *gest.T) {
			srv := New()
			names := listAllToolNames(t, srv)
			t.Expect(names).To(gest.Equal([]string{"Connect", "ConnectionStatus"}))
			for _, req := range protocol.RequestNames {
				t.Expect(names).NotTo(gest.Contain(req))
			}
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

func TestProtocolTools(t *testing.T) {
	gest.Run(t, "ProtocolTools", func(s *gest.S) {
		s.It("registers every official request after protocol tools are added", func(t *gest.T) {
			srv := mcp.NewServer(&mcp.Implementation{Name: "obs-test", Version: "0"}, &mcp.ServerOptions{PageSize: 1000})
			h := session.New()
			tools.RegisterAlways(srv, h)
			tools.RegisterConnected(srv, h)
			names := listAllToolNames(t, srv)
			for _, req := range protocol.RequestNames {
				t.Expect(names).To(gest.Contain(req))
			}
			for _, extra := range []string{"Connect", "ConnectionStatus", "Disconnect", "RequestBatch", "SubscribeEvents", "UnsubscribeEvents"} {
				t.Expect(names).To(gest.Contain(extra))
			}
		})

		s.It("does not register invented request names", func(t *gest.T) {
			srv := mcp.NewServer(&mcp.Implementation{Name: "obs-test", Version: "0"}, &mcp.ServerOptions{PageSize: 1000})
			h := session.New()
			tools.RegisterConnected(srv, h)
			names := listAllToolNames(t, srv)
			officialSet := make(map[string]struct{}, len(protocol.RequestNames))
			for _, n := range protocol.RequestNames {
				officialSet[n] = struct{}{}
			}
			sessionTools := map[string]struct{}{
				"Disconnect": {}, "RequestBatch": {}, "SubscribeEvents": {}, "UnsubscribeEvents": {},
			}
			for _, n := range names {
				if _, ok := sessionTools[n]; ok {
					continue
				}
				_, ok := officialSet[n]
				t.Expect(ok).To(gest.Equal(true))
			}
		})
	})
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
