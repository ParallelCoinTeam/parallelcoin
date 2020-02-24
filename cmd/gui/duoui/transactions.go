package duoui

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
)

var (
	transList = &layout.List{
		Axis: layout.Vertical,
	}
	allTxs      = new(controller.CheckBox)
	mintedTxs   = new(controller.CheckBox)
	immatureTxs = new(controller.CheckBox)
	sentTxs     = new(controller.CheckBox)
	receivedTxs = new(controller.CheckBox)
	itemValue   = &controller.DuoUIcounter{
		Value:        11,
		OperateValue: 1,
		From:         0,
		To:           15,
	}
)

func (ui *DuoUI) txsFilter() func() {
	return func() {
		layout.Flex{}.Layout(ui.ly.Context,
			layout.Rigid(ui.txsFilterItem("ALL", allTxs)),
			layout.Rigid(ui.txsFilterItem("MINTED", mintedTxs)),
			layout.Rigid(ui.txsFilterItem("IMATURE", immatureTxs)),
			layout.Rigid(ui.txsFilterItem("SENT", sentTxs)),
			layout.Rigid(ui.txsFilterItem("RECEIVED", receivedTxs)))
	}
}

func (ui *DuoUI) txsFilterItem(id string, c *controller.CheckBox) func() {
	return func() {
		ui.ly.Theme.DuoUIcheckBox(id, ui.ly.Theme.Color.Light, ui.ly.Theme.Color.Light).Layout(ui.ly.Context, c)
	}
}

func (ui *DuoUI) txsDetails(i int, t *model.DuoUItx) func() {
	return func() {
		layout.Flex{
			Axis: layout.Vertical,
		}.Layout(ui.ly.Context,
			layout.Rigid(func() {
				num := ui.ly.Theme.Body1(fmt.Sprint(i))
				num.Font.Typeface = ui.ly.Theme.Font.Primary
				num.Color = theme.HexARGB(ui.ly.Theme.Color.Hint)
				num.Layout(ui.ly.Context)
			}),
			layout.Rigid(func() {
				tim := ui.ly.Theme.Body1(t.TxID)
				tim.Font.Typeface = ui.ly.Theme.Font.Primary
				tim.Color = theme.HexARGB(ui.ly.Theme.Color.Hint)
				tim.Layout(ui.ly.Context)
			}),
			layout.Rigid(func() {
				amount := ui.ly.Theme.H5(fmt.Sprintf("%0.8f", t.Amount))
				amount.Font.Typeface = ui.ly.Theme.Font.Primary
				amount.Color = theme.HexARGB(ui.ly.Theme.Color.Hint)
				amount.Alignment = text.End
				amount.Font.Variant = "Mono"
				amount.Font.Weight = text.Bold
				amount.Layout(ui.ly.Context)
			}),
			layout.Rigid(func() {
				sat := ui.ly.Theme.Body1(t.Category)
				sat.Font.Typeface = ui.ly.Theme.Font.Primary
				sat.Color = theme.HexARGB(ui.ly.Theme.Color.Hint)
				sat.Layout(ui.ly.Context)
			}),
			layout.Rigid(func() {
				l := ui.ly.Theme.Body2(t.Time)
				l.Font.Typeface = ui.ly.Theme.Font.Primary
				l.Color = theme.HexARGB(ui.ly.Theme.Color.Hint)
				l.Layout(ui.ly.Context)
			}),
		)
	}
}

func (ui *DuoUI) contentHeader(b func()) func() {
	return func() {
		hmin := ui.ly.Context.Constraints.Width.Min
		vmin := ui.ly.Context.Constraints.Height.Min
		layout.Stack{Alignment: layout.Center}.Layout(ui.ly.Context,
			layout.Expanded(func() {
				clip.Rect{
					Rect: f32.Rectangle{Max: f32.Point{
						X: float32(ui.ly.Context.Constraints.Width.Min),
						Y: float32(ui.ly.Context.Constraints.Height.Min),
					}},
				}.Op(ui.ly.Context.Ops).Add(ui.ly.Context.Ops)
				fill(ui.ly.Context, theme.HexARGB(ui.ly.Theme.Color.Primary))
			}),
			layout.Stacked(func() {
				ui.ly.Context.Constraints.Width.Min = hmin
				ui.ly.Context.Constraints.Height.Min = vmin
				layout.UniformInset(unit.Dp(8)).Layout(ui.ly.Context, b)
			}),
		)
	}
}

func (ui *DuoUI) headerTransactions() func() {
	return func() {
		layout.Flex{
			Spacing: layout.SpaceBetween,
		}.Layout(ui.ly.Context,
			layout.Rigid(ui.txsFilter()),
			layout.Rigid(func() {
				ui.ly.Theme.DuoUIcounter().Layout(ui.ly.Context, itemValue)
			}),
		)
	}
}

func (ui *DuoUI) txsBody() func() {
	return func() {
		layout.UniformInset(unit.Dp(16)).Layout(ui.ly.Context, func() {
			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(ui.ly.Context,
				layout.Rigid(func() {
					cs := ui.ly.Context.Constraints
					transList.Layout(ui.ly.Context, len(ui.rc.Status.Wallet.LastTxs.Txs), func(i int) {
						t := ui.rc.Status.Wallet.LastTxs.Txs[i]
						theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 1, "ff535353", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
						layout.Flex{
							Spacing: layout.SpaceBetween,
						}.Layout(ui.ly.Context,
							layout.Rigid(ui.txsDetails(i, &t)),
							layout.Rigid(func() {
								sat := ui.ly.Theme.Body1(fmt.Sprintf("%0.8f", t.Amount))
								sat.Font.Typeface = ui.ly.Theme.Font.Primary
								sat.Color = theme.HexARGB(ui.ly.Theme.Color.Hint)
								sat.Layout(ui.ly.Context)
							}),
						)
					})
				}))
		})
	}
}
