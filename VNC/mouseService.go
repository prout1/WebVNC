package VNC

/*
#include <stdio.h>
#include <stdlib.h>

#include <windows.h>

#define MOUSE_TYPE_LEFT 1
#define MOUSE_TYPE_MIDDLE 2
#define MOUSE_TYPE_RIGHT 4
#define MOUSE_TYPE_SCROLL_UP 8
#define MOUSE_TYPE_SCROLL_DOWN 16

void HandleMouse(unsigned short x, unsigned short y, unsigned char flags, unsigned char state)
{
	DWORD lflags = 0;
	DWORD wheelMove = (flags & MOUSE_TYPE_SCROLL_UP) ? WHEEL_DELTA :
		((flags & MOUSE_TYPE_SCROLL_DOWN) ? -WHEEL_DELTA : 0);
	lflags |= ((flags & MOUSE_TYPE_LEFT) ?
		((state & MOUSE_TYPE_LEFT) ? 0 : MOUSEEVENTF_LEFTDOWN) :
		((state & MOUSE_TYPE_LEFT) ? 0 : MOUSEEVENTF_LEFTUP));
	lflags |= ((flags & MOUSE_TYPE_RIGHT) ?
		((state & MOUSE_TYPE_RIGHT) ? 0 : MOUSEEVENTF_RIGHTDOWN) :
		((state & MOUSE_TYPE_RIGHT) ? 0 : MOUSEEVENTF_RIGHTUP));
	lflags |= ((flags & MOUSE_TYPE_MIDDLE) ?
		((state & MOUSE_TYPE_MIDDLE) ? 0 : MOUSEEVENTF_MIDDLEDOWN) :
		((state & MOUSE_TYPE_MIDDLE) ? 0 : MOUSEEVENTF_MIDDLEUP));
	lflags |= (flags & MOUSE_TYPE_SCROLL_UP) ? MOUSEEVENTF_WHEEL : 0;
	lflags |= (flags & MOUSE_TYPE_SCROLL_DOWN) ? MOUSEEVENTF_WHEEL : 0;

	lflags |= MOUSEEVENTF_MOVE | MOUSEEVENTF_ABSOLUTE;

	mouse_event(
		flags,
		x,
		y,
		wheelMove,
		0
		);
}

*/
import "C"

// mouse input
type mouseService struct {
	requests   chan *pointerEvent
	stateFlags uint8 // flags indicating currently pressed keys and stuff
}

func (s *mouseService) Stop() {

}

func (s *mouseService) Init() {
	s.requests = make(chan *pointerEvent, ChanBufferSize)
	s.stateFlags = 0
	// TODO manage keycodes
}

func (s *mouseService) Run() {
	for {
		select {
		case req := <-s.requests:
			s.handleMouseEvent(req)
		}
	}
}

func (s *mouseService) Request(req *pointerEvent) {
	s.requests <- req
}

func (s *mouseService) handleMouseEvent(request *pointerEvent) {
	C.HandleMouse(C.ushort(request.x), C.ushort(request.y), C.uchar(request.buttonMask), C.uchar(s.stateFlags))
	s.stateFlags = request.buttonMask
}
