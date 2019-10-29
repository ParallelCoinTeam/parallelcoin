package core

import (
	"fmt"
	"github.com/p9c/pod/pkg/log"

	"github.com/p9c/pod/cmd/gui/mod"
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/pkg/chain/fork"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	database "github.com/p9c/pod/pkg/db"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/util"
)

func (d *DuOS) GetNetworkLastBlock() int32 {
	for _, g := range d.Cx.RPCServer.Cfg.ConnMgr.ConnectedPeers() {
		l := g.ToPeer().StatsSnapshot().LastBlock
		if l > d.Data.Status.NetworkLastBlock {
			d.Data.Status.NetworkLastBlock = l
		}
	}
	return d.Data.Status.NetworkLastBlock
}

// func (n *DuOSnode) GetBlocks() {
//	blks := []mod.Block{}
//	getBlockChain, err := rpc.HandleGetBlockChainInfo(n.rpc, nil, nil)
//	if err !=
//	}
//
//	n.Blocks = blks
// }

// func (n *DuOSnode) GetBlocks(per, page int) {
//	blks := []btcjson.GetBlockVerboseResult{}
//	getBlockChain, err := rpc.HandleGetBlockChainInfo(n.rpc, nil, nil)
//	if err != nil {
//		log.ERROR(err)
//		alert.Alert.Time = time.Now()
//		alert.Alert.Alert = err.Error()
//		alert.Alert.AlertType = "error"
//	}
//	blockChain := getBlockChain.(*btcjson.GetBlockChainInfoResult)
//	blockCount := int(blockChain.Blocks)
//	startBlock := blockCount - per*page
//	minusBlockStart := int(startBlock + per)
//	for ibh := minusBlockStart; ibh >= startBlock; {
//		block := btcjson.GetBlockVerboseResult{}
//		hcmd := btcjson.GetBlockHashCmd{
//			Index: int64(ibh),
//		}
//		hash, err := rpc.HandleGetBlockHash(n.rpc, &hcmd, nil)
//		if err != nil {
//		log.ERROR(err)
//			alert.Alert.Time = time.Now()
//			alert.Alert.Alert = err.Error()
//			alert.Alert.AlertType = "error"
//		}
//		if hash != nil {
//			verbose, verbosetx := true, true
//			bcmd := btcjson.GetBlockCmd{
//				Hash:      hash.(string),
//				Verbose:   &verbose,
//				VerboseTx: &verbosetx,
//			}
//			bl, err := rpc.HandleGetBlock(n.rpc, &bcmd, nil)
//			if err != nil {
//		log.ERROR(err)
//				alert.Alert.Time = time.Now()
//				alert.Alert.Alert = err.Error()
//				alert.Alert.AlertType = "error"
//			}
//			block = bl.(btcjson.GetBlockVerboseResult)
//			blks = append(blks, block)
//			ibh--
//		}
//	}
//	n.Blocks.Blocks = blks
//	n.Blocks.CurrentPage = page
//	n.Blocks.PageCount = blockCount / per
//
// }

