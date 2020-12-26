package gui

import (
	l "gioui.org/layout"
	
	"github.com/p9c/pod/pkg/gui/p9"
)

func (wg *WalletGUI) ReceivePage() l.Widget {
	return wg.th.VFlex().
		AlignMiddle().
		SpaceSides().
		Rigid(
			wg.th.Flex().
				Flexed(0.5, p9.EmptyMaxWidth()).
				Rigid(
					wg.th.H1("receive").Fn,
				).
				Flexed(0.5, p9.EmptyMaxWidth()).
				Fn,
		).
		Fn
}
