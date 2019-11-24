package gui

import (
	"github.com/p9c/pod/pkg/conte"
)

var rcv = &rcvar{
	alert:       DuOSalert{},
	status:      DuOStatus{},
	balance:     DuOSbalance{},
	lastxs:      DuOStransactions{},
	blockcount:  0,
	connections: 0,
}

type rcvar struct {
	cx           *conte.Xt
	alert        DuOSalert
	status       DuOStatus
	nethash      int64
	hashes       int64
	balance      DuOSbalance
	transactions DuOStransactions
	txs          DuOStransactionsExcerpts
	lastxs       DuOStransactions
	blockcount   int64
	netlastblock int32
	connections  int32
	sent         bool
	IsFirstRun   bool
	localhost    LocalDuOShost
}

type RcVar interface {
	GetDuOStransactions(sfrom, count int, cat string)
	GetDuOSbalance()
	GetDuOStransactionsExcerpts()
	DuoSend(wp string, ad string, am float64)
	GetDuOStatus()
	PushDuOSalert(t string, m interface{}, at string)
	GetDuOSblockCount()
	GetDuOSconnectionCount()
}
