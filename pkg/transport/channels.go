package transport

import (
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"runtime"
	"strings"
	"time"
	
	"github.com/p9c/pod/pkg/fec"
	"github.com/p9c/pod/pkg/gcm"
	"github.com/p9c/pod/pkg/log"
)

const (
	UDPMulticastAddress     = "224.0.0.1"
	success             int = iota // this is implicit zero of an int but starts the iota
	closed
	other
)

type (
	MsgBuffer struct {
		Buffers [][]byte
		First   time.Time
		Decoded bool
		Source  net.Addr
	}
	// HandlerFunc is a function that is used to process a received message
	HandlerFunc func(ctx interface{}, src *net.UDPAddr, dst string, b []byte) (err error)
	Handlers    map[string]HandlerFunc
	Channel     struct {
		Creator         string
		context         interface{}
		buffers         map[string]*MsgBuffer
		Sender          *net.UDPConn
		Receiver        *net.UDPConn
		MaxDatagramSize int
		ciph            cipher.AEAD
		lastSent        *time.Time
		firstSender     *string
	}
)

// SetDestination changes the address the outbound connection of a channel directs to
func (c *Channel) SetDestination(dst string) (err error) {
	log.DEBUG("sending to", dst)
	if c.Sender, err = NewSender(dst, c.MaxDatagramSize); log.Check(err) {
	}
	return
}

// Send fires off some data through the configured channel's outbound.
func (c *Channel) Send(magic []byte, nonce []byte, data []byte) (n int, err error) {
	if len(data) == 0 {
		err = errors.New("not sending empty packet")
		log.ERROR(err)
		return
	}
	var msg []byte
	if msg, err = EncryptMessage(c.Creator, c.ciph, magic, nonce, data); log.Check(err) {
	}
	n, err = c.Sender.Write(msg)
	// log.DEBUG(msg)
	return
}

