package chainclient

import (
	"errors"
	"github.com/p9c/matrjoska/pkg/btcaddr"
	"github.com/p9c/matrjoska/pkg/chaincfg"
	"github.com/p9c/matrjoska/pkg/txscript"
	"sync"
	"time"
	
	"github.com/p9c/qu"
	
	"github.com/p9c/matrjoska/pkg/btcjson"
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/gcs"
	"github.com/p9c/matrjoska/pkg/gcs/builder"
	"github.com/p9c/matrjoska/pkg/rpcclient"
	"github.com/p9c/matrjoska/pkg/util"
	"github.com/p9c/matrjoska/pkg/waddrmgr"
	"github.com/p9c/matrjoska/pkg/wire"
	"github.com/p9c/matrjoska/pkg/wtxmgr"
)

// RPCClient represents a persistent client connection to a bitcoin RPC server for information regarding the current
// best block chain.
type RPCClient struct {
	*rpcclient.Client
	connConfig          *rpcclient.ConnConfig // Work around unexported field
	chainParams         *chaincfg.Params
	reconnectAttempts   int
	enqueueNotification chan interface{}
	dequeueNotification chan interface{}
	currentBlock        chan *waddrmgr.BlockStamp
	quit                qu.C
	wg                  sync.WaitGroup
	started             bool
	quitMtx             sync.Mutex
}

// NewRPCClient creates a client connection to the server described by the connect string. If disableTLS is false, the
// remote RPC certificate must be provided in the certs slice. The connection is not established immediately, but must
// be done using the Start method. If the remote server does not operate on the same bitcoin network as described by the
// passed chain parameters, the connection will be disconnected.
func NewRPCClient(
	chainParams *chaincfg.Params,
	connect, user, pass string,
	certs []byte,
	tls bool,
	reconnectAttempts int,
	quit qu.C,
) (*RPCClient, error) {
	W.Ln("creating new RPC client")
	if reconnectAttempts < 0 {
		return nil, errors.New("reconnectAttempts must be positive")
	}
	client := &RPCClient{
		connConfig: &rpcclient.ConnConfig{
			Host:                 connect,
			Endpoint:             "ws",
			User:                 user,
			Pass:                 pass,
			Certificates:         certs,
			DisableAutoReconnect: false,
			DisableConnectOnNew:  true,
			TLS:                  tls,
		},
		chainParams:         chainParams,
		reconnectAttempts:   reconnectAttempts,
		enqueueNotification: make(chan interface{}),
		dequeueNotification: make(chan interface{}),
		currentBlock:        make(chan *waddrmgr.BlockStamp),
		quit:                quit,
	}
	ntfnCallbacks := &rpcclient.NotificationHandlers{
		OnClientConnected:   client.onClientConnect,
		OnBlockConnected:    client.onBlockConnected,
		OnBlockDisconnected: client.onBlockDisconnected,
		OnRecvTx:            client.onRecvTx,
		OnRedeemingTx:       client.onRedeemingTx,
		OnRescanFinished:    client.onRescanFinished,
		OnRescanProgress:    client.onRescanProgress,
	}
	W.Ln("*actually* creating rpc client")
	rpcClient, e := rpcclient.New(client.connConfig, ntfnCallbacks, client.quit)
	if e != nil {
		return nil, e
	}
	// defer W.Ln("*succeeded* in making rpc client")
	client.Client = rpcClient
	return client, nil
}

// BackEnd returns the name of the driver.
func (c *RPCClient) BackEnd() string {
	return "pod"
}

// Start attempts to establish a client connection with the remote server. If successful, handler goroutines are started
// to process notifications sent by the server. After a limited number of connection attempts, this function gives up,
// and therefore will not block forever waiting for the connection to be established to a server that may not exist.
func (c *RPCClient) Start() (e error) {
	// D.Ln(c.connConfig)
	e = c.Connect(c.reconnectAttempts)
	if e != nil {
		return e
	}
	// Verify that the server is running on the expected network.
	net, e := c.GetCurrentNet()
	if e != nil {
		c.Disconnect()
		return e
	}
	if net != c.chainParams.Net {
		c.Disconnect()
		return errors.New("mismatched networks")
	}
	c.quitMtx.Lock()
	c.started = true
	c.quitMtx.Unlock()
	c.wg.Add(1)
	go c.handler()
	return nil
}

// Stop disconnects the client and signals the shutdown of all goroutines started by Start.
func (c *RPCClient) Stop() {
	c.quitMtx.Lock()
	select {
	case <-c.quit.Wait():
	default:
		c.quit.Q()
		c.Client.Shutdown()
		if !c.started {
			close(c.dequeueNotification)
		}
	}
	c.quitMtx.Unlock()
}

// Rescan wraps the normal Rescan command with an additional parameter that allows us to map an outpoint to the address
// in the chain that it pays to. This is useful when using BIP 158 filters as they include the prev pkScript rather than
// the full outpoint.
func (c *RPCClient) Rescan(
	startHash *chainhash.Hash, addrs []btcaddr.Address,
	outPoints map[wire.OutPoint]btcaddr.Address,
) (e error) {
	flatOutpoints := make([]*wire.OutPoint, 0, len(outPoints))
	for ops := range outPoints {
		flatOutpoints = append(flatOutpoints, &ops)
	}
	return c.Client.Rescan(startHash, addrs, flatOutpoints)
}

