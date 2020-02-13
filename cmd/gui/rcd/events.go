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
	rc.GetDuoUIbalance(cx)
	rc.GetDuoUIunconfirmedBalance(cx)
	rc.GetDuoUITransactionsExcertps(cx)
	rc.GetDuoUIblockHeight(cx)
	rc.GetDuoUIstatus(cx)
	rc.GetDuoUIlocalLost()
	rc.GetDuoUIdifficulty(cx)
	rc.GetDuoUIlastTxs(cx)
	cx.RealNode.Chain.Subscribe(func(callback *blockchain.Notification) {
		switch callback.Type {
		case blockchain.NTBlockAccepted:
			rc.GetDuoUIbalance(cx)
			rc.GetDuoUIunconfirmedBalance(cx)
			rc.GetDuoUITransactionsExcertps(cx)
			rc.GetDuoUIblockHeight(cx)
			rc.GetDuoUIstatus(cx)
			rc.GetDuoUIlocalLost()
			rc.GetDuoUIdifficulty(cx)
			rc.GetDuoUIlastTxs(cx)
			rc.UpdateTrigger <- struct{}{}
		}
		rc.UpdateTrigger <- struct{}{}
	})
	
	return
}
