package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <stdio.h>
#import <Cocoa/Cocoa.h>

void
MoveMouseSomewhere(int x, int y) {
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
MouseDown(int x, int y) {
    CGPoint p = {x, y};
    CGEventRef theEvent = CGEventCreateMouseEvent(NULL, NX_LMOUSEDOWN, p, kCGMouseButtonLeft);
    CGEventSetType(theEvent, NX_LMOUSEDOWN);
    CGEventPost(kCGHIDEventTap, theEvent);
    CFRelease(theEvent);
}

void
MouseUp(int x, int y) {
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

func MoveMouse() (x, y C.int) {
	C.MoveMouseSomewhere(x, y)
	return
}

func mouseDown(x, y C.int) {
	C.MouseDown(x, y)
}

func mouseUp(x, y C.int) {
	C.MouseUp(x, y)
}

func Click(x, y C.int) {
	mouseDown(x, y)
	time.Sleep(200 * time.Millisecond)
	mouseUp(x, y)
}