// WaitForShutdown blocks until both the client has finished disconnecting and all handlers have exited.
func (c *RPCClient) WaitForShutdown() {
	c.Client.WaitForShutdown()
	c.wg.Wait()
}

// Notifications returns a channel of parsed notifications sent by the remote bitcoin RPC server. This channel must be
// continually read or the process may abort for running out memory, as unread notifications are queued for later reads.
func (c *RPCClient) Notifications() <-chan interface{} {
	return c.dequeueNotification
}

// BlockStamp returns the latest block notified by the client, or an error if the client has been shut down.
func (c *RPCClient) BlockStamp() (*waddrmgr.BlockStamp, error) {
	select {
	case bs := <-c.currentBlock:
		return bs, nil
	case <-c.quit.Wait():
		return nil, errors.New("disconnected")
	}
}

// buildFilterBlocksWatchList constructs a watchlist used for matching against a cfilter from a FilterBlocksRequest. The
// watchlist will be populated with all external addresses, internal addresses, and outpoints contained in the request.
func buildFilterBlocksWatchList(req *FilterBlocksRequest) ([][]byte, error) {
	// Construct a watch list containing the script addresses of all internal and external addresses that were
	// requested, in addition to the set of outpoints currently being watched.
	watchListSize := len(req.ExternalAddrs) +
		len(req.InternalAddrs) +
		len(req.WatchedOutPoints)
	watchList := make([][]byte, 0, watchListSize)
	for _, addr := range req.ExternalAddrs {
		p2shAddr, e := txscript.PayToAddrScript(addr)
		if e != nil {
			return nil, e
		}
		watchList = append(watchList, p2shAddr)
	}
	for _, addr := range req.InternalAddrs {
		p2shAddr, e := txscript.PayToAddrScript(addr)
		if e != nil {
			return nil, e
		}
		watchList = append(watchList, p2shAddr)
	}
	for _, addr := range req.WatchedOutPoints {
		addr, e := txscript.PayToAddrScript(addr)
		if e != nil {
			return nil, e
		}
		watchList = append(watchList, addr)
	}
	return watchList, nil
}

// FilterBlocks scans the blocks contained in the FilterBlocksRequest for any addresses of interest. For each requested
// block, the corresponding compact filter will first be checked for matches, skipping those that do not report
// anything. If the filter returns a positive match, the full block will be fetched and filtered. This method returns a
// FilterBlocksResponse for the first block containing a matching address. If no matches are found in the range of
// blocks requested, the returned response will be nil.
func (c *RPCClient) FilterBlocks(req *FilterBlocksRequest,) (*FilterBlocksResponse, error) {
	blockFilterer := NewBlockFilterer(c.chainParams, req)
	// Construct the watchlist using the addresses and outpoints contained in the filter blocks request.
	watchList, e := buildFilterBlocksWatchList(req)
	if e != nil {
		return nil, e
	}
	// Iterate over the requested blocks, fetching the compact filter for each one, and matching it against the
	// watchlist generated above. If the filter returns a positive match, the full block is then requested and scanned
	// for addresses using the block filterer.
	for i, blk := range req.Blocks {
		rawFilter, e := c.GetCFilter(&blk.Hash, wire.GCSFilterRegular)
		if e != nil {
			return nil, e
		}
		// Ensure the filter is large enough to be deserialized.
		if len(rawFilter.Data) < 4 {
			continue
		}
		filter, e := gcs.FromNBytes(
			builder.DefaultP, builder.DefaultM, rawFilter.Data,
		)
		if e != nil {
			return nil, e
		}
		// Skip any empty filters.
		if filter.N() == 0 {
			continue
		}
		key := builder.DeriveKey(&blk.Hash)
		matched, e := filter.MatchAny(key, watchList)
		if e != nil {
			return nil, e
		} else if !matched {
			continue
		}
		T.F(
			"fetching block height=%d hash=%v",
			blk.Height, blk.Hash,
		)
		rawBlock, e := c.GetBlock(&blk.Hash)
		if e != nil {
			return nil, e
		}
		if !blockFilterer.FilterBlock(rawBlock) {
			continue
		}
		// If any external or internal addresses were detected in this block, we return them to the caller so that the
		// rescan windows can widened with subsequent addresses. The `BatchIndex` is returned so that the caller can
		// compute the *next* block from which to begin again.
		resp := &FilterBlocksResponse{
			BatchIndex:         uint32(i),
			BlockMeta:          blk,
			FoundExternalAddrs: blockFilterer.FoundExternal,
			FoundInternalAddrs: blockFilterer.FoundInternal,
			FoundOutPoints:     blockFilterer.FoundOutPoints,
			RelevantTxns:       blockFilterer.RelevantTxns,
		}
		return resp, nil
	}
	// No addresses were found for this range.
	return nil, nil
}

