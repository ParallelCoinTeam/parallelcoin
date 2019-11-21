package transport

import (
	"context"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"github.com/p9c/pod/pkg/fec"
	"github.com/p9c/pod/pkg/gcm"
	"github.com/p9c/pod/pkg/log"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

type MsgBuffer struct {
	Buffers [][]byte
	First   time.Time
	Decoded bool
	Source  *net.UDPAddr
}

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
	sendConn        *net.UDPConn
	listenAddress   *net.UDPAddr
	listenConn      *net.UDPConn
	ciph            cipher.AEAD
	ctx             context.Context
	mx              *sync.Mutex
	receiveChan     chan []byte
}

// NewConnection creates a new connection with a defined default send
// connection and listener and pre shared key password for encryption on the
// local network
func NewConnection(send, listen, preSharedKey string,
	maxDatagramSize int) (c *Connection, cancel context.CancelFunc, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	sendAddr := GetUDPAddr(send)
	sendConn, err := net.DialUDP("udp", nil, sendAddr)
	if err != nil {
		log.ERROR(err)
		return
	}
	listenAddr := GetUDPAddr(listen)
	listenConn, err := net.DialUDP("udp", nil, listenAddr)
	if err != nil {
		log.ERROR(err)
		return
	}
	ciph := gcm.GetCipher(preSharedKey)
	return &Connection{
		maxDatagramSize: maxDatagramSize,
		buffers:         make(map[string]*MsgBuffer),
		sendAddress:     sendAddr,
		sendConn:        sendConn,
		listenAddress:   listenAddr,
		listenConn:      listenConn,
		ciph:            ciph, // gcm.GetCipher(*cx.Config.MinerPass),
		ctx:             ctx,
		mx:              &sync.Mutex{},
		receiveChan:     make(chan []byte),
	}, cancel, err
}

func (c *Connection) CreateShards(b, magic []byte) (shards [][]byte,
	err error) {
	magicLen := 4
	// get a nonce for the packet, it is both message ID and salt
	nonceLen := c.ciph.NonceSize()
	nonce := make([]byte, nonceLen)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.ERROR(err)
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

func send(shards [][]byte, sendConn *net.UDPConn) (err error) {
	for i := range shards {
		_, err = sendConn.Write(shards[i])
		if err != nil {
			log.ERROR(err)
		}
	}
	return
}

func (c *Connection) Send(b, magic []byte) (err error) {
	if len(magic) != 4 {
		err = errors.New("magic must be 4 bytes long")
		log.ERROR(err)
		return
	}
	var shards [][]byte
	shards, err = c.CreateShards(b, magic)
	err = send(shards, c.sendConn)
	if err != nil {
		log.ERROR(err)
	}
	return
}

func (c *Connection) SendTo(addr *net.UDPAddr, b, magic []byte) (err error) {
	if len(magic) != 4 {
		err = errors.New("magic must be 4 bytes long")
		log.ERROR(err)
		return
	}
	sendConn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.ERROR(err)
		return
	}
	shards, err := c.CreateShards(b, magic)
	err = send(shards, sendConn)
	if err != nil {
		log.ERROR(err)
	}
	return
}

func (c *Connection) SendShards(shards [][]byte) (err error) {
	err = send(shards, c.sendConn)
	if err != nil {
		log.ERROR(err)
	}
	return
}

func (c *Connection) SendShardsTo(shards [][]byte, addr *net.UDPAddr) (err error) {
	sendConn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.ERROR(err)
		return
	}
	err = send(shards, sendConn)
	if err != nil {
		log.ERROR(err)
	}
	return
}

func (c *Connection) Listen(handlers map[string]func(b []byte) (err error),
) (b []byte, err error) {
	log.TRACE("setting read buffer")
	err = c.listenConn.SetReadBuffer(c.maxDatagramSize)
	if err != nil {
		log.ERROR(err)
		return
	}
	buffer := make([]byte, c.maxDatagramSize)
	go func() {
		log.TRACE("starting connection handler")
	out:
		// read from socket until context is cancelled
		for {
			n, src, err := c.listenConn.ReadFromUDP(buffer)
			buf := buffer[:n]
			if err != nil {
				log.ERROR("ReadFromUDP failed:", err)
				continue
			}
			magic := string(buf[:4])
			if _, ok := handlers[magic]; ok {
				nonceBytes := buf[4:16]
				nonce := string(nonceBytes)
				// decipher
				shard, err := c.ciph.Open(nil, nonceBytes,
					buf[16:], nil)
				if err != nil {
					log.ERROR(err)
					continue
				}
				if bn, ok := c.buffers[nonce]; ok {
					if !bn.Decoded {
						bn.Buffers = append(bn.Buffers, shard)
						if len(bn.Buffers) >= 3 {
							// try to decode it
							var cipherText []byte
							log.SPEW(bn.Buffers)
							cipherText, err = fec.Decode(bn.Buffers)
							if err != nil {
								log.ERROR(err)
								continue
							}
							log.SPEW(cipherText)
							bn.Decoded = true
							c.receiveChan <- cipherText
						}
					} else {
						for i := range c.buffers {
							if i != nonce {
								// superseded messages can be deleted from the
								// buffers,
								// we don't add more data for the already
								// decoded.
								delete(c.buffers, i)
							}
						}
					}
				} else {
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

func GetUDPAddr(address string) (sendAddr *net.UDPAddr) {
	sendHost, sendPort, err := net.SplitHostPort(address)
	if err != nil {
		log.ERROR(err)
		return
	}
	sendPortI, err := strconv.ParseInt(sendPort, 10, 64)
	if err != nil {
		log.ERROR(err)
		return
	}
	sendAddr = &net.UDPAddr{IP: net.ParseIP(sendHost),
		Port: int(sendPortI)}
	return
}
