package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"strings"

	"github.com/lordnynex/mcp/lib/logger"
	"github.com/lordnynex/mcp/obs"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "obs",
	Short: "OBS Studio MCP server",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := initConfig(); err != nil {
			return err
		}
		return initLogger()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return obs.New().Run(cmd.Context(), &mcp.StdioTransport{})
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./obs.yaml or $HOME/.config/mcp/obs.yaml)")
	rootCmd.PersistentFlags().String("log-level", "info", "log level (debug, info, warn, error)")
	_ = viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.SetDefault("log-level", "info")
}

func initConfig() error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("obs")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.config/mcp")
	}

	viper.SetEnvPrefix("OBS")
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
