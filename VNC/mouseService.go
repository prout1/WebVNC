package VNC

import (
	"log"
)

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
	// TODO try post message way
	DWORD lflags = 0;
	INPUT input;
	DWORD lflags2 = 0;

	DWORD wheelMove = (flags & MOUSE_TYPE_SCROLL_UP) ? WHEEL_DELTA :
		((flags & MOUSE_TYPE_SCROLL_DOWN) ? -WHEEL_DELTA : 0);
	lflags |= (flags & MOUSE_TYPE_LEFT) ? MOUSEEVENTF_LEFTDOWN :
		((state & MOUSE_TYPE_LEFT) ? MOUSEEVENTF_LEFTUP : 0);
	lflags |= (flags & MOUSE_TYPE_RIGHT) ? MOUSEEVENTF_RIGHTDOWN :
		((state & MOUSE_TYPE_RIGHT) ? MOUSEEVENTF_RIGHTUP : 0);
	lflags |= (flags & MOUSE_TYPE_MIDDLE) ? MOUSEEVENTF_MIDDLEDOWN :
		((state & MOUSE_TYPE_MIDDLE) ? MOUSEEVENTF_MIDDLEUP : 0);
	lflags |= (flags & MOUSE_TYPE_SCROLL_UP) ? MOUSEEVENTF_WHEEL : 0;
	lflags |= (flags & MOUSE_TYPE_SCROLL_DOWN) ? MOUSEEVENTF_WHEEL : 0;

	lflags |= MOUSEEVENTF_MOVE | MOUSEEVENTF_ABSOLUTE;
	mouse_event(
		lflags,
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
	//log.Println(request)
	//log.Println("State: ", s.stateFlags)
	C.HandleMouse(C.ushort(request.X), C.ushort(request.Y), C.uchar(request.ButtonMask), C.uchar(s.stateFlags))
	s.stateFlags = request.ButtonMask
}

func TestRat() {
	var e = pointerEvent{}
	e.X = 30000
	e.Y = 30000
	e.ButtonMask = C.MOUSE_TYPE_LEFT
	var state uint8 = 0

	C.HandleMouse(C.ushort(e.X), C.ushort(e.Y), C.uchar(e.ButtonMask), C.uchar(state))

	C.Sleep(2000)
	e.ButtonMask = 0
	C.HandleMouse(C.ushort(e.X), C.ushort(e.Y), C.uchar(e.ButtonMask), C.uchar(state))
	log.Println("done")
	C.Sleep(1000)
}