func (d DuOS) GetBlockExcerpt(height int) (b mod.DuOSblock) {
	b = *new(mod.DuOSblock)
	hashHeight, err := d.Cx.RPCServer.Cfg.Chain.BlockHashByHeight(int32(height))
	if err != nil {
		log.ERROR(err)
}
	// Load the raw block bytes from the database.
	hash, err := chainhash.NewHashFromStr(hashHeight.String())
	if err != nil {
		log.ERROR(err)
}
	var blkBytes []byte
	err = d.Cx.RPCServer.Cfg.DB.View(func(dbTx database.Tx) error {
		var err error
		blkBytes, err = dbTx.FetchBlock(hash)
		return err
	})
	if err != nil {
		log.ERROR(err)
}
	// The verbose flag is set, so generate the JSON object and return it.
	// Deserialize the block.
	blk, err := util.NewBlockFromBytes(blkBytes)
	if err != nil {
		log.ERROR(err)
}
	// Get the block height from chain.
	blockHeight, err := d.Cx.RPCServer.Cfg.Chain.BlockHeightByHash(hash)
	if err != nil {
		log.ERROR(err)
}
	blk.SetHeight(blockHeight)
	params := d.Cx.RPCServer.Cfg.ChainParams
	blockHeader := &blk.MsgBlock().Header
	algoname := fork.GetAlgoName(blockHeader.Version, blockHeight)
	a := fork.GetAlgoVer(algoname, blockHeight)
	algoid := fork.GetAlgoID(algoname, blockHeight)
	// var value float64
	b.PowAlgoID = algoid
	b.Time = blockHeader.Timestamp.Unix()

	b.Height = int64(blockHeight)
	b.TxNum = len(blk.Transactions())
	b.Difficulty = rpc.GetDifficultyRatio(blockHeader.Bits, params, a)
	// txns := blk.Transactions()
	//
	// for _, tx := range txns {
	//	// Try to fetch the transaction from the memory pool and if that fails, try
	//	// the block database.
	//	var mtx *wire.MsgTx
	//
	//	// Look up the location of the transaction.
	//	blockRegion, err := b.rpc.Cfg.TxIndex.TxBlockRegion(tx.Hash())
	//	if err != nil {
		log.ERROR(err)
//	}
	//	if blockRegion == nil {
	//	}
	//	// Load the raw transaction bytes from the database.
	//	var txBytes []byte
	//	err = b.rpc.Cfg.DB.View(func(dbTx database.Tx) error {
	//		var err error
	//		txBytes, err = dbTx.FetchBlockRegion(blockRegion)
	//		return err
	//	})
	//	if err != nil {
		log.ERROR(err)
//	}
	//	// Deserialize the transaction
	//	var msgTx wire.MsgTx
	//	err = msgTx.Deserialize(bytes.NewReader(txBytes))
	//	if err != nil {
		log.ERROR(err)
//	}
	//	mtx = &msgTx
	//
	//	for _, vout := range rpc.CreateVoutList(mtx, b.rpc.Cfg.ChainParams, nil) {
	//
	//		value = value + vout.Value
	//	}
	//
	fmt.Println("Uzebekistanka malalalallalalaazsa")
	fmt.Println("Uzebekistanka malalalallalalaazsa")
	fmt.Println("Uzebekistanka malalalallalalaazsa")
	// fmt.Println("Uzebekistanka malalalallalalaazsa", b)
	fmt.Println("Uzebekistanka malalalallalalaazsa")
	// b.Amount = value
	// }
	return
}

func (d *DuOS) GetBlocksExcerpts(startBlock, blockHeight int) mod.DuOSblocks {
	for i := startBlock; i <= blockHeight; i++ {

		d.Data.Blocks = append(d.Data.Blocks, d.GetBlockExcerpt(i))
	}
	return d.Data.Blocks
}

