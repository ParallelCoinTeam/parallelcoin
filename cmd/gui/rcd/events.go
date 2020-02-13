package rcd

import (
	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
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
var UpdateTrigger = make(chan struct{}, 1)

func ListenInit(cx *conte.Xt, rc *RcVar) {
	rc.Events = EventsChan
	rc.UpdateTrigger = UpdateTrigger
	// first time starting up get all of these and trigger update
	rc.GetDuoUIbalance(cx)
	rc.GetDuoUIunconfirmedBalance(cx)
	rc.GetDuoUITransactionsExcertps(cx)
	rc.GetDuoUIblockHeight(cx)
	rc.GetDuoUIstatus(cx)
	rc.GetDuoUIlocalLost()
	rc.GetDuoUIdifficulty(cx)
	rc.GetDuoUIlastTxs(cx)
	log.DEBUG("sending trigger to populate data for first start")
	rc.UpdateTrigger <- struct{}{}
	log.DEBUG("sent trigger to populate data for first start")
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
			log.DEBUG("sending trigger to populate data for new block")
			rc.UpdateTrigger <- struct{}{}
			log.DEBUG("sent trigger to populate data for new block")
		}
	})
	
	return
}
