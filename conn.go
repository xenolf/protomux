package protomux

import (
	"bytes"
	"io"
	"net"
)

const (
	preReadBufferSize = 128
)

// Conn wraps a net.Conn for content introspection
type Conn struct {
	net.Conn
	buffer   *bytes.Buffer
	protocol Protocol
}

func wrapConn(conn net.Conn) *Conn {
	return &Conn{
		Conn:     conn,
		buffer:   bytes.NewBuffer(make([]byte, 0, preReadBufferSize)),
		protocol: None,
	}
}

func (c *Conn) Read(b []byte) (n int, err error) {
	// If the buffer is already drained, read from conn
	if c.buffer == nil {
		return c.Conn.Read(b)
	}

	n, err = c.buffer.Read(b)

	// If we encounter EOF while reading,
	// close the buffer and read the remaining bytes
	// from the connection directly.
	if err == io.EOF {
		c.buffer = nil
		var n2 int
		n2, err = c.Conn.Read(b[n:])
		n += n2
	}
	return
}

// Protocol returns the determined protocol.
func (c *Conn) Protocol() Protocol {
	return c.protocol
}

func (c *Conn) fillBuffer() {
	buf := make([]byte, preReadBufferSize)
	c.Conn.Read(buf)
	c.buffer.Write(buf)
}

func (c *Conn) setProtocol(newType Protocol) {
	c.protocol = newType
}
