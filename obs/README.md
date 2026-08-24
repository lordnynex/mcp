# obs

OBS Studio MCP server. It talks to OBS over the official websocket protocol through [goobs](https://github.com/andreykaipov/goobs). Tool names are the official obs-websocket request names from [protocol.md](https://github.com/obsproject/obs-websocket/blob/master/docs/generated/protocol.md).

Until `Connect` succeeds, only session tools are listed. After a successful identify, the server sends `notifications/tools/list_changed` and exposes every official request plus event subscription.

The process logs to stderr via [`lib/logger`](../lib/logger). Startup lines describe server identity, MCP tunables, OBS Connect defaults (never the password), and that OBS is not connected until `Connect`.

## Run

```bash
go run ./obs/cmd/obs
go run ./obs/cmd/obs --help
```

Run the package directory, not a single file. `go run obs/cmd/obs/main.go` only compiles `main.go` and fails with `undefined: rootCmd`.

The server speaks MCP on stdout. Logs go to stderr.

## Configuration

Runtime values come from Viper: flags, `OBS_` environment variables, and an optional config file. Host/password/timeout are **defaults for the Connect tool**, not an automatic OBS connection.

| Key | Flag | Env | Default |
| --- | --- | --- | --- |
| `log-level` | `--log-level` | `OBS_LOG_LEVEL` | `info` |
| `host` | `--host` | `OBS_HOST` | `localhost:4455` |
| `password` | `--password` | `OBS_PASSWORD` | empty |
| `response-timeout` | `--response-timeout` | `OBS_RESPONSE_TIMEOUT` | 10s |
| `page-size` | `--page-size` | `OBS_PAGE_SIZE` | 50 |
| `keep-alive` | `--keep-alive` | `OBS_KEEP_ALIVE` | 60s |
| `keep-alive-failure-threshold` | `--keep-alive-failure-threshold` | `OBS_KEEP_ALIVE_FAILURE_THRESHOLD` | 2 |
| `http-addr` | `--http-addr` | `OBS_HTTP_ADDR` | empty (HTTP off) |
| `http-path` | `--http-path` | `OBS_HTTP_PATH` | `/mcp` |
| `http-stateless` | `--http-stateless` | `OBS_HTTP_STATELESS` | false |
| `http-json` | `--http-json` | `OBS_HTTP_JSON` | false (SSE) |
| `http-session-timeout` | `--http-session-timeout` | `OBS_HTTP_SESSION_TIMEOUT` | 30m |
| `http-only` | `--http-only` | `OBS_HTTP_ONLY` | false |
| `http-bearer-token` | `--http-bearer-token` | `OBS_HTTP_BEARER_TOKEN` | empty |

`--config` selects a file. If omitted, the CLI looks for `obs.yaml` in the current directory and `$HOME/.config/mcp`. A missing file is ignored.

```bash
go run ./obs/cmd/obs --host 127.0.0.1:4455 --password secret
OBS_HOST=127.0.0.1:4455 OBS_PASSWORD=secret go run ./obs/cmd/obs
go run ./obs/cmd/obs --config ./obs.yaml
```

## Session tools

| Tool | When | Purpose |
| --- | --- | --- |
| `Connect` | always | Identify with OBS (`host`, `password`, optional `eventSubscriptions` bitmask). |
| `ConnectionStatus` | always | Whether the process is connected. |
| `Disconnect` | after connect | Close the OBS client and remove protocol tools. |
| `RequestBatch` | after connect | Client-side batch of official requests. `SerialFrame` is not supported (goobs cannot send OpCode 8). `Sleep` is only valid with `SerialRealtime`. |
| `SubscribeEvents` | after connect | Select official event names (or categories such as `Scenes`). Returns `obs://events/<EventType>` URIs. |
| `UnsubscribeEvents` | after connect | Drop those event selections. |

## Events

1. Call `SubscribeEvents` with official event type names (or an `EventSubscription` category: `General`, `Config`, `Scenes`, `Inputs`, `Transitions`, `Filters`, `Outputs`, `SceneItems`, `MediaInputs`, `Vendors`, `Ui`, `Canvases`). High-volume events (`InputVolumeMeters`, `InputActiveStateChanged`, `InputShowStateChanged`, `SceneItemTransformChanged`) must be named explicitly and enabled in `Connect.eventSubscriptions`.
2. Call `resources/subscribe` (legacy) or `subscriptions/listen` with the returned URIs.
3. Matching OBS events emit `notifications/resources/updated` on the MCP session (stdio JSON-RPC notifications). The payload is also in `_meta.obsEvent`. `resources/read` on `obs://events/<EventType>` returns the latest event; `obs://events` is a recent ring buffer. Clients can discover the per-type pattern via the `obs://events/{eventType}` resource template.

## Protocol tools

After connect, one tool exists per official request (147), including `GetVersion`, `SetCurrentProgramScene`, `CallVendorRequest` (vendor plugins), and the rest of the protocol request list. Descriptions and fields come from the official protocol document, via goobs generated types.

`Sleep` as a standalone tool returns the protocol restriction; use `RequestBatch` with `executionType` `0` (`SerialRealtime`).

`StartStream` asks for confirmation via MCP elicitation. **Recording tools never elicit** (`StartRecord`, `StopRecord`, `ToggleRecord`, and the rest run immediately) so unattended video capture is not blocked.

## HTTP

Stdio is the default. Set `--http-addr` (for example `127.0.0.1:8080`, not `0.0.0.0`) to serve [streamable HTTP](https://modelcontextprotocol.io/specification/2025-06-18/basic/transports#streamable-http). If stdin is a terminal, stdio is skipped so keepalive `ping` frames are not written to the console. When an MCP host pipes stdio (Cursor), the same process serves both. They share one OBS `Connect`. Use `--http-only` to skip stdio even when stdin is a pipe.

```bash
go run ./obs/cmd/obs --http-addr 127.0.0.1:8080
go run ./obs/cmd/obs --http-addr 127.0.0.1:8080 --http-only
```

- **Stateful (default).** After `initialize`, the server sets `Mcp-Session-Id`. Later POST/GET/DELETE must send that header. GET is the SSE notification stream; DELETE ends the session. Idle sessions close after `--http-session-timeout` (default 30m; `0` disables). A memory event store supports Last-Event-ID resume on protocol `2025-11-25` and earlier.
- **`--http-stateless`.** Required for HTTP clients on protocol `2026-07-28`. Session IDs are not used. Server-to-client requests cannot complete, so `StartStream` elicitation fails on this path. Stdio is unaffected.
- **`--http-json`.** Return `application/json` instead of `text/event-stream`.
- **`--http-bearer-token`.** Optional static bearer. Requests without `Authorization: Bearer …` get 401. The token is never logged (`tokenSet` only). This is not full OAuth / protected-resource metadata.

## MCP features

This server is written as a current-protocol reference (MCP 2026-07-28 and the Go SDK v1.7 options that are not deprecated).

- **Instructions** — connect-first workflow, event subscribe, and prompt names.
- **Implementation metadata** — `title`, `description`, and `websiteUrl` in addition to name/version.
- **Prompts** — `obs-connect`, `obs-switch-scene`, `obs-studio-status`, `obs-subscribe-events`, `obs-start-stream`, `obs-start-record`, `obs-request-batch`. After Connect, `obs-current-program` is added (`notifications/prompts/list_changed`).
- **Completions** — `completion/complete` for `obs-switch-scene` scene names (when connected), `obs-subscribe-events` categories, and `obs://events/{eventType}` event names.
- **Progress** — `notifications/progress` on `Connect`, `RequestBatch`, and `GetSourceScreenshot` when the client sends a progress token.
- **Elicitation** — `StartStream` only. Clients that cannot elicit can still record.
- **Pagination** — `tools/list` uses `--page-size` (default 50) so the full protocol tool list is cursor-paginated.
- **Keepalive** — server-originated `ping` on `--keep-alive` (default 60s), closing after `--keep-alive-failure-threshold` consecutive failures (default 2).
- **List cache hints** — `ttlMs` on tools/resources/prompts list results (SEP-2549).
- **Resource template** — `obs://events/{eventType}`.
- **Streamable HTTP** — opt-in `--http-addr`; stateful `Mcp-Session-Id` by default; `--http-stateless` for protocol 2026-07-28 HTTP clients.

Not implemented (deprecated as of protocol 2026-07-28 / SEP-2577): MCP logging to the client, roots, and sampling. Process logs stay on stderr via `lib/logger`.

Also not in this server: TLS, full OAuth / protected-resource metadata, icons, `x-mcp-header` tool annotations, custom JSON-RPC methods, per-MCP-session OBS clients, or a durable EventStore. `server/discover` is handled by the SDK.
