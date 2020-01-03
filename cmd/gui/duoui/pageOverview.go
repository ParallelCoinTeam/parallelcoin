package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/cmd/gui/components"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
)

func DuoUIoverview(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	rc.GetDuoUIbalance(duo, cx)
	rc.GetDuoUIunconfirmedBalance(duo, cx)
	rc.GetDuoUIblockHeight(duo, cx)
	rc.GetDuoUItatus(duo, cx)
	rc.GetDuoUIlocalLost(duo)
	rc.GetDuoUIdifficulty(duo, cx)

	duo.DuoUIcomponents.Overview.Layout.Layout(duo.DuoUIcontext,
		layout.Flexed(0.5, func() {
			cs := duo.DuoUIcontext.Constraints
			helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, helpers.HexARGB("ffcfcfcf"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
			components.DuoUIbalanceWidget(duo, rc)

		}),
		layout.Flexed(0.5, func() {
			cs := duo.DuoUIcontext.Constraints
			helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, helpers.HexARGB("ff424242"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
			components.DuoUIlatestTxsWidget(duo, cx, rc)
		}),
	)
}
