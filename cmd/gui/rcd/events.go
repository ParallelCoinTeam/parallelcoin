package rcd

import (
	"time"

	"go.uber.org/atomic"

	blockchain "github.com/p9c/pod/pkg/chain"
	log "github.com/p9c/logi"
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
	r.labelMiningAddreses()

	var ready atomic.Bool
	ready.Store(false)
	r.cx.RealNode.Chain.Subscribe(func(callback *blockchain.Notification) {
		switch callback.Type {
		case blockchain.NTBlockAccepted,
			blockchain.NTBlockConnected,
			blockchain.NTBlockDisconnected:
			if !ready.Load() {
				return
			}
			update(r)
			// go r.toastAdd("New block: "+fmt.Sprint(callback.Data.(*util.Block).Height()), callback.Data.(*util.Block).Hash().String())
		}
	})
	go func() {
		ticker := time.NewTicker(time.Second)
	out:
		for {
			select {
			case <-ticker.C:
				if !ready.Load() {
					if r.cx.IsCurrent() {
						ready.Store(true)
						r.cx.WalletServer.Rescan(nil, nil)
						r.UpdateTrigger <- struct{}{}
					}
				}
				r.GetDuoUIconnectionCount()
				r.UpdateTrigger <- struct{}{}
			// log.L.Warn("GetDuoUIconnectionCount")
			case <-r.cx.WalletServer.Update:
				update(r)
			case <-r.cx.KillAll:
				break out
			}
		}
	}()
	log.L.Warn("event update listener started")
	return
}

func update(r *RcVar) {
	// log.L.Warn("GetDuoUIbalance")
	r.GetDuoUIbalance()
	// log.L.Warn("GetDuoUIunconfirmedBalance")
	r.GetDuoUIunconfirmedBalance()
	// log.L.Warn("GetDuoUItransactionsNumber")
	r.GetDuoUItransactionsNumber()
	// r.GetTransactions()
	// log.L.Warn("GetLatestTransactions")
	r.GetLatestTransactions()
	// log.L.Info("")
	// log.L.Info("UPDATE")
	log.L.Info(r.Status.Wallet.LastTxs.Txs)
	// log.L.Info("")
	// r.GetDuoUIstatus()
	// r.GetDuoUIlocalLost()
	// r.GetDuoUIblockHeight()
	// log.L.Warn("GetDuoUIblockCount")
	r.GetDuoUIdifficulty()
	r.GetDuoUIblockCount()
	// log.L.Warn("GetDuoUIdifficulty")
	r.UpdateTrigger <- struct{}{}
}
