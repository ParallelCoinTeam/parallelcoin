package gui

import (
	"fmt"
	"strconv"

	"github.com/p9c/matrjoska/pkg/amt"

	"github.com/atotto/clipboard"

	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/text"

	"github.com/p9c/gel"
)

const Break1 = 48

type ReceivePage struct {
	wg                 *WalletGUI
	inputWidth, break1 float32
	sm, md, lg, xl     l.Widget
	urn                string
}

func (wg *WalletGUI) GetReceivePage() (rp *ReceivePage) {
	rp = &ReceivePage{
		wg:         wg,
		inputWidth: 17,
		break1:     48,
	}
	rp.sm = rp.SmallList
	return
}

func (rp *ReceivePage) Fn(gtx l.Context) l.Dimensions {
	wg := rp.wg
	return wg.Responsive(
		wg.Size.Load(), gel.Widgets{
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
		wg.Direction().Center().Embed(rp.QRButton()).Fn,
		rp.InputMessage(),
		rp.AmountInput(),
		rp.MessageInput(),
		rp.RegenerateButton(),
		rp.AddressbookHeader(),
	}
	smallWidgets = append(smallWidgets, rp.GetAddressbookHistoryCards("DocBg")...)
	le := func(gtx l.Context, index int) l.Dimensions {
		return wg.Inset(0.25, smallWidgets[index]).Fn(gtx)
	}
	return wg.VFlex().AlignStart().
		Flexed(
			1,
			wg.lists["receive"].
				Vertical().Start().
				Length(len(smallWidgets)).
				ListElement(le).Fn,
		).Fn(gtx)
}

func (rp *ReceivePage) InputMessage() l.Widget {
	return rp.wg.Body2("Input details to request a payment").Alignment(text.Middle).Fn
}

func (rp *ReceivePage) MediumList(gtx l.Context) l.Dimensions {
	wg := rp.wg
	qrWidget := []l.Widget{
		wg.Direction().Center().Embed(rp.QRButton()).Fn,
		rp.InputMessage(),
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
		return wg.Inset(
			0.25,
			historyWidget[index],
		).Fn(gtx)
	}
	return wg.Flex().AlignStart().
		Rigid(
			func(gtx l.Context) l.Dimensions {
				gtx.Constraints.Max.X, gtx.Constraints.Min.X = int(wg.TextSize.V*rp.inputWidth),
					int(wg.TextSize.V*rp.inputWidth)
				return wg.VFlex().
					Rigid(
						wg.lists["receiveMedium"].
							Vertical().
							Length(len(qrWidget)).
							ListElement(qrLE).Fn,
					).Fn(gtx)
			},
		).
		Rigid(
			wg.VFlex().AlignStart().
				Rigid(
					rp.AddressbookHeader(),
				).
				Rigid(
					wg.lists["receiveAddresses"].
						Vertical().
						Length(len(historyWidget)).
						ListElement(historyLE).Fn,
				).
				Fn,
		).Fn(gtx)
}

func (rp *ReceivePage) Spacer() l.Widget {
	return rp.wg.Flex().AlignMiddle().Flexed(1, rp.wg.Inset(0.25, gel.EmptySpace(0, 0)).Fn).Fn
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
							D.Ln("clicked receive address list item", j)
							if e := clipboard.WriteAll(qrText); E.Chk(e) {
							}
							wg.GetNewReceivingQRCode(qrText)
							rp.urn = qrText
						},
					),
				).
					Background(bg).
					Embed(
						wg.Inset(
							0.25,
							wg.VFlex().AlignStart().
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
	if !wg.WalletAndClientRunning() || wg.currentReceiveQRCode == nil {
		return func(gtx l.Context) l.Dimensions {
			return l.Dimensions{}
		}
	}
	return wg.VFlex().
		Rigid(
			wg.ButtonLayout(
				wg.currentReceiveCopyClickable.SetClick(
					func() {
						D.Ln("clicked qr code copy clicker")
						if e := clipboard.WriteAll(rp.urn); E.Chk(e) {
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
		).Rigid(
		rp.QRMessage(),
	).Fn
}

func (rp *ReceivePage) AddressbookHeader() l.Widget {
	wg := rp.wg
	return wg.Flex().
		Rigid(
			wg.Inset(
				0.25,
				wg.H5("Receive Address History").Alignment(text.Middle).Fn,
			).Fn,
		).Fn
}

func (rp *ReceivePage) AmountInput() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		wg := rp.wg
		// gtx.Constraints.Max.X, gtx.Constraints.Min.X = int(wg.TextSize.True*rp.inputWidth), int(wg.TextSize.True*rp.inputWidth)
		return wg.inputs["receiveAmount"].Fn(gtx)
	}
}

func (rp *ReceivePage) MessageInput() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		wg := rp.wg
		// gtx.Constraints.Max.X, gtx.Constraints.Min.X = int(wg.TextSize.True*rp.inputWidth), int(wg.TextSize.True*rp.inputWidth)
		return wg.inputs["receiveMessage"].Fn(gtx)
	}
}

func (rp *ReceivePage) RegenerateButton() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		wg := rp.wg
		if wg.inputs["receiveAmount"].GetText() == "" || wg.inputs["receiveMessage"].GetText() == "" {
			gtx.Queue = nil
		}
		// gtx.Constraints.Max.X, gtx.Constraints.Min.X = int(wg.TextSize.True*rp.inputWidth), int(wg.TextSize.True*rp.inputWidth)
		return wg.ButtonLayout(
			wg.currentReceiveRegenClickable.
				SetClick(
					func() {
						D.Ln("clicked regenerate button")
						var amount float64
						var am amt.Amount
						var e error
						if amount, e = strconv.ParseFloat(
							wg.inputs["receiveAmount"].GetText(),
							64,
						); !E.Chk(e) {
							if am, e = amt.NewAmount(amount); E.Chk(e) {
							}
						}
						msg := wg.inputs["receiveMessage"].GetText()
						if am == 0 || msg == "" {
							// never store an entry without both fields filled
							return
						}
						if len(wg.State.receiveAddresses) > 0 &&
							(wg.State.receiveAddresses[len(wg.State.receiveAddresses)-1].Amount == 0 ||
								wg.State.receiveAddresses[len(wg.State.receiveAddresses)-1].Message == "") {
							// the first entry has neither of these, and newly generated items without them are assumed to
							// not be intentional or used addresses so we don't generate a new entry for this case
							wg.State.receiveAddresses[len(wg.State.receiveAddresses)-1].Amount = am
							wg.State.receiveAddresses[len(wg.State.receiveAddresses)-1].Message = msg
						} else {
							// go func() {
							wg.GetNewReceivingAddress()
							msg := wg.inputs["receiveMessage"].GetText()
							if len(msg) > 64 {
								msg = msg[:64]
								// enforce the field length limit
								wg.inputs["receiveMessage"].SetText(msg)
							}
							qrText := fmt.Sprintf(
								"parallelcoin:%s?amount=%f&message=%s",
								wg.State.currentReceivingAddress.Load().EncodeAddress(),
								am.ToDUO(),
								msg,
							)
							rp.urn = qrText
							wg.GetNewReceivingQRCode(rp.urn)
							// }()
						}
						// force user to fill fields again after regenerate to stop duplicate entries especially from
						// accidental double clicks/taps
						wg.inputs["receiveAmount"].SetText("")
						wg.inputs["receiveMessage"].SetText("")
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
