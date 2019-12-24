package models

import (
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

// System Ststus
type
	DuoUIstatus struct {
		Version       string                           `json:"ver"`
		WalletVersion map[string]btcjson.VersionResult `json:"walletver"`
		UpTime        int64                            `json:"uptime"`
		CurrentNet    string                           `json:"net"`
		Chain         string                           `json:"chain"`
	}
type
	DuoUIhashes struct{ int64 }
type
	DuoUInetworkHash struct{ int64 }
type
	DuoUIheight struct{ int32 }
type
	DuoUIbestBlockHash struct{ string }
type
	DuoUIdifficulty struct{ float64 }

//type
// MempoolInfo      struct { string}
type
	DuoUIblockCount struct{ int64 }
type
	DuoUInetLastBlock struct{ int32 }
type
	DuoUIconnections struct{ int32 }
type
	DuoUIlocalHost struct {
		Cpu        []cpu.InfoStat        `json:"cpu"`
		CpuPercent []float64             `json:"cpupercent"`
		Memory     mem.VirtualMemoryStat `json:"mem"`
		Disk       disk.UsageStat        `json:"disk"`
	}
