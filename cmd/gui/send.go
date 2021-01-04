package gui

import (
	l "gioui.org/layout"
	
	"github.com/p9c/pod/pkg/gui"
)

type SendAddress struct {
	AddressInput      *gui.Input
	LabelInput        *gui.Input
	AddressBookBtn    *gui.Clickable
	PasteClipboardBtn *gui.Clickable
	ClearBtn          *gui.Clickable
	AmountInput       *gui.Input
	// AmountInput       *counter.Counter
	SubtractFee     *gui.Bool
	AllAvailableBtn *gui.Clickable
}

func (wg *WalletGUI) SendPage() l.Widget {
	return wg.MainApp.Placeholder("send")
}
