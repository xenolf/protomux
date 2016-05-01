package protomux

import (
	"bytes"
	"io"
	"net"
)

const (
	preReadBufferSize = 128
)

// ProtoConn wraps a net.Conn for content introspection
type ProtoConn struct {
	net.Conn
	buffer   *bytes.Buffer
	protocol Protocol
}

func wrapConn(conn net.Conn) *ProtoConn {
	return &ProtoConn{
		Conn:     conn,
		buffer:   bytes.NewBuffer(make([]byte, 0, preReadBufferSize)),
		protocol: None,
	}
}

func (c *ProtoConn) Read(b []byte) (n int, err error) {
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

// GetProtocol returns the determined protocol.
func (c *ProtoConn) GetProtocol() Protocol {
	return c.protocol
}

func (c *ProtoConn) fillBuffer() {
	buf := make([]byte, preReadBufferSize)
	c.Conn.Read(buf)
	c.buffer.Write(buf)
}

func (c *ProtoConn) setProtocol(newType Protocol) {
	c.protocol = newType
}
