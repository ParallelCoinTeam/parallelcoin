package gui

import (
	"fmt"

	"golang.org/x/exp/shiny/materialdesign/icons"

	l "gioui.org/layout"

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
	// create the header
	header := p9.TextTableHeader{
		{Text: "Amount", Priority: 0},
		{Text: "Transaction ID", Priority: 4},
		{Text: "Address", Priority: 2},
		{Text: "Category", Priority: 1},
		{Text: "Confirmations", Priority: 3},
		{Text: "Time", Priority: 5},
		{Text: "Comment", Priority: 6},
		{Text: "Fee", Priority: 7},
		{Text: "BlockHash", Priority: 8},
		{Text: "BlockTime", Priority: 9},
		{Text: "Generated", Priority: 10},
		{Text: "Abandoned", Priority: 11},
		{Text: "Time Received", Priority: 12},
		{Text: "Trusted", Priority: 13},
		{Text: "Vout", Priority: 14},
		{Text: "Wallet Conflicts", Priority: 15},
		{Text: "Account", Priority: 16},
		{Text: "Other Account", Priority: 17},
		{Text: "Involves Watch Only", Priority: 18},
	}
	body := p9.TextTableBody{}
	for i := range wg.State.allTxs {
		body = append(body, p9.TextTableRow{
			fmt.Sprintf("%v", wg.State.allTxs[i].Amount),
			wg.State.allTxs[i].TxID,
			wg.State.allTxs[i].Address,
			wg.State.allTxs[i].Category,
			fmt.Sprintf("%v", wg.State.allTxs[i].Confirmations),
			fmt.Sprintf("%v", wg.State.allTxs[i].Time),
			wg.State.allTxs[i].Comment,
			fmt.Sprintf("%v", wg.State.allTxs[i].Fee),
			wg.State.allTxs[i].BlockHash,
			fmt.Sprintf("%v", wg.State.allTxs[i].BlockTime),
			fmt.Sprintf("%v", wg.State.allTxs[i].Generated),
			fmt.Sprintf("%v", wg.State.allTxs[i].Abandoned),
			fmt.Sprintf("%v", wg.State.allTxs[i].Time),
			fmt.Sprintf("%v", wg.State.allTxs[i].Trusted),
			fmt.Sprintf("%v", wg.State.allTxs[i].Vout),
			fmt.Sprintf("%v", wg.State.allTxs[i].WalletConflicts),
			wg.State.allTxs[i].Account,
			wg.State.allTxs[i].OtherAccount,
			fmt.Sprintf("%v", wg.State.allTxs[i].InvolvesWatchOnly),
		})
	}
	table := &p9.TextTable{
		Header: header,
		Body:   body,
		Inset:  0.25,
	}

	return wg.th.Fill("DocBg",
		table.Fn,
		// p9.EmptySpace(0, 0),
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
