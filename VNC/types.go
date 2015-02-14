package VNC

type PixelFormat struct {
	BitsPerPixel uint8
	Depth        uint8
	BigEndian    uint8
	TrueColor    uint8
	RedMax       uint16
	GreenMax     uint16
	BlueMax      uint16
	RedShift     uint8
	GreenShift   uint8
	BlueShift    uint8
	P1, P2, P3   uint8 // padding bytes
}

var (
	PixelFormatRGBA PixelFormat
)

type serverInit struct {
	Width, Height uint16
	PixFormat     PixelFormat
	NameLen       uint32
}

type FrameBufferRequest struct {
	Incremental         uint8
	X, Y, Width, Height uint16
}

type FrameUpdateRect struct {
	X, Y          uint16
	Width, Height uint16
	EncodingType  uint32
}

type keyEvent struct {
	DownFlag uint8
	P1, P2   uint8
	KeyCode  uint32
}

type pointerEvent struct {
	ButtonMask uint8
	X, Y       uint16
}
