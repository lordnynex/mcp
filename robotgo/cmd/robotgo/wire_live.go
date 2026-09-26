//go:build purego || mac || x11 || wayland || libei || win || robotgocgo

package main

import (
	"github.com/lordnynex/mcp/robotgo/desktop"
	"github.com/lordnynex/mcp/robotgo/desktop/live"
)

func init() {
	desktop.SetDefault(live.New())
}
