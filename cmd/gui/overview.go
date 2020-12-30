package gui

import (
	"fmt"
	"strings"
	
	icons2 "golang.org/x/exp/shiny/materialdesign/icons"
	
	l "gioui.org/layout"
	
	"github.com/p9c/pod/pkg/gui/p9"
)

func (wg *WalletGUI) OverviewPage() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		return wg.th.Flex().
			Rigid(
				wg.th.Inset(0.25,
					wg.th.Inset(0.25,
						wg.th.VFlex().AlignMiddle().
							Rigid(
								wg.th.Fill("DocBg",
									wg.th.Inset(0.25,
										wg.th.H5("balances").Fn,
									).Fn,
								).Fn,
							).
							Rigid(
								wg.th.Fill("PanelBg",
									wg.th.Inset(0.25,
										wg.th.Flex().AlignEnd().
											Rigid(
												wg.th.Inset(0.5,
													wg.th.VFlex().AlignBaseline().
														Rigid(
															wg.th.Flex().AlignBaseline().
																Rigid(
																	wg.th.Body1("confirmed").Fn,
																).
																Rigid(
																	wg.th.H6(" ").Fn,
																).
																Fn,
														).
														Rigid(
															wg.th.Flex().AlignBaseline().
																Rigid(
																	wg.th.Body1("unconfirmed").Fn,
																).
																Rigid(
																	wg.th.H6(" ").Fn,
																).
																Fn,
														).
														Rigid(
															wg.th.Flex().AlignBaseline().
																Rigid(
																	wg.th.Body1("total").Fn,
																).
																Rigid(
																	wg.th.H6(" ").Fn,
																).
																Fn,
														).
														Fn,
												).Fn,
											).
											Rigid(
												wg.th.Inset(0.5,
													wg.th.VFlex().AlignBaseline().AlignEnd().
														Rigid(
															wg.th.Flex().AlignBaseline().
																Rigid(
																	wg.th.H6(" ").Fn,
																).
																Rigid(
																	wg.th.Body2(fmt.Sprintf("%6.8f",
																		wg.State.balance.Load()),
																	).Font("go regular").Fn,
																).Fn,
														).
														Rigid(
															wg.th.Flex().AlignBaseline().
																Rigid(
																	wg.th.H6(" ").Fn,
																).
																Rigid(
																	wg.th.Body2(fmt.Sprintf("%6.8f",
																		wg.State.balanceUnconfirmed.Load()),
																	).Font("go regular").Fn,
																).Fn,
														).
														Rigid(
															wg.th.Flex().AlignBaseline().
																Rigid(
																	wg.th.H6(" ").Fn,
																).
																Rigid(
																	wg.th.Body2(
																		fmt.Sprintf("%6.8f", wg.State.balance.Load()+wg.
																			State.balanceUnconfirmed.Load()),
																	).Font("go regular").Fn,
																).Fn,
														).
														Fn,
												).Fn,
											).Fn,
									).Fn,
								).Fn,
							).Fn,
					).Fn,
				).Fn,
			).
			Rigid(
				wg.th.Inset(0.25,
					wg.th.Fill("DocBg",
						wg.th.Inset(0.25,
							wg.th.VFlex().AlignMiddle().
								Rigid(wg.th.Inset(0.25,
									wg.th.H5("recent transactions").Fn).Fn).
								Flexed(1,
									p9.EmptyMaxWidth(),
								).
								Fn,
						).Fn,
					).Fn,
				).Fn,
			).
			Fn(gtx)
	}
}

