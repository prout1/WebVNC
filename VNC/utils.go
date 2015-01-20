package VNC

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"unsafe"
)

func GetBytes(x interface{}, order binary.ByteOrder) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 100))

	v := reflect.ValueOf(x)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i).Interface()
		if reflect.TypeOf(field).Kind() == reflect.String {
			fmt.Println("writing string")
			buf.WriteString(field.(string))

		} else {
			binary.Write(buf, order, field)
		}

	}
	fmt.Println(unsafe.Sizeof(PixelFormatRGBA))
	return buf.Bytes(), nil
}

func FromBytes(bts []byte) (interface{}, error) {
	// works only for fixed size types, meaning types that dont contain strings and arrays
	return nil, nil
}

func failFatal(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}

func SwapTwoBytes(x uint16) uint16 {
	return (x << 8) | (x >> 8)
}

func SwapFourBytes(x uint32) uint32 {
	x = ((x << 8) & 0xFF00FF00) | ((x >> 8) & 0xFF00FF)
	return (x << 16) | (x >> 16)
}
