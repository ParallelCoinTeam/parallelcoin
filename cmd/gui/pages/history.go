package pages

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gelook"
)

var (
//itemValue = &gel.DuoUIcounter{
//	Value:        11,
//	OperateValue: 1,
//	From:         0,
//	To:           15,
//	CounterInput: &gel.Editor{
//		Alignment:  text.Middle,
//		SingleLine: true,
//	},
//	CounterIncrease: new(gel.Button),
//	//CounterDecrease: new(controller.Button),
//	CounterReset: new(gel.Button),
//}
)

func History(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) *gelook.DuoUIpage {
	return th.DuoUIpage("HISTORY", 0, rc.GetDuoUItransactions(), component.ContentHeader(gtx, th, headerTransactions(rc, gtx, th)), txsBody(rc, gtx, th), func() {})
}

func headerTransactions(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.Flex{
			Spacing: layout.SpaceBetween,
		}.Layout(gtx,
			layout.Rigid(component.TransactionsFilter(rc, gtx, th)),
			layout.Rigid(func() {
				th.DuoUIcounter(rc.GetDuoUItransactions()).Layout(gtx, rc.History.PerPage, "TxNum per page: ", fmt.Sprint(rc.History.PerPage.Value))
			}),
			layout.Rigid(func() {
				th.DuoUIcounter(rc.GetDuoUItransactions()).Layout(gtx, rc.History.Page, "TxNum page: ", fmt.Sprint(rc.History.Page.Value))
			}),
		)
	}
}

func txsBody(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.UniformInset(unit.Dp(16)).Layout(gtx, func() {
			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(component.TransactionsList(rc, gtx, th)))
		})
	}
}
