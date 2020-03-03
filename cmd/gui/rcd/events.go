package rcd

import (
	blockchain "github.com/p9c/pod/pkg/chain"
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

func (r *RcVar) ListenInit(trigger chan struct{}) {
	r.Events = EventsChan
	r.UpdateTrigger = trigger
	// first time starting up get all of these and trigger update
	update(r)
	r.cx.RealNode.Chain.Subscribe(func(callback *blockchain.Notification) {
		switch callback.Type {
		case blockchain.NTBlockAccepted,
			blockchain.NTBlockConnected,
			blockchain.NTBlockDisconnected:
			go update(r)
			// go r.toastAdd("New block: "+fmt.Sprint(callback.Data.(*util.Block).Height()), callback.Data.(*util.Block).Hash().String())
		}
	})
	go func() {
	out:
		for {
			select {
			case <-r.cx.WalletServer.Update:
				update(r)
			case <-r.cx.KillAll:
				break out
			}
		}
	}()
	log.WARN("event update listener started")
	return
}

func update(r *RcVar) {
	// log.WARN("GetDuoUIbalance")
	r.GetDuoUIbalance()
	// log.WARN("GetDuoUIunconfirmedBalance")
	r.GetDuoUIunconfirmedBalance()
	// log.WARN("GetDuoUItransactionsNumber")
	r.GetDuoUItransactionsNumber()
	// r.GetTransactions()
	// log.WARN("GetLatestTransactions")
	r.GetLatestTransactions()
	log.INFO("")
	log.INFO("UPDATE")
	log.INFO(r.Status.Wallet.LastTxs.Txs)
	log.INFO("")
	// r.GetDuoUIstatus()
	// r.GetDuoUIlocalLost()
	// r.GetDuoUIblockHeight()
	// log.WARN("GetDuoUIblockCount")
	r.GetDuoUIblockCount()
	// log.WARN("GetDuoUIdifficulty")
	r.GetDuoUIdifficulty()
	// log.WARN("GetDuoUIconnectionCount")
	r.GetDuoUIconnectionCount()
	r.GetAddressBook()
	r.UpdateTrigger <- struct{}{}
}
