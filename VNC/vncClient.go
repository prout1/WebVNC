package VNC

import (
	"image"
	"net"
)

var (
	MessageTypeFramebufferUpdate uint8  = 0
	EncodingTypeRaw              uint32 = 0
)

type VNCClient struct {
	conn            *Conn
	frameUpdateChan chan *image.RGBA
	errorChan       chan struct{}
}

func NewClient(server net.Conn) *VNCClient {
	res := VNCClient{}
	res.conn = newConn(server)
	res.frameUpdateChan = make(chan *image.RGBA, ChanBufferSize)
	res.errorChan = make(chan struct{}, ChanBufferSize)
	return &res
}

func (c *VNCClient) ClientInit() {

}

func (c *VNCClient) Run() {
	for {
		var msgType uint8
		c.conn.read(&msgType)
		switch msgType {
		case MessageTypeFramebufferUpdate:
			c.handleFrameUpdate()
		default:
			// unsupported message
			return
		}
	}
}
func (c *VNCClient) SendFrameBufferUpdateRequest(x, y, width, height uint16) {
	req := FrameBufferRequest{}
	req.incremental = 0
	req.height = height
	req.width = width
	req.x = x
	req.y = y

	c.conn.write(&req)
}

func (c *VNCClient) SendSetPixFormat() {
	// ignore this shit
}

func (c *VNCClient) SendKeyEvent(keycode uint32, press bool) {
	req := keyEvent{}
	req.keyCode = keycode
	if press {
		req.downFlag = 1
	} else {
		req.downFlag = 0
	}

	c.conn.write(&req)
}

func (c *VNCClient) SendPointerEvent(x, y uint16, buttonMask uint8) {
	req := pointerEvent{}
	req.x = x
	req.y = y
	req.buttonMask = buttonMask

	c.conn.write(&req)
}

func (c *VNCClient) SendSetEncodings() {
	// ignore this shit
}

type FrameUpdateRect struct {
	x, y          uint16
	width, height uint16
	encodingType  uint32
}

func isSupportedEncoding(enc uint32) bool {
	return enc == EncodingTypeRaw
}

func (c *VNCClient) getPixels(r *FrameUpdateRect) []byte {
	return c.conn.readBytes(int(r.width * r.height * 4))
}

func (c *VNCClient) handleFrameUpdate() {
	c.conn.readPadding(1)
	var numRects uint16 // should be 1
	c.conn.read(&numRects)

	//assemble full frame
	for i := 0; i < int(numRects); i++ {
		rect := FrameUpdateRect{}
		c.conn.read(&rect)

		if isSupportedEncoding(rect.encodingType) {
			pixels := c.getPixels(&rect)
			img := image.NewRGBA(image.Rect(int(rect.x), int(rect.y), int(rect.x+rect.width), int(rect.y+rect.height)))
			img.Pix = []uint8(pixels)
			c.frameUpdateChan <- img
		}

		break // for now, because only raw encoding will be supported
	}
}
func (c *VNCClient) GetFrameUpdate() *image.RGBA {
	return <-c.frameUpdateChan
}
