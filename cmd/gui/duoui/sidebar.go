package duoui

import (
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/layout"
)

func DuoUIsidebar(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	duo.DuoUIcomponents.Sidebar.Layout.Layout(duo.DuoUIcontext,
		layout.Rigid(func() {
			DuoUImenu(duo, cx, rc)
		}),
	)
}
