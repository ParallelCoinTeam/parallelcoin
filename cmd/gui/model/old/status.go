package model

import (
	"time"

	"go.uber.org/atomic"

	"github.com/p9c/pod/pkg/data/ring"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

// System Ststus
type DuoUIstatus struct {
	Version    string
	StartTime  time.Time
	CurrentNet string
	Chain      string
	Node       *NodeStatus
	Wallet     *WalletStatus
	Kopach     *KopachStatus
}

type NodeStatus struct {
	// updates on new block and init
	NetHash          atomic.Uint64
	BlockHeight      atomic.Uint64
	BestBlock        atomic.String
	Difficulty       atomic.Float64
	BlockCount       atomic.Uint64
	NetworkLastBlock atomic.Int32
	// update every 5 seconds
	ConnectionCount atomic.Int32
}
type KopachStatus struct {
	// update every second
	Hashrate uint64
	Hps      *ring.BufferFloat64
}
type WalletStatus struct {
	// unchanging
	WalletVersion map[string]btcjson.VersionResult `json:"walletver"`
	// update on new block and at start
	Balance     atomic.String
	Unconfirmed atomic.String
	TxsNumber   atomic.Uint64
	// components
	LastTxs *DuoUItransactionsExcerpts
}

// slots of user interface
type DuoUIhashes struct{ int64 }
type DuoUInetworkHash struct{ int64 }
type DuoUIheight struct{ int32 }
type DuoUIbestBlockHash struct{ string }
type DuoUIdifficulty struct{ float64 }

// type
// MempoolInfo      struct { string}
type DuoUIblockCount struct{ int64 }
type DuoUInetLastBlock struct{ int32 }
type DuoUIconnections struct{ int32 }
type DuoUIlocalHost struct {
	// Cpu        []cpu.InfoStat        `json:"cpu"`
	// CpuPercent []float64             `json:"cpupercent"`
	// Memory     mem.VirtualMemoryStat `json:"mem"`
	// Disk       disk.UsageStat        `json:"disk"`
}
