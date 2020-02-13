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
(rc *RcVar) GetDuoUIstatus() {
	rc.Status = new(model.DuoUIstatus)
	v, err := rpc.HandleVersion(rc.Cx.RPCServer, nil, nil)
	if err != nil {
	}
	rc.Status.Version = "0.0.1"
	rc.Status.Wallet.WalletVersion = v.(map[string]btcjson.VersionResult)
	rc.Status.UpTime = time.Now().Unix() - rc.Cx.RPCServer.Cfg.StartupTime
	rc.Status.CurrentNet = rc.Cx.RPCServer.Cfg.ChainParams.Net.String()
	rc.Status.Chain = rc.Cx.RPCServer.Cfg.ChainParams.Name
	return
}
func
(rc *RcVar) GetDuoUIhashesPerSec() {
	rc.Status.Wallet.Hashes = int64(rc.Cx.RPCServer.Cfg.CPUMiner.HashesPerSecond())
	return
}
func
(rc *RcVar) GetDuoUInetworkHashesPerSec() {
	networkHashesPerSecIface, err := rpc.HandleGetNetworkHashPS(rc.Cx.RPCServer, btcjson.NewGetNetworkHashPSCmd(nil, nil), nil)
	if err != nil {
	}
	networkHashesPerSec, ok := networkHashesPerSecIface.(int64)
	if !ok {
	}
	rc.Status.Node.NetHash = networkHashesPerSec
	return
}
func
(rc *RcVar) GetDuoUIblockHeight() {
	rc.Status.Node.BlockHeight = rc.Cx.RPCServer.Cfg.Chain.BestSnapshot().Height
	return
}
func
(rc *RcVar) GetDuoUIbestBlockHash() {
	rc.Status.Node.BestBlock = rc.Cx.RPCServer.Cfg.Chain.BestSnapshot().Hash.String()
	return
}
func
(rc *RcVar) GetDuoUIdifficulty() {
	rc.Status.Node.Difficulty = rpc.GetDifficultyRatio(rc.Cx.RPCServer.Cfg.Chain.BestSnapshot().Bits, rc.Cx.RPCServer.Cfg.ChainParams, 2)
	return
}
func
(rc *RcVar) GetDuoUIblockCount() {
	getBlockCount, err := rpc.HandleGetBlockCount(rc.Cx.RPCServer, nil, nil)
	if err != nil {
		//rc.PushDuoUIalert("Error", err.Error(), "error")
	}
	rc.Status.Node.BlockCount = getBlockCount.(int64)
	return
}
func
(rc *RcVar) GetDuoUInetworkLastBlock() {
	for _, g := range rc.Cx.RPCServer.Cfg.ConnMgr.ConnectedPeers() {
		l := g.ToPeer().StatsSnapshot().LastBlock
		if l > rc.Status.Node.NetworkLastBlock {
			rc.Status.Node.NetworkLastBlock = l
		}
	}
	return
}
func
(rc *RcVar) GetDuoUIconnectionCount() {
	rc.Status.Node.ConnectionCount = rc.Cx.RPCServer.Cfg.ConnMgr.ConnectedCount()
	return
}
func
(rc *RcVar) GetDuoUIlocalLost() {
	rc.Localhost = *new(model.DuoUIlocalHost)
	sm, _ := mem.VirtualMemory()
	sc, _ := cpu.Info()
	sp, _ := cpu.Percent(0, true)
	sd, _ := disk.Usage("/")
	rc.Localhost.Cpu = sc
	rc.Localhost.CpuPercent = sp
	rc.Localhost.Memory = *sm
	rc.Localhost.Disk = *sd
	return
}
