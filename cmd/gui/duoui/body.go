package duoui

import (
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/layout"
)

func DuoUIbody(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	duo.Comp.Body.Layout.Layout(duo.Gc,
		layout.Rigid(func() {
			DuoUIsidebar(duo, cx, rc)
		}),
		layout.Flexed(1, func() {
			DuoUIcontent(duo,cx,rc)
		}),
	)
}
