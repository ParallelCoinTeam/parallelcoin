package rcd

import (
	"fmt"

	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/rpc/chainrpc"
)

func (r *RcVar) GetDuoUIhashesPerSec() {
	// r.Status.Wallet.Hashes = int64(r.cx.RPCServer.Cfg.CPUMiner.HashesPerSecond())
	Debug("centralise hash function stuff here") // cpuminer
	r.Status.Kopach.Hashrate = r.cx.Hashrate.Load()
	return
}
func (r *RcVar) GetDuoUInetworkHashesPerSec() {
	networkHashesPerSecIface, err := chainrpc.HandleGetNetworkHashPS(r.cx.RPCServer, btcjson.NewGetNetworkHashPSCmd(nil, nil), nil)
	if err != nil {
	}
	networkHashesPerSec, ok := networkHashesPerSecIface.(int64)
	if !ok {
	}
	r.Status.Node.NetHash.Store(uint64(networkHashesPerSec))
	return
}
func (r *RcVar) GetDuoUIblockHeight() {
	r.Status.Node.BlockHeight.Store(
		uint64(r.cx.RPCServer.Cfg.Chain.BestSnapshot().Height))
	return
}
func (r *RcVar) GetDuoUIbestBlockHash() {
	r.Status.Node.BestBlock.Store(r.cx.RPCServer.Cfg.Chain.BestSnapshot().Hash.String())
	return
}
func (r *RcVar) GetDuoUIdifficulty() {
	r.Status.Node.Difficulty.Store(
		chainrpc.GetDifficultyRatio(
			r.cx.RPCServer.Cfg.Chain.BestSnapshot().Bits, r.cx.RPCServer.Cfg.ChainParams, 2))
	return
}
func (r *RcVar) GetDuoUIblockCount() {
	getBlockCount, err := chainrpc.HandleGetBlockCount(r.cx.RPCServer, nil, nil)
	if err != nil {
		// r.PushDuoUIalert("BTCJSONError", err.BTCJSONError(), "error")
	}
	r.Status.Node.BlockCount.Store(uint64(getBlockCount.(int64)))
	// Info(getBlockCount)
	return
}
func (r *RcVar) GetDuoUInetworkLastBlock() {
	for _, g := range r.cx.RPCServer.Cfg.ConnMgr.ConnectedPeers() {
		l := g.ToPeer().StatsSnapshot().LastBlock
		if l > r.Status.Node.NetworkLastBlock.Load() {
			r.Status.Node.NetworkLastBlock.Store(l)
		}
	}
	return
}
func (r *RcVar) GetDuoUIconnectionCount() {
	r.Status.Node.ConnectionCount.Store(r.cx.RealNode.ConnectedCount())
	return
}

func (r *RcVar) GetDuoUIhashesPerSecList() {
	// // Create a new ring of size 5
	// hps := ring.New(3)
	// //GetDuoUIhashesPerSec
	// // Get the length of the ring
	// n := hps.Len()
	//
	// // Initialize the ring with some integer values
	// for i := 0; i < n; i++ {
	r.GetDuoUIhashesPerSec()
	// hps.Value = r.Status.Kopach.Hashrate
	//	hps = hps.Next()
	// }
	//
	// // Iterate through the ring and print its contents
	// hps.Do(func(p interface{}) {
	//	r.Status.Kopach.Hps = append(r.Status.Kopach.Hps, p.(float64))
	//
	fmt.Println(r.Status.Kopach.Hashrate)

	// })

}

func (r *RcVar) GetPeerInfo() func() {
	return func() {

		getPeers, err := chainrpc.HandleGetPeerInfo(r.cx.RPCServer, nil, nil)
		if err != nil {
			// dV.PushDuoVUEalert("BTCJSONError", err.BTCJSONError(), "error")
		}
		r.Network.Peers = getPeers.([]*btcjson.GetPeerInfoResult)
	}
}
