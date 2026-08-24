package obs

import (
	"testing"

	"github.com/lordnynex/gest"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestServer(t *testing.T) {
	gest.Run(t, "Server", func(s *gest.S) {
		s.Describe("New", func(s *gest.S) {
			s.It("returns a server", func(t *gest.T) {
				t.Expect(New()).NotTo(gest.BeNil())
			})
		})

		s.Describe("Register", func(s *gest.S) {
			s.It("can be called more than once", func(t *gest.T) {
				srv := mcp.NewServer(&mcp.Implementation{Name: "aggregate", Version: "0.0.0"}, nil)
				t.Expect(func() {
					Register(srv)
					Register(srv)
				}).NotTo(gest.Panic())
			})
		})
	})
}
