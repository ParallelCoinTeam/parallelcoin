package gui

import (
	"strconv"

	l "gioui.org/layout"
)

func (wg *WalletGUI) OverviewPage() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		return wg.th.Flex().
			Flexed(0.5,
				wg.panel("Balances", wg.balances()),
			).
			Flexed(0.5,
				wg.panel("Recent transactions", wg.th.Body1("transactions").Color("PanelText").Fn),
			).Fn(gtx)
	}
}

func (wg *WalletGUI) balances() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		return wg.Inset(0.25,
			wg.th.VFlex().
				Rigid(
					wg.row("Available:",
						strconv.FormatFloat(wg.State.balance, 'f', -1, 64)),
				).
				Rigid(
					wg.row("Pending:",
						strconv.FormatFloat(wg.State.balanceUnconfirmed, 'f', -1, 64)),
				).
				Rigid(
					wg.row("Total:",
						strconv.FormatFloat(wg.State.balanceUnconfirmed+wg.State.balance, 'f', 8, 64)),
				).Fn,
		).Fn(gtx)
	}
}

func (wg *WalletGUI) row(title, value string) l.Widget {
	return func(gtx l.Context) l.Dimensions {
		return wg.Inset(0.25,
			wg.th.Flex().
				SpaceBetween().
				Rigid(
					wg.Inset(0.0, wg.Inset(0.5, wg.Body1(title).Color("DocText").Fn).Fn).Fn,
				).Flexed(1,
				wg.Inset(0.0, wg.Inset(0.5, wg.Caption(value).Color("DocText").Fn).Fn).Fn,
			).Fn,
		).Fn(gtx)
	}
}

func (wg *WalletGUI) panel(title string, content l.Widget) l.Widget {
	return func(gtx l.Context) l.Dimensions {
		return wg.Inset(0.25,
			wg.th.VFlex().
				Rigid(
					wg.Inset(0.0, wg.Fill("DocText", wg.Inset(0.5, wg.H6(title).Color("DocBg").Fn).Fn).Fn).Fn,
				).Flexed(1,
				wg.Fill("DocBg",
					wg.Inset(0.25,
						content,
					).Fn,
				).Fn,
			).Fn,
		).Fn(gtx)
	}
}
