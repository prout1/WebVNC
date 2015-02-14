package VNC

import (
	"image"
	"log"
	"net"
)

var (
	MessageTypeFramebufferUpdate uint8  = FramebufferUpdateRequest
	EncodingTypeRaw              uint32 = 0
)

type VNCClient struct {
	conn            *Conn
	frameUpdateChan chan *image.RGBA
	errorChan       chan struct{}
}

func NewClient(server net.Conn) *VNCClient {
	res := &VNCClient{}
	res.conn = newConn(server)
	res.frameUpdateChan = make(chan *image.RGBA, ChanBufferSize)
	res.errorChan = make(chan struct{}, ChanBufferSize)
	if res.ClientInit() != nil {
		log.Println("ClientInit failed")
		return nil
	}
	return res
}

func (c *VNCClient) ClientInit() error {
	// handshake
	//version - request and response
	c.conn.writeString(v8)

	v := string(c.conn.readBytes(len([]byte(v8))))
	if v != v8 {
		log.Println("unsupported version, client")
	}

	// security type
	_ = c.conn.readBytes(2)

	c.conn.write(SecurityTypeNone)

	// security message
	var status uint32
	c.conn.read(&status)
	if status == 0 {
		log.Println("security type check failed")
	}

	// send client init
	var shared uint8 = 0
	c.conn.write(shared)

	// read server init
	servInit := serverInit{}
	c.conn.read(&servInit)
	name := string(c.conn.readBytes(int(servInit.NameLen)))
	log.Println(name)
	return nil
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
	req.Incremental = 0
	req.Height = height
	req.Width = width
	req.X = x
	req.Y = y

	c.conn.write(FramebufferUpdateRequest)
	c.conn.write(&req)
}

func (c *VNCClient) SendSetPixFormat() {
	// ignore this shit
}

func (c *VNCClient) SendKeyEvent(keycode uint32, press bool) {
	c.conn.write(KeyEvent)
	req := keyEvent{}
	req.KeyCode = keycode
	if press {
		req.DownFlag = 1
	} else {
		req.DownFlag = 0
	}

	c.conn.write(&req)
}

func (c *VNCClient) SendPointerEvent(x, y uint16, buttonMask uint8) {
	c.conn.write(PointerEvent)

	req := pointerEvent{}
	req.X = x
	req.Y = y
	req.ButtonMask = buttonMask

	c.conn.write(&req)
}

func (c *VNCClient) SendSetEncodings() {
	// ignore this shit
}

func isSupportedEncoding(enc uint32) bool {
	return enc == EncodingTypeRaw
}

func (c *VNCClient) getPixels(r *FrameUpdateRect) []byte {
	return c.conn.readBytes(int(r.Width) * int(r.Height) * 4)
}

func (c *VNCClient) handleFrameUpdate() {
	c.conn.readPadding(1)
	var numRects uint16 // should be 1
	c.conn.read(&numRects)

	//assemble full frame
	for i := 0; i < int(numRects); i++ {
		rect := FrameUpdateRect{}
		c.conn.read(&rect)
		if isSupportedEncoding(rect.EncodingType) {
			pixels := c.getPixels(&rect)
			img := image.NewRGBA(image.Rect(int(rect.X), int(rect.Y), int(rect.X+rect.Width), int(rect.Y+rect.Height)))
			img.Pix = []uint8(pixels)
			c.frameUpdateChan <- img
		}

		break // for now, because only raw encoding will be supported
	}
}
func (c *VNCClient) GetFrameUpdate() *image.RGBA {
	return <-c.frameUpdateChan
}
