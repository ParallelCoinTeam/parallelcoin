package duoui

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
)

var (
	transList = &layout.List{
		Axis: layout.Vertical,
	}
	allTxs              = new(controller.CheckBox)
	mintedTxs           = new(controller.CheckBox)
	immatureTxs         = new(controller.CheckBox)
	sentTxs             = new(controller.CheckBox)
	receivedTxs         = new(controller.CheckBox)
	transactionsCounter = new(controller.Counter)
)


func (ui *DuoUI)  DuoUItransactions()() {
	ui.rc.Status.Wallet.Txs.ModelTxsListNumber = 55
	layout.Flex{
		Axis: layout.Vertical,
	}.Layout(ui.ly.Context,
		layout.Rigid(func() {
			cs := ui.ly.Context.Constraints
			theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 48, ui.ly.Theme.Color.Primary, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

			in := layout.UniformInset(unit.Dp(8))
			in.Layout(ui.ly.Context, func() {

				layout.Flex{
					Spacing: layout.SpaceBetween,
				}.Layout(ui.ly.Context,
					layout.Rigid(func() {
						layout.Flex{}.Layout(ui.ly.Context,
							layout.Rigid(func() {
								ui.ly.Theme.DuoUIcheckBox("ALL").Layout(ui.ly.Context, allTxs)
							}),
							layout.Rigid(func() {
								ui.ly.Theme.DuoUIcheckBox("MINTED").Layout(ui.ly.Context, mintedTxs)
							}),
							layout.Rigid(func() {
								ui.ly.Theme.DuoUIcheckBox("IMATURE").Layout(ui.ly.Context, immatureTxs)
							}),
							layout.Rigid(func() {
								ui.ly.Theme.DuoUIcheckBox("SENT").Layout(ui.ly.Context, sentTxs)
							}),
							layout.Rigid(func() {
								ui.ly.Theme.DuoUIcheckBox("RECEIVED").Layout(ui.ly.Context, receivedTxs)
							}),
						)
					}),
					layout.Rigid(func() {
						layout.Flex{}.Layout(ui.ly.Context,
							layout.Rigid(func() {

								//view.DuoUIcounter(duo)

							}),
						)
					}),
				)
			})
		}),
		layout.Flexed(1, func() {

			in := layout.UniformInset(unit.Dp(16))
			in.Layout(ui.ly.Context, func() {
				layout.Flex{
					Axis: layout.Vertical,
				}.Layout(ui.ly.Context,
					// Balance status item
					layout.Rigid(func() {
						cs := ui.ly.Context.Constraints
						transList.Layout(ui.ly.Context, len(ui.rc.Status.Wallet.LastTxs.Txs), func(i int) {
							// Invert list
							//i = len(txs.Txs) - 1 - i
							//
							t := ui.rc.Status.Wallet.LastTxs.Txs[i]
							a := 1.0
							//const duration = 5
							theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 1, "ff535353", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

							layout.Flex{
								Spacing: layout.SpaceBetween,
							}.Layout(ui.ly.Context,
								layout.Rigid(func() {
									layout.Flex{
										Axis: layout.Vertical,
									}.Layout(ui.ly.Context,
										layout.Rigid(func() {
											num := ui.ly.Theme.Body1(fmt.Sprint(i))
											num.Font.Typeface = "bariol"
											num.Color = helpers.Alpha(a, num.Color)
											num.Layout(ui.ly.Context)
										}),
										layout.Rigid(func() {
											tim := ui.ly.Theme.Body1(t.TxID)
											tim.Font.Typeface = "bariol"
											tim.Color = helpers.Alpha(a, tim.Color)
											tim.Layout(ui.ly.Context)
										}),
										layout.Rigid(func() {
											amount := ui.ly.Theme.H5(fmt.Sprintf("%0.8f", t.Amount))
											amount.Font.Typeface = "bariol"
											amount.Color = helpers.RGB(0x003300)
											amount.Color = helpers.Alpha(a, amount.Color)
											amount.Alignment = text.End
											amount.Font.Variant = "Mono"
											amount.Font.Weight = text.Bold
											amount.Layout(ui.ly.Context)
										}),
										layout.Rigid(func() {
											sat := ui.ly.Theme.Body1(t.Category)
											sat.Font.Typeface = "bariol"
											sat.Color = helpers.Alpha(a, sat.Color)
											sat.Layout(ui.ly.Context)
										}),
										layout.Rigid(func() {
											l := ui.ly.Theme.Body2(t.Time)
											l.Font.Typeface = "bariol"
											l.Color = theme.HexARGB(ui.ly.Theme.Color.Hint)
											l.Color = helpers.Alpha(a, l.Color)
											l.Layout(ui.ly.Context)
										}),
									)
								}),
								layout.Rigid(func() {
									sat := ui.ly.Theme.Body1(fmt.Sprintf("%0.8f", t.Amount))
									sat.Font.Typeface = "bariol"
									sat.Color = helpers.Alpha(a, sat.Color)
									sat.Layout(ui.ly.Context)
								}),
							)
						})
					}))
			})
		}),
	)
}
