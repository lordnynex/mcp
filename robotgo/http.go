package robotgo

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const defaultHTTPSessionTimeout = 30 * time.Minute

// HTTPOptions configure streamable HTTP for a robotgo MCP server.
type HTTPOptions struct {
	Stateless      bool
	JSONResponse   bool
	SessionTimeout time.Duration
	BearerToken    string
}

// NewHTTPHandler wraps s in a streamable HTTP handler.
func NewHTTPHandler(s *mcp.Server, o HTTPOptions) http.Handler {
	opts := &mcp.StreamableHTTPOptions{
		Stateless:                    o.Stateless,
		JSONResponse:                 o.JSONResponse,
		Logger:                       slog.Default(),
		SessionTimeout:               o.SessionTimeout,
		PropagateRequestCancellation: true,
	}
	if !o.Stateless {
		opts.EventStore = mcp.NewMemoryEventStore(nil)
	}
	h := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return s
	}, opts)
	if o.BearerToken == "" {
		return h
	}
	want := o.BearerToken
	return auth.RequireBearerToken(func(_ context.Context, token string, _ *http.Request) (*auth.TokenInfo, error) {
		if subtle.ConstantTimeCompare([]byte(token), []byte(want)) != 1 {
			return nil, auth.ErrInvalidToken
		}
		return &auth.TokenInfo{Expiration: time.Now().Add(24 * time.Hour)}, nil
	}, nil)(h)
}

// ListenAndServeHTTP serves handler on addr until ctx is cancelled.
func ListenAndServeHTTP(ctx context.Context, addr string, handler http.Handler) error {
	if addr == "" {
		return fmt.Errorf("http listen address is empty")
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	srv := &http.Server{Handler: handler, BaseContext: func(net.Listener) context.Context { return ctx }}
	go func() {
		<-ctx.Done()
		shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shCtx)
	}()
	err = srv.Serve(ln)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

// DefaultHTTPSessionTimeout is used when the CLI flag is left at its default.
func DefaultHTTPSessionTimeout() time.Duration {
	return defaultHTTPSessionTimeout
}
