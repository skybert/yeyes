// Copyright (C) 2026 Torstein Krause Johansen
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version. See LICENSE for the full text.

// yeyes shows a pair of eyes that follow the mouse pointer around the screen,
// in the spirit of xeyes from X11.
package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
)

const appID = "net.skybert.yeyes"

// The size of the window, in Fyne's device independent units. Taller than it
// is wide per eye, so that the eyes come out as upright ovals.
const (
	windowWidth  = 200
	windowHeight = 150
)

func main() {
	a := app.NewWithID(appID)
	a.Settings().SetTheme(&transparentTheme{Theme: theme.DefaultTheme()})

	win := newEyesWindow(a)
	win.SetTitle("yeyes")
	win.SetPadded(false)
	win.Resize(fyne.NewSize(windowWidth, windowHeight))
	win.SetFixedSize(true)

	eyes := newEyes(win)
	win.SetContent(eyes)

	addQuitShortcuts(a, win)
	if desk, ok := win.(desktop.Window); ok {
		// We are a small overlay, so being buried under other windows would
		// leave nothing to look at.
		desk.RequestAlwaysOnTop()
	}

	// The window is shown from here rather than before Run, because the
	// transparency request only applies to windows created after it, and the
	// driver has to be up before we can make it.
	a.Lifecycle().SetOnStarted(func() {
		requestTransparentFramebuffer()
		win.Show()
		followPointer(eyes, win)
	})

	a.Run()
}

// newEyesWindow creates the window the eyes live in: borderless, so that only
// the eyes themselves are visible.
func newEyesWindow(a fyne.App) fyne.Window {
	if desk, ok := a.Driver().(desktop.Driver); ok {
		return desk.CreateSplashWindow()
	}
	return a.NewWindow("yeyes")
}

// addQuitShortcuts wires up Ctrl+Q and, for macOS, Cmd+Q. Without a window
// border there is no close button, so these are the only way out.
func addQuitShortcuts(a fyne.App, win fyne.Window) {
	quit := func(fyne.Shortcut) { a.Quit() }

	for _, mod := range []fyne.KeyModifier{fyne.KeyModifierControl, fyne.KeyModifierSuper} {
		win.Canvas().AddShortcut(&desktop.CustomShortcut{
			KeyName:  fyne.KeyQ,
			Modifier: mod,
		}, quit)
	}
}

// followPointer keeps the eyes aimed at the mouse pointer. The pointer spends
// nearly all of its time outside our little window, where Fyne receives no
// mouse events at all, so we have to go and look for it ourselves.
//
// The looking is done from an animation, because Fyne ticks those from its own
// run loop: that puts us on the main thread, where the windowing system wants
// to be asked, and stops us as soon as the app does.
func followPointer(eyes *eyes, win fyne.Window) {
	look := &fyne.Animation{
		Duration:    time.Second,
		RepeatCount: fyne.AnimationRepeatForever,
		Curve:       fyne.AnimationLinear,
		Tick: func(float32) {
			if state, ok := pointerInWindow(win); ok {
				eyes.LookAt(state.canvasPosition(win))
			}
		},
	}
	look.Start()
}
