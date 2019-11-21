package transport

import (
	"context"
	"crypto/cipher"
	"crypto/rand"
	"github.com/p9c/pod/pkg/controller/controllerold/fec"
	"github.com/p9c/pod/pkg/gcm"
	"github.com/p9c/pod/pkg/log"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

type MsgBuffer struct {
	Buffers    [][]byte
	First      time.Time
	Decoded    bool
	Superseded bool
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
	sendAddr := getUDPAddr(send)
	sendConn, err := net.DialUDP("udp", nil, sendAddr)
	if err != nil {
		log.ERROR(err)
		return
	}
	listenAddr := getUDPAddr(listen)
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

func (c *Connection) createShards(b, magic []byte) (shards [][]byte,
	err error) {
	magicLen := len(magic)
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
		copy(outBytes, magic)
		copy(outBytes[magicLen:], nonce)
		copy(outBytes[magicLen+nonceLen:], encryptedShard)
		shards[i] = outBytes
	}
	return
}

func (c *Connection) Send(b, magic []byte) (err error) {
	var shards [][]byte
	shards, err = c.createShards(b, magic)
	for i := range shards {
		_, err = c.sendConn.Write(shards[i])
		if err != nil {
			log.ERROR(err)
		}
	}
	return
}

func (c *Connection) Listen(handlers map[string]func(b []byte) (cancel context.
CancelFunc, err error)) (b []byte, err error) {
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
			// TODO: decrypt shards when they arrive and store the cleartext
			//  in the buffers
			_ = src
			if err != nil {
				log.ERROR("ReadFromUDP failed:", err)
				continue
			}
			magic := string(buffer[:])
			if _, ok := handlers[magic]; ok {
				nonce := string(buffer[:12])
				if bn, ok := c.buffers[nonce]; ok {
					if !bn.Decoded {
						payload := buffer[16:n]
						newP := make([]byte, len(payload))
						copy(newP, payload)
						bn.Buffers = append(bn.Buffers, newP)
						if len(bn.Buffers) >= 3 {
							// try to decode it
							var cipherText []byte
							log.SPEW(bn.Buffers)
							cipherText, err = fec.Decode(bn.Buffers)
							if err != nil {
								log.ERROR(err)
								return
							}
							log.SPEW(cipherText)
							msg, err := c.ciph.Open(nil, []byte(nonce),
								cipherText, nil)
							if err != nil {
								log.ERROR(err)
								return
							}
							bn.Decoded = true
							c.receiveChan <- msg
						}
					} else {
						for i := range c.buffers {
							if i != nonce {
								// superseded blocks can be deleted from the
								// buffers,
								// we don't add more data for the already
								// decoded
								c.buffers[i].Superseded = true
							}
						}
					}
				} else {
					c.buffers[nonce] = &MsgBuffer{[][]byte{},
						time.Now(), false, false}
					payload := buffer[16:n]
					newP := make([]byte, len(payload))
					copy(newP, payload)
					c.buffers[nonce].Buffers = append(c.buffers[nonce].
						Buffers, newP)
					//log.DEBUGF("%x", payload)
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

func getUDPAddr(address string) (sendAddr *net.UDPAddr) {
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
