package chainclient

import (
	"container/list"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/p9c/matrjoska/pkg/btcaddr"
	"github.com/p9c/matrjoska/pkg/chaincfg"
	"sync"
	"sync/atomic"
	"time"
	
	"github.com/p9c/qu"
	
	"github.com/p9c/matrjoska/pkg/btcjson"
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/txscript"
	"github.com/p9c/matrjoska/pkg/util"
	am "github.com/p9c/matrjoska/pkg/waddrmgr"
	"github.com/p9c/matrjoska/pkg/wire"
	tm "github.com/p9c/matrjoska/pkg/wtxmgr"
)

var (
	// ErrBitcoindClientShuttingDown is an error returned when we attempt to receive a notification for a specific item
	// and the bitcoind client is in the middle of shutting down.
	ErrBitcoindClientShuttingDown = errors.New("client is shutting down")
)

// BitcoindClient represents a persistent client connection to a bitcoind server for information regarding the current
// best block chain.
type BitcoindClient struct {
	started int32 // To be used atomically.
	stopped int32 // To be used atomically.
	// birthday is the earliest time for which we should begin scanning the chain.
	birthday time.Time
	// chainParams are the parameters of the current chain this client is active under.
	chainParams *chaincfg.Params
	// id is the unique ID of this client assigned by the backing bitcoind connection.
	id uint64
	// chainConn is the backing client to our rescan client that contains the RPC and ZMQ connections to a bitcoind
	// node.
	chainConn *BitcoindConn
	// bestBlock keeps track of the tip of the current best chain.
	bestBlockMtx sync.RWMutex
	bestBlock    am.BlockStamp
	// notifyBlocks signals whether the client is sending block notifications to the caller.
	notifyBlocks uint32
	// rescanUpdate is a channel will be sent items that we should match transactions against while processing a chain
	// rescan to determine if they are relevant to the client.
	rescanUpdate chan interface{}
	// watchedAddresses, watchedOutPoints, and watchedTxs are the set of items we should match transactions against
	// while processing a chain rescan to determine if they are relevant to the client.
	watchMtx         sync.RWMutex
	watchedAddresses map[string]struct{}
	watchedOutPoints map[wire.OutPoint]struct{}
	watchedTxs       map[chainhash.Hash]struct{}
	// mempool keeps track of all relevant transactions that have yet to be confirmed. This is used to shortcut the
	// filtering process of a transaction when a new confirmed transaction notification is received.
	//
	// NOTE: This requires the watchMtx to be held.
	mempool map[chainhash.Hash]struct{}
	// expiredMempool keeps track of a set of confirmed transactions along with the height at which they were included
	// in a block. These transactions will then be removed from the mempool after a period of 288 blocks. This is done
	// to ensure the transactions are safe from a reorg in the chain.
	//
	// NOTE: This requires the watchMtx to be held.
	expiredMempool map[int32]map[chainhash.Hash]struct{}
	// notificationQueue is a concurrent unbounded queue that handles dispatching notifications to the subscriber of
	// this client.
	//
	// TODO: Rather than leaving this as an unbounded queue for all types of notifications, try dropping ones where a
	//  later enqueued notification can fully invalidate one waiting to be processed. For example, BlockConnected
	//  notifications for greater block heights can remove the need to process earlier notifications still waiting to be
	//  processed.
	notificationQueue *ConcurrentQueue
	// zmqTxNtfns is a channel through which ZMQ transaction events will be retrieved from the backing bitcoind
	// connection.
	zmqTxNtfns chan *wire.MsgTx
	// zmqBlockNtfns is a channel through which ZMQ block events will be retrieved from the backing bitcoind connection.
	zmqBlockNtfns chan *wire.Block
	quit          qu.C
	wg            sync.WaitGroup
}

// A compile-time check to ensure that BitcoindClient satisfies the chainclient.Interface interface.
var _ Interface = (*BitcoindClient)(nil)

// BackEnd returns the name of the driver.
func (c *BitcoindClient) BackEnd() string {
	return "bitcoind"
}

// GetBestBlock returns the highest block known to bitcoind.
func (c *BitcoindClient) GetBestBlock() (*chainhash.Hash, int32, error) {
	bcInfo, e := c.chainConn.client.GetBlockChainInfo()
	if e != nil {
		return nil, 0, e
	}
	hash, e := chainhash.NewHashFromStr(bcInfo.BestBlockHash)
	if e != nil {
		return nil, 0, e
	}
	return hash, bcInfo.Blocks, nil
}

// GetBlockHeight returns the height for the hash, if known, or returns an error.
func (c *BitcoindClient) GetBlockHeight(hash *chainhash.Hash) (int32, error) {
	header, e := c.chainConn.client.GetBlockHeaderVerbose(hash)
	if e != nil {
		return 0, e
	}
	return header.Height, nil
}

// GetBlock returns a block from the hash.
func (c *BitcoindClient) GetBlock(hash *chainhash.Hash) (*wire.Block, error) {
	return c.chainConn.client.GetBlock(hash)
}

