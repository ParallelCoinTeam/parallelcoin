package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/cmd/gui/widgets"
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

	duo.Comp.Overview.Layout.Layout(duo.Gc,
		layout.Rigid(func() {
			// OverviewTop <<<
			duo.Comp.OverviewTop.Layout.Layout(duo.Gc,
				layout.Flexed(0.38, func() {
					helpers.DuoUIdrawRectangle(duo.Gc, duo.Cs.Width.Max-20, 180, helpers.HexARGB("ff303030"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
					widgets.DuoUIbalanceWidget(duo, rc)

				}),
				layout.Flexed(0.62, func() {
					widgets.DuoUIsendreceive(duo, cx, rc)
				}))
			// OverviewTop >>>
		}),
		layout.Flexed(1, func() {
			// OverviewBottom <<<
			in := layout.Inset{
				Top: unit.Dp(20),
			}
			in.Layout(duo.Gc, func() {
				cs := duo.Gc.Constraints
				duo.Comp.OverviewBottom.Layout.Layout(duo.Gc,
					layout.Flexed(0.76, func() {
						helpers.DuoUIdrawRectangle(duo.Gc, duo.Cs.Width.Max-30, cs.Height.Max, helpers.HexARGB("ff303030"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
						widgets.DuoUIlatestTxsWidget(duo, cx, rc)

					}),
					layout.Flexed(0.24, func() {
						helpers.DuoUIdrawRectangle(duo.Gc, duo.Cs.Width.Max, cs.Height.Max, helpers.HexARGB("ff303030"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
						widgets.DuoUIstatusWidget(duo, rc)
					}))
				// OverviewBottom >>>
			})
		}),
	)
}
