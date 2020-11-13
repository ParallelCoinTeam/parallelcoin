package gui

import (
	l "gioui.org/layout"
)

func (wg *WalletGUI) HistoryPage() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		return wg.th.VFlex().Fn(gtx)
	}
}
