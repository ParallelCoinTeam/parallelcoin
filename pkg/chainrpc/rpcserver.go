package chainrpc

import (
	"bytes"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	js "encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/p9c/qu"

	"github.com/btcsuite/websocket"
	uberatomic "go.uber.org/atomic"

	"github.com/p9c/matrjoska/pkg/amt"
	"github.com/p9c/matrjoska/pkg/bits"
	"github.com/p9c/matrjoska/pkg/block"
	"github.com/p9c/matrjoska/pkg/btcaddr"
	"github.com/p9c/matrjoska/pkg/chaincfg"
	"github.com/p9c/matrjoska/pkg/fork"
	"github.com/p9c/matrjoska/pod/config"

	"github.com/p9c/matrjoska/cmd/node/active"
	"github.com/p9c/matrjoska/pkg/blockchain"
	"github.com/p9c/matrjoska/pkg/btcjson"
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/database"
	"github.com/p9c/matrjoska/pkg/indexers"
	"github.com/p9c/matrjoska/pkg/mempool"
	"github.com/p9c/matrjoska/pkg/mining"
	p "github.com/p9c/matrjoska/pkg/peer"
	"github.com/p9c/matrjoska/pkg/txscript"
	"github.com/p9c/matrjoska/pkg/util"
	"github.com/p9c/matrjoska/pkg/wire"
)

const (
	deserialfail    = "Failed to deserialize transaction"
	blockheightfail = "Failed to obtain block height"
)

type CommandHandler struct {
	Fn     func(*Server, interface{}, qu.C) (interface{}, error)
	Call   chan API
	Result func() API
}

// GBTWorkState houses state that is used in between multiple RPC invocations to getblocktemplate.
type GBTWorkState struct {
	sync.Mutex
	LastTxUpdate  time.Time
	LastGenerated time.Time
	prevHash      *chainhash.Hash
	MinTimestamp  time.Time
	Template      *mining.BlockTemplate
	NotifyMap     map[chainhash.Hash]map[int64]qu.C
	TimeSource    blockchain.MedianTimeSource
	Algo          string
	StateCfg      *active.Config
	Config        *config.Config
}

// ParsedRPCCmd represents a JSON-RPC request object that has been parsed into a known concrete command along with any
// error that might have happened while parsing it.
type ParsedRPCCmd struct {
	ID     interface{}
	Method string
	Cmd    interface{}
	Err    *btcjson.RPCError
}

// RetrievedTx represents a transaction that was either loaded from the transaction memory pool or from the database.
//
// When a transaction is loaded from the database, it is loaded with the raw serialized bytes while the mempool has the
// fully deserialized structure.
//
// This structure therefore will have one of the two fields set depending on where is was retrieved from.
//
// This is mainly done for efficiency to avoid extra serialization steps when possible.
type RetrievedTx struct {
	TxBytes []byte
	BlkHash *chainhash.Hash // Only set when transaction is in a block.
	Tx      *util.Tx
}

// Server provides a concurrent safe RPC Server to a chain Server.
type Server struct {
	Cfg                             ServerConfig
	StateCfg                        *active.Config
	Config                          *config.Config
	NtfnMgr                         *WSNtfnMgr
	StatusLines                     map[int]string
	StatusLock                      sync.RWMutex
	WG                              sync.WaitGroup
	GBTWorkState                    *GBTWorkState
	HelpCacher                      *HelpCacher
	RequestProcessShutdown          qu.C
	Quit                            qu.C
	Started                         int32
	Shutdown                        int32
	NumClients                      int32
	AuthSHA                         [sha256.Size]byte
	LimitAuthSHA                    [sha256.Size]byte
	StartController, StopController qu.C
}

// ServerConfig is a descriptor containing the RPC Server configuration.
type ServerConfig struct {
	// Cx passes through the context variable for setting up a Server
	Cfg *config.Config
	// Listeners defines a slice of listeners for which the RPC Server will take ownership of and accept connections.
	//
	// Since the RPC Server takes ownership of these listeners, they will be closed when the RPC Server is stopped.
	Listeners []net.Listener
	// StartupTime is the unix timestamp for when the Server that is hosting the RPC Server started.
	StartupTime int64
	// ConnMgr defines the connection manager for the RPC Server to use.
	//
	// It provides the RPC Server with a means to do things such as add, remove, connect, disconnect, and query peers as
	// well as other connection-related data and tasks.
	ConnMgr ServerConnManager
	// SyncMgr defines the sync manager for the RPC Server to use.
	SyncMgr ServerSyncManager
	// These fields allow the RPC Server to interface with the local block chain data and state.
	TimeSource  blockchain.MedianTimeSource
	Chain       *blockchain.BlockChain
	ChainParams *chaincfg.Params
	DB          database.DB
	// TxMemPool defines the transaction memory pool to interact with.
	TxMemPool *mempool.TxPool
	// These fields allow the RPC Server to interface with mining.
	//
	// Generator produces block templates and the CPUMiner solves them using the CPU.
	//
	// CPU mining is typically only useful for test purposes when doing regression or simulation testing.
	// until there was divhash
	Generator *mining.BlkTmplGenerator
	// CPUMiner  *cpuminer.CPUMiner
	//
	// These fields define any optional indexes the RPC Server can make use of to provide additional data when queried.
	TxIndex   *indexers.TxIndex
	AddrIndex *indexers.AddrIndex
	CfIndex   *indexers.CFIndex
	// The fee estimator keeps track of how long transactions are left in the mempool before they are mined into blocks.
	FeeEstimator *mempool.FeeEstimator
	// Algo sets the algorithm expected from the RPC endpoint. This allows multiple ports to serve multiple types of
	// miners with one main node per algorithm. Currently 514 for Scrypt and anything else passes for SHA256d.
	Algo string
	// CPUMiner *exec.Cmd
	Hashrate                        uberatomic.Uint64
	Quit                            qu.C
	StartController, StopController qu.C
}

// ServerConnManager represents a connection manager for use with the RPC Server. The interface contract requires that
// all of these methods are safe for concurrent access.
type ServerConnManager interface {
	// Connect adds the provided address as a new outbound peer. The permanent flag indicates whether or not to make the
	// peer persistent and reconnect if the connection is lost. Attempting to connect to an already existing peer will
	// return an error.
	Connect(addr string, permanent bool) error
	// RemoveByID removes the peer associated with the provided id from the list of persistent peers.
	//
	// Attempting to remove an id that does not exist will return an error.
	RemoveByID(id int32) error
	// RemoveByAddr removes the peer associated with the provided address from the list of persistent peers.
	//
	// Attempting to remove an address that does not exist will return an error.
	RemoveByAddr(addr string) error
	// DisconnectByID disconnects the peer associated with the provided id. This applies to both inbound and outbound
	// peers.
	//
	// Attempting to remove an id that does not exist will return an error.
	DisconnectByID(id int32) error
	// DisconnectByAddr disconnects the peer associated with the provided address. This applies to both inbound and
	// outbound peers.
	//
	// Attempting to remove an address that does not exist will return an error.
	DisconnectByAddr(addr string) error
	// ConnectedCount returns the number of currently connected peers.
	ConnectedCount() int32
	// NetTotals returns the sum of all bytes received and sent across the network for all peers.
	NetTotals() (uint64, uint64)
	// ConnectedPeers returns an array consisting of all connected peers.
	ConnectedPeers() []ServerPeer
	// PersistentPeers returns an array consisting of all the persistent peers.
	PersistentPeers() []ServerPeer
	// BroadcastMessage sends the provided message to all currently connected
	// peers.
	BroadcastMessage(msg wire.Message)
	// AddRebroadcastInventory adds the provided inventory to the list of inventories to be rebroadcast at random
	// intervals until they show up in a block.
	AddRebroadcastInventory(iv *wire.InvVect, data interface{})
	// RelayTransactions generates and relays inventory vectors for all of the passed transactions to all connected
	// peers.
	RelayTransactions(txns []*mempool.TxDesc)
}

// ServerPeer represents a peer for use with the RPC Server.
//
// The interface contract requires that all of these methods are safe for concurrent access.
type ServerPeer interface {
	// ToPeer returns the underlying peer instance.
	ToPeer() *p.Peer
	// IsTxRelayDisabled returns whether or not the peer has disabled transaction relay.
	IsTxRelayDisabled() bool
	// GetBanScore returns the current integer value that represents how close the peer is to being banned.
	GetBanScore() uint32
	// GetFeeFilter returns the requested current minimum fee rate for which transactions should be announced.
	GetFeeFilter() int64
}

// ServerSyncManager represents a sync manager for use with the RPC Server.
//
// The interface contract requires that all of these methods are safe for concurrent access.
type ServerSyncManager interface {
	// IsCurrent returns whether or not the sync manager believes the chain is current as compared to the rest of the
	// network.
	IsCurrent() bool
	// SubmitBlock submits the provided block to the network after processing it locally.
	SubmitBlock(block *block.Block, flags blockchain.BehaviorFlags) (
		bool, error,
	)
	// Pause pauses the sync manager until the returned channel is closed.
	Pause() chan<- struct{}
	// SyncPeerID returns the ID of the peer that is currently the peer being used to sync from or 0 if there is none.
	SyncPeerID() int32
	// LocateHeaders returns the headers of the blocks after the first known block in the provided locators until the
	// provided stop hash or the current tip is reached, up to a max of wire.MaxBlockHeadersPerMsg hashes.
	LocateHeaders(
		locators []*chainhash.Hash, hashStop *chainhash.Hash,
	) []wire.BlockHeader
}

// API version constants
const (
	JSONRPCSemverString = "1.3.0"
	JSONRPCSemverMajor  = 1
	JSONRPCSemverMinor  = 3
	JSONRPCSemverPatch  = 0
	// RPCAuthTimeoutSeconds is the number of seconds a connection to the RPC Server is allowed to stay open without
	// authenticating before it is closed.
	RPCAuthTimeoutSeconds = 10
	// GBTNonceRange is two 32-bit big-endian hexadecimal integers which represent the valid ranges of nonces returned
	// by the getblocktemplate RPC.
	GBTNonceRange = "00000000ffffffff"
	// GBTRegenerateSeconds is the number of seconds that must pass before a new template is generated when the previous
	// block hash has not changed and there have been changes to the available transactions in the memory pool.
	GBTRegenerateSeconds = 60
	// MaxProtocolVersion is the max protocol version the Server supports.
	MaxProtocolVersion = 70002
)

func getIfc() chan interface{} {
	return make(chan interface{})
}

