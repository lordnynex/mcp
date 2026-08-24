# obs

OBS Studio MCP server. Tools are not implemented yet; this module currently exposes a mountable `*mcp.Server` and a Cobra CLI that serves it over stdio.

## Run

```bash
go run ./obs/cmd/obs
go run ./obs/cmd/obs --help
```

Run the package directory, not a single file. `go run obs/cmd/obs/main.go` only compiles `main.go` and fails with `undefined: rootCmd`.

The server speaks MCP on stdout. Logs go to stderr.

## Configuration

Runtime values come from Viper: flags, `OBS_` environment variables, and an optional config file.

| Key | Flag | Env | Default |
| --- | --- | --- | --- |
| `log-level` | `--log-level` | `OBS_LOG_LEVEL` | `info` |

`--config` selects a file. If omitted, the CLI looks for `obs.yaml` in the current directory and `$HOME/.config/mcp`. A missing file is ignored.

```bash
go run ./obs/cmd/obs --log-level debug
OBS_LOG_LEVEL=warn go run ./obs/cmd/obs
go run ./obs/cmd/obs --config ./obs.yaml
```
