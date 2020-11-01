package gui

import (
	"fmt"
	l "gioui.org/layout"
)

func (wg *WalletGUI) TransactionsPage() l.Widget {
	le := func(gtx l.Context, index int) l.Dimensions {
		return wg.singleTransaction(gtx, index)
	}
	return func(gtx l.Context) l.Dimensions {
		return wg.th.VFlex().
			Rigid(
				wg.Inset(0.25,
					wg.th.Flex().
						Rigid(
							wg.Inset(0.25, wg.Caption("Date:").Color("DocText").Fn).Fn,
						).
						Rigid(
							//wg.sendButton(wg.sendAddresses[index].AddressBookBtn, "AddressBook", func() {}),
							wg.Inset(0.25, wg.Caption("Type:").Color("DocText").Fn).Fn,
						).
						Flexed(1,
							//wg.sendButton(wg.sendAddresses[index].PasteClipboardBtn, "Paste", func() {}),
							wg.Inset(0.25, wg.Caption("Label:").Color("DocText").Fn).Fn,
						).
						Rigid(
							//wg.sendButton(wg.sendAddresses[index].ClearBtn, "Close", func() {}),
							wg.Inset(0.25, wg.Caption("Amount(DUO):").Color("DocText").Fn).Fn,
						).Fn,
				).Fn,
			).
			Flexed(1,
				wg.Inset(0.25, wg.Fill("DocBg", wg.Inset(0.25,
					wg.lists["transactions"].Vertical().Length(len(wg.txs)).ListElement(le).Fn,
				).Fn).Fn).Fn,
			).Fn(gtx)
	}
}

func (wg *WalletGUI) Transactions() {
	chainClient, err := wg.chainClient()
	if err != nil {
	}
	txs, err := chainClient.ListTransactions(nil)
	if err != nil {
	}
	wg.txs = txs
}

func (wg *WalletGUI) singleTransaction(gtx l.Context, i int) l.Dimensions {
	return wg.Inset(0.25,
		wg.Fill("DocBg",
			wg.Inset(0.25,
				wg.th.VFlex().
					Rigid(
						wg.Inset(0.25,
							wg.th.Flex().
								Rigid(
									wg.Inset(0.1, wg.Caption(fmt.Sprint(wg.txs[i].Time)).Color("DocText").Fn).Fn,
								).
								Rigid(
									wg.Inset(0.1, wg.Caption(fmt.Sprint(wg.txs[i].Category)).Color("DocText").Fn).Fn,
								).
								Flexed(1,
									wg.Inset(0.1, wg.Caption(fmt.Sprint(wg.txs[i].Comment)).Color("DocText").Fn).Fn,
								).
								Rigid(
									wg.Inset(0.1, wg.Caption(fmt.Sprint(wg.txs[i].Amount)).Color("DocText").Fn).Fn,
								).Fn,
						).Fn,
					).Fn,
			).Fn,
		).Fn,
	).Fn(gtx)
}

//func (wg *WalletGUI) ClearAddress(i int) {
//	wg.sendAddresses = remove(wg.sendAddresses, i)
//}
