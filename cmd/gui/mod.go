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
	cx    *conte.Xt
	alert DuOSalert

	status       DuOStatus
	hashes       int64
	nethash      int64
	height       int32
	bestblock    string
	difficulty   float64
	blockcount   int64
	netlastblock int32
	connections  int32

	balance      DuOSbalance
	transactions DuOStransactions
	txs          DuOStransactionsExcerpts
	lastxs       DuOStransactions

	sent       bool
	IsFirstRun bool
	localhost  DuOSlocalHost
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
