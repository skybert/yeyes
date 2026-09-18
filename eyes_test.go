package main

import (
	"math"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// layoutEyes builds a pair of eyes looking at target and hands back the
// renderer, so that the geometry can be inspected without a running app.
func layoutEyes(target fyne.Position) *eyesRenderer {
	e := newEyes(nil)
	r := e.CreateRenderer().(*eyesRenderer)
	r.Layout(fyne.NewSize(windowWidth, windowHeight))

	e.target = target
	r.aimPupils()
	return r
}

// insideEyeball reports how far out of its eyeball a pupil reaches, as a
// fraction of the room it has: at most 1 means it is still fully covered.
func insideEyeball(ball *canvas.Ellipse, pupil *canvas.Circle) float64 {
	ballSize, ballPos := ball.Size(), ball.Position()
	radiusX, radiusY := float64(ballSize.Width)/2, float64(ballSize.Height)/2
	centreX := float64(ballPos.X) + radiusX
	centreY := float64(ballPos.Y) + radiusY

	pupilRadius := float64(pupil.Size().Width) / 2
	pupilX := float64(pupil.Position().X) + pupilRadius
	pupilY := float64(pupil.Position().Y) + pupilRadius

	outX := (pupilX - centreX) / (radiusX - pupilRadius)
	outY := (pupilY - centreY) / (radiusY - pupilRadius)
	return math.Hypot(outX, outY)
}

func TestEyeballsAreUprightOvals(t *testing.T) {
	r := layoutEyes(fyne.NewPos(0, 0))

	for i, ball := range r.balls {
		size := ball.Size()
		if size.Height <= size.Width {
			t.Errorf("eyeball %d is %v, want taller than it is wide", i, size)
		}
	}
}

func TestPupilsStayInsideEyeballs(t *testing.T) {
	// Targets all around, well outside the window, plus a couple within it.
	targets := []fyne.Position{
		{X: -2000, Y: -2000}, {X: windowWidth / 2, Y: -2000}, {X: 2000, Y: -2000},
		{X: -2000, Y: windowHeight / 2}, {X: 2000, Y: windowHeight / 2},
		{X: -2000, Y: 2000}, {X: windowWidth / 2, Y: 2000}, {X: 2000, Y: 2000},
		{X: windowWidth / 2, Y: windowHeight / 2}, {X: 10, Y: 10},
	}

	for _, target := range targets {
		r := layoutEyes(target)
		for i, ball := range r.balls {
			if out := insideEyeball(ball, r.pupils[i]); out > 1 {
				t.Errorf("looking at %v, pupil %d reaches %.2f times its room", target, i, out)
			}
		}
	}
}

func TestPupilsPointAtTheTarget(t *testing.T) {
	middle := layoutEyes(fyne.NewPos(windowWidth/2, windowHeight/2))
	left := layoutEyes(fyne.NewPos(-2000, windowHeight/2))
	up := layoutEyes(fyne.NewPos(windowWidth/2, -2000))

	for i := range middle.pupils {
		if got, want := left.pupils[i].Position().X, middle.pupils[i].Position().X; got >= want {
			t.Errorf("pupil %d is at x %v looking left, want less than %v", i, got, want)
		}
		if got, want := up.pupils[i].Position().Y, middle.pupils[i].Position().Y; got >= want {
			t.Errorf("pupil %d is at y %v looking up, want less than %v", i, got, want)
		}
	}
}

// TestPupilsLookAtTheirOwnEye checks that an eye looks straight ahead when the
// target sits right on top of it, rather than jumping to one side.
func TestPupilsLookAtTheirOwnEye(t *testing.T) {
	r := layoutEyes(fyne.NewPos(0, 0))
	ball := r.balls[0]
	centre := fyne.NewPos(
		ball.Position().X+ball.Size().Width/2,
		ball.Position().Y+ball.Size().Height/2,
	)

	r = layoutEyes(centre)
	if out := insideEyeball(r.balls[0], r.pupils[0]); out > 0.01 {
		t.Errorf("pupil sits %.2f out of centre, want centred", out)
	}
}