var (
	// ErrRPCNoWallet is an error returned to RPC clients when the provided command is recognized as a wallet command.
	ErrRPCNoWallet = &btcjson.RPCError{
		Code:    btcjson.ErrRPCNoWallet,
		Message: "This implementation does not implement wallet commands",
	}
	// ErrRPCUnimplemented is an error returned to RPC clients when the provided command is recognized, but not
	// implemented.
	ErrRPCUnimplemented = &btcjson.RPCError{
		Code:    btcjson.ErrRPCUnimplemented,
		Message: "Command unimplemented",
	}
	// GBTCapabilities describes additional capabilities returned with a block template generated by the
	// getblocktemplate RPC.
	//
	// It is declared here to avoid the overhead of creating the slice on every invocation for constant data.
	GBTCapabilities = []string{"proposal"}
	// GBTCoinbaseAux describes additional data that miners should include in the coinbase signature script.
	//
	// It is declared here to avoid the overhead of creating a new object on every invocation for constant data.
	GBTCoinbaseAux = &btcjson.GetBlockTemplateResultAux{
		Flags: hex.EncodeToString(
			BuilderScript(
				txscript.
					NewScriptBuilder().
					AddData([]byte(mining.CoinbaseFlags)),
			),
		),
	}
	// GBTMutableFields are the manipulations the Server allows to be made to block templates generated by the
	// getblocktemplate RPC.
	//
	// It is declared here to avoid the overhead of creating the slice on every invocation for constant data.
	GBTMutableFields = []string{
		"time", "transactions/add", "prevblock", "coinbase/append",
	}

	// RPCAskWallet is list of commands that we recognize, but for which pod has no support because it lacks support for
	// wallet functionality. For these commands the user should ask a connected instance of the wallet.
	RPCAskWallet = map[string]CommandHandler{
		"addmultisigaddress":     {},
		"backupwallet":           {},
		"createencryptedwallet":  {},
		"createmultisig":         {},
		"dumpprivkey":            {},
		"dumpwallet":             {},
		"dropwallethistory":      {},
		"encryptwallet":          {},
		"getaccount":             {},
		"getaccountaddress":      {},
		"getaddressesbyaccount":  {},
		"getbalance":             {},
		"getnewaddress":          {},
		"getrawchangeaddress":    {},
		"getreceivedbyaccount":   {},
		"getreceivedbyaddress":   {},
		"gettransaction":         {},
		"gettxoutsetinfo":        {},
		"getunconfirmedbalance":  {},
		"getwalletinfo":          {},
		"importprivkey":          {},
		"importwallet":           {},
		"keypoolrefill":          {},
		"listaccounts":           {},
		"listaddressgroupings":   {},
		"listlockunspent":        {},
		"listreceivedbyaccount":  {},
		"listreceivedbyaddress":  {},
		"listsinceblock":         {},
		"listtransactions":       {},
		"listunspent":            {},
		"lockunspent":            {},
		"move":                   {},
		"sendfrom":               {},
		"sendmany":               {},
		"sendtoaddress":          {},
		"setaccount":             {},
		"settxfee":               {},
		"signmessage":            {},
		"signrawtransaction":     {},
		"walletlock":             {},
		"walletpassphrase":       {},
		"walletpassphrasechange": {},
	}

	// RPCHandlers maps RPC command strings to appropriate handler functions.
	//
	// This is set by init because help references RPCHandlers and thus causes a dependency loop.
	RPCHandlers map[string]CommandHandler

	// RPCLimited RPCHandlersBeforeInit is
	//
	// RPCHandlersBeforeInit = map[string]CommandHandler{
	// 	"addnode": {
	// 		HandleAddNode, make(chan API), func() API {
	// 			return API{
	// 				btcjson.AddNodeCmd{},
	// 				make(chan AddNodeRes),
	// 			}
	// 		},
	// 	},
	// 	"createrawtransaction": {
	// 		HandleCreateRawTransaction, make(chan API),
	// 		func() API {
	// 			return API{btcjson.CreateRawTransactionCmd{},
	// 				make(chan CreateRawTransactionRes)}
	// 		},
	// 	},
	// 	// "debuglevel":            handleDebugLevel,
	// 	"decoderawtransaction": {
	// 		HandleDecodeRawTransaction, make(chan API),
	// 		func() API {
	// 			return API{btcjson.DecodeRawTransactionCmd{},
	// 				make(chan DecodeRawTransactionRes)}
	// 		},
	// 	},
	// 	"decodescript": {
	// 		HandleDecodeScript, make(chan API),
	// 		func() API {
	// 			return API{btcjson.DecodeScriptCmd{},
	// 				make(chan DecodeScriptRes)}
	// 		},
	// 	},
	// 	"estimatefee": {
	// 		HandleEstimateFee, make(chan API),
	// 		func() API {
	// 			return API{btcjson.EstimateFeeCmd{},
	// 				make(chan EstimateFeeRes)}
	// 		},
	// 	},
	// 	"generate": {
	// 		HandleGenerate, make(chan API),
	// 		func() API { return API{nil, make(chan GenerateRes)} },
	// 	},
	// 	"getaddednodeinfo": {
	// 		HandleGetAddedNodeInfo, make(chan API),
	// 		func() API {
	// 			return API{btcjson.GetAddedNodeInfoCmd{},
	// 				make(chan GetAddedNodeInfoRes)}
	// 		},
	// 	},
	// 	"getbestblock": {
	// 		HandleGetBestBlock, make(chan API),
	// 		func() API { return API{nil, make(chan GetBestBlockRes)} },
	// 	},
	// 	"getbestblockhash": {
	// 		HandleGetBestBlockHash, make(chan API),
	// 		func() API { return API{nil, make(chan GetBestBlockHashRes)} },
	// 	},
	// 	"getblock": {
	// 		HandleGetBlock, make(chan API),
	// 		func() API {
	// 			return API{btcjson.GetBlockCmd{},
	// 				make(chan GetBlockRes)}
	// 		},
	// 	},
	// 	"getblockchaininfo": {
	// 		HandleGetBlockChainInfo, make(chan API),
	// 		func() API { return API{nil, make(chan GetBlockChainInfoRes)} },
	// 	},
	// 	"getblockcount": {
	// 		HandleGetBlockCount, make(chan API),
	// 		func() API { return API{nil, make(chan GetBlockCountRes)} },
	// 	},
	// 	"getblockhash": {
	// 		HandleGetBlockHash, make(chan API),
	// 		func() API {
	// 			return API{btcjson.GetBlockHashCmd{},
	// 				make(chan GetBlockHashRes)}
	// 		},
	// 	},
	// 	"getblockheader": {
	// 		HandleGetBlockHeader, make(chan API),
	// 		func() API {
	// 			return API{btcjson.GetBlockHeaderCmd{},
	// 				make(chan GetBlockHeaderRes)}
	// 		},
	// 	},
	// 	"getblocktemplate": {
	// 		HandleGetBlockTemplate, make(chan API),
	// 		func() API {
	// 			return API{btcjson.GetBlockTemplateCmd{},
	// 				make(chan GetBlockTemplateRes)}
	// 		},
	// 	},
	// 	"getcfilter": {
	// 		HandleGetCFilter, make(chan API),
	// 		func() API {
	// 			return API{btcjson.GetCFilterCmd{},
	// 				make(chan GetCFilterRes)}
	// 		},
	// 	},
	// 	"getcfilterheader": {
	// 		HandleGetCFilterHeader, make(chan API),
	// 		func() API {
	// 			return API{btcjson.GetCFilterHeaderCmd{},
	// 				make(chan GetCFilterHeaderRes)}
	// 		},
	// 	},
	// 	"getconnectioncount": {
	// 		HandleGetConnectionCount, make(chan API),
	// 		func() API { return API{nil, make(chan GetConnectionCountRes)} },
	// 	},
	// 	"getcurrentnet": {
	// 		HandleGetCurrentNet, make(chan API),
	// 		func() API { return API{nil, make(chan GetCurrentNetRes)} },
	// 	},
	// 	"getdifficulty": {
	// 		HandleGetDifficulty, make(chan API),
	// 		func() API {
	// 			return API{btcjson.GetDifficultyCmd{},
	// 				make(chan GetDifficultyRes)}
	// 		},
	// 	},
	// 	"getgenerate": {
	// 		HandleGetGenerate, make(chan API),
	// 		func() API {
	// 			return API{btcjson.GetHeadersCmd{},
	// 				make(chan GetGenerateRes)}
	// 		},
	// 	},
	// 	"gethashespersec": {
	// 		HandleGetHashesPerSec, make(chan API),
	// 		func() API { return API{nil, make(chan GetHashesPerSecRes)} },
	// 	},
	// 	"getheaders": {
	// 		HandleGetHeaders, make(chan API),
	// 		func() API {
	// 			return API{btcjson.GetHeadersCmd{},
	// 				make(chan GetHeadersRes)}
	// 		},
	// 	},
	// 	"getinfo": {
	// 		HandleGetInfo, make(chan API),
	// 		func() API { return API{nil, make(chan GetInfoRes)} },
	// 	},
	// 	"getmempoolinfo": {
	// 		HandleGetMempoolInfo, make(chan API),
	// 		func() API { return API{nil, make(chan GetMempoolInfoRes)} },
	// 	},
	// 	"getmininginfo": {
	// 		HandleGetMiningInfo, make(chan API),
	// 		func() API { return API{nil, make(chan GetMiningInfoRes)} },
	// 	},
	// 	"getnettotals": {
	// 		HandleGetNetTotals, make(chan API),
	// 		func() API { return API{nil, make(chan GetNetTotalsRes)} },
	// 	},
	// 	"getnetworkhashps": {
	// 		HandleGetNetworkHashPS, make(chan API),
	// 		func() API {
	// 			return API{btcjson.GetNetworkHashPSCmd{},
	// 				make(chan GetNetworkHashPSRes)}
	// 		},
	// 	},
	// 	"getpeerinfo": {
	// 		HandleGetPeerInfo, make(chan API),
	// 		func() API { return API{nil, make(chan GetPeerInfoRes)} },
	// 	},
	// 	"getrawmempool": {
	// 		HandleGetRawMempool, make(chan API),
	// 		func() API {
	// 			return API{btcjson.GetRawMempoolCmd{},
	// 				make(chan GetRawMempoolRes)}
	// 		},
	// 	},
	// 	"getrawtransaction": {
	// 		HandleGetRawTransaction, make(chan API),
	// 		func() API {
	// 			return API{btcjson.GetRawTransactionCmd{},
	// 				make(chan GetRawTransactionRes)}
	// 		},
	// 	},
	// 	"gettxout": {
	// 		HandleGetTxOut, make(chan API),
	// 		func() API {
	// 			return API{btcjson.GetTxOutCmd{},
	// 				make(chan GetTxOutRes)}
	// 		},
	// 	},
	// 	// "getwork":               HandleGetWork,
	// 	"help": {
	// 		HandleHelp, make(chan API),
	// 		func() API {
	// 			return API{btcjson.HelpCmd{},
	// 				make(chan HelpRes)}
	// 		},
	// 	},
	// 	"node": {
	// 		HandleNode, make(chan API),
	// 		func() API {
	// 			return API{btcjson.NodeCmd{},
	// 				make(chan NodeRes)}
	// 		},
	// 	},
	// 	"ping": {
	// 		HandlePing, make(chan API),
	// 		func() API { return API{nil, make(chan PingRes)} },
	// 	},
	// 	"searchrawtransactions": {
	// 		HandleSearchRawTransactions, make(chan API),
	// 		func() API {
	// 			return API{btcjson.SearchRawTransactionsCmd{},
	// 				make(chan SearchRawTransactionsRes)}
	// 		},
	// 	},
	// 	"sendrawtransaction": {
	// 		HandleSendRawTransaction, make(chan API),
	// 		func() API {
	// 			return API{btcjson.SendRawTransactionCmd{},
	// 				make(chan SendRawTransactionRes)}
	// 		},
	// 	},
	// 	"setgenerate": {
	// 		HandleSetGenerate, make(chan API),
	// 		func() API {
	// 			return API{btcjson.SetGenerateCmd{},
	// 				make(chan SetGenerateRes)}
	// 		},
	// 	},
	// 	"stop": {
	// 		HandleStop, make(chan API),
	// 		func() API { return API{nil, make(chan StopRes)} },
	// 	},
	// 	"restart": {
	// 		HandleRestart, make(chan API),
	// 		func() API { return API{nil, make(chan RestartRes)} },
	// 	},
	// 	"resetchain": {
	// 		HandleResetChain, make(chan API),
	// 		func() API { return API{nil, make(chan ResetChainRes)} },
	// 	},
	// 	// "dropwallethistory":     HandleDropWalletHistory,
	// 	"submitblock": {
	// 		HandleSubmitBlock, make(chan API),
	// 		func() API {
	// 			return API{btcjson.SubmitBlockCmd{},
	// 				make(chan SubmitBlockRes)}
	// 		},
	// 	},
	// 	"uptime": {
	// 		HandleUptime, make(chan API),
	// 		func() API { return API{nil, make(chan UptimeRes)} },
	// 	},
	// 	"validateaddress": {
	// 		HandleValidateAddress, make(chan API),
	// 		func() API {
	// 			return API{btcjson.ValidateAddressCmd{},
	// 				make(chan ValidateAddressRes)}
	// 		},
	// 	},
	// 	"verifychain": {
	// 		HandleVerifyChain, make(chan API),
	// 		func() API {
	// 			return API{btcjson.VerifyChainCmd{},
	// 				make(chan VerifyChainRes)}
	// 		},
	// 	},
	// 	"verifymessage": {
	// 		HandleVerifyMessage, make(chan API),
	// 		func() API {
	// 			return API{btcjson.VerifyMessageCmd{},
	// 				make(chan VerifyMessageRes)}
	// 		},
	// 	},
	// 	"version": {
	// 		HandleVersion, make(chan API),
	// 		func() API {
	// 			return API{btcjson.VersionCmd{},
	// 				make(chan VersionRes)}
	// 		},
	// 	},
	// }

	// RPCLimited isCommands that are available to a limited user
	RPCLimited = map[string]CommandHandler{
		// Websockets commands
		"loadtxfilter":          {},
		"notifyblocks":          {},
		"notifynewtransactions": {},
		"notifyreceived":        {},
		"notifyspent":           {},
		"rescan":                {},
		"rescanblocks":          {},
		"session":               {},
		// Websockets AND HTTP/S commands
		"help": {},
		// HTTP/S-only commands
		"createrawtransaction":  {},
		"decoderawtransaction":  {},
		"decodescript":          {},
		"estimatefee":           {},
		"getbestblock":          {},
		"getbestblockhash":      {},
		"getblock":              {},
		"getblockcount":         {},
		"getblockhash":          {},
		"getblockheader":        {},
		"getcfilter":            {},
		"getcfilterheader":      {},
		"getcurrentnet":         {},
		"getdifficulty":         {},
		"getheaders":            {},
		"getinfo":               {},
		"getnettotals":          {},
		"getnetworkhashps":      {},
		"getrawmempool":         {},
		"getrawtransaction":     {},
		"gettxout":              {},
		"searchrawtransactions": {},
		"sendrawtransaction":    {},
		"submitblock":           {},
		"uptime":                {},
		"validateaddress":       {},
		"verifymessage":         {},
		"version":               {},
	}
	// RPCUnimplemented is commands that are currently unimplemented, but should ultimately be.
	RPCUnimplemented = map[string]struct{}{
		"estimatepriority": {},
		"getchaintips":     {},
		"getmempoolentry":  {},
		"getnetworkinfo":   {},
		"getwork":          {},
		"invalidateblock":  {},
		"preciousblock":    {},
		"reconsiderblock":  {},
	}
)

