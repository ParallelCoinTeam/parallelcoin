// Package broadcast is a simple udp broadcast
package broadcast

import (
	"context"
	"crypto/cipher"
	"net"

	"github.com/p9c/pod/pkg/log"
)

const (
	MaxDatagramSize = 8192
	DefaultAddress  = "239.0.0.0:11042"
)

// for fast elimination of irrelevant messages a magic 64 bit word is used to
// identify relevant types of messages and 64 bits so the buffer is aligned
var (
	SolBlock = "solblock"
	TplBlock = "tplblock"
	Solution = []byte(SolBlock)
	Template = []byte(TplBlock)
)

// Send broadcasts bytes on the given multicast address
func Send(addr *net.UDPAddr, bytes []byte, ciph cipher.AEAD,
	typ []byte) (err error) {
	var shards [][]byte
	shards, err = Encode(ciph, bytes, typ)
	if err != nil {
		return
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return
	}
	var n, cumulative int
	for i := range shards {
		n, err = conn.WriteToUDP(shards[i], addr)
		if err != nil {
			log.ERROR(err, len(shards[i]))
			return
		}
		cumulative += n
	}
	log.TRACE("wrote", n, "bytes to multicast address", addr.IP, "port",
		addr.Port)
	err = conn.Close()
	if err != nil {
		log.ERROR(err)
	}
	return
}

// New creates a new UDP multicast address on which to broadcast
func New(address string) (addr *net.UDPAddr, err error) {
	addr, err = net.ResolveUDPAddr("udp", address)
	return
}

// Listen binds to the UDP address and port given and writes packets received
// from that address to a buffer which is passed to a handler
func Listen(address string, handler func(*net.UDPAddr, int,
	[]byte)) (cancel context.CancelFunc) {
	var ctx context.Context
	ctx, cancel = context.WithCancel(context.Background())
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.ERROR(err)
	}
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.ERROR(err)
	}

	err = conn.SetReadBuffer(MaxDatagramSize)
	if err != nil {
		log.ERROR(err)
	}
	go func() {
	out:
		// read from socket until context is cancelled
		for {
			buffer := make([]byte, MaxDatagramSize)
			numBytes, src, err := conn.ReadFromUDP(buffer)
			if err != nil {
				log.ERROR("ReadFromUDP failed:", err)
				continue
			}
			handler(src, numBytes, buffer)
			select {
			case <-ctx.Done():
				break out
			default:
			}
		}
	}()
	return
}
