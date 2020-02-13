package rcd

import (
	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/conte"
)

const (
	NewBlock uint32 = iota
	// Add new events here
	EventCount
)

type Event struct {
	Type    uint32
	Payload []byte
}

var EventsChan = make(chan Event, 1)

func ListenInit(cx *conte.Xt, rc *RcVar, trigger chan struct{}) {
	rc.Events = EventsChan
	rc.UpdateTrigger = trigger
	// first time starting up get all of these and trigger update
	rc.GetDuoUIbalance()
	rc.GetDuoUIunconfirmedBalance()
	rc.ComTransactions()
	rc.GetDuoUIblockHeight()
	rc.GetDuoUIstatus()
	rc.GetDuoUIlocalLost()
	rc.GetDuoUIdifficulty()
	rc.ComLatestTransactions()
	cx.RealNode.Chain.Subscribe(func(callback *blockchain.Notification) {
		switch callback.Type {
		case blockchain.NTBlockAccepted:
			rc.GetDuoUIbalance()
			rc.GetDuoUIunconfirmedBalance()
			rc.ComTransactions()
			rc.GetDuoUIblockHeight()
			rc.GetDuoUIstatus()
			rc.GetDuoUIlocalLost()
			rc.GetDuoUIdifficulty()
			rc.ComLatestTransactions()
			rc.UpdateTrigger <- struct{}{}
		}
		rc.UpdateTrigger <- struct{}{}
	})
	
	return
}
