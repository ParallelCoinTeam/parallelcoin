package duoui

import (
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/cmd/gui/widgets"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/layout"
)

func DuoUIhistory(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	layout.Flex{}.Layout(duo.Gc,
		layout.Flexed(1, func() {

			widgets.DuoUItransactionsWidget(duo, cx, rc)
				// Overview >>>

		}),
	)
}