// GetBlockVerbose returns a verbose block from the hash.
func (c *BitcoindClient) GetBlockVerbose(
	hash *chainhash.Hash,
) (*btcjson.GetBlockVerboseResult, error) {
	return c.chainConn.client.GetBlockVerbose(hash)
}

// GetBlockHash returns a block hash from the height.
func (c *BitcoindClient) GetBlockHash(height int64) (*chainhash.Hash, error) {
	return c.chainConn.client.GetBlockHash(height)
}

// GetBlockHeader returns a block header from the hash.
func (c *BitcoindClient) GetBlockHeader(
	hash *chainhash.Hash,
) (*wire.BlockHeader, error) {
	return c.chainConn.client.GetBlockHeader(hash)
}

// GetBlockHeaderVerbose returns a block header from the hash.
func (c *BitcoindClient) GetBlockHeaderVerbose(
	hash *chainhash.Hash,
) (*btcjson.GetBlockHeaderVerboseResult, error) {
	return c.chainConn.client.GetBlockHeaderVerbose(hash)
}

// GetRawTransactionVerbose returns a transaction from the tx hash.
func (c *BitcoindClient) GetRawTransactionVerbose(
	hash *chainhash.Hash,
) (*btcjson.TxRawResult, error) {
	return c.chainConn.client.GetRawTransactionVerbose(hash)
}

// GetTxOut returns a txout from the outpoint info provided.
func (c *BitcoindClient) GetTxOut(
	txHash *chainhash.Hash, index uint32,
	mempool bool,
) (*btcjson.GetTxOutResult, error) {
	return c.chainConn.client.GetTxOut(txHash, index, mempool)
}

// SendRawTransaction sends a raw transaction via bitcoind.
func (c *BitcoindClient) SendRawTransaction(
	tx *wire.MsgTx,
	allowHighFees bool,
) (*chainhash.Hash, error) {
	return c.chainConn.client.SendRawTransaction(tx, allowHighFees)
}

// Notifications returns a channel to retrieve notifications from.
//
// NOTE: This is part of the chainclient.Interface interface.
func (c *BitcoindClient) Notifications() <-chan interface{} {
	return c.notificationQueue.ChanOut()
}

// NotifyReceived allows the chain backend to notify the caller whenever a transaction pays to any of the given
// addresses.
//
// NOTE: This is part of the chainclient.Interface interface.
func (c *BitcoindClient) NotifyReceived(addrs []btcaddr.Address) (e error) {
	if e = c.NotifyBlocks(); E.Chk(e) {
	}
	select {
	case c.rescanUpdate <- addrs:
	case <-c.quit.Wait():
		return ErrBitcoindClientShuttingDown
	}
	return nil
}

// NotifySpent allows the chain backend to notify the caller whenever a transaction spends any of the given outpoints.
func (c *BitcoindClient) NotifySpent(outPoints []*wire.OutPoint) (e error) {
	if e = c.NotifyBlocks(); E.Chk(e) {
	}
	select {
	case c.rescanUpdate <- outPoints:
	case <-c.quit.Wait():
		return ErrBitcoindClientShuttingDown
	}
	return nil
}

// NotifyTx allows the chain backend to notify the caller whenever any of the given transactions confirm within the
// chain.
func (c *BitcoindClient) NotifyTx(txids []chainhash.Hash) (e error) {
	if e = c.NotifyBlocks(); E.Chk(e) {
	}
	select {
	case c.rescanUpdate <- txids:
	case <-c.quit.Wait():
		return ErrBitcoindClientShuttingDown
	}
	return nil
}

// NotifyBlocks allows the chain backend to notify the caller whenever a block is connected or disconnected.
//
// NOTE: This is part of the chainclient.Interface interface.
func (c *BitcoindClient) NotifyBlocks() (e error) {
	atomic.StoreUint32(&c.notifyBlocks, 1)
	return nil
}

// shouldNotifyBlocks determines whether the client should send block notifications to the caller.
func (c *BitcoindClient) shouldNotifyBlocks() bool {
	return atomic.LoadUint32(&c.notifyBlocks) == 1
}

// LoadTxFilter uses the given filters to what we should match transactions against to determine if they are relevant to
// the client. The reset argument is used to reset the current filters.
//
// The current filters supported are of the following types:
//	[]util.Address
//	[]wire.OutPoint
//	[]*wire.OutPoint
//	map[wire.OutPoint]util.Address
//	[]chainhash.Hash
//	[]*chainhash.Hash
func (c *BitcoindClient) LoadTxFilter(reset bool, filters ...interface{}) (e error) {
	if reset {
		select {
		case c.rescanUpdate <- struct{}{}:
		case <-c.quit.Wait():
			return ErrBitcoindClientShuttingDown
		}
	}
	updateFilter := func(filter interface{}) (e error) {
		select {
		case c.rescanUpdate <- filter:
		case <-c.quit.Wait():
			return ErrBitcoindClientShuttingDown
		}
		return nil
	}
	// In order to make this operation atomic, we'll iterate through the filters twice: the first to ensure there aren't
	// any unsupported filter types, and the second to actually update our filters.
	for _, filter := range filters {
		switch filter := filter.(type) {
		case []btcaddr.Address, []wire.OutPoint, []*wire.OutPoint,
			map[wire.OutPoint]btcaddr.Address, []chainhash.Hash,
			[]*chainhash.Hash:
			// Proceed to check the next filter type.
		default:
			return fmt.Errorf("unsupported filter type %Ter", filter)
		}
	}
	for _, filter := range filters {
		if e := updateFilter(filter); E.Chk(e) {
			return e
		}
	}
	return nil
}

