package transport

import (
	"crypto/cipher"
	"errors"
	"fmt"
	"net"
	"runtime"
	"strings"
	"time"

	"github.com/p9c/log"

	"github.com/p9c/qu"

	"github.com/p9c/matrjoska/pkg/fec"
	"github.com/p9c/matrjoska/pkg/gcm"
	"github.com/p9c/matrjoska/pkg/multicast"
)

const (
	UDPMulticastAddress     = "224.0.0.1"
	success             int = iota // this is implicit zero of an int but starts the iota
	closed
	other
	DefaultPort = 11049
)

var DefaultIP = net.IPv4(224, 0, 0, 1)
var MulticastAddress = &net.UDPAddr{IP: DefaultIP, Port: DefaultPort}

type (
	MsgBuffer struct {
		Buffers [][]byte
		First   time.Time
		Decoded bool
		Source  net.Addr
	}
	// HandlerFunc is a function that is used to process a received message
	HandlerFunc func(
		ctx interface{}, src net.Addr, dst string, b []byte,
	) (e error)
	Handlers    map[string]HandlerFunc
	Channel     struct {
		buffers         map[string]*MsgBuffer
		Ready           qu.C
		context         interface{}
		Creator         string
		firstSender     *string
		lastSent        *time.Time
		MaxDatagramSize int
		receiveCiph     cipher.AEAD
		Receiver        *net.UDPConn
		sendCiph        cipher.AEAD
		Sender          *net.UDPConn
	}
)

// SetDestination changes the address the outbound connection of a multicast
// directs to
func (c *Channel) SetDestination(dst string) (e error) {
	D.Ln("sending to", dst)
	if c.Sender, e = NewSender(dst, c.MaxDatagramSize); E.Chk(e) {
	}
	return
}

// Send fires off some data through the configured multicast's outbound.
func (c *Channel) Send(magic []byte, nonce []byte, data []byte) (
	n int, e error,
) {
	if len(data) == 0 {
		e = errors.New("not sending empty packet")
		E.Ln(e)
		return
	}
	var msg []byte
	if msg, e = EncryptMessage(c.Creator, c.sendCiph, magic, nonce, data); E.Chk(e) {
	}
	n, e = c.Sender.Write(msg)
	// D.Ln(msg)
	return
}

// SendMany sends a BufIter of shards as produced by GetShards
func (c *Channel) SendMany(magic []byte, b [][]byte) (e error) {
	D.Ln("magic", string(magic), log.Caller("sending from", 1))
	var nonce []byte
	if nonce, e = GetNonce(c.sendCiph); E.Chk(e) {
	} else {
		for i := 0; i < len(b); i++ {
			// D.Ln(i)
			// D.Ln("segment length", len(b[i]))
			if _, e = c.Send(magic, nonce, b[i]); E.Chk(e) {
				// debug.PrintStack()
			}
		}
		// T.Ln(c.Creator, "sent packets", string(magic), hex.EncodeToString(nonce), c.Sender.LocalAddr(), c.Sender.RemoteAddr())
	}
	return
}

// Close the multicast
func (c *Channel) Close() (e error) {
	// if e = c.Sender.Close(); E.Chk(e) {
	// }
	// if e = c.Receiver.Close(); E.Chk(e) {
	// }
	return
}

// GetShards returns a buffer iterator to feed to Channel.SendMany containing
// fec encoded shards built from the provided buffer
func GetShards(data []byte) (shards [][]byte) {
	var e error
	if shards, e = fec.Encode(data); E.Chk(e) {
	}
	return
}

