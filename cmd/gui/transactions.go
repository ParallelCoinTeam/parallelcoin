package gui

import (
	"fmt"
	l "gioui.org/layout"
)

func (wg *WalletGUI) TransactionsPage() l.Widget {
	fmt.Print("sss", wg.sendAddresses)
	le := func(gtx l.Context, index int) l.Dimensions {
		return wg.Caption("BalaaaaaaaaaaaaaaaO_" + fmt.Sprint(index)).Color("DocBg").Fn(gtx)
	}
	return wg.th.VFlex().
		Flexed(1,
			wg.Inset(0.0, wg.Fill("DocText", wg.Inset(0.5,
				wg.lists["transactions"].Vertical().Length(len(wg.sendAddresses)).ListElement(le).Fn,
			).Fn).Fn).Fn,
		).Fn
}
