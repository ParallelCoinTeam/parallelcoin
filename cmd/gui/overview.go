package gui

import (
	"fmt"
	"strings"
	"time"

	icons2 "golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/gio/text"

	l "github.com/p9c/gio/layout"

	"github.com/p9c/gel"
	"github.com/p9c/matrjoska/pkg/btcjson"
)

func (wg *WalletGUI) balanceCard() func(gtx l.Context) l.Dimensions {

	// gtx.Constraints.Min.X = int(wg.TextSize.True * sp.inputWidth)
	return func(gtx l.Context) l.Dimensions {
		gtx.Constraints.Min.X =
			int(wg.TextSize.V * 16)
		gtx.Constraints.Max.X =
			int(wg.TextSize.V * 16)
		return wg.VFlex().
			AlignStart().
			Rigid(
				wg.Inset(
					0.25,
					wg.H5("balances").
						Fn,
				).Fn,
			).
			Rigid(
				wg.Fill(
					"Primary", l.E, 0, 0,
					wg.Inset(
						0.25,
						wg.Flex().
							AlignEnd().
							Flexed(
								1,
								wg.VFlex().AlignEnd().
									Rigid(
										wg.ButtonLayout(wg.clickables["balanceConfirmed"]).SetClick(
											func() {
												go wg.WriteClipboard(
													fmt.Sprintf(
														"%6.8f",
														wg.State.balance.Load(),
													),
												)
											},
										).Background("Transparent").Embed(
											wg.Flex().AlignEnd().Flexed(
												1,
												wg.Inset(
													0.5,
													wg.Caption(
														"confirmed"+leftPadTo(
															14, 14,
															fmt.Sprintf(
																"%6.8f",
																wg.State.balance.Load(),
															),
														),
													).
														Font("go regular").
														Alignment(text.End).
														Color("DocText").Fn,
												).Fn,
											).Fn,
										).Fn,
									).
									Rigid(
										wg.ButtonLayout(wg.clickables["balanceUnconfirmed"]).SetClick(
											func() {
												go wg.WriteClipboard(
													fmt.Sprintf(
														"%6.8f",
														wg.State.balanceUnconfirmed.Load(),
													),
												)
											},
										).Background("Transparent").Embed(
											wg.Flex().AlignEnd().Flexed(
												1,
												wg.Inset(
													0.5,
													wg.Caption(
														"unconfirmed"+leftPadTo(
															14, 14,
															fmt.Sprintf(
																"%6.8f",
																wg.State.balanceUnconfirmed.Load(),
															),
														),
													).
														Font("go regular").
														Alignment(text.End).
														Color("DocText").Fn,

												).Fn,
											).Fn,
										).Fn,
									).
									Rigid(
										wg.ButtonLayout(wg.clickables["balanceTotal"]).SetClick(
											func() {
												go wg.WriteClipboard(
													fmt.Sprintf(
														"%6.8f",
														wg.State.balance.Load()+wg.State.balanceUnconfirmed.Load(),
													),
												)
											},
										).Background("Transparent").Embed(
											wg.Flex().AlignEnd().Flexed(
												1,
												wg.Inset(
													0.5,
													wg.H5(
														"total"+leftPadTo(
															14, 14, fmt.Sprintf(
																"%6.8f", wg.State.balance.Load()+wg.
																	State.balanceUnconfirmed.Load(),
															),
														),
													).
														Alignment(text.End).
														Color("DocText").Fn,
												).
													Fn,
											).Fn,
										).Fn,
									).Fn,
							).Fn,
					).Fn,
				).Fn,
			).Fn(gtx)
	}
}

