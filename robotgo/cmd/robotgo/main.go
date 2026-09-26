package main

import (
	"context"
	"log/slog"
	"os"
)

func main() {
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		slog.Error("command failed", "err", err)
		os.Exit(1)
	}
}
