package elem

import (
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/p9c/pod/cmd/gui/assets/ico"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
)

func DuoUIsidebar(cx *conte.Xt) layout.FlexChild {

	return cx.DuoUI.Layouts.Main.Rigid(cx.Gtx , func() {
		helpers.DuoUIdrawRect(cx.Gtx , 64, cs.Height.Max, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf})

		flm := layout.Flex{Axis: layout.Vertical}
		overview := flm.Rigid(cx.Gtx , func() {
			th.IconButton(u.Ico.Overview).Layout(cx.Gtx , u.Buttons.Logo)
		})
		history := flm.Rigid(cx.Gtx , func() {
			th.IconButton(u.Ico.History).Layout(cx.Gtx , u.Buttons.Logo)
		})
		network := flm.Rigid(cx.Gtx , func() {
			th.IconButton(u.Ico.Network).Layout(cx.Gtx , u.Buttons.Logo)
		})
		settings := flm.Rigid(cx.Gtx , func() {
			th.IconButton(u.Ico.Settings).Layout(cx.Gtx , u.Buttons.Logo)
		})
		flm.Layout(cx.Gtx , overview, history, network, settings)

	})
}
