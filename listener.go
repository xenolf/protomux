// Package protomux implements a simple means to multiplex TCP protocols over the same port.
package protomux

import "net"

// Protocol ...
type Protocol int

// Interface for protocol matching
type probe interface {
	apply([]byte) bool
}

const (
	// None - no protocol has been matched
	None Protocol = iota
	// HTTP protocol has been matched
	HTTP Protocol = iota
	// TLS protocol has been matched
	TLS Protocol = iota
)

var standardProbes = map[Protocol]probe{
	HTTP: httpProbe{},
	TLS:  tlsProbe{},
}

type Listener struct {
	net.Listener
	probes map[Protocol]probe
}

// Listen will start a TCP listener on the supplied address.
func Listen(addr string) (net.Listener, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Listener{listener, standardProbes}, nil
}

// Accept waits for and returns the next connection to the listener
// after applying protocol detection probes to it. The net.Conn returned
// will be untouched with all bytes ready to read.
func (p *Listener) Accept() (net.Conn, error) {
	conn, err := p.Listener.Accept()
	if err != nil {
		return nil, err
	}

	pConn := wrapConn(conn)
	pConn.fillBuffer()

	// in the case when no probe matches, the proto will simply
	// remain at None and the user can decide how to proceed.
	for proto, probe := range p.probes {
		if probe.apply(pConn.buffer.Bytes()) {
			pConn.setProtocol(proto)
			break
		}
	}

	return pConn, nil
}
