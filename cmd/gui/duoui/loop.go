package duoui

import (
	"errors"
	"github.com/p9c/pod/cmd/gui/mvc/view"
	"github.com/p9c/pod/pkg/gui/io/system"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func DuoUImainLoop(sys *view.DuOS) error {
	ui := &DuoUI{
		ly:sys.Duo,
		rc:sys.Rc,
	}
	sys.Duo.Pages = ui.LoadPages()
	for {
		select {
		case <-ui.ly.Ready:
			ui.ly.IsReady = true
		case <-ui.ly.Quit:
			log.DEBUG("quit signal received")
			interrupt.Request()
			// This case is for handling when some external application is controlling the GUI and to gracefully
			// handle the back-end servers being shut down by the interrupt library receiving an interrupt signal
			// Probably nothing needs to be run between starting it and shutting down
			<-interrupt.HandlersDone
			log.DEBUG("closing GUI from interrupt/quit signal")
			return errors.New("shutdown triggered from back end")
		case e := <-ui.ly.Window.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				log.DEBUG("destroy event received")
				interrupt.Request()
				// Here do cleanup like are you sure (optional) modal or shutting down indefinite spinner
				<-interrupt.HandlersDone
				return e.Err
			case system.FrameEvent:
				ui.ly.Context.Reset(e.Config, e.Size)
				//go func() {
				//sys.Rc.GetDuoUIbalance()
				//sys.Rc.GetDuoUIunconfirmedBalance()
				//sys.Rc.ComTransactions()
				//
				//sys.Rc.GetDuoUIblockHeight()
				//sys.Rc.GetDuoUIstatus()
				//rc.GetDuoUIlocalLost()
				//rc.GetDuoUIdifficulty()

				//rc.GetDuoUIlastTxs()
				//time.Sleep(1 * time.Second)
				//}()

				//if rc.Boot.IsBoot {
				//d.DuoUImainScreen()
				//e.Frame(d.mod.Context.Ops)
				//} else {
				//	d.mod.Context.Reset(e.Config, e.Size)
				//	if rc.Boot.IsFirstRun {
				//		//DuoUIloaderCreateWallet(duo.m, cx, rc)
				//	} else {
				ui.DuoUImainScreen()
				//		if rc.Dialog.Show {
				//			d.DuoUIdialog(rc)
				//		}
				//		d.DuoUItoastSys()
				//
				//		go func() {
				//			time.Sleep(1 * time.Second)
				//
				//			//rc.GetDuoUIbalance()
				//			//rc.GetDuoUIunconfirmedBalance()
				//rc.GetDuoUITransactionsExcertps()
				//
				//			//rc.GetDuoUIblockHeight()
				//			//rc.GetDuoUIstatus()
				//			//rc.GetDuoUIlocalLost()
				//			//rc.GetDuoUIdifficulty()
				//
				//			//rc.GetDuoUIlastTxs()
				//		}()
				//	}
				e.Frame(ui.ly.Context.Ops)
				sys.Duo.Context.Reset(e.Config, e.Size)
				//}
			}
		}
	}
}