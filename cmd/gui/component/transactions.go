package component

import (
	"fmt"
	"gioui.org/layout"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
)

func TransactionsList(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		rc.History.TransList.Layout(gtx, len(rc.History.Txs.Txs), func(i int) {
			t := rc.History.Txs.Txs[i]
			HorizontalLine(gtx, 1, th.Colors["Hint"])()
			layout.Flex{
				Spacing: layout.SpaceBetween,
			}.Layout(gtx,
				layout.Rigid(txsDetails(gtx, th, i, &t)),
				layout.Rigid(Label(gtx, th, th.Fonts["Mono"], 12, th.Colors["Dark"], fmt.Sprintf("%0.8f", t.Amount))))
		})
	}
}

func txsDetails(gtx *layout.Context, th *gelook.DuoUItheme, i int, t *model.DuoUItransactionExcerpt) func() {
	return func() {
		layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprint(i))),
			layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], t.TxID)),
			layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprintf("%0.8f", t.Amount))),
			layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], t.Category)),
			layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], t.Time)),
		)
	}
}

func TransactionsFilter(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.Flex{}.Layout(gtx,
			layout.Rigid(txsFilterItem(gtx, th, "ALL", rc.History.Categories.AllTxs)),
			layout.Rigid(txsFilterItem(gtx, th, "MINTED", rc.History.Categories.MintedTxs)),
			layout.Rigid(txsFilterItem(gtx, th, "IMATURE", rc.History.Categories.ImmatureTxs)),
			layout.Rigid(txsFilterItem(gtx, th, "SENT", rc.History.Categories.SentTxs)),
			layout.Rigid(txsFilterItem(gtx, th, "RECEIVED", rc.History.Categories.ReceivedTxs)))
		switch c := true; c {
		case rc.History.Categories.AllTxs.Checked(gtx):
			rc.History.Category = "all"
		case rc.History.Categories.MintedTxs.Checked(gtx):
			rc.History.Category = "generate"
		case rc.History.Categories.ImmatureTxs.Checked(gtx):
			rc.History.Category = "immature"
		case rc.History.Categories.SentTxs.Checked(gtx):
			rc.History.Category = "sent"
		case rc.History.Categories.ReceivedTxs.Checked(gtx):
			rc.History.Category = "received"
		}
	}
}

func txsFilterItem(gtx *layout.Context, th *gelook.DuoUItheme, id string, c *gel.CheckBox) func() {
	return func() {
		th.DuoUIcheckBox(id, th.Colors["Light"], th.Colors["Light"]).Layout(gtx, c)
	}
}
