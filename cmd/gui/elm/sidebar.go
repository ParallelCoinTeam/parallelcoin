package elm

import (
	"gioui.org/layout"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/mod"
	"image/color"
)

func DuoUIsidebar(d *mod.DuoUI) layout.FlexChild {
	return d.Layouts.Main.Rigid(&d.Gtx, func() {
		hlp.DuoUIdrawRect(&d.Gtx, 64, d.Gtx.Constraints.Height.Max, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf})
		flm := layout.Flex{Axis: layout.Vertical}
		overview := flm.Rigid(&d.Gtx, func() {
			d.Theme.IconButton(d.Ico.Overview).Layout(&d.Gtx, d.Buttons.Logo)
		})
		history := flm.Rigid(&d.Gtx, func() {
			d.Theme.IconButton(d.Ico.History).Layout(&d.Gtx, d.Buttons.Logo)
		})
		network := flm.Rigid(&d.Gtx, func() {
			d.Theme.IconButton(d.Ico.Network).Layout(&d.Gtx, d.Buttons.Logo)
		})
		settings := flm.Rigid(&d.Gtx, func() {
			d.Theme.IconButton(d.Ico.Settings).Layout(&d.Gtx, d.Buttons.Logo)
		})
		flm.Layout(&d.Gtx, overview, history, network, settings)

	})
}