// NewUnicastChannel sets up a listener and sender for a specified destination
func NewUnicastChannel(
	creator string, ctx interface{}, key []byte, sender, receiver string,
	maxDatagramSize int,
	handlers Handlers, quit qu.C,
) (channel *Channel, e error) {
	channel = &Channel{
		Creator:         creator,
		MaxDatagramSize: maxDatagramSize,
		buffers:         make(map[string]*MsgBuffer),
		context:         ctx,
	}
	var magics []string
	bytes := make([]byte, len(key))
	copy(bytes, key)
	for i := range handlers {
		magics = append(magics, i)
	}
	if channel.sendCiph, e = gcm.GetCipher(bytes); E.Chk(e) {
	}
	copy(bytes, key)
	if channel.receiveCiph, e = gcm.GetCipher(bytes); E.Chk(e) {
	}
	for i := range bytes {
		bytes[i] = 0
		key[i] = 0
	}
	if channel.Receiver, e = Listen(receiver, channel, maxDatagramSize, handlers, quit); E.Chk(e) {
	}
	if channel.Sender, e = NewSender(sender, maxDatagramSize); E.Chk(e) {
	}
	D.Ln("starting unicast multicast:", channel.Creator, sender, receiver, magics)
	return
}

// NewSender creates a new UDP connection to a specified address
func NewSender(address string, maxDatagramSize int) (
	conn *net.UDPConn, e error,
) {
	var addr *net.UDPAddr
	if addr, e = net.ResolveUDPAddr("udp4", address); E.Chk(e) {
		return
	} else if conn, e = net.DialUDP("udp4", nil, addr); E.Chk(e) {
		// debug.PrintStack()
		return
	}
	D.Ln("started new sender on", conn.LocalAddr(), "->", conn.RemoteAddr())
	if e = conn.SetWriteBuffer(maxDatagramSize); E.Chk(e) {
	}
	return
}

// Listen binds to the UDP Address and port given and writes packets received
// from that Address to a buffer which is passed to a handler
func Listen(
	address string, channel *Channel, maxDatagramSize int, handlers Handlers,
	quit qu.C,
) (conn *net.UDPConn, e error) {
	var addr *net.UDPAddr
	if addr, e = net.ResolveUDPAddr("udp4", address); E.Chk(e) {
		return
	} else if conn, e = net.ListenUDP("udp4", addr); E.Chk(e) {
		return
	} else if conn == nil {
		return nil, errors.New("unable to start connection ")
	}
	D.Ln("starting listener on", conn.LocalAddr(), "->", conn.RemoteAddr())
	if e = conn.SetReadBuffer(maxDatagramSize); E.Chk(e) {
		// not a critical error but should not happen
	}
	go Handle(address, channel, handlers, maxDatagramSize, quit)
	return
}

// NewBroadcastChannel returns a broadcaster and listener with a given handler
// on a multicast address and specified port. The handlers define the messages
// that will be processed and any other messages are ignored
func NewBroadcastChannel(
	creator string, ctx interface{}, key []byte, port int, maxDatagramSize int,
	handlers Handlers,
	quit qu.C,
) (channel *Channel, e error) {
	channel = &Channel{
		Creator:         creator,
		MaxDatagramSize: maxDatagramSize,
		buffers:         make(map[string]*MsgBuffer),
		context:         ctx,
		Ready:           qu.T(),
	}
	bytes := make([]byte, len(key))
	copy(bytes, key)
	if channel.sendCiph, e = gcm.GetCipher(bytes); E.Chk(e) {
		panic(e)
	}
	copy(bytes, key)
	if channel.receiveCiph, e = gcm.GetCipher(bytes); E.Chk(e) {
		panic(e)
	}
	for i := range bytes {
		key[i] = 0
		bytes[i] = 0
	}
	if channel.Receiver, e = ListenBroadcast(port, channel, maxDatagramSize, handlers, quit); E.Chk(e) {
	}
	if channel.Sender, e = NewBroadcaster(port, maxDatagramSize); E.Chk(e) {
	}
	channel.Ready.Q()
	return
}

// NewBroadcaster creates a new UDP multicast connection on which to broadcast
func NewBroadcaster(port int, maxDatagramSize int) (
	conn *net.UDPConn, e error,
) {
	address := net.JoinHostPort(UDPMulticastAddress, fmt.Sprint(port))
	if conn, e = NewSender(address, maxDatagramSize); E.Chk(e) {
	}
	return
}

