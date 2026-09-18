//go:build (linux || freebsd || netbsd || openbsd) && !wayland

package main

/*
#cgo LDFLAGS: -lX11

#include <stdbool.h>
#include <stdint.h>
#include <X11/Xlib.h>

typedef struct {
  double x, y;
  double width, height;
  double windowX, windowY;
} geyesPointer;

// geyesPointerInWindow fills in out for the given X11 window, in pixels.
// XQueryPointer answers even when the pointer is over some other window, which
// is what lets the eyes follow it across the whole screen.
//
// Only ever called from the main thread, so the cached display needs no lock.
static bool geyesPointerInWindow(uintptr_t xWindow, geyesPointer *out) {
  static Display *display = NULL;
  if (display == NULL) {
    display = XOpenDisplay(NULL);
    if (display == NULL) {
      return false;
    }
  }

  Window window = (Window)xWindow;
  Window root, child;
  int rootX, rootY, windowX, windowY;
  unsigned int buttons;
  if (!XQueryPointer(display, window, &root, &child, &rootX, &rootY, &windowX,
                     &windowY, &buttons)) {
    return false;
  }

  XWindowAttributes attributes;
  if (!XGetWindowAttributes(display, window, &attributes)) {
    return false;
  }

  out->x = windowX;
  out->y = windowY;
  out->width = attributes.width;
  out->height = attributes.height;

  // The pointer's position is known both relative to the window and to the
  // root, and the difference between them is where the window is.
  out->windowX = rootX - windowX;
  out->windowY = rootY - windowY;
  return true;
}
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
		x11, isX11 := context.(driver.X11WindowContext)
		if !isX11 || x11.WindowHandle == 0 {
			return
		}
		ok = bool(C.geyesPointerInWindow(C.uintptr_t(x11.WindowHandle), &found))
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
