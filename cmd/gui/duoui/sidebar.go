package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
)

func DuoUIsidebar(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	duo.Comp.Sidebar.Layout.Layout(duo.Gc,
		layout.Rigid(func() {
			helpers.DuoUIdrawRectangle(duo.Gc, 64, duo.Cs.Height.Max, helpers.HexARGB("ffcfcfcf"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
			DuoUImenu(duo, cx, rc)
		}),
	)
}