// RescanBlocks rescans any blocks passed, returning only the blocks that matched as []json.BlockDetails.
func (c *BitcoindClient) RescanBlocks(
	blockHashes []chainhash.Hash,
) ([]btcjson.RescannedBlock, error) {
	rescannedBlocks := make([]btcjson.RescannedBlock, 0, len(blockHashes))
	for _, hash := range blockHashes {
		header, e := c.GetBlockHeaderVerbose(&hash)
		if e != nil {
			E.F(
				"unable to get header %s from bitcoind: %s",
				hash, e,
			)
			continue
		}
		block, e := c.GetBlock(&hash)
		if e != nil {
			E.F(
				"unable to get block %s from bitcoind: %s",
				hash, e,
			)
			continue
		}
		relevantTxs, e := c.filterBlock(block, header.Height, false)
		if e != nil {
		}
		if len(relevantTxs) > 0 {
			rescannedBlock := btcjson.RescannedBlock{
				Hash: hash.String(),
			}
			for _, tx := range relevantTxs {
				rescannedBlock.Transactions = append(
					rescannedBlock.Transactions,
					hex.EncodeToString(tx.SerializedTx),
				)
			}
			rescannedBlocks = append(rescannedBlocks, rescannedBlock)
		}
	}
	return rescannedBlocks, nil
}

// Rescan rescans from the block with the given hash until the current block, after adding the passed addresses and
// outpoints to the client's watch list.
func (c *BitcoindClient) Rescan(
	blockHash *chainhash.Hash,
	addresses []btcaddr.Address, outPoints map[wire.OutPoint]btcaddr.Address,
) (e error) {
	// A block hash is required to use as the starting point of the rescan.
	if blockHash == nil {
		return errors.New("rescan requires a starting block hash")
	}
	// We'll then update our filters with the given outpoints and addresses.
	select {
	case c.rescanUpdate <- addresses:
	case <-c.quit.Wait():
		return ErrBitcoindClientShuttingDown
	}
	select {
	case c.rescanUpdate <- outPoints:
	case <-c.quit.Wait():
		return ErrBitcoindClientShuttingDown
	}
	// Once the filters have been updated, we can begin the rescan.
	select {
	case c.rescanUpdate <- *blockHash:
	case <-c.quit.Wait():
		return ErrBitcoindClientShuttingDown
	}
	return nil
}

// Start initializes the bitcoind rescan client using the backing bitcoind connection and starts all goroutines
// necessary in order to process rescans and ZMQ notifications.
//
// NOTE: This is part of the chainclient.Interface interface.
func (c *BitcoindClient) Start() (e error) {
	if !atomic.CompareAndSwapInt32(&c.started, 0, 1) {
		return nil
	}
	// Start the notification queue and immediately dispatch a ClientConnected notification to the caller. This is
	// needed as some of the callers will require this notification before proceeding.
	c.notificationQueue.Start()
	c.notificationQueue.ChanIn() <- ClientConnected{}
	// Retrieve the best block of the chain.
	bestHash, bestHeight, e := c.GetBestBlock()
	if e != nil {
		return fmt.Errorf("unable to retrieve best block: %v", e)
	}
	bestHeader, e := c.GetBlockHeaderVerbose(bestHash)
	if e != nil {
		return fmt.Errorf(
			"unable to retrieve header for best block: "+
				"%v", e,
		)
	}
	c.bestBlockMtx.Lock()
	c.bestBlock = am.BlockStamp{
		Hash:      *bestHash,
		Height:    bestHeight,
		Timestamp: time.Unix(bestHeader.Time, 0),
	}
	c.bestBlockMtx.Unlock()
	// Once the client has started successfully, we'll include it in the set of rescan clients of the backing bitcoind
	// connection in order to received ZMQ event notifications.
	c.chainConn.AddClient(c)
	c.wg.Add(2)
	go c.rescanHandler()
	go c.ntfnHandler()
	return nil
}

// Stop stops the bitcoind rescan client from processing rescans and ZMQ notifications.
//
// NOTE: This is part of the chainclient.Interface interface.
func (c *BitcoindClient) Stop() {
	if !atomic.CompareAndSwapInt32(&c.stopped, 0, 1) {
		return
	}
	c.quit.Q()
	// Remove this client's reference from the bitcoind connection to prevent sending notifications to it after it's
	// been stopped.
	c.chainConn.RemoveClient(c.id)
	c.notificationQueue.Stop()
}