// func (v *DuOSnode) Addnode(a *btcjson.AddNodeCmd) {
// 	r, err := v.Cx.RPCServer.HandleAddNode(v.Cx.RPCServer, a, nil)
// 	return
// }
// func (v *DuOSnode) Createrawtransaction(a *btcjson.CreateRawTransactionCmd) {
// 	r, err := v.Cx.RPCServer.HandleCreateRawTransaction(v.Cx.RPCServer, a, nil)
// 	r = ""
// 	return
// }
// func (v *DuOSnode) Decoderawtransaction(a *btcjson.DecodeRawTransactionCmd) {
// 	r, err := v.Cx.RPCServer.HandleDecodeRawTransaction(v.Cx.RPCServer, a, nil)
// 	r = btcjson.TxRawDecodeResult{}
// 	return
// }
// func (v *DuOSnode) Decodescript(a *btcjson.DecodeScriptCmd) {
// 	r, err := v.Cx.RPCServer.HandleDecodeScript(v.Cx.RPCServer, a, nil)
// 	return
// }
// func (v *DuOSnode) Estimatefee(a *btcjson.EstimateFeeCmd) {
// 	r, err := v.Cx.RPCServer.HandleEstimateFee(v.Cx.RPCServer, a, nil)
// 	r = 0.0
// 	return
// }
// func (v *DuOSnode) Generate(a *btcjson.GenerateCmd) {
// 	r, err := v.Cx.RPCServer.HandleGenerate(v.Cx.RPCServer, a, nil)
// 	r = []string{}
// 	return
// }
// func (v *DuOSnode) Getaddednodeinfo(a *btcjson.GetAddedNodeInfoCmd) {
// 	r, err := v.Cx.RPCServer.HandleGetAddedNodeInfo(v.Cx.RPCServer, a, nil)
// 	r = []string{}
// 	return
// }
// func (v *DuOSnode) Getbestblock() {
// 	r, err := v.Cx.RPCServer.HandleGetBestBlock(v.Cx.RPCServer, a, nil)
// 	r = btcjson.GetBestBlockResult{}
// 	return
// }
// func (v *DuOSnode) Getbestblockhash() {
// 	r, err := v.Cx.RPCServer.HandleGetBestBlockHash(v.Cx.RPCServer, a, nil)
// 	r = ""
// 	return
// }
// func (v *DuOSnode) Getblock(a *btcjson.GetBlockCmd) {
// 	r, err := v.Cx.RPCServer.HandleGetBlock(v.Cx.RPCServer, a, nil)
// 	r = btcjson.GetBlockVerboseResult{}
// 	return
// }
// func (d *DuOS) GetBlockChainInfo() {
//	getBlockChainInfo, err := rpc.HandleGetBlockChainInfo(d.Cx.RPCServer, nil, nil)
//	if err != nil {
//		log.ERROR(err)
//		d.PushDuOSalert("Error",err.Error(), "error")
//	}
//	var ok bool
//	d.Core.Node.BlockChainInfo, ok = getBlockChainInfo.(*btcjson.
//	GetBlockChainInfoResult)
//	if !ok {
//		d.Core.Node.BlockChainInfo = &btcjson.GetBlockChainInfoResult{}
//	}
//
// }

func (d *DuOS) GetBlockCount() int64 {
	getBlockCount, err := rpc.HandleGetBlockCount(d.Cx.RPCServer, nil, nil)
	if err != nil {
		log.ERROR(err)
d.PushDuOSalert("Error", err.Error(), "error")
	}
	d.Data.Status.BlockCount = getBlockCount.(int64)
	return d.Data.Status.BlockCount
}
func (d *DuOS) GetBlockHash(blockHeight int) string {
	hcmd := btcjson.GetBlockHashCmd{
		Index: int64(blockHeight),
	}
	hash, err := rpc.HandleGetBlockHash(d.Cx.RPCServer, &hcmd, nil)
	if err != nil {
		log.ERROR(err)
d.PushDuOSalert("Error", err.Error(), "error")
	}
	return hash.(string)
}
func (d *DuOS) GetBlock(hash string) btcjson.GetBlockVerboseResult {
	verbose, verbosetx := true, true
	bcmd := btcjson.GetBlockCmd{
		Hash:      hash,
		Verbose:   &verbose,
		VerboseTx: &verbosetx,
	}
	bl, err := rpc.HandleGetBlock(d.Cx.RPCServer, &bcmd, nil)
	if err != nil {
		log.ERROR(err)
d.PushDuOSalert("Error", err.Error(), "error")
	}
	return bl.(btcjson.GetBlockVerboseResult)
}

// func (v *DuOSnode) Getblockheader(a *btcjson.GetBlockHeaderCmd) {
// 	r, err := v.Cx.RPCServer.HandleGetBlockHeader(v.Cx.RPCServer, a, nil)
// 	r = btcjson.GetBlockHeaderVerboseResult{}
// 	return
// }

func (d *DuOS) GetConnectionCount() int32 {
	d.Data.Status.ConnectionCount = d.Cx.RPCServer.Cfg.ConnMgr.ConnectedCount()
	return d.Data.Status.ConnectionCount
}

