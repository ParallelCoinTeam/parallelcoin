package gui

import (
	l "gioui.org/layout"
	"github.com/p9c/pod/pkg/gui/p9"
)

type Market struct {
	SubtractFee     *p9.Bool
	AllAvailableBtn *p9.Clickable
}

func (wg *WalletGUI) MarketPage() l.Widget {
	// le := func(gtx l.Context, index int) l.Dimensions {
	// 	return wg.singleSendAddress(gtx, index)
	// }
	return wg.th.VFlex().
		Flexed(1,
			// wg.Inset(0.25,
			func(gtx l.Context) l.Dimensions {
				return l.Dimensions{}
			},
			// ).Fn,
		).Fn
}
