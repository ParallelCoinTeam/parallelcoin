package gui

import (
	l "gioui.org/layout"
	
	"github.com/p9c/pod/pkg/gui/p9"
)

type SendAddress struct {
	AddressInput      *p9.Input
	LabelInput        *p9.Input
	AddressBookBtn    *p9.Clickable
	PasteClipboardBtn *p9.Clickable
	ClearBtn          *p9.Clickable
	AmountInput       *p9.Input
	// AmountInput       *counter.Counter
	SubtractFee     *p9.Bool
	AllAvailableBtn *p9.Clickable
}

func (wg *WalletGUI) SendPage() l.Widget {
	return wg.th.VFlex().
		AlignMiddle().
		SpaceSides().
		Rigid(
			wg.th.Flex().
				Flexed(0.5, p9.EmptyMaxWidth()).
				Rigid(
					wg.th.H1("send").Fn,
				).
				Flexed(0.5, p9.EmptyMaxWidth()).
				Fn,
		).
		Fn
}
