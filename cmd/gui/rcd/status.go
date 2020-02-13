package rcd

import (
	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"time"
)

// System Ststus

func
(r *RcVar) GetDuoUIstatus() {
	v, err := rpc.HandleVersion(r.Cx.RPCServer, nil, nil)
	if err != nil {
	}
	r.Status.Version = "0.0.1"
	r.Status.Wallet.WalletVersion = v.(map[string]btcjson.VersionResult)
	r.Status.UpTime = time.Now().Unix() - r.Cx.RPCServer.Cfg.StartupTime
	r.Status.CurrentNet = r.Cx.RPCServer.Cfg.ChainParams.Net.String()
	r.Status.Chain = r.Cx.RPCServer.Cfg.ChainParams.Name
	return
}
func
(r *RcVar) GetDuoUIhashesPerSec() {
	r.Status.Wallet.Hashes = int64(r.Cx.RPCServer.Cfg.CPUMiner.HashesPerSecond())
	return
}
func
(r *RcVar) GetDuoUInetworkHashesPerSec() {
	networkHashesPerSecIface, err := rpc.HandleGetNetworkHashPS(r.Cx.RPCServer, btcjson.NewGetNetworkHashPSCmd(nil, nil), nil)
	if err != nil {
	}
	networkHashesPerSec, ok := networkHashesPerSecIface.(int64)
	if !ok {
	}
	r.Status.Node.NetHash = networkHashesPerSec
	return
}
func
(r *RcVar) GetDuoUIblockHeight() {
	r.Status.Node.BlockHeight = r.Cx.RPCServer.Cfg.Chain.BestSnapshot().Height
	return
}
func
(r *RcVar) GetDuoUIbestBlockHash() {
	r.Status.Node.BestBlock = r.Cx.RPCServer.Cfg.Chain.BestSnapshot().Hash.String()
	return
}
func
(r *RcVar) GetDuoUIdifficulty() {
	r.Status.Node.Difficulty = rpc.GetDifficultyRatio(r.Cx.RPCServer.Cfg.Chain.BestSnapshot().Bits, r.Cx.RPCServer.Cfg.ChainParams, 2)
	return
}
func
(r *RcVar) GetDuoUIblockCount() {
	getBlockCount, err := rpc.HandleGetBlockCount(r.Cx.RPCServer, nil, nil)
	if err != nil {
		//r.PushDuoUIalert("Error", err.Error(), "error")
	}
	r.Status.Node.BlockCount = getBlockCount.(int64)
	return
}
func
(r *RcVar) GetDuoUInetworkLastBlock() {
	for _, g := range r.Cx.RPCServer.Cfg.ConnMgr.ConnectedPeers() {
		l := g.ToPeer().StatsSnapshot().LastBlock
		if l > r.Status.Node.NetworkLastBlock {
			r.Status.Node.NetworkLastBlock = l
		}
	}
	return
}
func
(r *RcVar) GetDuoUIconnectionCount() {
	r.Status.Node.ConnectionCount = r.Cx.RPCServer.Cfg.ConnMgr.ConnectedCount()
	return
}
func
(r *RcVar) GetDuoUIlocalLost() {
	r.Localhost = *new(model.DuoUIlocalHost)
	sm, _ := mem.VirtualMemory()
	sc, _ := cpu.Info()
	sp, _ := cpu.Percent(0, true)
	sd, _ := disk.Usage("/")
	r.Localhost.Cpu = sc
	r.Localhost.CpuPercent = sp
	r.Localhost.Memory = *sm
	r.Localhost.Disk = *sd
	return
}