func (wg *WalletGUI) OverviewPage() l.Widget {
	if wg.RecentTxsWidget == nil {
		wg.RecentTxsWidget = func(gtx l.Context) l.Dimensions {
			return l.Dimensions{Size: gtx.Constraints.Max}
		}
	}
	return func(gtx l.Context) l.Dimensions {
		return wg.Responsive(
			wg.Size.Load(), gel.Widgets{
				{
					Size: 0,
					Widget:
					wg.VFlex().AlignStart().
						Rigid(
							// wg.ButtonInset(0.25,
							wg.VFlex().
								Rigid(
									wg.Inset(
										0.25,
										wg.balanceCard(),
									).Fn,
								).Fn,
							// ).Fn,
						).
						// Rigid(wg.Inset(0.25, gel.EmptySpace(0, 0)).Fn).
						Flexed(
							1,
							wg.VFlex().AlignStart().
								Rigid(
									wg.Inset(
										0.25,
										wg.H5("Recent Transactions").Fn,
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
						).
						Fn,
				},
				{
					Size: 64,
					Widget: wg.Flex().AlignStart().
						Rigid(
							// wg.ButtonInset(0.25,
							wg.VFlex(). // SpaceSides().AlignStart().
								Rigid(
									wg.Inset(
										0.25,
										wg.balanceCard(),
									).Fn,
								).Fn,
							// ).Fn,
						).
						// Rigid(wg.Inset(0.25, gel.EmptySpace(0, 0)).Fn).
						Flexed(
							1,
							// wg.Inset(
							// 	0.25,
							wg.VFlex().AlignStart().
								Rigid(
									wg.Inset(
										0.25,
										wg.H5("Recent transactions").Fn,
									).Fn,
								).
								Flexed(
									1,
									// wg.Fill("DocBg", l.W, wg.TextSize.True, 0, wg.Inset(0.25,
									wg.RecentTxsWidget,
									// p9.EmptyMaxWidth(),
									// ).Fn).Fn,
								).
								Fn,
							// ).
							// 	Fn,
						).
						Fn,
				},
			},
		).Fn(gtx)
	}
}

func (wg *WalletGUI) recentTxCardStub(txs *btcjson.ListTransactionsResult) l.Widget {
	return wg.Inset(
		0.25,
		wg.Flex().
			// AlignBaseline().
			// AlignStart().
			// SpaceEvenly().
			SpaceBetween().
			// Flexed(
			// 	1,
			// 	wg.Inset(
			// 		0.25,
			// 		wg.Caption(txs.Address).
			// 			Font("go regular").
			// 			Color("PanelText").
			// 			TextScale(0.66).
			// 			Alignment(text.End).
			// 			Fn,
			// 	).Fn,
			// ).
			Rigid(
				wg.Caption(fmt.Sprintf("%-6.8f DUO", txs.Amount)).Font("go regular").Color("DocText").Fn,
			).
			Rigid(
				wg.Flex().
					Rigid(
						wg.Icon().Color("PanelText").Scale(1).Src(&icons2.DeviceWidgets).Fn,
					).
					Rigid(
						wg.Caption(fmt.Sprint(txs.BlockIndex)).Fn,
						// wg.buttonIconText(txs.clickBlock,
						// 	fmt.Sprint(*txs.BlockIndex),
						// 	&icons2.DeviceWidgets,
						// 	wg.blockPage(*txs.BlockIndex)),
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
			// Flexed(1, gel.EmptyMaxWidth()).
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
						).Font("go regular").
							// Alignment(text.End).
							Color("PanelText").Fn,
					).
					Fn,
			).
			Fn,
	).
		Fn
}

func (wg *WalletGUI) recentTxCardSummary(txs *btcjson.ListTransactionsResult) l.Widget {
	return wg.VFlex().AlignStart().SpaceBetween().
		Rigid(
			// wg.Inset(
			// 	0.25,
			wg.Flex().AlignStart().SpaceBetween().
				Rigid(
					wg.H6(fmt.Sprintf("%-.8f DUO", txs.Amount)).Alignment(text.Start).Color("PanelText").Fn,
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
			// ).Fn,
		).
		Rigid(
			// wg.Inset(
			// 	0.25,
			wg.Flex().
				Flexed(
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
			// ).Fn,
		).Fn
}

func (wg *WalletGUI) recentTxCardSummaryButton(
	txs *btcjson.ListTransactionsResult,
	clickable *gel.Clickable,
	bgColor string, back bool,
) l.Widget {
	return wg.ButtonLayout(
		clickable.SetClick(
			func() {
				D.Ln("clicked tx")
				// D.S(txs)
				curr := wg.openTxID.Load()
				if curr == txs.TxID {
					wg.prevOpenTxID.Store(wg.openTxID.Load())
					wg.openTxID.Store("")
					moveto := wg.originTxDetail
					if moveto == "" {
						moveto = wg.MainApp.ActivePageGet()
					}
					wg.MainApp.ActivePage(moveto)
				} else {
					if wg.MainApp.ActivePageGet() == "home" {
						wg.originTxDetail = "home"
						wg.MainApp.ActivePage("history")
					} else {
						wg.originTxDetail = "history"
					}
					wg.openTxID.Store(txs.TxID)
				}
			},
		),
	).
		Background(bgColor).
		Embed(
			gel.If(
				back,
				wg.Flex().
					Rigid(
						wg.Icon().Color("PanelText").Scale(4).Src(&icons2.NavigationArrowBack).Fn,
					).
					Rigid(
						wg.Inset(0.5, gel.EmptyMinWidth()).Fn,
					).
					Flexed(
						1,
						wg.Fill(
							"DocBg", l.Center, 0, 0, wg.Inset(
								0.5,
								wg.recentTxCardSummary(txs),
							).Fn,
						).Fn,
					).
					Fn,
				wg.Flex().
					Rigid(
						wg.Inset(0.5, gel.EmptyMaxHeight()).Fn,
					).
					Flexed(
						1,
						wg.Fill(
							"DocBg", l.Center, 0, 0, wg.Inset(
								0.5,
								wg.recentTxCardSummary(txs),
							).Fn,
						).Fn,
					).
					Fn,
			),
		).Fn
}

func (wg *WalletGUI) recentTxCardSummaryButtonGenerate(
	txs *btcjson.ListTransactionsResult,
	clickable *gel.Clickable,
	bgColor string, back bool,
) l.Widget {
	return wg.ButtonLayout(
		clickable.SetClick(
			func() {
				D.Ln("clicked tx")
				// D.S(txs)
				curr := wg.openTxID.Load()
				if curr == txs.TxID {
					wg.prevOpenTxID.Store(wg.openTxID.Load())
					wg.openTxID.Store("")
					moveto := wg.originTxDetail
					if moveto == "" {
						moveto = wg.MainApp.ActivePageGet()
					}
					wg.MainApp.ActivePage(moveto)
				} else {
					if wg.MainApp.ActivePageGet() == "home" {
						wg.originTxDetail = "home"
						wg.MainApp.ActivePage("history")
					} else {
						wg.originTxDetail = "history"
					}
					wg.openTxID.Store(txs.TxID)
				}
			},
		),
	).
		Background(bgColor).
		Embed(
			wg.Flex().AlignStart().
				Rigid(
					// wg.Fill(
					// 	"Primary", l.W, 0, 0, wg.Inset(
					// 		0.5,
					gel.If(
						back,
						wg.Flex().AlignStart().
							Rigid(
								wg.Icon().Color("PanelText").Scale(4).Src(&icons2.NavigationArrowBack).Fn,
							).
							Flexed(
								1,
								wg.recentTxCardSummary(txs),
							).
							Fn,
						wg.Flex().AlignStart().
							Flexed(
								1,
								wg.recentTxCardStub(txs),
							).
							Fn,
						// wg.Flex().
						// 	Rigid(
						// 		wg.Inset(0.5, gel.EmptyMaxHeight()).Fn,
						// 	).
						// 	Flexed(
						// 		1,
						// 		wg.Fill(
						// 			"DocBg", l.Center, 0, 0, wg.Inset(
						// 				0.5,
						// 				wg.recentTxCardSummary(txs),
						// 			).Fn,
						// 		).Fn,
						// 	).
						// 	Fn,
					),
				).Fn,
			// ).Fn,
			// ).Fn,
		).Fn
}

func (wg *WalletGUI) recentTxCardDetail(txs *btcjson.ListTransactionsResult, clickable *gel.Clickable) l.Widget {
	return wg.VFlex().
		Rigid(
			wg.Fill(
				"Primary", l.Center, wg.TextSize.V, 0,
				wg.recentTxCardSummaryButton(txs, clickable, "Primary", false),
			).Fn,
			// ).
			// Rigid(
			// 	wg.Fill(
			// 		"DocBg", l.Center, wg.TextSize.True, 0,
			// 		wg.Flex().
			// 			Flexed(
			// 				1,
			// 				wg.Inset(
			// 					0.25,
			// 					wg.VFlex().
			// 						Rigid(wg.Inset(0.25, gel.EmptySpace(0, 0)).Fn).
			// 						Rigid(
			// 							wg.H6("Transaction Details").
			// 								Color("PanelText").
			// 								Fn,
			// 						).
			// 						Rigid(
			// 							wg.Inset(
			// 								0.25,
			// 								wg.VFlex().
			// 									Rigid(
			// 										wg.txDetailEntry("Transaction ID", txs.TxID),
			// 									).
			// 									Rigid(
			// 										wg.txDetailEntry("Address", txs.Address),
			// 									).
			// 									Rigid(
			// 										wg.txDetailEntry("Amount", fmt.Sprintf("%0.8f", txs.Amount)),
			// 									).
			// 									Rigid(
			// 										wg.txDetailEntry("In Block", fmt.Sprint(txs.BlockIndex)),
			// 									).
			// 									Rigid(
			// 										wg.txDetailEntry("First Mined", fmt.Sprint(txs.BlockTime)),
			// 									).
			// 									Rigid(
			// 										wg.txDetailEntry("Category", txs.Category),
			// 									).
			// 									Rigid(
			// 										wg.txDetailEntry("Confirmations", fmt.Sprint(txs.Confirmations)),
			// 									).
			// 									Rigid(
			// 										wg.txDetailEntry("Fee", fmt.Sprintf("%0.8f", txs.Fee)),
			// 									).
			// 									Rigid(
			// 										wg.txDetailEntry("Confirmations", fmt.Sprint(txs.Confirmations)),
			// 									).
			// 									Rigid(
			// 										wg.txDetailEntry("Involves Watch Only", fmt.Sprint(txs.InvolvesWatchOnly)),
			// 									).
			// 									Rigid(
			// 										wg.txDetailEntry("Time", fmt.Sprint(txs.Time)),
			// 									).
			// 									Rigid(
			// 										wg.txDetailEntry("Time Received", fmt.Sprint(txs.TimeReceived)),
			// 									).
			// 									Rigid(
			// 										wg.txDetailEntry("Trusted", fmt.Sprint(txs.Trusted)),
			// 									).
			// 									Rigid(
			// 										wg.txDetailEntry("Abandoned", fmt.Sprint(txs.Abandoned)),
			// 									).
			// 									Rigid(
			// 										wg.txDetailEntry("BIP125 Replaceable", fmt.Sprint(txs.BIP125Replaceable)),
			// 									).
			// 									Fn,
			// 							).Fn,
			// 						).Fn,
			// 				).Fn,
			// 			).Fn,
			// 	).Fn,
		).Fn
}

func (wg *WalletGUI) txDetailEntry(name, detail string, bgColor string, small bool) l.Widget {
	content := wg.Body1
	if small {
		content = wg.Caption
	}
	return wg.Fill(
		bgColor, l.Center, wg.TextSize.V, 0,
		wg.Flex().AlignBaseline().
			Flexed(
				0.25,
				wg.Inset(
					0.25,
					wg.Body1(name).
						Color("PanelText").
						Font("bariol bold").
						Fn,
				).Fn,
			).
			Flexed(
				0.75,
				wg.Flex().SpaceStart().Rigid(
					wg.Inset(
						0.25,
						content(detail).Font("go regular").
							Color("PanelText").
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
	// first := true
	// out = append(out)
	var txList []btcjson.ListTransactionsResult
	var clickables []*gel.Clickable
	txList = wg.txHistoryList
	switch listName {
	case "history":
		clickables = wg.txHistoryClickables
	case "recent":
		// txList = wg.txRecentList
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
	D.Ln(">>>>>>>>>>>>>>>> iterating transactions", n, listName)
	var collected int
	for x := range txList {
		if collected >= n && n > 0 {
			break
		}
		txs := txList[x]
		switch listName {
		case "history":
			collected++
		case "recent":
			if txs.Category == "generate" || txs.Category == "immature" || txs.Amount < 0 && txs.Fee == 0 {
				continue
			} else {
				collected++
			}
		}
		// spacer
		// if !first {
		// 	out = append(
		// 		out,
		// 		wg.Inset(0.25, gel.EmptyMaxWidth()).Fn,
		// 	)
		// } else {
		// 	first = false
		// }

		ck := clickables[x]
		out = append(
			out,
			func(gtx l.Context) l.Dimensions {
				return gel.If(
					txs.Category == "immature",
					wg.recentTxCardSummaryButtonGenerate(&txs, ck, "DocBg", false),
					gel.If(
						txs.Category == "send",
						wg.recentTxCardSummaryButton(&txs, ck, "Danger", false),
						gel.If(
							txs.Category == "receive",
							wg.recentTxCardSummaryButton(&txs, ck, "Success", false),
							gel.If(
								txs.Category == "generate",
								wg.recentTxCardSummaryButtonGenerate(&txs, ck, "DocBg", false),
								gel.If(
									wg.prevOpenTxID.Load() == txs.TxID,
									wg.recentTxCardSummaryButton(&txs, ck, "Primary", false),
									wg.recentTxCardSummaryButton(&txs, ck, "DocBg", false),
								),
							),
						),
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
		return wg.Inset(
			0.25,
			out[index],
		).Fn(gtx)
	}
	wo := func(gtx l.Context) l.Dimensions {
		return wg.VFlex().AlignStart().
			Rigid(
				wg.lists[listName].
					Vertical().
					Length(len(out)).
					ListElement(le).
					Fn,
			).Fn(gtx)
	}
	D.Ln(">>>>>>>>>>>>>>>> history widget completed", n, listName)
	switch listName {
	case "history":
		wg.TxHistoryWidget = wo
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
