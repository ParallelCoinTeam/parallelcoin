package sub

import (
	"net"
	"time"
)

var (
	uNet = "udp4"
	// // Maximum of 9 packets per message, so 16kb round is enough
	// defaultBufferSize = 16384
	// FEC expands message by 150%, we don't split message chunks over more than one packet
	maxMessageSize = 3072
	// default channel buffer sizes for Base
	baseChanBufs = 128
	// latency maximum
	latencyMax = time.Millisecond * 250
)

type (
	// BaseInterface is the core functions required for a Base
	BaseInterface interface {
		SetupListener()
	}
	// BaseCfg is the configuration for a Base
	BaseCfg struct {
		Handler    func(message Message)
		Listener   string
		Password   []byte
		BufferSize int
	}
	// Base is the common structure between a worker and a node
	Base struct {
		cfg       BaseCfg
		listener  *net.UDPConn
		packets   chan Packet
		incoming  chan Bundle
		returning chan Bundle
		trash     chan Bundle
		doneRet   chan bool
		message   chan Message
		quit      chan bool
	}
	// A Node is a server with some number of subscribers
	Node struct {
		Base
		// subscribers []*net.UDPAddr
	}
	// A Worker is a node that subscribes to a Node's messages
	Worker struct {
		Base
		// node *net.UDPAddr
	}
	// Packet is the structure of individual encoded packets of the message. These are made from a 9/3 Reed Solomon code and 9 are sent in distinct packets and only 3 are required to guarantee retransmit-free delivery.
	Packet struct {
		sender string // address packet was received from
		bytes  []byte // raw FEC encoded bytes of packet
	}
	// A Bundle is a collection of the received packets received from the same sender with up to 9 pieces.
	Bundle struct {
		uuid     int32
		sender   string
		received time.Time
		packets  [][]byte
	}
	// Message is the data reconstructed from a complete Bundle, containing data in messagepack format
	Message struct {
		uuid      int32
		sender    string
		timestamp time.Time
		bytes     []byte
	}
	// Subscription is the message sent by a worker node to request updates from the node
	Subscription struct {
		// address string
		// pubKey  []byte
	}
	// Confirmation is the reply message for a subscription request
	Confirmation struct {
		// subscriber string // confirming address of subscriber
		// pubKey     []byte // public key of server for message verification
	}
)
