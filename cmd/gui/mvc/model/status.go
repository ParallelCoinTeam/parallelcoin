package model

import (
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

// System Ststus
type
	DuoUIstatus struct {
		Version    string `json:"ver"`
		UpTime     int64  `json:"uptime"`
		CurrentNet string `json:"net"`
		Chain      string `json:"chain"`
		Node       *NodeStatus
		Wallet     *WalletStatus
		Kopach     *KopachStatus
	}

type NodeStatus struct {
	NetHash          int64
	BlockHeight      int32
	BestBlock        string
	Difficulty       float64
	BlockCount       int64
	NetworkLastBlock int32
	ConnectionCount  int32
}
type KopachStatus struct {
}
type WalletStatus struct {
	WalletVersion map[string]btcjson.VersionResult `json:"walletver"`
	Hashes        int64
	Balance       string
	Unconfirmed   string
	TxsNumber     int
	Transactions  *DuoUItransactions
	Txs           *DuoUItransactionsExcerpts
	LastTxs       *DuoUItransactions
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
		//Cpu        []cpu.InfoStat        `json:"cpu"`
		//CpuPercent []float64             `json:"cpupercent"`
		//Memory     mem.VirtualMemoryStat `json:"mem"`
		//Disk       disk.UsageStat        `json:"disk"`
	}
