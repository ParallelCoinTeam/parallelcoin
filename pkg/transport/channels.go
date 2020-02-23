package transport

import (
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	
	"github.com/p9c/pod/pkg/fec"
	"github.com/p9c/pod/pkg/gcm"
	"github.com/p9c/pod/pkg/log"
)

const (
	UDPMulticastAddress = "224.0.0.1"
)

type (
	// HandlerFunc is a function that is used to process a received message
	HandlerFunc func(src *net.UDPAddr, dst string, count int, data []byte) (err error)
	Handlers    map[string]HandlerFunc
	Channel     struct {
		Sender          *net.UDPConn
		Receiver        *net.UDPConn
		MaxDatagramSize int
		ciph            cipher.AEAD
	}
)

// SetDestination changes the address the outbound connection of a channel directs to
func (c *Channel) SetDestination(dst string) (err error) {
	if c.Sender, err = NewSender(dst, c.MaxDatagramSize); log.Check(err) {
	}
	return
}

// Send fires off some data through the configured channel's outbound
func (c *Channel) Send(magic []byte, data []byte) (n int, err error) {
	if len(data) == 0 {
		err = errors.New("not sending empty packet")
		log.ERROR(err)
		return
	}
	var msg []byte
	if msg, err = encryptMessage(magic, c.ciph, data); log.Check(err) {
	}
	n, err = c.Sender.Write(msg)
	return
}

func (c *Channel) SendMany(magic []byte, b BufIter) (err error) {
	for ; b.More(); b.Next() {
		if _, err = c.Send(magic, b.Get()); log.Check(err) {
		}
	}
	return
}

// Close the channel
func (c *Channel) Close() (err error) {
	if err = c.Sender.Close(); log.Check(err) {
	}
	if err = c.Receiver.Close(); log.Check(err) {
	}
	return
}

// GetShards returns a buffer iterator to feed to Channel.SendMany containing
// fec encoded shards built from the provided buffer
func GetShards(data []byte) (bI BufIter) {
	var err error
	if bI.Buf, err = fec.Encode(data); log.Check(err) {
	}
	return
}

// NewUnicastChannel sets up a listener and sender for a specified destination
func NewUnicastChannel(key, sender, receiver string, maxDatagramSize int,
	handlers Handlers) (channel *Channel, err error) {
	channel = &Channel{MaxDatagramSize: maxDatagramSize}
	if key != "" {
		if channel.ciph, err = gcm.GetCipher(key); log.Check(err) {
		}
	}
	ready := make(chan struct{})
	go func() {
		channel.Receiver, err = Listen(channel.ciph, receiver, maxDatagramSize, handlers)
		ready <- struct{}{}
	}()
	<-ready
	channel.Sender, err = NewSender(sender, maxDatagramSize)
	if err != nil {
		log.ERROR(err)
	}
	return
}

// NewSender creates a new UDP connection to a specified address
func NewSender(address string, maxDatagramSize int) (conn *net.UDPConn, err error) {
	var addr *net.UDPAddr
	if addr, err = net.ResolveUDPAddr("udp", address); log.Check(err) {
		return
	} else if conn, err = net.DialUDP("udp", nil, addr); log.Check(err) {
		return
	}
	log.DEBUG("started new sender on", conn.LocalAddr(), "->", conn.RemoteAddr())
	if err = conn.SetWriteBuffer(maxDatagramSize); log.Check(err) {
	}
	return
}

// Listen binds to the UDP Address and port given and writes packets received
// from that Address to a buffer which is passed to a handler
func Listen(ciph cipher.AEAD, address string, maxDatagramSize int, handlers Handlers) (conn *net.UDPConn, err error) {
	var addr *net.UDPAddr
	if addr, err = net.ResolveUDPAddr("udp", address); log.Check(err) {
		return
	} else if conn, err = net.ListenUDP("udp", addr); log.Check(err) {
		return
	} else if conn == nil {
		return nil, errors.New("unable to start connection ")
	}
	log.DEBUG("starting listener on", conn.LocalAddr(), "->", conn.RemoteAddr())
	if err = conn.SetReadBuffer(maxDatagramSize); log.Check(err) {
		// not a critical error but should not happen
	}
	go handlerFunc(ciph, address, conn, handlers, maxDatagramSize)
	return
}

