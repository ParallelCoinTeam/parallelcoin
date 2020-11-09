package gui

import (
	l "gioui.org/layout"
)



func (wg *WalletGUI) WalletPage(gtx l.Context) l.Dimensions {
	return wg.walletPage.Fn()(gtx)
}
