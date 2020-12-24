package gui

import (
	"strconv"

	"golang.org/x/exp/shiny/materialdesign/icons"

	l "gioui.org/layout"
	"gioui.org/text"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/coding/base58"
	"github.com/p9c/pod/pkg/util"

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
	// le := func(gtx l.Context, index int) l.Dimensions {
	// 	return wg.singleSendAddress(gtx, index)
	// }
	return wg.th.VFlex().
		Flexed(1,
			// wg.Inset(0.25,
			func(gtx l.Context) l.Dimensions {
				return wg.lists["send"].Vertical().Length(len(wg.sendAddresses)).ListElement(wg.singleSendAddress).Fn(gtx)
			},
			// ).Fn,
		).
		Rigid(
			wg.sendFooter(),
		).Fn
}

func (wg *WalletGUI) CreateSendAddressItem() {
	wg.sendAddresses = append(wg.sendAddresses,
		SendAddress{
			AddressInput: wg.th.Input("", "Enter a ParallelCoin address (e.g. 9ef0sdjifvmlkdsfnsdlkg)", "Primary", "DocText", "DocBg", func(txt string) {
			
			}),
			LabelInput:  wg.th.Input("", "Enter a label for this address to add it to the list of used addresses", "Primary", "DocText", "DocBg", func(pass string) {}),
			AmountInput: wg.th.Input("", "Enter amount", "Primary", "DocText", "DocBg", func(pass string) {}),
			// AmountInput: &counter.Counter{
			//	Value:        1,
			//	OperateValue: 1,
			//	From:         1,
			//	To:           999,
			//	CounterInput: &p9.Editor{
			//		//Alignment:  text.Middle,
			//		//SingleLine: true,
			//		//Submit:     true,
			//	},
			//	//PageFunction:    w.PrikazaniElementSumaRacunica(),
			//	CounterIncrease: new(p9.Clickable),
			//	CounterDecrease: new(p9.Clickable),
			//	CounterReset:    new(p9.Clickable),
			// },
			AddressBookBtn:    new(p9.Clickable),
			PasteClipboardBtn: new(p9.Clickable),
			ClearBtn:          new(p9.Clickable),
			SubtractFee:       new(p9.Bool),
			AllAvailableBtn:   new(p9.Clickable),
		})
}

func (wg *WalletGUI) checkSendItem(address, amount string) (check bool) {
	switch {
	case address == "" && amount != "":
		go wg.toasts.AddToast("Address error", "Address field is empty", "Danger")
	case address != "" && amount == "":
		go wg.toasts.AddToast("Amount error", "Amount is 0", "Danger")
		_, _, err := base58.CheckDecode(address)
		if err != nil {
			Error(err)
			if err == base58.ErrChecksum {
				go wg.toasts.AddToast("Address error", "Address checksum mismatch", "Danger")
			}
			go wg.toasts.AddToast("Address error", "Address is of unknown format", "Danger")
		}
	case address == "" && amount == "":
		go wg.toasts.AddToast("Address and amount error", "Address is empty and amount is 0", "Danger")
	case address != "" && amount != "":
		_, _, err := base58.CheckDecode(address)
		if err != nil {
			Error(err)
			if err == base58.ErrChecksum {
				go wg.toasts.AddToast("Address error", "Address checksum mismatch", "Danger")
			}
			go wg.toasts.AddToast("Address error", "Address is of unknown format", "Danger")
		} else {
			check = true
		}
	}
	return
}

func (wg *WalletGUI) Send() {
	// ToDo Send RPC command
	// TODO: yes, do one like the runner in run.go
	if wg.WalletClient != nil {
		if len(wg.sendAddresses) < 2 {
			if wg.checkSendItem(wg.sendAddresses[0].AddressInput.GetText(), wg.sendAddresses[0].AmountInput.GetText()) {
				address, err := util.DecodeAddress(wg.sendAddresses[0].AddressInput.GetText(), nil)
				if err != nil {
				}
				var h *chainhash.Hash
				amountFloat, err := strconv.ParseFloat(wg.sendAddresses[0].AmountInput.GetText(), 32)
				if err != nil {
				}
				amount, err := util.NewAmount(amountFloat)
				if err != nil {
				}
				if h, err = wg.ChainClient.SendToAddress(address, amount); Check(err) {
					go wg.toasts.AddToast("TxID", h.String(), "Danger")
				}
				// TODO: this is the txid hash
			}
		} else {
			//		for _, sendAddress := range wg.sendAddresses {
			//			fmt.Println(sendAddress.AmountInput.GetText())
			//			address, err := util.DecodeAddress("sendAddress.AmountInput.GetText()", nil)
			//			if err != nil {
			//			}
			//			var h *chainhash.Hash
			//			if h, err = wg.ChainClient.SendToAddress(address, 1); Check(err) {
			//			}
			//			// TODO: this is the txid hash
			//			_ = h
			//		}
		}
	}
}