// WaitForShutdown blocks until the client has finished disconnecting and all handlers have exited.
//
// NOTE: This is part of the chainclient.Interface interface.
func (c *BitcoindClient) WaitForShutdown() {
	c.wg.Wait()
}

// rescanHandler handles the logic needed for the caller to trigger a chain rescan.
//
// NOTE: This must be called as a goroutine.
func (c *BitcoindClient) rescanHandler() {
	defer c.wg.Done()
	for {
		select {
		case update := <-c.rescanUpdate:
			switch update := update.(type) {
			// We're clearing the filters.
			case struct{}:
				c.watchMtx.Lock()
				c.watchedOutPoints = make(map[wire.OutPoint]struct{})
				c.watchedAddresses = make(map[string]struct{})
				c.watchedTxs = make(map[chainhash.Hash]struct{})
				c.watchMtx.Unlock()
			// We're adding the addresses to our filter.
			case []btcaddr.Address:
				c.watchMtx.Lock()
				for _, addr := range update {
					c.watchedAddresses[addr.String()] = struct{}{}
				}
				c.watchMtx.Unlock()
			// We're adding the outpoints to our filter.
			case []wire.OutPoint:
				c.watchMtx.Lock()
				for _, op := range update {
					c.watchedOutPoints[op] = struct{}{}
				}
				c.watchMtx.Unlock()
			case []*wire.OutPoint:
				c.watchMtx.Lock()
				for _, op := range update {
					c.watchedOutPoints[*op] = struct{}{}
				}
				c.watchMtx.Unlock()
			// We're adding the outpoints that map to the scripts that we should scan for to our filter.
			case map[wire.OutPoint]btcaddr.Address:
				c.watchMtx.Lock()
				for op := range update {
					c.watchedOutPoints[op] = struct{}{}
				}
				c.watchMtx.Unlock()
			// We're adding the transactions to our filter.
			case []chainhash.Hash:
				c.watchMtx.Lock()
				for _, txid := range update {
					c.watchedTxs[txid] = struct{}{}
				}
				c.watchMtx.Unlock()
			case []*chainhash.Hash:
				c.watchMtx.Lock()
				for _, txid := range update {
					c.watchedTxs[*txid] = struct{}{}
				}
				c.watchMtx.Unlock()
			// We're starting a rescan from the hash.
			case chainhash.Hash:
				if e := c.rescan(update); E.Chk(e) {
					E.Ln(
						"unable to complete chain rescan:", e,
					)
				}
			default:
				W.F(
					"received unexpected filter type %Ter", update,
				)
			}
		case <-c.quit.Wait():
			return
		}
	}
}

// ntfnHandler handles the logic to retrieve ZMQ notifications from the backing bitcoind connection.
//
// NOTE: This must be called as a goroutine.
func (c *BitcoindClient) ntfnHandler() {
	defer c.wg.Done()
	for {
		select {
		case tx := <-c.zmqTxNtfns:
			var e error
			if _, _, e = c.filterTx(tx, nil, true); E.Chk(e) {
				E.F(
					"unable to filter transaction %v: %v %s",
					tx.TxHash(), e,
				)
			}
		case newBlock := <-c.zmqBlockNtfns:
			// If the new block's previous hash matches the best hash known to us, then the new block is the next
			// successor, so we'll update our best block to reflect this and determine if this new block matches any of
			// our existing filters.
			c.bestBlockMtx.Lock()
			bestBlock := c.bestBlock
			c.bestBlockMtx.Unlock()
			if newBlock.Header.PrevBlock == bestBlock.Hash {
				newBlockHeight := bestBlock.Height + 1
				_, e := c.filterBlock(
					newBlock, newBlockHeight, true,
				)
				if e != nil {
					E.F(
						"unable to filter block %v: %v",
						newBlock.BlockHash(), e,
					)
					continue
				}
				// With the block successfully filtered, we'll make it our new best block.
				bestBlock.Hash = newBlock.BlockHash()
				bestBlock.Height = newBlockHeight
				bestBlock.Timestamp = newBlock.Header.Timestamp
				c.bestBlockMtx.Lock()
				c.bestBlock = bestBlock
				c.bestBlockMtx.Unlock()
				continue
			}
			// Otherwise, we've encountered a reorg.
			if e := c.reorg(bestBlock, newBlock); E.Chk(e) {
				E.F(
					"unable to process chain reorg:", e,
				)
			}
		case <-c.quit.Wait():
			return
		}
	}
}

// SetBirthday sets the birthday of the bitcoind rescan client.
//
// NOTE: This should be done before the client has been started in order for it to properly carry its duties.
func (c *BitcoindClient) SetBirthday(t time.Time) {
	c.birthday = t
}

// BlockStamp returns the latest block notified by the client, or an error if the client has been shut down.
func (c *BitcoindClient) BlockStamp() (*am.BlockStamp, error) {
	c.bestBlockMtx.RLock()
	bestBlock := c.bestBlock
	c.bestBlockMtx.RUnlock()
	return &bestBlock, nil
}

