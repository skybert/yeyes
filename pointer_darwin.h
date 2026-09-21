#ifndef GEYES_POINTER_DARWIN_H
#define GEYES_POINTER_DARWIN_H

#include <stdbool.h>
#include <stdint.h>

// yeyesPointer holds where the mouse pointer is relative to the top left corner
// of a window's content, the size of that content, and where it sits on the
// desktop. All values are in points.
typedef struct {
  double x, y;
  double width, height;
  double windowX, windowY;
} yeyesPointer;

// yeyesPointerInWindow fills in out for the given NSWindow. Must be called on
// the main thread.
bool yeyesPointerInWindow(uintptr_t nsWindow, yeyesPointer *out);

#endif
