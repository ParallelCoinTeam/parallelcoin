// SPDX-License-Identifier: Unlicense OR MIT

package dap

import (
	"errors"
	"fmt"
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"github.com/p9c/pod/pkg/gui/wallet/dap/box"
	"github.com/p9c/pod/pkg/gui/wallet/dap/res"
	"github.com/p9c/pod/pkg/gui/wallet/dap/win"
	"github.com/p9c/pod/pkg/gui/wallet/lyt"
	"github.com/p9c/pod/pkg/util/interrupt"
	"log"
	"os"
	"sync/atomic"
)

//func (d *dap) DAP() {
//	defer os.Exit(0)
//	if err := d.loop(); err != nil {
//		log.Fatal(err)
//	}
//}

func (d *dap) DAP() {
	_ = d.loop()
	//if d.boot.Rc.Boot.IsBoot {
	//	fmt.Println("BootIsBoot")
	//	if d.boot.Rc.Boot.IsFirstRun {
	//		ui.DuoUIloaderCreateWallet()
	//fmt.Println("DuoUIloaderCreateWallet")
	//} else {
	//	d.newWindow(d.SplashScreen)
	//	fmt.Println("DuoUIsplashScreen")
	//}
	//} else {
	//	ui.DuoUImainScreen()
	//}

	app.Main()
}

func (d *dap) loop() error {
	d.newWindow(d.SplashScreen)
	//Debug("starting up GUI2-=0-=0-=22")

	for {
		select {
		case <-d.boot.Rc.Ready:
			updateTrigger := make(chan struct{}, 1)
			go func() {
			quitTrigger:
				for {
					select {
					case <-updateTrigger:
						//log.L.Trace("repaint forced")
						d.boot.UI.W.W["main"].W.Invalidate()
					case <-d.boot.Rc.Quit:
						break quitTrigger
					}
				}
			}()
			d.boot.Rc.ListenInit(updateTrigger)
			d.boot.Rc.IsReady = true
		case <-d.boot.Rc.Quit:
			//log.L.Debug("quit signal received")
			if !interrupt.Requested() {
				interrupt.Request()
			}
			// This case is for handling when some external application is controlling the GUI and to gracefully
			// handle the back-end servers being shut down by the interrupt library receiving an interrupt signal
			// Probably nothing needs to be run between starting it and shutting down
			<-interrupt.HandlersDone
			//log.L.Debug("closing GUI from interrupt/quit signal")
			return errors.New("shutdown triggered from back end")
		//TODO events of gui
		//case e := <-d.boot.Rc.Commands.Events:
		//	switch e := e.(type) {
		//	case mod.CommandEvent:
		//		d.boot.Rc.Commands.History = append(d.boot.Rc.Commands.History, e.Command)
		//		d.boot.UI.Window.Invalidate()
		//}

		case e := <-d.boot.UI.W.W["main"].W.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				//log.L.Debug("destroy event received")
				interrupt.Request()
				// Here do cleanup like are you sure (optional) modal or shutting down indefinite spinner
				<-interrupt.HandlersDone
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&d.boot.UI.Ops, e)
				//d.boot.UI.G.Reset(e.Config, e.Size)
				gtx.Reset()
				//if d.boot.Rc.Boot.IsBoot {
				//d.newWindow(d.SplashScreen)
				fmt.Println("BootIsBoot")
				//if d.boot.Rc.Boot.IsFirstRun {
				//ui.DuoUIloaderCreateWallet()
				//fmt.Println("DuoUIloaderCreateWallet")
				//} else {
				//ui.DuoUIsplashScreen()
				//fmt.Println("DuoUIsplashScreen")
				//}
				//e.Frame(d.boot.UI.G.Ops)
				//} else {
				//ui.DuoUImainScreen()
				d.BeforeMain(gtx)
				d.Main(gtx)
				//lyt.Format(gtx, "max(inset(0dp0dp0dp0dp,_))", d.Main())
				d.AfterMain(gtx)
				//e.Frame(gtx.Ops)
				//if d.boot.Rc.Dialog.Show {
				//	component.DuoUIdialog(ui.rc, ui.ly.Context, ui.ly.Theme)
				//	ui.DuoUItoastSys()
				//}
				e.Frame(gtx.Ops)
				//}
			}
			//d.boot.UI.Window.Invalidate()
		}
	}
}

