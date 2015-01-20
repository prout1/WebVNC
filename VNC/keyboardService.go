package VNC

import ()

/*
#include <stdio.h>
#include <stdlib.h>

#include <windows.h>

void PressKey(unsigned int key)
{
  // Simulate a key press
     keybd_event( key,
                  0,
                  KEYEVENTF_EXTENDEDKEY,
                  0 );
}

void ReleaseKey(unsigned int key)
{
  // Simulate a key release
     keybd_event( key,
                  0,
                  KEYEVENTF_KEYUP,
                  0);
}
*/
import "C"

type keybdService struct {
	server     *VNCServer
	requests   chan *keyEvent
	xToVKeyMap map[uint32]uint32
}

func (s *keybdService) Run() {
	for {
		select {
		case req := <-s.requests:
			s.handleKeyEvent(req.keyCode, req.downFlag == 0)
		}
	}
}

func (s *keybdService) Request(event *keyEvent) {
	s.requests <- event
}

func (s *keybdService) Stop() {

}

func (s *keybdService) Init() {
	s.xToVKeyMap = make(map[uint32]uint32)
	s.requests = make(chan *keyEvent, ChanBufferSize)
	// TODO manage keycodes
}

func PressKey(keyCode uint32) {
	//6to ne ba4ka6 maika ti deiba
	C.PressKey(C.uint(keyCode))
}

func ReleaseKey(keyCode uint32) {
	C.ReleaseKey(C.uint(keyCode))
}

func (s *keybdService) getVKey(keyCode uint32) uint32 {
	return keyCode
}

func (s *keybdService) handleKeyEvent(keyCode uint32, pressFlag bool) {
	vKey := s.getVKey(keyCode)
	if pressFlag {
		PressKey(vKey)
	} else {
		ReleaseKey(vKey)
	}
}
