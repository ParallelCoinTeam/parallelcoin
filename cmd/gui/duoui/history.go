package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"

	"image/color"
)

func DuoUIhistory(duo *DuoUI) {
	layout.Flex{}.Layout(duo.gc,
		layout.Flexed(1, func() {
			duo.comp.History.Inset.Layout(duo.gc, func() {
				helpers.DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0xcf, B: 0xcf}, 0, 0, 0, 0)
				// Overview <<<
				in := layout.UniformInset(unit.Dp(60))
				in.Layout(duo.gc, func() {
					helpers.DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0xcf, G: 0x30, B: 0xcf}, 0, 0, 0, 0)

					duo.th.H5("history :").Layout(duo.gc)
				})
				// Overview >>>
			})
		}),
	)
}
