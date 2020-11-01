package gui

import (
	"fmt"
	"strings"

	l "gioui.org/layout"
	"gioui.org/op"

	"github.com/p9c/pod/pkg/gui/p9"
)

func (wg *WalletGUI) OverviewPage() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		balanceColumn := wg.th.Column(p9.Rows{
			{Label: "Available:", W: wg.balanceWidget(wg.State.balance)},
			{Label: "Unconfirmed:", W: wg.balanceWidget(wg.State.balanceUnconfirmed)},
			{Label: "Total:", W: wg.balanceWidget(wg.State.balance + wg.State.balanceUnconfirmed)},
		}, "bariol bold", 1).List
		return wg.th.Responsive(*wg.App.Size, p9.Widgets{
			{
				Widget: wg.th.VFlex().
					Rigid(
						func(gtx l.Context) l.Dimensions {
							_, bc := balanceColumn(gtx)
							return wg.th.Inset(0.25,
								wg.th.Fill("DocBg",
									wg.th.SliceToWidget(
										append([]l.Widget{
											func(gtx l.Context) l.Dimensions {
												_, bc = balanceColumn(gtx)
												// // render the widgets onto a second context to get their dimensions
												// gtx1 := p9.CopyContextDimensions(gtx, gtx.Constraints.Max, l.Vertical)
												// dim := p9.GetDimension(gtx1, wg.th.SliceToWidget(bc, l.Vertical))
												// gtx.Constraints.Max.X = dim.Size.X
												// gtx.Constraints.Min.X = 0
												// gtx.Constraints.Max.X = dim.Size.Y
												// gtx.Constraints.Min.X = 0
												gtx.Constraints.Max.X = *wg.Size
												return wg.th.Fill("DocText",
													wg.th.Flex().
														Rigid(
															wg.th.Inset(0.5,
																wg.th.H6("Balances").
																	// Font("bariol bold").
																	Color("DocBg").
																	Fn,
															).Fn,
														).Fn,
												).Fn(gtx)
											},
										},
											bc...), l.Vertical),
								).Fn,
							).Fn(gtx)
						},
					).
					Rigid(
						wg.panel("Recent transactions", true,
							wg.th.Body1("transactions").Color("PanelText").Fn),
					).
					Fn,
			},
			{
				Size: 1280,
				Widget: wg.th.Flex().
					Rigid(
						func(gtx l.Context) l.Dimensions {
							_, bc := balanceColumn(gtx)
							return wg.th.Inset(0.25,
								wg.th.Fill("DocBg",
									wg.th.SliceToWidget(
										append([]l.Widget{
											func(gtx l.Context) l.Dimensions {
												// render the widgets onto a second context to get their dimensions
												gtx1 := p9.CopyContextDimensions(gtx, gtx.Constraints.Max, l.Vertical)
												dim := p9.GetDimension(gtx1, wg.th.SliceToWidget(bc, l.Vertical))
												gtx.Constraints.Max.X = dim.Size.X
												gtx.Constraints.Min.X = dim.Size.X
												return wg.th.Fill("DocText",
													wg.th.Flex().
														Flexed(1,
															wg.th.Inset(0.5,
																wg.th.H6("Balances").
																	// Font("bariol bold").
																	Color("DocBg").
																	Fn,
															).Fn,
														).Fn,
												).Fn(gtx)
											},
										},
											bc...), l.Vertical),
								).Fn,
							).Fn(gtx)
						},
					).
					Flexed(1,
						wg.panel("Recent transactions", true,
							wg.th.Body1("transactions").Color("PanelText").Fn),
					).Fn,
			},
		}).Fn(gtx)
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
	bal := leftPadTo(15, 15, fmt.Sprintf("%6.8f", balance))
	return wg.th.Inset(0.25,
		wg.th.Flex().AlignEnd().
			Rigid(wg.th.Body1(" ").Fn).
			Rigid(
				wg.th.Caption(bal).
					Font("go regular").
					Fn,
			).
			Fn,
	).Fn
}

func (wg *WalletGUI) panel(title string, fill bool, content l.Widget) l.Widget {
	return func(gtx l.Context) l.Dimensions {
		w := wg.Inset(0.25,
			wg.Fill("DocBg",
				wg.th.VFlex().
					Rigid(
						wg.Fill("DocText",
							wg.th.Flex().
								Rigid(
									wg.Inset(0.5,
										wg.H6(title).Color("DocBg").Fn,
									).Fn,
								).Fn,
						).Fn,
					).
					Rigid(
						wg.Fill("DocBg",
							wg.Inset(0.25,
								content,
							).Fn,
						).Fn,
					).Fn,
			).Fn,
		).Fn
		if !fill {
			// render the widgets onto a second context to get their dimensions
			gtx1 := p9.CopyContextDimensions(gtx, gtx.Constraints.Max, l.Vertical)
			// generate the dimensions for all the list elements
			child := op.Record(gtx1.Ops)
			d := w(gtx1)
			_ = child.Stop()
			gtx.Constraints.Max.X = d.Size.X
			gtx.Constraints.Max.Y = d.Size.Y
			gtx.Constraints.Min = gtx.Constraints.Max
			w = wg.Inset(0.25,
				wg.th.VFlex().
					Rigid(
						wg.Fill("DocText",
							wg.th.Flex().
								Flexed(1,
									wg.Inset(0.5,
										wg.H6(title).Color("DocBg").Fn,
									).Fn,
								).Fn,
						).Fn,
					).
					Rigid(
						wg.Fill("DocBg",
							wg.Inset(0.25,
								content,
							).Fn,
						).Fn,
					).Fn,
			).Fn
		}
		return w(gtx)
	}
}
