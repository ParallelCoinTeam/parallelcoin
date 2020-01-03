package duoui

import (
	"github.com/p9c/pod/cmd/gui/componentsWidgets"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/layout"
)

func DuoUIhistory(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	layout.Flex{}.Layout(duo.DuoUIcontext,
		layout.Flexed(1, func() {
			componentsWidgets.DuoUItransactionsWidget(duo, cx, rc)
		}),
	)
}
