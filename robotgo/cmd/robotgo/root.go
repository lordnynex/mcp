package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lordnynex/mcp/lib/logger"
	rgomcp "github.com/lordnynex/mcp/robotgo"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "robotgo",
	Short: "robotgo desktop MCP server",
	Long: `MCP server that drives the local desktop through go-vgo/robotgo.

The live backend is selected at build time. Linux must name x11, wayland, or libei:

  go run -tags purego ./robotgo/cmd/robotgo
  go run -tags "purego,x11" ./robotgo/cmd/robotgo
  go run -tags "purego,wayland" ./robotgo/cmd/robotgo
  go run -tags "purego,libei" ./robotgo/cmd/robotgo

Bare -tags purego on Linux selects wayland (GNOME/KDE need libei; libei cannot screenshot).`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := initConfig(); err != nil {
			return err
		}
		return initLogger()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		rgomcp.SetMCPDefaults(rgomcp.MCPDefaults{
			PageSize:                  viper.GetInt("page-size"),
			KeepAlive:                 viper.GetDuration("keep-alive"),
			KeepAliveFailureThreshold: viper.GetInt("keep-alive-failure-threshold"),
		})
		rgomcp.SetDesktopDefaults(rgomcp.DesktopDefaults{
			KeySleep:           viper.GetInt("key-sleep"),
			MouseSleep:         viper.GetInt("mouse-sleep"),
			ScreenshotMaxWidth: viper.GetInt("screenshot-max-width"),
			ScreenshotFormat:   viper.GetString("screenshot-format"),
		})

		mcpOpts := rgomcp.MCPDefaults{
			PageSize:                  viper.GetInt("page-size"),
			KeepAlive:                 viper.GetDuration("keep-alive"),
			KeepAliveFailureThreshold: viper.GetInt("keep-alive-failure-threshold"),
		}
		if mcpOpts.PageSize <= 0 {
			mcpOpts.PageSize = 50
		}
		if mcpOpts.KeepAlive <= 0 {
			mcpOpts.KeepAlive = 60 * time.Second
		}
		if mcpOpts.KeepAliveFailureThreshold <= 0 {
			mcpOpts.KeepAliveFailureThreshold = 2
		}

		slog.Info("starting robotgo MCP server",
			"name", rgomcp.Name,
			"version", rgomcp.Version,
			"title", rgomcp.Title,
		)
		slog.Info("MCP options",
			"pageSize", mcpOpts.PageSize,
			"keepAlive", mcpOpts.KeepAlive,
			"keepAliveFailureThreshold", mcpOpts.KeepAliveFailureThreshold,
		)
		slog.Info("desktop options",
			"keySleep", viper.GetInt("key-sleep"),
			"mouseSleep", viper.GetInt("mouse-sleep"),
			"screenshotMaxWidth", viper.GetInt("screenshot-max-width"),
			"screenshotFormat", viper.GetString("screenshot-format"),
		)

		httpAddr := strings.TrimSpace(viper.GetString("http-addr"))
		httpOnly := viper.GetBool("http-only")
		if httpAddr != "" && !httpOnly && stdinIsTerminal() {
			httpOnly = true
			slog.Info("stdin is a terminal; skipping stdio so keepalive pings are not written to the console (use a piped MCP client to serve both)")
		}
		if httpOnly && httpAddr == "" {
			return fmt.Errorf("http-only requires http-addr")
		}

		s := rgomcp.New()
		if httpAddr != "" {
			path := httpPath(viper.GetString("http-path"))
			stateless := viper.GetBool("http-stateless")
			jsonResp := viper.GetBool("http-json")
			sessionTimeout := viper.GetDuration("http-session-timeout")
			bearer := viper.GetString("http-bearer-token")
			mode := "sse"
			if jsonResp {
				mode = "json"
			}
			slog.Info("HTTP transport listening",
				"addr", httpAddr,
				"path", path,
				"stateless", stateless,
				"response", mode,
				"sessionTimeout", sessionTimeout,
				"tokenSet", bearer != "",
			)
			if stateless {
				slog.Info("HTTP is stateless; protocol 2026-07-28 is accepted; Mcp-Session-Id is not used; KillProcess elicitation cannot complete over this transport")
			} else {
				slog.Info("HTTP is stateful; Mcp-Session-Id is required after initialize; protocol 2026-07-28 HTTP clients must use --http-stateless")
			}
			mux := http.NewServeMux()
			mux.Handle(path, rgomcp.NewHTTPHandler(s, rgomcp.HTTPOptions{
				Stateless:      stateless,
				JSONResponse:   jsonResp,
				SessionTimeout: sessionTimeout,
				BearerToken:    bearer,
			}))
			if httpOnly {
				slog.Info("http-only; stdio transport is not started")
				return rgomcp.ListenAndServeHTTP(cmd.Context(), httpAddr, mux)
			}
			go func() {
				if err := rgomcp.ListenAndServeHTTP(cmd.Context(), httpAddr, mux); err != nil {
					slog.Error("http server stopped", "err", err)
				}
			}()
		}

		if httpAddr == "" {
			slog.Info("stdio transport ready; waiting for MCP client")
		} else {
			slog.Info("stdio and HTTP transports ready; waiting for MCP clients")
		}
		return s.Run(cmd.Context(), &mcp.StdioTransport{})
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./robotgo.yaml or $HOME/.config/mcp/robotgo.yaml)")
	rootCmd.PersistentFlags().String("log-level", "info", "log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().Int("page-size", 50, "MCP list pagination page size")
	rootCmd.PersistentFlags().Duration("keep-alive", 60*time.Second, "MCP server ping interval")
	rootCmd.PersistentFlags().Int("keep-alive-failure-threshold", 2, "consecutive keepalive ping failures before closing the session")
	rootCmd.PersistentFlags().String("http-addr", "", "listen address for streamable HTTP (empty disables HTTP)")
	rootCmd.PersistentFlags().String("http-path", "/mcp", "URL path for streamable HTTP")
	rootCmd.PersistentFlags().Bool("http-stateless", false, "stateless HTTP (required for protocol 2026-07-28; disables session IDs)")
	rootCmd.PersistentFlags().Bool("http-json", false, "return application/json instead of SSE on HTTP")
	rootCmd.PersistentFlags().Duration("http-session-timeout", rgomcp.DefaultHTTPSessionTimeout(), "idle HTTP session timeout (0 disables)")
	rootCmd.PersistentFlags().Bool("http-only", false, "serve HTTP only (do not start stdio)")
	rootCmd.PersistentFlags().String("http-bearer-token", "", "optional bearer token required on HTTP requests")
	rootCmd.PersistentFlags().Int("key-sleep", 10, "robotgo KeySleep milliseconds")
	rootCmd.PersistentFlags().Int("mouse-sleep", 0, "robotgo MouseSleep milliseconds")
	rootCmd.PersistentFlags().Int("screenshot-max-width", 1280, "max encoded screenshot width")
	rootCmd.PersistentFlags().String("screenshot-format", "png", "screenshot format (png or jpeg)")
	_ = viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	_ = viper.BindPFlag("page-size", rootCmd.PersistentFlags().Lookup("page-size"))
	_ = viper.BindPFlag("keep-alive", rootCmd.PersistentFlags().Lookup("keep-alive"))
	_ = viper.BindPFlag("keep-alive-failure-threshold", rootCmd.PersistentFlags().Lookup("keep-alive-failure-threshold"))
	_ = viper.BindPFlag("http-addr", rootCmd.PersistentFlags().Lookup("http-addr"))
	_ = viper.BindPFlag("http-path", rootCmd.PersistentFlags().Lookup("http-path"))
	_ = viper.BindPFlag("http-stateless", rootCmd.PersistentFlags().Lookup("http-stateless"))
	_ = viper.BindPFlag("http-json", rootCmd.PersistentFlags().Lookup("http-json"))
	_ = viper.BindPFlag("http-session-timeout", rootCmd.PersistentFlags().Lookup("http-session-timeout"))
	_ = viper.BindPFlag("http-only", rootCmd.PersistentFlags().Lookup("http-only"))
	_ = viper.BindPFlag("http-bearer-token", rootCmd.PersistentFlags().Lookup("http-bearer-token"))
	_ = viper.BindPFlag("key-sleep", rootCmd.PersistentFlags().Lookup("key-sleep"))
	_ = viper.BindPFlag("mouse-sleep", rootCmd.PersistentFlags().Lookup("mouse-sleep"))
	_ = viper.BindPFlag("screenshot-max-width", rootCmd.PersistentFlags().Lookup("screenshot-max-width"))
	_ = viper.BindPFlag("screenshot-format", rootCmd.PersistentFlags().Lookup("screenshot-format"))
	viper.SetDefault("log-level", "info")
	viper.SetDefault("page-size", 50)
	viper.SetDefault("keep-alive", 60*time.Second)
	viper.SetDefault("keep-alive-failure-threshold", 2)
	viper.SetDefault("http-path", "/mcp")
	viper.SetDefault("http-session-timeout", rgomcp.DefaultHTTPSessionTimeout())
	viper.SetDefault("key-sleep", 10)
	viper.SetDefault("mouse-sleep", 0)
	viper.SetDefault("screenshot-max-width", 1280)
	viper.SetDefault("screenshot-format", "png")
}

func stdinIsTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

func httpPath(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return "/mcp"
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return p
}

func initConfig() error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("robotgo")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.config/mcp")
	}

	viper.SetEnvPrefix("ROBOTGO")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if missingConfig(err) {
			return nil
		}
		return fmt.Errorf("read config: %w", err)
	}
	return nil
}

func initLogger() error {
	level, err := parseLevel(viper.GetString("log-level"))
	if err != nil {
		return err
	}
	logger.SetDefault(logger.Options{Level: level})
	return nil
}

func parseLevel(s string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug, nil
	case "info", "":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("invalid log-level %q", s)
	}
}

func missingConfig(err error) bool {
	var notFound viper.ConfigFileNotFoundError
	if errors.As(err, &notFound) {
		return true
	}
	return errors.Is(err, fs.ErrNotExist) || os.IsNotExist(err)
}
