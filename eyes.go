package main

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

const (
	// eyeGap is how much of the total width is left empty between the eyes.
	eyeGap = 0.06

	// borderWidth is the thickness of the black border drawn around the white of
	// an eye. Fyne draws an ellipse's stroke inside its bounds, so this eats into
	// the white rather than adding to the size of the eye.
	borderWidth = 8

	// pupilSize is the pupil radius as a fraction of the eye's narrower radius,
	// which is the horizontal one on our upright ovals. Small, as xeyes has it.
	pupilSize = 0.22

	// pupilMargin keeps the pupil from quite touching the border, so that the
	// two never look like they have merged.
	pupilMargin = 1

	// eyeMargin keeps the eyes just clear of the window edge. An ellipse that
	// touches it loses the smoothing on its outermost pixels and comes out
	// looking as though it has been cut off.
	eyeMargin = 2
)

var (
	blackColour = color.NRGBA{A: 0xff}
	whiteColour = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
)

// Declare conformity with the interfaces that let the eyes be dragged around.
var (
	_ desktop.Mouseable = (*eyes)(nil)
	_ fyne.Draggable    = (*eyes)(nil)
)

// eyes is a pair of upright oval eyes whose pupils follow a target position.
type eyes struct {
	widget.BaseWidget

	// target is where the eyes are looking, in the widget's own coordinates.
	target fyne.Position

	// drag carries the window around, as the eyes are the only thing there is
	// to take hold of in a window without a border.
	drag *mover
}

func newEyes(win fyne.Window) *eyes {
	e := &eyes{drag: &mover{win: win}}
	e.ExtendBaseWidget(e)
	return e
}

// LookAt points both pupils at pos, given in the widget's own coordinate
// space. Positions outside the widget are expected and normal.
func (e *eyes) LookAt(pos fyne.Position) {
	if pos == e.target {
		return
	}

	e.target = pos
	e.Refresh()
}

func (e *eyes) MouseDown(*desktop.MouseEvent) { e.drag.take() }

func (e *eyes) MouseUp(*desktop.MouseEvent) {}

func (e *eyes) Dragged(*fyne.DragEvent) { e.drag.move() }

func (e *eyes) DragEnd() { e.drag.release() }

func (e *eyes) CreateRenderer() fyne.WidgetRenderer {
	r := &eyesRenderer{eyes: e}
	for i := range r.pair {
		r.pair[i] = newEye()
		r.objects = append(r.objects, r.pair[i].objects()...)
	}
	return r
}

// eye is a single eye: the white of it inside a black border, and the pupil.
type eye struct {
	ball  *canvas.Ellipse
	pupil *canvas.Circle

	// pos and size are the area the whole eye occupies.
	pos  fyne.Position
	size fyne.Size
}

func newEye() *eye {
	ball := canvas.NewEllipse(whiteColour)
	ball.StrokeColor = blackColour
	ball.StrokeWidth = borderWidth

	return &eye{ball: ball, pupil: canvas.NewCircle(blackColour)}
}

// objects lists the eye's parts in the order they are drawn, pupil last so that
// it sits on top of the white.
func (e *eye) objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{e.ball, e.pupil}
}

// resize puts the eye in the given area.
func (e *eye) resize(pos fyne.Position, size fyne.Size) {
	e.pos, e.size = pos, size

	e.ball.Move(pos)
	e.ball.Resize(size)
}

// lookAt moves the pupil as close to target as it can get without straying into
// the border, the way xeyes does it.
func (e *eye) lookAt(target fyne.Position) {
	radiusX, radiusY := e.size.Width/2, e.size.Height/2
	centreX, centreY := e.pos.X+radiusX, e.pos.Y+radiusY
	pupilRadius := fyne.Min(radiusX, radiusY) * pupilSize

	// How far the pupil's centre may travel from the eye's centre. Following
	// the shape of the eye lets the pupil go further up and down than sideways.
	reachX := radiusX - borderWidth - pupilRadius - pupilMargin
	reachY := radiusY - borderWidth - pupilRadius - pupilMargin

	towardsX := target.X - centreX
	towardsY := target.Y - centreY

	var offsetX, offsetY float32
	if away := float32(math.Hypot(float64(towardsX), float64(towardsY))); away > 0.5 {
		// Look all the way out once the target has left the eye, and
		// proportionally less when it is inside.
		out := fyne.Min(1, away/radiusX)
		offsetX = towardsX / away * reachX * out
		offsetY = towardsY / away * reachY * out
	}

	e.pupil.Resize(fyne.NewSize(pupilRadius*2, pupilRadius*2))
	e.pupil.Move(fyne.NewPos(centreX+offsetX-pupilRadius, centreY+offsetY-pupilRadius))
}

type eyesRenderer struct {
	eyes *eyes

	pair    [2]*eye
	objects []fyne.CanvasObject
}

func (r *eyesRenderer) Layout(size fyne.Size) {
	area := fyne.NewSize(size.Width-2*eyeMargin, size.Height-2*eyeMargin)
	gap := area.Width * eyeGap
	eyeSize := fyne.NewSize((area.Width-gap)/2, area.Height)

	for i, e := range r.pair {
		e.resize(fyne.NewPos(eyeMargin+float32(i)*(eyeSize.Width+gap), eyeMargin), eyeSize)
	}

	r.aim()
}

func (r *eyesRenderer) MinSize() fyne.Size {
	return fyne.NewSize(60, 45)
}

func (r *eyesRenderer) Refresh() {
	r.aim()
}

func (r *eyesRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *eyesRenderer) Destroy() {}

func (r *eyesRenderer) aim() {
	for _, e := range r.pair {
		e.lookAt(r.eyes.target)
	}
}
