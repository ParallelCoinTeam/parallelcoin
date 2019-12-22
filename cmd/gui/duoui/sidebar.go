package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/pkg/gio/layout"
	"image/color"
)

func DuoUIsidebar(duo *DuoUI) {
	layout.Flex{}.Layout(duo.gc,
		layout.Rigid(func() {
			helpers.DuoUIdrawRect(duo.gc, 64, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}, 0, 0, 0, 0)
			DuoUImenu(duo)
		}),
	)
}
