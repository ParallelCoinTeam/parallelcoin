package gui

import (
	"fmt"

	l "gioui.org/layout"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

func (wg *WalletGUI) HistoryPage() l.Widget {

	return func(gtx l.Context) l.Dimensions {
		return wg.th.VFlex().
			Rigid(
				wg.th.Inset(0.25,
					wg.th.Fill("PanelBg",
						wg.th.Responsive(*wg.Size, p9.Widgets{
							{
								Widget: wg.th.VFlex().
									Flexed(1, wg.HistoryPageView()).
									Rigid(
										wg.th.Fill("DocBg",
											wg.th.Flex().AlignMiddle().SpaceBetween().
												Flexed(0.5, p9.EmptyMaxWidth()).
												Rigid(wg.HistoryPageStatusFilter()).
												Flexed(0.5, p9.EmptyMaxWidth()).
												Fn,
										).Fn,
									).
									Rigid(
										wg.th.Fill("DocBg",
											wg.th.Flex().AlignMiddle().SpaceBetween().
												Rigid(wg.HistoryPager()).
												Rigid(wg.HistoryPagePerPageCount()).
												Fn,
										).Fn,
									).
									Fn,
							},
							{
								Size: 1280,
								Widget: wg.th.VFlex().
									Flexed(1, wg.HistoryPageView()).
									Rigid(
										wg.th.Fill("DocBg",
											wg.th.Flex().AlignMiddle().SpaceBetween().
												Rigid(wg.HistoryPager()).
												Rigid(wg.HistoryPageStatusFilter()).
												Rigid(wg.HistoryPagePerPageCount()).
												Fn,
										).Fn,
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
	gen := wg.bools["showGenerate"].GetValue()
	sent := wg.bools["showSent"].GetValue()
	recv := wg.bools["showReceived"].GetValue()
	imma := wg.bools["showImmature"].GetValue()
	current := wg.incdecs["transactionsPerPage"].GetCurrent()
	cursor := wg.historyCurPage * current
	var out []btcjson.ListTransactionsResult
	for i := 0; i < wg.incdecs["transactionsPerPage"].GetCurrent(); i++ {
		for ; cursor < len(wg.State.allTxs); cursor++ {
			if wg.State.allTxs[cursor].Generated && gen ||
				wg.State.allTxs[cursor].Category == "send" && sent ||
				wg.State.allTxs[cursor].Category == "generate" && gen ||
				wg.State.allTxs[cursor].Category == "immature" && imma ||
				wg.State.allTxs[cursor].Category == "receive" && recv ||
				wg.State.allTxs[cursor].Category == "unknown" {
				out = append(out, wg.State.allTxs[cursor])
				break
			}
		}
		if cursor == len(wg.State.allTxs)-1 {
			break
		}
	}
	Debugs(out)
	return wg.th.Fill("DocBg",

		p9.EmptySpace(0, 0),
	).Fn
}

func (wg *WalletGUI) HistoryPager() l.Widget {
	v := wg.incdecs["transactionsPerPage"].GetCurrent()
	vd := len(wg.State.allTxs) / v
	vm := len(wg.State.allTxs) % v
	if vm != 0 {
		vd++
	}
	return wg.th.Flex().AlignMiddle().
		Rigid(
			wg.th.IconButton(wg.clickables["txPageBack"]).
				Background("Transparent").
				Color("DocText").
				Scale(1).
				Icon(
					wg.th.Icon().Color("DocText").
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
			wg.th.IconButton(wg.clickables["txPageForward"]).
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
