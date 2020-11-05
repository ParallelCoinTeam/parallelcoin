package gui

import (
	"fmt"
	l "gioui.org/layout"
	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	icons2 "golang.org/x/exp/shiny/materialdesign/icons"
)

type tx struct {
	time       string
	data       btcjson.ListTransactionsResult
	clickTx    *p9.Clickable
	clickBlock *p9.Clickable
	list       *p9.List
}

func (wg *WalletGUI) TransactionsPage() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		return wg.th.VFlex().
			Rigid(
				wg.Inset(0.25,
					wg.th.Flex().
						Rigid(
							wg.Inset(0.1,
								wg.Inset(0.25, wg.Caption("Number of displayed transactions:").Color("DocText").Fn).Fn,
							).Fn,
						).
						Rigid(
							wg.buttonText(wg.clickables["transactions10"], "10", wg.Transactions),
						).
						Rigid(
							//wg.sendButton(wg.sendAddresses[index].PasteClipboardBtn, "Paste", func() {}),
							wg.buttonText(wg.clickables["transactions30"], "30", wg.Transactions),
						).
						Rigid(
							//wg.sendButton(wg.sendAddresses[index].ClearBtn, "Close", func() {}),
							wg.buttonText(wg.clickables["transactions50"], "50", wg.Transactions),
						).Fn,
				).Fn,
			).
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
					wg.lists["transactions"].Vertical().Length(len(wg.txs)).ListElement(wg.singleTransaction).Fn,
				).Fn).Fn).Fn,
			).Fn(gtx)
	}
}

