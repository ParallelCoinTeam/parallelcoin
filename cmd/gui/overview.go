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
		return wg.th.Responsive(*wg.App.Size, p9.Widgets{
			{
				Widget: wg.th.VFlex().
					Rigid(
						wg.panel("Balances", true, wg.balances()),
					).
					Rigid(
						wg.panel("Recent transactions", true, wg.th.Body1("transactions").Color("PanelText").Fn),
					).Fn,
			},
			{
				Size: 1280,
				Widget: wg.th.Flex().
					Rigid(
						wg.panel("Balances", false, wg.balances()),
					).
					Flexed(1,
						wg.panel("Recent transactions", true, wg.th.Body1("transactions").Color("PanelText").Fn),
					).Fn},
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

func (wg *WalletGUI) balances() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		balance := leftPadTo(15, 15, fmt.Sprintf("%6.8f", wg.State.balance))
		balanceUnconfirmed := leftPadTo(15, 15, fmt.Sprintf("%6.8f", wg.State.balanceUnconfirmed))
		balanceTotal := leftPadTo(15, 15, fmt.Sprintf("%6.8f", wg.State.balanceUnconfirmed+wg.State.balance))
		return wg.Inset(0.25,
			wg.th.Flex().
				Rigid(
					wg.th.VFlex().
						Rigid(
							wg.th.Inset(0.25,
								wg.th.Body1("Available:").Font("bariol bold").Fn,
							).Fn,
						).
						Rigid(
							wg.th.Inset(0.25,
								wg.th.Body1("Pending:").Font("bariol bold").Fn,
							).Fn,
						).
						Rigid(
							wg.th.Inset(0.25,
								wg.th.Body1("Total:").Font("bariol bold").Fn,
							).Fn,
						).
						Fn,
				).
				Rigid(
					wg.th.VFlex().
						Rigid(
							wg.th.Inset(0.25,
								wg.th.Flex().AlignBaseline().
									Rigid(wg.th.Body1(" ").Fn).
									Rigid(
										wg.th.Caption(balance).
											Font("go regular").
											Fn,
									).
									Fn,
							).Fn,
						).
						Rigid(
							wg.th.Inset(0.25,
								wg.th.Flex().AlignBaseline().
									Rigid(wg.th.Body1(" ").Fn).
									Rigid(
										wg.th.Caption(balanceUnconfirmed).
											Font("go regular").
											Fn,
									).Fn,
							).Fn,
						).
						Rigid(
							wg.th.Inset(0.25,
								wg.th.Flex().AlignBaseline().
									Rigid(wg.th.Body1(" ").Fn).
									Rigid(
										wg.th.Caption(balanceTotal).
											Font("go regular").
											Fn,
									).Fn,
							).Fn,
						).Fn,
				).Fn,
		).Fn(gtx)
	}
}

func (wg *WalletGUI) panel(title string, fill bool, content l.Widget) l.Widget {
	return func(gtx l.Context) l.Dimensions {
		w := wg.Inset(0.25,
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
		).Fn
		if !fill {
			// render the widgets onto a second context to get their dimensions
			gtx1 := p9.CopyContextDimensions(gtx, gtx.Constraints.Max, l.Horizontal)
			// generate the dimensions for all the list elements
			child := op.Record(gtx1.Ops)
			d := w(gtx1)
			_ = child.Stop()
			gtx.Constraints.Max.X = d.Size.X
			gtx.Constraints.Max.Y = d.Size.Y
			gtx.Constraints.Min = gtx.Constraints.Max
		}
		return w(gtx)
	}
}
