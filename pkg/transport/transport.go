package transport

import (
	"context"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"net"
	"sync"
	"time"
	
	"github.com/p9c/matrjoska/pkg/fec"
)

// HandleFunc is a map of handlers for working on received, decoded packets
type HandleFunc map[string]func(ctx interface{}) func(b []byte) (e error)

// Connection is the state and working memory references for a simple reliable UDP lan transport, encrypted by a GCM AES
// cipher, with the simple protocol of sending out 9 packets containing encrypted FEC shards containing a slice of
// bytes.
//
// This protocol probably won't work well outside of a multicast lan in adverse conditions but it is designed for local
// network control systems todo: it is if the updated fec segmenting code is put in
type Connection struct {
	maxDatagramSize int
	buffers         map[string]*MsgBuffer
	sendAddress     *net.UDPAddr
	SendConn        net.Conn
	listenAddress   *net.UDPAddr
	listenConn      *net.PacketConn
	ciph            cipher.AEAD
	ctx             context.Context
	mx              *sync.Mutex
}

//
// // NewConnection creates a new connection with a defined default send
// // connection and listener and pre shared key password for encryption on the
// // local network
// func NewConnection(send, listen, preSharedKey string,
// 	maxDatagramSize int, ctx context.Context) (c *Connection, e error) {
// 	sendAddr := &net.UDPAddr{}
// 	var sendConn net.Conn
// 	listenAddr := &net.UDPAddr{}
// 	var listenConn net.PacketConn
// 	if listen != "" {
// 		config := &net.ListenConfig{Control: reusePort}
// 		listenConn, e = config.ListenPacket(context.Background(), "udp4", listen)
// 		if e != nil  {
// 			E.Ln(e)
// 		}
// 	}
// 	if send != "" {
// 		// sendAddr, e = net.ResolveUDPAddr("udp4", send)
// 		// if e != nil  {
// 		// 	E.Ln(e)
// 		// }
// 		sendConn, e = net.Dial("udp4", send)
// 		if e != nil  {
// 			Error(err, sendAddr)
// 		}
// 		// L.Spew(sendConn)
// 	}
// 	var ciph cipher.AEAD
// 	if ciph, e = gcm.GetCipher(preSharedKey); E.Chk(e) {
// 	}
// 	return &Connection{
// 		maxDatagramSize: maxDatagramSize,
// 		buffers:         make(map[string]*MsgBuffer),
// 		sendAddress:     sendAddr,
// 		SendConn:        sendConn,
// 		listenAddress:   listenAddr,
// 		listenConn:      &listenConn,
// 		ciph:            ciph, // gcm.GetCipher(*cx.Config.MinerPass),
// 		ctx:             ctx,
// 		mx:              &sync.Mutex{},
// 	}, err
// }

// SetSendConn sets up an outbound connection
func (c *Connection) SetSendConn(ad string) (e error) {
	// c.sendAddress, e = net.ResolveUDPAddr("udp4", ad)
	// if e != nil  {
	// 		// }
	var sC net.Conn
	if sC, e = net.Dial("udp4", ad); !E.Chk(e) {
		c.SendConn = sC
	}
	return
}

// CreateShards takes a slice of bites and generates 3
func (c *Connection) CreateShards(b, magic []byte) (
	shards [][]byte,
	e error,
) {
	magicLen := 4
	// get a nonce for the packet, it is both message ID and salt
	nonceLen := c.ciph.NonceSize()
	nonce := make([]byte, nonceLen)
	if _, e = io.ReadFull(rand.Reader, nonce); E.Chk(e) {
		return
	}
	// generate the shards
	if shards, e = fec.Encode(b); E.Chk(e) {
	}
	for i := range shards {
		encryptedShard := c.ciph.Seal(nil, nonce, shards[i], nil)
		shardLen := len(encryptedShard)
		// assemble the packet: magic, nonce, and encrypted shard
		outBytes := make([]byte, shardLen+magicLen+nonceLen)
		copy(outBytes, magic[:magicLen])
		copy(outBytes[magicLen:], nonce)
		copy(outBytes[magicLen+nonceLen:], encryptedShard)
		shards[i] = outBytes
	}
	return
}

func send(shards [][]byte, sendConn net.Conn) (e error) {
	for i := range shards {
		if _, e = sendConn.Write(shards[i]); E.Chk(e) {
		}
	}
	return
}

