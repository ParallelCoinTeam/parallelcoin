package transport

import (
	"errors"
	"net"
	"strings"
	
	"github.com/p9c/pod/pkg/log"
)

var Address = "224.0.0.1:11049"

// NewBroadcaster creates a new UDP multicast connection on which to broadcast
func NewBroadcaster(maxDatagramSize int) (conn *net.UDPConn, err error) {
	addr, err := net.ResolveUDPAddr("udp", Address)
	if err != nil {
		return nil, err
	}
	conn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}
	conn.SetWriteBuffer(maxDatagramSize)
	log.DEBUG("created new broadcaster connection on", conn.LocalAddr(),
		"->", conn.RemoteAddr())
	return
}

// Listen binds to the UDP Address and port given and writes packets received
// from that Address to a buffer which is passed to a hander
func Listen(maxDatagramSize int, handler func(*net.UDPAddr, int, []byte)) (
	conn *net.UDPConn, err error) {
	// Parse the string Address
	addr, err := net.ResolveUDPAddr("udp", Address)
	if err != nil {
		log.ERROR(err)
	}
	// Open up a connection
	conn, err = net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.ERROR(err)
	}
	if conn == nil {
		return nil, errors.New("connection closed")
	}
	err = conn.SetReadBuffer(maxDatagramSize)
	if err != nil {
		log.ERROR(err)
	}
	go func() {
		log.DEBUG("starting listener on ", Address)
		// Loop forever reading from the socket until it is closed
		for {
			buffer := make([]byte, maxDatagramSize)
			numBytes, src, err := conn.ReadFromUDP(buffer)
			if err != nil {
				ll := len(strings.Split(err.Error(), "use of closed network connection"))
				if ll == 2 {
					log.DEBUG("listener closed", Address)
					break
				}
				log.ERRORF("ReadFromUDP failed: '%s'", err)
			}
			handler(src, numBytes, buffer[:numBytes])
		}
	}()
	return
}