func (wg *WalletGUI) Transactions() {
	walletClient, err := wg.walletClient()
	if err != nil {
	}
	txs, err := walletClient.ListTransactionsCount("default", 20)
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
					).Rigid(
					wg.th.Fill("DocBg",
						wg.th.Flex().AlignMiddle(). // SpaceBetween().
										Rigid(
								wg.th.Flex().AlignMiddle().
									Rigid(
										wg.Icon().Color("DocText").Scale(1).Src(&icons2.DeviceWidgets).Fn,
									).
									Rigid(
										wg.th.Caption(fmt.Sprintf("%d ", *wg.State.txs[i].data.BlockIndex)).Fn,
									).
									Fn,
							).
							Rigid(
								wg.th.Flex().AlignMiddle().
									Rigid(
										wg.Icon().Color("DocText").Scale(1).Src(&icons2.ActionCheckCircle).Fn,
									).
									Rigid(
										wg.th.Caption(fmt.Sprintf("%d ", wg.State.txs[i].data.Confirmations)).Fn,
									).
									Rigid(
										wg.Inset(0.1, wg.buttonText(wg.State.txs[i].clickTx, "details", wg.txPage(i))).Fn,
									).
									Fn,
							).
							Rigid(
								wg.th.Flex().AlignMiddle().
									Rigid(
										func(gtx l.Context) l.Dimensions {
											switch wg.State.txs[i].data.Category {
											case "generate":
												return wg.Icon().Color("DocText").Scale(1).Src(&icons2.ActionStars).Fn(gtx)
											case "immature":
												return wg.Icon().Color("DocText").Scale(1).Src(&icons2.ImageTimeLapse).Fn(gtx)
											case "receive":
												return wg.Icon().Color("DocText").Scale(1).Src(&icons2.ActionPlayForWork).Fn(gtx)
											case "unknown":
												return wg.Icon().Color("DocText").Scale(1).Src(&icons2.AVNewReleases).Fn(gtx)
											}
											return l.Dimensions{}
										},
									).
									Rigid(
										wg.th.Caption(wg.State.txs[i].data.Category+" ").Fn,
									).
									Fn,
							).
							Rigid(
								wg.th.Flex().AlignMiddle().
									Rigid(
										wg.Icon().Color("DocText").Scale(1).Src(&icons2.DeviceAccessTime).Fn,
									).
									Rigid(
										wg.th.Caption(
											wg.State.txs[i].time,
										).Color("DocText").Fn,
									).
									Fn,
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

func (wg *WalletGUI) txIitem(label, data string) l.Widget {
	if data != "" {
		return wg.Inset(0.25,
			wg.th.VFlex().
				Rigid(
					wg.Inset(0.0, wg.Fill("PanelBg", wg.Inset(0.2, wg.H6(label).Color("DocText").Fn).Fn).Fn).Fn,
				).
				Rigid(
					wg.Inset(0.0, wg.Fill("DocBg", wg.Inset(0.2, wg.Body1(data).Color("DocText").Font("go regular").Fn).Fn).Fn).Fn,
				).Fn,
		).Fn
	} else {
		return p9.EmptyMaxWidth()
	}
}

func (wg *WalletGUI) txPage(i int) func() {
	txLayout := []l.Widget{
		wg.txIitem("TxId:", wg.State.txs[i].data.TxID),
		wg.txIitem("Comment:", wg.State.txs[i].data.Comment),
		wg.txIitem("Category:", wg.State.txs[i].data.Category),
		wg.txIitem("Address:", wg.State.txs[i].data.Address),
		wg.txIitem("Generated:", fmt.Sprint(wg.State.txs[i].data.Generated)),
		wg.txIitem("BIP125Replaceable:", wg.State.txs[i].data.BIP125Replaceable),
		wg.txIitem("Block Hash:", wg.State.txs[i].data.BlockHash),
		wg.txIitem("Block Index:", fmt.Sprint(wg.State.txs[i].data.BlockIndex)),
		wg.txIitem("BlockTime:", fmt.Sprint(wg.State.txs[i].data.BlockTime)),
		wg.txIitem("Category:", wg.State.txs[i].data.Category),
		wg.txIitem("Confirmations:", fmt.Sprint(wg.State.txs[i].data.Confirmations)),
		wg.txIitem("Fee:", fmt.Sprint(wg.State.txs[i].data.Fee)),
		wg.txIitem("InvolvesWatchOnly:", fmt.Sprint(wg.State.txs[i].data.InvolvesWatchOnly)),
		wg.txIitem("Time:", fmt.Sprint(wg.State.txs[i].data.Time)),
		wg.txIitem("TimeReceived:", fmt.Sprint(wg.State.txs[i].data.TimeReceived)),
		wg.txIitem("Vout:", fmt.Sprint(wg.State.txs[i].data.Vout)),
		wg.txIitem("WalletConflicts:", fmt.Sprint(wg.State.txs[i].data.WalletConflicts)),
		wg.txIitem("Comment:", wg.State.txs[i].data.Comment),
		wg.txIitem("OtherAccount:", wg.State.txs[i].data.OtherAccount),
	}
	le := func(gtx l.Context, index int) l.Dimensions {
		return txLayout[index](gtx)
	}

	return func() {
		wg.newWindow(wg.State.txs[i].data.TxID, "Tx: "+wg.State.txs[i].data.TxID, 600, 800,
			wg.th.VFlex().
				Rigid(
					wg.Inset(0.0, wg.Fill("Primary", wg.Inset(0.5, wg.Caption(wg.State.txs[i].data.TxID).Color("DocBg").Fn).Fn).Fn).Fn,
				).Flexed(1,
				wg.Inset(0,
					func(gtx l.Context) l.Dimensions {
						return wg.State.txs[i].list.Vertical().Length(len(txLayout)).ListElement(le).Fn(gtx)
					},
				).Fn,
			).
				Rigid(
					wg.Button(
						wg.State.txs[i].clickTx.SetClick(func() {
							wg.w[wg.State.txs[i].data.TxID].Window.Close()
						})).
						CornerRadius(0).
						Background("Primary").
						Color("Dark").
						Font("bariol bold").
						TextScale(1).
						Text("CLOSE").
						Inset(0.5).
						Fn,
				).Fn,
		)
	}
}
