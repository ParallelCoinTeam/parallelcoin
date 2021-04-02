package netsync

import (
	"github.com/p9c/duod/pkg/blockchain"
	"github.com/p9c/duod/pkg/chaincfg"
	"github.com/p9c/duod/pkg/chainhash"
	"github.com/p9c/duod/pkg/mempool"
	"github.com/p9c/duod/pkg/peer"
	"github.com/p9c/duod/pkg/util"
	"github.com/p9c/duod/pkg/wire"
)

// PeerNotifier exposes methods to notify peers of status changes to transactions, blocks, etc. Currently server (in the
// main package) implements this interface.
type PeerNotifier interface {
	AnnounceNewTransactions(newTxs []*mempool.TxDesc)
	UpdatePeerHeights(latestBlkHash *chainhash.Hash, latestHeight int32, updateSource *peer.Peer)
	RelayInventory(invVect *wire.InvVect, data interface{})
	TransactionConfirmed(tx *util.Tx)
}

// Config is a configuration struct used to initialize a new SyncManager.
type Config struct {
	PeerNotifier       PeerNotifier
	Chain              *blockchain.BlockChain
	TxMemPool          *mempool.TxPool
	ChainParams        *chaincfg.Params
	DisableCheckpoints bool
	MaxPeers           int
	FeeEstimator       *mempool.FeeEstimator
}
