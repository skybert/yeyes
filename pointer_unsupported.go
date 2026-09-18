//go:build !darwin && (wayland || (!linux && !freebsd && !netbsd && !openbsd))

package main

import "fyne.io/fyne/v2"

// pointerInWindow has nobody to ask on this windowing system, so the eyes just
// stare straight ahead.
func pointerInWindow(fyne.Window) (pointerState, bool) {
	return pointerState{}, false
}
