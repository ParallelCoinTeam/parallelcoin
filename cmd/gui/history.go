package gui

import (
	"fmt"
	l "gioui.org/layout"
	
	"github.com/p9c/pod/pkg/gui"
)

func (wg *WalletGUI) HistoryPage() l.Widget {
	if wg.TxHistoryWidget == nil {
		wg.TxHistoryWidget = func(gtx l.Context) l.Dimensions {
			return l.Dimensions{Size: gtx.Constraints.Max}
		}
	}
	return func(gtx l.Context) l.Dimensions {
		if wg.openTxID.Load() != "" {
			for i := range wg.txHistoryList {
				if wg.txHistoryList[i].TxID == wg.openTxID.Load() {
					txs := wg.txHistoryList[i]
					// instead return detail view
					var out []l.Widget
					out = []l.Widget{
						wg.txDetailEntry("Abandoned", fmt.Sprint(txs.Abandoned), "DocBg"),
						wg.txDetailEntry("Account", fmt.Sprint(txs.Account), "DocBgDim"),
						wg.txDetailEntry("Address", txs.Address, "DocBg"),
						wg.txDetailEntry("Block Hash", txs.BlockHash, "DocBgDim"),
						wg.txDetailEntry("Block Index", fmt.Sprint(txs.BlockIndex), "DocBg"),
						wg.txDetailEntry("Block Time", fmt.Sprint(txs.BlockTime), "DocBgDim"),
						wg.txDetailEntry("Category", txs.Category, "DocBg"),
						wg.txDetailEntry("Confirmations", fmt.Sprint(txs.Confirmations), "DocBgDim"),
						wg.txDetailEntry("Fee", fmt.Sprintf("%0.8f", txs.Fee), "DocBg"),
						wg.txDetailEntry("Generated", fmt.Sprint(txs.Generated), "DocBgDim"),
						wg.txDetailEntry("Involves Watch Only", fmt.Sprint(txs.InvolvesWatchOnly), "DocBg"),
						wg.txDetailEntry("Time", fmt.Sprint(txs.Time), "DocBgDim"),
						wg.txDetailEntry("Time Received", fmt.Sprint(txs.TimeReceived), "DocBg"),
						wg.txDetailEntry("Trusted", fmt.Sprint(txs.Trusted), "DocBgDim"),
						wg.txDetailEntry("TxID", txs.TxID, "DocBg"),
						// todo: add WalletConflicts here
						wg.txDetailEntry("Comment", fmt.Sprintf("%0.8f", txs.Amount), "DocBgDim"),
						wg.txDetailEntry("OtherAccount", fmt.Sprint(txs.BlockTime), "DocBg"),
					}
					le := func(gtx l.Context, index int) l.Dimensions {
						return out[index](gtx)
					}
					return wg.VFlex().
						Rigid(
							wg.recentTxCardSummaryButton(&txs, wg.clickables["txPageBack"], "Primary", true),
							// wg.H6(wg.openTxID.Load()).Fn,
						).
						Flexed(
							1, wg.lists["txdetail"].
								Vertical().
								Length(len(out)).
								ListElement(le).
								Fn,
						).
						Fn(gtx)
					
					// return wg.Flex().Flexed(
					// 	1,
					// 	wg.H3(wg.openTxID.Load()).Fn,
					// ).Fn(gtx)
				}
			}
			// if we got to here, the tx was not found
			if wg.originTxDetail != "" {
				wg.MainApp.ActivePage(wg.originTxDetail)
				wg.originTxDetail = ""
			}
		}
		return wg.VFlex().
			Rigid(
				// wg.Fill("DocBg", l.Center, 0, 0,
				// 	wg.Inset(0.25,
				wg.Responsive(
					*wg.Size, gui.Widgets{
						{
							Widget: wg.VFlex().
								Flexed(1, wg.HistoryPageView()).
								// Rigid(
								// 	// 	wg.Fill("DocBg",
								// 	wg.Flex().AlignMiddle().SpaceBetween().
								// 		Flexed(0.5, gui.EmptyMaxWidth()).
								// 		Rigid(wg.HistoryPageStatusFilter()).
								// 		Flexed(0.5, gui.EmptyMaxWidth()).
								// 		Fn,
								// 	// 	).Fn,
								// ).
								// Rigid(
								// 	wg.Fill("DocBg",
								// 		wg.Flex().AlignMiddle().SpaceBetween().
								// 			Rigid(wg.HistoryPager()).
								// 			Rigid(wg.HistoryPagePerPageCount()).
								// 			Fn,
								// 	).Fn,
								// ).
								Fn,
						},
						{
							Size: 64,
							Widget: wg.VFlex().
								Flexed(1, wg.HistoryPageView()).
								// Rigid(
								// 	// 	wg.Fill("DocBg",
								// 	wg.Flex().AlignMiddle().SpaceBetween().
								// 		// 			Rigid(wg.HistoryPager()).
								// 		Flexed(0.5, gui.EmptyMaxWidth()).
								// 		Rigid(wg.HistoryPageStatusFilter()).
								// 		Flexed(0.5, gui.EmptyMaxWidth()).
								// 		// 			Rigid(wg.HistoryPagePerPageCount()).
								// 		Fn,
								// 	// 	).Fn,
								// ).
								Fn,
						},
					},
				).Fn,
				// ).Fn,
				// ).Fn,
			).Fn(gtx)
	}
}

func (wg *WalletGUI) HistoryPageView() l.Widget {
	return wg.VFlex().
		Rigid(
			// wg.Fill("DocBg", l.Center, wg.TextSize.V, 0,
			// 	wg.Inset(0.25,
			wg.TxHistoryWidget,
			// ).Fn,
			// ).Fn,
		).Fn
}

func (wg *WalletGUI) HistoryPageStatusFilter() l.Widget {
	return wg.Flex().AlignMiddle().
		Rigid(
			wg.Inset(
				0.25,
				wg.Caption("show").Fn,
			).Fn,
		).
		Rigid(
			wg.Inset(
				0.25,
				func(gtx l.Context) l.Dimensions {
					return wg.CheckBox(wg.bools["showGenerate"]).
						TextColor("DocText").
						TextScale(1).
						Text("generate").
						IconScale(1).
						Fn(gtx)
				},
			).Fn,
		).
		Rigid(
			wg.Inset(
				0.25,
				func(gtx l.Context) l.Dimensions {
					return wg.CheckBox(wg.bools["showSent"]).
						TextColor("DocText").
						TextScale(1).
						Text("sent").
						IconScale(1).
						Fn(gtx)
				},
			).Fn,
		).
		Rigid(
			wg.Inset(
				0.25,
				func(gtx l.Context) l.Dimensions {
					return wg.CheckBox(wg.bools["showReceived"]).
						TextColor("DocText").
						TextScale(1).
						Text("received").
						IconScale(1).
						Fn(gtx)
				},
			).Fn,
		).
		Rigid(
			wg.Inset(
				0.25,
				func(gtx l.Context) l.Dimensions {
					return wg.CheckBox(wg.bools["showImmature"]).
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
