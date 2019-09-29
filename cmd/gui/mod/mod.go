package mod

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"

	"github.com/p9c/pod/pkg/rpc/btcjson"
)

type DuOStemplates struct {
	App  map[string][]byte            `json:"app"`
	Data map[string]map[string][]byte `json:"data"`
}

type DbAddress string

type Address struct {
	Index   int     `json:"num"`
	Label   string  `json:"label"`
	Account string  `json:"account"`
	Address string  `json:"address"`
	Amount  float64 `json:"amount"`
}
type Send struct {
	// Phrase string  `json:"phrase"`
	// Addr   string  `json:"addr"`
	// Amount float64 `json:"amount"`
	//exit
}

type AddBook struct {
	Address string `json:"address"`
	Label   string `json:"label"`
}

// Items

type DuOSitem struct {
	Enabled  bool        `json:"enabled"`
	Name     string      `json:"name"`
	Slug     string      `json:"slug"`
	Version  string      `json:"ver"`
	CompType string      `json:"comptype"`
	SubType  string      `json:"subtype"`
	Data     interface{} `json:"data"`
}
type DuOSitems struct {
	Slug  string              `json:"slug"`
	Items map[string]DuOSitem `json:"items"`
}

type DuOScomps []DuOScomp

//  Vue App Model
type DuOScomp struct {
	IsApp       bool   `json:"isapp"`
	Name        string `json:"name"`
	ID          string `json:"id"`
	Version     string `json:"ver"`
	Description string `json:"desc"`
	State       string `json:"state"`
	Image       string `json:"img"`
	URL         string `json:"url"`
	CompType    string `json:"comtype"`
	SubType     string `json:"subtype"`
	Js          string `json:"js"`
	Template    string `json:"template"`
	Css         string `json:"css"`
}

// System Ststus
type DuOSstatus struct {
	Version          string                           `json:"ver"`
	WalletVersion    map[string]btcjson.VersionResult `json:"walletver"`
	UpTime           int64                            `json:"uptime"`
	Cpu              []cpu.InfoStat                   `json:"cpu"`
	CpuPercent       []float64                        `json:"cpupercent"`
	Memory           mem.VirtualMemoryStat            `json:"mem"`
	Disk             disk.UsageStat                   `json:"disk"`
	CurrentNet       string                           `json:"net"`
	Chain            string                           `json:"chain"`
	HashesPerSec     int64                            `json:"hashrate"`
	Height           int32                            `json:"height"`
	BestBlockHash    string                           `json:"bestblockhash"`
	NetworkHashPS    int64                            `json:"networkhashrate"`
	Difficulty       float64                          `json:"diff"`
	Balance          DuOSbalance                      `json:"balance"`
	BlockCount       int64                            `json:"blockcount"`
	ConnectionCount  int32                            `json:"connectioncount"`
	NetworkLastBlock int32                            `json:"networklastblock"`
	TxsNumber        int                              `json:"txsnumber"`
	//MempoolInfo      string                        `json:"ver"`
}

type DuOSdata struct {
	Status               DuOSstatus                   `json:"status"`
	Addressbook          DuOSaddressBook              `json:"addressbook"`
	Peers                []*btcjson.GetPeerInfoResult `json:"peers"`
	TransactionsExcerpts DuOStransactionsExcerpts     `json:"txsex"`
	Blocks               DuOSblocks                   `json:"blocks"`
	Send                 Send                         `json:"send"`
	Icons                map[string]string            `json:"ico"`
}