// NotifyBlockConnected uses the newly-connected block to notify any long poll clients with a new block template when
// their existing block template is stale due to the newly connected block.
func (state *GBTWorkState) NotifyBlockConnected(blockHash *chainhash.Hash) {
	go func() {
		state.Lock()
		statelasttxupdate := state.LastTxUpdate
		state.Unlock()
		state.NotifyLongPollers(blockHash, statelasttxupdate)
	}()
}

// NotifyMempoolTx uses the new last updated time for the transaction memory pool to notify any long poll clients with a
// new block template when their existing block template is stale due to enough time passing and the contents of the
// memory pool changing.
func (state *GBTWorkState) NotifyMempoolTx(lastUpdated time.Time) {
	go func() {
		state.Lock()
		defer state.Unlock()
		// No need to notify anything if no block templates have been generated yet.
		if state.prevHash == nil || state.LastGenerated.IsZero() {
			return
		}
		if time.Now().After(state.LastGenerated.Add(time.Second * GBTRegenerateSeconds)) {
			state.NotifyLongPollers(state.prevHash, lastUpdated)
		}
	}()
}

// BlockTemplateResult returns the current block template associated with the state as a json.GetBlockTemplateResult
// that is ready to be encoded to JSON and returned to the caller.
//
// This function MUST be called with the state locked.
func (state *GBTWorkState) BlockTemplateResult(
	useCoinbaseValue bool, submitOld *bool,
) (
	*btcjson.GetBlockTemplateResult,
	error,
) {
	// Ensure the timestamps are still in valid range for the template.
	//
	// This should really only ever happen if the local clock is changed after the template is generated, but it's
	// important to avoid serving invalid block templates.
	template := state.Template
	msgBlock := template.Block
	header := &msgBlock.Header
	adjustedTime := state.TimeSource.AdjustedTime()
	maxTime := adjustedTime.Add(time.Second * blockchain.MaxTimeOffsetSeconds)
	if header.Timestamp.After(maxTime) {
		return nil, &btcjson.RPCError{
			Code: btcjson.ErrRPCOutOfRange,
			Message: fmt.Sprintf(
				"The template time is after the "+
					"maximum allowed time for a block - template "+
					"time %v, maximum time %v", adjustedTime,
				maxTime,
			),
		}
	}
	// Convert each transaction in the block template to a template result transaction. The result does not include the
	// coinbase, so notice the adjustments to the various lengths and indices.
	numTx := len(msgBlock.Transactions)
	transactions := make([]btcjson.GetBlockTemplateResultTx, 0, numTx-1)
	txIndex := make(map[chainhash.Hash]int64, numTx)
	for i, tx := range msgBlock.Transactions {
		txHash := tx.TxHash()
		txIndex[txHash] = int64(i)
		// Skip the coinbase transaction.
		if i == 0 {
			continue
		}
		// Create an array of 1-based indices to transactions that come before this one in the transactions list which
		// this one depends on.
		//
		// This is necessary since the created block must ensure proper ordering of the dependencies.
		//
		// A map is used before creating the final array to prevent duplicate entries when multiple inputs reference the
		// same transaction.
		dependsMap := make(map[int64]struct{})
		for _, txIn := range tx.TxIn {
			if idx, ok := txIndex[txIn.PreviousOutPoint.Hash]; ok {
				dependsMap[idx] = struct{}{}
			}
		}
		depends := make([]int64, 0, len(dependsMap))
		for idx := range dependsMap {
			depends = append(depends, idx)
		}
		// Serialize the transaction for later conversion to hex.
		txBuf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
		if e := tx.Serialize(txBuf); E.Chk(e) {
			context := "Failed to serialize transaction"
			return nil, InternalRPCError(e.Error(), context)
		}
		bTx := util.NewTx(tx)
		resultTx := btcjson.GetBlockTemplateResultTx{
			Data:    hex.EncodeToString(txBuf.Bytes()),
			Hash:    txHash.String(),
			Depends: depends,
			Fee:     template.Fees[i],
			SigOps:  template.SigOpCosts[i],
			Weight:  blockchain.GetTransactionWeight(bTx),
		}
		transactions = append(transactions, resultTx)
	}
	// Generate the block template reply.
	//
	// Note that following mutations are implied by the included or omission of fields:
	//
	//   Including MinTime -> time/ decrement
	//
	//   Omitting CoinbaseTxn -> coinbase, generation
	targetDifficulty := fmt.Sprintf("%064x", bits.CompactToBig(header.Bits))
	templateID := EncodeTemplateID(state.prevHash, state.LastGenerated)
	reply := btcjson.GetBlockTemplateResult{
		Bits:         strconv.FormatInt(int64(header.Bits), 16),
		CurTime:      header.Timestamp.Unix(),
		Height:       int64(template.Height),
		PreviousHash: header.PrevBlock.String(),
		WeightLimit:  blockchain.MaxBlockWeight,
		SigOpLimit:   blockchain.MaxBlockSigOpsCost,
		SizeLimit:    wire.MaxBlockPayload,
		Transactions: transactions,
		Version:      header.Version,
		LongPollID:   templateID,
		SubmitOld:    submitOld,
		Target:       targetDifficulty,
		MinTime:      state.MinTimestamp.Unix(),
		MaxTime:      maxTime.Unix(),
		Mutable:      GBTMutableFields,
		NonceRange:   GBTNonceRange,
		Capabilities: GBTCapabilities,
	}
	// // If the generated block template includes transactions with witness data, then
	// // include the witness commitment in the GBT result.
	// if template.WitnessCommitment != nil {
	// 	reply.DefaultWitnessCommitment = hex.EncodeToString(template.WitnessCommitment)
	// }
	if useCoinbaseValue {
		reply.CoinbaseAux = GBTCoinbaseAux
		reply.CoinbaseValue = &msgBlock.Transactions[0].TxOut[0].Value
	} else {
		// Ensure the template has a valid payment address associated with it when a full coinbase is requested.
		if !template.ValidPayAddress {
			return nil, &btcjson.RPCError{
				Code: btcjson.ErrRPCInternal.Code,
				Message: "A coinbase transaction has been " +
					"requested, but the Server has not " +
					"been configured with any payment " +
					"addresses via --miningaddr",
			}
		}
		// Serialize the transaction for conversion to hex.
		tx := msgBlock.Transactions[0]
		txBuf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
		if e := tx.Serialize(txBuf); E.Chk(e) {
			context := "Failed to serialize transaction"
			return nil, InternalRPCError(e.Error(), context)
		}
		resultTx := btcjson.GetBlockTemplateResultTx{
			Data:    hex.EncodeToString(txBuf.Bytes()),
			Hash:    tx.TxHash().String(),
			Depends: []int64{},
			Fee:     template.Fees[0],
			SigOps:  template.SigOpCosts[0],
		}
		reply.CoinbaseTxn = &resultTx
	}
	return &reply, nil
}

