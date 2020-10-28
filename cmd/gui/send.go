package gui

import (
	"fmt"
	l "gioui.org/layout"
	"gioui.org/text"
	"github.com/p9c/pod/pkg/gui/p9"
)

var (
	sendAddresses     []*SendAddress
	sendAddressesList = &l.List{
		Axis: l.Vertical,
	}
	sendChooseFeeBtn = new(p9.Clickable)
	sendClearAllBtn  = new(p9.Clickable)
	sendBtn          = new(p9.Clickable)
	addRecipientBtn  = new(p9.Clickable)
	//amountCnt        = &counter.Counter{
	//	Value:        1,
	//	OperateValue: 1,
	//	From:         1,
	//	To:           999,
	//	CounterInput: &p9.Editor{
	//		//Alignment:  text.Middle,
	//		//SingleLine: true,
	//		//Submit:     true,
	//	},
	//	PageFunction:    func() {},
	//	CounterIncrease: new(p9.Clickable),
	//	CounterDecrease: new(p9.Clickable),
	//	CounterReset:    new(p9.Clickable),
	//}
)

type SendAddress struct {
	AddressInput      *p9.Editor
	AddressBookBtn    *p9.Clickable
	PasteClipboardBtn *p9.Clickable
	ClearBtn          *p9.Clickable
	LabelInput        *p9.Editor
	//AmountInput       *counter.Counter
	SubtractFee     *p9.Bool
	AllAvailableBtn *p9.Clickable
}

func (wg *WalletGUI) SendPage() l.Widget {
	//var out []l.Widget
	fmt.Print("sss", sendAddresses)
	le := func(gtx l.Context, index int) l.Dimensions {
		return wg.Caption("BalaaaaaaaaaaaaaaaO").Color("DocText").Fn(gtx)
	}
	return wg.th.VFlex().
		Flexed(1,

			//wg.Inset(0.0, wg.Fill("DocText", wg.Inset(0.5, wg.H6("title").Color("DocBg").Fn).Fn).Fn).Fn,
			wg.lists["send"].Vertical().Length(len(sendAddresses)).ListElement(le).Fn,

		//func(gtx l.Context) l.Dimensions {
		//return sendAddressesList.Layout(gtx, len(sendAddresses), wg.singleSendAddress())
		//return l.Dimensions{}
		//}
		).
		Rigid(
			wg.sendFooter(),
		).Fn
}

func (wg *WalletGUI) GetSend() {

}

func (wg *WalletGUI) CreateSendAddressItem() {
	sendAddresses = append(sendAddresses,
		&SendAddress{
			AddressInput: &p9.Editor{
				//SingleLine: true,
				//Submit:     true,
			},
			LabelInput: &p9.Editor{
				//SingleLine: true,
				//Submit:     true,
			},
			//AmountInput: &counter.Counter{
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
			//},
			AddressBookBtn:    new(p9.Clickable),
			PasteClipboardBtn: new(p9.Clickable),
			ClearBtn:          new(p9.Clickable),
			SubtractFee:       new(p9.Bool),
			AllAvailableBtn:   new(p9.Clickable),
		})
}

//func ClearAddress(i int) func() {
//	return func() {
//		sendAddresses = remove(sendAddresses, i)
//	}
//}
//
//func ClearAllAddresses() func() {
//	return func() {
//		sendAddresses = []*SendAddress{}
//		CreateSendAddressItem()()
//	}
//}
//

