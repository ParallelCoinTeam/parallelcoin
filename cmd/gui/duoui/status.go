package duoui

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
(duo *DuoUI) GetDuOStatus() {
	duo.rc.Status = *new(DuOStatus)
	v, err := rpc.HandleVersion(duo.cx.RPCServer, nil, nil)
	if err != nil {
	}
	duo.rc.Status.Version = "0.0.1"
	duo.rc.Status.WalletVersion = v.(map[string]btcjson.VersionResult)
	duo.rc.Status.UpTime = time.Now().Unix() - duo.cx.RPCServer.Cfg.StartupTime
	duo.rc.Status.CurrentNet = duo.cx.RPCServer.Cfg.ChainParams.Net.String()
	duo.rc.Status.Chain = duo.cx.RPCServer.Cfg.ChainParams.Name
	return
}
func
(duo *DuoUI) GetDuOShashesPerSec() {
	duo.rc.Hashes = int64(duo.cx.RPCServer.Cfg.CPUMiner.HashesPerSecond())
	return
}
func
(duo *DuoUI) GetDuOSnetworkHashesPerSec() {
	networkHashesPerSecIface, err := rpc.HandleGetNetworkHashPS(duo.cx.RPCServer, btcjson.NewGetNetworkHashPSCmd(nil, nil), nil)
	if err != nil {
	}
	networkHashesPerSec, ok := networkHashesPerSecIface.(int64)
	if !ok {
	}
	duo.rc.NetHash = networkHashesPerSec
	return
}
func
(duo *DuoUI) GetDuOSblockHeight() {
	duo.rc.BlockHeight = duo.cx.RPCServer.Cfg.Chain.BestSnapshot().Height
	return
}
func
(duo *DuoUI) GetDuOSbestBlockHash() {
	duo.rc.BestBlock = duo.cx.RPCServer.Cfg.Chain.BestSnapshot().Hash.String()
	return
}
func
(duo *DuoUI) GetDuOSdifficulty() {
	duo.rc.Difficulty = rpc.GetDifficultyRatio(duo.cx.RPCServer.Cfg.Chain.BestSnapshot().Bits, duo.cx.RPCServer.Cfg.ChainParams, 2)
	return
}
func
(duo *DuoUI) GetDuOSblockCount() {
	getBlockCount, err := rpc.HandleGetBlockCount(duo.cx.RPCServer, nil, nil)
	if err != nil {
		duo.rc.PushDuOSalert("Error", err.Error(), "error")
	}
	duo.rc.BlockCount = getBlockCount.(int64)
	return
}
func
(duo *DuoUI) GetDuOSnetworkLastBlock() {
	for _, g := range duo.cx.RPCServer.Cfg.ConnMgr.ConnectedPeers() {
		l := g.ToPeer().StatsSnapshot().LastBlock
		if l > duo.rc.NetLastBlock {
			duo.rc.NetLastBlock = l
		}
	}
	return
}
func
(duo *DuoUI) GetDuOSconnectionCount() {
	duo.rc.Connections = duo.cx.RPCServer.Cfg.ConnMgr.ConnectedCount()
	return
}
func
(duo *DuoUI) GetDuOSlocalLost() {
	duo.rc.Localhost = *new(DuOSlocalHost)
	sm, _ := mem.VirtualMemory()
	sc, _ := cpu.Info()
	sp, _ := cpu.Percent(0, true)
	sd, _ := disk.Usage("/")
	duo.rc.Localhost.Cpu = sc
	duo.rc.Localhost.CpuPercent = sp
	duo.rc.Localhost.Memory = *sm
	duo.rc.Localhost.Disk = *sd
	return
}
