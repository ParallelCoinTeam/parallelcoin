package model

import (
	"gioui.org/layout"
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
	TransList  *layout.List
	Category   string
	Categories *DuoUIhistoryCategories
	Page       *gel.DuoUIcounter
	PerPage    *gel.DuoUIcounter
	Txs        *DuoUItransactionsExcerpts
	SingleTx   btcjson.GetTransactionDetailsResult
}
type DuoUIhistoryCategories struct {
	AllTxs      *gel.CheckBox
	MintedTxs   *gel.CheckBox
	ImmatureTxs *gel.CheckBox
	SentTxs     *gel.CheckBox
	ReceivedTxs *gel.CheckBox
}
