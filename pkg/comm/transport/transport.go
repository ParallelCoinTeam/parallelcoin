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

	"github.com/p9c/pod/pkg/coding/fec"
)

type HandleFunc map[string]func(ctx interface{}) func(b []byte) (err error)

// Connection is the state and working memory references for a simple
// reliable UDP lan transport, encrypted by a GCM AES cipher,
// with the simple protocol of sending out 9 packets containing encrypted FEC
// shards containing a slice of bytes.
// This protocol probably won't work well outside of a multicast lan in
// adverse conditions but it is designed for local network control systems
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
// 	maxDatagramSize int, ctx context.Context) (c *Connection, err error) {
// 	sendAddr := &net.UDPAddr{}
// 	var sendConn net.Conn
// 	listenAddr := &net.UDPAddr{}
// 	var listenConn net.PacketConn
// 	if listen != "" {
// 		config := &net.ListenConfig{Control: reusePort}
// 		listenConn, err = config.ListenPacket(context.Background(), "udp4", listen)
// 		if err != nil {
// 			Error(err)
// 		}
// 	}
// 	if send != "" {
// 		// sendAddr, err = net.ResolveUDPAddr("udp4", send)
// 		// if err != nil {
// 		// 	Error(err)
// 		// }
// 		sendConn, err = net.Dial("udp4", send)
// 		if err != nil {
// 			Error(err, sendAddr)
// 		}
// 		// L.Spew(sendConn)
// 	}
// 	var ciph cipher.AEAD
// 	if ciph, err = gcm.GetCipher(preSharedKey); Check(err) {
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

func (c *Connection) SetSendConn(ad string) (err error) {
	// c.sendAddress, err = net.ResolveUDPAddr("udp4", ad)
	// if err != nil {
	// 	Error(err)
	// }
	var sC net.Conn
	sC, err = net.Dial("udp4", ad)
	if err != nil {
		Error(err)
		return
	}
	c.SendConn = sC
	return
}

func (c *Connection) CreateShards(b, magic []byte) (shards [][]byte,
	err error) {
	magicLen := 4
	// get a nonce for the packet, it is both message ID and salt
	nonceLen := c.ciph.NonceSize()
	nonce := make([]byte, nonceLen)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		Error(err)
		return
	}
	// generate the shards
	shards, err = fec.Encode(b)
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

func send(shards [][]byte, sendConn net.Conn) (err error) {
	for i := range shards {
		_, err = sendConn.Write(shards[i])
		if err != nil {
			Error(err)
		}
	}
	return
}

func (c *Connection) Send(b, magic []byte) (err error) {
	if len(magic) != 4 {
		err = errors.New("magic must be 4 bytes long")
		Error(err)
		return
	}
	var shards [][]byte
	shards, err = c.CreateShards(b, magic)
	err = send(shards, c.SendConn)
	if err != nil {
		Error(err)
	}
	return
}

func (c *Connection) SendTo(addr *net.UDPAddr, b, magic []byte) (err error) {
	if len(magic) != 4 {
		err = errors.New("magic must be 4 bytes long")
		Error(err)
		return
	}
	sendConn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		Error(err)
		return
	}
	shards, err := c.CreateShards(b, magic)
	err = send(shards, sendConn)
	if err != nil {
		Error(err)
	}
	return
}

func (c *Connection) SendShards(shards [][]byte) (err error) {
	err = send(shards, c.SendConn)
	if err != nil {
		Error(err)
	}
	return
}

func (c *Connection) SendShardsTo(shards [][]byte, addr *net.UDPAddr) (err error) {
	sendConn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		Error(err)
		return
	}
	err = send(shards, sendConn)
	if err != nil {
		Error(err)
	}
	return
}

func (c *Connection) Listen(handlers HandleFunc, ifc interface{},
	lastSent *time.Time, firstSender *string) (err error) {
	Trace("setting read buffer")
	buffer := make([]byte, c.maxDatagramSize)
	go func() {
		Trace("starting connection handler")
	out:
		// read from socket until context is cancelled
		for {
			n, src, err := (*c.listenConn).ReadFrom(buffer)
			buf := buffer[:n]
			if err != nil {
				// Error("ReadFromUDP failed:", err)
				continue
			}
			magic := string(buf[:4])
			if _, ok := handlers[magic]; ok {
				// if caller needs to know the liveness status of the
				// controller it is working on, the code below
				if lastSent != nil && firstSender != nil {
					*lastSent = time.Now()
				}
				nonceBytes := buf[4:16]
				nonce := string(nonceBytes)
				// decipher
				shard, err := c.ciph.Open(nil, nonceBytes,
					buf[16:], nil)
				if err != nil {
					// Error(err)
					// corrupted or irrelevant message
					continue
				}
				if bn, ok := c.buffers[nonce]; ok {
					if !bn.Decoded {
						bn.Buffers = append(bn.Buffers, shard)
						if len(bn.Buffers) >= 3 {
							// try to decode it
							var cipherText []byte
							cipherText, err = fec.Decode(bn.Buffers)
							if err != nil {
								Error(err)
								continue
							}
							bn.Decoded = true
							err = handlers[magic](ifc)(cipherText)
							if err != nil {
								Error(err)
								continue
							}
						}
					} else {
						for i := range c.buffers {
							if i != nonce {
								// superseded messages can be deleted from the
								// buffers,
								// we don't add more data for the already
								// decoded.
								// Trace("deleting superseded buffer",
								// 	hex.EncodeToString([]byte(i)))
								delete(c.buffers, i)
							}
						}
					}
				} else {
					// Trace("new message arriving",
					// 	hex.EncodeToString([]byte(nonce)))
					c.buffers[nonce] = &MsgBuffer{[][]byte{},
						time.Now(), false, src}
					c.buffers[nonce].Buffers = append(c.buffers[nonce].
						Buffers, shard)
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
// 	sendHost, sendPort, err := net.SplitHostPort(address)
// 	if err != nil {
// 		Error(err)
// 		return
// 	}
// 	sendPortI, err := strconv.ParseInt(sendPort, 10, 64)
// 	if err != nil {
// 		Error(err)
// 		return
// 	}
// 	sendAddr = &net.UDPAddr{IP: net.ParseIP(sendHost),
// 		Port: int(sendPortI)}
// 	// Debug("multicast", Address)
// 	// L.Spew(sendAddr)
// 	return
// }