// onBlockConnected is a callback that's executed whenever a new block has been detected. This will queue a
// BlockConnected notification to the caller.
func (c *BitcoindClient) onBlockConnected(
	hash *chainhash.Hash, height int32,
	timestamp time.Time,
) {
	if c.shouldNotifyBlocks() {
		select {
		case c.notificationQueue.ChanIn() <- BlockConnected{
			Block: tm.Block{
				Hash:   *hash,
				Height: height,
			},
			Time: timestamp,
		}:
		case <-c.quit.Wait():
		}
	}
}

// onFilteredBlockConnected is an alternative callback that's executed whenever a new block has been detected. It serves
// the same purpose as onBlockConnected, but it also includes a list of the relevant transactions found within the block
// being connected. This will queue a FilteredBlockConnected notification to the caller.
func (c *BitcoindClient) onFilteredBlockConnected(
	height int32,
	header *wire.BlockHeader, relevantTxs []*tm.TxRecord,
) {
	if c.shouldNotifyBlocks() {
		select {
		case c.notificationQueue.ChanIn() <- FilteredBlockConnected{
			Block: &tm.BlockMeta{
				Block: tm.Block{
					Hash:   header.BlockHash(),
					Height: height,
				},
				Time: header.Timestamp,
			},
			RelevantTxs: relevantTxs,
		}:
		case <-c.quit.Wait():
		}
	}
}

// onBlockDisconnected is a callback that's executed whenever a block has been disconnected. This will queue a
// BlockDisconnected notification to the caller with the details of the block being disconnected.
func (c *BitcoindClient) onBlockDisconnected(
	hash *chainhash.Hash, height int32,
	timestamp time.Time,
) {
	if c.shouldNotifyBlocks() {
		select {
		case c.notificationQueue.ChanIn() <- BlockDisconnected{
			Block: tm.Block{
				Hash:   *hash,
				Height: height,
			},
			Time: timestamp,
		}:
		case <-c.quit.Wait():
		}
	}
}

// onRelevantTx is a callback that's executed whenever a transaction is relevant to the caller. This means that the
// transaction matched a specific item in the client's different filters. This will queue a RelevantTx notification to
// the caller.
func (c *BitcoindClient) onRelevantTx(
	tx *tm.TxRecord,
	blockDetails *btcjson.BlockDetails,
) {
	block, e := parseBlock(blockDetails)
	if e != nil {
		E.Ln(
			"unable to send onRelevantTx notification, failed parse block:",
			e,
		)
		return
	}
	select {
	case c.notificationQueue.ChanIn() <- RelevantTx{
		TxRecord: tx,
		Block:    block,
	}:
	case <-c.quit.Wait():
	}
}

// onRescanProgress is a callback that's executed whenever a rescan is in progress. This will queue a RescanProgress
// notification to the caller with the current rescan progress details.
func (c *BitcoindClient) onRescanProgress(
	hash *chainhash.Hash, height int32,
	timestamp time.Time,
) {
	select {
	case c.notificationQueue.ChanIn() <- &RescanProgress{
		Hash:   hash,
		Height: height,
		Time:   timestamp,
	}:
	case <-c.quit.Wait():
	}
}

// onRescanFinished is a callback that's executed whenever a rescan has finished. This will queue a RescanFinished
// notification to the caller with the details of the last block in the range of the rescan.
func (c *BitcoindClient) onRescanFinished(
	hash *chainhash.Hash, height int32,
	timestamp time.Time,
) {
	I.F(
		"rescan finished at %d (%s)",
		height, hash,
	)
	select {
	case c.notificationQueue.ChanIn() <- &RescanFinished{
		Hash:   hash,
		Height: height,
		Time:   timestamp,
	}:
	case <-c.quit.Wait():
	}
}

