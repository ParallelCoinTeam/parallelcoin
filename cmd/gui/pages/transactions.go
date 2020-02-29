package pages

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/component"
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

func History(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) *theme.DuoUIpage {
	return th.DuoUIpage("HISTORY", 0, func() {}, component.ContentHeader(gtx, th, headerTransactions(rc, gtx, th)), txsBody(rc, gtx, th), func() {})
}
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
						component.HorizontalLine(gtx, 1, th.Color.Hint)()
						layout.Flex{
							Spacing: layout.SpaceBetween,
						}.Layout(gtx,
							layout.Rigid(txsDetails(gtx, th, i, &t)),
							layout.Rigid(component.Label(gtx, th, th.Font.Mono, fmt.Sprintf("%0.8f", t.Amount))))
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
			layout.Rigid(component.Label(gtx, th, th.Font.Primary, fmt.Sprint(i))),
			layout.Rigid(component.Label(gtx, th, th.Font.Primary, t.TxID)),
			layout.Rigid(component.Label(gtx, th, th.Font.Primary, fmt.Sprintf("%0.8f", t.Amount))),
			layout.Rigid(component.Label(gtx, th, th.Font.Primary, t.Category)),
			layout.Rigid(component.Label(gtx, th, th.Font.Primary, t.Time)),
		)
	}
}
