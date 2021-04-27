package chainclient

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
	
	"github.com/p9c/qu"
	
	"github.com/tstranex/gozmq"
	
	"github.com/p9c/matrjoska/pkg/chaincfg"
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/rpcclient"
	"github.com/p9c/matrjoska/pkg/wire"
)

// BitcoindConn represents a persistent client connection to a bitcoind node that listens for events read from a ZMQ
// connection.
type BitcoindConn struct {
	started int32 // To be used atomically.
	stopped int32 // To be used atomically.
	// rescanClientCounter is an atomic counter that assigns a unique ID to each new bitcoind rescan client using the
	// current bitcoind connection.
	rescanClientCounter uint64
	// chainParams identifies the current network the bitcoind node is running on.
	chainParams *chaincfg.Params
	// client is the RPC client to the bitcoind node.
	client *rpcclient.Client
	// zmqBlockHost is the host listening for ZMQ connections that will be responsible for delivering raw transaction
	// events.
	zmqBlockHost string
	// zmqTxHost is the host listening for ZMQ connections that will be responsible for delivering raw transaction
	// events.
	zmqTxHost string
	// zmqPollInterval is the interval at which we'll attempt to retrieve an event from the ZMQ connection.
	zmqPollInterval time.Duration
	// rescanClients is the set of active bitcoind rescan clients to which ZMQ event notfications will be sent to.
	rescanClientsMtx sync.Mutex
	rescanClients    map[uint64]*BitcoindClient
	quit             qu.C
	wg               sync.WaitGroup
}

// NewBitcoindConn creates a client connection to the node described by the host string. The connection is not
// established immediately, but must be done using the Start method. If the remote node does not operate on the same
// bitcoin network as described by the passed chain parameters, the connection will be disconnected.
func NewBitcoindConn(
	chainParams *chaincfg.Params,
	host, user, pass, zmqBlockHost, zmqTxHost string,
	zmqPollInterval time.Duration,
) (*BitcoindConn, error) {
	clientCfg := &rpcclient.ConnConfig{
		Host:                 host,
		User:                 user,
		Pass:                 pass,
		DisableAutoReconnect: false,
		DisableConnectOnNew:  true,
		TLS:                  false,
		HTTPPostMode:         true,
	}
	client, e := rpcclient.New(clientCfg, nil, qu.T())
	if e != nil {
		return nil, e
	}
	conn := &BitcoindConn{
		chainParams:     chainParams,
		client:          client,
		zmqBlockHost:    zmqBlockHost,
		zmqTxHost:       zmqTxHost,
		zmqPollInterval: zmqPollInterval,
		rescanClients:   make(map[uint64]*BitcoindClient),
		quit:            qu.T(),
	}
	return conn, nil
}

// Start attempts to establish a RPC and ZMQ connection to a bitcoind node. If successful, a goroutine is spawned to
// read events from the ZMQ connection. It's possible for this function to fail due to a limited number of connection
// attempts. This is done to prevent waiting forever on the connection to be established in the case that the node is
// down.
func (c *BitcoindConn) Start() (e error) {
	if !atomic.CompareAndSwapInt32(&c.started, 0, 1) {
		return nil
	}
	// Verify that the node is running on the expected network.
	netw, e := c.getCurrentNet()
	if e != nil {
		c.client.Disconnect()
		return e
	}
	if netw != c.chainParams.Net {
		c.client.Disconnect()
		return fmt.Errorf(
			"expected network %v, got %v",
			c.chainParams.Net, netw,
		)
	}
	// Establish two different ZMQ connections to bitcoind to retrieve block and transaction event notifications. We'll
	// use two as a separation of concern to ensure one type of event isn't dropped from the connection queue due to
	// another type of event filling it up.
	zmqBlockConn, e := gozmq.Subscribe(
		c.zmqBlockHost, []string{"rawblock"},
	)
	if e != nil {
		c.client.Disconnect()
		return fmt.Errorf(
			"unable to subscribe for zmq block events: "+
				"%v", e,
		)
	}
	zmqTxConn, e := gozmq.Subscribe(
		c.zmqTxHost, []string{"rawtx"},
	)
	if e != nil {
		c.client.Disconnect()
		return fmt.Errorf(
			"unable to subscribe for zmq tx events: %v",
			e,
		)
	}
	c.wg.Add(2)
	go c.blockEventHandler(zmqBlockConn)
	go c.txEventHandler(zmqTxConn)
	return nil
}

