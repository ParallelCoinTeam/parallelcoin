package rcd

import (
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/p9c/pod/pkg/conte"
	"time"
)

// System Ststus
type
	DuOStatus struct {
		Version       string                           `json:"ver"`
		WalletVersion map[string]btcjson.VersionResult `json:"walletver"`
		UpTime        int64                            `json:"uptime"`
		CurrentNet    string                           `json:"net"`
		Chain         string                           `json:"chain"`
	}
type
	DuOShashes struct{ int64 }
type
	DuOSnetworkHash struct{ int64 }
type
	DuOSheight struct{ int32 }
type
	DuOSbestBlockHash struct{ string }
type
	DuOSdifficulty struct{ float64 }

//type
// MempoolInfo      struct { string}
type
	DuOSblockCount struct{ int64 }
type
	DuOSnetLastBlock struct{ int32 }
type
	DuOSconnections struct{ int32 }
type
	DuOSlocalHost struct {
		Cpu        []cpu.InfoStat        `json:"cpu"`
		CpuPercent []float64             `json:"cpupercent"`
		Memory     mem.VirtualMemoryStat `json:"mem"`
		Disk       disk.UsageStat        `json:"disk"`
	}

func
(r *RcVar) GetDuOStatus(cx *conte.Xt) {
	r.Status = *new(DuOStatus)
	v, err := rpc.HandleVersion(cx.RPCServer, nil, nil)
	if err != nil {
	}
	r.Status.Version = "0.0.1"
	r.Status.WalletVersion = v.(map[string]btcjson.VersionResult)
	r.Status.UpTime = time.Now().Unix() - cx.RPCServer.Cfg.StartupTime
	r.Status.CurrentNet = cx.RPCServer.Cfg.ChainParams.Net.String()
	r.Status.Chain = cx.RPCServer.Cfg.ChainParams.Name
	return
}
func
(r *RcVar) GetDuOShashesPerSec(cx *conte.Xt) {
	r.Hashes = int64(cx.RPCServer.Cfg.CPUMiner.HashesPerSecond())
	return
}
func
(r *RcVar) GetDuOSnetworkHashesPerSec(cx *conte.Xt) {
	networkHashesPerSecIface, err := rpc.HandleGetNetworkHashPS(cx.RPCServer, btcjson.NewGetNetworkHashPSCmd(nil, nil), nil)
	if err != nil {
	}
	networkHashesPerSec, ok := networkHashesPerSecIface.(int64)
	if !ok {
	}
	r.NetHash = networkHashesPerSec
	return
}
func
(r *RcVar) GetDuOSblockHeight(cx *conte.Xt) {
	r.BlockHeight = cx.RPCServer.Cfg.Chain.BestSnapshot().Height
	return
}
func
(r *RcVar) GetDuOSbestBlockHash(cx *conte.Xt) {
	r.BestBlock = cx.RPCServer.Cfg.Chain.BestSnapshot().Hash.String()
	return
}
func
(r *RcVar) GetDuOSdifficulty(cx *conte.Xt) {
	r.Difficulty = rpc.GetDifficultyRatio(cx.RPCServer.Cfg.Chain.BestSnapshot().Bits, cx.RPCServer.Cfg.ChainParams, 2)
	return
}
func
(r *RcVar) GetDuOSblockCount(cx *conte.Xt) {
	getBlockCount, err := rpc.HandleGetBlockCount(cx.RPCServer, nil, nil)
	if err != nil {
		r.PushDuOSalert("Error", err.Error(), "error")
	}
	r.BlockCount = getBlockCount.(int64)
	return
}
func
(r *RcVar) GetDuOSnetworkLastBlock(cx *conte.Xt) {
	for _, g := range cx.RPCServer.Cfg.ConnMgr.ConnectedPeers() {
		l := g.ToPeer().StatsSnapshot().LastBlock
		if l > r.NetLastBlock {
			r.NetLastBlock = l
		}
	}
	return
}
func
(r *RcVar) GetDuOSconnectionCount(cx *conte.Xt) {
	r.Connections = cx.RPCServer.Cfg.ConnMgr.ConnectedCount()
	return
}
func
(r *RcVar) GetDuOSlocalLost(cx *conte.Xt) {
	r.Localhost = *new(DuOSlocalHost)
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
