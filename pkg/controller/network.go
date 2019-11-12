package controller

import (
	"context"
	"crypto/rand"
	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/fec"
	"github.com/p9c/pod/pkg/log"
	"io"
	"net"
)

const (
	// MaxDatagramSize is the largest a packet could be,
	// it is a little larger but this is easier to calculate.
	// There is only one listening thread but it needs a buffer this size for
	// worst case largest block possible.
	// Note also this is why FEC is used on the packets in case some get lost it
	// has to puncture 6 of the 9.
	// This protocol is connectionless and stateless so if one misses,
	// the next one probably won't, usually a second or 3 later
	MaxDatagramSize      = blockchain.MaxBlockBaseSize / 3
	UDP6MulticastAddress = "ff02::1"
	UDP4MulticastAddress = "224.0.0.1"
)

var (
	MCAddresses = []*net.UDPAddr{
		{IP: net.ParseIP(UDP6MulticastAddress), Port: 11049},
		{IP: net.ParseIP(UDP4MulticastAddress), Port: 11049},
	}
)

// Send broadcasts bytes on the given multicast address with each shard
// labeled with a random 32 bit nonce to identify its group to the listener's
// handler function
func Send(addr *net.UDPAddr, bytes []byte, magic [4]byte) (err error) {
	var shards [][]byte
	shards, err = fec.Encode(bytes)
	if err != nil {
		return
	}
	//log.SPEW(shards)
	// nonce is a batch identifier for the FEC encoded shards which are sent
	// out as individual packets
	nonce := make([]byte, 4)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.ERROR(err)
		return
	}
	prefix := append(nonce, magic[:]...)
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return
	}
	var n, cumulative int
	for i := range shards {
		n, err = conn.WriteToUDP(append(prefix, shards[i]...), addr)
		if err != nil {
			log.ERROR(err, len(shards[i]))
			return
		}
		cumulative += n
	}
	log.TRACE("wrote", cumulative, "bytes to multicast address", addr.IP,
		"port",
		addr.Port)
	err = conn.Close()
	if err != nil {
		log.ERROR(err)
	}
	return
}

// Listen binds to the UDP address and port given and writes packets received
// from that address to a buffer which is passed to a handler
func Listen(address *net.UDPAddr, handler func(*net.UDPAddr, int,
	[]byte)) (cancel context.CancelFunc, err error) {
	var ctx context.Context
	ctx, cancel = context.WithCancel(context.Background())
	log.DEBUG("resolving", address)
	//addr, err := net.ResolveUDPAddr("udp", address)
	//if err != nil {
	//	log.ERROR(err)
	//	cancel()
	//	return
	//}
	//log.DEBUG("resolved", addr.IP, addr.Port, addr.String())
	var conn *net.UDPConn
	conn, err = net.ListenUDP("udp", address)
	if err != nil {
		log.ERROR(err)
		cancel()
		return
	}

	err = conn.SetReadBuffer(MaxDatagramSize)
	if err != nil {
		log.ERROR(err)
	}
	buffer := make([]byte, MaxDatagramSize)
	go func() {
	out:
		// read from socket until context is cancelled
		for {
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