// parseBlock parses a btcws definition of the block a tx is mined it to the Block structure of the wtxmgr package, and the
// block index. This is done here since rpcclient doesn't parse this nicely for us.
func parseBlock(block *btcjson.BlockDetails) (*wtxmgr.BlockMeta, error) {
	if block == nil {
		return nil, nil
	}
	blkHash, e := chainhash.NewHashFromStr(block.Hash)
	if e != nil {
		return nil, e
	}
	blk := &wtxmgr.BlockMeta{
		Block: wtxmgr.Block{
			Height: block.Height,
			Hash:   *blkHash,
		},
		Time: time.Unix(block.Time, 0),
	}
	return blk, nil
}
func (c *RPCClient) onClientConnect() {
	select {
	case c.enqueueNotification <- ClientConnected{}:
	case <-c.quit.Wait():
	}
}
func (c *RPCClient) onBlockConnected(hash *chainhash.Hash, height int32, time time.Time) {
	select {
	case c.enqueueNotification <- BlockConnected{
		Block: wtxmgr.Block{
			Hash:   *hash,
			Height: height,
		},
		Time: time,
	}:
	case <-c.quit.Wait():
	}
}
func (c *RPCClient) onBlockDisconnected(hash *chainhash.Hash, height int32, time time.Time) {
	select {
	case c.enqueueNotification <- BlockDisconnected{
		Block: wtxmgr.Block{
			Hash:   *hash,
			Height: height,
		},
		Time: time,
	}:
	case <-c.quit.Wait():
	}
}
func (c *RPCClient) onRecvTx(tx *util.Tx, block *btcjson.BlockDetails) {
	blk, e := parseBlock(block)
	if e != nil {
		// Log and drop improper notification.
		E.Ln(
			"recvtx notification bad block:", e,
		)
		return
	}
	rec, e := wtxmgr.NewTxRecordFromMsgTx(tx.MsgTx(), time.Now())
	if e != nil {
		E.Ln("cannot create transaction record for relevant tx:", e)
		return
	}
	select {
	case c.enqueueNotification <- RelevantTx{rec, blk}:
	case <-c.quit.Wait():
	}
}
func (c *RPCClient) onRedeemingTx(tx *util.Tx, block *btcjson.BlockDetails) {
	// Handled exactly like recvtx notifications.
	c.onRecvTx(tx, block)
}
func (c *RPCClient) onRescanProgress(hash *chainhash.Hash, height int32, blkTime time.Time) {
	select {
	case c.enqueueNotification <- &RescanProgress{hash, height, blkTime}:
	case <-c.quit.Wait():
	}
}
func (c *RPCClient) onRescanFinished(hash *chainhash.Hash, height int32, blkTime time.Time) {
	select {
	case c.enqueueNotification <- &RescanFinished{hash, height, blkTime}:
	case <-c.quit.Wait():
	}
}

// handler maintains a queue of notifications and the current state (best block) of the chain.
func (c *RPCClient) handler() {
	hash, height, e := c.GetBestBlock()
	if e != nil {
		E.Ln("failed to receive best block from chain server:", e)
		c.Stop()
		c.wg.Done()
		return
	}
	bs := &waddrmgr.BlockStamp{Hash: *hash, Height: height}
	// TODO: Rather than leaving this as an unbounded queue for all types of notifications, try dropping ones where a
	//  later enqueued notification can fully invalidate one waiting to be processed. For example, blockconnected
	//  notifications for greater block heights can remove the need to process earlier blockconnected notifications still
	//  waiting here.
	var notifications []interface{}
	enqueue := c.enqueueNotification
	var dequeue chan interface{}
	var next interface{}
out:
	for {
		select {
		case n, ok := <-enqueue:
			if !ok {
				// If no notifications are queued for handling, the queue is finished.
				if len(notifications) == 0 {
					break out
				}
				// nil channel so no more reads can occur.
				enqueue = nil
				continue
			}
			if len(notifications) == 0 {
				next = n
				dequeue = c.dequeueNotification
			}
			notifications = append(notifications, n)
		case dequeue <- next:
			if n, ok := next.(BlockConnected); ok {
				bs = &waddrmgr.BlockStamp{
					Height: n.Height,
					Hash:   n.Hash,
				}
			}
			notifications[0] = nil
			notifications = notifications[1:]
			if len(notifications) != 0 {
				next = notifications[0]
			} else {
				// If no more notifications can be enqueued, the queue is finished.
				if enqueue == nil {
					break out
				}
				dequeue = nil
			}
		case c.currentBlock <- bs:
		case <-c.quit.Wait():
			D.Ln("legacy rpc handler stopping on quit channel close")
			break out
		}
	}
	c.Stop()
	close(c.dequeueNotification)
	c.wg.Done()
}

// POSTClient creates the equivalent HTTP POST rpcclient.Client.
func (c *RPCClient) POSTClient() (*rpcclient.Client, error) {
	configCopy := *c.connConfig
	configCopy.HTTPPostMode = true
	return rpcclient.New(&configCopy, nil, qu.T())
}
