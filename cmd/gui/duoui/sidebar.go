package duoui

import (
	"gioui.org/layout"
	"image/color"
)

func DuoUIsidebar(duo *DuoUI) layout.FlexChild {
	return duo.comp.body.l.Rigid(duo.gc, func() {
		DuoUIdrawRect(duo.gc, 64, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}, 0, 0, 0, 0)

		duo.comp.sidebar.l.Layout(duo.gc, DuoUImenu(duo))
	})
}
