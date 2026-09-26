//go:build !purego && !mac && !x11 && !wayland && !libei && !win && !robotgocgo

package main

import "github.com/lordnynex/mcp/robotgo/desktop"

func init() {
	desktop.SetDefault(desktop.Missing())
}
