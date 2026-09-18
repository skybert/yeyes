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
	// eyeStroke is the thickness of the outline around each eyeball.
	eyeStroke = 3

	// eyeGap is how much of the total width is left empty between the eyes.
	eyeGap = 0.06

	// pupilSize is the pupil radius as a fraction of the eyeball's narrower
	// radius, which is the horizontal one on our upright ovals.
	pupilSize = 0.36
)

var (
	eyeballColour = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	outlineColour = color.NRGBA{A: 0xff}
	pupilColour   = color.NRGBA{A: 0xff}
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
	for i := range r.balls {
		ball := canvas.NewEllipse(eyeballColour)
		ball.StrokeColor = outlineColour
		ball.StrokeWidth = eyeStroke

		r.balls[i] = ball
		r.pupils[i] = canvas.NewCircle(pupilColour)
	}

	// The eyeballs come first so that the pupils are drawn on top of them.
	r.objects = []fyne.CanvasObject{r.balls[0], r.balls[1], r.pupils[0], r.pupils[1]}
	return r
}

type eyesRenderer struct {
	eyes *eyes

	balls   [2]*canvas.Ellipse
	pupils  [2]*canvas.Circle
	objects []fyne.CanvasObject
}

func (r *eyesRenderer) Layout(size fyne.Size) {
	gap := size.Width * eyeGap
	ball := fyne.NewSize((size.Width-gap)/2-eyeStroke, size.Height-eyeStroke)

	for i, eye := range r.balls {
		x := eyeStroke/2 + float32(i)*(ball.Width+eyeStroke+gap)
		eye.Move(fyne.NewPos(x, eyeStroke/2))
		eye.Resize(ball)
	}

	r.aimPupils()
}

func (r *eyesRenderer) MinSize() fyne.Size {
	return fyne.NewSize(60, 45)
}

func (r *eyesRenderer) Refresh() {
	r.aimPupils()
}

func (r *eyesRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *eyesRenderer) Destroy() {}

// aimPupils moves each pupil as close to the target as it can get without
// leaving its eyeball, the way xeyes does it.
func (r *eyesRenderer) aimPupils() {
	for i, ball := range r.balls {
		size, pos := ball.Size(), ball.Position()
		radiusX, radiusY := size.Width/2, size.Height/2
		centreX, centreY := pos.X+radiusX, pos.Y+radiusY
		pupilRadius := fyne.Min(radiusX, radiusY) * pupilSize

		// How far the pupil's centre may travel from the eyeball's centre,
		// keeping the pupil clear of the outline. Following the shape of the
		// eyeball means the pupil can go further up and down than sideways.
		reachX := radiusX - pupilRadius - eyeStroke
		reachY := radiusY - pupilRadius - eyeStroke

		toTargetX := r.eyes.target.X - centreX
		toTargetY := r.eyes.target.Y - centreY

		var offsetX, offsetY float32
		if away := float32(math.Hypot(float64(toTargetX), float64(toTargetY))); away > 0.5 {
			// Look all the way out once the target has left the eyeball, and
			// proportionally less when it is inside.
			out := fyne.Min(1, away/radiusX)
			offsetX = toTargetX / away * reachX * out
			offsetY = toTargetY / away * reachY * out
		}

		pupil := r.pupils[i]
		pupil.Resize(fyne.NewSize(pupilRadius*2, pupilRadius*2))
		pupil.Move(fyne.NewPos(centreX+offsetX-pupilRadius, centreY+offsetY-pupilRadius))
	}
}
