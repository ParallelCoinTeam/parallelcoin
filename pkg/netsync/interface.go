package netsync

import (
	"github.com/p9c/monorepo/node/mempool"
	"github.com/p9c/monorepo/pkg/blockchain"
	"github.com/p9c/monorepo/pkg/chaincfg"
	"github.com/p9c/monorepo/pkg/chainhash"
	"github.com/p9c/monorepo/pkg/peer"
	"github.com/p9c/monorepo/pkg/util"
	"github.com/p9c/monorepo/pkg/wire"
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
