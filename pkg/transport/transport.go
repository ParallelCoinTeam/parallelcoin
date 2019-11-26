package transport

import (
	"context"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/ipv4"

	"github.com/p9c/pod/pkg/fec"
	"github.com/p9c/pod/pkg/gcm"
	"github.com/p9c/pod/pkg/log"
)

type HandleFunc map[string]func(ctx interface{}) func(b []byte) (err error)

type MsgBuffer struct {
	Buffers [][]byte
	First   time.Time
	Decoded bool
	Source  *net.Addr
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
	SendConn        []*net.UDPConn
	listenAddress   *net.UDPAddr
	listenConn      net.PacketConn
	ciph            cipher.AEAD
	ctx             context.Context
	mx              *sync.Mutex
}

// NewConnection creates a new connection with a defined default send
// connection and listener and pre shared key password for encryption on the
// local network
func NewConnection(send, listen, preSharedKey string, maxDatagramSize int, ctx context.Context, multicast bool) (c *Connection, err error) {
	var sendAddr *net.UDPAddr
	sendConn := []*net.UDPConn{}
	var sC *net.UDPConn
	var listenAddr *net.UDPAddr
	var listenConn net.PacketConn
	var mcInterface net.Interface
	var ifi []net.Interface
	ifi, err = net.Interfaces()
	if err != nil {
		log.ERROR(err)
	}
	for i := range ifi {
		ad, _ := ifi[i].Addrs()
		if ifi[i].Flags&net.FlagMulticast != 0 &&
			ifi[i].HardwareAddr != nil &&
			ad != nil {
			mcInterface = ifi[i]
			break
		}
	}
	if listen != "" {
		if multicast {
			var conn net.PacketConn
			conn, err = net.ListenPacket("udp", listen)
			if err != nil {
				log.ERROR(err)
			}
			pc := ipv4.NewPacketConn(conn)
			err = pc.JoinGroup(&mcInterface, &net.UDPAddr{IP: net.IPv4(
				224, 0, 0, 1)})
			if err != nil {
				log.ERROR(err)
				err = conn.Close()
				if err != nil {
					log.ERROR(err)
					return
				}
			}
			listenConn = conn
		} else {
			listenAddr = GetUDPAddr(listen)
			listenConn, err = net.ListenUDP("udp", listenAddr)
			if err != nil {
				log.ERROR(err)
				return
			}
			log.SPEW(listenConn)
		}
	}
	if send != "" {
		mI, _ := mcInterface.Addrs()
		log.SPEW(mI)
		var listenWithoutEveryInterface string
		for i := range mI {
			_ = mI[i]
			log.DEBUG("ADDRESSS", mI[i])
			a := strings.Split(mI[i].String(), "/")[0]
			if strings.Count(a, ":") == 0 {
				listenWithoutEveryInterface = net.JoinHostPort(a, "0")
			}
		}
		log.DEBUG(listenWithoutEveryInterface)
		var laddr *net.UDPAddr
		laddr, err = net.ResolveUDPAddr("udp", listenWithoutEveryInterface)
		if err != nil {
			log.ERROR(err)
			return
		}
		sendAddr = GetUDPAddr(send)
		sC, err = net.DialUDP("udp", laddr, sendAddr)
		if err != nil {
			log.ERROR(err)
			return
		}
		sendConn = append(sendConn, sC)
	}
	ciph := gcm.GetCipher(preSharedKey)
	return &Connection{
		maxDatagramSize: maxDatagramSize,
		buffers:         make(map[string]*MsgBuffer),
		sendAddress:     sendAddr,
		SendConn:        sendConn,
		listenAddress:   listenAddr,
		listenConn:      listenConn,
		ciph:            ciph, // gcm.GetCipher(*cx.Config.MinerPass),
		ctx:             ctx,
		mx:              &sync.Mutex{},
	}, err
}

func (c *Connection) SetSendConn(ad ...string) (err error) {
	c.SendConn = []*net.UDPConn{}
	var sC *net.UDPConn
	for i := range ad {
		sC, err = net.DialUDP("udp", nil, GetUDPAddr(ad[i]))
		if err != nil {
			log.ERROR(err)
			return
		}
		c.SendConn = append(c.SendConn, sC)
	}
	return
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
	for _, sC := range c.SendConn {
		err = send(shards, sC)
		if err != nil {
			log.ERROR(err)
		}
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
	for _, sC := range c.SendConn {
		err = send(shards, sC)
		if err != nil {
			log.ERROR(err)
		}
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

func (c *Connection) Listen(handlers HandleFunc, ifc interface{},
	lastSent *time.Time, firstSender *string) (err error) {
	log.TRACE("setting read buffer")
	//err = c.listenConn.SetReadBuffer(c.maxDatagramSize)
	//if err != nil {
	//	log.ERROR(err)
	//	return
	//}
	buffer := make([]byte, c.maxDatagramSize)
	go func() {
		log.TRACE("starting connection handler")
	out:
		// read from socket until context is cancelled
		for {
			n, src, err := c.listenConn.ReadFrom(buffer)
			buf := buffer[:n]
			if err != nil {
				log.ERROR("ReadFromUDP failed:", err)
				continue
			}
			log.DEBUG("received message on UDP connection")
			log.SPEW(handlers)
			magic := string(buf[:4])
			if _, ok := handlers[magic]; ok {
				// if caller needs to know the liveness status of the
				// controller it is working on, the code below
				if lastSent != nil && firstSender != nil {
					//log.DEBUG("src", src.String(), "last", *firstSender)
					*lastSent = time.Now()
				}
				//log.DEBUG("received packet with magic:", magic)
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
					//log.DEBUG("new shard for",
					//	hex.EncodeToString([]byte(nonce)))
					if !bn.Decoded {
						bn.Buffers = append(bn.Buffers, shard)
						if len(bn.Buffers) >= 3 {
							// try to decode it
							var cipherText []byte
							//log.SPEW(bn.Buffers)
							cipherText, err = fec.Decode(bn.Buffers)
							if err != nil {
								log.ERROR(err)
								continue
							}
							//log.DEBUG("magic", magic, handlers[magic])
							//log.SPEW(cipherText)
							bn.Decoded = true
							err = handlers[magic](ifc)(cipherText)
							if err != nil {
								log.ERROR(err)
								continue
							}
							//log.DEBUG("called handler", magic)
						}
					} else {
						for i := range c.buffers {
							if i != nonce {
								// superseded messages can be deleted from the
								// buffers,
								// we don't add more data for the already
								// decoded.
								log.TRACE("deleting superseded buffer",
									hex.EncodeToString([]byte(i)))
								delete(c.buffers, i)
							}
						}
					}
				} else {
					log.TRACE("new message arriving",
						hex.EncodeToString([]byte(nonce)))
					c.buffers[nonce] = &MsgBuffer{[][]byte{},
						time.Now(), false, &src}
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
	log.DEBUG("multicast", address)
	log.SPEW(sendAddr)
	return
}
