package gui

import (
	"fmt"
	"strconv"
	
	l "gioui.org/layout"
	"gioui.org/text"
	"github.com/atotto/clipboard"
	
	"github.com/p9c/pod/pkg/gui"
	"github.com/p9c/pod/pkg/util"
)

const Break1 = 48

type ReceivePage struct {
	wg                 *WalletGUI
	inputWidth, break1 float32
	sm, md, lg, xl     l.Widget
}

func (wg *WalletGUI) GetReceivePage() (rp *ReceivePage) {
	rp = &ReceivePage{
		wg:         wg,
		inputWidth: 20,
		break1:     48,
	}
	rp.sm = rp.SmallList
	return
}

func (rp *ReceivePage) Fn(gtx l.Context) l.Dimensions {
	wg := rp.wg
	return wg.Responsive(
		*wg.Size, gui.Widgets{
			{
				Widget: rp.SmallList,
			},
			{
				Size:   rp.break1,
				Widget: rp.MediumList,
			},
		},
	).Fn(gtx)
}

func (rp *ReceivePage) SmallList(gtx l.Context) l.Dimensions {
	wg := rp.wg
	smallWidgets := []l.Widget{
		rp.QRMessage(),
		wg.Direction().Center().Embed(rp.QRButton()).Fn,
		rp.AmountInput(),
		rp.MessageInput(),
		rp.RegenerateButton(),
		rp.AddressbookHeader(),
	}
	smallWidgets = append(smallWidgets, rp.GetAddressbookHistoryCards("DocBg")...)
	le := func(gtx l.Context, index int) l.Dimensions {
		return wg.Inset(0.25, smallWidgets[index]).Fn(gtx)
	}
	return wg.lists["receive"].
		Vertical().
		Length(len(smallWidgets)).
		ListElement(le).Fn(gtx)
}

func (rp *ReceivePage) MediumList(gtx l.Context) l.Dimensions {
	wg := rp.wg
	qrWidget := []l.Widget{
		rp.QRMessage(),
		wg.Direction().Center().Embed(rp.QRButton()).Fn,
		rp.AmountInput(),
		rp.MessageInput(),
		rp.RegenerateButton(),
		// rp.AddressbookHeader(),
	}
	qrLE := func(gtx l.Context, index int) l.Dimensions {
		return wg.Inset(0.25, qrWidget[index]).Fn(gtx)
	}
	var historyWidget []l.Widget
	
	historyWidget = append(historyWidget, rp.GetAddressbookHistoryCards("DocBg")...)
	historyLE := func(gtx l.Context, index int) l.Dimensions {
		return wg.Inset(0.25,
			historyWidget[index],
		).Fn(gtx)
	}
	return wg.Flex().
		Rigid(
			func(gtx l.Context) l.Dimensions {
				gtx.Constraints.Max.X, gtx.Constraints.Min.X = int(wg.TextSize.V*rp.inputWidth),
					int(wg.TextSize.V*rp.inputWidth)
				return wg.lists["receiveMedium"].
					Vertical().
					Length(len(qrWidget)).
					ListElement(qrLE).Fn(gtx)
			},
		).
		Flexed(
			1,
			wg.VFlex().Rigid(
				rp.AddressbookHeader(),
			).Flexed(
				1,
				wg.lists["receiveAddresses"].
					Vertical().
					Length(len(historyWidget)).
					ListElement(historyLE).Fn,
			).Fn,
		).Fn(gtx)
}

func (rp *ReceivePage) Spacer() l.Widget {
	return rp.wg.Flex().AlignMiddle().Flexed(1, rp.wg.Inset(0.5, gui.EmptySpace(0, 0)).Fn).Fn
}