// SendMany sends a BufIter of shards as produced by GetShards
func (c *Channel) SendMany(magic []byte, b [][]byte) (err error) {
	if nonce, err := GetNonce(c.ciph); log.Check(err) {
	} else {
		for i := 0; i < len(b); i++ {
			// log.DEBUG(i)
			if _, err = c.Send(magic, nonce, b[i]); log.Check(err) {
				// debug.PrintStack()
			}
		}
		log.DEBUG(c.Creator, "sent packets", string(magic), hex.EncodeToString(nonce), c.Sender.LocalAddr(), c.Sender.RemoteAddr())
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
func GetShards(data []byte) (shards [][]byte) {
	var err error
	if shards, err = fec.Encode(data); log.Check(err) {
	}
	return
}

// NewUnicastChannel sets up a listener and sender for a specified destination
func NewUnicastChannel(creator string, ctx interface{}, key, sender, receiver string, maxDatagramSize int,
	handlers Handlers) (channel *Channel, err error) {
	channel = &Channel{Creator: creator, MaxDatagramSize: maxDatagramSize, buffers: make(map[string]*MsgBuffer),
		context: ctx}
	var magics []string
	for i := range handlers {
		magics = append(magics, i)
	}
	// log.DEBUG("magics", magics, PrevCallers())
	if key != "" {
		if channel.ciph, err = gcm.GetCipher(key); log.Check(err) {
		}
	}
	ready := make(chan struct{})
	go func() {
		channel.Receiver, err = Listen(receiver, channel, maxDatagramSize, handlers)
		ready <- struct{}{}
	}()
	<-ready
	channel.Sender, err = NewSender(sender, maxDatagramSize)
	if err != nil {
		log.ERROR(err)
	}
	log.WARN("starting unicast channel:", channel.Creator, sender, receiver, magics)
	return
}

// NewSender creates a new UDP connection to a specified address
func NewSender(address string, maxDatagramSize int) (conn *net.UDPConn, err error) {
	var addr *net.UDPAddr
	if addr, err = net.ResolveUDPAddr("udp4", address); log.Check(err) {
		return
	} else if conn, err = net.DialUDP("udp4", nil, addr); log.Check(err) {
		return
	}
	log.DEBUG("started new sender on", conn.LocalAddr(), "->", conn.RemoteAddr())
	if err = conn.SetWriteBuffer(maxDatagramSize); log.Check(err) {
	}
	return
}

// Listen binds to the UDP Address and port given and writes packets received
// from that Address to a buffer which is passed to a handler
func Listen(address string, channel *Channel, maxDatagramSize int, handlers Handlers) (conn *net.UDPConn, err error) {
	var addr *net.UDPAddr
	if addr, err = net.ResolveUDPAddr("udp4", address); log.Check(err) {
		return
	} else if conn, err = net.ListenUDP("udp4", addr); log.Check(err) {
		return
	} else if conn == nil {
		return nil, errors.New("unable to start connection ")
	}
	log.DEBUG("starting listener on", conn.LocalAddr(), "->", conn.RemoteAddr())
	if err = conn.SetReadBuffer(maxDatagramSize); log.Check(err) {
		// not a critical error but should not happen
	}
	go Handle(address, channel, handlers, maxDatagramSize)
	return
}

// NewBroadcastChannel returns a broadcaster and listener with a given handler on a multicast
// address and specified port. The handlers define the messages that will be processed and
// any other messages are ignored
func NewBroadcastChannel(creator string, ctx interface{}, key string, port int, maxDatagramSize int, handlers Handlers) (channel *Channel, err error) {
	channel = &Channel{Creator: creator, MaxDatagramSize: maxDatagramSize, buffers: make(map[string]*MsgBuffer),
		context: ctx}
	if key != "" {
		if channel.ciph, err = gcm.GetCipher(key); log.Check(err) {
		}
	}
	if channel.Receiver, err = ListenBroadcast(port, channel, maxDatagramSize, handlers); log.Check(err) {
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
func ListenBroadcast(port int, channel *Channel, maxDatagramSize int, handlers Handlers) (conn *net.UDPConn, err error) {
	address := net.JoinHostPort(UDPMulticastAddress, fmt.Sprint(port))
	var addr *net.UDPAddr
	// Parse the string Address
	if addr, err = net.ResolveUDPAddr("udp4", address); log.Check(err) {
		return
		// Open up a connection
	} else if conn, err = net.ListenMulticastUDP("udp4", nil, addr); log.Check(err) {
		return
	} else if conn == nil {
		return nil, errors.New("unable to start connection ")
	}
	var magics []string
	for i := range handlers {
		magics = append(magics, i)
	}
	// log.DEBUG("magics", magics, PrevCallers())
	log.INFO("starting broadcast listener", channel.Creator, address, magics)
	if err = conn.SetReadBuffer(maxDatagramSize); log.Check(err) {
	}
	channel.Receiver = conn
	go Handle(address, channel, handlers, maxDatagramSize)
	return
}

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

// Handle listens for messages, decodes them, aggregates them, recovers the data from the
// reed solomon fec shards received and invokes the handler provided matching the magic
// on the complete received messages
func Handle(address string, channel *Channel, handlers Handlers, maxDatagramSize int) {
	buffer := make([]byte, maxDatagramSize)
	log.WARN("starting handler for", channel.Creator, "listener")
	// Loop forever reading from the socket until it is closed
out:
	for {
		if numBytes, src, err := channel.Receiver.ReadFromUDP(buffer); err != nil {
			switch handleNetworkError(address, err) {
			case closed:
				break out
			case other:
				continue
			case success:
			}
			// Filter messages by magic, if there is no match in the map the packet is ignored
		} else if numBytes > 4 {
			magic := string(buffer[:4])
			// log.DEBUG("magic", magic)
			// for i := range handlers {
			// 	log.INFO(i)
			// }
			if handler, ok := handlers[magic]; ok {
				// if caller needs to know the liveness status of the
				// controller it is working on, the code below
				if channel.lastSent != nil && channel.firstSender != nil {
					*channel.lastSent = time.Now()
				}
				msg := buffer[:numBytes]
				nL := channel.ciph.NonceSize()
				nonceBytes := msg[4 : 4+nL]
				nonce := string(nonceBytes)
				// decipher
				var shard []byte
				if shard, err = channel.ciph.Open(nil, nonceBytes, msg[4+len(nonceBytes):], nil); log.Check(err) {
					continue
				}
				if bn, ok := channel.buffers[nonce]; ok {
					// log.DEBUGF("%s adding shard to %"+fmt.Sprint(nL*2)+"x", channel.Creator, nonceBytes)
					if !bn.Decoded {
						// log.DEBUG(PrevCallers())
						bn.Buffers = append(bn.Buffers, shard)
						if len(bn.Buffers) >= 3 {
							// try to decode it
							var cipherText []byte
							cipherText, err = fec.Decode(bn.Buffers)
							if err != nil {
								log.ERROR(err)
								continue
							}
							log.DEBUG(hex.EncodeToString(cipherText))
							bn.Decoded = true
							// if channel.ciph != nil {
							// 	if msg, err = DecryptMessage(channel.Creator, channel.ciph, cipherText); log.Check(err) {
							// 		log.WARN(PrevCallers())
							// 		continue
							// 	}
							if err = handler(channel.context, src, address, cipherText); log.Check(err) {
								// err = handler(channel.context, src, channel.Sender.RemoteAddr().String(), cipherText)
								// if err != nil {
								// 	log.ERROR(err)
								// 	continue
								// }
							}
							// }
						}
					} else {
						for i := range channel.buffers {
							if i != nonce {
								// superseded messages can be deleted from the
								// buffers, we don't add more data for the already
								// decoded.
								delete(channel.buffers, i)
							}
						}
					}
				} else {
					log.DEBUGF("%s adding %s %x to buffers", channel.Creator, magic, nonceBytes)
					channel.buffers[nonce] = &MsgBuffer{[][]byte{},
						time.Now(), false, src}
					channel.buffers[nonce].Buffers = append(channel.buffers[nonce].
						Buffers, shard)
				}
			} else {
				// log.DEBUGF("ignoring irrelevant message %s\n%s\n%s", magic, PrevCallers(), string(debug.Stack()))
				// log.DEBUG("ignoring irrelevant message", magic, channel.Creator)
				continue
			}
		} else {
			log.DEBUG(channel.Creator, "short message")
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
