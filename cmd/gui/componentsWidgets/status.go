package componentsWidgets

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
)

func DuoUIstatusWidget(duo *models.DuoUI, rc *rcd.RcVar){
	in := layout.UniformInset(unit.Dp(30))
	in.Layout(duo.DuoUIcontext, func() {
		duo.DuoUIcomponents.Status.Layout.Layout(duo.DuoUIcontext,
			// Balance status item
			layout.Rigid(func() {
				duo.DuoUIcomponents.StatusItem.Layout.Layout(duo.DuoUIcontext,
					layout.Rigid(func() {
						balanceTxt := duo.DuoUItheme.H5("Balance :")
						balanceTxt.Color = duo.DuoUIconfiguration.PrimaryTextColor
						balanceTxt.Layout(duo.DuoUIcontext)
					}),
					layout.Rigid(func() {
						balance := duo.DuoUItheme.H5(rc.Balance + " " + duo.DuoUIconfiguration.Abbrevation)
						balance.Color = duo.DuoUIconfiguration.PrimaryTextColor
						balance.Layout(duo.DuoUIcontext)
					}),

				)
			}),
			// Block height status item
			layout.Rigid(func() {
				duo.DuoUIcomponents.StatusItem.Layout.Layout(duo.DuoUIcontext,
					layout.Rigid(func() {
						blockheightTxt := duo.DuoUItheme.H5("Block Height :")
						blockheightTxt.Color = duo.DuoUIconfiguration.PrimaryTextColor
						blockheightTxt.Layout(duo.DuoUIcontext)
					}),
					layout.Rigid(func() {
						blockheightVal := duo.DuoUItheme.H5(fmt.Sprint(rc.BlockHeight))
						blockheightVal.Color = duo.DuoUIconfiguration.PrimaryTextColor
						blockheightVal.Layout(duo.DuoUIcontext)
					}),

				)
			}),

			// Difficulty height status item
			layout.Rigid(func() {
				duo.DuoUIcomponents.StatusItem.Layout.Layout(duo.DuoUIcontext,
					layout.Rigid(func() {
						difficulty := duo.DuoUItheme.H5("Difficulty :")
						difficulty.Color = duo.DuoUIconfiguration.PrimaryTextColor
						difficulty.Layout(duo.DuoUIcontext)
					}),
					layout.Rigid(func() {
						difficulty := duo.DuoUItheme.H5(fmt.Sprintf("%f", rc.Difficulty))
						difficulty.Color = duo.DuoUIconfiguration.PrimaryTextColor
						difficulty.Layout(duo.DuoUIcontext)
					}),

				)
			}),

			// Connections status item
			layout.Rigid(func() {
				duo.DuoUIcomponents.StatusItem.Layout.Layout(duo.DuoUIcontext,
					layout.Rigid(func() {
						connections := duo.DuoUItheme.H5("Connections :")
						connections.Color = duo.DuoUIconfiguration.PrimaryTextColor
						connections.Layout(duo.DuoUIcontext)
					}),
					layout.Rigid(func() {
						connections := duo.DuoUItheme.H5(fmt.Sprint(rc.ConnectionCount))
						connections.Color = duo.DuoUIconfiguration.PrimaryTextColor
						connections.Layout(duo.DuoUIcontext)
					}),

				)
			}),

		)
	})
}