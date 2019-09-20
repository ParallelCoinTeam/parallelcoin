package vue

import (
	"fmt"
	"github.com/parallelcointeam/parallelcoin/cmd/node/rpc"
	"github.com/parallelcointeam/parallelcoin/pkg/chain/fork"
	chainhash "github.com/parallelcointeam/parallelcoin/pkg/chain/hash"
	database "github.com/parallelcointeam/parallelcoin/pkg/db"
	"github.com/parallelcointeam/parallelcoin/pkg/util"
)

func (d *DuoVUE) GetNetworkLastBlock() int32 {
	for _, g := range d.cx.RPCServer.Cfg.ConnMgr.ConnectedPeers() {
		l := g.ToPeer().StatsSnapshot().LastBlock
		if l > d.Status.NetworkLastBlock {
			d.Status.NetworkLastBlock = l
		}
	}
	return d.Status.NetworkLastBlock
}

//func (n *DuoVUEnode) GetBlocks() {
//	blks := []mod.Block{}
//	getBlockChain, err := rpc.HandleGetBlockChainInfo(n.rpc, nil, nil)
//	if err != nil {
//		alert.Alert.Time = time.Now()
//		alert.Alert.Alert = err.Error()
//		alert.Alert.AlertType = "error"
//	}
//	blockChain := getBlockChain.(*json.GetBlockChainInfoResult)
//	blockCount := int(blockChain.Blocks)
//	for ibh := blockCount; ibh >= 0; {
//		block := json.GetBlockVerboseResult{}
//		hcmd := json.GetBlockHashCmd{
//			Index: int64(ibh),
//		}
//		hash, err := rpc.HandleGetBlockHash(n.rpc, &hcmd, nil)
//		if err != nil {
//			alert.Alert.Time = time.Now()
//			alert.Alert.Alert = err.Error()
//			alert.Alert.AlertType = "error"
//		}
//		if hash != nil {
//			verbose, verbosetx := true, true
//			bcmd := json.GetBlockCmd{
//				Hash:      hash.(string),
//				Verbose:   &verbose,
//				VerboseTx: &verbosetx,
//			}
//			bl, err := rpc.HandleGetBlock(n.rpc, &bcmd, nil)
//			if err != nil {
//				alert.Alert.Time = time.Now()
//				alert.Alert.Alert = err.Error()
//				alert.Alert.AlertType = "error"
//			}
//			block = bl.(json.GetBlockVerboseResult)
//			blks = append(blks, mod.Block{
//				Hash:          block.Hash,
//				Confirmations: block.Confirmations,
//				Height:        block.Height,
//				TxNum:         block.TxNum,
//				Time:          block.Time,
//			})
//			ibh--
//		}
//	}
//
//	n.Blocks = blks
//}

//func (n *DuoVUEnode) GetBlocks(per, page int) {
//	blks := []json.GetBlockVerboseResult{}
//	getBlockChain, err := rpc.HandleGetBlockChainInfo(n.rpc, nil, nil)
//	if err != nil {
//		alert.Alert.Time = time.Now()
//		alert.Alert.Alert = err.Error()
//		alert.Alert.AlertType = "error"
//	}
//	blockChain := getBlockChain.(*json.GetBlockChainInfoResult)
//	blockCount := int(blockChain.Blocks)
//	startBlock := blockCount - per*page
//	minusBlockStart := int(startBlock + per)
//	for ibh := minusBlockStart; ibh >= startBlock; {
//		block := json.GetBlockVerboseResult{}
//		hcmd := json.GetBlockHashCmd{
//			Index: int64(ibh),
//		}
//		hash, err := rpc.HandleGetBlockHash(n.rpc, &hcmd, nil)
//		if err != nil {
//			alert.Alert.Time = time.Now()
//			alert.Alert.Alert = err.Error()
//			alert.Alert.AlertType = "error"
//		}
//		if hash != nil {
//			verbose, verbosetx := true, true
//			bcmd := json.GetBlockCmd{
//				Hash:      hash.(string),
//				Verbose:   &verbose,
//				VerboseTx: &verbosetx,
//			}
//			bl, err := rpc.HandleGetBlock(n.rpc, &bcmd, nil)
//			if err != nil {
//				alert.Alert.Time = time.Now()
//				alert.Alert.Alert = err.Error()
//				alert.Alert.AlertType = "error"
//			}
//			block = bl.(json.GetBlockVerboseResult)
//			blks = append(blks, block)
//			ibh--
//		}
//	}
//	n.Blocks.Blocks = blks
//	n.Blocks.CurrentPage = page
//	n.Blocks.PageCount = blockCount / per
//
//}

func (d DuoVUE) GetBlockExcerpt(height int) (b DuoVUEblock) {
	hashHeight, err := d.cx.RPCServer.Cfg.Chain.BlockHashByHeight(int32(height))
	if err != nil {
	}
	// Load the raw block bytes from the database.
	hash, err := chainhash.NewHashFromStr(hashHeight.String())
	if err != nil {
	}
	var blkBytes []byte
	err = d.cx.RPCServer.Cfg.DB.View(func(dbTx database.Tx) error {
		var err error
		blkBytes, err = dbTx.FetchBlock(hash)
		return err
	})
	if err != nil {
	}
	// The verbose flag is set, so generate the JSON object and return it.
	// Deserialize the block.
	blk, err := util.NewBlockFromBytes(blkBytes)
	if err != nil {
	}
	// Get the block height from chain.
	blockHeight, err := d.cx.RPCServer.Cfg.Chain.BlockHeightByHash(hash)
	if err != nil {
	}
	blk.SetHeight(blockHeight)
	params := d.cx.RPCServer.Cfg.ChainParams
	blockHeader := &blk.MsgBlock().Header
	algoname := fork.GetAlgoName(blockHeader.Version, blockHeight)
	a := fork.GetAlgoVer(algoname, blockHeight)
	algoid := fork.GetAlgoID(algoname, blockHeight)
	//var value float64
	b.PowAlgoID = algoid
	b.Time = blockHeader.Timestamp.Unix()

	b.Height = int64(blockHeight)
	b.TxNum = len(blk.Transactions())
	b.Difficulty = rpc.GetDifficultyRatio(blockHeader.Bits, params, a)
	//txns := blk.Transactions()
	//
	//for _, tx := range txns {
	//	// Try to fetch the transaction from the memory pool and if that fails, try
	//	// the block database.
	//	var mtx *wire.MsgTx
	//
	//	// Look up the location of the transaction.
	//	blockRegion, err := b.rpc.Cfg.TxIndex.TxBlockRegion(tx.Hash())
	//	if err != nil {
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
	//	}
	//	// Deserialize the transaction
	//	var msgTx wire.MsgTx
	//	err = msgTx.Deserialize(bytes.NewReader(txBytes))
	//	if err != nil {
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
	//fmt.Println("Uzebekistanka malalalallalalaazsa", b)
	fmt.Println("Uzebekistanka malalalallalalaazsa")
	//b.Amount = value
	//}
	return
}

func (d *DuoVUE) GetBlocksExcerpts(startBlock, blockHeight int) (b *DuoVUEchain) {
	for i := startBlock; i <= blockHeight; i++ {

		b.Blocks = append(b.Blocks, d.GetBlockExcerpt(i))
	}
	return
}
