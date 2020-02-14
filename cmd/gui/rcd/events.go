package rcd

import (
	blockchain "github.com/p9c/pod/pkg/chain"
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

func (r *RcVar)ListenInit(trigger chan struct{}){
	r.Events = EventsChan
	r.UpdateTrigger = trigger
	// first time starting up get all of these and trigger update
	r.GetDuoUIbalance()
	r.GetDuoUIunconfirmedBalance()
	r.GetDuoUItransactionsNumber()
	r.GetTransactions()
	r.GetLatestTransactions()
	//r.GetDuoUIstatus()
	//r.GetDuoUIlocalLost()
	r.GetDuoUIblockHeight()
	r.GetDuoUIdifficulty()
	r.GetDuoUIconnectionCount()
	r.cx.RealNode.Chain.Subscribe(func(callback *blockchain.Notification) {
		switch callback.Type {
		case blockchain.NTBlockAccepted:
			r.GetDuoUIbalance()
			r.GetDuoUIunconfirmedBalance()
			r.GetDuoUItransactionsNumber()
			r.GetTransactions()
			r.GetLatestTransactions()
			//r.GetDuoUIstatus()
			//r.GetDuoUIlocalLost()
			r.GetDuoUIblockHeight()
			r.GetDuoUIdifficulty()
			r.GetDuoUIconnectionCount()
			r.UpdateTrigger <- struct{}{}
		}
		r.UpdateTrigger <- struct{}{}
	})
	
	return
}
