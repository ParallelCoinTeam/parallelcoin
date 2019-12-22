package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
	"image/color"
)

var (
	inLogo = layout.UniformInset(unit.Dp(4))
)

func DuoUIheader(duo *DuoUI) {
	layout.Flex{}.Layout(duo.gc,
		layout.Rigid(func() {
			helpers.DuoUIdrawRect(duo.gc, duo.cs.Width.Max, 64, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}, 0, 0, 0, 0)
			// Header <<<
			layout.Flex{}.Layout(duo.gc,
				layout.Rigid(func() {
					helpers.DuoUIdrawRect(duo.gc, 64, 64, color.RGBA{A: 0xff, R: 0x30, B: 0x30, G: 0x30}, 0, 0, 0, 0)
					//inLogo.Layout(duo.gc, func() {

					duo.ico.Logo.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
					duo.ico.Logo.Layout(duo.gc, unit.Dp(64))
					//})
				}),
			)
			layout.Flex{}.Layout(duo.gc,
				layout.Rigid(func() {
					duo.th.H5(duo.rc.Balance + " DUO").Layout(duo.gc)
				}),
			)

			// Header >>>
		}),
	)
}
