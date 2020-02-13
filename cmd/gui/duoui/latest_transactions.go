package duoui

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
)

var (
	latestTransList = &layout.List{
		Axis: layout.Vertical,
	}
)

func (ui *DuoUI) DuoUIlatestTransactions() {
	layout.Flex{
		Axis: layout.Vertical,
	}.Layout(ui.ly.Context,
		layout.Rigid(func() {
			cs := ui.ly.Context.Constraints
			theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 48, ui.ly.Theme.Color.Primary, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

			in := layout.UniformInset(unit.Dp(8))
			in.Layout(ui.ly.Context, func() {

				latestx := ui.ly.Theme.H5("LATEST TRANSACTIONS")
				latestx.Color = theme.HexARGB("ffcfcfcf")
				latestx.Alignment = text.Start
				latestx.Layout(ui.ly.Context)
			})
		}),
		layout.Flexed(1, func() {

			in := layout.UniformInset(unit.Dp(16))
			in.Layout(ui.ly.Context, func() {
				layout.Flex{Axis: layout.Vertical}.Layout(ui.ly.Context,
					layout.Rigid(func() {
						cs := ui.ly.Context.Constraints
						latestTransList.Layout(ui.ly.Context, len(ui.rc.Status.Wallet.Txs.Txs), func(i int) {
							// Invert list
							//i = len(txs.Txs) - 1 - i
							t := ui.rc.Status.Wallet.Txs.Txs[i]
							a := 1.0
							//const duration = 5
							theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, cs.Height.Max, ui.ly.Theme.Color.Text, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
							layout.Flex{
								Spacing: layout.SpaceBetween,
							}.Layout(ui.ly.Context,
								layout.Rigid(func() {
									in := layout.UniformInset(unit.Dp(15))
									in.Layout(ui.ly.Context, func() {

										layout.Flex{
											Axis: layout.Vertical,
										}.Layout(ui.ly.Context,
											layout.Rigid(func() {

												tim := ui.ly.Theme.Caption(t.TxID)
												tim.Font.Typeface = "bariol"
												tim.Color = helpers.RGB(0xcfcfcf)
												tim.Color = helpers.Alpha(a, tim.Color)
												tim.Layout(ui.ly.Context)
											}),
											layout.Rigid(func() {
												amount := ui.ly.Theme.H5(fmt.Sprintf("%0.8f", t.Amount))
												amount.Color = helpers.RGB(0x003300)
												amount.Color = helpers.Alpha(a, amount.Color)
												amount.Alignment = text.End
												amount.Font.Typeface = "bariol"
												amount.Font.Variant = "Bold"
												amount.Font.Weight = text.Bold
												amount.Layout(ui.ly.Context)
											}),
											layout.Rigid(func() {
												sat := ui.ly.Theme.Body1(t.Category)
												sat.Color = helpers.RGB(0xcfcfcf)
												sat.Font.Typeface = "bariol"

												sat.Color = helpers.Alpha(a, sat.Color)
												sat.Layout(ui.ly.Context)
											}),
											layout.Rigid(func() {

												l := ui.ly.Theme.Body1(t.Time)
												l.Font.Typeface = "bariol"
												l.Color = helpers.RGB(0xcfcfcf)
												l.Color = theme.HexARGB(ui.ly.Theme.Color.Hint)
												l.Color = helpers.Alpha(a, l.Color)
												l.Layout(ui.ly.Context)
											}),
										)
									})
								}),
								layout.Rigid(func() {
									in := layout.UniformInset(unit.Dp(15))
									in.Layout(ui.ly.Context, func() {
										sat := ui.ly.Theme.H6(fmt.Sprintf("%0.8f", t.Amount))
										sat.Font.Typeface = "bariol"
										sat.Color = helpers.RGB(0xcfcfcf)
										sat.Color = helpers.Alpha(a, sat.Color)
										sat.Layout(ui.ly.Context)
									})

								}),
							)
						})
					}))
			})
		}),
	)
}
