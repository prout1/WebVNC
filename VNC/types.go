package VNC

type PixelFormat struct {
	bitsPerPixel uint8
	depth        uint8
	bigEndian    uint8
	trueColor    uint8
	redMax       uint16
	greenMax     uint16
	blueMax      uint16
	redShift     uint8
	greenShift   uint8
	blueShift    uint8
	p1, p2, p3   uint8 // padding bytes
}

var (
	PixelFormatRGBA PixelFormat
)

type serverInit struct {
	width, height uint16
	pixFormat     PixelFormat
	nameLen       uint32
}

type FrameBufferRequest struct {
	incremental         uint8
	x, y, width, height uint16
}

type keyEvent struct {
	downFlag uint8
	p1, p2   uint8
	keyCode  uint32
}

type pointerEvent struct {
	buttonMask uint8
	x, y       uint16
}
