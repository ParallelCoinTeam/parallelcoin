package gui

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/p9c/pod/cmd/gui/assets/ico"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/mod"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"golang.org/x/exp/shiny/materialdesign/icons"
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
			rc.GetDuOSheight(cx)
			rc.GetDuOStatus(cx)
			rc.GetDuOSlocalLost(cx)
			rc.GetDuOSdifficulty(cx)
			u := &mod.DuoUI{}
			i := &mod.DuoUIicons{}
			var err error
			i.Logo, err = material.NewIcon(ico.ParalleCoin)
			if err != nil {
				log.FATAL(err)
			}

			i.Overview, err = material.NewIcon(icons.ActionHome)
			if err != nil {
				log.FATAL(err)
			}
			i.History, err = material.NewIcon(icons.ActionHistory)
			if err != nil {
				log.FATAL(err)
			}
			i.Network, err = material.NewIcon(icons.DeviceNetworkCell)
			if err != nil {
				log.FATAL(err)
			}
			i.Settings, err = material.NewIcon(icons.ActionSettings)
			if err != nil {
				log.FATAL(err)
			}

			u.Ico = *i
			u.Buttons.Logo = new(widget.Button)

			cs := cx.DuoUI.Gtx.Constraints
			cx.DuoUI.Layouts.View = &layout.Flex{Axis: layout.Vertical}
			cx.DuoUI.Layouts.Main = &layout.Flex{Axis: layout.Horizontal}
			in := layout.UniformInset(unit.Dp(30))

			header := cx.DuoUI.Layouts.View.Rigid(&cx.DuoUI.Gtx , func() {
				helpers.DuoUIdrawRect(&cx.DuoUI.Gtx , cs.Width.Max, 64, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf})
				helpers.DuoUIdrawRect(&cx.DuoUI.Gtx , 64, 64, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30})
				cx.DuoUI.Theme.IconButton(u.Ico.Logo).Layout(&cx.DuoUI.Gtx , u.Buttons.Logo)
			})

			//header := cx.DuoUI.Layouts.View.Rigid(&cx.DuoUI.Gtx, func(){ elem.DuoUIheader(cx.DuoUI)})

			main := cx.DuoUI.Layouts.View.Rigid(&cx.DuoUI.Gtx , func() {
				//balance := flh.Rigid(gtx, func() {
				//	in.Layout(gtx, func() {
				//		th.H3("balance :" + r.balance).Layout(gtx)
				//	})
				//})

				sidebar := cx.DuoUI.Layouts.Main.Rigid(&cx.DuoUI.Gtx , func() {
				
					helpers.DuoUIdrawRect(&cx.DuoUI.Gtx , 64, cs.Height.Max, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf})
				
					flm := layout.Flex{Axis: layout.Vertical}
					overview := flm.Rigid(&cx.DuoUI.Gtx , func() {
						cx.DuoUI.Theme.IconButton(u.Ico.Overview).Layout(&cx.DuoUI.Gtx , u.Buttons.Logo)
					})
					history := flm.Rigid(&cx.DuoUI.Gtx , func() {
						cx.DuoUI.Theme.IconButton(u.Ico.History).Layout(&cx.DuoUI.Gtx , u.Buttons.Logo)
					})
					network := flm.Rigid(&cx.DuoUI.Gtx , func() {
						cx.DuoUI.Theme.IconButton(u.Ico.Network).Layout(&cx.DuoUI.Gtx , u.Buttons.Logo)
					})
					settings := flm.Rigid(&cx.DuoUI.Gtx , func() {
						cx.DuoUI.Theme.IconButton(u.Ico.Settings).Layout(&cx.DuoUI.Gtx , u.Buttons.Logo)
					})
					flm.Layout(&cx.DuoUI.Gtx , overview, history, network, settings)
				
				})


				content := cx.DuoUI.Layouts.Main.Rigid(&cx.DuoUI.Gtx , func() {
					in.Layout(&cx.DuoUI.Gtx , func() {
						helpers.DuoUIdrawRect(&cx.DuoUI.Gtx , cs.Width.Max, cs.Height.Max, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf})
						cx.DuoUI.Theme.H3("balance :" + rc.Balance).Layout(&cx.DuoUI.Gtx )
					})
				})

				cx.DuoUI.Layouts.Main.Layout(&cx.DuoUI.Gtx , sidebar, content)
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

			cx.DuoUI.Layouts.View.Layout(&cx.DuoUI.Gtx , header, main)
			//cx.DuoUI.Layouts.View.Layout(&cx.DuoUI.Gtx , header, main)
			e.Frame(cx.DuoUI.Gtx .Ops)
		}
	}
}
