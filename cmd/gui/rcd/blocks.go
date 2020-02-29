package rcd

import (
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

func (r *RcVar) GetSingleBlock(hash string) func() {
	return func() {
		r.SingleBlock = r.GetBlock(hash)
	}
}

func (r *RcVar) GetBlock(hash string) btcjson.GetBlockVerboseResult {
	verbose, verbosetx := true, true
	bcmd := btcjson.GetBlockCmd{
		Hash:      hash,
		Verbose:   &verbose,
		VerboseTx: &verbosetx,
	}
	bl, err := rpc.HandleGetBlock(r.cx.RPCServer, &bcmd, nil)
	if err != nil {
		//dv.PushDuoVUEalert("Error", err.Error(), "error")
	}
	return bl.(btcjson.GetBlockVerboseResult)
}

func (r *RcVar) GetNetworkLastBlock() int32 {
	for _, g := range r.cx.RPCServer.Cfg.ConnMgr.ConnectedPeers() {
		l := g.ToPeer().StatsSnapshot().LastBlock
		if l > r.Status.Node.NetworkLastBlock {
			r.Status.Node.NetworkLastBlock = l
		}
	}
	return r.Status.Node.NetworkLastBlock
}

func (r *RcVar) GetBlockExcerpt(height int) (b model.DuoUIblock) {
	b = *new(model.DuoUIblock)
	hashHeight, err := r.cx.RPCServer.Cfg.Chain.BlockHashByHeight(int32(height))
	if err != nil {
		log.ERROR("Block Hash By Height:", err)
	}

	verbose, verbosetx := true, true
	bcmd := btcjson.GetBlockCmd{
		Hash:      hashHeight.String(),
		Verbose:   &verbose,
		VerboseTx: &verbosetx,
	}
	bl, err := rpc.HandleGetBlock(r.cx.RPCServer, &bcmd, nil)
	if err != nil {
		//dv.PushDuoVUEalert("Error", err.Error(), "error")
	}
	block := bl.(btcjson.GetBlockVerboseResult)
	b.Height = block.Height
	b.BlockHash = block.Hash
	b.Confirmations = block.Confirmations
	b.TxNum = block.TxNum
	b.Time = string(block.Time)
	b.Link = &controller.Button{}
	return
}

func (r *RcVar) GetBlocksExcerpts(page, perPage int) func() {
	return func() {
		//pages := r.Status.Node.BlockCount / perPage
		startBlock := page * perPage
		endBlock := page*perPage + perPage

		blocks := *new([]model.DuoUIblock)
		for i := startBlock; i <= endBlock; i++ {
			blocks = append(blocks, r.GetBlockExcerpt(i))
			log.INFO("trazo")
			log.INFO(r.Status.Node.BlockHeight)
		}
		r.Blocks = blocks
		return
	}
}

func (r *RcVar) GetBlockCount() {
	getBlockCount, err := rpc.HandleGetBlockCount(r.cx.RPCServer, nil, nil)
	if err != nil {
		//dv.PushDuoVUEalert("Error", err.Error(), "error")
	}
	r.Status.Node.BlockCount = int(getBlockCount.(int64))
	return
}
func (r *RcVar) GetBlockHash(blockHeight int) string {
	hcmd := btcjson.GetBlockHashCmd{
		Index: int64(blockHeight),
	}
	hash, err := rpc.HandleGetBlockHash(r.cx.RPCServer, &hcmd, nil)
	if err != nil {
		//dv.PushDuoVUEalert("Error", err.Error(), "error")
	}
	return hash.(string)
}

func (r *RcVar) GetConnectionCount() {
	r.Status.Node.ConnectionCount = r.cx.RPCServer.Cfg.ConnMgr.ConnectedCount()
	return
}

func (r *RcVar) GetDifficulty() {
	c := btcjson.GetDifficultyCmd{}
	diff, err := rpc.HandleGetDifficulty(r.cx.RPCServer, c, nil)
	if err != nil {
		//dv.PushDuoVUEalert("Error", err.Error(), "error")
	}
	r.Status.Node.Difficulty = diff.(float64)
	return
}

// func (v *DuoVUEnode) Gethashespersec() {
// 	r, err := v.r.cx.RPCServer.HandleGetHashesPerSec(v.r.cx.RPCServer, a, nil)
// 	r = int64(0)
// 	return
// }
// func (v *DuoVUEnode) Getheaders(a *btcjson.GetHeadersCmd) {
// 	r, err := v.r.cx.RPCServer.HandleGetHeaders(v.r.cx.RPCServer, a, nil)
// 	r = []string{}
// 	return
// }
// func (v *DuoVUEnode) Getinfo() {
// 	r, err := v.r.cx.RPCServer.HandleGetInfo(v.r.cx.RPCServer, a, nil)
// 	r = btcjson.InfoChainResult{}
// 	return
// }
// func (v *DuoVUEnode) Getmempoolinfo() {
// 	r, err := v.r.cx.RPCServer.HandleGetMempoolInfo(v.r.cx.RPCServer, a, nil)
// 	r = btcjson.GetMempoolInfoResult{}
// 	return
// }
// func (v *DuoVUEnode) Getmininginfo() {
// 	r, err := v.r.cx.RPCServer.HandleGetMiningInfo(v.r.cx.RPCServer, a, nil)
// 	r = btcjson.GetMiningInfoResult{}
// 	return
// }
// func (v *DuoVUEnode) Getnettotals() {
// 	r, err := v.r.cx.RPCServer.HandleGetNetTotals(v.r.cx.RPCServer, a, nil)
// 	r = btcjson.GetNetTotalsResult{}
// 	return
// }
// func (v *DuoVUEnode) Getnetworkhashps(a *btcjson.GetNetworkHashPSCmd) {
// 	r, err := v.r.cx.RPCServer.HandleGetNetworkHashPS(v.r.cx.RPCServer, a, nil)
// 	r = int64(0)
// 	return
// }
func (r *RcVar) GetPeerInfo() {
	getPeers, err := rpc.HandleGetPeerInfo(r.cx.RPCServer, nil, nil)
	if err != nil {
		//dV.PushDuoVUEalert("Error", err.Error(), "error")
	}
	r.Peers = getPeers.([]*btcjson.GetPeerInfoResult)
	return
}

// func (v *DuoVUEnode) Stop() {
// 	r, err := v.r.cx.RPCServer.HandleStop(v.r.cx.RPCServer, a, nil)
// 	r = ""
// 	return
// }
func (r *RcVar) GetUptime() {
	rRaw, err := rpc.HandleUptime(r.cx.RPCServer, nil, nil)
	if err != nil {
	}
	// rRaw = int64(0)
	r.Uptime = rRaw.(int)
	return
}

// func (v *DuoVUEnode) Validateaddress(a *btcjson.ValidateAddressCmd) {
// 	r, err := v.r.cx.RPCServer.HandleValidateAddress(v.r.cx.RPCServer, a, nil)
// 	r = btcjson.ValidateAddressChainResult{}
// 	return
// }
// func (v *DuoVUEnode) Verifychain(a *btcjson.VerifyChainCmd) {
// 	r, err := v.r.cx.RPCServer.HandleVerifyChain(v.r.cx.RPCServer, a, nil)
// }
// func (v *DuoVUEnode) Verifymessage(a *btcjson.VerifyMessageCmd) {
// 	r, err := v.r.cx.RPCServer.HandleVerifyMessage(v.r.cx.RPCServer, a, nil)
// 	r = ""
// 	return
// }
func (r *RcVar) GetWalletVersion() map[string]btcjson.VersionResult {
	v, err := rpc.HandleVersion(r.cx.RPCServer, nil, nil)
	if err != nil {
	}
	return v.(map[string]btcjson.VersionResult)
}
