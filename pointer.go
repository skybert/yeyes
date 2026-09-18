package main

import "fyne.io/fyne/v2"

// pointerState says where the mouse pointer is in relation to our window, in
// whatever units the windowing system uses: points on macOS, pixels on X11.
//
// It is filled in by pointerInWindow, which has one implementation per
// windowing system in pointer_*.go.
type pointerState struct {
	// x and y locate the pointer relative to the top left corner of the
	// window's content. They fall outside the window whenever the pointer is
	// elsewhere on the desktop, which is nearly always, and that is the point:
	// the eyes need to know which way to look, not whether they are hovered.
	x, y float64

	// width and height are the size of the window's content.
	width, height float64

	// windowX and windowY locate the top left corner of the window's content on
	// the desktop, in the coordinates desktop.Window.RequestPosition expects.
	windowX, windowY float64
}

// canvasPosition converts the pointer position into the coordinate space of
// win's canvas.
func (s pointerState) canvasPosition(win fyne.Window) fyne.Position {
	size := win.Canvas().Size()

	pos := fyne.NewPos(float32(s.x), float32(s.y))
	if s.width > 0 {
		pos.X = float32(s.x / s.width * float64(size.Width))
	}
	if s.height > 0 {
		pos.Y = float32(s.y / s.height * float64(size.Height))
	}
	return pos
}
