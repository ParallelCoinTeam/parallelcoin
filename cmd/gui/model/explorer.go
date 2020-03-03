package model

import (
	"github.com/p9c/pod/pkg/gui/controller"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

type DuoUIexplorer struct {
	Page        *controller.DuoUIcounter
	PerPage     *controller.DuoUIcounter
	Blocks      []DuoUIblock
	SingleBlock btcjson.GetBlockVerboseResult
}

type DuoUIhistory struct {
	Page     *controller.DuoUIcounter
	PerPage  *controller.DuoUIcounter
	Txs      *DuoUItransactionsExcerpts
	SingleTx btcjson.GetTransactionDetailsResult
}