func (d *dap) BesssforeMain(gtx C) {

	////case e := <-d.boot.UI.W.W["main"].W.Events():
	////switch e := e.(
	////type
	////) {
	//case system.DestroyEvent:
	//	//log.L.Debug("destroy event received")
	//	//interrupt.Request()
	//	// Here do cleanup like are you sure (optional) modal or shutting down indefinite spinner
	//	//<-interrupt.HandlersDone
	//	return e.Err
	//case system.FrameEvent:
	//	gtx := layout.NewContext(&d.boot.UI.Ops, e)
	//	//d.boot.UI.G.Reset(e.Config, e.Size)
	//	gtx.Reset()
	//	if d.boot.Rc.Boot.IsBoot {
	//		fmt.Println("BootIsBoot")
	//		if d.boot.Rc.Boot.IsFirstRun {
	//			//ui.DuoUIloaderCreateWallet()
	//			fmt.Println("DuoUIloaderCreateWallet")
	//		} else {
	//			//ui.DuoUIsplashScreen()
	//			fmt.Println("DuoUIsplashScreen")
	//		}
	//		e.Frame(d.boot.UI.G.Ops)
	//	} else {
	//		//ui.DuoUImainScreen()
	//		d.BeforeMain(gtx)
	//		d.Main(gtx)
	//		//lyt.Format(gtx, "max(inset(0dp0dp0dp0dp,_))", d.Main())
	//		d.AfterMain(gtx)
	//		//e.Frame(gtx.Ops)
	//		//if d.boot.Rc.Dialog.Show {
	//		//	component.DuoUIdialog(ui.rc, ui.ly.Context, ui.ly.Theme)
	//		//	ui.DuoUItoastSys()
	//	}
	//	e.Frame(gtx.Ops)
	//	//}
	//}
	//d.boot.UI.Window.Invalidate()
}

func (d *dap) BeforeMain(gtx C) {
	d.boot.UI.R = res.Resposnsivity(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
	d.boot.UI.N.Mode = d.boot.UI.R.Mode
	d.boot.UI.N.NavLayout = d.boot.UI.R.Mod["Nav"].(string)
	d.boot.UI.N.ItemLayout = d.boot.UI.R.Mod["NavIconAndLabel"].(string)
	d.boot.UI.N.Axis = d.boot.UI.R.Mod["NavItemsAxis"].(layout.Axis)
	d.boot.UI.N.Size = d.boot.UI.R.Mod["NavSize"].(int)
	d.boot.UI.N.NoContent = d.boot.UI.N.Wide
	d.boot.UI.N.LogoWidget = d.boot.UI.N.LogoLayout(d.boot.UI.Theme)
}

func (d *dap) AfterMain(gtx C) {
	//pop.Popup(gtx, d.boot.UI.Theme, func(gtx C)D{return material.H3(d.boot.UI.Theme.T,"tetstette").Layout(gtx)})
	//return lyt.Format(gtx, "hflex(middle,f(1,inset(8dp8dp8dp8dp,_)))",
	//pop.Popup(d.boot.UI.Theme, func(gtx C) D {
	//	title := theme.Body(d.boot.UI.Theme, "Requested payments history")
	//	title.Alignment = text.Start
	//	return title.Layout(gtx)
	//	}),
	//})

}

func (d *dap) Main(gtx C) D {
	return lyt.Format(gtx, "max(inset(0dp0dp0dp0dp,_))", func(gtx C) D {
		return lyt.Format(gtx, d.boot.UI.R.Mod["Container"].(string),
			box.BoxBase(d.boot.UI.Theme.Colors["NavBg"], d.boot.UI.N.Nav(d.boot.UI.Theme, gtx)),
			func(gtx C) D {
				return lyt.Format(gtx, d.boot.UI.R.Mod["Main"].(string),
					d.boot.UI.N.CurrentPage.P(d.boot.UI.Theme, d.boot.UI.R.Mod["Page"].(string)),
					d.boot.UI.F,
				)
			})
	})
}

func (d *dap) newWindow(lyt func(gtx C) D) {
	atomic.AddInt32(&d.boot.UI.W.WindowCount, +1)
	go func() {
		w := &win.Window{
			L: lyt,
		}
		w.W = app.NewWindow()
		if err := w.Loop(w.W.Events()); err != nil {
			log.Fatal(err)
		}
		if c := atomic.AddInt32(&d.boot.UI.W.WindowCount, -1); c == 0 {
			os.Exit(0)
		}
	}()
}
