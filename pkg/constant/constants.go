package constant

import (
	"github.com/p9c/matrjoska/pkg/amt"
	"github.com/p9c/matrjoska/pkg/appdata"
	"github.com/p9c/matrjoska/pkg/blockchain"
	"github.com/p9c/matrjoska/pkg/peer"
	"time"
)

// A lotta constants that probably aren't being used
const (
	// defaultConfigFilename        = "conf.json"
	// defaultDataDirname           = "node"
	
	DefaultP2PPort               = 11047
	DefaultRPCListener           = "127.0.0.1"
	DefaultMaxPeers              = 23
	DefaultBanDuration           = time.Hour * 24
	DefaultBanThreshold          = 100
	DefaultMaxRPCClients         = 10
	DefaultMaxRPCWebsockets      = 25
	DefaultMaxRPCConcurrentReqs  = 20
	DefaultDbType                = "ffldb"
	DefaultFreeTxRelayLimit      = 15.0
	DefaultTrickleInterval       = peer.DefaultTrickleInterval
	DefaultBlockMaxSize          = 200000
	DefaultBlockMaxWeight        = 3000000
	BlockMaxSizeMin              = 1000
	BlockMaxSizeMax              = blockchain.MaxBlockBaseSize - 1000
	BlockMaxWeightMin            = 4000
	BlockMaxWeightMax            = blockchain.MaxBlockWeight - 4000
	DefaultMaxOrphanTransactions = 100
	DefaultSigCacheMaxSize       = 100000
	// DefaultBlockPrioritySize is the default size in bytes for high - priority / low-fee transactions. It is used to
	// help determine which are allowed into the mempool and consequently affects their relay and inclusion when
	// generating block templates.
	DefaultBlockPrioritySize = 50000
	// DefaultMinRelayTxFee is the minimum fee in satoshi that is required for a
	// transaction to be treated as free for relay and mining purposes. It is also
	// used to help determine if a transaction is considered dust and as a base for
	// calculating minimum required fees for larger transactions. This value is in
	// Satoshi/1000 bytes.
	DefaultMinRelayTxFee    = amt.Amount(1000)
	DefaultDataDirname      = "wallet"
	DefaultCAFilename       = "wallet.cert"
	DefaultConfigFilename   = "conf.json"
	DefaultLogLevel         = "info"
	DefaultLogDirname       = ""
	DefaultLogFilename      = "wallet/log"
	DefaultRPCMaxClients    = 10
	DefaultRPCMaxWebsockets = 25
	WalletDbName            = "wallet.db"
	// DbName is
	DbName = "wallet.db"
)

var DefaultHomeDir = appdata.Dir("pod", false)
