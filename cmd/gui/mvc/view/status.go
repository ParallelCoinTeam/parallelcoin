package view

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/unit"
)

type DuoUIstatus struct {
	*model.DuOScomponent
}

func DuoCOMstatus() *DuoUIstatus {
	settings := *new(DuoUIstatus)
	cs := &model.DuOScomponent{
		Name:    "status",
		Version: "0.1",
		//Model:      ,
		//Controller: c,
	}
	*settings.DuOScomponent = *cs
	return &settings
}

func (b *DuoUIstatus) View(gtx *layout.Context, th *theme.DuoUItheme, stat *model.DuoUIstatus, abbr string) func() {
	return func() {
		in := layout.UniformInset(unit.Dp(30))
		in.Layout(gtx, func() {
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				// Balance status item
				layout.Rigid(func() {
					layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func() {
							balanceTxt := th.H5("Balance :")
							balanceTxt.Color = th.Color.Primary
							balanceTxt.Layout(gtx)
						}),
						layout.Rigid(func() {
							balance := th.H5(stat.Wallet.Balance + " " + abbr)
							balance.Color = th.Color.Primary
							balance.Layout(gtx)
						}),

					)
				}),
				// Block height status item
				layout.Rigid(func() {
					layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func() {
							blockheightTxt := th.H5("Block Height :")
							blockheightTxt.Color = th.Color.Primary
							blockheightTxt.Layout(gtx)
						}),
						layout.Rigid(func() {
							blockheightVal := th.H5(fmt.Sprint(stat.Node.BlockHeight))
							blockheightVal.Color = th.Color.Primary
							blockheightVal.Layout(gtx)
						}),

					)
				}),

				// Difficulty height status item
				layout.Rigid(func() {
					layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func() {
							difficulty := th.H5("Difficulty :")
							difficulty.Color = th.Color.Primary
							difficulty.Layout(gtx)
						}),
						layout.Rigid(func() {
							difficulty := th.H5(fmt.Sprintf("%f", stat.Node.Difficulty))
							difficulty.Color = th.Color.Primary
							difficulty.Layout(gtx)
						}),

					)
				}),

				// Connections status item
				layout.Rigid(func() {
					layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func() {
							connections := th.H5("Connections :")
							connections.Color = th.Color.Primary
							connections.Layout(gtx)
						}),
						layout.Rigid(func() {
							connections := th.H5(fmt.Sprint(stat.Node.ConnectionCount))
							connections.Color = th.Color.Primary
							connections.Layout(gtx)
						}),

					)
				}),

			)
		})
	}
}
