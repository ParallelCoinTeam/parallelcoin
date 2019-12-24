package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/cmd/gui/widgets"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
	"image/color"
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
					helpers.DuoUIdrawRectangle(duo.Gc, duo.Cs.Width.Max-20, 180, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 9.9, 9.9, 9.9, 9.9, unit.Dp(0))
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
						helpers.DuoUIdrawRectangle(duo.Gc, duo.Cs.Width.Max-30, cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 9.9, 9.9, 9.9, 9.9, unit.Dp(0))
						widgets.DuoUIlatestTxsWidget(duo, cx, rc)

					}),
					layout.Flexed(0.24, func() {
						helpers.DuoUIdrawRectangle(duo.Gc, duo.Cs.Width.Max, cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 9.9, 9.9, 9.9, 9.9, unit.Dp(0))
						widgets.DuoUIstatusWidget(duo, rc)
					}))
				// OverviewBottom >>>
			})
		}),
	)
}
