package pages

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/stalker-loki/pod/cmd/gui/model"
	"github.com/stalker-loki/pod/pkg/gui/gel"

	"github.com/stalker-loki/pod/cmd/gui/component"
	"github.com/stalker-loki/pod/cmd/gui/rcd"
	"github.com/stalker-loki/pod/pkg/gui/gelook"
)

var (
	// itemValue = &gel.DuoUIcounter{
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
	// }
	transactionsPanelElement = gel.NewPanel()
)

func History(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "HISTORY",
		Command:       rc.GetDuoUItransactions(),
		Border:        4,
		BorderColor:   th.Colors["Light"],
		Header:        historyHeader(rc, gtx, th),
		HeaderPadding: 4,
		Body:          historyBody(rc, gtx, th),
		BodyBgColor:   th.Colors["Light"],
		BodyPadding:   0,
		Footer:        func() {},
		FooterPadding: 0,
	}
	return th.DuoUIpage(page)
}

func historyHeader(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.Flex{
			Spacing: layout.SpaceBetween,
		}.Layout(gtx,
			layout.Rigid(component.
				TransactionsFilter(rc, gtx, th)),
			layout.Rigid(func() {
				th.DuoUIcounter(rc.GetDuoUItransactions()).
					Layout(gtx, rc.History.PerPage, "TxNum per page: ",
						fmt.Sprint(rc.History.PerPage.Value))
			}),
			layout.Rigid(func() {
				th.DuoUIcounter(rc.GetDuoUItransactions()).
					Layout(gtx, rc.History.Page, "TxNum page: ",
						fmt.Sprint(rc.History.Page.Value))
			}),
		)
	}
}

func historyBody(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.UniformInset(unit.Dp(16)).Layout(gtx, func() {
			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(
					func() {
						transactionsPanel := th.DuoUIpanel()
						transactionsPanel.PanelObject = rc.History.Txs.Txs
						transactionsPanel.ScrollBar = th.ScrollBar(16)
						transactionsPanelElement.PanelObjectsNumber = len(rc.History.Txs.Txs)
						transactionsPanel.Layout(gtx, transactionsPanelElement,
							func(i int, in interface{}) {
								txs := in.([]model.DuoUItransactionExcerpt)
								t := txs[i]
								th.DuoUIline(gtx, 0, 0, 1, th.Colors["Hint"])()
								for t.Link.Clicked(gtx) {
									rc.ShowPage = fmt.Sprintf("TRANSACTION %s", t.TxID)
									rc.GetSingleTx(t.TxID)()
									component.SetPage(rc, txPage(rc, gtx, th, t.TxID))
								}
								width := gtx.Constraints.Width.Max
								button := th.DuoUIbutton("", "",
									"", "", "", "",
									"", "", 0, 0,
									0, 0, 0, 0,
									0, 0)
								button.InsideLayout(gtx, t.Link, func() {
									gtx.Constraints.Width.Min = width
									layout.Flex{
										Spacing: layout.SpaceBetween,
									}.Layout(gtx,
										layout.Rigid(component.TxsDetails(gtx,
											th, i, &t)),
										layout.Rigid(component.Label(gtx,
											th, th.Fonts["Mono"], 12,
											th.Colors["Secondary"],
											fmt.Sprintf("%0.8f", t.Amount))))
								})
							})
					}))
		})
	}
}
