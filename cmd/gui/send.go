package gui

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	
	l "gioui.org/layout"
	"gioui.org/text"
	"github.com/atotto/clipboard"
	
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/gui"
	"github.com/p9c/pod/pkg/util"
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
	wg.inputs["sendAddress"].SetPasteFunc = sp.pasteFunction
	wg.inputs["sendAmount"].SetPasteFunc = sp.pasteFunction
	wg.inputs["sendMessage"].SetPasteFunc = sp.pasteFunction
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
				return wg.lists["sendMedium"].
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
		if wg.inputs["sendAmount"].GetText() == "" || wg.inputs["sendMessage"].GetText() == "" ||
			wg.inputs["sendAddress"].GetText() == "" {
			gtx.Queue = nil
		}
		return wg.ButtonLayout(
			wg.clickables["sendSend"].
				SetClick(
					func() {
						Debug("clicked send button")
						go func() {
							if wg.WalletAndClientRunning() {
								var amt float64
								var am util.Amount
								var err error
								if amt, err = strconv.ParseFloat(
									wg.inputs["sendAmount"].GetText(),
									64,
								); !Check(err) {
									if am, err = util.NewAmount(amt); Check(err) {
										// todo: indicate this to the user somehow
										return
									}
								} else {
									// todo: indicate this to the user somehow
									return
								}
								var addr util.Address
								if addr, err = util.DecodeAddress(wg.inputs["sendAddress"].GetText(),
									wg.cx.ActiveNet); Check(err) {
									Debug("invalid address")
									// TODO: indicate this to the user somehow
									return
								}
								if err= wg.WalletClient.WalletPassphrase(*wg.cx.Config.WalletPass, 5); Check(err){
									return
								}
								var txid *chainhash.Hash
								if txid, err = wg.WalletClient.SendToAddress(addr, am); Check(err) {
									// TODO: indicate send failure to user somehow
									Debug("******")
									return
								}
								wg.RecentTransactions(10, "recent")
								wg.RecentTransactions(-1, "history")
								Debug("transaction successful", txid)
								// prevent accidental double clicks recording the same entry again
								wg.inputs["sendAmount"].SetText("")
								wg.inputs["sendMessage"].SetText("")
								wg.inputs["sendAddress"].SetText("")
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
					func() {
						Debug("clicked save button")
						amtS := wg.inputs["sendAmount"].GetText()
						var err error
						var amt float64
						if amt, err = strconv.ParseFloat(amtS, 64); Check(err) {
							return
						}
						if amt == 0 {
							return
						}
						var ua util.Amount
						if ua, err = util.NewAmount(amt); Check(err) {
							return
						}
						msg := wg.inputs["sendMessage"].GetText()
						if msg == "" {
							return
						}
						addr := wg.inputs["sendAddress"].GetText()
						var ad util.Address
						if ad, err = util.DecodeAddress(addr, wg.cx.ActiveNet); Check(err) {
							return
						}
						wg.State.sendAddresses = append(wg.State.sendAddresses, AddressEntry{
							Address: ad.EncodeAddress(),
							Label:   msg,
							Amount:  ua,
							Created: time.Now(),
						})
						// prevent accidental double clicks recording the same entry again
						wg.inputs["sendAmount"].SetText("")
						wg.inputs["sendMessage"].SetText("")
						wg.inputs["sendAddress"].SetText("")
					},
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
				SetClick(func() { sp.pasteFunction() })).
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
	Debug("clicked paste button")
	var urn string
	var err error
	if urn, err = clipboard.ReadAll(); Check(err) {
		return
	}
	if !strings.HasPrefix(urn, "parallelcoin:") {
		if err = clipboard.WriteAll(urn); Check(err) {
		}
		return
	}
	split1 := strings.Split(urn, "parallelcoin:")
	split2 := strings.Split(split1[1], "?")
	addr := split2[0]
	var ua util.Address
	if ua, err = util.DecodeAddress(addr, wg.cx.ActiveNet); Check(err) {
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
		Debug(split4)
		if len(split4) > 1 {
			switch split4[0] {
			case "amount":
				wg.inputs["sendAmount"].SetText(split4[1])
				// Debug("############ amount", split4[1])
			case "message", "label":
				msg := split4[i]
				if len(msg) > 64 {
					msg = msg[:64]
				}
				wg.inputs["sendMessage"].SetText(msg)
				// Debug("############ message", split4[1])
			}
		}
	}
	return
}

func (sp *SendPage) AddressbookHeader() l.Widget {
	wg := sp.wg
	return wg.Flex().Flexed(
		1,
		wg.Inset(
			0.25,
			wg.H6("Address Book").Alignment(text.Middle).Fn,
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
									wg.Caption(wg.State.sendAddresses[i].Label).MaxLines(1).Fn,
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