func (rp *ReceivePage) GetAddressbookHistoryCards(bg string) (widgets []l.Widget) {
	wg := rp.wg
	avail := len(wg.receiveAddressbookClickables)
	req := len(wg.State.receiveAddresses)
	if req > avail {
		for i := 0; i < req-avail; i++ {
			wg.receiveAddressbookClickables = append(wg.receiveAddressbookClickables, wg.WidgetPool.GetClickable())
		}
	}
	for x := range wg.State.receiveAddresses {
		j := x
		i := len(wg.State.receiveAddresses) - 1 - x
		widgets = append(
			widgets, func(gtx l.Context) l.Dimensions {
				return wg.ButtonLayout(
					wg.receiveAddressbookClickables[i].SetClick(
						func() {
							msg := wg.State.receiveAddresses[i].Message
							if len(msg) > 64 {
								msg = msg[:64]
							}
							qrText := fmt.Sprintf(
								"parallelcoin:%s?amount=%8.8f&message=%s",
								wg.State.receiveAddresses[i].Address,
								wg.State.receiveAddresses[i].Amount.ToDUO(),
								msg,
							)
							Debug("clicked receive address list item", j)
							if err := clipboard.WriteAll(qrText); Check(err) {
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
											wg.Caption(wg.State.receiveAddresses[i].Address).
												Font("go regular").Fn,
										).
										Flexed(
											1,
											wg.Body1(wg.State.receiveAddresses[i].Amount.String()).
												Alignment(text.End).Fn,
										).
										Fn,
								).
								Rigid(
									wg.Caption(wg.State.receiveAddresses[i].Message).MaxLines(1).Fn,
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

func (rp *ReceivePage) QRMessage() l.Widget {
	return rp.wg.Body2("Scan to send or click to copy").Alignment(text.Middle).Fn
}

func (rp *ReceivePage) GetQRText() string {
	wg := rp.wg
	msg := wg.inputs["receiveMessage"].GetText()
	if len(msg) > 64 {
		msg = msg[:64]
	}
	return fmt.Sprintf(
		"parallelcoin:%s?amount=%s&message=%s",
		wg.State.currentReceivingAddress.Load().EncodeAddress(),
		wg.inputs["receiveAmount"].GetText(),
		msg,
	)
}

func (rp *ReceivePage) QRButton() l.Widget {
	wg := rp.wg
	if !wg.WalletAndClientRunning() {
		return func(gtx l.Context) l.Dimensions {
			return l.Dimensions{}
		}
	}
	if wg.currentReceiveQRCode == nil {
		wg.GetNewReceivingAddress()
		wg.GetNewReceivingQRCode()
	}
	return wg.Flex().Rigid(
		wg.ButtonLayout(
			wg.currentReceiveCopyClickable.SetClick(
				func() {
					Debug("clicked qr code copy clicker")
					if err := clipboard.WriteAll(rp.GetQRText()); Check(err) {
					}
				},
			),
		).
			Background("white").
			Embed(
				wg.Inset(
					0.125,
					wg.Image().Src(*wg.currentReceiveQRCode).Scale(1).Fn,
				).Fn,
			).Fn,
	).Fn
}

func (rp *ReceivePage) AddressbookHeader() l.Widget {
	wg := rp.wg
	return wg.Flex().Flexed(
		1,
		wg.Inset(
			0.25,
			wg.H6("Receive Address History").Alignment(text.Middle).Fn,
		).Fn,
	).Fn
}

func (rp *ReceivePage) AmountInput() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		wg := rp.wg
		// gtx.Constraints.Max.X, gtx.Constraints.Min.X = int(wg.TextSize.V*rp.inputWidth), int(wg.TextSize.V*rp.inputWidth)
		return wg.inputs["receiveAmount"].Fn(gtx)
	}
}

func (rp *ReceivePage) MessageInput() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		wg := rp.wg
		// gtx.Constraints.Max.X, gtx.Constraints.Min.X = int(wg.TextSize.V*rp.inputWidth), int(wg.TextSize.V*rp.inputWidth)
		return wg.inputs["receiveMessage"].Fn(gtx)
	}
}

func (rp *ReceivePage) RegenerateButton() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		wg := rp.wg
		// gtx.Constraints.Max.X, gtx.Constraints.Min.X = int(wg.TextSize.V*rp.inputWidth), int(wg.TextSize.V*rp.inputWidth)
		return wg.ButtonLayout(
			wg.currentReceiveRegenClickable.
				SetClick(
					func() {
						Debug("clicked regenerate button")
						var amt float64
						var am util.Amount
						var err error
						if amt, err = strconv.ParseFloat(
							wg.inputs["receiveAmount"].GetText(),
							64,
						); !Check(err) {
							if am, err = util.NewAmount(amt); Check(err) {
							}
						}
						msg := wg.inputs["receiveMessage"].GetText()
						if am == 0 || msg == "" {
							// never store an entry without both fields filled
							return
						}
						if wg.State.receiveAddresses[len(wg.State.receiveAddresses)-1].Amount == 0 ||
							wg.State.receiveAddresses[len(wg.State.receiveAddresses)-1].Message == "" {
							// the first entry has neither of these, and newly generated items without them are assumed to
							// not be intentional or used addresses so we don't generate a new entry for this case
							wg.State.receiveAddresses[len(wg.State.receiveAddresses)-1].Amount = am
							wg.State.receiveAddresses[len(wg.State.receiveAddresses)-1].Message = msg
						} else {
							go func() {
								wg.GetNewReceivingAddress()
								wg.GetNewReceivingQRCode()
							}()
						}
						wg.invalidate <- struct{}{}
					},
				),
		).
			Background("Primary").
			Embed(
				wg.Inset(
					0.5,
					wg.H6("regenerate").Color("Light").Fn,
				).
					Fn,
			).
			Fn(gtx)
	}
}
