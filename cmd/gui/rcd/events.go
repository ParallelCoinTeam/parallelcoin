package rcd

import (
	"fmt"
	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/util"
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

func (r *RcVar) ListenInit(trigger chan struct{}) {
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

	r.GetAddressBook()
	r.cx.RealNode.Chain.Subscribe(func(callback *blockchain.Notification) {
		switch callback.Type {
		case blockchain.NTBlockAccepted:
			go r.GetDuoUIbalance()
			go r.GetDuoUIunconfirmedBalance()
			go r.GetDuoUItransactionsNumber()
			go r.GetTransactions()
			go r.GetLatestTransactions()
			//r.GetDuoUIstatus()
			//r.GetDuoUIlocalLost()
			go r.GetDuoUIblockHeight()
			go r.GetDuoUIdifficulty()
			go r.GetDuoUIconnectionCount()
			r.UpdateTrigger <- struct{}{}
			go r.toastAdd("New block: "+fmt.Sprint(callback.Data.(*util.Block).Height()), callback.Data.(*util.Block).Hash().String())
		}
		r.UpdateTrigger <- struct{}{}
	})

	return
}