func (wg *WalletGUI) sendFooter() l.Widget {
	return wg.th.VFlex().
		Rigid(
			wg.th.Inset(0.25,
				wg.th.Flex().
					SpaceBetween().
					Rigid(
						wg.th.Inset(0.0, wg.th.Fill("DocBg",
							wg.th.Inset(0.5,
								wg.th.Caption("Transaction Fee:").
									Color("DocText").Fn,
							).Fn,
						).Fn,
						).Fn,
					).
					Rigid(
						wg.th.Inset(0.0, wg.th.Fill("DocBg",
							wg.th.Inset(0.5,
								wg.th.Caption("0.00000 DUO/kb").
									Color("DocText").Fn,
							).Fn,
						).Fn,
						).Fn,
					).
					Rigid(
						wg.th.Inset(0.0, wg.th.Fill("DocBg",
							wg.th.Inset(0.5,
								wg.th.Caption("net").
									Color("DocText").Fn,
							).Fn,
						).Fn,
						).Fn,
					).
					Flexed(1,
						wg.th.Inset(0.0, wg.th.Fill("DocBg",
							wg.th.Inset(0.5,
								wg.th.Caption("Balance 0.00000 DUO").
									Color("DocText").Fn,
							).Fn,
						).Fn,
						).Fn,
					).Fn,
			).Fn,
		).Rigid(
		wg.th.Inset(0.25,
			wg.th.Flex().
				SpaceBetween().
				Rigid(
					wg.th.Inset(0.25,
						wg.buttonText(wg.clickables["sendSend"],
							"Send", wg.Send)).Fn,
				).
				Rigid(
					wg.th.Inset(0.25,
						wg.buttonText(wg.clickables["sendClearAll"],
							"Clear All", wg.ClearAllAddresses)).Fn,
				).
				Rigid(
					wg.th.Inset(0.25,
						wg.buttonText(wg.clickables["sendAddRecipient"],
							"Add Recipient", wg.CreateSendAddressItem)).Fn,
				).
				Flexed(1,
					wg.th.Inset(0.25,
						wg.th.Caption("Balance:0.00000000").Alignment(text.End).Color("DocText").Fn).Fn,
				).Fn,
		).Fn,
	).Fn
}

func (wg *WalletGUI) singleSendAddress(gtx l.Context, i int) l.Dimensions {
	return wg.th.Inset(0.25,
		wg.th.Fill("DocBg",
			wg.th.Inset(0.25,
				wg.th.VFlex().
					Rigid(
						wg.th.Inset(0.25,
							wg.th.Flex().
								Rigid(
									wg.rowLabel("Pay to:"),
								).
								Rigid(
									wg.th.Flex().
										Rigid(
											wg.sendAddresses[i].AddressInput.Fn,
										).
										Rigid(
											// wg.sendButton(wg.sendAddresses[index].AddressBookBtn, "AddressBook", func() {}),
											// wg.sendIconButton("settings", 2, &icons.ActionBook),
											wg.buttonIcon(wg.sendAddresses[i].AddressBookBtn, "settings", &icons.ActionBook),
										).
										Rigid(
											// wg.sendButton(wg.sendAddresses[index].PasteClipboardBtn, "Paste", func() {}),
											// wg.sendIconButton("settings", 2, &icons.ActionSettings),
											wg.buttonIcon(wg.sendAddresses[i].PasteClipboardBtn, "settings", &icons.ActionSettings),
										).
										Rigid(
											// wg.sendButton(wg.sendAddresses[index].ClearBtn, "Close", func() {}),
											// wg.sendIconButton("settings", 2, &icons.ActionSettings),
											wg.buttonIcon(wg.sendAddresses[i].ClearBtn, "settings", &icons.ActionSettings),
										).Fn,
								).Fn,
						).Fn,
					).
					Rigid(
						wg.th.Inset(0.25,
							wg.th.Flex().
								Rigid(
									wg.rowLabel("Label:"),
								).
								Rigid(
									wg.th.Flex().
										Rigid(
											wg.sendAddresses[i].LabelInput.Fn,
										).Fn,
								).Fn,
						).Fn,
					).
					Rigid(
						wg.th.Inset(0.25,
							wg.th.Flex().
								Rigid(
									wg.rowLabel("Amount:"),
								).
								Rigid(
									wg.th.Flex().
										Rigid(
											wg.sendAddresses[i].AmountInput.Fn,
										).
										Rigid(
											wg.th.Inset(0.25,
												wg.buttonText(wg.sendAddresses[i].PasteClipboardBtn,
													"Subtract fee from amount", func() {})).Fn,
										).
										Rigid(
											wg.th.Inset(0.25,
												wg.buttonText(wg.sendAddresses[i].ClearBtn,
													"Use available balance", func() {})).Fn,
										).Fn,
								).Fn,
						).Fn,
					).Fn,
			).Fn,
		).Fn,
	).Fn(gtx)
}

