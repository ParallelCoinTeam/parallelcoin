package gui

import (
	l "gioui.org/layout"
	
	"github.com/p9c/pod/pkg/gui"
)

func (wg *WalletGUI) ReceivePage() l.Widget {
	return wg.th.VFlex().
		AlignMiddle().
		SpaceSides().
		Rigid(
			wg.th.Flex().
				Flexed(0.5, gui.EmptyMaxWidth()).
				Rigid(
					wg.th.H1("receive").Fn,
				).
				Flexed(0.5, gui.EmptyMaxWidth()).
				Fn,
		).
		Fn
}
