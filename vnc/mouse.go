package vnc

// This package wraps some Cocoa cursor handlers (objective-C) in Go
// Compile with CC=clang

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <stdio.h>
#import <Cocoa/Cocoa.h>

void
MoveMouseSomewhere(double x, double y) {
    CGPoint p = {x, y};
    CGError e = CGWarpMouseCursorPosition(p);
}

NSPoint
getCurPosition(void) {
NSPoint mouseLoc = [NSEvent mouseLocation];
return mouseLoc;
}

void
getMouse(double *x, double *y) {
    NSPoint mouseLoc = [NSEvent mouseLocation];
    *x = (double)mouseLoc.x;
    *y = (double)mouseLoc.y;
    return;
}

void
MouseDown(double x, double y) {
    CGPoint p = {x, y};
    CGEventRef theEvent = CGEventCreateMouseEvent(NULL, NX_LMOUSEDOWN, p, kCGMouseButtonLeft);
    CGEventSetType(theEvent, NX_LMOUSEDOWN);
    CGEventPost(kCGHIDEventTap, theEvent);
    CFRelease(theEvent);
}

void
MouseUp(double x, double y) {
    CGPoint p = {x, y};
    CGEventRef theEvent2 = CGEventCreateMouseEvent(NULL, NX_LMOUSEUP, p, kCGMouseButtonLeft);
    CGEventSetType(theEvent2, NX_LMOUSEUP);
    CGEventPost(kCGHIDEventTap, theEvent2);
    CFRelease(theEvent2);
}
*/
import "C"
import "time"

func getMouse() (x, y float64) {
	var cx, cy _Ctype_double
	C.getMouse(&cx, &cy)
	x = float64(cx)
	y = float64(cy)
	return
}

func MoveMouse(x, y float64) {
    cx, cy := C.double(x), C.double(y)
	C.MoveMouseSomewhere(cx, cy)
	return
}

func mouseDown(x, y float64) {
    cx, cy := C.double(x), C.double(y)
	C.MouseDown(cx, cy)
}

func mouseUp(x, y float64) {
    cx, cy := C.double(x), C.double(y)
	C.MouseUp(cx, cy)
}

func Click(x, y float64) {
	mouseDown(x, y)
	time.Sleep(200 * time.Millisecond)
	mouseUp(x, y)
}
