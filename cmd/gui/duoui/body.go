package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"

	"image/color"
)

func DuoUIbody(duo *DuoUI) {
	duo.comp.Body.Layout.Layout(duo.gc,
		layout.Rigid(func() {
			cs := duo.gc.Constraints
			helpers.DuoUIdrawRectangle(duo.gc, cs.Width.Max, 64, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}, 0, 0, 0, 0, unit.Dp(0))
			DuoUIsidebar(duo)
		}),
		layout.Flexed(1, func() {
			DuoUIcontent(duo)
		}),
	)
}