// NewBroadcastChannel returns a broadcaster and listener with a given handler on a multicast
// address and specified port, this is just a convenience to eliminate boilerplate control
// channels and shorten the syntax, and ensure that the listener is ready before the caller
// can start sending.
//
// It is the responsibility of the caller to ignore their own messages, using magic bytes would be
// the usual way.
func NewBroadcastChannel(key string, port int, maxDatagramSize int, handlers Handlers) (channel *Channel, err error) {
	channel = &Channel{MaxDatagramSize: maxDatagramSize,}
	if key != "" {
		if channel.ciph, err = gcm.GetCipher(key); log.Check(err) {
		}
	}
	if channel.Receiver, err = ListenBroadcast(channel.ciph, port, maxDatagramSize, handlers); log.Check(err) {
	} else if channel.Sender, err = NewBroadcaster(port, maxDatagramSize); log.Check(err) {
	}
	return
}

// NewBroadcaster creates a new UDP multicast connection on which to broadcast
func NewBroadcaster(port int, maxDatagramSize int) (conn *net.UDPConn, err error) {
	address := net.JoinHostPort(UDPMulticastAddress, fmt.Sprint(port))
	if conn, err = NewSender(address, maxDatagramSize); log.Check(err) {
	}
	return
}

// ListenBroadcast binds to the UDP Address and port given and writes packets received
// from that Address to a buffer which is passed to a handler
func ListenBroadcast(ciph cipher.AEAD, port int, maxDatagramSize int, handlers Handlers) (conn *net.UDPConn, err error) {
	address := net.JoinHostPort(UDPMulticastAddress, fmt.Sprint(port))
	var addr *net.UDPAddr
	// Parse the string Address
	if addr, err = net.ResolveUDPAddr("udp", address); log.Check(err) {
		return
		// Open up a connection
	} else if conn, err = net.ListenMulticastUDP("udp", nil, addr); log.Check(err) {
		return
	} else if conn == nil {
		return nil, errors.New("unable to start connection ")
	}
	log.DEBUG("starting broadcast listener", address)
	if err = conn.SetReadBuffer(maxDatagramSize); log.Check(err) {
	}
	go handlerFunc(ciph, address, conn, handlers, maxDatagramSize)
	return
}

func handlerFunc(ciph cipher.AEAD, address string, conn *net.UDPConn, handlers Handlers, maxDatagramSize int) {
	buffer := make([]byte, maxDatagramSize)
	// Loop forever reading from the socket until it is closed
out:
	for {
		if numBytes, src, err := conn.ReadFromUDP(buffer); err != nil {
			switch handleNetworkError(address, err) {
			case closed:
				break out
			case other:
				continue
			case success:
			}
			// Filter messages by magic, if there is no match in the map the packet is ignored
		} else if numBytes > 4 {
			if handler, ok := handlers[string(buffer[:4])]; ok {
				msg := buffer[4:numBytes]
				if ciph != nil {
					if msg, err = decryptMessage(ciph, msg); log.Check(err) {
						continue
					}
				}
				if err = handler(src, address, len(msg), msg); log.Check(err) {
				}
			} else {
				log.DEBUG("ignoring irrelevant message")
				continue
			}
		} else {
			log.DEBUG("short message")
		}
	}
}

const (
	success int = iota // this is implicit zero of an int but starts the iota
	closed
	other
)

func handleNetworkError(address string, err error) (result int) {
	if len(strings.Split(err.Error(), "use of closed network connection")) >= 2 {
		log.DEBUG("connection closed", address)
		result = closed
	} else {
		log.ERRORF("ReadFromUDP failed: '%s'", err)
		result = other
	}
	return
}

func getNonce(ciph cipher.AEAD) (nonce []byte, err error) {
	// get a nonce for the packet, it is both message ID and salt
	nonce = make([]byte, ciph.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); log.Check(err) {
	}
	return
}

func decryptMessage(ciph cipher.AEAD, data []byte) (msg []byte, err error) {
	nonceSize := ciph.NonceSize()
	msg, err = ciph.Open(nil, data[:nonceSize], data[nonceSize:], nil)
	return
}

func encryptMessage(magic []byte, ciph cipher.AEAD, data []byte) (msg []byte, err error) {
	if ciph != nil {
		var nonce []byte
		nonce, err = getNonce(ciph)
		msg = append(append(magic, nonce...), ciph.Seal(nil, nonce, data, nil)...)
		return
	} else {
		return append(magic, data...), err
	}
}
