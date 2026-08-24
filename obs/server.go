package obs

import "github.com/modelcontextprotocol/go-sdk/mcp"

const (
	Name    = "obs"
	Version = "0.1.0"
)

// New returns a standalone OBS MCP server with tools registered.
func New() *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: Name, Version: Version}, nil)
	Register(s)
	return s
}

// Register mounts OBS tools onto an existing MCP server so this package can
// run standalone or be composed into an aggregate server.
func Register(s *mcp.Server) {
	// Tools added in a later plan.
}
