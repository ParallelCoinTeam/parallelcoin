package mod

import (
	"fmt"
)

const (
	NewBlock uint32 = iota
	// Add new events here
	EventCount
)

// type Event struct {
//	Type    uint32
//	Payload []byte
// }

// var EventsChan = make(chan Event, 1)

// func (r *RcVar) ListenInit(trigger chan struct{}) {
//	//L.Debug("listeninit")
//	//r.Events = EventsChan
//	//r.UpdateTrigger = trigger
//
//	// first time starting up get all of these and trigger update
//	update(r)
//	//r.labelMiningAddreses()
//
//	var ready atomic.Bool
//	ready.Store(false)
//	r.Cx.RealNode.Chain.Subscribe(func(callback *blockchain.Notification) {
//		switch callback.Type {
//		case blockchain.NTBlockAccepted,
//			blockchain.NTBlockConnected,
//			blockchain.NTBlockDisconnected:
//			if !ready.Load() {
//				return
//			}
//			update(r)
//			// go r.toastAdd("New block: "+fmt.Sprint(callback.Data.(*util.Block).Height()), callback.Data.(*util.Block).Hash().String())
//		}
//	})
//	go func() {
//		ticker := time.NewTicker(time.Second)
//	out:
//		for {
//			select {
//			case <-ticker.C:
//				if !ready.Load() {
//					if r.Cx.IsCurrent() {
//						ready.Store(true)
//						// 		go func() {
//						// 			r.cx.WalletServer.Rescan(nil, nil)
//						// 			r.Ready <- struct{}{}
//						// 			r.UpdateTrigger <- struct{}{}
//						// 		}()
//					}
//				}
//				//r.GetDuoUIconnectionCount()
//				r.UpdateTrigger <- struct{}{}
//			// L.Warn("GetDuoUIconnectionCount")
//			case <-r.Cx.WalletServer.Update:
//				update(r)
//			case <-r.Cx.KillAll:
//				break out
//			}
//		}
//	}()
//	//L.Warn("event update listener started")
//	return
// }

func update(r *RcVar) {
	fmt.Println("kakakak")
	// L.Warn("GetDuoUIbalance")
	// r.GetDuoUIbalance()
	// L.Warn("GetDuoUIunconfirmedBalance")
	// r.GetDuoUIunconfirmedBalance()
	// L.Warn("GetDuoUItransactionsNumber")
	// r.GetDuoUItransactionsNumber()
	// r.GetTransactions()
	// L.Warn("GetLatestTransactions")
	// r.GetLatestTransactions()
	// L.Info("")
	// L.Info("UPDATE")
	// L.Trace(r.History.PerPage)
	// L.Info("")
	// r.GetDuoUIstatus()
	// r.GetDuoUIlocalLost()
	// r.GetDuoUIblockHeight()
	// L.Warn("GetDuoUIblockCount")
	// r.GetDuoUIdifficulty()
	// r.GetDuoUIblockCount()
	// r.GetPeerInfo()
	// L.Warn("GetDuoUIdifficulty")
	// r.UpdateTrigger <- struct{}{}
}