// NotifyLongPollers notifies any channels that have been registered to be notified when block templates are stale.
//
// This function MUST be called with the state locked.
func (state *GBTWorkState) NotifyLongPollers(
	latestHash *chainhash.Hash,
	lastGenerated time.Time,
) {
	// Notify anything that is waiting for a block template update from a hash which is not the hash of the tip of the
	// best chain since their work is now invalid.
	for hash, channels := range state.NotifyMap {
		if !hash.IsEqual(latestHash) {
			for _, c := range channels {
				c.Q()
			}
			delete(state.NotifyMap, hash)
		}
	}
	// Return now if the provided last generated timestamp has not been initialized.
	if lastGenerated.IsZero() {
		return
	}
	// Return now if there is nothing registered for updates to the current best block hash.
	channels, ok := state.NotifyMap[*latestHash]
	if !ok {
		return
	}
	// Notify anything that is waiting for a block template update from a block template generated before the most
	// recently generated block template.
	lastGeneratedUnix := lastGenerated.Unix()
	for lastGen, c := range channels {
		if lastGen < lastGeneratedUnix {
			c.Q()
			delete(channels, lastGen)
		}
	}
	// Remove the entry altogether if there are no more registered channels.
	if len(channels) == 0 {
		delete(state.NotifyMap, *latestHash)
	}
}

// TemplateUpdateChan returns a channel that will be closed once the block template associated with the passed previous
// hash and last generated time is stale.
//
// The function will return existing channels for duplicate parameters which allows to wait for the same block template
// without requiring a different channel for each client.
//
// This function MUST be called with the state locked.
func (state *GBTWorkState) TemplateUpdateChan(
	prevHash *chainhash.Hash, lastGenerated int64,
) qu.C {
	// Either get the current list of channels waiting for updates about changes to block template for the previous hash
	// or create a new one.
	channels, ok := state.NotifyMap[*prevHash]
	if !ok {
		m := make(map[int64]qu.C)
		state.NotifyMap[*prevHash] = m
		channels = m
	}
	// Get the current channel associated with the time the block template was last generated or create a new one.
	c, ok := channels[lastGenerated]
	if !ok {
		c = qu.T()
		channels[lastGenerated] = c
	}
	return c
}

// UpdateBlockTemplate creates or updates a block template for the work state.
//
// A new block template will be generated when the current best block has changed or the transactions in the memory pool
// have been updated and it has been long enough since the last template was generated.
//
// Otherwise, the timestamp for the existing block template is updated (and possibly the difficulty on testnet per the
// consensus rules).
//
// Finally, if the useCoinbaseValue flag is false and the existing block template does not already contain a valid
// payment address, the block template will be updated with a randomly selected payment address from the list of
// configured addresses.
//
// This function MUST be called with the state locked.
func (state *GBTWorkState) UpdateBlockTemplate(
	s *Server,
	useCoinbaseValue bool,
) (e error) {
	generator := s.Cfg.Generator
	lastTxUpdate := generator.GetTxSource().LastUpdated()
	if lastTxUpdate.IsZero() {
		lastTxUpdate = time.Now()
	}
	// Generate a new block template when the current best block has changed or the transactions in the memory pool have
	// been updated and it has been at least gbtRegenerateSecond since the last template was generated.
	var msgBlock *wire.Block
	var targetDifficulty string
	latestHash := &s.Cfg.Chain.BestSnapshot().Hash
	template := state.Template
	if template == nil || state.prevHash == nil ||
		!state.prevHash.IsEqual(latestHash) ||
		(state.LastTxUpdate != lastTxUpdate &&
			time.Now().After(
				state.LastGenerated.Add(
					time.Second*
						GBTRegenerateSeconds,
				),
			)) {
		// Reset the previous best hash the block template was generated against so any errors below cause the next
		// invocation to try again.
		state.prevHash = nil
		// Choose a payment address at random if the caller requests a full coinbase as opposed to only the pertinent
		// details needed to create their own coinbase.
		var payAddr btcaddr.Address
		if !useCoinbaseValue {
			payAddr = s.StateCfg.ActiveMiningAddrs[rand.Intn(
				len(
					s.StateCfg.
						ActiveMiningAddrs,
				),
			)]
		}
		// Create a new block template that has a coinbase which anyone can redeem.
		//
		// This is only acceptable because the returned block template doesn't include the coinbase, so the caller will
		// ultimately create their own coinbase which pays to the appropriate address(es).
		blkTemplate, e := generator.NewBlockTemplate(payAddr, state.Algo)
		if e != nil {
			return InternalRPCError(
				"(rpcserver.go) Failed to create new block "+
					"template: "+e.Error(), "",
			)
		}
		template = blkTemplate
		msgBlock = template.Block
		targetDifficulty = fmt.Sprintf(
			"%064x",
			bits.CompactToBig(msgBlock.Header.Bits),
		)
		// Get the minimum allowed timestamp for the block based on the median timestamp of the last several blocks per
		// the chain consensus rules.
		best := s.Cfg.Chain.BestSnapshot()
		minTimestamp := mining.MinimumMedianTime(best)
		// Update work state to ensure another block template isn't generated until needed.
		state.Template = template
		state.LastGenerated = time.Now()
		state.LastTxUpdate = lastTxUpdate
		state.prevHash = latestHash
		state.MinTimestamp = minTimestamp
		D.F(
			"generated block template (timestamp %v, target %s, merkle root %s)",
			msgBlock.Header.Timestamp,
			targetDifficulty,
			msgBlock.Header.MerkleRoot,
		)

		// Notify any clients that are long polling about the new template.
		state.NotifyLongPollers(latestHash, lastTxUpdate)
	} else {
		// At this point, there is a saved block template and another request for a template was made, but either the
		// available transactions haven't change or it hasn't been long enough to trigger a new block template to be
		// generated.
		//
		// So, update the existing block template.
		//
		// When the caller requires a full coinbase as opposed to only the pertinent details needed to create their own
		// coinbase, add a payment address to the output of the coinbase of the template if it doesn't already have one.
		//
		// Since this requires mining addresses to be specified via the config, an error is returned if none have been
		// specified.
		if !useCoinbaseValue && !template.ValidPayAddress {
			// Choose a payment address at random.
			payToAddr := s.StateCfg.ActiveMiningAddrs[rand.Intn(
				len(
					s.
						StateCfg.ActiveMiningAddrs,
				),
			)]
			// Update the block coinbase output of the template to pay to the randomly selected payment address.
			pkScript, e := txscript.PayToAddrScript(payToAddr)
			if e != nil {
				context := "Failed to create pay-to-addr script"
				return InternalRPCError(e.Error(), context)
			}
			template.Block.Transactions[0].TxOut[0].PkScript = pkScript
			template.ValidPayAddress = true
			// Update the merkle root.
			block := block.NewBlock(template.Block)
			merkles := blockchain.BuildMerkleTreeStore(block.Transactions(), false)
			template.Block.Header.MerkleRoot = *merkles.GetRoot()
		}
		// Set locals for convenience.
		msgBlock = template.Block
		targetDifficulty = fmt.Sprintf(
			"%064x",
			bits.CompactToBig(msgBlock.Header.Bits),
		)
		// Update the time of the block template to the current time while accounting for the median time of the past
		// several blocks per the chain consensus rules.
		e := generator.UpdateBlockTime(0, msgBlock)
		if e != nil {
			D.Ln(e)

		}
		msgBlock.Header.Nonce = 0
		D.F(
			"updated block template (timestamp %v, target %s)",
			msgBlock.Header.Timestamp,
			targetDifficulty,
		)

	}
	return nil
}

// NotifyNewTransactions notifies both websocket and getblocktemplate long poll clients of the passed transactions.
//
// This function should be called whenever new transactions are added to the mempool.
func (s *Server) NotifyNewTransactions(txns []*mempool.TxDesc) {
	for _, txD := range txns {
		// Notify websocket clients about mempool transactions.
		s.NtfnMgr.SendNotifyMempoolTx(txD.Tx, true)
		// Potentially notify any getblocktemplate long poll clients about stale block templates due to the new
		// transaction.
		s.GBTWorkState.NotifyMempoolTx(s.Cfg.TxMemPool.LastUpdated())
	}
}

// RequestedProcessShutdown returns a channel that is sent to when an authorized RPC client requests the process to
// shutdown. If the request can not be read immediately, it is dropped.
func (s *Server) RequestedProcessShutdown() qu.C {
	return s.RequestProcessShutdown
}

