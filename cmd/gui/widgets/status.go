package widgets

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
)

func DuoUIstatusWidget(duo *models.DuoUI, rc *rcd.RcVar){
	in := layout.UniformInset(unit.Dp(30))
	in.Layout(duo.Gc, func() {
		duo.Comp.Status.Layout.Layout(duo.Gc,
			// Balance status item
			layout.Rigid(func() {
				duo.Comp.StatusItem.Layout.Layout(duo.Gc,
					layout.Rigid(func() {
						balanceTxt := duo.Th.H5("Balance :")
						balanceTxt.Color = duo.Conf.StatusTextColor
						balanceTxt.Layout(duo.Gc)
					}),
					layout.Rigid(func() {
						balance := duo.Th.H5(rc.Balance + " " + duo.Conf.Abbrevation)
						balance.Color = duo.Conf.StatusTextColor
						balance.Layout(duo.Gc)
					}),

				)
			}),
			// Block height status item
			layout.Rigid(func() {
				duo.Comp.StatusItem.Layout.Layout(duo.Gc,
					layout.Rigid(func() {
						blockheightTxt := duo.Th.H5("Block Height :")
						blockheightTxt.Color = duo.Conf.StatusTextColor
						blockheightTxt.Layout(duo.Gc)
					}),
					layout.Rigid(func() {
						blockheightVal := duo.Th.H5(fmt.Sprint(rc.BlockHeight))
						blockheightVal.Color = duo.Conf.StatusTextColor
						blockheightVal.Layout(duo.Gc)
					}),

				)
			}),

			// Difficulty height status item
			layout.Rigid(func() {
				duo.Comp.StatusItem.Layout.Layout(duo.Gc,
					layout.Rigid(func() {
						difficulty := duo.Th.H5("Difficulty :")
						difficulty.Color = duo.Conf.StatusTextColor
						difficulty.Layout(duo.Gc)
					}),
					layout.Rigid(func() {
						difficulty := duo.Th.H5(fmt.Sprintf("%f", rc.Difficulty))
						difficulty.Color = duo.Conf.StatusTextColor
						difficulty.Layout(duo.Gc)
					}),

				)
			}),

			// Connections status item
			layout.Rigid(func() {
				duo.Comp.StatusItem.Layout.Layout(duo.Gc,
					layout.Rigid(func() {
						connections := duo.Th.H5("Connections :")
						connections.Color = duo.Conf.StatusTextColor
						connections.Layout(duo.Gc)
					}),
					layout.Rigid(func() {
						connections := duo.Th.H5(fmt.Sprint(rc.Connections))
						connections.Color = duo.Conf.StatusTextColor
						connections.Layout(duo.Gc)
					}),

				)
			}),

		)
	})
}