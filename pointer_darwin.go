//go:build darwin

package main

/*
#cgo CFLAGS: -x objective-c -I${SRCDIR}
#cgo LDFLAGS: -framework AppKit -framework CoreGraphics

#include "pointer_darwin.h"
*/
import "C"

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver"
)

func pointerInWindow(win fyne.Window) (pointerState, bool) {
	native, ok := win.(driver.NativeWindow)
	if !ok {
		return pointerState{}, false
	}

	var found C.geyesPointer
	ok = false
	native.RunNative(func(context any) {
		mac, isMac := context.(driver.MacWindowContext)
		if !isMac || mac.NSWindow == 0 {
			return
		}
		ok = bool(C.geyesPointerInWindow(C.uintptr_t(mac.NSWindow), &found))
	})
	if !ok {
		return pointerState{}, false
	}

	return pointerState{
		x:       float64(found.x),
		y:       float64(found.y),
		width:   float64(found.width),
		height:  float64(found.height),
		windowX: float64(found.windowX),
		windowY: float64(found.windowY),
	}, true
}
