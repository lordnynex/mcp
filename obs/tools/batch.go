package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/andreykaipov/goobs"
	"github.com/lordnynex/mcp/obs/protocol"
	"github.com/lordnynex/mcp/obs/session"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	batchSerialRealtime = 0
	batchSerialFrame    = 1
	batchParallel       = 2
)

type batchInput struct {
	HaltOnFailure bool           `json:"haltOnFailure,omitempty" jsonschema:"Stop the batch on first failure (default false)"`
	ExecutionType *int           `json:"executionType,omitempty" jsonschema:"RequestBatchExecutionType: 0 SerialRealtime, 1 SerialFrame (unsupported), 2 Parallel"`
	Requests      []batchRequest `json:"requests" jsonschema:"Sub-requests using official requestType names"`
}

type batchRequest struct {
	RequestType string          `json:"requestType"`
	RequestID   string          `json:"requestId,omitempty"`
	RequestData json.RawMessage `json:"requestData,omitempty"`
}

type batchResult struct {
	Results []batchItemResult `json:"results"`
}

type batchItemResult struct {
	RequestType string `json:"requestType"`
	RequestID   string `json:"requestId,omitempty"`
	OK          bool   `json:"ok"`
	Error       string `json:"error,omitempty"`
	Response    any    `json:"response,omitempty"`
}

func registerBatchTool(s *mcp.Server, h *session.Host) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "RequestBatch",
		Title:       "Request batch",
		Description: "Run a batch of official obs-websocket requests on the connected client. This is a client-side dispatcher (goobs cannot send OpCode 8). SerialFrame is not supported. Sleep is only valid with SerialRealtime.",
		Annotations: annotationsFor("RequestBatch"),
	}, func(ctx context.Context, req *mcp.CallToolRequest, in batchInput) (*mcp.CallToolResult, batchResult, error) {
		c, err := h.RequireClient()
		if err != nil {
			return nil, batchResult{}, err
		}
		out, err := runRequestBatch(ctx, req, c, in)
		return nil, out, err
	})
}

func runRequestBatch(ctx context.Context, req *mcp.CallToolRequest, c *goobs.Client, in batchInput) (batchResult, error) {
	if len(in.Requests) == 0 {
		return batchResult{}, fmt.Errorf("requests must not be empty")
	}
	exec := batchSerialRealtime
	if in.ExecutionType != nil {
		exec = *in.ExecutionType
	}
	if exec == batchSerialFrame {
		return batchResult{}, fmt.Errorf("executionType SerialFrame is not supported: goobs cannot send RequestBatch OpCode 8")
	}
	for _, r := range in.Requests {
		if !protocol.IsKnownRequest(r.RequestType) {
			return batchResult{}, fmt.Errorf("UnknownRequestType: %s", r.RequestType)
		}
		if r.RequestType == "Sleep" && exec != batchSerialRealtime {
			return batchResult{}, fmt.Errorf("Sleep is only valid in RequestBatch with SerialRealtime")
		}
	}
	switch exec {
	case batchParallel:
		return runBatchParallel(ctx, req, c, in)
	case batchSerialRealtime:
		return runBatchSerial(ctx, req, c, in)
	default:
		return batchResult{}, fmt.Errorf("unsupported executionType %d", exec)
	}
}

func runBatchSerial(ctx context.Context, req *mcp.CallToolRequest, c *goobs.Client, in batchInput) (batchResult, error) {
	total := float64(len(in.Requests))
	out := batchResult{Results: make([]batchItemResult, 0, len(in.Requests))}
	for i, r := range in.Requests {
		notifyProgress(ctx, req, float64(i), total, r.RequestType)
		item := runOne(c, r)
		out.Results = append(out.Results, item)
		if in.HaltOnFailure && !item.OK {
			break
		}
	}
	notifyProgress(ctx, req, total, total, "batch complete")
	return out, nil
}

func runBatchParallel(ctx context.Context, req *mcp.CallToolRequest, c *goobs.Client, in batchInput) (batchResult, error) {
	total := float64(len(in.Requests))
	out := batchResult{Results: make([]batchItemResult, len(in.Requests))}
	var mu sync.Mutex
	var halt bool
	var done int
	var wg sync.WaitGroup
	for i, r := range in.Requests {
		if ctx.Err() != nil {
			return out, ctx.Err()
		}
		wg.Add(1)
		go func(i int, r batchRequest) {
			defer wg.Done()
			if in.HaltOnFailure {
				mu.Lock()
				stopped := halt
				mu.Unlock()
				if stopped {
					return
				}
			}
			item := runOne(c, r)
			mu.Lock()
			out.Results[i] = item
			if in.HaltOnFailure && !item.OK {
				halt = true
			}
			done++
			notifyProgress(ctx, req, float64(done), total, r.RequestType)
			mu.Unlock()
		}(i, r)
	}
	wg.Wait()
	return out, nil
}

func runOne(c *goobs.Client, req batchRequest) batchItemResult {
	item := batchItemResult{RequestType: req.RequestType, RequestID: req.RequestID}
	fn, ok := lookupDispatcher(req.RequestType)
	if !ok {
		item.Error = fmt.Sprintf("UnknownRequestType: %s", req.RequestType)
		return item
	}
	resp, err := fn(c, req.RequestData)
	if err != nil {
		item.Error = err.Error()
		return item
	}
	item.OK = true
	item.Response = resp
	return item
}
