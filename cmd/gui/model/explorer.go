package model

import (
	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

type DuoUIexplorer struct {
	Page        *gel.DuoUIcounter
	PerPage     *gel.DuoUIcounter
	Blocks      []DuoUIblock
	SingleBlock btcjson.GetBlockVerboseResult
}

type DuoUIhistory struct {
	Page     *gel.DuoUIcounter
	PerPage  *gel.DuoUIcounter
	Txs      *DuoUItransactionsExcerpts
	SingleTx btcjson.GetTransactionDetailsResult
}
