package gui

import (
	l "gioui.org/layout"
)

func (wg *WalletGUI) ReceivePage() l.Widget {
	return wg.MainApp.Placeholder("receive")
}
