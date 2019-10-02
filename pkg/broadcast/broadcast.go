// Package broadcast is a simple udp broadcast
package broadcast

import (
	"net"

	"github.com/p9c/pod/pkg/log"
)

const (
	maxDatagramSize = 8192
)

// New creates a new UDP multicast connection on which to broadcast
func New(address string) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Listen binds to the UDP address and port given and writes packets received
// from that address to a buffer which is passed to a handler
func Listen(address string, handler func(*net.UDPAddr, int, []byte)) {
	// Parse the string address
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.ERROR(err)
	}

	// Open up a connection
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.ERROR(err)
	}

	err = conn.SetReadBuffer(maxDatagramSize)
	if err != nil {
		log.ERROR(err)
	}

	// Loop forever reading from the socket
	for {
		buffer := make([]byte, maxDatagramSize)
		numBytes, src, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.ERROR("ReadFromUDP failed:", err)
		}

		handler(src, numBytes, buffer)
	}
}
