package duoui

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/cmd/gui/rcd"
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
		Value:           11,
		OperateValue:    1,
		From:            0,
		To:              15,
		CounterIncrease: new(controller.Button),
		CounterDecrease: new(controller.Button),
		CounterReset:    new(controller.Button),
	}
)

func txsFilter(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		layout.Flex{}.Layout(gtx,
			layout.Rigid(txsFilterItem(gtx, th, "ALL", allTxs)),
			layout.Rigid(txsFilterItem(gtx, th, "MINTED", mintedTxs)),
			layout.Rigid(txsFilterItem(gtx, th, "IMATURE", immatureTxs)),
			layout.Rigid(txsFilterItem(gtx, th, "SENT", sentTxs)),
			layout.Rigid(txsFilterItem(gtx, th, "RECEIVED", receivedTxs)))
	}
}

func txsFilterItem(gtx *layout.Context, th *theme.DuoUItheme, id string, c *controller.CheckBox) func() {
	return func() {
		th.DuoUIcheckBox(id, th.Color.Light, th.Color.Light).Layout(gtx, c)
	}
}

func headerTransactions(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		layout.Flex{
			Spacing: layout.SpaceBetween,
		}.Layout(gtx,
			layout.Rigid(txsFilter(rc, gtx, th)),
			layout.Rigid(func() {
				th.DuoUIcounter().Layout(gtx, itemValue)
			}),
		)
	}
}

func txsBody(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		layout.UniformInset(unit.Dp(16)).Layout(gtx, func() {
			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func() {
					transList.Layout(gtx, len(rc.Status.Wallet.Transactions.Txs), func(i int) {
						t := rc.Status.Wallet.Transactions.Txs[i]
						line(gtx, th.Color.Hint)()
						layout.Flex{
							Spacing: layout.SpaceBetween,
						}.Layout(gtx,
							layout.Rigid(txsDetails(gtx, th, i, &t)),
							layout.Rigid(func() {
								sat := th.Body1(fmt.Sprintf("%0.8f", t.Amount))
								sat.Font.Typeface = th.Font.Primary
								sat.Color = theme.HexARGB(th.Color.Hint)
								sat.Layout(gtx)
							}),
						)
					})
				}))
		})
	}
}

func txsDetails(gtx *layout.Context, th *theme.DuoUItheme, i int, t *model.DuoUItx) func() {
	return func() {
		layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func() {
				num := th.Body1(fmt.Sprint(i))
				num.Font.Typeface = th.Font.Primary
				num.Color = theme.HexARGB(th.Color.Hint)
				num.Layout(gtx)
			}),
			layout.Rigid(func() {
				tim := th.Body1(t.TxID)
				tim.Font.Typeface = th.Font.Primary
				tim.Color = theme.HexARGB(th.Color.Hint)
				tim.Layout(gtx)
			}),
			layout.Rigid(func() {
				amount := th.H5(fmt.Sprintf("%0.8f", t.Amount))
				amount.Font.Typeface = th.Font.Primary
				amount.Color = theme.HexARGB(th.Color.Hint)
				amount.Alignment = text.End
				amount.Font.Variant = "Mono"
				amount.Font.Weight = text.Bold
				amount.Layout(gtx)
			}),
			layout.Rigid(func() {
				sat := th.Body1(t.Category)
				sat.Font.Typeface = th.Font.Primary
				sat.Color = theme.HexARGB(th.Color.Hint)
				sat.Layout(gtx)
			}),
			layout.Rigid(func() {
				l := th.Body2(t.Time)
				l.Font.Typeface = th.Font.Primary
				l.Color = theme.HexARGB(th.Color.Hint)
				l.Layout(gtx)
			}),
		)
	}
}