func (wg *WalletGUI) sendFooter() l.Widget {
	return wg.Inset(0.25,
		wg.th.VFlex().
			Rigid(
				wg.Inset(0.25,
					wg.th.Flex().
						SpaceBetween().
						Rigid(
							wg.Inset(0.0, wg.Fill("DocBg", wg.Inset(0.5, wg.Caption("Transaction Fee:").Color("DocText").Fn).Fn).Fn).Fn,
						).
						Rigid(
							wg.Inset(0.0, wg.Fill("DocBg", wg.Inset(0.5, wg.Caption("0.00000 DUO/kb").Color("DocText").Fn).Fn).Fn).Fn,
						).
						Rigid(
							wg.Inset(0.0, wg.Fill("DocBg", wg.Inset(0.5, wg.Caption("net").Color("DocText").Fn).Fn).Fn).Fn,
						).
						Flexed(1,
							wg.Inset(0.0, wg.Fill("DocBg", wg.Inset(0.5, wg.Caption("Balance 0.00000 DUO").Color("DocText").Fn).Fn).Fn).Fn,
						).Fn,
				).Fn,
			).Rigid(
			wg.Inset(0.25,
				wg.th.Flex().
					SpaceBetween().
					Rigid(
						wg.sendButton(wg.clickables["send"], "Send", func() {}),
					).
					Rigid(
						wg.sendButton(wg.clickables["send"], "Clear All", func() {}),
					).
					Rigid(
						wg.sendButton(wg.clickables["send"], "Add Recipient", func() {}),
					).
					Flexed(1,
						wg.Inset(0.0, wg.Caption("Balance:0.00000000").Alignment(text.End).Color("DocText").Fn).Fn,
					).Fn,
			).Fn,
		).Fn,
	).Fn

	//return func(gtx gui.C) gui.D {
	//				return lyt.Format(gtx, "hmax(hflex(middle,r(_),r(inset(0dp8dp0dp8dp,_)),r(_),f(1,inset(8dp8dp8dp8dp,_))))",
	//					sendButton(th, sendBtn, "Send", "hflex(r(_),r(_)))", "Send", func() {}),
	//					sendButton(th, sendClearAllBtn, "Close", "hflex(r(_),r(_)))", "Clear All", ClearAllAddresses()),
	//					sendButton(th, addRecipientBtn, "CounterPlus", "hflex(r(_),r(_)))", "Add Recipient", CreateSendAddressItem()),
	//					func(gtx gui.C) gui.D {
	//						title := theme.H6(th, "Balance: 15.26656.5664654 DUO")
	//						title.Alignment = text.Start
	//						return title.Layout(gtx)
	//					},
	//				)
	//			}
	//		}
}

//func (wg *WalletGUI) sendBody() func(gtx gui.C) gui.D {
//	return box.BoxPanel(g.ui.Theme, func(gtx gui.C) gui.D {
//		return lyt.Format(gtx, "max(hflex(middle,f(1,_)))",
//			func(gtx gui.C) gui.D {
//				return sendAddressesList.Layout(gtx, len(sendAddresses), singleAddress(g.ui.Theme))
//			})
//	})
//}
//
func (wg *WalletGUI) singleSendAddress(gtx l.Context, i int) l.Widget {
	return wg.Inset(0.25,
		wg.th.Flex().
			SpaceBetween().
			Rigid(
				wg.sendButton(wg.clickables["send"], "Send", func() {}),
			).
			Rigid(
				wg.sendButton(wg.clickables["send"], "Clear All", func() {}),
			).
			Rigid(
				wg.sendButton(wg.clickables["send"], "Add Recipient", func() {}),
			).
			Rigid(
				wg.Inset(0.5, wg.Caption("Balance:0.00000000").Alignment(text.End).Color("DocText").Fn).Fn,
			).Fn,
	).Fn
	//return func(gtx gui.C, i int) gui.D {
	//	return lyt.Format(gtx, "vflex(start,r(inset(0dp0dp4dp0dp,_)),r(inset(0dp0dp4dp0dp,_)),r(inset(0dp0dp4dp0dp,_)),r(inset(0dp0dp0dp0dp,_))))",
	//		gui.labeledRow(th, "Pay to:",
	//			func(gtx gui.C) gui.D {
	//				return lyt.Format(gtx, "hflex(middle,f(1,inset(0dp0dp0dp0dp,_)),r(inset(0dp0dp0dp4dp,_)),r(inset(0dp4dp0dp4dp,_)),r(_))",
	//					box.BoxEditor(th, func(gtx gui.C) gui.D {
	//						gtx.Constraints.Min.X = gtx.Constraints.Max.X
	//						e := material.Editor(th.T, sendAddresses[i].AddressInput, "Enter a ParallelCoin address (e.g. 9ef0sdjifvmlkdsfnsdlkg)")
	//						return e.Layout(gtx)
	//					}),
	//					sendButton(th, sendAddresses[i].AddressBookBtn, "AddressBook", "max(inset(0dp0dp0dp0dp,_))", "", func() {}),
	//					sendButton(th, sendAddresses[i].PasteClipboardBtn, "Paste", "max(inset(0dp0dp0dp0dp,_))", "", func() {}),
	//					sendButton(th, sendAddresses[i].ClearBtn, "Close", "max(inset(0dp0dp0dp0dp,_))", "", ClearAddress(i)),
	//				)
	//			}),
	//		gui.labeledRow(th, "Label:",
	//			func(gtx gui.C) gui.D {
	//				return lyt.Format(gtx, "hflex(middle,f(1,inset(0dp0dp0dp0dp,_)))",
	//					box.BoxEditor(th, func(gtx gui.C) gui.D {
	//						gtx.Constraints.Min.X = gtx.Constraints.Max.X
	//						e := material.Editor(th.T, sendAddresses[i].LabelInput, "Enter a label for this address to add it to the list of used addresses")
	//						return e.Layout(gtx)
	//					}),
	//				)
	//			}),
	//		gui.labeledRow(th, "Amount:",
	//			func(gtx gui.C) gui.D {
	//				return lyt.Format(gtx, "hflex(middle,r(_),r(_),r(_),r(_))",
	//					counter.CounterSt(th, sendAddresses[i].AmountInput).Layout(th, fmt.Sprint(sendAddresses[i].AmountInput.Value)),
	//					//func(gtx C) D {return D{}},
	//					func(gtx gui.C) gui.D {
	//						btn := material.IconButton(th.T, gui.connectionsBtn, th.Icons["networkIcon"])
	//						btn.Inset = layout.Inset{unit.Dp(2), unit.Dp(2), unit.Dp(2), unit.Dp(2)}
	//						btn.Size = unit.Dp(21)
	//						btn.Background = helper.HexARGB(th.Colors["Secondary"])
	//						for gui.connectionsBtn.Clicked() {
	//							//ui.N.CurrentPage = "Welcome"
	//						}
	//						return btn.Layout(gtx)
	//					},
	//					func(gtx gui.C) gui.D {
	//						return material.CheckBox(th.T, sendAddresses[i].SubtractFee, "Subtract fee from amount").Layout(gtx)
	//					},
	//					func(gtx gui.C) gui.D {
	//						btn := material.Button(th.T, sendAddresses[i].AllAvailableBtn, "Use available balance")
	//						btn.Inset = layout.Inset{unit.Dp(2), unit.Dp(2), unit.Dp(2), unit.Dp(2)}
	//						btn.Background = helper.HexARGB(th.Colors["Secondary"])
	//						for gui.connectionsBtn.Clicked() {
	//							//ui.N.CurrentPage = "Welcome"
	//						}
	//						return btn.Layout(gtx)
	//					},
	//				)
	//			}),
	//		helper.DuoUIline(false, 8, 0, 1, th.Colors["Border"]),
	//	)
	//}
}

