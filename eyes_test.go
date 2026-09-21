package main

import (
	"image/color"
	"math"
	"os"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestMain(m *testing.M) {
	// Canvas objects ask the running app to repaint them when they are moved, so
	// one has to exist even though nothing here is ever drawn.
	test.NewApp()
	os.Exit(m.Run())
}

// layoutEyes builds a pair of eyes looking at target and hands back the
// renderer, so that the geometry can be inspected without a running app.
func layoutEyes(target fyne.Position) *eyesRenderer {
	e := newEyes(nil)
	r := e.CreateRenderer().(*eyesRenderer)
	r.Layout(fyne.NewSize(windowWidth, windowHeight))

	e.target = target
	r.aim()
	return r
}

func eyeCentre(e *eye) fyne.Position {
	return fyne.NewPos(e.pos.X+e.size.Width/2, e.pos.Y+e.size.Height/2)
}

func pupilCentre(e *eye) fyne.Position {
	radius := e.pupil.Size().Width / 2
	return fyne.NewPos(e.pupil.Position().X+radius, e.pupil.Position().Y+radius)
}

// pupilOffset is how far the pupil has moved from the centre of its eye.
func pupilOffset(e *eye) float64 {
	centre, pupil := eyeCentre(e), pupilCentre(e)
	return math.Hypot(float64(pupil.X-centre.X), float64(pupil.Y-centre.Y))
}

// pupilClearance is the gap between the pupil and the black border. A negative
// gap means the pupil has strayed into the border.
//
// The pupil is round and the white it moves in is an oval, so the reach used to
// place it is not by itself proof that it fits: on the diagonals a point on the
// reach oval sits closer to the border than it does straight up or sideways.
// This measures the real distance instead.
func pupilClearance(e *eye) float64 {
	radiusX := float64(e.size.Width)/2 - borderWidth
	radiusY := float64(e.size.Height)/2 - borderWidth
	centre, pupil := eyeCentre(e), pupilCentre(e)
	offsetX := float64(pupil.X - centre.X)
	offsetY := float64(pupil.Y - centre.Y)

	// The nearest point on an oval has no tidy closed form, so walk around it.
	nearest := math.MaxFloat64
	const steps = 3600
	for step := range steps {
		around := float64(step) * 2 * math.Pi / steps
		nearest = math.Min(nearest, math.Hypot(
			radiusX*math.Cos(around)-offsetX,
			radiusY*math.Sin(around)-offsetY,
		))
	}

	return nearest - float64(e.pupil.Size().Width)/2
}

func TestEyesAreUprightOvals(t *testing.T) {
	r := layoutEyes(fyne.NewPos(0, 0))

	for i, e := range r.pair {
		if e.size.Height <= e.size.Width {
			t.Errorf("eye %d is %v, want taller than it is wide", i, e.size)
		}
	}
}

// TestEyesAreWhiteBehindABlackBorder checks the xeyes look: opaque white,
// ringed by a black border thick enough to read as one.
func TestEyesAreWhiteBehindABlackBorder(t *testing.T) {
	r := layoutEyes(fyne.NewPos(0, 0))

	for i, e := range r.pair {
		if got := e.ball.FillColor; got != color.Color(whiteColour) {
			t.Errorf("eye %d fills with %v, want %v", i, got, whiteColour)
		}
		if got := e.ball.StrokeColor; got != color.Color(blackColour) {
			t.Errorf("eye %d is bordered %v, want %v", i, got, blackColour)
		}

		// A hairline would not look like xeyes, and a border wider than the white
		// it surrounds would swallow the eye.
		border := e.ball.StrokeWidth
		if thinnest := e.size.Width / 20; border < thinnest {
			t.Errorf("eye %d border is %v thick, want at least %v", i, border, thinnest)
		}
		if thickest := e.size.Width / 4; border > thickest {
			t.Errorf("eye %d border is %v thick, want at most %v", i, border, thickest)
		}
	}
}

// TestPupilsAreSmall keeps the pupils dots rather than the wide discs that made
// the eyes look startled.
func TestPupilsAreSmall(t *testing.T) {
	r := layoutEyes(fyne.NewPos(0, 0))

	for i, e := range r.pair {
		pupil, white := e.pupil.Size().Width, e.size.Width-2*borderWidth
		if pupil > white/3 {
			t.Errorf("pupil %d is %v across in %v of white, want no more than a third", i, pupil, white)
		}
	}
}

func TestPupilsStayInsideEyes(t *testing.T) {
	// Targets all around, well outside the window, plus a couple within it. The
	// diagonals matter most: that is where the pupil comes closest to the border.
	targets := []fyne.Position{
		{X: -2000, Y: -2000}, {X: windowWidth / 2, Y: -2000}, {X: 2000, Y: -2000},
		{X: -2000, Y: windowHeight / 2}, {X: 2000, Y: windowHeight / 2},
		{X: -2000, Y: 2000}, {X: windowWidth / 2, Y: 2000}, {X: 2000, Y: 2000},
		{X: windowWidth / 2, Y: windowHeight / 2}, {X: 10, Y: 10},
	}

	for _, target := range targets {
		r := layoutEyes(target)
		for i, e := range r.pair {
			if gap := pupilClearance(e); gap < 0 {
				t.Errorf("looking at %v, pupil %d crosses the border by %.2f", target, i, -gap)
			}
		}
	}
}

// TestPupilsPointAtTheTarget checks each pupil against the centre of its own
// eye, since an eye off to one side is already looking sideways at anything in
// the middle of the window.
func TestPupilsPointAtTheTarget(t *testing.T) {
	middle := float32(windowHeight / 2)
	left := layoutEyes(fyne.NewPos(-2000, middle))
	right := layoutEyes(fyne.NewPos(2000, middle))
	up := layoutEyes(fyne.NewPos(windowWidth/2, -2000))
	down := layoutEyes(fyne.NewPos(windowWidth/2, 2000))

	for i := range left.pair {
		centre := eyeCentre(left.pair[i])

		if got := pupilCentre(left.pair[i]).X; got >= centre.X {
			t.Errorf("pupil %d is at x %v looking left, want left of %v", i, got, centre.X)
		}
		if got := pupilCentre(right.pair[i]).X; got <= centre.X {
			t.Errorf("pupil %d is at x %v looking right, want right of %v", i, got, centre.X)
		}
		if got := pupilCentre(up.pair[i]).Y; got >= centre.Y {
			t.Errorf("pupil %d is at y %v looking up, want above %v", i, got, centre.Y)
		}
		if got := pupilCentre(down.pair[i]).Y; got <= centre.Y {
			t.Errorf("pupil %d is at y %v looking down, want below %v", i, got, centre.Y)
		}
	}
}

// TestPupilsLookAtTheirOwnEye checks that an eye looks straight ahead when the
// target sits right on top of it, rather than jumping to one side.
func TestPupilsLookAtTheirOwnEye(t *testing.T) {
	r := layoutEyes(fyne.NewPos(0, 0))
	first := r.pair[0]
	centre := fyne.NewPos(
		first.pos.X+first.size.Width/2,
		first.pos.Y+first.size.Height/2,
	)

	r = layoutEyes(centre)
	if off := pupilOffset(r.pair[0]); off > 0.5 {
		t.Errorf("pupil sits %.2f from the centre of its eye, want centred", off)
	}
}
