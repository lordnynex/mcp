# mcp

A Go 1.27 workspace of Model Context Protocol (MCP) servers. The root module is `github.com/lordnynex/mcp`. Each top-level directory except `lib/` is one server package that can run on its own or be mounted into a future aggregate server.

```bash
go test ./...
go run ./obs/cmd/obs --help
go run -tags purego ./robotgo/cmd/robotgo --help
```

Shared libraries live under `lib/`. Process logs use `log/slog` via [`lib/logger`](lib/logger). CLIs are built with Cobra and configured with Viper.

## Servers

| Server | Description | Docs |
| --- | --- | --- |
| obs | OBS Studio MCP server | [obs/README.md](obs/README.md) |
| robotgo | Desktop automation (keyboard, mouse, screenshots) | [robotgo/README.md](robotgo/README.md) |
