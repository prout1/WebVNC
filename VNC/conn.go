package VNC

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"sync"
	//"unsafe"
)

type Conn struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	wm     sync.Mutex
}

func newConn(conn net.Conn) *Conn {
	var res Conn
	res.conn = conn
	res.reader = bufio.NewReader(conn)
	res.writer = bufio.NewWriter(conn)
	return &res
}

func (c *Conn) readBytes(count int) (res []byte) {
	res = make([]byte, count)
	n, e := io.ReadFull(c.reader, res)
	if e != nil || n != count {
		failFatal("failed to read %d bytes from connection", count)
	}
	return res
}

func (c *Conn) read(value interface{}) {
	e := binary.Read(c.reader, binary.BigEndian, value)
	if e != nil {
		failFatal("failed to read value from connection: ", e)
	}
}

func (c *Conn) readPadding(count int) {
	_ = c.readBytes(count)
}

func (c *Conn) write(value interface{}) {
	e := binary.Write(c.writer, binary.BigEndian, value)
	if e != nil {
		failFatal("write failed: ", e)
	}
	c.writer.Flush()
}

func (c *Conn) writeBytes(data []byte) {
	n, e := c.writer.Write(data)
	if e != nil || n != len(data) {
		failFatal("writeBytes failed: ", e)
	}
	c.writer.Flush()
}

func (c *Conn) writeString(str string) {
	n, e := c.writer.WriteString(str)
	if e != nil || n != len(str) {
		failFatal("writeBytes failed: ", e)
	}
	c.writer.Flush()
}
