package VNC

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	//"unsafe"
)

const (
	v8                             = "RFB 003.008\n"
	SetPixelFormat           uint8 = 0
	SetEncodings             uint8 = 2
	FramebufferUpdateRequest uint8 = 3
	KeyEvent                 uint8 = 4
	PointerEvent             uint8 = 5
	SecurityTypeNone         uint8 = 1
)

func init() {
	PixelFormatRGBA = PixelFormat{} // init this shiet
}

// vnc client connection

type Client struct {
	conn             *Conn
	server           *VNCServer
	disconnectChan   chan struct{}
	errorChan        chan struct{}
	pixFormat        *PixelFormat
	securityType     uint8
	encoding         int32
	frameUpdatesChan chan []byte // here frames from the main servers screenshot service are received and
	writeChan        chan []byte
}

func newClient(server *VNCServer, conn net.Conn) *Client {
	var res Client
	res.server = server
	res.conn = newConn(conn)

	res.disconnectChan = make(chan struct{}, ChanBufferSize)
	res.errorChan = make(chan struct{}, ChanBufferSize)
	res.writeChan = make(chan []byte, ChanBufferSize)
	res.frameUpdatesChan = make(chan []byte, ChanBufferSize)
	return &res
}

func (c *Client) ReadBytes(count int) (res []byte) {
	return c.conn.readBytes(count)
}

func (c *Client) Read(value interface{}) {
	c.conn.read(value)
}

func (c *Client) ReadPadding(count int) {
	c.conn.readPadding(count)
}

func (c *Client) Write(value interface{}) {
	c.conn.write(value)
}

func (c *Client) WriteBytes(data []byte) {
	c.conn.writeBytes(data)
}

func (c *Client) WriteString(str string) {
	c.conn.writeString(str)
}

func (c *Client) SendFrameUpdate(img []byte) {
	c.Write(FramebufferUpdateRequest) // message type

	c.Write(FramebufferUpdateRequest) // padding
	var numRects uint16 = 1
	c.Write(numRects)

	rect := FrameUpdateRect{}
	rect.EncodingType = EncodingTypeRaw
	width, height, _ := getScreenDimensions()
	rect.Width = uint16(width)
	rect.Height = uint16(height)
	rect.X = uint16(0)
	rect.Y = uint16(0)
	c.Write(&rect)
	c.WriteBytes(img)
}

func (c *Client) writeService() {
	for {
		select {
		case img := <-c.frameUpdatesChan:
			//log.Println("frame update serv conn")
			// send frame update
			c.SendFrameUpdate(img)
		}
	}
}

func (this *Client) readService() {
	var messageType uint8
	var e error

	for {
		// handle further requests
		// reading from the connection must be synchronous
		this.Read(&messageType)
		switch messageType {
		case SetPixelFormat:
			e = this.handlePixelFormat()
		case SetEncodings:
			e = this.handleEncodings()
		case FramebufferUpdateRequest:
			e = this.handleFrameBufferUpdateRequest()
		case KeyEvent:
			e = this.handleKeyEvent()
		case PointerEvent:
			e = this.handlePointerEvent()
		default:
			log.Println("unsupported message type: ", messageType)
			return
		}

		if e != nil {
			failFatal("handle message failed")
		}
	}
}

func (this *Client) handle() {
	// TODO handle connection errors
	var e error
	e = this.handshake()
	if e != nil {
		failFatal("handshake failed")
	}
	e = this.init()
	if e != nil {
		failFatal("init client failed")
	}

	go this.writeService()
	this.readService()
}

func isSupportedSecurityType(t uint8) bool {
	return t == SecurityTypeNone
}

func (this *Client) handshake() error {
	// send version message
	this.WriteBytes([]byte(v8))

	// wait for response
	res := this.ReadBytes(len([]byte(v8)))
	if string(res) != v8 {
		log.Println("unsupported version %s", string(res))
		return nil
	}

	// construct and send security types message
	var securityTypeMessage = make([]uint8, 2)
	securityTypeMessage[0] = 1
	securityTypeMessage[1] = SecurityTypeNone
	this.WriteBytes(securityTypeMessage)

	// wait for security response
	var securityType uint8
	this.Read(&securityType)
	if !isSupportedSecurityType(securityType) {
		log.Println("unsupported security type")
		// send error string to the client
		// disconnect the client
		return nil
	}

	// send security result
	var result uint32 = 1
	this.conn.write(&result)
	log.Println("server connected")
	return nil
}

func constructServerInit() []byte {
	// construct fields
	sInit := serverInit{}
	width, height, _ := getScreenDimensions()
	sInit.Width = uint16(width)
	sInit.Height = uint16(height)
	sInit.PixFormat = PixelFormatRGBA
	name := "BaiHoi"
	sInit.NameLen = uint32(len([]byte(name)))

	// write to buffer
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, &sInit)
	buf.WriteString(name)
	return buf.Bytes()
}

func (this *Client) init() error {
	//client init
	var shared byte
	this.Read(&shared) // doesnt matter what this is

	// server init
	servInit := constructServerInit()
	this.WriteBytes(servInit)

	return nil
}

func (this *Client) handlePixelFormat() error {
	pFormat := PixelFormat{}
	this.Read(&pFormat)
	// TODO: figure what to do with this message, for now just ignoring it and assume the format is RGBA
	return nil
}

func (this *Client) chooseEncoding(encodings []int32) int32 {
	for i := 0; i < len(encodings); i++ {
		if this.server.isSupportedEncoding(encodings[i]) {
			return encodings[i]
		}
	}
	return -1
}
func (this *Client) handleEncodings() error {
	this.ReadPadding(1)
	var encodingCount uint16
	this.Read(&encodingCount)

	encodings := make([]int32, encodingCount)
	this.Read(encodings)

	this.encoding = this.chooseEncoding(encodings)
	return nil
}

func (this *Client) handleFrameBufferUpdateRequest() error {
	request := FrameBufferRequest{}
	this.Read(&request)
	this.server.screenShotService.Request(&request)
	return nil
}

func (this *Client) handleKeyEvent() error {
	event := keyEvent{}
	this.Read(&event)
	this.server.keyboardService.Request(&event)
	return nil
}

func (this *Client) handlePointerEvent() error {
	event := pointerEvent{}
	this.Read(&event)

	this.server.pointerService.Request(&event)
	return nil
}
