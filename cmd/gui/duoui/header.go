package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
	"image/color"
)

var (
	inLogo = layout.Stack{Alignment: layout.Center}
)

func DuoUIheader(duo *DuoUI) {
	// Header <<<
	duo.comp.Header.Layout.Layout(duo.gc,
		layout.Rigid(func() {
			helpers.DuoUIdrawRectangle(duo.gc, 64, 64, color.RGBA{A: 0xff, R: 0x30, B: 0x30, G: 0x30}, 0, 0, 0, 0, unit.Dp(0))
			layout.Align(layout.Center).Layout(duo.gc, func() {
				layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(5), Right: unit.Dp(4)}.Layout(duo.gc, func() {

					duo.ico.Logo.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
					duo.ico.Logo.Layout(duo.gc, unit.Dp(64))
				})
			})

		}),
		layout.Rigid(func() {
			layout.Align(layout.Center).Layout(duo.gc, func() {
				layout.Inset{Top: unit.Dp(16), Bottom: unit.Dp(16), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(duo.gc, func() {
					duo.th.H5(duo.rc.Balance + " " + duo.conf.Abbrevation).Layout(duo.gc)
				})
			})
		}))

	// Header >>>

}