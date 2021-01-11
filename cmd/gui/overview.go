package gui

import (
	"fmt"
	"strings"
	"time"
	
	icons2 "golang.org/x/exp/shiny/materialdesign/icons"
	
	l "gioui.org/layout"
	
	"github.com/p9c/pod/pkg/gui"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

func (wg *WalletGUI) balanceCard(corners int) func(gtx l.Context) l.Dimensions {
	return wg.VFlex().AlignMiddle().
		Rigid(
			// wg.ButtonInset(0.25,
			wg.H5("balances").Fn,
			// ).Fn,
		).
		Rigid(
			wg.Fill("Primary", l.Center, wg.TextSize.V, corners,
				// wg.Flex().Flexed(1,
					wg.Flex().SpaceEvenly().
						Rigid(
							wg.Inset(0.25,
								wg.VFlex().AlignBaseline().
									Rigid(
										wg.Flex().AlignBaseline().
											Rigid(
												wg.Body1("confirmed").Color("Light").Fn,
											).
											Rigid(
												wg.H6(" ").Fn,
											).
											Fn,
									).
									Rigid(
										wg.Flex().AlignBaseline().
											Rigid(
												wg.Body1("unconfirmed").Color("Light").Fn,
											).
											Rigid(
												wg.H6(" ").Fn,
											).
											Fn,
									).
									Rigid(
										wg.Flex().AlignBaseline().
											Rigid(
												wg.Body1("total").Color("Light").Fn,
											).
											Rigid(
												wg.H6(" ").Fn,
											).
											Fn,
									).
									Fn,
							).Fn,
						).
						Rigid(
							wg.Inset(0.25,
								wg.VFlex().AlignBaseline().AlignEnd().
									Rigid(
										wg.Flex().AlignBaseline().
											Rigid(
												wg.H6(" ").Fn,
											).
											Rigid(
												wg.Caption(leftPadTo(14, 14,
													fmt.Sprintf("%6.8f",
														wg.State.balance.Load())),
												).Color("Light").Font("go regular").Fn,
											).Fn,
									).
									Rigid(
										wg.Flex().AlignBaseline().
											Rigid(
												wg.H6(" ").Fn,
											).
											Rigid(
												wg.Caption(leftPadTo(14, 14,
													fmt.Sprintf("%6.8f",
														wg.State.balanceUnconfirmed.Load())),
												).Color("Light").Font("go regular").Fn,
											).Fn,
									).
									Rigid(
										wg.Flex().AlignBaseline().
											Rigid(
												wg.H6(" ").Fn,
											).
											Rigid(
												wg.Caption(
													leftPadTo(14, 14, fmt.Sprintf("%6.8f", wg.State.balance.Load()+wg.
														State.balanceUnconfirmed.Load())),
												).Color("Light").Font("go regular").Fn,
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
	if wg.RecentTransactionsWidget == nil {
		wg.RecentTransactionsWidget = func(gtx l.Context) l.Dimensions {
			return l.Dimensions{Size: gtx.Constraints.Max}
		}
	}
	return func(gtx l.Context) l.Dimensions {
		return wg.Responsive(*wg.Size, gui.Widgets{
			{
				Size: 0,
				Widget:
				wg.VFlex().SpaceAround().AlignStart().
					Rigid(
						// wg.ButtonInset(0.25,
						wg.VFlex().SpaceSides().
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
								wg.Inset(0.25,
									wg.H5("recent transactions").Fn).Fn,
							).
							Flexed(1,
								wg.Fill("DocBg", l.Center, wg.TextSize.V, 0, wg.Inset(0.5,
									wg.RecentTransactionsWidget,
									// p9.EmptyMaxWidth(),
								).Fn).Fn,
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
								wg.Inset(0.25,
									wg.H5("recent transactions").Fn,
								).Fn,
							).
							Flexed(1,
								wg.Fill("DocBg", l.Center, wg.TextSize.V, 0, wg.Inset(0.25,
									wg.RecentTransactionsWidget,
									// p9.EmptyMaxWidth(),
								).Fn).Fn,
							).
							Fn,
						// ).
						// Fn,
					).
					Fn,
			},
		}).Fn(gtx)
	}
}

// RecentTransactions generates a display showing recent transactions
//
// fields to use: Address, Amount, BlockIndex, BlockTime, Category, Confirmations, Generated
func (wg *WalletGUI) RecentTransactions(n int, listName string) l.Widget {
	wg.txMx.Lock()
	defer wg.txMx.Unlock()
	var out []l.Widget
	first := true
	// out = append(out)
	var wga []btcjson.ListTransactionsResult
	switch listName {
	case "history":
		wga = wg.txHistoryList
	case "recent":
		wga = wg.txRecentList
	}
	if len(wga) == 0 {
		return func(gtx l.Context) l.Dimensions {
			return l.Dimensions{Size: gtx.Constraints.Max}
		}
	}
	for x := range wga {
		if x > n && n > 0 {
			break
		}
		
		i := x
		txs := wga[i]
		// spacer
		if !first {
			out = append(out,
				wg.Inset(0.25, gui.EmptyMaxWidth()).Fn,
			)
		} else {
			first = false
		}
		out = append(out,
			wg.Body1(fmt.Sprintf("%-6.8f DUO", txs.Amount)).Color("PanelText").Fn,
		)
		
		out = append(out,
			wg.Caption(txs.Address).
				Font("go regular").
				Color("PanelText").
				TextScale(0.66).Fn,
		)
		
		// out = append(out,
		// 	wg.Caption(txs.TxID).
		// 		Font("go regular").
		// 		Color("PanelText").
		// 		TextScale(0.5).Fn,
		// )
		out = append(out,
			wg.Flex().
				Rigid(
					wg.Flex().
						Rigid(
							wg.Icon().Color("DocText").Scale(1).Src(&icons2.DeviceWidgets).Fn,
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
							wg.Icon().Color("DocText").Scale(1).Src(&icons2.ActionCheckCircle).Fn,
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
									return wg.Icon().Color("DocText").Scale(1).Src(&icons2.ActionStars).Fn(gtx)
								case "immature":
									return wg.Icon().Color("DocText").Scale(1).Src(&icons2.ImageTimeLapse).Fn(gtx)
								case "receive":
									return wg.Icon().Color("DocText").Scale(1).Src(&icons2.ActionPlayForWork).Fn(gtx)
								case "unknown":
									return wg.Icon().Color("DocText").Scale(1).Src(&icons2.AVNewReleases).Fn(gtx)
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
							wg.Icon().Color("DocText").Scale(1).Src(&icons2.DeviceAccessTime).Fn,
						).
						Rigid(
							wg.Caption(
								time.Unix(txs.Time,
									0).Format("02 Jan 06 15:04:05 MST"),
							).Color("DocText").Fn,
						).
						Fn,
				).Fn,
		)
	}
	le := func(gtx l.Context, index int) l.Dimensions {
		return out[index](gtx)
	}
	// if listName == "recent" {
	// 	wg.lists[listName].LeftSide(true)
	// }
	corners := 0 // gui.NW | gui.SW | gui.NE
	// if listName == "" {
	// 	corners = 0
	// }
	wo := func(gtx l.Context) l.Dimensions {
		// clip.UniformRRect(f32.Rectangle{
		// 	Max: f32.Pt(float32(gtx.Constraints.Max.X), float32(gtx.Constraints.Max.Y)),
		// }, wg.TextSize.V/4).Add(gtx.Ops)
		return wg.Fill("DocBg", l.Center, wg.TextSize.V/2, corners, wg.lists[listName].
			Vertical().
			Length(len(out)).
			ListElement(le).
			Fn).Fn(gtx)
	}
	switch listName {
	case "history":
		wg.HistoryWidget = wo
		if !wg.txReady.Load() {
			wg.txReady.Store(true)
		}
	case "recent":
		wg.RecentTransactionsWidget = wo
	}
	return func(gtx l.Context) l.Dimensions {
		return wo(gtx)
	}
}

func leftPadTo(length, limit int, txt string) string {
	if len(txt) > limit {
		return txt[limit-len(txt):]
	}
	pad := length - len(txt)
	return strings.Repeat(" ", pad) + txt
}

func (wg *WalletGUI) balanceWidget(balance float64) l.Widget {
	bal :=
	// leftPadTo(15, 15,
		fmt.Sprintf("%6.8f", balance)
	// )
	return wg.Flex().AlignEnd().
		Rigid(wg.H6(" ").Fn).
		Rigid(
			wg.Caption(bal).
				Font("go regular").
				Fn,
		).
		Fn
}

//
// func (wg *WalletGUI) panel(title string, fill bool, content l.Widget) l.Widget {
// 	return func(gtx l.Context) l.Dimensions {
// 		Window := wg.ButtonInset(0.25,
// 			wg.Fill("DocBg",
// 				wg.VFlex().
// 					Rigid(
// 						wg.Fill("DocText",
// 							wg.Flex().
// 								Rigid(
// 									wg.ButtonInset(0.5,
// 										wg.H6(title).Color("DocBg").Fn,
// 									).Fn,
// 								).Fn,
// 						).Fn,
// 					).
// 					Rigid(
// 						wg.Fill("DocBg",
// 							wg.ButtonInset(0.25,
// 								content,
// 							).Fn,
// 						).Fn,
// 					).Fn,
// 			).Fn,
// 		).Fn
// 		if !fill {
// 			// render the widgets onto a second context to get their dimensions
// 			gtx1 := p9.CopyContextDimensionsWithMaxAxis(gtx, gtx.Constraints.Max, l.Vertical)
// 			// generate the dimensions for all the list elements
// 			child := op.Record(gtx1.Ops)
// 			d := Window(gtx1)
// 			_ = child.Stop()
// 			gtx.Constraints.Max.X = d.Size.X
// 			gtx.Constraints.Max.Y = d.Size.Y
// 			gtx.Constraints.Min = gtx.Constraints.Max
// 			Window = wg.ButtonInset(0.25,
// 				wg.VFlex().
// 					Rigid(
// 						wg.Fill("DocText",
// 							wg.Flex().
// 								Flexed(1,
// 									wg.ButtonInset(0.5,
// 										wg.H6(title).Color("DocBg").Fn,
// 									).Fn,
// 								).Fn,
// 						).Fn,
// 					).
// 					Rigid(
// 						wg.Fill("DocBg",
// 							wg.ButtonInset(0.25,
// 								content,
// 							).Fn,
// 						).Fn,
// 					).Fn,
// 			).Fn
// 		}
// 		return Window(gtx)
// 	}
// }