// RecentTransactions generates a display showing recent transactions
//
// fields to use: Address, Amount, BlockIndex, BlockTime, Category, Confirmations, Generated
func (wg *WalletGUI) RecentTransactions() l.Widget {
	var out []l.Widget
	first := true
	// out = append(out)
	for x := 0; x < 10; x++ {
		i := x
		if len(wg.State.AllTimeStrings.Load().([]string)) <= i {
			break
		}
		txs := wg.State.AllTxs[i]
		times := wg.State.AllTimeStrings
		// spacer
		if !first {
			out = append(out,
				wg.th.Fill("DocBg",
					wg.th.Inset(0.25, p9.EmptyMaxWidth()).Fn,
				).Fn,
			)
		} else {
			first = false
		}
		out = append(out,
			wg.th.Fill("DocBg",
				wg.th.Body1(fmt.Sprintf("%-6.8f DUO", txs.Amount)).Color("PanelText").Fn,
			).Fn,
		)
		
		out = append(out,
			wg.th.Fill("DocBg",
				wg.th.Caption(txs.Address).
					Font("go regular").
					Color("PanelText").
					TextScale(0.66).Fn,
			).Fn,
		)
		
		out = append(out,
			wg.th.Fill("DocBg",
				wg.th.Caption(txs.TxID).
					Font("go regular").
					Color("PanelText").
					TextScale(0.5).Fn,
			).Fn,
		)
		out = append(out,
			func(gtx l.Context) l.Dimensions {
				return wg.th.Fill("DocBg",
					wg.th.Flex().AlignMiddle(). // SpaceBetween().
						Rigid(
							wg.th.Flex().AlignMiddle().
								Rigid(
									wg.th.Icon().Color("DocText").Scale(1).Src(&icons2.DeviceWidgets).Fn,
								).
								// Rigid(
								// 	wg.th.Caption(fmt.Sprint(*txs.BlockIndex)).Fn,
								// 	// wg.buttonIconText(txs.clickBlock,
								// 	// 	fmt.Sprint(*txs.BlockIndex),
								// 	// 	&icons2.DeviceWidgets,
								// 	// 	wg.blockPage(*txs.BlockIndex)),
								// ).
								Rigid(
									wg.th.Caption(fmt.Sprintf("%d ", *txs.BlockIndex)).Fn,
								).
								Fn,
						).
						Rigid(
							wg.th.Flex().AlignMiddle().
								Rigid(
									wg.th.Icon().Color("DocText").Scale(1).Src(&icons2.ActionCheckCircle).Fn,
								).
								Rigid(
									wg.th.Caption(fmt.Sprintf("%d ", txs.Confirmations)).Fn,
								).
								Fn,
						).
						Rigid(
							wg.th.Flex().AlignMiddle().
								Rigid(
									func(gtx l.Context) l.Dimensions {
										switch txs.Category {
										case "generate":
											return wg.th.Icon().Color("DocText").Scale(1).Src(&icons2.ActionStars).Fn(gtx)
										case "immature":
											return wg.th.Icon().Color("DocText").Scale(1).Src(&icons2.ImageTimeLapse).Fn(gtx)
										case "receive":
											return wg.th.Icon().Color("DocText").Scale(1).Src(&icons2.ActionPlayForWork).Fn(gtx)
										case "unknown":
											return wg.th.Icon().Color("DocText").Scale(1).Src(&icons2.AVNewReleases).Fn(gtx)
										}
										return l.Dimensions{}
									},
								).
								Rigid(
									wg.th.Caption(txs.Category+" ").Fn,
								).
								Fn,
						).
						Rigid(
							wg.th.Flex().AlignMiddle().
								Rigid(
									wg.th.Icon().Color("DocText").Scale(1).Src(&icons2.DeviceAccessTime).Fn,
								).
								Rigid(
									wg.th.Caption(
										times.Load().([]string)[i],
									).Color("DocText").Fn,
								).
								Fn,
						).
						Fn,
				).
					Fn(gtx)
			})
	}
	le := func(gtx l.Context, index int) l.Dimensions {
		return out[index](gtx)
	}
	return func(gtx l.Context) l.Dimensions {
		return wg.lists["recent"].
			Vertical().
			Length(len(out)).
			ListElement(le).
			Fn(gtx)
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
	return wg.th.Flex().AlignEnd().
		Rigid(wg.th.H6(" ").Fn).
		Rigid(
			wg.th.Caption(bal).
				Font("go regular").
				Fn,
		).
		Fn
}

//
// func (wg *WalletGUI) panel(title string, fill bool, content l.Widget) l.Widget {
// 	return func(gtx l.Context) l.Dimensions {
// 		w := wg.Inset(0.25,
// 			wg.Fill("DocBg",
// 				wg.th.VFlex().
// 					Rigid(
// 						wg.Fill("DocText",
// 							wg.th.Flex().
// 								Rigid(
// 									wg.Inset(0.5,
// 										wg.H6(title).Color("DocBg").Fn,
// 									).Fn,
// 								).Fn,
// 						).Fn,
// 					).
// 					Rigid(
// 						wg.Fill("DocBg",
// 							wg.Inset(0.25,
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
// 			d := w(gtx1)
// 			_ = child.Stop()
// 			gtx.Constraints.Max.X = d.Size.X
// 			gtx.Constraints.Max.Y = d.Size.Y
// 			gtx.Constraints.Min = gtx.Constraints.Max
// 			w = wg.Inset(0.25,
// 				wg.th.VFlex().
// 					Rigid(
// 						wg.Fill("DocText",
// 							wg.th.Flex().
// 								Flexed(1,
// 									wg.Inset(0.5,
// 										wg.H6(title).Color("DocBg").Fn,
// 									).Fn,
// 								).Fn,
// 						).Fn,
// 					).
// 					Rigid(
// 						wg.Fill("DocBg",
// 							wg.Inset(0.25,
// 								content,
// 							).Fn,
// 						).Fn,
// 					).Fn,
// 			).Fn
// 		}
// 		return w(gtx)
// 	}
// }
