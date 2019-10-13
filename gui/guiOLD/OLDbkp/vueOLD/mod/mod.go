package mod

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"

	"github.com/p9c/pod/pkg/rpc/btcjson"
)

type DuoGuiTemplates struct {
	App  map[string][]byte            `json:"app"`
	Data map[string]map[string][]byte `json:"data"`
}

type Blocks struct {
	// Per         int                          `json:"per"`
	// Page        int                          `json:"page"`
	CurrentPage int                             `json:"currentpage"`
	PageCount   int                             `json:"pagecount"`
	Blocks      []btcjson.GetBlockVerboseResult `json:"blocks"`
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

type DuoGuiItem struct {
	Enabled  bool        `json:"enabled"`
	Name     string      `json:"name"`
	Slug     string      `json:"slug"`
	Version  string      `json:"ver"`
	CompType string      `json:"comptype"`
	SubType  string      `json:"subtype"`
	Data     interface{} `json:"data"`
}
type DuoGuiItems struct {
	Slug  string                `json:"slug"`
	Items map[string]DuoGuiItem `json:"items"`
}

type DuoVUEcomps []DuoVUEcomp

//  Vue App Model
type DuoVUEcomp struct {
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

type DuoVUEbalance struct {
	Balance     string `json:"balance"`
	Unconfirmed string `json:"unconfirmed"`
}

type DuoVUEtransactions struct {
	Txs       []btcjson.ListTransactionsResult `json:"txs"`
	TxsNumber int                              `json:"txsnumber"`
}
type DuoVUEtransactionsExcerpts struct {
	Txs           []TransactionExcerpt `json:"txs"`
	TxsNumber     int                  `json:"txsnumber"`
	Balance       float64              `json:"balance"`
	BalanceHeight float64              `json:"balanceheight"`
}

type TransactionExcerpt struct {
	Balance       float64 `json:"balance"`
	Amount        float64 `json:"amount"`
	Category      string  `json:"category"`
	Confirmations int64   `json:"confirmations"`
	Time          string  `json:"time"`
	TxID          string  `json:"txid"`
	Comment       string  `json:"comment,omitempty"`
}
type DuoVUEaddressBook struct {
	Num       int       `json:"num"`
	Addresses []Address `json:"addresses"`
}

type DuoVUEblocks []DuoVUEblock

type DuoVUEblock struct {
	Height     int64   `json:"height"`
	PowAlgoID  uint32  `json:"pow"`
	Difficulty float64 `json:"diff"`
	Amount     float64 `json:"amount"`
	TxNum      int     `json:"txnum"`
	Time       int64   `json:"time"`
}

type DuoVUEchain struct {
	LastPushed int64         `json:"lastpushed"`
	Blocks     []DuoVUEblock `json:"blocks"`
}

// System Ststus
type DuoVUEstatus struct {
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
	Balance          DuoVUEbalance                    `json:"balance"`
	BlockCount       int64                            `json:"blockcount"`
	ConnectionCount  int32                            `json:"connectioncount"`
	NetworkLastBlock int32                            `json:"networklastblock"`
	TxsNumber        int                              `json:"txsnumber"`
	//MempoolInfo      string                        `json:"ver"`
}
