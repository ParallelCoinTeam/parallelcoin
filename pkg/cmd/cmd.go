package cmd

import (
	core2 "github.com/p9c/pod/__OLDgui/____BEZI/test/pkg/duos/core"
	mod2 "github.com/p9c/pod/__OLDgui/____BEZI/test/pkg/duos/mod"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

type rcvar struct {
	d *core2.DuOS `json:"duos"`
}

type CmdInterface interface {
	GetBalance() float64
	GetTransactions() mod2.DuOStransactions
	GetTransactionsExcertps() mod2.DuOStransactionsExcerpts
	GetAddressBook() mod2.DuOSaddressBook
	DuoSend() string
	CreateNewAddress() string
	SaveAddressLabel()
}

type NodeInterface interface {
	GetNetworkLastBlock() int32
	GetBlocks()
	GetBlockExcerpt(height int) mod2.DuOSblock
	GetBlocksExcerpts(startBlock, blockHeight int) mod2.DuOSblocks
	Addnode(a *btcjson.AddNodeCmd)
	Createrawtransaction(a *btcjson.CreateRawTransactionCmd)
	Decoderawtransaction(a *btcjson.DecodeRawTransactionCmd)
	Decodescript(a *btcjson.DecodeScriptCmd)
	Estimatefee(a *btcjson.EstimateFeeCmd)
	Generate(a *btcjson.GenerateCmd)
	Getaddednodeinfo(a *btcjson.GetAddedNodeInfoCmd)
	Getbestblock() int64
	Getbestblockhash() string
	Getblock(a *btcjson.GetBlockCmd)
	GetBlockChainInfo()
	GetBlockCount() int64
	GetBlockHash(blockHeight int) string
	GetBlock(hash string) btcjson.GetBlockVerboseResult
	Getblockheader(a *btcjson.GetBlockHeaderCmd)
	GetConnectionCount() int32
	GetDifficulty() float64
	Gethashespersec()
	Getheaders(a *btcjson.GetHeadersCmd)
	Getinfo()
	Getmempoolinfo()
	Getmininginfo()
	Getnettotals()
	Getnetworkhashps(a *btcjson.GetNetworkHashPSCmd)
	GetPeerInfo() []*btcjson.GetPeerInfoResult
	Stop()
	Uptime()
	Validateaddress(a *btcjson.ValidateAddressCmd)
	Verifychain(a *btcjson.VerifyChainCmd)
	Verifymessage(a *btcjson.VerifyMessageCmd)
	GetWalletVersion() map[string]btcjson.VersionResult
}
