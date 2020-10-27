package gui

import (
	l "gioui.org/layout"
)

func (wg *WalletGUI) OverviewPage() l.Widget {
	return wg.th.Flex().
		Flexed(0.5,
			wg.panel("Balance", wg.balances()),
		).
		Flexed(0.5,
			wg.panel("Recent transactions", wg.th.Body1("transactions").Color("PanelText").Fn),
		).Fn
}

func (wg *WalletGUI) panel(title string, content l.Widget) l.Widget {
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
	).Fn
}

func (wg *WalletGUI) balances() l.Widget {
	return wg.Inset(0.25,
		wg.th.VFlex().
			Rigid(
				wg.Inset(0.0, wg.Fill("DocText", wg.Inset(0.5, wg.H6("title").Color("DocBg").Fn).Fn).Fn).Fn,
			).Fn,
	).Fn
}
