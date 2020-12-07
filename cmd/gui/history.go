package gui

import (
	"fmt"

	"golang.org/x/exp/shiny/materialdesign/icons"

	l "gioui.org/layout"

	"github.com/p9c/pod/pkg/gui/p9"
)

func (wg *WalletGUI) HistoryPage() l.Widget {

	return func(gtx l.Context) l.Dimensions {
		return wg.th.VFlex().
			Rigid(
				wg.th.Inset(0.25,
					wg.th.Fill("DocBg",
						wg.th.Responsive(*wg.Size, p9.Widgets{
							{
								Widget: wg.th.VFlex().
									Flexed(1, wg.HistoryPageView()).
									Rigid(
										// 	wg.th.Fill("DocBg",
										wg.th.Flex().AlignMiddle().SpaceBetween().
											Flexed(0.5, p9.EmptyMaxWidth()).
											Rigid(wg.HistoryPageStatusFilter()).
											Flexed(0.5, p9.EmptyMaxWidth()).
											Fn,
										// 	).Fn,
									).
									// Rigid(
									// 	wg.th.Fill("DocBg",
									// 		wg.th.Flex().AlignMiddle().SpaceBetween().
									// 			Rigid(wg.HistoryPager()).
									// 			Rigid(wg.HistoryPagePerPageCount()).
									// 			Fn,
									// 	).Fn,
									// ).
									Fn,
							},
							{
								Size: 1280,
								Widget: wg.th.VFlex().
									Flexed(1, wg.HistoryPageView()).
									Rigid(
										// 	wg.th.Fill("DocBg",
										wg.th.Flex().AlignMiddle().SpaceBetween().
											// 			Rigid(wg.HistoryPager()).
											Flexed(0.5, p9.EmptyMaxWidth()).
											Rigid(wg.HistoryPageStatusFilter()).
											Flexed(0.5, p9.EmptyMaxWidth()).
											// 			Rigid(wg.HistoryPagePerPageCount()).
											Fn,
										// 	).Fn,
									).
									Fn,
							},
						}).Fn,
						// ).Fn,
					).Fn,
				).Fn,
			).Fn(gtx)
	}
}

func (wg *WalletGUI) HistoryPageView() l.Widget {
	// gen := wg.bools["showGenerate"].GetValue()
	// sent := wg.bools["showSent"].GetValue()
	// recv := wg.bools["showReceived"].GetValue()
	// imma := wg.bools["showImmature"].GetValue()
	// // current := wg.incdecs["transactionsPerPage"].GetCurrent()
	// cursor := 0 // wg.historyCurPage * current
	// // Debug(cursor, wg.historyCurPage, current, *wg.Size)
	// var out []btcjson.ListTransactionsResult
	// // for i := 0; i < wg.incdecs["transactionsPerPage"].GetCurrent(); i++ {
	// // Debugs(wg.State.AllTxs)
	// ws := wg.State.AllTxs
	// for ; cursor < len(ws)-1; cursor++ {
	// 	wsa := ws[cursor]
	// 	if wsa.Generated && gen ||
	// 		wsa.Category == "send" && sent ||
	// 		wsa.Category == "generate" && gen ||
	// 		wsa.Category == "immature" && imma ||
	// 		wsa.Category == "receive" && recv ||
	// 		wsa.Category == "unknown" {
	// 		out = append(out, wsa)
	// 		// break
	// 	}
	//
	// }
	// if cursor == len(wg.State.AllTxs)-1 {
	// 	break
	// }
	// }
	// Debugs(out)

	return wg.th.VFlex().Flexed(1, wg.historyTable.Fn).Fn
	// wg.th.Fill("DocBg",
	// p9.EmptySpace(0, 0),
	// ).Fn
}

func (wg *WalletGUI) HistoryPager() l.Widget {
	v := wg.incdecs["transactionsPerPage"].GetCurrent()
	vd := len(wg.State.AllTxs) / v
	vm := len(wg.State.AllTxs) % v
	if vm != 0 {
		vd++
	}
	return wg.th.Flex().AlignMiddle().
		Rigid(
			wg.th.IconButton(wg.clickables["txPageBack"].
				SetClick(func() {
					wg.historyCurPage--
					if wg.historyCurPage < 0 {
						wg.historyCurPage = 0
					}
					// if wg.historyCurPage > wg.hist
				})).
				Background("Transparent").
				Color("DocText").
				Scale(1).
				Icon(
					wg.th.Icon().
						Color("DocText").
						Scale(1).
						Src(&icons.NavigationArrowBack),
				).
				Fn,
		).
		Rigid(
			wg.th.Inset(0.25,
				wg.th.Caption(fmt.Sprintf("page %d/%d", wg.historyCurPage, vd)).Fn,
			).Fn,
		).
		Rigid(
			wg.th.IconButton(wg.clickables["txPageForward"].
				SetClick(func() {
					wg.historyCurPage++
					pLen := len(wg.State.AllTxs) / wg.incdecs["transactionsPerPage"].GetCurrent()
					if wg.historyCurPage > pLen {
						wg.historyCurPage = pLen
					}
				})).
				Background("Transparent").
				Color("DocText").
				Scale(1).
				Icon(
					wg.th.Icon().Color("DocText").
						Scale(1).
						Src(&icons.NavigationArrowForward),
				).
				Fn,
		).Fn
}

func (wg *WalletGUI) HistoryPagePerPageCount() l.Widget {
	return wg.th.Flex().AlignMiddle().
		Rigid(
			wg.incdecs["transactionsPerPage"].
				Color("DocText").Background("DocBg").Scale(p9.Scales["Caption"]).Fn,
		).
		Rigid(
			wg.th.Inset(0.25,
				wg.th.Caption("tx/page").Fn,
			).Fn,
		).Fn
}

func (wg *WalletGUI) HistoryPageStatusFilter() l.Widget {
	return wg.th.Flex().AlignMiddle().
		Rigid(
			wg.th.Inset(0.25,
				wg.th.Caption("show").Fn,
			).Fn,
		).
		Rigid(
			wg.th.Inset(0.25,
				func(gtx l.Context) l.Dimensions {
					return wg.th.CheckBox(wg.bools["showGenerate"]).
						TextColor("DocText").
						TextScale(1).
						Text("generate").
						IconScale(1).
						Fn(gtx)
				},
			).Fn,
		).
		Rigid(
			wg.th.Inset(0.25,
				func(gtx l.Context) l.Dimensions {
					return wg.th.CheckBox(wg.bools["showSent"]).
						TextColor("DocText").
						TextScale(1).
						Text("sent").
						IconScale(1).
						Fn(gtx)
				},
			).Fn,
		).
		Rigid(
			wg.th.Inset(0.25,
				func(gtx l.Context) l.Dimensions {
					return wg.th.CheckBox(wg.bools["showReceived"]).
						TextColor("DocText").
						TextScale(1).
						Text("received").
						IconScale(1).
						Fn(gtx)
				},
			).Fn,
		).
		Rigid(
			wg.th.Inset(0.25,
				func(gtx l.Context) l.Dimensions {
					return wg.th.CheckBox(wg.bools["showImmature"]).
						TextColor("DocText").
						TextScale(1).
						Text("immature").
						IconScale(1).
						Fn(gtx)
				},
			).Fn,
		).
		Fn
}
