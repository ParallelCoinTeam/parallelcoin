package pages

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/cmd/gui/controller"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/cmd/gui/theme"
)

var (
	itemValue = &controller.DuoUIcounter{
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

func headerTransactions(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		layout.Flex{
			Spacing: layout.SpaceBetween,
		}.Layout(gtx,
			layout.Rigid(component.TransactionsFilter(rc, gtx, th)),
			layout.Rigid(func() {
				th.DuoUIcounter(rc.GetTransactions()).Layout(gtx, itemValue)
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
				layout.Rigid(component.TransactionsList(rc, gtx, th)))
		})
	}
}