// Start is used by Server to start the rpc listener.
func (s *Server) Start() {
	if atomic.AddInt32(&s.Started, 1) != 1 {
		return
	}
	rpcServeMux := http.NewServeMux()
	httpServer := &http.Server{
		Handler: rpcServeMux,
		// Timeout connections which don't complete the initial handshake within the allowed timeframe.
		ReadTimeout: time.Second * RPCAuthTimeoutSeconds,
	}
	rpcServeMux.HandleFunc(
		"/", func(w http.ResponseWriter, r *http.Request) {
			I.Ln("handling new connection from", r.RemoteAddr)
			w.Header().Set("Connection", "close")
			w.Header().Set("Content-Type", "application/json")
			r.Close = true
			// Limit the number of connections to max allowed.
			if s.LimitConnections(w, r.RemoteAddr) {
				return
			}
			// Keep track of the number of connected clients.
			s.IncrementClients()
			defer s.DecrementClients()
			_, isAdmin, e := s.CheckAuth(r, true)
			if e != nil {
				JSONAuthFail(w)
				return
			}
			// Read and respond to the request.
			s.JSONRPCRead(w, r, isAdmin)
		},
	)
	// Websocket endpoint.
	rpcServeMux.HandleFunc(
		"/ws", func(w http.ResponseWriter, r *http.Request) {
			I.Ln("handling new connection from", r.RemoteAddr)
			authenticated, isAdmin, e := s.CheckAuth(r, false)
			if e != nil {
				JSONAuthFail(w)
				return
			}
			// Attempt to upgrade the connection to a websocket connection using the default size for read/write buffers.
			ws, e := websocket.Upgrade(w, r, nil, 0, 0)
			if e != nil {
				if _, ok := e.(websocket.HandshakeError); !ok {
					E.Ln("unexpected websocket error:", e)

				}
				http.Error(w, "400 Bad Request.", http.StatusBadRequest)
				return
			}
			s.WebsocketHandler(ws, r.RemoteAddr, authenticated, isAdmin)
		},
	)
	for _, listener := range s.Cfg.Listeners {
		s.WG.Add(1)
		go func(listener net.Listener) {
			I.Ln("chain RPC Server listening on ", listener.Addr())
			e := httpServer.Serve(listener)
			if e != nil {
				D.Ln(e)
			}
			D.Ln("chain RPC listener done for", listener.Addr())
			if e := listener.Close(); E.Chk(e) {
			}
			s.WG.Done()
		}(listener)
	}
	s.NtfnMgr.WG.Add(2)
	s.NtfnMgr.Start()
}

// Stop is used by Server.go_ to stop the rpc listener.
func (s *Server) Stop() (e error) {
	if atomic.AddInt32(&s.Shutdown, 1) != 1 {
		W.Ln("RPC Server is already in the process of shutting down")
		return nil
	}
	I.Ln("RPC Server shutting down")
	s.Quit.Q()
	for _, listener := range s.Cfg.Listeners {
		e := listener.Close()
		if e != nil {
			E.Ln("problem shutting down RPC:", e)
			return e
		}
	}
	s.NtfnMgr.Shutdown()
	s.NtfnMgr.WaitForShutdown()
	s.WG.Wait()
	D.Ln("RPC Server shutdown complete")
	return nil
}

// CheckAuth checks the HTTP Basic authentication supplied by a wallet or RPC client in the HTTP request r.
//
// If the supplied authentication does not match the username and password expected, a non-nil error is returned. This
// check is time-constant.
//
// The first bool return value signifies auth success ( true if successful) and the second bool return value specifies
// whether the user can change the state of the Server (true) or whether the user is limited (false).
//
// The second is always false if the first is.
func (s *Server) CheckAuth(r *http.Request, require bool) (bool, bool, error) {
	authhdr := r.Header["Authorization"]
	if len(authhdr) == 0 {
		if require {
			W.Ln("RPC authentication failure from", r.RemoteAddr)

			return false, false, errors.New("auth failure")
		}
		return false, false, nil
	}
	authsha := sha256.Sum256([]byte(authhdr[0]))
	// Chk for limited auth first as in environments with limited users, those are probably expected to have a higher
	// volume of calls
	limitcmp := subtle.ConstantTimeCompare(authsha[:], s.LimitAuthSHA[:])
	if limitcmp == 1 {
		return true, false, nil
	}
	// Chk for admin-level auth
	cmp := subtle.ConstantTimeCompare(authsha[:], s.AuthSHA[:])
	if cmp == 1 {
		return true, true, nil
	}
	// Request's auth doesn't match either user
	W.Ln("RPC authentication failure from", r.RemoteAddr)

	return false, false, errors.New("auth failure")
}

// DecrementClients subtracts one from the number of connected RPC clients. Note this only applies to standard clients.
//
// Websocket clients have their own limits and are tracked separately. This function is safe for concurrent access.
func (s *Server) DecrementClients() {
	atomic.AddInt32(&s.NumClients, -1)
}

// HandleBlockchainNotification handles callbacks for notifications from blockchain. It notifies clients that are long
// polling for changes or subscribed to websockets notifications.
func (s *Server) HandleBlockchainNotification(notification *blockchain.Notification) {
	if s.Cfg.Chain.IsCurrent() {
		switch notification.Type {
		case blockchain.NTBlockAccepted:
			block, ok := notification.Data.(*block.Block)
			if !ok {
				W.Ln("chain accepted notification is not a block")
				break
			}
			// Allow any clients performing long polling via the getblocktemplate RPC to be notified when the new block
			// causes their old block template to become stale.
			s.GBTWorkState.NotifyBlockConnected(block.Hash())
		case blockchain.NTBlockConnected:
			block, ok := notification.Data.(*block.Block)
			if !ok {
				W.Ln("chain connected notification is not a block")
				break
			}
			// Notify registered websocket clients of incoming block.
			s.NtfnMgr.SendNotifyBlockConnected(block)
		case blockchain.NTBlockDisconnected:
			block, ok := notification.Data.(*block.Block)
			if !ok {
				W.Ln("chain disconnected notification is not a block.")
				break
			}
			// Notify registered websocket clients.
			s.NtfnMgr.SendNotifyBlockDisconnected(block)
		}
	}
}

// HTTPStatusLine returns a response Status-Line (RFC 2616 Section 6.1) for the given request and response status code.
//
// This function was lifted and adapted from the standard library HTTP Server code since it's not exported.
func (s *Server) HTTPStatusLine(req *http.Request, code int) string {
	// Fast path:
	key := code
	proto11 := req.ProtoAtLeast(1, 1)
	if !proto11 {
		key = -key
	}
	s.StatusLock.RLock()
	line, ok := s.StatusLines[key]
	s.StatusLock.RUnlock()
	if ok {
		return line
	}
	// Slow path:
	proto := "HTTP/1.0"
	if proto11 {
		proto = "HTTP/1.1"
	}
	codeStr := strconv.Itoa(code)
	text := http.StatusText(code)
	if text != "" {
		line = proto + " " + codeStr + " " + text + "\r\n"
		s.StatusLock.Lock()
		s.StatusLines[key] = line
		s.StatusLock.Unlock()
	} else {
		text = "status code " + codeStr
		line = proto + " " + codeStr + " " + text + "\r\n"
	}
	return line
}

// IncrementClients adds one to the number of connected RPC clients. Note this only applies to standard clients.
//
// Websocket clients have their own limits and are tracked separately.
//
// This function is safe for concurrent access.
func (s *Server) IncrementClients() {
	atomic.AddInt32(&s.NumClients, 1)
}

// JSONRPCRead handles reading and responding to RPC messages.
func (s *Server) JSONRPCRead(
	w http.ResponseWriter, r *http.Request, isAdmin bool,
) {
	if atomic.LoadInt32(&s.Shutdown) != 0 {
		return
	}
	// Read and close the JSON-RPC request body from the caller.
	var e error
	var body []byte
	if body, e = ioutil.ReadAll(r.Body); E.Chk(e) {
	}
	I.Ln("read JSONRPC request", string(body))
	if e = r.Body.Close(); E.Chk(e) {
	}
	if e != nil {
		errCode := http.StatusBadRequest
		http.Error(
			w, fmt.Sprintf(
				"%d error reading JSON message: %v",
				errCode, e,
			), errCode,
		)
		return
	}
	// Unfortunately, the http Server doesn't provide the ability to change the read deadline for the new connection and
	// having one breaks long polling.
	//
	// However, not having a read deadline on the initial connection would mean clients can connect and idle forever.
	//
	// Thus, hijack the connection from the HTTP Server, clear the read deadline, and handle writing the response
	// manually.
	hj, ok := w.(http.Hijacker)
	if !ok {
		errMsg := "webserver doesn't support hijacking"
		W.F(errMsg)
		errCode := http.StatusInternalServerError
		http.Error(w, strconv.Itoa(errCode)+" "+errMsg, errCode)
		return
	}
	conn, buf, e := hj.Hijack()
	if e != nil {
		W.Ln("failed to hijack HTTP connection:", e)
		errCode := http.StatusInternalServerError
		http.Error(w, strconv.Itoa(errCode)+" "+e.Error(), errCode)
		return
	}
	defer func() {
		if e = buf.Flush(); E.Chk(e) {
		}
		if e = conn.Close(); E.Chk(e) {
		}
	}()
	e = conn.SetReadDeadline(TimeZeroVal)
	if e != nil {
		D.Ln(e)
	}
	// Attempt to parse the raw body into a JSON-RPC request.
	var responseID interface{}
	var jsonErr error
	var result interface{}
	var request btcjson.Request
	if e = js.Unmarshal(body, &request); E.Chk(e) {
		jsonErr = &btcjson.RPCError{
			Code:    btcjson.ErrRPCParse.Code,
			Message: "Failed to parse request: " + e.Error(),
		}
	}
	I.S(request)
	if jsonErr == nil {
		// The JSON-RPC 1.0 spec defines that notifications must have their "id"
		// set to null and states that notifications do not have a response.
		//
		// A JSON-RPC 2.0 notification is a request with "json-rpc":"2.0", and without an "id" member.
		//
		// The specification states that notifications must not be responded to. JSON-RPC 2.0 permits the null value as
		// a valid request id, therefore such requests are not notifications.
		//
		// Bitcoin Core serves requests with "id":null or even an absent "id", and responds to such requests with
		// "id":null in the response. Save does not respond to any request without and "id" or "id":null, regardless the
		// indicated JSON-RPC protocol version unless RPC quirks are enabled.
		//
		// With RPC quirks enabled, such requests will be responded to if the request does not indicate JSON-RPC
		// version. RPC quirks can be enabled by the user to avoid compatibility issues with software relying on Core's
		// behavior.
		if request.ID == nil && !(s.Config.RPCQuirks.True() && request.Jsonrpc == "") {
			return
		}
		// The parse was at least successful enough to have an ID so set it for the response.
		responseID = request.ID
		// Setup a close notifier. Since the connection is hijacked, the CloseNotifer on the ResponseWriter is not
		// available.
		closeChan := qu.Ts(1)
		go func() {
			_, e = conn.Read(make([]byte, 1))
			if e != nil {
				closeChan.Q()
			}
		}()
		// Chk if the user is limited and set error if method unauthorized
		if !isAdmin {
			if _, ok := RPCLimited[request.Method]; !ok {
				jsonErr = &btcjson.RPCError{
					Code:    btcjson.ErrRPCInvalidParams.Code,
					Message: "limited user not authorized for this method",
				}
			}
		}
		if jsonErr == nil {
			// Attempt to parse the JSON-RPC request into a known concrete command.
			parsedCmd := ParseCmd(&request)
			if parsedCmd.Err != nil {
				jsonErr = parsedCmd.Err
			} else {
				I.Ln("calling method:", parsedCmd.Method)
				result, jsonErr = s.StandardCmdResult(parsedCmd, closeChan)
			}
		}
	}
	// Marshal the response.
	var msg []byte
	msg, e = CreateMarshalledReply(responseID, result, jsonErr)
	if e != nil {
		E.Ln("failed to marshal reply:", e)
		return
	}
	I.Ln("\n" + string(msg))
	// Write the response.
	e = s.WriteHTTPResponseHeaders(r, w.Header(), http.StatusOK, buf)
	if e != nil {
		E.Ln(e)
		return
	}
	if _, e = buf.Write(msg); E.Chk(e) {
		E.Ln("failed to write marshalled reply:", e)
	}
	// Terminate with newline to maintain compatibility with Bitcoin Core.
	if e := buf.WriteByte('\n'); E.Chk(e) {
		E.Ln("failed to append terminating newline to reply:", e)

	}
}

