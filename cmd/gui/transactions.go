package gui

import (
	"fmt"

	l "gioui.org/layout"
	icons2 "golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/pod/pkg/gui/f"
	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

func (wg *WalletGUI) OldTransactionsPage() l.Widget {
	// // TODO: this page doesn't have data being populated yet
	// if true {
	// 	return func(l.Context) l.Dimensions {
	// 		return l.Dimensions{}
	// 	}
	// }
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
							// wg.sendButton(wg.sendAddresses[index].PasteClipboardBtn, "Paste", func() {}),
							wg.buttonText(wg.clickables["transactions30"], "30", wg.Transactions),
						).
						Rigid(
							// wg.sendButton(wg.sendAddresses[index].ClearBtn, "Close", func() {}),
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
							wg.Inset(0.25, wg.Caption("Type:").Color("DocText").Fn).Fn,
						).
						Flexed(1,
							wg.Inset(0.25, wg.Caption("Label:").Color("DocText").Fn).Fn,
						).
						Rigid(
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
	// walletClient, err := wg.walletClient()
	// if err != nil {
	// }
	if wg.WalletClient == nil {
		Debug("not connected to wallet yet")
		return
	}
	var txs []btcjson.ListTransactionsResult
	var err error
	if txs, err = wg.WalletClient.ListTransactionsCount("default", 20); Check(err) {
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
										wg.th.Caption(fmt.Sprintf("%d ", *wg.txs[i].BlockIndex)).Fn,
									).
									Fn,
							).
							Rigid(
								wg.th.Flex().AlignMiddle().
									Rigid(
										wg.Icon().Color("DocText").Scale(1).Src(&icons2.ActionCheckCircle).Fn,
									).
									Rigid(
										wg.th.Caption(fmt.Sprintf("%d ", wg.txs[i].Confirmations)).Fn,
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
											switch wg.txs[i].Category {
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
										wg.th.Caption(wg.txs[i].Category+" ").Fn,
									).
									Fn,
							).
							Rigid(
								wg.th.Flex().AlignMiddle().
									Rigid(
										wg.Icon().Color("DocText").Scale(1).Src(&icons2.DeviceAccessTime).Fn,
									).
									// Rigid(
									// 	wg.th.Caption(
									// 		wg.txs[i].time,
									// 	).Color("DocText").Fn,
									// ).
									Fn,
							).Fn,
					).Fn,
				).Fn,
			).Fn,
		).Fn,
	).Fn(gtx)
}

// func (wg *WalletGUI) ClearAddress(i int) {
//	wg.sendAddresses = remove(wg.sendAddresses, i)
// }

func (wg *WalletGUI) txItem(label, data string) l.Widget {
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
	// TODO: this page doesn't have data being populated yet
	if true {
		return func() {}
	}
	txLayout := []l.Widget{
		wg.txItem("TxId:", wg.State.txs[i].data.TxID),
		wg.txItem("Comment:", wg.State.txs[i].data.Comment),
		wg.txItem("Category:", wg.State.txs[i].data.Category),
		wg.txItem("Address:", wg.State.txs[i].data.Address),
		wg.txItem("Generated:", fmt.Sprint(wg.State.txs[i].data.Generated)),
		wg.txItem("BIP125Replaceable:", wg.State.txs[i].data.BIP125Replaceable),
		wg.txItem("Block Hash:", wg.State.txs[i].data.BlockHash),
		wg.txItem("Block Index:", fmt.Sprint(wg.State.txs[i].data.BlockIndex)),
		wg.txItem("BlockTime:", fmt.Sprint(wg.State.txs[i].data.BlockTime)),
		wg.txItem("Category:", wg.State.txs[i].data.Category),
		wg.txItem("Confirmations:", fmt.Sprint(wg.State.txs[i].data.Confirmations)),
		wg.txItem("Fee:", fmt.Sprint(wg.State.txs[i].data.Fee)),
		wg.txItem("InvolvesWatchOnly:", fmt.Sprint(wg.State.txs[i].data.InvolvesWatchOnly)),
		wg.txItem("Time:", fmt.Sprint(wg.State.txs[i].data.Time)),
		wg.txItem("TimeReceived:", fmt.Sprint(wg.State.txs[i].data.TimeReceived)),
		wg.txItem("Vout:", fmt.Sprint(wg.State.txs[i].data.Vout)),
		wg.txItem("WalletConflicts:", fmt.Sprint(wg.State.txs[i].data.WalletConflicts)),
		wg.txItem("Comment:", wg.State.txs[i].data.Comment),
		wg.txItem("OtherAccount:", wg.State.txs[i].data.OtherAccount),
	}
	le := func(gtx l.Context, index int) l.Dimensions {
		return txLayout[index](gtx)
	}

	return func() {
		wg.w[wg.State.txs[i].data.TxID] = f.NewWindow(wg.th)
		go func() {
			if err := wg.w[wg.State.txs[i].data.TxID].
				Size(64, 32).
				Title("Tx: "+wg.State.txs[i].data.TxID).
				Open().
				Run(
					wg.th.VFlex().
						Rigid(
							wg.Inset(0.0, wg.Fill("Primary", wg.Inset(0.5, wg.Caption(wg.State.txs[i].data.TxID).Color("DocBg").Fn).Fn).Fn).Fn,
						).
						Flexed(1,
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
					func(gtx l.Context) {},
					func() {
						Debug("closing tx window", wg.State.txs[i].data.TxID)
					},
					wg.quit,
				); Check(err) {
			}

		}()
	}
}
