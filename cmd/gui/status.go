package gui

import (
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
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
(r *rcvar) GetDuOStatus() {
	r.status = *new(DuOStatus)
	v, err := rpc.HandleVersion(r.cx.RPCServer, nil, nil)
	if err != nil {
	}
	r.status.Version = "0.0.1"
	r.status.WalletVersion = v.(map[string]btcjson.VersionResult)
	r.status.UpTime = time.Now().Unix() - r.cx.RPCServer.Cfg.StartupTime
	r.status.CurrentNet = r.cx.RPCServer.Cfg.ChainParams.Net.String()
	r.status.Chain = r.cx.RPCServer.Cfg.ChainParams.Name
	return
}
func
(r *rcvar) GetDuOShashesPerSec() {
	r.hashes = int64(r.cx.RPCServer.Cfg.CPUMiner.HashesPerSecond())
	return
}
func
(r *rcvar) GetDuOSnetworkHashesPerSec() {
	networkHashesPerSecIface, err := rpc.HandleGetNetworkHashPS(r.cx.RPCServer, btcjson.NewGetNetworkHashPSCmd(nil, nil), nil)
	if err != nil {
	}
	networkHashesPerSec, ok := networkHashesPerSecIface.(int64)
	if !ok {
	}
	r.nethash = networkHashesPerSec
	return
}
func
(r *rcvar) GetDuOSheight() {
	r.height = r.cx.RPCServer.Cfg.Chain.BestSnapshot().Height
	return
}
func
(r *rcvar) GetDuOSbestBlockHash() {
	r.bestblock = r.cx.RPCServer.Cfg.Chain.BestSnapshot().Hash.String()
	return
}
func
(r *rcvar) GetDuOSdifficulty() {
	r.difficulty = rpc.GetDifficultyRatio(r.cx.RPCServer.Cfg.Chain.BestSnapshot().Bits, r.cx.RPCServer.Cfg.ChainParams, 2)
	return
}
func
(r *rcvar) GetDuOSblockCount() {
	getBlockCount, err := rpc.HandleGetBlockCount(r.cx.RPCServer, nil, nil)
	if err != nil {
		r.PushDuOSalert("Error", err.Error(), "error")
	}
	r.blockcount = getBlockCount.(int64)
	return
}
func
(r *rcvar) GetDuOSnetworkLastBlock() {
	for _, g := range r.cx.RPCServer.Cfg.ConnMgr.ConnectedPeers() {
		l := g.ToPeer().StatsSnapshot().LastBlock
		if l > r.netlastblock {
			r.netlastblock = l
		}
	}
	return
}
func
(r *rcvar) GetDuOSconnectionCount() {
	r.connections = r.cx.RPCServer.Cfg.ConnMgr.ConnectedCount()
	return
}
func
(r *rcvar) GetDuOSlocalLost() {
	r.localhost = *new(DuOSlocalHost)
	sm, _ := mem.VirtualMemory()
	sc, _ := cpu.Info()
	sp, _ := cpu.Percent(0, true)
	sd, _ := disk.Usage("/")
	r.localhost.Cpu = sc
	r.localhost.CpuPercent = sp
	r.localhost.Memory = *sm
	r.localhost.Disk = *sd
	return
}