// LimitConnections responds with a 503 service unavailable and returns true if adding another client would exceed the
// maximum allow RPC clients.
//
// This function is safe for concurrent access.
func (s *Server) LimitConnections(
	w http.ResponseWriter, remoteAddr string,
) bool {
	if int(atomic.LoadInt32(&s.NumClients)+1) > s.Config.RPCMaxClients.V() {
		I.F(
			"max RPC clients exceeded [%d] - disconnecting client %s",
			s.Config.RPCMaxClients, remoteAddr,
		)

		http.Error(
			w, "503 Too busy.  Try again later.",
			http.StatusServiceUnavailable,
		)
		return true
	}
	return false
}

// StandardCmdResult checks that a parsed command is a standard Bitcoin JSON-RPC command and runs the appropriate
// handler to reply to the command.
//
// Any commands which are not recognized or not implemented will return an error suitable for use in replies.
func (s *Server) StandardCmdResult(
	cmd *ParsedRPCCmd,
	closeChan qu.C,
) (interface{}, error) {
	handler, ok := RPCHandlers[cmd.Method]
	if ok {
		goto handled
	}
	_, ok = RPCAskWallet[cmd.Method]
	if ok {
		handler.Fn = HandleAskWallet
		goto handled
	}
	_, ok = RPCUnimplemented[cmd.Method]
	if ok {
		handler.Fn = HandleUnimplemented
		goto handled
	}
	return nil, btcjson.ErrRPCMethodNotFound
handled:
	return handler.Fn(s, cmd.Cmd, closeChan)
}

// WriteHTTPResponseHeaders writes the necessary response headers prior to writing an HTTP body given a request to use
// for protocol negotiation, headers to write, a status code, and a writer.
func (s *Server) WriteHTTPResponseHeaders(
	req *http.Request,
	headers http.Header, code int, w io.Writer,
) (e error) {
	_, e = io.WriteString(w, s.HTTPStatusLine(req, code))
	if e != nil {
		return e
	}
	e = headers.Write(w)
	if e != nil {
		return e
	}
	_, e = io.WriteString(w, "\r\n")
	return e
}

// BuilderScript is a convenience function which is used for hard-coded scripts built with the script builder. Any
// errors are converted to a panic since it is only, and must only, be used with hard-coded, and therefore, known good,
// scripts.
func BuilderScript(builder *txscript.ScriptBuilder) []byte {
	script, e := builder.Script()
	if e != nil {
		panic(e)
	}
	return script
}

// ChainErrToGBTErrString converts an error returned from btcchain to a string which matches the reasons and format
// described in BIP0022 for rejection reasons.
func ChainErrToGBTErrString(e error) string {
	// When the passed error is not a RuleError, just return a generic rejected string with the error text.
	ruleErr, ok := e.(blockchain.RuleError)
	if !ok {
		return "rejected: " + e.Error()
	}
	switch ruleErr.ErrorCode {
	case blockchain.ErrDuplicateBlock:
		return "duplicate"
	case blockchain.ErrBlockTooBig:
		return "bad-blk-length"
	case blockchain.ErrBlockWeightTooHigh:
		return "bad-blk-weight"
	case blockchain.ErrBlockVersionTooOld:
		return "bad-version"
	case blockchain.ErrInvalidTime:
		return "bad-time"
	case blockchain.ErrTimeTooOld:
		return "time-too-old"
	case blockchain.ErrTimeTooNew:
		return "time-too-new"
	case blockchain.ErrDifficultyTooLow:
		return "bad-diffbits"
	case blockchain.ErrUnexpectedDifficulty:
		return "bad-diffbits"
	case blockchain.ErrHighHash:
		return "high-hash"
	case blockchain.ErrBadMerkleRoot:
		return "bad-txnmrklroot"
	case blockchain.ErrBadCheckpoint:
		return "bad-checkpoint"
	case blockchain.ErrForkTooOld:
		return "fork-too-old"
	case blockchain.ErrCheckpointTimeTooOld:
		return "checkpoint-time-too-old"
	case blockchain.ErrNoTransactions:
		return "bad-txns-none"
	case blockchain.ErrNoTxInputs:
		return "bad-txns-noinputs"
	case blockchain.ErrNoTxOutputs:
		return "bad-txns-nooutputs"
	case blockchain.ErrTxTooBig:
		return "bad-txns-size"
	case blockchain.ErrBadTxOutValue:
		return "bad-txns-outputvalue"
	case blockchain.ErrDuplicateTxInputs:
		return "bad-txns-dupinputs"
	case blockchain.ErrBadTxInput:
		return "bad-txns-badinput"
	case blockchain.ErrMissingTxOut:
		return "bad-txns-missinginput"
	case blockchain.ErrUnfinalizedTx:
		return "bad-txns-unfinalizedtx"
	case blockchain.ErrDuplicateTx:
		return "bad-txns-duplicate"
	case blockchain.ErrOverwriteTx:
		return "bad-txns-overwrite"
	case blockchain.ErrImmatureSpend:
		return "bad-txns-maturity"
	case blockchain.ErrSpendTooHigh:
		return "bad-txns-highspend"
	case blockchain.ErrBadFees:
		return "bad-txns-fees"
	case blockchain.ErrTooManySigOps:
		return "high-sigops"
	case blockchain.ErrFirstTxNotCoinbase:
		return "bad-txns-nocoinbase"
	case blockchain.ErrMultipleCoinbases:
		return "bad-txns-multicoinbase"
	case blockchain.ErrBadCoinbaseScriptLen:
		return "bad-cb-length"
	case blockchain.ErrBadCoinbaseValue:
		return "bad-cb-value"
	case blockchain.ErrMissingCoinbaseHeight:
		return "bad-cb-height"
	case blockchain.ErrBadCoinbaseHeight:
		return "bad-cb-height"
	case blockchain.ErrScriptMalformed:
		return "bad-script-malformed"
	case blockchain.ErrScriptValidation:
		return "bad-script-validate"
	// case blockchain.ErrUnexpectedWitness:
	// 	return "unexpected-witness"
	// case blockchain.ErrInvalidWitnessCommitment:
	// 	return "bad-witness-nonce-size"
	// case blockchain.ErrWitnessCommitmentMismatch:
	// 	return "bad-witness-merkle-match"
	case blockchain.ErrPreviousBlockUnknown:
		return "prev-blk-not-found"
	case blockchain.ErrInvalidAncestorBlock:
		return "bad-prevblk"
	case blockchain.ErrPrevBlockNotBest:
		return "inconclusive-not-best-prvblk"
	}
	return "rejected: " + e.Error()
}

// CreateMarshalledReply returns a new marshalled JSON-RPC response given the passed parameters. It will automatically
// convert errors that are not of the type *json.RPCError to the appropriate type as needed.
func CreateMarshalledReply(id, result interface{}, replyErr error) (
	[]byte, error,
) {
	var jsonErr *btcjson.RPCError
	if replyErr != nil {
		if jErr, ok := replyErr.(*btcjson.RPCError); ok {
			jsonErr = jErr
		} else {
			jsonErr = InternalRPCError(replyErr.Error(), "")
		}
	}
	return btcjson.MarshalResponse(id, result, jsonErr)
}

// CreateTxRawResult converts the passed transaction and associated parameters to a raw transaction JSON object.
func CreateTxRawResult(
	chainParams *chaincfg.Params, mtx *wire.MsgTx,
	txHash string, blkHeader *wire.BlockHeader, blkHash string, blkHeight int32,
	chainHeight int32,
) (*btcjson.TxRawResult, error) {
	mtxHex, e := MessageToHex(mtx)
	if e != nil {
		return nil, e
	}
	txReply := &btcjson.TxRawResult{
		Hex:      mtxHex,
		Txid:     txHash,
		Hash:     mtx.TxHash().String(), // WitnessHash().String(),
		Size:     int32(mtx.SerializeSize()),
		Vsize:    int32(mempool.GetTxVirtualSize(util.NewTx(mtx))),
		Vin:      CreateVinList(mtx),
		Vout:     CreateVoutList(mtx, chainParams, nil),
		Version:  mtx.Version,
		LockTime: mtx.LockTime,
	}
	if blkHeader != nil {
		// This is not a typo, they are identical in bitcoind as well.
		txReply.Time = blkHeader.Timestamp.Unix()
		txReply.Blocktime = blkHeader.Timestamp.Unix()
		txReply.BlockHash = blkHash
		txReply.Confirmations = uint64(1 + chainHeight - blkHeight)
	}
	return txReply, nil
}

// CreateVinList returns a slice of JSON objects for the inputs of the passed transaction.
func CreateVinList(mtx *wire.MsgTx) []btcjson.Vin {
	// Coinbase transactions only have a single txin by definition.
	vinList := make([]btcjson.Vin, len(mtx.TxIn))
	if blockchain.IsCoinBaseTx(mtx) {
		txIn := mtx.TxIn[0]
		vinList[0].Coinbase = hex.EncodeToString(txIn.SignatureScript)
		vinList[0].Sequence = txIn.Sequence
		// vinList[0].Witness = WitnessToHex(txIn.Witness)
		return vinList
	}
	for i, txIn := range mtx.TxIn {
		// The disassembled string will contain [error] inline if the script doesn't
		// fully parse, so ignore the error here.
		disbuf, _ := txscript.DisasmString(txIn.SignatureScript)
		vinEntry := &vinList[i]
		vinEntry.Txid = txIn.PreviousOutPoint.Hash.String()
		vinEntry.Vout = txIn.PreviousOutPoint.Index
		vinEntry.Sequence = txIn.Sequence
		vinEntry.ScriptSig = &btcjson.ScriptSig{
			Asm: disbuf,
			Hex: hex.EncodeToString(txIn.SignatureScript),
		}
		// if mtx.HasWitness() {
		// 	vinEntry.Witness = WitnessToHex(txIn.Witness)
		// }
	}
	return vinList
}

