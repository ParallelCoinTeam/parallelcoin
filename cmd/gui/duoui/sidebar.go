package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"

	"image/color"
)

func DuoUIsidebar(duo *DuoUI) {
	duo.comp.Sidebar.Layout.Layout(duo.gc,
		layout.Rigid(func() {
			helpers.DuoUIdrawRectangle(duo.gc, 64, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}, 0, 0, 0, 0, unit.Dp(0))
			DuoUImenu(duo)
		}),
	)
}
