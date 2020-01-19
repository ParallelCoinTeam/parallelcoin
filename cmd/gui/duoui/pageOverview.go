package duoui

import (
	"github.com/p9c/pod/cmd/gui/componentsWidgets"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/layout"
)

func DuoUIoverview(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	rc.GetDuoUIbalance(duo, cx)
	rc.GetDuoUIunconfirmedBalance(duo, cx)
	rc.GetDuoUIblockHeight(duo, cx)
	rc.GetDuoUIstatus(duo, cx)
	rc.GetDuoUIlocalLost(duo)
	rc.GetDuoUIdifficulty(duo, cx)

	viewport := layout.Flex{Axis: layout.Horizontal}

	if duo.DuoUIcontext.Constraints.Width.Max < 1024 {
		viewport = layout.Flex{Axis: layout.Vertical}
	}

	viewport.Layout(duo.DuoUIcontext,
		layout.Flexed(0.5, func() {
			cs := duo.DuoUIcontext.Constraints
			helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, helpers.HexARGB("ffcfcfcf"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			componentsWidgets.DuoUIbalanceWidget(duo, rc)

		}),
		layout.Flexed(0.5, func() {
			cs := duo.DuoUIcontext.Constraints
			helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, helpers.HexARGB("ff424242"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			componentsWidgets.DuoUIlatestTxsWidget(duo, cx, rc)
		}),
	)
}