// CreateVinListPrevOut returns a slice of JSON objects for the inputs of the
// passed transaction.
func CreateVinListPrevOut(
	s *Server, mtx *wire.MsgTx,
	chainParams *chaincfg.Params, vinExtra bool,
	filterAddrMap map[string]struct{},
) ([]btcjson.VinPrevOut, error) {
	// Coinbase transactions only have a single txin by definition.
	if blockchain.IsCoinBaseTx(mtx) {
		// Only include the transaction if the filter map is empty because a coinbase
		// input has no addresses and so would never match a non-empty filter.
		if len(filterAddrMap) != 0 {
			return nil, nil
		}
		txIn := mtx.TxIn[0]
		vinList := make([]btcjson.VinPrevOut, 1)
		vinList[0].Coinbase = hex.EncodeToString(txIn.SignatureScript)
		vinList[0].Sequence = txIn.Sequence
		return vinList, nil
	}
	// Use a dynamically sized list to accommodate the address filter.
	vinList := make([]btcjson.VinPrevOut, 0, len(mtx.TxIn))
	// Lookup all of the referenced transaction outputs needed to populate the previous output information if requested.
	var originOutputs map[wire.OutPoint]wire.TxOut
	if vinExtra || len(filterAddrMap) > 0 {
		var e error
		originOutputs, e = FetchInputTxos(s, mtx)
		if e != nil {
			return nil, e
		}
	}
	for _, txIn := range mtx.TxIn {
		// The disassembled string will contain [error] inline if the script doesn't
		// fully parse, so ignore the error here.
		disbuf, _ := txscript.DisasmString(txIn.SignatureScript)
		// Create the basic input entry without the additional optional previous output
		// details which will be added later if requested and available.
		prevOut := &txIn.PreviousOutPoint
		vinEntry := btcjson.VinPrevOut{
			Txid:     prevOut.Hash.String(),
			Vout:     prevOut.Index,
			Sequence: txIn.Sequence,
			ScriptSig: &btcjson.ScriptSig{
				Asm: disbuf,
				Hex: hex.EncodeToString(txIn.SignatureScript),
			},
		}
		// if len(txIn.Witness) != 0 {
		// 	vinEntry.Witness = WitnessToHex(txIn.Witness)
		// }
		// Add the entry to the list now if it already passed the filter since the
		// previous output might not be available.
		passesFilter := len(filterAddrMap) == 0
		if passesFilter {
			vinList = append(vinList, vinEntry)
		}
		// Only populate previous output information if requested and available.
		if len(originOutputs) == 0 {
			continue
		}
		originTxOut, ok := originOutputs[*prevOut]
		if !ok {
			continue
		}
		// Ignore the error here since an error means the script couldn't parse and there is no additional information
		// about it anyways.
		_, addrs, _, _ := txscript.ExtractPkScriptAddrs(originTxOut.PkScript, chainParams)
		// Encode the addresses while checking if the address passes the filter when needed.
		encodedAddrs := make([]string, len(addrs))
		for j, addr := range addrs {
			encodedAddr := addr.EncodeAddress()
			encodedAddrs[j] = encodedAddr
			// No need to check the map again if the filter already passes.
			if passesFilter {
				continue
			}
			if _, exists := filterAddrMap[encodedAddr]; exists {
				passesFilter = true
			}
		}
		// Ignore the entry if it doesn't pass the filter.
		if !passesFilter {
			continue
		}
		// Add entry to the list if it wasn't already done above.
		if len(filterAddrMap) != 0 {
			vinList = append(vinList, vinEntry)
		}
		// Update the entry with previous output information if requested.
		if vinExtra {
			vinListEntry := &vinList[len(vinList)-1]
			vinListEntry.PrevOut = &btcjson.PrevOut{
				Addresses: encodedAddrs,
				Value:     amt.Amount(originTxOut.Value).ToDUO(),
			}
		}
	}
	return vinList, nil
}

// CreateVoutList returns a slice of JSON objects for the outputs of the passed transaction.
func CreateVoutList(
	mtx *wire.MsgTx, chainParams *chaincfg.Params,
	filterAddrMap map[string]struct{},
) []btcjson.Vout {
	voutList := make([]btcjson.Vout, 0, len(mtx.TxOut))
	for i, v := range mtx.TxOut {
		// The disassembled string will contain [error] inline if the script doesn't fully parse, so ignore the error
		// here.
		disbuf, _ := txscript.DisasmString(v.PkScript)
		// Ignore the error here since an error means the script couldn't parse and there is no additional information
		// about it anyways.
		scriptClass, addrs, reqSigs, _ := txscript.ExtractPkScriptAddrs(v.PkScript, chainParams)
		// Encode the addresses while checking if the address passes the filter when needed.
		passesFilter := len(filterAddrMap) == 0
		encodedAddrs := make([]string, len(addrs))
		for j, addr := range addrs {
			encodedAddr := addr.EncodeAddress()
			encodedAddrs[j] = encodedAddr
			// No need to check the map again if the filter already passes.
			if passesFilter {
				continue
			}
			if _, exists := filterAddrMap[encodedAddr]; exists {
				passesFilter = true
			}
		}
		if !passesFilter {
			continue
		}
		var vout btcjson.Vout
		vout.N = uint32(i)
		vout.Value = amt.Amount(v.Value).ToDUO()
		vout.ScriptPubKey.Addresses = encodedAddrs
		vout.ScriptPubKey.Asm = disbuf
		vout.ScriptPubKey.Hex = hex.EncodeToString(v.PkScript)
		vout.ScriptPubKey.Type = scriptClass.String()
		vout.ScriptPubKey.ReqSigs = int32(reqSigs)
		voutList = append(voutList, vout)
	}
	return voutList
}

// DecodeTemplateID decodes an ID that is used to uniquely identify a block template. This is mainly used as a mechanism
// to track when to update clients that are using long polling for block templates. The ID consists of the previous
// block hash for the associated template and the time the associated template was generated.
func DecodeTemplateID(templateID string) (*chainhash.Hash, int64, error) {
	fields := strings.Split(templateID, "-")
	if len(fields) != 2 {
		return nil, 0, errors.New("invalid longpollid format")
	}
	prevHash, e := chainhash.NewHashFromStr(fields[0])
	if e != nil {
		return nil, 0, errors.New("invalid longpollid format")
	}
	lastGenerated, e := strconv.ParseInt(fields[1], 10, 64)
	if e != nil {
		return nil, 0, errors.New("invalid longpollid format")
	}
	return prevHash, lastGenerated, nil
}

// EncodeTemplateID encodes the passed details into an ID that can be used to uniquely identify a block template.
func EncodeTemplateID(
	prevHash *chainhash.Hash, lastGenerated time.Time,
) string {
	return fmt.Sprintf("%s-%d", prevHash.String(), lastGenerated.Unix())
}

// FetchInputTxos fetches the outpoints from all transactions referenced by the inputs to the passed transaction by
// checking the transaction mempool first then the transaction index for those already mined into blocks.
func FetchInputTxos(s *Server, tx *wire.MsgTx) (
	map[wire.OutPoint]wire.TxOut, error,
) {
	mp := s.Cfg.TxMemPool
	originOutputs := make(map[wire.OutPoint]wire.TxOut)
	for txInIndex, txIn := range tx.TxIn {
		// Attempt to fetch and use the referenced transaction from the memory pool.
		origin := &txIn.PreviousOutPoint
		originTx, e := mp.FetchTransaction(&origin.Hash)
		if e == nil {
			txOuts := originTx.MsgTx().TxOut
			if origin.Index >= uint32(len(txOuts)) {
				errStr := fmt.Sprintf(
					"unable to find output %v referenced from transaction %s:%d", origin, tx.TxHash(), txInIndex,
				)
				return nil, InternalRPCError(errStr, "")
			}
			originOutputs[*origin] = *txOuts[origin.Index]
			continue
		}
		// Look up the location of the transaction.
		blockRegion, e := s.Cfg.TxIndex.TxBlockRegion(&origin.Hash)
		if e != nil {
			context := "Failed to retrieve transaction location"
			return nil, InternalRPCError(e.Error(), context)
		}
		if blockRegion == nil {
			return nil, NoTxInfoError(&origin.Hash)
		}
		// Load the raw transaction bytes from the database.
		var txBytes []byte
		e = s.Cfg.DB.View(
			func(dbTx database.Tx) (e error) {
				txBytes, e = dbTx.FetchBlockRegion(blockRegion)
				return e
			},
		)
		if e != nil {
			return nil, NoTxInfoError(&origin.Hash)
		}
		// Deserialize the transaction
		var msgTx wire.MsgTx
		e = msgTx.Deserialize(bytes.NewReader(txBytes))
		if e != nil {
			context := deserialfail
			return nil, InternalRPCError(e.Error(), context)
		}
		// Add the referenced output to the map.
		if origin.Index >= uint32(len(msgTx.TxOut)) {
			errStr := fmt.Sprintf(
				"unable to find output %v "+
					"referenced from transaction %s:%d", origin,
				tx.TxHash(), txInIndex,
			)
			return nil, InternalRPCError(errStr, "")
		}
		originOutputs[*origin] = *msgTx.TxOut[origin.Index]
	}
	return originOutputs, nil
}

// FetchMempoolTxnsForAddress queries the address index for all unconfirmed transactions that involve the provided
// address. The results will be limited by the number to skip and the number requested.
func FetchMempoolTxnsForAddress(
	s *Server, addr btcaddr.Address, numToSkip,
	numRequested uint32,
) ([]*util.Tx, uint32) {
	// There are no entries to return when there are less available than the
	// number being skipped.
	mpTxns := s.Cfg.AddrIndex.UnconfirmedTxnsForAddress(addr)
	numAvailable := uint32(len(mpTxns))
	if numToSkip > numAvailable {
		return nil, numAvailable
	}
	// Filter the available entries based on the number to skip and number requested.
	rangeEnd := numToSkip + numRequested
	if rangeEnd > numAvailable {
		rangeEnd = numAvailable
	}
	return mpTxns[numToSkip:rangeEnd], numToSkip
}