// ListenBroadcast binds to the UDP Address and port given and writes packets
// received from that Address to a buffer which is passed to a handler
func ListenBroadcast(
	port int,
	channel *Channel,
	maxDatagramSize int,
	handlers Handlers,
	quit qu.C,
) (conn *net.UDPConn, e error) {
	if conn, e = multicast.Conn(port); E.Chk(e) {
		return
	}
	address := conn.LocalAddr().String()
	var magics []string
	for i := range handlers {
		magics = append(magics, i)
	}
	// D.S(handlers)
	// D.Ln("magics", magics, PrevCallers())
	D.Ln("starting broadcast listener", channel.Creator, address, magics)
	if e = conn.SetReadBuffer(maxDatagramSize); E.Chk(e) {
	}
	channel.Receiver = conn
	go Handle(address, channel, handlers, maxDatagramSize, quit)
	return
}

func handleNetworkError(address string, e error) (result int) {
	if len(strings.Split(e.Error(), "use of closed network connection")) >= 2 {
		D.Ln("connection closed", address)
		result = closed
	} else {
		E.F("ReadFromUDP failed: '%s'", e)
		result = other
	}
	return
}

// Handle listens for messages, decodes them, aggregates them, recovers the data
// from the reed solomon fec shards received and invokes the handler provided
// matching the magic on the complete received messages
func Handle(
	address string, channel *Channel,
	handlers Handlers, maxDatagramSize int, quit qu.C,
) {
	buffer := make([]byte, maxDatagramSize)
	T.Ln("starting handler for", channel.Creator, "listener")
	// Loop forever reading from the socket until it is closed
	// seenNonce := ""
	var e error
	var numBytes int
	var src net.Addr
	// var seenNonce string
	<-channel.Ready
out:
	for {
		select {
		case <-quit.Wait():
			break out
		default:
		}
		if numBytes, src, e = channel.Receiver.ReadFromUDP(buffer); E.Chk(e) {
			switch handleNetworkError(address, e) {
			case closed:
				break out
			case other:
				continue
			case success:
			}
		}
		// Filter messages by magic, if there is no match in the map the packet is
		// ignored
		magic := string(buffer[:4])
		if handler, ok := handlers[magic]; ok {
			if channel.lastSent != nil && channel.firstSender != nil {
				*channel.lastSent = time.Now()
			}
			msg := buffer[:numBytes]
			nL := channel.receiveCiph.NonceSize()
			nonceBytes := msg[4 : 4+nL]
			nonce := string(nonceBytes)
			var shard []byte
			if shard, e = channel.receiveCiph.Open(nil, nonceBytes, msg[4+len(nonceBytes):], nil); e != nil {
				continue
			}
			// D.Ln("read", numBytes, "from", src, e, hex.EncodeToString(msg))
			if bn, ok := channel.buffers[nonce]; ok {
				if !bn.Decoded {
					bn.Buffers = append(bn.Buffers, shard)
					if len(bn.Buffers) >= 3 {
						// try to decode it
						var cipherText []byte
						if cipherText, e = fec.Decode(bn.Buffers); E.Chk(e) {
							continue
						}
						// D.F("received packet with magic %s from %s len %d bytes", magic, src.String(), len(cipherText))
						bn.Decoded = true
						if e = handler(channel.context, src, address, cipherText); E.Chk(e) {
							continue
						}
						// src = nil
						// buffer = buffer[:0]
					}
				} else {
					for i := range channel.buffers {
						if i != nonce && channel.buffers[i].Decoded {
							// superseded messages can be deleted from the buffers, we don't add more data
							// for the already decoded. todo: this will be changed to track stats for the
							// puncture rate and redundancy scaling
							delete(channel.buffers, i)
						}
					}
				}
			} else {
				channel.buffers[nonce] = &MsgBuffer{
					[][]byte{},
					time.Now(), false, src,
				}
				channel.buffers[nonce].Buffers = append(
					channel.buffers[nonce].
						Buffers, shard,
				)
			}
		}
	}
}

func PrevCallers() (out string) {
	for i := 0; i < 10; i++ {
		_, loc, iline, _ := runtime.Caller(i)
		out += fmt.Sprintf("%s:%d \n", loc, iline)
	}
	return
}
