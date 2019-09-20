package netsync

import (
	"github.com/parallelcointeam/parallelcoin/cmd/node/mempool"
	blockchain "github.com/parallelcointeam/parallelcoin/pkg/chain"
   `github.com/parallelcointeam/parallelcoin/pkg/chain/config/netparams`
   chainhash "github.com/parallelcointeam/parallelcoin/pkg/chain/hash"
	"github.com/parallelcointeam/parallelcoin/pkg/chain/wire"
	"github.com/parallelcointeam/parallelcoin/pkg/peer"
	"github.com/parallelcointeam/parallelcoin/pkg/util"
)

// PeerNotifier exposes methods to notify peers of status changes to transactions, blocks, etc. Currently server (in the main package) implements this interface.
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
	ChainParams        *netparams.Params
	DisableCheckpoints bool
	MaxPeers           int
	FeeEstimator       *mempool.FeeEstimator
}