// Stop terminates the RPC and ZMQ connection to a bitcoind node and removes any active rescan clients.
func (c *BitcoindConn) Stop() {
	if !atomic.CompareAndSwapInt32(&c.stopped, 0, 1) {
		return
	}
	for _, client := range c.rescanClients {
		client.Stop()
	}
	c.quit.Q()
	c.client.Shutdown()
	c.client.WaitForShutdown()
	c.wg.Wait()
}

// blockEventHandler reads raw blocks events from the ZMQ block socket and forwards them along to the current rescan
// clients.
//
// NOTE: This must be run as a goroutine.
func (c *BitcoindConn) blockEventHandler(conn *gozmq.Conn) {
	defer c.wg.Done()
	defer func() {
		if e := conn.Close(); E.Chk(e) {
		}
	}()
	I.Ln(
		"started listening for bitcoind block notifications via ZMQ on", c.zmqBlockHost,
	)
	for {
		// Before attempting to read from the ZMQ socket, we'll make sure to check if we've been requested to shut down.
		select {
		case <-c.quit.Wait():
			return
		default:
		}
		// Poll an event from the ZMQ socket.
		msgBytes, e := conn.Receive()
		if e != nil {
			// It's possible that the connection to the socket continuously times out, so we'll prevent logging this
			// error to prevent spamming the logs.
			netErr, ok := e.(net.Error)
			if ok && netErr.Timeout() {
				continue
			}
			E.Ln(
				"unable to receive ZMQ rawblock message:", e,
			)
			continue
		}
		// We have an event! We'll now ensure it is a block event, deserialize it, and report it to the different rescan
		// clients.
		eventType := string(msgBytes[0])
		switch eventType {
		case "rawblock":
			block := &wire.Block{}
			r := bytes.NewReader(msgBytes[1])
			if e := block.Deserialize(r); E.Chk(e) {
				E.Ln(
					"unable to deserialize block:", e,
				)
				continue
			}
			c.rescanClientsMtx.Lock()
			for _, client := range c.rescanClients {
				select {
				case client.zmqBlockNtfns <- block:
				case <-client.quit.Wait():
				case <-c.quit.Wait():
					c.rescanClientsMtx.Unlock()
					return
				}
			}
			c.rescanClientsMtx.Unlock()
		default:
			// It's possible that the message wasn't fully read if bitcoind shuts down, which will produce an unreadable
			// event type. To prevent from logging it, we'll make sure it conforms to the ASCII standard.
			if eventType == "" || !isASCII(eventType) {
				continue
			}
			W.Ln(
				"received unexpected event type from rawblock subscription:",
				eventType,
			)
		}
	}
}