// reorg processes a reorganization during chain synchronization. This is separate from a rescan's handling of a reorg.
// This will rewind back until it finds a common ancestor and notify all the new blocks since then.
func (c *BitcoindClient) reorg(
	currentBlock am.BlockStamp,
	reorgBlock *wire.Block,
) (e error) {
	D.Ln("possible reorg at block", reorgBlock.BlockHash())
	// Retrieve the best known height based on the block which caused the reorg. This way, we can preserve the chain of
	// blocks we need to retrieve.
	bestHash := reorgBlock.BlockHash()
	var bestHeight int32
	bestHeight, e = c.GetBlockHeight(&bestHash)
	if e != nil {
		return e
	}
	if bestHeight < currentBlock.Height {
		D.Ln("detected multiple reorgs")
		return nil
	}
	// We'll now keep track of all the blocks known to the *chain*, starting from the best block known to us until the
	// best block in the chain. This will let us fast-forward despite any future reorgs.
	blocksToNotify := list.New()
	blocksToNotify.PushFront(reorgBlock)
	previousBlock := reorgBlock.Header.PrevBlock
	for i := bestHeight - 1; i >= currentBlock.Height; i-- {
		var block *wire.Block
		block, e = c.GetBlock(&previousBlock)
		if e != nil {
			return e
		}
		blocksToNotify.PushFront(block)
		previousBlock = block.Header.PrevBlock
	}
	// Rewind back to the last common ancestor block using the previous block hash from each header to avoid any race
	// conditions. If we encounter more reorgs, they'll be queued and we'll repeat the cycle.
	//
	// We'll start by retrieving the header to the best block known to us.
	currentHeader, e := c.GetBlockHeader(&currentBlock.Hash)
	if e != nil {
		return e
	}
	// Then, we'll walk backwards in the chain until we find our common ancestor.
	for previousBlock != currentHeader.PrevBlock {
		// Since the previous hashes don't match, the current block has been reorged out of the chain, so we should send
		// a BlockDisconnected notification for it.
		D.F(
			"disconnecting block %d (%v) %s",
			currentBlock.Height,
			currentBlock.Hash,
		)
		c.onBlockDisconnected(
			&currentBlock.Hash, currentBlock.Height,
			currentBlock.Timestamp,
		)
		// Our current block should now reflect the previous one to continue the common ancestor search.
		currentHeader, e = c.GetBlockHeader(&currentHeader.PrevBlock)
		if e != nil {
			return e
		}
		currentBlock.Height--
		currentBlock.Hash = currentHeader.PrevBlock
		currentBlock.Timestamp = currentHeader.Timestamp
		// Store the correct block in our list in order to notify it once we've found our common ancestor.
		block, e := c.GetBlock(&previousBlock)
		if e != nil {
			return e
		}
		blocksToNotify.PushFront(block)
		previousBlock = block.Header.PrevBlock
	}
	// Disconnect the last block from the old chain. Since the previous block remains the same between the old and new
	// chains, the tip will now be the last common ancestor.
	D.F(
		"disconnecting block %d (%v) %s",
		currentBlock.Height, currentBlock.Hash,
	)
	c.onBlockDisconnected(
		&currentBlock.Hash, currentBlock.Height, currentHeader.Timestamp,
	)
	currentBlock.Height--
	// Now we fast-forward to the new block, notifying along the way.
	for blocksToNotify.Front() != nil {
		nextBlock := blocksToNotify.Front().Value.(*wire.Block)
		nextHeight := currentBlock.Height + 1
		nextHash := nextBlock.BlockHash()
		nextHeader, e := c.GetBlockHeader(&nextHash)
		if e != nil {
			return e
		}
		_, e = c.filterBlock(nextBlock, nextHeight, true)
		if e != nil {
			return e
		}
		currentBlock.Height = nextHeight
		currentBlock.Hash = nextHash
		currentBlock.Timestamp = nextHeader.Timestamp
		blocksToNotify.Remove(blocksToNotify.Front())
	}
	c.bestBlockMtx.Lock()
	c.bestBlock = currentBlock
	c.bestBlockMtx.Unlock()
	return nil
}

