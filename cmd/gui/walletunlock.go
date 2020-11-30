package gui

import (
	l "gioui.org/layout"
)

func (wg *WalletGUI) WalletUnlockPage(gtx l.Context) l.Dimensions {
	return wg.th.Fill("PanelBg",
		wg.th.Inset(0.5,
			wg.th.H4("unlock wallet").Fn,
			// p9.EmptyMaxWidth(),
		).Fn,
	).Fn(gtx)
}