//
// func (wg *WalletGUI) sendButton(b *p9.Clickable, title string, click func()) func(gtx l.Context) l.Dimensions {
// 	return func(gtx l.Context) l.Dimensions {
// 		gtx.Constraints.Max.X = int(wg.TextSize.Scale(10).V)
// 		gtx.Constraints.Min.X = gtx.Constraints.Max.X
//
// 		return wg.ButtonLayout(b).Embed(
// 			func(gtx l.Context) l.Dimensions {
// 				background := "DocText"
// 				color := "DocBg"
// 				var inPad, outPad float32 = 0.5, 0
// 				return wg.Inset(outPad,
// 					wg.Fill(background,
// 						wg.Flex().
// 							Flexed(1,
// 								wg.Inset(inPad,
// 									wg.Caption(title).
// 										Color(color).
// 										Fn,
// 								).Fn,
// 							).Fn,
// 					).Fn,
// 				).Fn(gtx)
// 			},
// 		).
// 			Background("Transparent").
// 			SetClick(click).
// 			Fn(gtx)
// 	}
// }
//
// func (wg *WalletGUI) sendIconButton(name string, index int, ico *[]byte) func(gtx l.Context) l.Dimensions {
// 	return func(gtx l.Context) l.Dimensions {
// 		background := wg.TitleBarBackgroundGet()
// 		color := wg.MenuColorGet()
// 		if wg.ActivePageGet() == name {
// 			color = "PanelText"
// 			background = "PanelBg"
// 		}
// 		ic := wg.Icon().
// 			Scale(p9.Scales["H5"]).
// 			Color(color).
// 			Src(ico).
// 			Fn
// 		return wg.Flex().Rigid(
// 			// wg.Inset(0.25,
// 			wg.ButtonLayout(wg.buttonBarButtons[index]).
// 				CornerRadius(0).
// 				Embed(
// 					wg.Inset(0.375,
// 						ic,
// 					).Fn,
// 				).
// 				Background(background).
// 				SetClick(
// 					func() {
// 						if wg.MenuOpen {
// 							wg.MenuOpen = false
// 						}
// 						wg.ActivePage(name)
// 					}).
// 				Fn,
// 			// ).Fn,
// 		).Fn(gtx)
// 	}
// }

func (wg *WalletGUI) ClearAddress(i int) {
	wg.sendAddresses = remove(wg.sendAddresses, i)
}

func (wg *WalletGUI) ClearAllAddresses() {
	wg.sendAddresses = []SendAddress{}
	wg.CreateSendAddressItem()
}

func (wg *WalletGUI) rowLabel(label string) l.Widget {
	return func(gtx l.Context) l.Dimensions {
		gtx.Constraints.Max.X = int(wg.th.TextSize.Scale(3).V)
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return wg.th.Caption(label).Color("Primary").Alignment(text.End).Fn(gtx)
	}
}
func remove(slice []SendAddress, s int) []SendAddress {
	return append(slice[:s], slice[s+1:]...)
}
