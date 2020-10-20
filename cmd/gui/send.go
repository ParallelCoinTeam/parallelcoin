package gui

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gioapp/gel/helper"
	"github.com/p9c/pod/pkg/gui/wallet/btn"
	"github.com/p9c/pod/pkg/gui/wallet/counter"
	"github.com/p9c/pod/pkg/gui/wallet/dap/box"
	"github.com/p9c/pod/pkg/gui/wallet/lyt"
	"github.com/p9c/pod/pkg/gui/wallet/theme"
)

var (
	sendAddresses     []*SendAddress
	sendAddressesList = &layout.List{
		Axis: layout.Vertical,
	}
	sendChooseFeeBtn = new(widget.Clickable)
	sendClearAllBtn  = new(widget.Clickable)
	sendBtn          = new(widget.Clickable)
	addRecipientBtn  = new(widget.Clickable)
	amountCnt        = &counter.Counter{
		Value:        1,
		OperateValue: 1,
		From:         1,
		To:           999,
		CounterInput: &widget.Editor{
			Alignment:  text.Middle,
			SingleLine: true,
			Submit:     true,
		},
		PageFunction:    func() {},
		CounterIncrease: new(widget.Clickable),
		CounterDecrease: new(widget.Clickable),
		CounterReset:    new(widget.Clickable),
	}
)

type SendAddress struct {
	AddressInput      *widget.Editor
	AddressBookBtn    *widget.Clickable
	PasteClipboardBtn *widget.Clickable
	ClearBtn          *widget.Clickable
	LabelInput        *widget.Editor
	AmountInput       *counter.Counter
	SubtractFee       *widget.Bool
	AllAvailableBtn   *widget.Clickable
}

func (g *GuiAppModel) GetSend() {

}

func CreateSendAddressItem() func() {
	return func() {
		sendAddresses = append(sendAddresses,
			&SendAddress{
				AddressInput: &widget.Editor{
					SingleLine: true,
					Submit:     true,
				},
				LabelInput: &widget.Editor{
					SingleLine: true,
					Submit:     true,
				},
				AmountInput: &counter.Counter{
					Value:        1,
					OperateValue: 1,
					From:         1,
					To:           999,
					CounterInput: &widget.Editor{
						Alignment:  text.Middle,
						SingleLine: true,
						Submit:     true,
					},
					//PageFunction:    w.PrikazaniElementSumaRacunica(),
					CounterIncrease: new(widget.Clickable),
					CounterDecrease: new(widget.Clickable),
					CounterReset:    new(widget.Clickable),
				},
				AddressBookBtn:    new(widget.Clickable),
				PasteClipboardBtn: new(widget.Clickable),
				ClearBtn:          new(widget.Clickable),
				SubtractFee:       new(widget.Bool),
				AllAvailableBtn:   new(widget.Clickable),
			})
	}
}
func ClearAddress(i int) func() {
	return func() {
		sendAddresses = remove(sendAddresses, i)
	}
}
func ClearAllAddresses() func() {
	return func() {
		sendAddresses = []*SendAddress{}
		CreateSendAddressItem()()
	}
}

func (g *GuiAppModel) sendHeader() func(gtx C) D {
	return box.BoxPanel(g.ui.Theme, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		helper.Fill(gtx, helper.HexARGB(g.ui.Theme.Colors["PanelBg"]))
		return lyt.Format(gtx, "vflex(middle,r(inset(5dp0dp5dp0dp,_))))",
			func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				title := theme.H6(g.ui.Theme, "Send Header")
				title.Alignment = text.Start
				return title.Layout(gtx)
			})
	})
}

func (g *GuiAppModel) sendFooter() func(gtx C) D {
	return func(gtx C) D {
		return lyt.Format(gtx, "vflex(middle,r(inset(5dp0dp5dp0dp,_)),r(inset(5dp0dp5dp0dp,_))))",
			sendFooterTop(g.ui.Theme),
			sendFooterBottom(g.ui.Theme))
	}
}

func sendFooterBottom(th *theme.Theme) func(gtx C) D {
	return func(gtx C) D {
		return lyt.Format(gtx, "hmax(hflex(middle,r(_),r(inset(0dp8dp0dp8dp,_)),r(_),f(1,inset(8dp8dp8dp8dp,_))))",
			sendButton(th, sendBtn, "Send", "hflex(r(_),r(_)))", "Send", func() {}),
			sendButton(th, sendClearAllBtn, "Close", "hflex(r(_),r(_)))", "Clear All", ClearAllAddresses()),
			sendButton(th, addRecipientBtn, "CounterPlus", "hflex(r(_),r(_)))", "Add Recipient", CreateSendAddressItem()),
			func(gtx C) D {
				title := theme.H6(th, "Balance: 15.26656.5664654 DUO")
				title.Alignment = text.Start
				return title.Layout(gtx)
			},
		)
	}
}

