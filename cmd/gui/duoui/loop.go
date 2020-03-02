package duoui

import (
	"errors"

	"gioui.org/io/system"
	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func DuoUImainLoop(d *model.DuoUI, r *rcd.RcVar) error {
	ui := new(DuoUI)
	ui = &DuoUI{
		ly: d,
		rc: r,
	}
	for {
		select {
		case <-ui.rc.Ready:
			updateTrigger := make(chan struct{}, 1)
			go func() {
			quitTrigger:
				for {
					select {
					case <-updateTrigger:
						log.TRACE("repaint forced")
						//ui.ly.Window.Invalidate()
					case <-ui.rc.Quit:
						break quitTrigger
					}
				}
			}()
			ui.rc.ListenInit(updateTrigger)
			ui.rc.IsReady = true
		case <-ui.rc.Quit:
			log.DEBUG("quit signal received")
			if !interrupt.Requested() {
				interrupt.Request()
			}
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
				if ui.rc.Boot.IsBoot {
					ui.DuoUIsplashScreen()
					e.Frame(ui.ly.Context.Ops)
				} else {
					if ui.rc.Boot.IsFirstRun {
						ui.DuoUIloaderCreateWallet()
					} else {
						ui.DuoUImainScreen()
						if ui.rc.Dialog.Show {
							ui.DuoUIdialog()
						}
						// ui.DuoUItoastSys()
					}
					e.Frame(ui.ly.Context.Ops)
				}
			}
		}
	}
}