//
//func sendButton(th *theme.Theme, c *widget.Clickable, icon, lay, label string, onClick func()) func(gtx gui.C) gui.D {
//	return box.BoxButton(th, func(gtx gui.C) gui.D {
//		noLabel := true
//		if label != "" {
//			noLabel = true
//		}
//		b := btn.IconTextBtn(th, c, lay, noLabel, label)
//		b.TextSize = unit.Dp(8)
//		b.Background = th.Colors["ButtonBg"]
//		b.Icon = th.Icons[icon]
//		b.IconSize = unit.Dp(15)
//		b.CornerRadius = unit.Dp(0)
//		//b.Background = th.Colors["NavBg"]
//		b.TextColor = th.Colors["ButtonText"]
//		for c.Clicked() {
//			//n.CurrentPage = item.Page
//			onClick()
//		}
//		return b.Layout(gtx)
//	})
//}
//
//func remove(slice []*SendAddress, s int) []*SendAddress {
//	return append(slice[:s], slice[s+1:]...)
//}
func (wg *WalletGUI) sendButton(b *p9.Clickable, title string, click func()) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		gtx.Constraints.Max.X = int(wg.TextSize.Scale(8).V)
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return wg.ButtonLayout(b).Embed(
			func(gtx l.Context) l.Dimensions {
				background := "DocText"
				color := "DocBg"
				var inPad, outPad float32 = 0.5, 0.25
				if *wg.Size >= 800 {
					inPad, outPad = 0.75, 0
				}
				return wg.Inset(outPad,
					wg.Fill(background,
						wg.Flex().
							Flexed(1,
								wg.Inset(inPad,
									wg.Caption(title).
										Color(color).
										Fn,
								).Fn,
							).Fn,
					).Fn,
				).Fn(gtx)
			},
		).
			Background("Transparent").
			SetClick(click).
			Fn(gtx)
	}
}
