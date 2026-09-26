// Package live wraps go-vgo/robotgo behind desktop.Driver.
//
// The implementation is compiled only with a backend tag:
//
//	-tags purego              // macOS (or wayland on Linux)
//	-tags "purego,x11"        // Linux X11
//	-tags "purego,wayland"    // Linux wlroots
//	-tags "purego,libei"      // Linux GNOME/KDE (input only)
package live
