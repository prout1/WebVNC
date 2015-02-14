package VNC

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"image"
	"image/jpeg"
	"log"
	"net"
	"net/http"
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

var addr = flag.String("addr", ":8080", "http service address")

func CreateProxy(tcpConn net.Conn) *WSProxy {
	res := WSProxy{}
	res.client = NewClient(tcpConn)
	go res.client.Run()
	res.errorChan = make(chan struct{}, ChanBufferSize)
	return &res
}

func dataToJson(data []byte) map[string]interface{} {
	//fmt.Println(string(data))
	res := make(map[string]interface{})
	json.Unmarshal(data, &res)
	return res
}

func (p *WSProxy) Run() {
	// getting the ws conn
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			}}

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		// close the connection maybe
		p.wsConn = ws
		go p.writer()
		p.reader()
	})

	log.Fatal(http.ListenAndServe(*addr, nil))
}

func (p *WSProxy) reader() {
	for {
		_, data, err := p.wsConn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
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

func (p *WSProxy) writer() {
	for {
		p.sendFrameUpdate()
	}
}

func getButtonMask(jsonEntry interface{}) uint8 {
	var buttonMask uint8 = 0
	buttonMap := jsonEntry.(map[string]interface{})
	if buttonMap["left"].(bool) {
		buttonMask |= (1 << 0)
	}
	if buttonMap["mid"].(bool) {
		buttonMask |= (1 << 1)
	}
	if buttonMap["right"].(bool) {
		buttonMask |= (1 << 2)
	}
	if buttonMap["up"].(bool) {
		buttonMask |= (1 << 3)
	}
	if buttonMap["down"].(bool) {
		buttonMask |= (1 << 4)
	}
	return buttonMask
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

func getJpegBuffer(img *image.RGBA) string {
	buf := bytes.Buffer{}
	jpeg.Encode(&buf, img, nil)
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}
func (p *WSProxy) sendFrameUpdate() {
	update := p.client.GetFrameUpdate()
	png := getJpegBuffer(update)
	if err := p.wsConn.WriteMessage(websocket.TextMessage, []byte(png)); err != nil {
		fmt.Println(err)
		return
	}
}