func (d *DuOS) GetDifficulty() float64 {
	c := btcjson.GetDifficultyCmd{}
	r, err := rpc.HandleGetDifficulty(d.Cx.RPCServer, c, nil)
	if err != nil {
		log.ERROR(err)
d.PushDuOSalert("Error", err.Error(), "error")
	}
	d.Data.Status.Difficulty = r.(float64)
	return d.Data.Status.Difficulty
}

// func (v *DuOSnode) Gethashespersec() {
// 	r, err := v.Cx.RPCServer.HandleGetHashesPerSec(v.Cx.RPCServer, a, nil)
// 	r = int64(0)
// 	return
// }
// func (v *DuOSnode) Getheaders(a *btcjson.GetHeadersCmd) {
// 	r, err := v.Cx.RPCServer.HandleGetHeaders(v.Cx.RPCServer, a, nil)
// 	r = []string{}
// 	return
// }
// func (v *DuOSnode) Getinfo() {
// 	r, err := v.Cx.RPCServer.HandleGetInfo(v.Cx.RPCServer, a, nil)
// 	r = btcjson.InfoChainResult{}
// 	return
// }
// func (v *DuOSnode) Getmempoolinfo() {
// 	r, err := v.Cx.RPCServer.HandleGetMempoolInfo(v.Cx.RPCServer, a, nil)
// 	r = btcjson.GetMempoolInfoResult{}
// 	return
// }
// func (v *DuOSnode) Getmininginfo() {
// 	r, err := v.Cx.RPCServer.HandleGetMiningInfo(v.Cx.RPCServer, a, nil)
// 	r = btcjson.GetMiningInfoResult{}
// 	return
// }
// func (v *DuOSnode) Getnettotals() {
// 	r, err := v.Cx.RPCServer.HandleGetNetTotals(v.Cx.RPCServer, a, nil)
// 	r = btcjson.GetNetTotalsResult{}
// 	return
// }
// func (v *DuOSnode) Getnetworkhashps(a *btcjson.GetNetworkHashPSCmd) {
// 	r, err := v.Cx.RPCServer.HandleGetNetworkHashPS(v.Cx.RPCServer, a, nil)
// 	r = int64(0)
// 	return
// }
func (dV *DuOS) GetPeerInfo() []*btcjson.GetPeerInfoResult {
	getPeers, err := rpc.HandleGetPeerInfo(dV.Cx.RPCServer, nil, nil)
	if err != nil {
		log.ERROR(err)
dV.PushDuOSalert("Error", err.Error(), "error")
	}
	dV.Data.Peers = getPeers.([]*btcjson.GetPeerInfoResult)
	return dV.Data.Peers
}

// func (v *DuOSnode) Stop() {
// 	r, err := v.Cx.RPCServer.HandleStop(v.Cx.RPCServer, a, nil)
// 	r = ""
// 	return
// }
func (d *DuOS) Uptime() (r int64) {
	rRaw, err := rpc.HandleUptime(d.Cx.RPCServer, nil, nil)
	if err != nil {
		log.ERROR(err)
}
	// rRaw = int64(0)
	d.Data.Status.UpTime = rRaw.(int64)
	return d.Data.Status.UpTime
}

// func (v *DuOSnode) Validateaddress(a *btcjson.ValidateAddressCmd) {
// 	r, err := v.Cx.RPCServer.HandleValidateAddress(v.Cx.RPCServer, a, nil)
// 	r = btcjson.ValidateAddressChainResult{}
// 	return
// }
// func (v *DuOSnode) Verifychain(a *btcjson.VerifyChainCmd) {
// 	r, err := v.Cx.RPCServer.HandleVerifyChain(v.Cx.RPCServer, a, nil)
// }
// func (v *DuOSnode) Verifymessage(a *btcjson.VerifyMessageCmd) {
// 	r, err := v.Cx.RPCServer.HandleVerifyMessage(v.Cx.RPCServer, a, nil)
// 	r = ""
// 	return
// }
func (d *DuOS) GetWalletVersion() map[string]btcjson.VersionResult {
	v, err := rpc.HandleVersion(d.Cx.RPCServer, nil, nil)
	if err != nil {
		log.ERROR(err)
}
	return v.(map[string]btcjson.VersionResult)
}