// GenCertPair generates a key/cert pair to the paths provided.
func GenCertPair(certFile, keyFile string) (e error) {
	I.Ln("generating TLS certificates...")
	org := "pod autogenerated cert"
	validUntil := time.Now().Add(10 * 365 * 24 * time.Hour)
	cert, key, e := util.NewTLSCertPair(org, validUntil, nil)
	if e != nil {
		return e
	}
	// Write cert and key files.
	if e = ioutil.WriteFile(certFile, cert, 0666); E.Chk(e) {
		return e
	}
	if e = ioutil.WriteFile(keyFile, key, 0600); E.Chk(e) {
		if e = os.Remove(certFile); E.Chk(e) {
		}
		return e
	}
	I.Ln("Done generating TLS certificates")
	return nil
}

// GetDifficultyRatio returns the proof-of-work difficulty as a multiple of the
// minimum difficulty using the passed bits field from the header of a block.
func GetDifficultyRatio(
	b uint32, params *chaincfg.Params,
	algo int32,
) float64 {
	// The minimum difficulty is the max possible proof-of-work limit bits converted
	// back to a number. Note this is not the same as the proof of work limit
	// directly because the block difficulty is encoded in a block with the compact
	// form which loses precision.
	max := bits.CompactToBig(0x1d00ffff)
	target := bits.CompactToBig(b)
	difficulty := new(big.Rat).SetFrac(max, target)
	outString := difficulty.FloatString(8)
	diff, e := strconv.ParseFloat(outString, 64)
	if e != nil {
		E.Ln("cannot get difficulty:", e)

		return 0
	}
	return diff
}

// NormalizeAddress returns addr with the passed default port appended if there
// is not already a port specified.
func NormalizeAddress(addr, defaultPort string) string {
	var e error
	_, _, e = net.SplitHostPort(addr)
	if e != nil {
		return net.JoinHostPort(addr, defaultPort)
	}
	return addr
}

func init() {

	RPCHandlers = RPCHandlersBeforeInit
	rand.Seed(time.Now().UnixNano())
}

// InternalRPCError is a convenience function to convert an internal error to an
// RPC error with the appropriate code set. It also logs the error to the RPC
// Server subsystem since internal errors really should not occur.
//
// The context parameter is only used in the l mayo be empty if it's not needed.
func InternalRPCError(errStr, context string) *btcjson.RPCError {
	logStr := errStr
	if context != "" {
		logStr = context + ": " + errStr
	}
	E.Ln(logStr)

	return btcjson.NewRPCError(btcjson.ErrRPCInternal.Code, errStr)
}

// JSONAuthFail sends a message back to the client if the http auth is rejected.
func JSONAuthFail(w http.ResponseWriter) {
	w.Header().Add("WWW-Authenticate", `Basic realm="pod RPC"`)
	http.Error(w, "401 Unauthorized.", http.StatusUnauthorized)
}

// MessageToHex serializes a message to the wire protocol encoding using the
// latest protocol version and returns a hex-encoded string of the result.
func MessageToHex(msg wire.Message) (string, error) {
	var buf bytes.Buffer
	if e := msg.BtcEncode(
		&buf, MaxProtocolVersion,
		wire.BaseEncoding,
	); E.Chk(e) {
		context := fmt.Sprintf("Failed to encode msg of type %Ter", msg)
		return "", InternalRPCError(e.Error(), context)
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

// NewGbtWorkState returns a new instance of a GBTWorkState with all internal
// fields initialized and ready to use.
func NewGbtWorkState(
	timeSource blockchain.MedianTimeSource,
	algoName string,
) *GBTWorkState {
	return &GBTWorkState{
		NotifyMap:  make(map[chainhash.Hash]map[int64]qu.C),
		TimeSource: timeSource,
		Algo:       algoName,
	}
}

// NewRPCServer returns a new instance of the RPCServer struct.
func NewRPCServer(
	config *ServerConfig, statecfg *active.Config,
	podcfg *config.Config,
) (*Server, error) {
	rpc := Server{
		Cfg:                    *config,
		Config:                 podcfg,
		StateCfg:               statecfg,
		StatusLines:            make(map[int]string),
		GBTWorkState:           NewGbtWorkState(config.TimeSource, config.Algo),
		HelpCacher:             NewHelpCacher(),
		RequestProcessShutdown: qu.T(),
		Quit:                   qu.T(),
	}
	if podcfg.Username.V() != "" && podcfg.Password.V() != "" {
		login := podcfg.Username.V() + ":" + podcfg.Password.V()
		// I.Ln("login:", login)
		auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(login))
		rpc.AuthSHA = sha256.Sum256([]byte(auth))
	}
	if podcfg.LimitUser.V() != "" && podcfg.LimitPass.V() != "" {
		login := podcfg.LimitUser.V() + ":" + podcfg.LimitPass.V()
		// I.Ln("login:", login)
		auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(login))
		rpc.LimitAuthSHA = sha256.Sum256([]byte(auth))
	}
	rpc.NtfnMgr = NewWSNotificationManager(&rpc)
	rpc.Cfg.Chain.Subscribe(rpc.HandleBlockchainNotification)
	return &rpc, nil
}

// ParseCmd parses a JSON-RPC request object into known concrete command. The
// err field of the returned ParsedRPCCmd struct will contain an RPC error that
// is suitable for use in replies if the command is invalid in some way such as
// an unregistered command or invalid parameters.
func ParseCmd(request *btcjson.Request) *ParsedRPCCmd {
	var parsedCmd ParsedRPCCmd
	parsedCmd.ID = request.ID
	parsedCmd.Method = request.Method
	cmd, e := btcjson.UnmarshalCmd(request)
	if e != nil {
		// When the error is because the method is not registered, produce a method not found RPC error.
		if jerr, ok := e.(btcjson.GeneralError); ok &&
			jerr.ErrorCode == btcjson.ErrUnregisteredMethod {
			parsedCmd.Err = btcjson.ErrRPCMethodNotFound
			return &parsedCmd
		}
		// Otherwise, some type of invalid parameters is the cause, so produce the equivalent RPC error.
		parsedCmd.Err = btcjson.NewRPCError(
			btcjson.ErrRPCInvalidParams.Code, e.Error(),
		)
		return &parsedCmd
	}
	parsedCmd.Cmd = cmd
	return &parsedCmd
}

// PeerExists determines if a certain peer is currently connected given
// information about all currently connected peers.
//
// Peer existence is determined using either a target address or node id.
func PeerExists(connMgr ServerConnManager, addr string, nodeID int32) bool {
	for _, pp := range connMgr.ConnectedPeers() {
		if pp.ToPeer().ID() == nodeID || pp.ToPeer().Addr() == addr {
			return true
		}
	}
	return false
}

// DecodeHexError is a convenience function for returning a nicely formatted RPC
// error which indicates the provided hex string failed to decode.
func DecodeHexError(gotHex string) *btcjson.RPCError {
	return btcjson.NewRPCError(
		btcjson.ErrRPCDecodeHexString,
		fmt.Sprintf(
			"Argument must be hexadecimal string (not %q)",
			gotHex,
		),
	)
}

// NoTxInfoError is a convenience function for returning a nicely formatted RPC
// error which indicates there is no information available for the provided
// transaction hash.
func NoTxInfoError(txHash *chainhash.Hash) *btcjson.RPCError {
	return btcjson.NewRPCError(
		btcjson.ErrRPCNoTxInfo,
		fmt.Sprintf(
			"No information available about transaction %v",
			txHash,
		),
	)
}

// SoftForkStatus converts a ThresholdState state into a human readable string
// corresponding to the particular state.
func SoftForkStatus(state blockchain.ThresholdState) (string, error) {
	switch state {
	case blockchain.ThresholdDefined:
		return "defined", nil
	case blockchain.ThresholdStarted:
		return "started", nil
	case blockchain.ThresholdLockedIn:
		return "lockedin", nil
	case blockchain.ThresholdActive:
		return "active", nil
	case blockchain.ThresholdFailed:
		return "failed", nil
	default:
		return "", fmt.Errorf("unknown deployment state: %v", state)
	}
}

// VerifyChain does?
func VerifyChain(s *Server, level, depth int32) (e error) {
	best := s.Cfg.Chain.BestSnapshot()
	finishHeight := best.Height - depth
	if finishHeight < 0 {
		finishHeight = 0
	}
	I.F(
		"verifying chain for %d blocks at level %d",
		best.Height-finishHeight,
		level,
	)

	var blk *block.Block
	for height := best.Height; height > finishHeight; height-- {
		// Level 0 just looks up the block.
		blk, e = s.Cfg.Chain.BlockByHeight(height)
		if e != nil {
			E.F("verify is unable to fetch block at height %d: %v", height, e)
			return e
		}
		powLimit := fork.GetMinDiff(
			fork.GetAlgoName(
				blk.WireBlock().Header.
					Version, height,
			), height,
		)
		var pb *block.Block
		if pb, e = s.Cfg.Chain.BlockByHash(&blk.WireBlock().Header.PrevBlock); E.Chk(e) {
			return
		}
		// Level 1 does basic chain sanity checks.
		if level > 0 {
			e = blockchain.CheckBlockSanity(
				blk,
				powLimit,
				s.Cfg.TimeSource,
				true,
				blk.Height(),
				pb.WireBlock().Header.Timestamp,
			)
			if e != nil {
				E.F(
					"verify is unable to validate block at hash %v height %d: %v %s",
					blk.Hash(), height, e,
				)
				return e
			}
		}
	}
	I.Ln("chain verify completed successfully")
	return nil
}

/*
// handleDebugLevel handles debuglevel commands.
func handleDebugLevel(	s *RPCServer, cmd interface{}, closeChan <-qu.C) (interface{}, error) {
	c := cmd.(*json.DebugLevelCmd)
	// Special show command to list supported subsystems.
	if c.LevelSpec == "show" {
		return fmt.Sprintf("Supported subsystems %v",
			supportedSubsystems()), nil
	}
	e := parseAndSetDebugLevels(c.LevelSpec)
	if e != nil  {
		return nil, &json.RPCError{
			Code:    json.ErrRPCInvalidParams.Code,
			Message: err.ScriptError(),
		}
	}
	return "Done.", nil
}
*/

// // WitnessToHex formats the passed witness stack as a slice of hex-encoded
// // strings to be used in a JSON response.
// func WitnessToHex(witness wire.TxWitness) []string {
// 	// Ensure nil is returned when there are no entries versus an empty slice so it
// 	// can properly be omitted as necessary.
// 	if len(witness) == 0 {
// 		return nil
// 	}
// 	result := make([]string, 0, len(witness))
// 	for _, wit := range witness {
// 		result = append(result, hex.EncodeToString(wit))
// 	}
// 	return result
// }
