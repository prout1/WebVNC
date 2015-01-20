package VNC

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	//"log"
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"net"
)

var (
	MsgPointerEvent       string = "pe"
	MsgKeyBoardEvent             = "ke"
	MsgFrameUpdateRequest        = "re"
	FullScreen            uint16 = 65535
)

type WSProxy struct {
	wsConn    *websocket.Conn
	client    *VNCClient
	errorChan chan struct{}
}

func CreateProxy(wsConn *websocket.Conn, tcpConn net.Conn) *WSProxy {
	res := WSProxy{}
	res.client = NewClient(tcpConn)
	res.errorChan = make(chan struct{}, ChanBufferSize)
	res.wsConn = wsConn
	return &res
}

func dataToJson(data []byte) map[string]interface{} {
	res := make(map[string]interface{})
	json.Unmarshal(data, &res)
	return res
}

func (p *WSProxy) Run() {
	for {
		var data []byte
		_ = websocket.Message.Receive(p.wsConn, &data)
		json := dataToJson(data)

		t, ok := json["type"]
		if !ok {
			//error, abort
		}

		switch t.(string) {
		case MsgPointerEvent:
			p.handlePointerEvent(json)
		case MsgKeyBoardEvent:
			p.handleKeybdEvent(json)
		case MsgFrameUpdateRequest:
			p.handleFrameUpdateRequest(json)
		default:
			//error
			return
		}
	}
}

func getButtonMask(jsonEntry interface{}) uint8 {
	return 0
}

func (p *WSProxy) handlePointerEvent(eventJson map[string]interface{}) {
	x := uint16(eventJson["x"].(float64))
	y := uint16(eventJson["y"].(float64))

	buttonMask := getButtonMask(eventJson["buttons"])
	p.client.SendPointerEvent(x, y, buttonMask)
}

func (p *WSProxy) handleKeybdEvent(eventJson map[string]interface{}) {
	keyCode := uint32(eventJson["keyCode"].(float64))
	press := eventJson["press"].(bool)
	p.client.SendKeyEvent(keyCode, press)
}

func (p *WSProxy) handleFrameUpdateRequest(eventJson map[string]interface{}) {
	p.client.SendFrameBufferUpdateRequest(0, 0, FullScreen, FullScreen)
}

func getPngBuffer(img *image.RGBA) string {
	buf := bytes.Buffer{}
	png.Encode(&buf, img)
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}
func (p *WSProxy) sendFrameUpdate() {
	update := p.client.GetFrameUpdate()
	png := getPngBuffer(update)
	websocket.Message.Send(p.wsConn, png)
}
