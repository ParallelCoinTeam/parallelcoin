package rcd

import (
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"time"
)

// System Ststus

func
(rc *RcVar) GetDuoUIstatus(duo *models.DuoUI, cx *conte.Xt) {
	rc.Status = *new(models.DuoUIstatus)
	v, err := rpc.HandleVersion(cx.RPCServer, nil, nil)
	if err != nil {
	}
	rc.Status.Version = "0.0.1"
	rc.Status.WalletVersion = v.(map[string]btcjson.VersionResult)
	rc.Status.UpTime = time.Now().Unix() - cx.RPCServer.Cfg.StartupTime
	rc.Status.CurrentNet = cx.RPCServer.Cfg.ChainParams.Net.String()
	rc.Status.Chain = cx.RPCServer.Cfg.ChainParams.Name
	return
}
func
(rc *RcVar) GetDuoUIhashesPerSec(duo *models.DuoUI, cx *conte.Xt) {
	rc.Hashes = int64(cx.RPCServer.Cfg.CPUMiner.HashesPerSecond())
	return
}
func
(rc *RcVar) GetDuoUInetworkHashesPerSec(duo *models.DuoUI, cx *conte.Xt) {
	networkHashesPerSecIface, err := rpc.HandleGetNetworkHashPS(cx.RPCServer, btcjson.NewGetNetworkHashPSCmd(nil, nil), nil)
	if err != nil {
	}
	networkHashesPerSec, ok := networkHashesPerSecIface.(int64)
	if !ok {
	}
	rc.NetHash = networkHashesPerSec
	return
}
func
(rc *RcVar) GetDuoUIblockHeight(duo *models.DuoUI, cx *conte.Xt) {
	rc.BlockHeight = cx.RPCServer.Cfg.Chain.BestSnapshot().Height
	return
}
func
(rc *RcVar) GetDuoUIbestBlockHash(duo *models.DuoUI, cx *conte.Xt) {
	rc.BestBlock = cx.RPCServer.Cfg.Chain.BestSnapshot().Hash.String()
	return
}
func
(rc *RcVar) GetDuoUIdifficulty(duo *models.DuoUI, cx *conte.Xt) {
	rc.Difficulty = rpc.GetDifficultyRatio(cx.RPCServer.Cfg.Chain.BestSnapshot().Bits, cx.RPCServer.Cfg.ChainParams, 2)
	return
}
func
(rc *RcVar) GetDuoUIblockCount(duo *models.DuoUI, cx *conte.Xt) {
	getBlockCount, err := rpc.HandleGetBlockCount(cx.RPCServer, nil, nil)
	if err != nil {
		rc.PushDuoUIalert("Error", err.Error(), "error")
	}
	rc.BlockCount = getBlockCount.(int64)
	return
}
func
(rc *RcVar) GetDuoUInetworkLastBlock(duo *models.DuoUI, cx *conte.Xt) {
	for _, g := range cx.RPCServer.Cfg.ConnMgr.ConnectedPeers() {
		l := g.ToPeer().StatsSnapshot().LastBlock
		if l > rc.NetworkLastBlock {
			rc.NetworkLastBlock = l
		}
	}
	return
}
func
(rc *RcVar) GetDuoUIconnectionCount(duo *models.DuoUI, cx *conte.Xt) {
	rc.ConnectionCount = cx.RPCServer.Cfg.ConnMgr.ConnectedCount()
	return
}
func
(rc *RcVar) GetDuoUIlocalLost(duo *models.DuoUI) {
	rc.Localhost = *new(models.DuoUIlocalHost)
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
