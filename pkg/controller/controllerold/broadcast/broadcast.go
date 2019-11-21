// Package broadcast is a simple udp broadcast
package broadcast

import (
	"context"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/p9c/pod/pkg/controller/controllerold/fec"
	"io"
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
func Listen(address string, handler func(*net.UDPAddr, int, []byte)) (cancel context.CancelFunc) {
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

// Encode creates Reed Solomon shards and encrypts them using
// the provided GCM cipher function (from pkg/gcm).
// Message type is given in the first byte of each shard so nodes can quickly
// eliminate erroneous or irrelevant messages
func Encode(ciph cipher.AEAD, bytes []byte, typ []byte) (shards [][]byte,
	err error) {
	if len(bytes) > 1<<32 {
		log.WARN("GCM ciphers should only encode a maximum of 4gb per nonce" +
			" per key")
	}
	var clearText [][]byte
	clearText, err = fec.Encode(bytes)
	if err != nil {
		log.ERROR(err)
		return
	}
	// the nonce groups a broadcast's pieces,
	// the listener will gather them by this criteria.
	// The decoder assumes this but a message can be identified by
	// its' nonce due to using the same for each piece of the message
	nonce := make([]byte, ciph.NonceSize())
	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.ERROR(err)
		return
	}
	for i := range clearText {
		shards = append(shards, append(append(typ, nonce...), ciph.Seal(nil, nonce,
			clearText[i], nil)...))
	}
	log.SPEW(shards)
	return
}

func Decode(ciph cipher.AEAD, shards [][]byte) (bytes []byte, err error) {
	plainShards := make([][]byte, len(shards))
	nonceSize := ciph.NonceSize()
	for i := range shards {
		if len(shards[i]) < nonceSize {
			errMsg := []interface{}{"shard size too small, got",
				len(shards[i]), "expected minimum", nonceSize}
			log.ERROR(errMsg...)
			return nil, errors.New(fmt.Sprintln(errMsg...))
		}
		nonce, cipherText := shards[i][:nonceSize], shards[i][nonceSize:]
		var plaintext []byte
		plaintext, err = ciph.Open(nil, nonce, cipherText, nil)
		if err != nil {
			log.ERROR(err)
			return
		}
		plainShards[i] = plaintext
	}
	return fec.Decode(plainShards)
}
