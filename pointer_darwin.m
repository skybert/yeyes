#import <Cocoa/Cocoa.h>

#include "pointer_darwin.h"

bool yeyesPointerInWindow(uintptr_t nsWindow, yeyesPointer *out) {
  NSWindow *window = (NSWindow *)(void *)nsWindow;
  if (window == nil) {
    return false;
  }

  @autoreleasepool {
    // Cocoa hands out the pointer location whether or not it is over one of our
    // windows, which is what lets the eyes follow it across the whole screen.
    NSPoint pointer = [NSEvent mouseLocation];
    NSRect content = [window contentRectForFrameRect:[window frame]];

    // The pointer and the window are both in screen coordinates with the origin
    // in the bottom left corner, so y is flipped to the top left one that Fyne
    // and GLFW use. GLFW measures window positions from the top of the main
    // display, so that is where the flip starts from.
    double screenHeight = CGDisplayBounds(CGMainDisplayID()).size.height;

    out->x = pointer.x - NSMinX(content);
    out->y = NSMaxY(content) - pointer.y;
    out->width = NSWidth(content);
    out->height = NSHeight(content);
    out->windowX = NSMinX(content);
    out->windowY = screenHeight - NSMaxY(content);
  }

  return true;
}
