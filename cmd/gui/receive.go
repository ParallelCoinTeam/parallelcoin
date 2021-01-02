package gui

import (
	l "gioui.org/layout"
	
	"github.com/p9c/pod/pkg/gui"
)

func (wg *WalletGUI) ReceivePage() l.Widget {
	return wg.VFlex().
		AlignMiddle().
		SpaceSides().
		Rigid(
			wg.Flex().
				Flexed(0.5, gui.EmptyMaxWidth()).
				Rigid(
					wg.H1("receive").Fn,
				).
				Flexed(0.5, gui.EmptyMaxWidth()).
				Fn,
		).
		Fn
}
