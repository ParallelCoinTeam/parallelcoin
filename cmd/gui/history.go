package gui

import (
	"fmt"
	"time"

	l "github.com/p9c/gio/layout"

	"github.com/p9c/gel"
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
						wg.txDetailEntry("Abandoned", fmt.Sprint(txs.Abandoned), "DocBg", false),
						wg.txDetailEntry("Account", fmt.Sprint(txs.Account), "DocBgDim", false),
						wg.txDetailEntry("Address", txs.Address, "DocBg", false),
						wg.txDetailEntry("Block Hash", txs.BlockHash, "DocBgDim", true),
						wg.txDetailEntry("Block Index", fmt.Sprint(txs.BlockIndex), "DocBg", false),
						wg.txDetailEntry("Block Time", fmt.Sprint(time.Unix(txs.BlockTime, 0)), "DocBgDim", false),
						wg.txDetailEntry("Category", txs.Category, "DocBg", false),
						wg.txDetailEntry("Confirmations", fmt.Sprint(txs.Confirmations), "DocBgDim", false),
						wg.txDetailEntry("Fee", fmt.Sprintf("%0.8f", txs.Fee), "DocBg", false),
						wg.txDetailEntry("Generated", fmt.Sprint(txs.Generated), "DocBgDim", false),
						wg.txDetailEntry("Involves Watch Only", fmt.Sprint(txs.InvolvesWatchOnly), "DocBg", false),
						wg.txDetailEntry("Time", fmt.Sprint(time.Unix(txs.Time, 0)), "DocBgDim", false),
						wg.txDetailEntry("Time Received", fmt.Sprint(time.Unix(txs.TimeReceived, 0)), "DocBg", false),
						wg.txDetailEntry("Trusted", fmt.Sprint(txs.Trusted), "DocBgDim", false),
						wg.txDetailEntry("TxID", txs.TxID, "DocBg", true),
						// todo: add WalletConflicts here
						wg.txDetailEntry("Comment", fmt.Sprintf("%0.8f", txs.Amount), "DocBgDim", false),
						wg.txDetailEntry("OtherAccount", fmt.Sprint(txs.OtherAccount), "DocBg", false),
					}
					le := func(gtx l.Context, index int) l.Dimensions {
						return out[index](gtx)
					}
					return wg.VFlex().AlignStart().
						Rigid(
							wg.recentTxCardSummaryButton(&txs, wg.clickables["txPageBack"], "Primary", true),
							// wg.H6(wg.openTxID.Load()).Fn,
						).
						Rigid(
							wg.lists["txdetail"].
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
					wg.Size.Load(), gel.Widgets{
						{
							Widget: wg.VFlex().
								Flexed(1, wg.HistoryPageView()).
								// Rigid(
								// 	// 	wg.Fill("DocBg",
								// 	wg.Flex().AlignMiddle().SpaceBetween().
								// 		Flexed(0.5, gel.EmptyMaxWidth()).
								// 		Rigid(wg.HistoryPageStatusFilter()).
								// 		Flexed(0.5, gel.EmptyMaxWidth()).
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
								// 		Flexed(0.5, gel.EmptyMaxWidth()).
								// 		Rigid(wg.HistoryPageStatusFilter()).
								// 		Flexed(0.5, gel.EmptyMaxWidth()).
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
			// wg.Fill("DocBg", l.Center, wg.TextSize.True, 0,
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