// FilterBlocks scans the blocks contained in the FilterBlocksRequest for any addresses of interest. Each block will be
// fetched and filtered sequentially, returning a FilterBlocksResponse for the first block containing a matching
// address. If no matches are found in the range of blocks requested, the returned response will be nil.
//
// NOTE: This is part of the chainclient.Interface interface.
func (c *BitcoindClient) FilterBlocks(
	req *FilterBlocksRequest,
) (*FilterBlocksResponse, error) {
	blockFilterer := NewBlockFilterer(c.chainParams, req)
	// Iterate over the requested blocks, fetching each from the rpc client. Each block will scanned using the reverse
	// addresses indexes generated above, breaking out early if any addresses are found.
	for i, block := range req.Blocks {
		// TODO(conner): add prefetching, since we already know we'll be
		// fetching *every* block
		rawBlock, e := c.GetBlock(&block.Hash)
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
			BlockMeta:          block,
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

// rescan performs a rescan of the chain using a bitcoind backend, from the specified hash to the best known hash, while
// watching out for reorgs that happen during the rescan. It uses the addresses and outputs being tracked by the client
// in the watch list. This is called only within a queue processing loop.
func (c *BitcoindClient) rescan(start chainhash.Hash) (e error) {
	I.Ln("starting rescan from block", start)
	// We start by getting the best already processed block. We only use the height, as the hash can change during a
	// reorganization, which we catch by testing connectivity from known blocks to the previous block.
	bestHash, bestHeight, e := c.GetBestBlock()
	if e != nil {
		return e
	}
	bestHeader, e := c.GetBlockHeaderVerbose(bestHash)
	if e != nil {
		return e
	}
	bestBlock := am.BlockStamp{
		Hash:      *bestHash,
		Height:    bestHeight,
		Timestamp: time.Unix(bestHeader.Time, 0),
	}
	// Create a list of headers sorted in forward order. We'll use this in the event that we need to backtrack due to a
	// chain reorg.
	headers := list.New()
	previousHeader, e := c.GetBlockHeaderVerbose(&start)
	if e != nil {
		return e
	}
	previousHash, e := chainhash.NewHashFromStr(previousHeader.Hash)
	if e != nil {
		return e
	}
	headers.PushBack(previousHeader)
	// Queue a RescanFinished notification to the caller with the last block processed throughout the rescan once done.
	defer c.onRescanFinished(
		previousHash, previousHeader.Height,
		time.Unix(previousHeader.Time, 0),
	)
	// Cycle through all of the blocks known to bitcoind, being mindful of reorgs.
	for i := previousHeader.Height + 1; i <= bestBlock.Height; i++ {
		var hash *chainhash.Hash
		hash, e = c.GetBlockHash(int64(i))
		if e != nil {
			return e
		}
		// If the previous header is before the wallet birthday, fetch the current header and construct a dummy block,
		// rather than fetching the whole block itself. This speeds things up as we no longer have to fetch the whole
		// block when we know it won't match any of our filters.
		var block *wire.Block
		afterBirthday := previousHeader.Time >= c.birthday.Unix()
		if !afterBirthday {
			var header *wire.BlockHeader
			header, e = c.GetBlockHeader(hash)
			if e != nil {
				return e
			}
			block = &wire.Block{
				Header: *header,
			}
			afterBirthday = c.birthday.Before(header.Timestamp)
			if afterBirthday {
				c.onRescanProgress(
					previousHash, i,
					block.Header.Timestamp,
				)
			}
		}
		if afterBirthday {
			block, e = c.GetBlock(hash)
			if e != nil {
				return e
			}
		}
		for block.Header.PrevBlock.String() != previousHeader.Hash {
			// If we're in this for loop, it looks like we've been reorganized. We now walk backwards to the common
			// ancestor between the best chain and the known chain.
			//
			// First, we signal a disconnected block to rewind the rescan state.
			c.onBlockDisconnected(
				previousHash, previousHeader.Height,
				time.Unix(previousHeader.Time, 0),
			)
			// Get the previous block of the best chain.
			hash, e = c.GetBlockHash(int64(i - 1))
			if e != nil {
				return e
			}
			block, e = c.GetBlock(hash)
			if e != nil {
				return e
			}
			// Then, we'll the get the header of this previous block.
			if headers.Back() != nil {
				// If it's already in the headers list, we can just get it from there and remove the current hash.
				headers.Remove(headers.Back())
				if headers.Back() != nil {
					previousHeader = headers.Back().
						Value.(*btcjson.GetBlockHeaderVerboseResult)
					previousHash, e = chainhash.NewHashFromStr(
						previousHeader.Hash,
					)
					if e != nil {
						return e
					}
				}
			} else {
				// Otherwise, we get it from bitcoind.
				previousHash, e = chainhash.NewHashFromStr(
					previousHeader.PreviousHash,
				)
				if e != nil {
					return e
				}
				previousHeader, e = c.GetBlockHeaderVerbose(
					previousHash,
				)
				if e != nil {
					return e
				}
			}
		}
		// Now that we've ensured we haven't come across a reorg, we'll add the current block header to our list of
		// headers.
		blockHash := block.BlockHash()
		previousHash = &blockHash
		previousHeader = &btcjson.GetBlockHeaderVerboseResult{
			Hash:         blockHash.String(),
			Height:       i,
			PreviousHash: block.Header.PrevBlock.String(),
			Time:         block.Header.Timestamp.Unix(),
		}
		headers.PushBack(previousHeader)
		// Notify the block and any of its relevant transactions.
		if _, e = c.filterBlock(block, i, true); E.Chk(e) {
			return e
		}
		if i%10000 == 0 {
			c.onRescanProgress(
				previousHash, i, block.Header.Timestamp,
			)
		}
		// If we've reached the previously best known block, check to make sure the underlying node hasn't synchronized
		// additional blocks. If it has, update the best known block and continue to rescan to that point.
		if i == bestBlock.Height {
			bestHash, bestHeight, e = c.GetBestBlock()
			if e != nil {
				return e
			}
			bestHeader, e = c.GetBlockHeaderVerbose(bestHash)
			if e != nil {
				return e
			}
			bestBlock.Hash = *bestHash
			bestBlock.Height = bestHeight
			bestBlock.Timestamp = time.Unix(bestHeader.Time, 0)
		}
	}
	return nil
}

// filterBlock filters a block for watched outpoints and addresses, and returns any matching transactions, sending
// notifications along the way.
func (c *BitcoindClient) filterBlock(
	block *wire.Block, height int32,
	notify bool,
) ([]*tm.TxRecord, error) {
	// If this block happened before the client's birthday, then we'll skip it entirely.
	if block.Header.Timestamp.Before(c.birthday) {
		return nil, nil
	}
	if c.shouldNotifyBlocks() {
		D.F(
			"filtering block %d (%s) with %d transactions %s",
			height, block.BlockHash(), len(block.Transactions),
		)
	}
	// Create a block details template to use for all of the confirmed transactions found within this block.
	blockHash := block.BlockHash()
	blockDetails := &btcjson.BlockDetails{
		Hash:   blockHash.String(),
		Height: height,
		Time:   block.Header.Timestamp.Unix(),
	}
	// Now, we'll through all of the transactions in the block keeping track of any relevant to the caller.
	var relevantTxs []*tm.TxRecord
	confirmedTxs := make(map[chainhash.Hash]struct{})
	for i, tx := range block.Transactions {
		// Update the index in the block details with the index of this
		// transaction.
		blockDetails.Index = i
		isRelevant, rec, e := c.filterTx(tx, blockDetails, notify)
		if e != nil {
			W.F(
				"Unable to filter transaction %v: %v",
				tx.TxHash(), e,
			)
			continue
		}
		if isRelevant {
			relevantTxs = append(relevantTxs, rec)
			confirmedTxs[tx.TxHash()] = struct{}{}
		}
	}
	// Update the expiration map by setting the block's confirmed transactions and deleting any in the mempool that were
	// confirmed over 288 blocks ago.
	c.watchMtx.Lock()
	c.expiredMempool[height] = confirmedTxs
	if oldBlock, ok := c.expiredMempool[height-288]; ok {
		for txHash := range oldBlock {
			delete(c.mempool, txHash)
		}
		delete(c.expiredMempool, height-288)
	}
	c.watchMtx.Unlock()
	if notify {
		c.onFilteredBlockConnected(height, &block.Header, relevantTxs)
		c.onBlockConnected(&blockHash, height, block.Header.Timestamp)
	}
	return relevantTxs, nil
}

// filterTx determines whether a transaction is relevant to the client by inspecting the client's different filters.
func (c *BitcoindClient) filterTx(
	tx *wire.MsgTx,
	blockDetails *btcjson.BlockDetails,
	notify bool,
) (bool, *tm.TxRecord, error) {
	txDetails := util.NewTx(tx)
	if blockDetails != nil {
		txDetails.SetIndex(blockDetails.Index)
	}
	rec, e := tm.NewTxRecordFromMsgTx(txDetails.MsgTx(), time.Now())
	if e != nil {
		E.Ln(
			"Cannot create transaction record for relevant tx:", e,
		)
		return false, nil, e
	}
	if blockDetails != nil {
		rec.Received = time.Unix(blockDetails.Time, 0)
	}
	// We'll begin the filtering process by holding the lock to ensure we match exactly against what's currently in the
	// filters.
	c.watchMtx.Lock()
	defer c.watchMtx.Unlock()
	// If we've already seen this transaction and it's now been confirmed, then we'll shortcut the filter process by
	// immediately sending a notification to the caller that the filter matches.
	if _, ok := c.mempool[tx.TxHash()]; ok {
		if notify && blockDetails != nil {
			c.onRelevantTx(rec, blockDetails)
		}
		return true, rec, nil
	}
	// Otherwise, this is a new transaction we have yet to see. We'll need to determine if this transaction is somehow
	// relevant to the caller.
	var isRelevant bool
	// We'll start by cycling through its outputs to determine if it pays to any of the currently watched addresses. If
	// an output matches, we'll add it to our watch list.
	for i, out := range tx.TxOut {
		var addrs []btcaddr.Address
		_, addrs, _, e = txscript.ExtractPkScriptAddrs(
			out.PkScript, c.chainParams,
		)
		if e != nil {
			D.F(
				"Unable to parse output script in %s:%d: %v %s",
				tx.TxHash(), i, e,
			)
			continue
		}
		for _, addr := range addrs {
			if _, ok := c.watchedAddresses[addr.String()]; ok {
				isRelevant = true
				op := wire.OutPoint{
					Hash:  tx.TxHash(),
					Index: uint32(i),
				}
				c.watchedOutPoints[op] = struct{}{}
			}
		}
	}
	// If the transaction didn't pay to any of our watched addresses, we'll check if we're currently watching for the
	// hash of this transaction.
	if !isRelevant {
		if _, ok := c.watchedTxs[tx.TxHash()]; ok {
			isRelevant = true
		}
	}
	// If the transaction didn't pay to any of our watched hashes, we'll check if it spends any of our watched
	// outpoints.
	if !isRelevant {
		for _, in := range tx.TxIn {
			if _, ok := c.watchedOutPoints[in.PreviousOutPoint]; ok {
				isRelevant = true
				break
			}
		}
	}
	// If the transaction is not relevant to us, we can simply exit.
	if !isRelevant {
		return false, rec, nil
	}
	// Otherwise, the transaction matched our filters, so we should dispatch a notification for it. If it's still
	// unconfirmed, we'll include it in our mempool so that it can also be notified as part of FilteredBlockConnected
	// once it confirms.
	if blockDetails == nil {
		c.mempool[tx.TxHash()] = struct{}{}
	}
	c.onRelevantTx(rec, blockDetails)
	return true, rec, nil
}