func (c *Connection) Send(b, magic []byte) (e error) {
	if len(magic) != 4 {
		e = errors.New("magic must be 4 bytes long")
		return
	}
	var shards [][]byte
	shards, e = c.CreateShards(b, magic)
	if e = send(shards, c.SendConn); E.Chk(e) {
	}
	return
}

func (c *Connection) SendTo(addr *net.UDPAddr, b, magic []byte) (e error) {
	if len(magic) != 4 {
		if e = errors.New("magic must be 4 bytes long"); E.Chk(e) {
			return
		}
	}
	var sendConn *net.UDPConn
	if sendConn, e = net.DialUDP("udp", nil, addr); E.Chk(e) {
		return
	}
	var shards [][]byte
	if shards, e = c.CreateShards(b, magic); E.Chk(e) {
	}
	if e = send(shards, sendConn); E.Chk(e) {
	}
	return
}

func (c *Connection) SendShards(shards [][]byte) (e error) {
	if e = send(shards, c.SendConn); E.Chk(e) {
	}
	return
}

func (c *Connection) SendShardsTo(shards [][]byte, addr *net.UDPAddr) (e error) {
	var sendConn *net.UDPConn
	if sendConn, e = net.DialUDP("udp", nil, addr); !E.Chk(e) {
		if e = send(shards, sendConn); E.Chk(e) {
		}
	}
	return
}

// Listen runs a goroutine that collects and attempts to decode the FEC shards
// once it has enough intact pieces
func (c *Connection) Listen(handlers HandleFunc, ifc interface{}, lastSent *time.Time, firstSender *string,) (e error) {
	F.Ln("setting read buffer")
	buffer := make([]byte, c.maxDatagramSize)
	go func() {
		F.Ln("starting connection handler")
	out:
		// read from socket until context is cancelled
		for {
			var src net.Addr
			var n int
			n, src, e = (*c.listenConn).ReadFrom(buffer)
			buf := buffer[:n]
			if E.Chk(e) {
				// Error("ReadFromUDP failed:", e)
				continue
			}
			magic := string(buf[:4])
			if _, ok := handlers[magic]; ok {
				// if caller needs to know the liveness status of the controller it is working on, the code below
				if lastSent != nil && firstSender != nil {
					*lastSent = time.Now()
				}
				nonceBytes := buf[4:16]
				nonce := string(nonceBytes)
				// decipher
				var shard []byte
				if shard, e = c.ciph.Open(nil, nonceBytes, buf[16:], nil); E.Chk(e) {
					// corrupted or irrelevant message
					continue
				}
				var bn *MsgBuffer
				if bn, ok = c.buffers[nonce]; ok {
					if !bn.Decoded {
						bn.Buffers = append(bn.Buffers, shard)
						if len(bn.Buffers) >= 3 {
							// try to decode it
							var cipherText []byte
							if cipherText, e = fec.Decode(bn.Buffers); E.Chk(e) {
								continue
							}
							bn.Decoded = true
							if e = handlers[magic](ifc)(cipherText); E.Chk(e) {
								continue
							}
						}
					} else {
						for i := range c.buffers {
							if i != nonce {
								// superseded messages can be deleted from the buffers, we don't add more data
								// for the already decoded.
								// F.Ln("deleting superseded buffer", hex.EncodeToString([]byte(i)))
								delete(c.buffers, i)
							}
						}
					}
				} else {
					// F.Ln("new message arriving",
					// 	hex.EncodeToString([]byte(nonce)))
					c.buffers[nonce] = &MsgBuffer{
						[][]byte{},
						time.Now(), false, src,
					}
					c.buffers[nonce].Buffers = append(
						c.buffers[nonce].
							Buffers, shard,
					)
				}
			}
			select {
			case <-c.ctx.Done():
				break out
			default:
			}
		}
	}()
	return
}

//
// func GetUDPAddr(address string) (sendAddr *net.UDPAddr) {
// 	sendHost, sendPort, e := net.SplitHostPort(address)
// 	if e != nil  {
// 		// 		return
// 	}
// 	sendPortI, e := strconv.ParseInt(sendPort, 10, 64)
// 	if e != nil  {
// 		// 		return
// 	}
// 	sendAddr = &net.UDPAddr{IP: net.ParseIP(sendHost),
// 		Port: int(sendPortI)}
// 	// D.Ln("multicast", Address)
// 	// L.Spew(sendAddr)
// 	return
// }
