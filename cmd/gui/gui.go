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
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
)

func WalletGUI(cx *conte.Xt) (err error) {
	r := &rcvar{
		cx:         cx,
		blockcount: 0,
	}
	go func() {
		cx.Window = app.NewWindow()
		if err := loop(cx, r); err != nil {
			log.FATAL(err)
		}
		//runUI()
	}()
	app.Main()
	return
}

func loop(cx *conte.Xt, r *rcvar) error {
	gofont.Register()
	cx.DuoUI.Theme = material.NewTheme()
	cx.Gtx = layout.NewContext(cx.Window.Queue())
	for {
		e := <-cx.Window.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			cx.Gtx.Reset(e.Config, e.Size)
			r.GetDuOSbalance()
			r.GetDuOSheight()
			r.GetDuOStatus()
			r.GetDuOSlocalLost()
			r.GetDuOSdifficulty()
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

			cs := cx.Gtx.Constraints
			cx.DuoUI.Layouts.View = &layout.Flex{Axis: layout.Vertical}
			cx.DuoUI.Layouts.Main = &layout.Flex{Axis: layout.Horizontal}
			in := layout.UniformInset(unit.Dp(30))
			header := cx.DuoUI.Layouts.View.Rigid(cx.Gtx , func() {
				helpers.DuoUIdrawRect(cx.Gtx , cs.Width.Max, 64, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf})
				helpers.DuoUIdrawRect(cx.Gtx , 64, 64, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30})
				cx.DuoUI.Theme.IconButton(u.Ico.Logo).Layout(cx.Gtx , u.Buttons.Logo)
				//th.Image(Duo)

			})
			main := cx.DuoUI.Layouts.View.Rigid(cx.Gtx , func() {
				//balance := flh.Rigid(gtx, func() {
				//	in.Layout(gtx, func() {
				//		th.H3("balance :" + r.balance).Layout(gtx)
				//	})
				//})
				sidebar := cx.DuoUI.Layouts.Main.Rigid(cx.Gtx , func() {

					helpers.DuoUIdrawRect(cx.Gtx , 64, cs.Height.Max, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf})

					flm := layout.Flex{Axis: layout.Vertical}
					overview := flm.Rigid(cx.Gtx , func() {
						cx.DuoUI.Theme.IconButton(u.Ico.Overview).Layout(cx.Gtx , u.Buttons.Logo)
					})
					history := flm.Rigid(cx.Gtx , func() {
						cx.DuoUI.Theme.IconButton(u.Ico.History).Layout(cx.Gtx , u.Buttons.Logo)
					})
					network := flm.Rigid(cx.Gtx , func() {
						cx.DuoUI.Theme.IconButton(u.Ico.Network).Layout(cx.Gtx , u.Buttons.Logo)
					})
					settings := flm.Rigid(cx.Gtx , func() {
						cx.DuoUI.Theme.IconButton(u.Ico.Settings).Layout(cx.Gtx , u.Buttons.Logo)
					})
					flm.Layout(cx.Gtx , overview, history, network, settings)

				})
				content := cx.DuoUI.Layouts.Main.Rigid(cx.Gtx , func() {
					in.Layout(cx.Gtx , func() {
						helpers.DuoUIdrawRect(cx.Gtx , cs.Width.Max, cs.Height.Max, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf})
						cx.DuoUI.Theme.H3("balance :" + r.balance).Layout(cx.Gtx )
					})
				})

				cx.DuoUI.Layouts.Main.Layout(cx.Gtx , sidebar, content)
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

			cx.DuoUI.Layouts.View.Layout(cx.Gtx , header, main)
			e.Frame(cx.Gtx .Ops)
		}
	}
}
