package gui

import (
	"github.com/p9c/pod/pkg/conte"
)

type rcvar struct {
	cx     *conte.Xt
	alert  *DuOSalert
	status *DuOStatus
	txs    *DuOStransactionsExcerpts
	lastxs *DuOStransactions
}


type RcVar interface{
	GetTransactions(sfrom, count int, cat string) (txs DuOStransactions)
	GetBalance() (b DuOSbalance)
	GetTransactionsExcertps() (txse DuOStransactionsExcerpts)
	DuoSend(wp string, ad string, am float64) string
	GetDuOStatus() (s DuOStatus)
	PushDuOSalert(t string, m interface{}, at string) (d *DuOSalert)
	GetBlockCount() int64
	GetConnectionCount() int32
}