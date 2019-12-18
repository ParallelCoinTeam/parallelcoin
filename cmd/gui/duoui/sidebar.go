package duoui

import (
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/cmd/gui/helpers"
	"image/color"
)

func DuoUIsidebar(duo *DuoUI) layout.FlexChild {
	return duo.comp.Body.Layout.Rigid(duo.gc, func() {
		helpers.DuoUIdrawRect(duo.gc, 64, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}, 0, 0, 0, 0)

		duo.comp.Sidebar.Layout.Layout(duo.gc, DuoUImenu(duo))
	})
}
