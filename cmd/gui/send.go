package gui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/p9c/matrjoska/pkg/amt"
	"github.com/p9c/matrjoska/pkg/btcaddr"

	"github.com/atotto/clipboard"

	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/text"

	"github.com/p9c/gel"
	"github.com/p9c/matrjoska/pkg/chainhash"
)

type SendPage struct {
	wg                 *WalletGUI
	inputWidth, break1 float32
}

func (wg *WalletGUI) GetSendPage() (sp *SendPage) {
	sp = &SendPage{
		wg:         wg,
		inputWidth: 17,
		break1:     48,
	}
	wg.inputs["sendAddress"].SetPasteFunc = sp.pasteFunction
	wg.inputs["sendAmount"].SetPasteFunc = sp.pasteFunction
	wg.inputs["sendMessage"].SetPasteFunc = sp.pasteFunction
	return
}

func (sp *SendPage) Fn(gtx l.Context) l.Dimensions {
	wg := sp.wg
	return wg.Responsive(
		wg.Size.Load(), gel.Widgets{
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
		wg.Flex().Rigid(wg.balanceCard()).Fn,
		sp.InputMessage(),
		sp.AddressInput(),
		sp.AmountInput(),
		sp.MessageInput(),
		wg.Flex().
			Flexed(
				1,
				sp.SendButton(),
			).
			Rigid(
				wg.Inset(0.5, gel.EmptySpace(0, 0)).Fn,
			).
			Rigid(
				sp.PasteButton(),
			).
			Rigid(
				wg.Inset(0.5, gel.EmptySpace(0, 0)).Fn,
			).
			Rigid(
				sp.SaveButton(),
			).Fn,
		sp.AddressbookHeader(),
	}
	smallWidgets = append(smallWidgets, sp.GetAddressbookHistoryCards("DocBg")...)
	le := func(gtx l.Context, index int) l.Dimensions {
		return wg.Inset(
			0.25,
			smallWidgets[index],
		).Fn(gtx)
	}
	return wg.lists["send"].
		Vertical().
		Length(len(smallWidgets)).
		ListElement(le).Fn(gtx)
}

func (sp *SendPage) InputMessage() l.Widget {
	return sp.wg.Body2("Enter or paste the details for a payment").Alignment(text.Start).Fn
}

func (sp *SendPage) MediumList(gtx l.Context) l.Dimensions {
	wg := sp.wg
	sendFormWidget := []l.Widget{
		wg.balanceCard(),
		sp.InputMessage(),
		sp.AddressInput(),
		sp.AmountInput(),
		sp.MessageInput(),
		wg.Flex().
			Flexed(
				1,
				sp.SendButton(),
			).
			Rigid(
				wg.Inset(0.5, gel.EmptySpace(0, 0)).Fn,
			).
			Rigid(
				sp.PasteButton(),
			).
			Rigid(
				wg.Inset(0.5, gel.EmptySpace(0, 0)).Fn,
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
		return wg.Inset(
			0.25,
			historyWidget[index],
		).Fn(gtx)
	}
	return wg.Flex().AlignStart().
		Rigid(
			func(gtx l.Context) l.Dimensions {
				gtx.Constraints.Max.X =
					int(wg.TextSize.V * sp.inputWidth)
				// gtx.Constraints.Min.X = int(wg.TextSize.True * sp.inputWidth)

				return wg.VFlex().AlignStart().
					Rigid(
						wg.lists["sendMedium"].
							Vertical().
							Length(len(sendFormWidget)).
							ListElement(sendLE).Fn,
					).Fn(gtx)
			},
		).
		// Rigid(wg.Inset(0.25, gel.EmptySpace(0, 0)).Fn).
		Flexed(
			1,
			wg.VFlex().AlignStart().
				Rigid(
					sp.AddressbookHeader(),
				).
				Rigid(
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
		if wg.inputs["sendAmount"].GetText() == "" || wg.inputs["sendMessage"].GetText() == "" ||
			wg.inputs["sendAddress"].GetText() == "" {
			gtx.Queue = nil
		}
		return wg.ButtonLayout(
			wg.clickables["sendSend"].
				SetClick(
					func() {
						D.Ln("clicked send button")
						go func() {
							if wg.WalletAndClientRunning() {
								var amount float64
								var am amt.Amount
								var e error
								if amount, e = strconv.ParseFloat(
									wg.inputs["sendAmount"].GetText(),
									64,
								); !E.Chk(e) {
									if am, e = amt.NewAmount(amount); E.Chk(e) {
										D.Ln(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", e)
										// todo: indicate this to the user somehow
										return
									}
								} else {
									// todo: indicate this to the user somehow
									return
								}
								var addr btcaddr.Address
								if addr, e = btcaddr.Decode(
									wg.inputs["sendAddress"].GetText(),
									wg.cx.ActiveNet,
								); E.Chk(e) {
									D.Ln(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", e)
									D.Ln("invalid address")
									// TODO: indicate this to the user somehow
									return
								}
								if e = wg.WalletClient.WalletPassphrase(wg.cx.Config.WalletPass.V(), 5); E.Chk(e) {
									D.Ln(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", e)
									return
								}
								var txid *chainhash.Hash
								if txid, e = wg.WalletClient.SendToAddress(addr, am); E.Chk(e) {
									// TODO: indicate send failure to user somehow
									D.Ln(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", e)
									return
								}
								wg.RecentTransactions(10, "recent")
								wg.RecentTransactions(-1, "history")
								wg.Invalidate()
								D.Ln("transaction successful", txid)
								sp.saveForm(txid.String())
								select {
								case <-time.After(time.Second * 5):
								case <-wg.quit:
								}
							}
						}()
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
func (sp *SendPage) saveForm(txid string) {
	wg := sp.wg
	D.Ln("processing form data to save")
	amtS := wg.inputs["sendAmount"].GetText()
	var e error
	var amount float64
	if amount, e = strconv.ParseFloat(amtS, 64); E.Chk(e) {
		return
	}
	if amount == 0 {
		return
	}
	var ua amt.Amount
	if ua, e = amt.NewAmount(amount); E.Chk(e) {
		return
	}
	msg := wg.inputs["sendMessage"].GetText()
	if msg == "" {
		return
	}
	addr := wg.inputs["sendAddress"].GetText()
	var ad btcaddr.Address
	if ad, e = btcaddr.Decode(addr, wg.cx.ActiveNet); E.Chk(e) {
		return
	}
	wg.State.sendAddresses = append(
		wg.State.sendAddresses, AddressEntry{
			Address: ad.EncodeAddress(),
			Label:   msg,
			Amount:  ua,
			Created: time.Now(),
			TxID:    txid,
		},
	)
	// prevent accidental double clicks recording the same entry again
	wg.inputs["sendAmount"].SetText("")
	wg.inputs["sendMessage"].SetText("")
	wg.inputs["sendAddress"].SetText("")
}

func (sp *SendPage) SaveButton() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		wg := sp.wg
		if wg.inputs["sendAmount"].GetText() == "" || wg.inputs["sendMessage"].GetText() == "" ||
			wg.inputs["sendAddress"].GetText() == "" {
			gtx.Queue = nil
		}
		return wg.ButtonLayout(
			wg.clickables["sendSave"].
				SetClick(
					func() { sp.saveForm("") },
				),
		).
			Background("Primary").
			Embed(
				wg.Inset(
					0.5,
					wg.H6("save").Color("Light").Fn,
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
				SetClick(func() { sp.pasteFunction() }),
		).
			Background("Primary").
			Embed(
				wg.Inset(
					0.5,
					wg.H6("paste").Color("Light").Fn,
				).
					Fn,
			).
			Fn(gtx)
	}
}

func (sp *SendPage) pasteFunction() (b bool) {
	wg := sp.wg
	D.Ln("clicked paste button")
	var urn string
	var e error
	if urn, e = clipboard.ReadAll(); E.Chk(e) {
		return
	}
	if !strings.HasPrefix(urn, "parallelcoin:") {
		if e = clipboard.WriteAll(urn); E.Chk(e) {
		}
		return
	}
	split1 := strings.Split(urn, "parallelcoin:")
	split2 := strings.Split(split1[1], "?")
	addr := split2[0]
	var ua btcaddr.Address
	if ua, e = btcaddr.Decode(addr, wg.cx.ActiveNet); E.Chk(e) {
		return
	}
	_ = ua
	b = true
	wg.inputs["sendAddress"].SetText(addr)
	if len(split2) <= 1 {
		return
	}
	split3 := strings.Split(split2[1], "&")
	for i := range split3 {
		var split4 []string
		split4 = strings.Split(split3[i], "=")
		D.Ln(split4)
		if len(split4) > 1 {
			switch split4[0] {
			case "amount":
				wg.inputs["sendAmount"].SetText(split4[1])
				// D.Ln("############ amount", split4[1])
			case "message", "label":
				msg := split4[i]
				if len(msg) > 64 {
					msg = msg[:64]
				}
				wg.inputs["sendMessage"].SetText(msg)
				// D.Ln("############ message", split4[1])
			}
		}
	}
	return
}

func (sp *SendPage) AddressbookHeader() l.Widget {
	wg := sp.wg
	return wg.Flex().AlignStart().
		Rigid(
			wg.Inset(
				0.25,
				wg.H5("Send Address Book").Fn,
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
								wg.State.sendAddresses[i].Label,
							)
							D.Ln("clicked send address list item", j)
							if e := clipboard.WriteAll(sendText); E.Chk(e) {
							}
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
									wg.Inset(
										0.25,
										wg.Body1(wg.State.sendAddresses[i].Label).MaxLines(1).Fn,
									).Fn,
								).
								Rigid(
									gel.If(
										wg.State.sendAddresses[i].TxID != "",
										func(ctx l.Context) l.Dimensions {
											for j := range wg.txHistoryList {
												if wg.txHistoryList[j].TxID == wg.State.sendAddresses[i].TxID {
													return wg.Flex().Flexed(
														1,
														wg.VFlex().
															Rigid(
																wg.Flex().Flexed(
																	1,
																	wg.Caption(wg.State.sendAddresses[i].TxID).MaxLines(1).Fn,
																).Fn,
															).
															Rigid(
																wg.Body1(
																	fmt.Sprint(
																		"Confirmations: ",
																		wg.txHistoryList[j].Confirmations,
																	),
																).Fn,
															).Fn,
													).Fn(gtx)
												}
											}
											return func(ctx l.Context) l.Dimensions { return l.Dimensions{} }(gtx)
										},
										func(ctx l.Context) l.Dimensions { return l.Dimensions{} },
									),
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
