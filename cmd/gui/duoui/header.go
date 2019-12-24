package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
	"image/color"
)

var (
	inLogo = layout.Stack{Alignment: layout.Center}
)

func DuoUIheader(duo *models.DuoUI, rc *rcd.RcVar) {
	// Header <<<
	duo.Comp.Header.Layout.Layout(duo.Gc,
		layout.Rigid(func() {
			helpers.DuoUIdrawRectangle(duo.Gc, 64, 64, color.RGBA{A: 0xff, R: 0x30, B: 0x30, G: 0x30}, 0, 0, 0, 0, unit.Dp(0))
			layout.Align(layout.Center).Layout(duo.Gc, func() {
				layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(5), Right: unit.Dp(4)}.Layout(duo.Gc, func() {

					duo.Ico.Logo.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
					duo.Ico.Logo.Layout(duo.Gc, unit.Dp(64))
				})
			})

		}),
		layout.Rigid(func() {
			layout.Align(layout.Center).Layout(duo.Gc, func() {
				layout.Inset{Top: unit.Dp(16), Bottom: unit.Dp(16), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(duo.Gc, func() {
					duo.Th.H5(rc.Balance + " " + duo.Conf.Abbrevation).Layout(duo.Gc)
				})
			})
		}))

	// Header >>>

}