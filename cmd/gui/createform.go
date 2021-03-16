package gui

import (
	"encoding/hex"
	l "gioui.org/layout"
	"gioui.org/text"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/gui"
	p9icons "github.com/p9c/pod/pkg/gui/ico/svg"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"strings"
)

func (wg *WalletGUI) centered(w l.Widget) l.Widget {
	return wg.Flex().
		Flexed(0.5, gui.EmptyMaxWidth()).
		Rigid(
			wg.VFlex().
				AlignMiddle().
				Rigid(
					w,
				).
				Fn,
		).
		Flexed(0.5, gui.EmptyMaxWidth()).
		Fn
}

func (wg *WalletGUI) createWalletFormWidgets() (out []l.Widget) {
	out = append(
		out,
		wg.centered(
			wg.Icon().
				Scale(gui.Scales["H2"]).
				Color("DocText").
				Src(&p9icons.ParallelCoin).Fn,
		),
		wg.centered(
			wg.H4("create new wallet").
				Color("PanelText").
				Fn,
		),
		wg.H6("wallet password").
			Color("PanelText").
			Fn,
		// wg.Caption("(min 8 characters)").
		// 	Color("PanelText").
		// 	Fn,
		wg.passwords["passEditor"].
			Fn,
		wg.passwords["confirmPassEditor"].
			Fn,
		wg.Flex().SpaceBetween().
			Rigid(
				wg.H6("wallet seed").
					Color("PanelText").
					Fn,
			).
			Fn,
		func(gtx l.Context) (o l.Dimensions) {
			verifyState := wg.Button(
				wg.clickables["createVerify"].SetClick(
					func() {
						wg.createVerifying = false
					},
				),
			).Text("back").Fn
			if wg.createWords == wg.createMatch {
				verifyState = wg.Inset(0.25, wg.Body1("match").Color("Success").Fn).Fn
			}
			if wg.createVerifying {
				return wg.Flex().
					Flexed(
						1,
						wg.inputs["walletWords"].Fn,
					).
					Rigid(
						verifyState,
					).
					Rigid(
						wg.IconButton(
							wg.clickables["createShuffle"].SetClick(
								func() {
									wg.ShuffleSeed()
									wg.inputs["walletWords"].SetText("") // wg.createWords)
									wg.createVerifying = false
								},
							),
						).Background("Transparent").
							ButtonInset(0).
							// Scale(gui.Scales["Caption"]).
							Icon(
								wg.Icon().
									Scale(gui.Scales["H5"]).
									Color("DocText").
									Src(&icons.NavigationRefresh),
							).Fn,
					).
					Fn(gtx)
			} else {
				
				// var b []byte
				// if b, e = hex.DecodeString(wg.inputs["walletSeed"].GetText()); err.Chk(e) {
				// 	return
				// }
				col := "DocText"
				if wg.createWords == wg.createMatch {
					col = "Success"
				}
				return wg.Flex().
					Flexed(
						1,
						wg.ButtonLayout(
							wg.clickables["createVerify"].SetClick(
								func() {
									wg.createVerifying = true
								},
							),
						).Background("Transparent").Embed(
							wg.VFlex().
								Rigid(
									wg.Caption("Write the following words down, then click to re-enter and verify transcription").
										Color("PanelText").
										Fn,
								).
								Rigid(
									wg.Flex().Flexed(
										1,
										wg.Body1(wg.showWords).Alignment(text.Middle).Color(col).Fn,
									).Fn,
								).Fn,
						).Fn,
					).
					Rigid(
						wg.IconButton(
							wg.clickables["createShuffle"].SetClick(
								func() {
									wg.ShuffleSeed()
									wg.inputs["walletWords"].SetText("") // wg.createWords)
								},
							),
						).Background("Transparent").
							ButtonInset(0).
							// Scale(gui.Scales["Caption"]).
							Icon(
								wg.Icon().
									Scale(gui.Scales["H5"]).
									Color("DocText").
									Src(&icons.NavigationRefresh),
							).Fn,
					).
					Fn(gtx)
			}
		},
		// wg.inputs["walletSeed"].
		// 	Fn,
		// func(gtx l.Context) l.Dimensions {
		// 	// gtx.Constraints.Max.X = int(wg.TextSize.Scale(22).V)
		// 	return wg.Inset(
		// 		0.25,
		// 		wg.Caption(wg.inputs["walletSeed"].GetText()).
		// 			Font("go regular").
		// 			// TextScale(0.66).
		// 			Fn,
		// 	).Fn(gtx)
		// },
		wg.Flex().
			Rigid(
				func(gtx l.Context) l.Dimensions {
					return wg.CheckBox(
						wg.bools["testnet"].SetOnChange(
							func(b bool) {
								if !b {
									wg.bools["solo"].Value(false)
									wg.bools["lan"].Value(false)
									*wg.cx.Config.Solo, *wg.cx.Config.LAN = false, false
									wg.Invalidate()
								}
								wg.createWalletTestnetToggle(b)
							},
						),
					).
						IconColor("Primary").
						TextColor("DocText").
						Text("Use Testnet").
						Fn(gtx)
				},
			).
			Rigid(
				func(gtx l.Context) l.Dimensions {
					checkColor, textColor := "Primary", "DocText"
					if !wg.bools["testnet"].GetValue() {
						gtx = gtx.Disabled()
						checkColor, textColor = "scrim", "scrim"
					}
					return wg.CheckBox(
						wg.bools["lan"].SetOnChange(
							func(b bool) {
								dbg.Ln("lan now set to", b)
								*wg.cx.Config.LAN = b
								if *wg.cx.Config.Solo {
									*wg.cx.Config.Solo = false
									*wg.cx.Config.MinerPass = "pa55word"
									wg.bools["solo"].Value(false)
									wg.Invalidate()
								}
								save.Pod(wg.cx.Config)
							},
							// wg.createWalletTestnetToggle,
						),
					).
						IconColor(checkColor).
						TextColor(textColor).
						Text("LAN only").
						Fn(gtx)
				},
			).
			Rigid(
				func(gtx l.Context) l.Dimensions {
					checkColor, textColor := "Primary", "DocText"
					if !wg.bools["testnet"].GetValue() {
						gtx = gtx.Disabled()
						checkColor, textColor = "scrim", "scrim"
					}
					return wg.CheckBox(
						wg.bools["solo"].SetOnChange(
							func(b bool) {
								dbg.Ln("solo now set to", b)
								*wg.cx.Config.Solo = b
								if *wg.cx.Config.LAN {
									*wg.cx.Config.LAN = false
									wg.bools["lan"].Value(false)
									wg.Invalidate()
								}
								save.Pod(wg.cx.Config)
							},
						),
					).
						IconColor(checkColor).
						TextColor(textColor).
						Text("Solo (mine without peers)").
						Fn(gtx)
				},
			).
			Rigid(wg.Inset(0.25, gui.EmptySpace(0, 0)).Fn).
			Rigid(
				func(gtx l.Context) l.Dimensions {
					if !wg.bools["testnet"].GetValue() {
						gtx = gtx.Disabled()
					}
					return wg.Button(
						wg.clickables["genesis"].SetClick(
							func() {
								seedString := "f4d2c4c542bb52512ed9e6bbfa2d000e576a0c8b4ebd1acafd7efa37247366bc"
								var e error
								if wg.createSeed, e = hex.DecodeString(seedString); ftl.Chk(e) {
									panic(e)
								}
								var wk string
								if wk, e = bip39.NewMnemonic(wg.createSeed); err.Chk(e) {
									panic(e)
								}
								wks := strings.Split(wk, " ")
								var out string
								for i := 0; i < 24; i += 4 {
									out += strings.Join(wks[i:i+4], " ")
									if i+4 < 24 {
										out += "\n"
									}
								}
								wg.showWords = out
								wg.createWords = wk
								wg.createMatch = wk
								wg.inputs["walletWords"].SetText(wk)
								wg.createVerifying = true
							},
						),
					).Text("genesis").Fn(gtx)
				},
			).
			Rigid(wg.Inset(0.25, gui.EmptySpace(0, 0)).Fn).
			Rigid(
				func(gtx l.Context) l.Dimensions {
					if !wg.bools["testnet"].GetValue() {
						gtx = gtx.Disabled()
					}
					return wg.Button(
						wg.clickables["autofill"].SetClick(
							func() {
								wk := wg.createWords
								wg.createMatch = wk
								wg.inputs["walletWords"].SetText(wk)
								wg.createVerifying = true
							},
						),
					).Text("autofill").Fn(gtx)
				},
			).
			Fn,
		func(gtx l.Context) l.Dimensions {
			return wg.CheckBox(
				wg.bools["ihaveread"].SetOnChange(
					func(b bool) {
						dbg.Ln("confirmed read", b)
						// if the password has been entered, we need to copy it to the variable
						if wg.createWalletPasswordsMatch() {
							wg.cx.Config.Lock()
							*wg.cx.Config.WalletPass = wg.passwords["confirmPassEditor"].GetPassword()
							wg.cx.Config.Unlock()
						}
					},
				),
			).
				IconColor("Primary").
				TextColor("DocText").
				Text(
					"I have stored the seed and password safely " +
						"and understand it cannot be recovered",
				).
				Fn(gtx)
		},
	)
	return
}
