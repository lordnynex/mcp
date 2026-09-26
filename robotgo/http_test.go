package robotgo

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lordnynex/gest"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestHTTP(t *testing.T) {
	gest.Run(t, "HTTP", func(s *gest.S) {
		s.It("stateful initialize sets Mcp-Session-Id and GetScreenInfo works", func(t *gest.T) {
			withFake(t)
			httpSrv := httptest.NewServer(NewHTTPHandler(New(), HTTPOptions{JSONResponse: true}))
			t.Cleanup(httpSrv.Close)

			sid, status := postJSON(t, httpSrv.URL, "", initializeBody("2025-11-25"))
			t.Expect(status).To(gest.Equal(http.StatusOK))
			t.Expect(sid).NotTo(gest.Equal(""))

			client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "0"}, nil)
			sess, err := client.Connect(t.Context(), &mcp.StreamableClientTransport{Endpoint: httpSrv.URL}, nil)
			t.Require(err).To(gest.BeNil())
			t.Cleanup(func() { _ = sess.Close() })
			t.Expect(sess.ID()).NotTo(gest.Equal(""))

			res, err := sess.CallTool(t.Context(), &mcp.CallToolParams{Name: "GetScreenInfo"})
			t.Require(err).To(gest.BeNil())
			t.Expect(res.IsError).To(gest.Equal(false))
		})

		s.It("rejects a stateful GET without Mcp-Session-Id", func(t *gest.T) {
			httpSrv := httptest.NewServer(NewHTTPHandler(New(), HTTPOptions{JSONResponse: true}))
			t.Cleanup(httpSrv.Close)

			sid, status := postJSON(t, httpSrv.URL, "", initializeBody("2025-11-25"))
			t.Require(status).To(gest.Equal(http.StatusOK))
			t.Require(sid).NotTo(gest.Equal(""))

			req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, httpSrv.URL, nil)
			t.Require(err).To(gest.BeNil())
			req.Header.Set("Accept", "text/event-stream")
			resp, err := http.DefaultClient.Do(req)
			t.Require(err).To(gest.BeNil())
			defer resp.Body.Close()
			t.Expect(resp.StatusCode).To(gest.Equal(http.StatusBadRequest))
		})

		s.It("stateless HTTP does not require Mcp-Session-Id", func(t *gest.T) {
			withFake(t)
			httpSrv := httptest.NewServer(NewHTTPHandler(New(), HTTPOptions{Stateless: true, JSONResponse: true}))
			t.Cleanup(httpSrv.Close)

			sid, status := postJSON(t, httpSrv.URL, "", initializeBody("2025-11-25"))
			t.Expect(status).To(gest.Equal(http.StatusOK))
			t.Expect(sid).To(gest.Equal(""))

			client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "0"}, nil)
			sess, err := client.Connect(t.Context(), &mcp.StreamableClientTransport{Endpoint: httpSrv.URL}, nil)
			t.Require(err).To(gest.BeNil())
			t.Cleanup(func() { _ = sess.Close() })

			res, err := sess.CallTool(t.Context(), &mcp.CallToolParams{Name: "GetScreenInfo"})
			t.Require(err).To(gest.BeNil())
			t.Expect(res.IsError).To(gest.Equal(false))
		})

		s.It("returns 401 without a bearer token when one is configured", func(t *gest.T) {
			httpSrv := httptest.NewServer(NewHTTPHandler(New(), HTTPOptions{
				JSONResponse: true,
				BearerToken:  "secret",
			}))
			t.Cleanup(httpSrv.Close)

			req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, httpSrv.URL, bytes.NewReader(initializeBody("2025-11-25")))
			t.Require(err).To(gest.BeNil())
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json, text/event-stream")
			resp, err := http.DefaultClient.Do(req)
			t.Require(err).To(gest.BeNil())
			defer resp.Body.Close()
			t.Expect(resp.StatusCode).To(gest.Equal(http.StatusUnauthorized))
		})
	})
}

func initializeBody(version string) []byte {
	return []byte(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"` + version + `","capabilities":{},"clientInfo":{"name":"test","version":"0"}}}`)
}

func postJSON(t *gest.T, url, sessionID string, body []byte) (string, int) {
	t.Helper()
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, url, bytes.NewReader(body))
	t.Require(err).To(gest.BeNil())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	if sessionID != "" {
		req.Header.Set("Mcp-Session-Id", sessionID)
	}
	resp, err := http.DefaultClient.Do(req)
	t.Require(err).To(gest.BeNil())
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp.Header.Get("Mcp-Session-Id"), resp.StatusCode
}

func TestListenAndServeHTTP(t *testing.T) {
	gest.Run(t, "ListenAndServeHTTP", func(s *gest.S) {
		s.It("rejects an empty address", func(t *gest.T) {
			err := ListenAndServeHTTP(t.Context(), "", http.NotFoundHandler())
			t.Expect(err).NotTo(gest.BeNil())
		})
	})
}
