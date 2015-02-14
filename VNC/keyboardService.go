package VNC

import (
//"log"
)

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
			s.handleKeyEvent(req.KeyCode, req.DownFlag == 0)
		}
	}
}

func (s *keybdService) Request(event *keyEvent) {
	s.requests <- event
}

func (s *keybdService) Stop() {

}

func FillMap(m map[uint32]uint32) {
	m[8] = C.VK_BACK     // backspace
	m[9] = C.VK_TAB      // tab
	m[13] = C.VK_RETURN  // enter
	m[16] = C.VK_SHIFT   // shift
	m[17] = C.VK_CONTROL // ctrl
	m[18] = C.VK_MENU    // alt key
	m[19] = C.VK_PAUSE
	m[20] = C.VK_CAPITAL // caps lock
	m[27] = C.VK_ESCAPE
	m[32] = C.VK_SPACE //space bar
	m[33] = C.VK_PRIOR //page up
	m[34] = C.VK_NEXT  //page down
	m[35] = C.VK_END
	m[36] = C.VK_HOME
	m[37] = C.VK_LEFT //arrows
	m[38] = C.VK_UP
	m[39] = C.VK_RIGHT
	m[40] = C.VK_DOWN
	m[45] = C.VK_INSERT
	m[46] = C.VK_DELETE
	m[48] = 0x30 // digits 0-9
	m[49] = 0x31
	m[50] = 0x32
	m[51] = 0x33
	m[52] = 0x34
	m[53] = 0x35
	m[54] = 0x36
	m[55] = 0x37
	m[56] = 0x38
	m[57] = 0x39 // 9
	m[65] = 0x41 // letters a-z
	m[66] = 0x42
	m[67] = 0x43
	m[68] = 0x44
	m[69] = 0x45
	m[70] = 0x46
	m[71] = 0x47
	m[72] = 0x48
	m[73] = 0x49
	m[74] = 0x4A
	m[75] = 0x4b
	m[76] = 0x4c
	m[77] = 0x4d
	m[78] = 0x4e
	m[79] = 0x4f
	m[80] = 0x50
	m[81] = 0x51
	m[82] = 0x52
	m[83] = 0x53
	m[84] = 0x54
	m[85] = 0x55
	m[86] = 0x56
	m[87] = 0x57
	m[88] = 0x58
	m[89] = 0x59
	m[90] = 0x5a      // z
	m[91] = C.VK_LWIN // left window key
	m[92] = C.VK_RWIN
	m[93] = C.VK_SELECT
	m[96] = C.VK_NUMPAD0
	m[97] = C.VK_NUMPAD1
	m[98] = C.VK_NUMPAD2
	m[99] = C.VK_NUMPAD3
	m[100] = C.VK_NUMPAD4
	m[101] = C.VK_NUMPAD5
	m[102] = C.VK_NUMPAD6
	m[103] = C.VK_NUMPAD7
	m[104] = C.VK_NUMPAD8
	m[105] = C.VK_NUMPAD9
	m[106] = C.VK_MULTIPLY // numpad operations
	m[107] = C.VK_ADD
	m[109] = C.VK_SUBTRACT
	m[110] = C.VK_DECIMAL // decimal point
	m[111] = C.VK_DIVIDE
	m[112] = C.VK_F1
	m[113] = C.VK_F2
	m[114] = C.VK_F3
	m[115] = C.VK_F4
	m[116] = C.VK_F5
	m[117] = C.VK_F6
	m[118] = C.VK_F7
	m[119] = C.VK_F8
	m[120] = C.VK_F9
	m[121] = C.VK_F10
	m[122] = C.VK_F11
	m[123] = C.VK_F12
	m[144] = C.VK_NUMLOCK
	m[145] = C.VK_SCROLL
	// tuka stava mnogo strashno
	m[186] = C.VK_OEM_1    // semicolon
	m[187] = C.VK_OEM_PLUS // equal sign or plus
	m[188] = C.VK_OEM_COMMA
	m[189] = C.VK_OEM_MINUS
	m[190] = C.VK_OEM_PERIOD
	m[191] = C.VK_OEM_2 // forward slash ('/')
	m[192] = C.VK_OEM_3 // grave accent key ('~')
	m[219] = C.VK_OEM_4 // open bracket ('[')
	m[220] = C.VK_OEM_5 // backslach ('\')
	m[221] = C.VK_OEM_6 // close bracket (']')
	m[222] = C.VK_OEM_7 // single quote (''' lol)
}

func (s *keybdService) Init() {
	s.xToVKeyMap = make(map[uint32]uint32)
	FillMap(s.xToVKeyMap)
	s.requests = make(chan *keyEvent, ChanBufferSize)
	// TODO manage keycodes
}

func PressKey(keyCode uint32) {
	C.PressKey(C.uint(keyCode))
}

func ReleaseKey(keyCode uint32) {
	C.ReleaseKey(C.uint(keyCode))
}

func (s *keybdService) getVKey(keyCode uint32) uint32 {
	return s.xToVKeyMap[keyCode]
}

func (s *keybdService) handleKeyEvent(keyCode uint32, pressFlag bool) {
	//log.Println(keyCode, pressFlag)
	vKey := s.getVKey(keyCode)
	if pressFlag {
		PressKey(vKey)
	} else {
		ReleaseKey(vKey)
	}
}
