package main

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// mover drags a window around the desktop.
//
// A window without a border has no title bar to grab, and the gestures a
// desktop offers for moving such a window are hit and miss: Ctrl+Cmd+drag, for
// one, leaves borderless macOS windows where they are. So we do the moving
// ourselves, from wherever the drag lands on the window.
type mover struct {
	win fyne.Window

	// held is where in the window the drag was started, in native units. The
	// window is moved so that this point stays under the pointer.
	held    pointerState
	holding bool
}

// take remembers where the window was taken hold of.
func (m *mover) take() {
	m.held, m.holding = pointerInWindow(m.win)
}

// move brings the window along with the pointer.
func (m *mover) move() {
	if !m.holding {
		return
	}

	desk, ok := m.win.(desktop.Window)
	if !ok {
		return
	}

	now, ok := pointerInWindow(m.win)
	if !ok {
		return
	}

	// How far the pointer has drifted from the point it took hold of is how far
	// the window has to follow it.
	desk.RequestPosition(
		int(math.Round(now.windowX+now.x-m.held.x)),
		int(math.Round(now.windowY+now.y-m.held.y)),
	)
}

func (m *mover) release() {
	m.holding = false
}
