package gui

import (
	"fmt"
	"strings"
	"time"
	
	"gioui.org/text"
	icons2 "golang.org/x/exp/shiny/materialdesign/icons"
	
	l "gioui.org/layout"
	
	"github.com/p9c/pod/pkg/gui"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

func (wg *WalletGUI) balanceCard(corners int) func(gtx l.Context) l.Dimensions {
	return wg.VFlex().AlignMiddle().
		Rigid(
			// wg.ButtonInset(0.25,
			wg.H5("balances").Alignment(text.Middle).Fn,
			// ).Fn,
		).
		Rigid(
			wg.Fill(
				"Primary", l.W, wg.TextSize.V, corners,
				// wg.Flex().Flexed(1,
				wg.Flex().SpaceEvenly().
					Rigid(
						wg.Inset(
							0.25,
							wg.VFlex().AlignBaseline().
								Rigid(
									wg.Inset(
										0.25,
										wg.Flex().AlignBaseline().
											Rigid(
												wg.Body1("confirmed").Color("Light").Fn,
											).
											Rigid(
												wg.H6(" ").Fn,
											).
											Fn,
									).Fn,
								).
								Rigid(
									wg.Inset(
										0.25,
										
										wg.Flex().AlignBaseline().
											Rigid(
												wg.Body1("unconfirmed").Color("Light").Fn,
											).
											Rigid(
												wg.H6(" ").Fn,
											).
											Fn,
									).Fn,
								).
								Rigid(
									wg.Inset(
										0.5,
										wg.Flex().AlignBaseline().
											Rigid(
												wg.H6("total").Color("Light").Fn,
											).
											Rigid(
												wg.H6(" ").Fn,
											).
											Fn,
									).Fn,
								).
								Fn,
						).Fn,
					).
					Rigid(
						wg.Inset(
							0.25,
							wg.VFlex().AlignBaseline().AlignEnd().
								Rigid(
									wg.Inset(
										0.25,
										wg.Flex().AlignBaseline().
											Rigid(
												wg.H6(" ").Fn,
											).
											Rigid(
												wg.Caption(
													leftPadTo(
														14, 14,
														fmt.Sprintf(
															"%6.8f",
															wg.State.balance.Load(),
														),
													),
												).Color("Light").Font("go regular").Fn,
											).Fn,
									).Fn,
								).
								Rigid(
									wg.Inset(
										0.25,
										wg.Flex().AlignBaseline().
											Rigid(
												wg.H6(" ").Fn,
											).
											Rigid(
												wg.Caption(
													leftPadTo(
														14, 14,
														fmt.Sprintf(
															"%6.8f",
															wg.State.balanceUnconfirmed.Load(),
														),
													),
												).Color("Light").Font("go regular").Fn,
											).Fn,
									).Fn,
								).
								Rigid(
									wg.Inset(
										0.5,
										wg.Flex().AlignBaseline().
											Rigid(
												wg.H6(" ").Fn,
											).
											Rigid(
												wg.H6(
													leftPadTo(
														14, 14, fmt.Sprintf(
															"%6.8f", wg.State.balance.Load()+wg.
																State.balanceUnconfirmed.Load(),
														),
													),
												).Color("Light").Fn,
											).Fn,
									).Fn,
								).
								Fn,
						).Fn,
					).Fn,
				// ).Fn,
			).Fn,
		).Fn
}

func (wg *WalletGUI) OverviewPage() l.Widget {
	if wg.RecentTxsWidget == nil {
		wg.RecentTxsWidget = func(gtx l.Context) l.Dimensions {
			return l.Dimensions{Size: gtx.Constraints.Max}
		}
	}
	return func(gtx l.Context) l.Dimensions {
		return wg.Responsive(
			*wg.Size, gui.Widgets{
				{
					Size: 0,
					Widget:
					wg.VFlex().AlignMiddle().
						Rigid(
							// wg.ButtonInset(0.25,
							wg.VFlex().
								Rigid(
									// wg.Inset(0.25,
									wg.balanceCard(0),
									// ).Fn,
								).Fn,
							// ).Fn,
						).
						Rigid(
							// wg.Inset(0.25,
							wg.VFlex().AlignMiddle().
								Rigid(
									wg.Inset(
										0.25,
										wg.H5("recent transactions").Fn,
									).Fn,
								).
								Flexed(
									1,
									// wg.Inset(0.5,
									wg.RecentTxsWidget,
									// p9.EmptyMaxWidth(),
									// ).Fn,
								).
								Fn,
							// ).Fn,
						).
						Fn,
				},
				{
					Size: 64,
					Widget: wg.Flex().SpaceAround().AlignMiddle(). // SpaceSides().AlignMiddle().
						Rigid(
							// wg.ButtonInset(0.25,
							wg.VFlex().SpaceSides().AlignMiddle().
								Rigid(
									// wg.Inset(0.25,
									wg.balanceCard(0),
									// ).Fn,
								).Fn,
							// ).Fn,
						).
						Rigid(
							// wg.Inset(0.25,
							wg.VFlex().SpaceSides().AlignMiddle().
								Rigid(
									wg.Inset(
										0.25,
										wg.H5("recent transactions").Fn,
									).Fn,
								).
								Flexed(
									1,
									// wg.Fill("DocBg", l.W, wg.TextSize.V, 0, wg.Inset(0.25,
									wg.RecentTxsWidget,
									// p9.EmptyMaxWidth(),
									// ).Fn).Fn,
								).
								Fn,
							// ).
							// Fn,
						).
						Fn,
				},
			},
		).Fn(gtx)
	}
}

func (wg *WalletGUI) recentTxCardSummary(txs *btcjson.ListTransactionsResult, clickable *gui.Clickable) l.Widget {
	return wg.ButtonLayout(
		clickable.SetClick(
			func() {
				dbg.Ln("clicked tx")
				// dbg.S(txs)
				wg.openTxID.Store(txs.TxID)
			},
		),
	).Background("DocBg").Embed(
		wg.VFlex().Rigid(
			wg.Inset(
				0.25,
				wg.Flex().
					Rigid(
						wg.Body1(fmt.Sprintf("%-6.8f DUO", txs.Amount)).Color("PanelText").Fn,
					).
					Flexed(
						1,
						wg.Inset(
							0.25,
							wg.Caption(txs.Address).
								Font("go regular").
								Color("PanelText").
								TextScale(0.66).
								Alignment(text.End).
								Fn,
						).Fn,
					).Fn,
			).Fn,
		).Rigid(
			wg.Inset(
				0.25,
				wg.Flex().Flexed(
					1,
					wg.Flex().
						Rigid(
							wg.Flex().
								Rigid(
									wg.Icon().Color("PanelText").Scale(1).Src(&icons2.DeviceWidgets).Fn,
								).
								// Rigid(
								// 	wg.Caption(fmt.Sprint(*txs.BlockIndex)).Fn,
								// 	// wg.buttonIconText(txs.clickBlock,
								// 	// 	fmt.Sprint(*txs.BlockIndex),
								// 	// 	&icons2.DeviceWidgets,
								// 	// 	wg.blockPage(*txs.BlockIndex)),
								// ).
								Rigid(
									wg.Caption(fmt.Sprintf("%d ", txs.BlockIndex)).Fn,
								).
								Fn,
						).
						Rigid(
							wg.Flex().
								Rigid(
									wg.Icon().Color("PanelText").Scale(1).Src(&icons2.ActionCheckCircle).Fn,
								).
								Rigid(
									wg.Caption(fmt.Sprintf("%d ", txs.Confirmations)).Fn,
								).
								Fn,
						).
						Rigid(
							wg.Flex().
								Rigid(
									func(gtx l.Context) l.Dimensions {
										switch txs.Category {
										case "generate":
											return wg.Icon().Color("PanelText").Scale(1).Src(&icons2.ActionStars).Fn(gtx)
										case "immature":
											return wg.Icon().Color("PanelText").Scale(1).Src(&icons2.ImageTimeLapse).Fn(gtx)
										case "receive":
											return wg.Icon().Color("PanelText").Scale(1).Src(&icons2.ActionPlayForWork).Fn(gtx)
										case "unknown":
											return wg.Icon().Color("PanelText").Scale(1).Src(&icons2.AVNewReleases).Fn(gtx)
										}
										return l.Dimensions{}
									},
								).
								Rigid(
									wg.Caption(txs.Category+" ").Fn,
								).
								Fn,
						).
						Rigid(
							wg.Flex().
								Rigid(
									wg.Icon().Color("PanelText").Scale(1).Src(&icons2.DeviceAccessTime).Fn,
								).
								Rigid(
									wg.Caption(
										time.Unix(
											txs.Time,
											0,
										).Format("02 Jan 06 15:04:05 MST"),
									).Color("PanelText").Fn,
								).
								Fn,
						).Fn,
				).Fn,
			).Fn,
		).Fn,
	).Fn
}

func (wg *WalletGUI) recentTxCardDetail(txs *btcjson.ListTransactionsResult) l.Widget {
	return wg.VFlex().
		Rigid(
			wg.Body1("details").Color("PanelText").Fn,
		).
		Rigid(
			wg.Inset(
				0.25,
				wg.Flex().
					Rigid(
						wg.Body1(fmt.Sprintf("%-6.8f DUO", txs.Amount)).Color("PanelText").Fn,
					).
					Flexed(
						1,
						wg.Inset(
							0.25,
							wg.Caption(txs.Address).
								Font("go regular").
								Color("PanelText").
								TextScale(0.66).
								Alignment(text.End).
								Fn,
						).Fn,
					).Fn,
			).Fn,
		).Rigid(
		wg.Inset(
			0.25,
			wg.Flex().Flexed(
				1,
				wg.Flex().
					Rigid(
						wg.Flex().
							Rigid(
								wg.Icon().Color("PanelText").Scale(1).Src(&icons2.DeviceWidgets).Fn,
							).
							// Rigid(
							// 	wg.Caption(fmt.Sprint(*txs.BlockIndex)).Fn,
							// 	// wg.buttonIconText(txs.clickBlock,
							// 	// 	fmt.Sprint(*txs.BlockIndex),
							// 	// 	&icons2.DeviceWidgets,
							// 	// 	wg.blockPage(*txs.BlockIndex)),
							// ).
							Rigid(
								wg.Caption(fmt.Sprintf("%d ", txs.BlockIndex)).Fn,
							).
							Fn,
					).
					Rigid(
						wg.Flex().
							Rigid(
								wg.Icon().Color("PanelText").Scale(1).Src(&icons2.ActionCheckCircle).Fn,
							).
							Rigid(
								wg.Caption(fmt.Sprintf("%d ", txs.Confirmations)).Fn,
							).
							Fn,
					).
					Rigid(
						wg.Flex().
							Rigid(
								func(gtx l.Context) l.Dimensions {
									switch txs.Category {
									case "generate":
										return wg.Icon().Color("PanelText").Scale(1).Src(&icons2.ActionStars).Fn(gtx)
									case "immature":
										return wg.Icon().Color("PanelText").Scale(1).Src(&icons2.ImageTimeLapse).Fn(gtx)
									case "receive":
										return wg.Icon().Color("PanelText").Scale(1).Src(&icons2.ActionPlayForWork).Fn(gtx)
									case "unknown":
										return wg.Icon().Color("PanelText").Scale(1).Src(&icons2.AVNewReleases).Fn(gtx)
									}
									return l.Dimensions{}
								},
							).
							Rigid(
								wg.Caption(txs.Category+" ").Fn,
							).
							Fn,
					).
					Rigid(
						wg.Flex().
							Rigid(
								wg.Icon().Color("PanelText").Scale(1).Src(&icons2.DeviceAccessTime).Fn,
							).
							Rigid(
								wg.Caption(
									time.Unix(
										txs.Time,
										0,
									).Format("02 Jan 06 15:04:05 MST"),
								).Color("PanelText").Fn,
							).
							Fn,
					).Fn,
			).Fn,
		).Fn,
	).Fn
}

// RecentTransactions generates a display showing recent transactions
//
// fields to use: Address, Amount, BlockIndex, BlockTime, Category, Confirmations, Generated
func (wg *WalletGUI) RecentTransactions(n int, listName string) l.Widget {
	wg.txMx.Lock()
	defer wg.txMx.Unlock()
	// wg.ready.Store(false)
	var out []l.Widget
	first := true
	// out = append(out)
	var txList []btcjson.ListTransactionsResult
	var clickables []*gui.Clickable
	switch listName {
	case "history":
		txList = wg.txHistoryList
		clickables = wg.txHistoryClickables
	case "recent":
		txList = wg.txRecentList
		clickables = wg.recentTxsClickables
	}
	ltxl := len(txList)
	ltc := len(clickables)
	if ltxl > ltc {
		count := ltxl - ltc
		for ; count > 0; count-- {
			clickables = append(clickables, wg.Clickable())
		}
	}
	if len(clickables) == 0 {
		return func(gtx l.Context) l.Dimensions {
			return l.Dimensions{Size: gtx.Constraints.Max}
		}
	}
	dbg.Ln(">>>>>>>>>>>>>>>> iterating transactions", n, listName)
	for x := range txList {
		if x > n && n > 0 {
			break
		}
		
		txs := txList[x]
		// spacer
		if !first {
			out = append(
				out,
				wg.Inset(0.25, gui.EmptyMaxWidth()).Fn,
			)
		} else {
			first = false
		}
		ck := clickables[x]
		out = append(
			out,
			func(gtx l.Context) l.Dimensions {
				return gui.If(
					wg.openTxID.Load() == txs.TxID,
					wg.recentTxCardDetail(&txs),
					wg.recentTxCardSummary(
						&txs, ck,
					),
				)(gtx)
			},
		)
		// out = append(out,
		// 	wg.Caption(txs.TxID).
		// 		Font("go regular").
		// 		Color("PanelText").
		// 		TextScale(0.5).Fn,
		// )
		// out = append(
		// 	out,
		// 	wg.Fill(
		// 		"DocBg", l.W, 0, 0,
		//
		// 	).Fn,
		// )
	}
	le := func(gtx l.Context, index int) l.Dimensions {
		return out[index](gtx)
	}
	wo := func(gtx l.Context) l.Dimensions {
		return wg.lists[listName].
			Vertical().
			Length(len(out)).
			ListElement(le).
			Fn(gtx)
	}
	dbg.Ln(">>>>>>>>>>>>>>>> history widget completed", n, listName)
	switch listName {
	case "history":
		wg.TxHistoryWidget = wo
		if !wg.ready.Load() {
			wg.ready.Store(true)
		}
	case "recent":
		wg.RecentTxsWidget = wo
	}
	return func(gtx l.Context) l.Dimensions {
		return wo(gtx)
	}
}

func leftPadTo(length, limit int, txt string) string {
	if len(txt) > limit {
		return txt[:limit]
	}
	if len(txt) == limit {
		return txt
	}
	pad := length - len(txt)
	return strings.Repeat(" ", pad) + txt
}
