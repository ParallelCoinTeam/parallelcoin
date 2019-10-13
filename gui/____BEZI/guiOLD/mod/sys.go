package mod

import (
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

type Blocks struct {
	// Per         int                          `json:"per"`
	// Page        int                          `json:"page"`
	CurrentPage int                             `json:"currentpage"`
	PageCount   int                             `json:"pagecount"`
	Blocks      []btcjson.GetBlockVerboseResult `json:"blocks"`
}

type DuOSbalance struct {
	Balance     string `json:"balance"`
	Unconfirmed string `json:"unconfirmed"`
}

type DuOStransactions struct {
	Txs       []btcjson.ListTransactionsResult `json:"txs"`
	TxsNumber int                              `json:"txsnumber"`
}
type DuOStransactionsExcerpts struct {
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
type DuOSaddressBook struct {
	Num       int       `json:"num"`
	Addresses []Address `json:"addresses"`
}

type DuOSblocks []DuOSblock

type DuOSblock struct {
	Height     int64   `json:"height"`
	PowAlgoID  uint32  `json:"pow"`
	Difficulty float64 `json:"diff"`
	Amount     float64 `json:"amount"`
	TxNum      int     `json:"txnum"`
	Time       int64   `json:"time"`
}

type DuOSchain struct {
	LastPushed int64       `json:"lastpushed"`
	Blocks     []DuOSblock `json:"blocks"`
}
