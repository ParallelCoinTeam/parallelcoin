//+build !headless

package gui

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/p9c/pod/cmd/gui/ast/ico"
	"github.com/p9c/pod/cmd/gui/hlp"
	"github.com/p9c/pod/cmd/gui/mod"
	"github.com/p9c/pod/cmd/gui/lyt"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"image/color"
)

func WalletGUI(rc *rcd.RcVar, cx *conte.Xt) (err error) {
	go func() {
		cx.DuoUI.Window = app.NewWindow()
		if err := loop(rc, cx); err != nil {
			log.FATAL(err)
		}
		//runUI()
	}()
	app.Main()
	return
}

func loop(rc *rcd.RcVar, cx *conte.Xt) error {
	gofont.Register()
	cx.DuoUI.Theme = *material.NewTheme()
	cx.DuoUI.Gtx = *layout.NewContext(cx.DuoUI.Window.Queue())
	for {
		e := <-cx.DuoUI.Window.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			cx.DuoUI.Gtx.Reset(e.Config, e.Size)
			rc.GetDuOSbalance(cx)
			rc.GetDuOSblockHeight(cx)
			rc.GetDuOStatus(cx)
			rc.GetDuOSlocalLost(cx)
			rc.GetDuOSdifficulty(cx)
			u := &mod.DuoUI{}
			i := &ico.DuoUIicons{}
			l := &lyt.DuoUIlayouts{}

			i.DuoUIicons()
			l.DuoUIlayouts()

			u.Ico = *i
			u.Layouts = *l

			u.Buttons.Logo = new(widget.Button)

			cs := cx.DuoUI.Gtx.Constraints




			header := cx.DuoUI.Layouts.View.Rigid(&cx.DuoUI.Gtx, func() {
				hlp.DuoUIdrawRect(&cx.DuoUI.Gtx, cs.Width.Max, 64, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf})
				hlp.DuoUIdrawRect(&cx.DuoUI.Gtx, 64, 64, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30})
				cx.DuoUI.Theme.IconButton(u.Ico.Logo).Layout(&cx.DuoUI.Gtx, u.Buttons.Logo)
			})

			//header := cx.DuoUI.Layouts.View.Rigid(&cx.DuoUI.Gtx, func(){ elem.DuoUIheader(cx.DuoUI)})

			main := cx.DuoUI.Layouts.View.Rigid(&cx.DuoUI.Gtx, func() {
				//balance := flh.Rigid(gtx, func() {
				//	in.Layout(gtx, func() {
				//		th.H3("balance :" + r.balance).Layout(gtx)
				//	})
				//})

				sidebar := cx.DuoUI.Layouts.Main.Rigid(&cx.DuoUI.Gtx, func() {

					hlp.DuoUIdrawRect(&cx.DuoUI.Gtx, 64, cs.Height.Max, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf})

					overview := cx.DuoUI.Layouts.Menu.Rigid(&cx.DuoUI.Gtx, func() {
						cx.DuoUI.Theme.IconButton(u.Ico.Overview).Layout(&cx.DuoUI.Gtx, u.Buttons.Logo)
					})
					history := cx.DuoUI.Layouts.Menu.Rigid(&cx.DuoUI.Gtx, func() {
						cx.DuoUI.Theme.IconButton(u.Ico.History).Layout(&cx.DuoUI.Gtx, u.Buttons.Logo)
					})
					network := cx.DuoUI.Layouts.Menu.Rigid(&cx.DuoUI.Gtx, func() {
						cx.DuoUI.Theme.IconButton(u.Ico.Network).Layout(&cx.DuoUI.Gtx, u.Buttons.Logo)
					})
					settings := cx.DuoUI.Layouts.Menu.Rigid(&cx.DuoUI.Gtx, func() {
						cx.DuoUI.Theme.IconButton(u.Ico.Settings).Layout(&cx.DuoUI.Gtx, u.Buttons.Logo)
					})
					cx.DuoUI.Layouts.Menu.Layout(&cx.DuoUI.Gtx, overview, history, network, settings)

				})

				content := cx.DuoUI.Layouts.Main.Rigid(&cx.DuoUI.Gtx, func() {
					in := layout.UniformInset(unit.Dp(11))
					in.Layout(&cx.DuoUI.Gtx, func() {
						hlp.DuoUIdrawRect(&cx.DuoUI.Gtx, cs.Width.Max, cs.Height.Max, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf})

						balance := cx.DuoUI.Layouts.Status.Rigid(&cx.DuoUI.Gtx, func() {
							cx.DuoUI.Theme.H5("balance :" + rc.Balance).Layout(&cx.DuoUI.Gtx)
						})
						blockheight := cx.DuoUI.Layouts.Status.Rigid(&cx.DuoUI.Gtx, func() {
							cx.DuoUI.Theme.H5("blockheight :" + fmt.Sprint(rc.BlockHeight)).Layout(&cx.DuoUI.Gtx)
						})
						difficulty := cx.DuoUI.Layouts.Status.Rigid(&cx.DuoUI.Gtx, func() {
							cx.DuoUI.Theme.H5("difficulty :" + fmt.Sprintf("%f", rc.Difficulty)).Layout(&cx.DuoUI.Gtx)
						})
						connections := cx.DuoUI.Layouts.Status.Rigid(&cx.DuoUI.Gtx, func() {
							cx.DuoUI.Theme.H5("connections :" + fmt.Sprint(rc.Connections)).Layout(&cx.DuoUI.Gtx)
						})

						cx.DuoUI.Layouts.Status.Layout(&cx.DuoUI.Gtx, balance, blockheight, difficulty, connections)
					})
				})

				cx.DuoUI.Layouts.Main.Layout(&cx.DuoUI.Gtx, sidebar, content)
			})
			//block := fl.Rigid(gtx, func() {
			//	th.H3("Block height :" + fmt.Sprint(r.height)).Layout(gtx)
			//})
			//
			//difficulty := fl.Rigid(gtx, func() {
			//	th.H3("difficulty :" + fmt.Sprint(r.difficulty)).Layout(gtx)
			//})
			////block := fl.Rigid(gtx, func() {
			////	th.H3("Block height :" + fmt.Sprint(r.height)).Layout(gtx)
			////})
			//status := fl.Rigid(gtx, func() {
			//	th.H3("Block height :" + fmt.Sprint(r.status)).Layout(gtx)
			//})
			//

			//cx.DuoUI.Layouts.View.Layout(&cx.DuoUI.Gtx , elem.DuoUIheader(cx.DuoUI), main)
			cx.DuoUI.Layouts.View.Layout(&cx.DuoUI.Gtx, header, main)
			e.Frame(cx.DuoUI.Gtx.Ops)
		}
	}
}
