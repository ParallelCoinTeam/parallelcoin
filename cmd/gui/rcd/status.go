package rcd

import (
	"fmt"
	"time"

	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

// System Ststus

func (r *RcVar) GetDuoUIstatus() {
	v, err := rpc.HandleVersion(r.cx.RPCServer, nil, nil)
	if err != nil {
	}
	r.Status.Version = "0.0.1"
	r.Status.Wallet.WalletVersion = v.(map[string]btcjson.VersionResult)
	r.Status.UpTime = time.Now().Unix() - r.cx.RPCServer.Cfg.StartupTime
	r.Status.CurrentNet = r.cx.RPCServer.Cfg.ChainParams.Net.String()
	r.Status.Chain = r.cx.RPCServer.Cfg.ChainParams.Name
	return
}
func (r *RcVar) GetDuoUIhashesPerSec() {
	// r.Status.Wallet.Hashes = int64(r.cx.RPCServer.Cfg.CPUMiner.HashesPerSecond())
	log.DEBUG("centralise hash function stuff here") // cpuminer
	r.Status.Kopach.Hashrate = r.cx.Hashrate.Load()
	return
}
func (r *RcVar) GetDuoUInetworkHashesPerSec() {
	networkHashesPerSecIface, err := rpc.HandleGetNetworkHashPS(r.cx.RPCServer, btcjson.NewGetNetworkHashPSCmd(nil, nil), nil)
	if err != nil {
	}
	networkHashesPerSec, ok := networkHashesPerSecIface.(int64)
	if !ok {
	}
	r.Status.Node.NetHash = networkHashesPerSec
	return
}
func (r *RcVar) GetDuoUIblockHeight() {
	r.Status.Node.BlockHeight = int(r.cx.RPCServer.Cfg.Chain.BestSnapshot().Height)
	return
}
func (r *RcVar) GetDuoUIbestBlockHash() {
	r.Status.Node.BestBlock = r.cx.RPCServer.Cfg.Chain.BestSnapshot().Hash.String()
	return
}
func (r *RcVar) GetDuoUIdifficulty() {
	r.Status.Node.Difficulty = rpc.GetDifficultyRatio(r.cx.RPCServer.Cfg.Chain.BestSnapshot().Bits, r.cx.RPCServer.Cfg.ChainParams, 2)
	return
}
func (r *RcVar) GetDuoUIblockCount() {
	getBlockCount, err := rpc.HandleGetBlockCount(r.cx.RPCServer, nil, nil)
	if err != nil {
		//r.PushDuoUIalert("Error", err.Error(), "error")
	}
	r.Status.Node.BlockCount = int(getBlockCount.(int64))
	// log.INFO(getBlockCount)
	return
}
func (r *RcVar) GetDuoUInetworkLastBlock() {
	for _, g := range r.cx.RPCServer.Cfg.ConnMgr.ConnectedPeers() {
		l := g.ToPeer().StatsSnapshot().LastBlock
		if l > r.Status.Node.NetworkLastBlock {
			r.Status.Node.NetworkLastBlock = l
		}
	}
	return
}
func (r *RcVar) GetDuoUIconnectionCount() {
	r.Status.Node.ConnectionCount = r.cx.RPCServer.Cfg.ConnMgr.ConnectedCount()
	return
}
func (r *RcVar) GetDuoUIlocalLost() {
	r.Localhost = *new(model.DuoUIlocalHost)
	//sm, _ := mem.VirtualMemory()
	//sc, _ := cpu.Info()
	//sp, _ := cpu.Percent(0, true)
	//sd, _ := disk.Usage("/")
	//r.Localhost.Cpu = sc
	//r.Localhost.CpuPercent = sp
	//r.Localhost.Memory = *sm
	//r.Localhost.Disk = *sd
	return
}

func (r *RcVar) GetDuoUIhashesPerSecList() {
	//// Create a new ring of size 5
	//hps := ring.New(3)
	////GetDuoUIhashesPerSec
	//// Get the length of the ring
	//n := hps.Len()
	//
	//// Initialize the ring with some integer values
	//for i := 0; i < n; i++ {
	r.GetDuoUIhashesPerSec()
	//hps.Value = r.Status.Kopach.Hashrate
	//	hps = hps.Next()
	//}
	//
	//// Iterate through the ring and print its contents
	//hps.Do(func(p interface{}) {
	//	r.Status.Kopach.Hps = append(r.Status.Kopach.Hps, p.(float64))
	//
	fmt.Println(r.Status.Kopach.Hashrate)

	//})

}
