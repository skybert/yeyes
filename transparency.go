package main

import "github.com/go-gl/glfw/v3.4/glfw"

// requestTransparentFramebuffer asks for the next window to be created with an
// alpha channel that the desktop compositor blends, instead of the usual
// opaque one. A transparent theme colour alone is not enough: without this the
// compositor paints our zeroed pixels black.
//
// Fyne has no API for this, but its desktop driver is built on GLFW, whose
// window hints are global state applying to the next window created. That
// window is ours, as long as this is called after the driver has started (GLFW
// is initialised by then) and before the window is shown.
func requestTransparentFramebuffer() {
	glfw.WindowHint(glfw.TransparentFramebuffer, glfw.True)
}
