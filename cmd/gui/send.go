package gui

import (
	"fmt"
	
	l "gioui.org/layout"
	"gioui.org/text"
	"github.com/atotto/clipboard"
	
	"github.com/p9c/pod/pkg/gui"
)

type SendPage struct {
	wg                 *WalletGUI
	inputWidth, break1 float32
}

func (wg *WalletGUI) GetSendPage() (sp *SendPage) {
	sp = &SendPage{
		wg:         wg,
		inputWidth: 20,
		break1:     48,
	}
	return
}

func (sp *SendPage) Fn(gtx l.Context) l.Dimensions {
	wg := sp.wg
	return wg.Responsive(
		*wg.Size, gui.Widgets{
			{
				Widget: sp.SmallList,
			},
			{
				Size:   sp.break1,
				Widget: sp.MediumList,
			},
		},
	).Fn(gtx)
}

func (sp *SendPage) SmallList(gtx l.Context) l.Dimensions {
	wg := sp.wg
	smallWidgets := []l.Widget{
		sp.AddressInput(),
		sp.AmountInput(),
		sp.MessageInput(),
		wg.Flex().
			Flexed(1,
				sp.SendButton(),
			).
			Rigid(
				wg.Inset(0.5, gui.EmptySpace(0, 0)).Fn,
			).
			Rigid(
				sp.PasteButton(),
			).
			Rigid(
				wg.Inset(0.5, gui.EmptySpace(0, 0)).Fn,
			).
			Rigid(
				sp.SaveButton(),
			).Fn,
		sp.AddressbookHeader(),
	}
	smallWidgets = append(smallWidgets, sp.GetAddressbookHistoryCards("DocBg")...)
	le := func(gtx l.Context, index int) l.Dimensions {
		return wg.Inset(0.25, smallWidgets[index]).Fn(gtx)
	}
	return wg.lists["send"].
		Vertical().
		Length(len(smallWidgets)).
		ListElement(le).Fn(gtx)
}

func (sp *SendPage) MediumList(gtx l.Context) l.Dimensions {
	wg := sp.wg
	sendFormWidget := []l.Widget{
		sp.AddressInput(),
		sp.AmountInput(),
		sp.MessageInput(),
		wg.Flex().
			Flexed(1,
				sp.SendButton(),
			).
			Rigid(
				wg.Inset(0.5, gui.EmptySpace(0, 0)).Fn,
			).
			Rigid(
				sp.PasteButton(),
			).
			Rigid(
				wg.Inset(0.5, gui.EmptySpace(0, 0)).Fn,
			).
			Rigid(
				sp.SaveButton(),
			).Fn,
	}
	sendLE := func(gtx l.Context, index int) l.Dimensions {
		return wg.Inset(0.25, sendFormWidget[index]).Fn(gtx)
	}
	var historyWidget []l.Widget
	historyWidget = append(historyWidget, sp.GetAddressbookHistoryCards("DocBg")...)
	historyLE := func(gtx l.Context, index int) l.Dimensions {
		return wg.Inset(0.25,
			historyWidget[index],
		).Fn(gtx)
	}
	return wg.Flex().
		Rigid(
			func(gtx l.Context) l.Dimensions {
				gtx.Constraints.Max.X, gtx.Constraints.Min.X = int(wg.TextSize.V*sp.inputWidth),
					int(wg.TextSize.V*sp.inputWidth)
				return wg.lists["send"].
					Vertical().
					Length(len(sendFormWidget)).
					ListElement(sendLE).Fn(gtx)
			},
		).
		Flexed(
			1,
			wg.VFlex().Rigid(
				sp.AddressbookHeader(),
			).Flexed(
				1,
				wg.lists["sendAddresses"].
					Vertical().
					Length(len(historyWidget)).
					ListElement(historyLE).Fn,
			).Fn,
		).Fn(gtx)
}

func (sp *SendPage) AddressInput() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		wg := sp.wg
		return wg.inputs["sendAddress"].Fn(gtx)
	}
}

func (sp *SendPage) AmountInput() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		wg := sp.wg
		return wg.inputs["sendAmount"].Fn(gtx)
	}
}

func (sp *SendPage) MessageInput() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		wg := sp.wg
		return wg.inputs["sendMessage"].Fn(gtx)
	}
}

func (sp *SendPage) SendButton() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		wg := sp.wg
		return wg.ButtonLayout(
			wg.clickables["sendSend"].
				SetClick(
					func() {
						Debug("clicked send button")
					},
				),
		).
			Background("Primary").
			Embed(
				wg.Inset(
					0.5,
					wg.H6("send").Color("Light").Fn,
				).
					Fn,
			).
			Fn(gtx)
	}
}

func (sp *SendPage) SaveButton() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		wg := sp.wg
		return wg.ButtonLayout(
			wg.clickables["sendSave"].
				SetClick(
					func() {
						Debug("clicked save button")
						
					},
				),
		).
			Background("DocBg").
			Embed(
				wg.Inset(
					0.5,
					wg.H6("save").Color("DocText").Fn,
				).
					Fn,
			).
			Fn(gtx)
	}
}

func (sp *SendPage) PasteButton() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		wg := sp.wg
		return wg.ButtonLayout(
			wg.clickables["sendFromRequest"].
				SetClick(
					func() {
						Debug("clicked paste button")
						
					},
				),
		).
			Background("DocBg").
			Embed(
				wg.Inset(
					0.5,
					wg.H6("paste").Color("DocText").Fn,
				).
					Fn,
			).
			Fn(gtx)
	}
}

func (sp *SendPage) AddressbookHeader() l.Widget {
	wg := sp.wg
	return wg.Flex().Flexed(
		1,
		wg.Inset(
			0.25,
			wg.H6("Send Address History").Alignment(text.Middle).Fn,
		).Fn,
	).Fn
}

func (sp *SendPage) GetAddressbookHistoryCards(bg string) (widgets []l.Widget) {
	wg := sp.wg
	avail := len(wg.sendAddressbookClickables)
	req := len(wg.State.sendAddresses)
	if req > avail {
		for i := 0; i < req-avail; i++ {
			wg.sendAddressbookClickables = append(wg.sendAddressbookClickables, wg.WidgetPool.GetClickable())
		}
	}
	for x := range wg.State.sendAddresses {
		j := x
		i := len(wg.State.sendAddresses) - 1 - x
		widgets = append(
			widgets, func(gtx l.Context) l.Dimensions {
				return wg.ButtonLayout(
					wg.sendAddressbookClickables[i].SetClick(
						func() {
							sendText := fmt.Sprintf(
								"parallelcoin:%s?amount=%8.8f&message=%s",
								wg.State.sendAddresses[i].Address,
								wg.State.sendAddresses[i].Amount.ToDUO(),
								wg.State.sendAddresses[i].Message,
							)
							Debug("clicked send address list item", j)
							if err := clipboard.WriteAll(sendText); Check(err) {
							}
						},
					),
				).
					Background(bg).
					Embed(
						wg.Inset(
							0.25,
							wg.VFlex().
								Rigid(
									wg.Flex().AlignBaseline().
										Rigid(
											wg.Caption(wg.State.sendAddresses[i].Address).
												Font("go regular").Fn,
										).
										Flexed(
											1,
											wg.Body1(wg.State.sendAddresses[i].Amount.String()).
												Alignment(text.End).Fn,
										).
										Fn,
								).
								Rigid(
									wg.Caption(wg.State.sendAddresses[i].Message).Fn,
								).
								Fn,
						).
							Fn,
					).Fn(gtx)
			},
		)
	}
	return
}
