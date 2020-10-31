package explorer

import (
	"fmt"
	"strings"

	l "gioui.org/layout"
	"gioui.org/op"

	"github.com/p9c/pod/pkg/gui/p9"
)

func (ex *Explorer) OverviewPage() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		return ex.th.Responsive(*ex.App.Size, p9.Widgets{
			{
				Widget: ex.th.VFlex().
					Rigid(
						ex.panel("Balances", true, ex.balances()),
					).
					Rigid(
						ex.panel("Recent transactions", true, ex.th.Body1("transactions").Color("PanelText").Fn),
					).Fn,
			},
			{
				Size: 1280,
				Widget: ex.th.Flex().
					Rigid(
						ex.panel("Balances", false, ex.balances()),
					).
					Flexed(1,
						ex.panel("Recent transactions", true, ex.th.Body1("transactions").Color("PanelText").Fn),
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

func (ex *Explorer) balances() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		balance := leftPadTo(15, 15, fmt.Sprintf("%6.8f", ex.State.balance))
		balanceUnconfirmed := leftPadTo(15, 15, fmt.Sprintf("%6.8f", ex.State.balanceUnconfirmed))
		balanceTotal := leftPadTo(15, 15, fmt.Sprintf("%6.8f", ex.State.balanceUnconfirmed+ex.State.balance))
		return ex.Inset(0.25,
			ex.th.Flex().
				Rigid(
					ex.th.VFlex().
						Rigid(
							ex.th.Inset(0.25,
								ex.th.Body1("Available:").Font("bariol bold").Fn,
							).Fn,
						).
						Rigid(
							ex.th.Inset(0.25,
								ex.th.Body1("Pending:").Font("bariol bold").Fn,
							).Fn,
						).
						Rigid(
							ex.th.Inset(0.25,
								ex.th.Body1("Total:").Font("bariol bold").Fn,
							).Fn,
						).
						Fn,
				).
				Rigid(
					ex.th.VFlex().
						Rigid(
							ex.th.Inset(0.25,
								ex.th.Flex().AlignBaseline().
									Rigid(ex.th.Body1(" ").Fn).
									Rigid(
										ex.th.Caption(balance).
											Font("go regular").
											Fn,
									).
									Fn,
							).Fn,
						).
						Rigid(
							ex.th.Inset(0.25,
								ex.th.Flex().AlignBaseline().
									Rigid(ex.th.Body1(" ").Fn).
									Rigid(
										ex.th.Caption(balanceUnconfirmed).
											Font("go regular").
											Fn,
									).Fn,
							).Fn,
						).
						Rigid(
							ex.th.Inset(0.25,
								ex.th.Flex().AlignBaseline().
									Rigid(ex.th.Body1(" ").Fn).
									Rigid(
										ex.th.Caption(balanceTotal).
											Font("go regular").
											Fn,
									).Fn,
							).Fn,
						).Fn,
				).Fn,
		).Fn(gtx)
	}
}

func (ex *Explorer) panel(title string, fill bool, content l.Widget) l.Widget {
	return func(gtx l.Context) l.Dimensions {
		w := ex.Inset(0.25,
			ex.th.VFlex().
				Rigid(
					ex.Fill("DocText",
						ex.th.Flex().
							Rigid(
								ex.Inset(0.5,
									ex.H6(title).Color("DocBg").Fn,
								).Fn,
							).Fn,
					).Fn,
				).
				Rigid(
					ex.Fill("DocBg",
						ex.Inset(0.25,
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