func sendFooterTop(th *theme.Theme) func(gtx C) D {
	return box.BoxPanel(th, func(gtx C) D {
		return lyt.Format(gtx, "hflex(middle,r(_),r(inset(0dp8dp0dp8dp,_)),r(_),f(1,inset(8dp8dp8dp8dp,_)))",
			func(gtx C) D {
				title := theme.H6(th, "Transaction Fee:")
				title.Alignment = text.Start
				return title.Layout(gtx)
			},
			func(gtx C) D {
				title := theme.H6(th, "0.00000 DUO/kb")
				title.Alignment = text.Start
				return title.Layout(gtx)
			},
			func(gtx C) D {
				btn := material.IconButton(th.T, connectionsBtn, th.Icons["networkIcon"])
				btn.Inset = layout.Inset{unit.Dp(2), unit.Dp(2), unit.Dp(2), unit.Dp(2)}
				btn.Size = unit.Dp(21)
				btn.Background = helper.HexARGB(th.Colors["Secondary"])
				for connectionsBtn.Clicked() {
					//ui.N.CurrentPage = "Welcome"
				}
				return btn.Layout(gtx)
			},
			func(gtx C) D {
				title := theme.H6(th, "Warning: Fee estimation is currently possible.")
				title.Alignment = text.Start
				return title.Layout(gtx)
			},
		)
	})
}
func (g *GuiAppModel) sendBody() func(gtx C) D {
	return box.BoxPanel(g.ui.Theme, func(gtx C) D {
		return lyt.Format(gtx, "max(hflex(middle,f(1,_)))",
			func(gtx C) D {
				return sendAddressesList.Layout(gtx, len(sendAddresses), singleAddress(g.ui.Theme))
			})
	})
}

func singleAddress(th *theme.Theme) func(gtx C, i int) D {
	return func(gtx C, i int) D {
		return lyt.Format(gtx, "vflex(start,r(inset(0dp0dp4dp0dp,_)),r(inset(0dp0dp4dp0dp,_)),r(inset(0dp0dp4dp0dp,_)),r(inset(0dp0dp0dp0dp,_))))",
			labeledRow(th, "Pay to:",
				func(gtx C) D {
					return lyt.Format(gtx, "hflex(middle,f(1,inset(0dp0dp0dp0dp,_)),r(inset(0dp0dp0dp4dp,_)),r(inset(0dp4dp0dp4dp,_)),r(_))",
						box.BoxEditor(th, func(gtx C) D {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							e := material.Editor(th.T, sendAddresses[i].AddressInput, "Enter a ParallelCoin address (e.g. 9ef0sdjifvmlkdsfnsdlkg)")
							return e.Layout(gtx)
						}),
						sendButton(th, sendAddresses[i].AddressBookBtn, "AddressBook", "max(inset(0dp0dp0dp0dp,_))", "", func() {}),
						sendButton(th, sendAddresses[i].PasteClipboardBtn, "Paste", "max(inset(0dp0dp0dp0dp,_))", "", func() {}),
						sendButton(th, sendAddresses[i].ClearBtn, "Close", "max(inset(0dp0dp0dp0dp,_))", "", ClearAddress(i)),
					)
				}),
			labeledRow(th, "Label:",
				func(gtx C) D {
					return lyt.Format(gtx, "hflex(middle,f(1,inset(0dp0dp0dp0dp,_)))",
						box.BoxEditor(th, func(gtx C) D {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							e := material.Editor(th.T, sendAddresses[i].LabelInput, "Enter a label for this address to add it to the list of used addresses")
							return e.Layout(gtx)
						}),
					)
				}),
			labeledRow(th, "Amount:",
				func(gtx C) D {
					return lyt.Format(gtx, "hflex(middle,r(_),r(_),r(_),r(_))",
						counter.CounterSt(th, sendAddresses[i].AmountInput).Layout(th, fmt.Sprint(sendAddresses[i].AmountInput.Value)),
						//func(gtx C) D {return D{}},
						func(gtx C) D {
							btn := material.IconButton(th.T, connectionsBtn, th.Icons["networkIcon"])
							btn.Inset = layout.Inset{unit.Dp(2), unit.Dp(2), unit.Dp(2), unit.Dp(2)}
							btn.Size = unit.Dp(21)
							btn.Background = helper.HexARGB(th.Colors["Secondary"])
							for connectionsBtn.Clicked() {
								//ui.N.CurrentPage = "Welcome"
							}
							return btn.Layout(gtx)
						},
						func(gtx C) D {
							return material.CheckBox(th.T, sendAddresses[i].SubtractFee, "Subtract fee from amount").Layout(gtx)
						},
						func(gtx C) D {
							btn := material.Button(th.T, sendAddresses[i].AllAvailableBtn, "Use available balance")
							btn.Inset = layout.Inset{unit.Dp(2), unit.Dp(2), unit.Dp(2), unit.Dp(2)}
							btn.Background = helper.HexARGB(th.Colors["Secondary"])
							for connectionsBtn.Clicked() {
								//ui.N.CurrentPage = "Welcome"
							}
							return btn.Layout(gtx)
						},
					)
				}),
			helper.DuoUIline(false, 8, 0, 1, th.Colors["Border"]),
		)
	}
}

func sendButton(th *theme.Theme, c *widget.Clickable, icon, lay, label string, onClick func()) func(gtx C) D {
	return box.BoxButton(th, func(gtx C) D {
		noLabel := true
		if label != "" {
			noLabel = true
		}
		b := btn.IconTextBtn(th, c, lay, noLabel, label)
		b.TextSize = unit.Dp(8)
		b.Background = th.Colors["ButtonBg"]
		b.Icon = th.Icons[icon]
		b.IconSize = unit.Dp(15)
		b.CornerRadius = unit.Dp(0)
		//b.Background = th.Colors["NavBg"]
		b.TextColor = th.Colors["ButtonText"]
		for c.Clicked() {
			//n.CurrentPage = item.Page
			onClick()
		}
		return b.Layout(gtx)
	})
}

func remove(slice []*SendAddress, s int) []*SendAddress {
	return append(slice[:s], slice[s+1:]...)
}