// txEventHandler reads raw blocks events from the ZMQ block socket and forwards them along to the current rescan
// clients.
//
// NOTE: This must be run as a goroutine.
func (c *BitcoindConn) txEventHandler(conn *gozmq.Conn) {
	defer c.wg.Done()
	defer func() {
		if e := conn.Close(); E.Chk(e) {
		}
	}()
	I.Ln(
		"started listening for bitcoind transaction notifications via ZMQ on",
		c.zmqTxHost,
	)
	for {
		// Before attempting to read from the ZMQ socket, we'll make sure to check if we've been requested to shut down.
		select {
		case <-c.quit.Wait():
			return
		default:
		}
		// Poll an event from the ZMQ socket.
		msgBytes, e := conn.Receive()
		if e != nil {
			// It's possible that the connection to the socket continuously times out, so we'll prevent logging this
			// error to prevent spamming the logs.
			netErr, ok := e.(net.Error)
			if ok && netErr.Timeout() {
				continue
			}
			E.Ln(
				"unable to receive ZMQ rawtx message:", e,
			)
			continue
		}
		// We have an event! We'll now ensure it is a transaction event, deserialize it, and report it to the different
		// rescan clients.
		eventType := string(msgBytes[0])
		switch eventType {
		case "rawtx":
			tx := &wire.MsgTx{}
			r := bytes.NewReader(msgBytes[1])
			if e := tx.Deserialize(r); E.Chk(e) {
				E.Ln(
					"unable to deserialize transaction:", e,
				)
				continue
			}
			c.rescanClientsMtx.Lock()
			for _, client := range c.rescanClients {
				select {
				case client.zmqTxNtfns <- tx:
				case <-client.quit.Wait():
				case <-c.quit.Wait():
					c.rescanClientsMtx.Unlock()
					return
				}
			}
			c.rescanClientsMtx.Unlock()
		default:
			// It's possible that the message wasn't fully read if bitcoind shuts down, which will produce an unreadable
			// event type. To prevent from logging it, we'll make sure it conforms to the ASCII standard.
			if eventType == "" || !isASCII(eventType) {
				continue
			}
			W.Ln(
				"received unexpected event type from rawtx subscription:",
				eventType,
			)
		}
	}
}

// getCurrentNet returns the network on which the bitcoind node is running.
func (c *BitcoindConn) getCurrentNet() (wire.BitcoinNet, error) {
	hash, e := c.client.GetBlockHash(0)
	if e != nil {
		return 0, e
	}
	switch *hash {
	case *chaincfg.TestNet3Params.GenesisHash:
		return chaincfg.TestNet3Params.Net, nil
	case *chaincfg.RegressionTestParams.GenesisHash:
		return chaincfg.RegressionTestParams.Net, nil
	case *chaincfg.MainNetParams.GenesisHash:
		return chaincfg.MainNetParams.Net, nil
	default:
		return 0, fmt.Errorf("unknown network with genesis hash %v", hash)
	}
}

// NewBitcoindClient returns a bitcoind client using the current bitcoind connection. This allows us to share the same
// connection using multiple clients.
func (c *BitcoindConn) NewBitcoindClient() *BitcoindClient {
	return &BitcoindClient{
		quit:              qu.T(),
		id:                atomic.AddUint64(&c.rescanClientCounter, 1),
		chainParams:       c.chainParams,
		chainConn:         c,
		rescanUpdate:      make(chan interface{}),
		watchedAddresses:  make(map[string]struct{}),
		watchedOutPoints:  make(map[wire.OutPoint]struct{}),
		watchedTxs:        make(map[chainhash.Hash]struct{}),
		notificationQueue: NewConcurrentQueue(20),
		zmqTxNtfns:        make(chan *wire.MsgTx),
		zmqBlockNtfns:     make(chan *wire.Block),
		mempool:           make(map[chainhash.Hash]struct{}),
		expiredMempool:    make(map[int32]map[chainhash.Hash]struct{}),
	}
}

// AddClient adds a client to the set of active rescan clients of the current chain connection. This allows the
// connection to include the specified client in its notification delivery.
//
// NOTE: This function is safe for concurrent access.
func (c *BitcoindConn) AddClient(client *BitcoindClient) {
	c.rescanClientsMtx.Lock()
	defer c.rescanClientsMtx.Unlock()
	c.rescanClients[client.id] = client
}

// RemoveClient removes the client with the given ID from the set of active rescan clients. Once removed, the client
// will no longer receive block and transaction notifications from the chain connection.
//
// NOTE: This function is safe for concurrent access.
func (c *BitcoindConn) RemoveClient(id uint64) {
	c.rescanClientsMtx.Lock()
	defer c.rescanClientsMtx.Unlock()
	delete(c.rescanClients, id)
}

// isASCII is a helper method that checks whether all bytes in `data` would be printable ASCII characters if interpreted
// as a string.
func isASCII(s string) bool {
	for _, c := range s {
		if c < 32 || c > 126 {
			return false
		}
	}
	return true
}
