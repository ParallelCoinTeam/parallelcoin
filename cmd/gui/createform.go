package gui

import (
	"encoding/hex"
	"strings"

	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/text"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/gel"
	"github.com/p9c/matrjoska/pkg/p9icons"
)

func (wg *WalletGUI) centered(w l.Widget) l.Widget {
	return wg.Flex().
		Flexed(0.5, gel.EmptyMaxWidth()).
		Rigid(
			wg.VFlex().
				AlignMiddle().
				Rigid(
					w,
				).
				Fn,
		).
		Flexed(0.5, gel.EmptyMaxWidth()).
		Fn
}

func (wg *WalletGUI) cwfLogoHeader() l.Widget {
	return wg.centered(
		wg.Icon().
			Scale(gel.Scales["H2"]).
			Color("DocText").
			Src(&p9icons.ParallelCoin).Fn,
	)
}

func (wg *WalletGUI) cwfHeader() l.Widget {
	return wg.centered(
		wg.H4("create new wallet").
			Color("PanelText").
			Fn,
	)
}

func (wg *WalletGUI) cwfPasswordHeader() l.Widget {
	return wg.H5("password").
		Color("PanelText").
		Fn
}

func (wg *WalletGUI) cwfShuffleButton() l.Widget {
	return wg.ButtonLayout(
		wg.clickables["createShuffle"].SetClick(
			func() {
				wg.ShuffleSeed()
				wg.inputs["walletWords"].SetText("") // wg.createWords)
				wg.inputs["walletWords"].SetText("") // wg.createWords)
				wg.restoring = false
				wg.createVerifying = false
			},
		),
	).
		CornerRadius(0).
		Corners(0).
		Background("Primary").
		Embed(
			wg.Inset(
				0.25,
				wg.Flex().AlignMiddle().
					Rigid(
						wg.Icon().
							Scale(
								gel.Scales["H6"],
							).
							Color("DocText").
							Src(
								&icons.NavigationRefresh,
							).Fn,
					).
					Rigid(
						wg.Body1("new").Color("DocText").Fn,
					).
					Fn,
			).Fn,
		).Fn
}

func (wg *WalletGUI) cwfRestoreButton() l.Widget {
	return wg.ButtonLayout(
		wg.clickables["createRestore"].SetClick(
			func() {
				D.Ln("clicked restore button")
				if !wg.restoring {
					wg.inputs["walletRestore"].SetText("")
					wg.createMatch = ""
					wg.restoring = true
					wg.createVerifying = false
				} else {
					wg.createMatch = ""
					wg.restoring = false
					wg.createVerifying = false
				}
			},
		),
	).
		CornerRadius(0).
		Corners(0).
		Background("Primary").
		Embed(
			wg.Inset(
				0.25,
				wg.Flex().AlignMiddle().
					Rigid(
						wg.Icon().
							Scale(
								gel.Scales["H6"],
							).
							Color("DocText").
							Src(
								&icons.ActionRestore,
							).Fn,
					).
					Rigid(
						wg.Body1("restore").Color("DocText").Fn,
					).
					Fn,
			).Fn,
		).Fn
}

func (wg *WalletGUI) cwfSetGenesis() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		if !wg.bools["testnet"].GetValue() {
			return l.Dimensions{}
		} else {
			return wg.ButtonLayout(
				wg.clickables["genesis"].SetClick(
					func() {
						seedString := "f4d2c4c542bb52512ed9e6bbfa2d000e576a0c8b4ebd1acafd7efa37247366bc"
						var e error
						if wg.createSeed, e = hex.DecodeString(seedString); F.Chk(e) {
							panic(e)
						}
						var wk string
						if wk, e = bip39.NewMnemonic(wg.createSeed); E.Chk(e) {
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
			).
				CornerRadius(0).
				Corners(0).
				Background("Primary").
				Embed(
					wg.Inset(
						0.25,
						wg.Flex().AlignMiddle().
							Rigid(
								wg.Icon().
									Scale(
										gel.Scales["H6"],
									).
									Color("DocText").
									Src(
										&icons.ActionOpenInNew,
									).Fn,
							).
							Rigid(
								wg.Body1("genesis").Color("DocText").Fn,
							).
							Fn,
					).Fn,
				).Fn(gtx)
		}
	}
}

func (wg *WalletGUI) cwfSetAutofill() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		if !wg.bools["testnet"].GetValue() {
			return l.Dimensions{}
		} else {
			return wg.ButtonLayout(
				wg.clickables["autofill"].SetClick(
					func() {
						wk := wg.createWords
						wg.createMatch = wk
						wg.inputs["walletWords"].SetText(wk)
						wg.createVerifying = true
					},
				),
			).
				CornerRadius(0).
				Corners(0).
				Background("Primary").
				Embed(
					wg.Inset(
						0.25,
						wg.Flex().AlignMiddle().
							Rigid(
								wg.Icon().
									Scale(
										gel.Scales["H6"],
									).
									Color("DocText").
									Src(
										&icons.ActionOpenInNew,
									).Fn,
							).
							Rigid(
								wg.Body1("autofill").Color("DocText").Fn,
							).
							Fn,
					).Fn,
				).Fn(gtx)
		}
	}
}

func (wg *WalletGUI) cwfSeedHeader() l.Widget {
	return wg.Flex(). //AlignMiddle().
		Rigid(
			wg.Inset(
				0.25,
				wg.H5("seed").
					Color("PanelText").
					Fn,
			).Fn,
		).
		Rigid(wg.Inset(0.25, gel.EmptySpace(0, 0)).Fn).
		Rigid(wg.cwfShuffleButton()).
		Rigid(wg.Inset(0.25, gel.EmptySpace(0, 0)).Fn).
		Rigid(wg.cwfRestoreButton()).
		Rigid(wg.Inset(0.25, gel.EmptySpace(0, 0)).Fn).
		Rigid(wg.cwfSetGenesis()).
		Rigid(wg.Inset(0.25, gel.EmptySpace(0, 0)).Fn).
		Rigid(wg.cwfSetAutofill()).
		Fn
}

func (wg *WalletGUI) cfwWords() (w l.Widget) {
	if !wg.createVerifying {
		col := "DocText"
		if wg.createWords == wg.createMatch {
			col = "Success"
		}
		return wg.Flex().
			Rigid(
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
			Fn
	}
	return nil
}

func (wg *WalletGUI) cfwWordsVerify() (w l.Widget) {
	if wg.createVerifying {
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
		return wg.Flex().
			Rigid(
				verifyState,
			).
			Rigid(
				wg.inputs["walletWords"].Fn,
			).
			Fn
	}
	return nil
}

func (wg *WalletGUI) cfwRestore() (w l.Widget) {
	w = func(l.Context) l.Dimensions {
		return l.Dimensions{}
	}
	if wg.restoring {
		// restoreState := wg.Button(
		// 	wg.clickables["createRestore"].SetClick(
		// 		func() {
		// 			wg.restoring = false
		// 		},
		// 	),
		// ).Text("back").Fn
		if wg.createWords == wg.createMatch {
			w = wg.Flex().AlignMiddle().
				Rigid(
					wg.Inset(0.25, wg.H5("valid").Color("Success").Fn).Fn,
				).Fn
		}
		return wg.Flex().
			Rigid(
				w,
			).
			Rigid(
				wg.inputs["walletRestore"].Fn,
			).
			Fn
	}
	return
}

func (wg *WalletGUI) cwfTestnetSettings() (out l.Widget) {
	return wg.Flex().
		Rigid(
			func(gtx l.Context) l.Dimensions {
				return wg.CheckBox(
					wg.bools["testnet"].SetOnChange(
						func(b bool) {
							if !b {
								wg.bools["solo"].Value(false)
								wg.bools["lan"].Value(false)
								// wg.cx.Config.MulticastPass.Set("pa55word")
								wg.cx.Config.Solo.F()
								wg.cx.Config.LAN.F()
								wg.ShuffleSeed()
								wg.createVerifying = false
								wg.inputs["walletWords"].SetText("")
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
							D.Ln("lan now set to", b)
							wg.cx.Config.LAN.Set(b)
							if b && wg.cx.Config.Solo.True() {
								wg.cx.Config.Solo.F()
								wg.cx.Config.DisableDNSSeed.T()
								wg.cx.Config.AutoListen.F()
								wg.bools["solo"].Value(false)
								// wg.cx.Config.MulticastPass.Set("pa55word")
								wg.Invalidate()
							} else {
								wg.cx.Config.Solo.F()
								wg.cx.Config.DisableDNSSeed.F()
								// wg.cx.Config.MulticastPass.Set("pa55word")
								wg.cx.Config.AutoListen.T()
							}
							_ = wg.cx.Config.WriteToFile(wg.cx.Config.ConfigFile.V())
						},
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
							D.Ln("solo now set to", b)
							wg.cx.Config.Solo.Set(b)
							if b && wg.cx.Config.LAN.True() {
								wg.cx.Config.LAN.F()
								wg.cx.Config.DisableDNSSeed.T()
								wg.cx.Config.AutoListen.F()
								// wg.cx.Config.MulticastPass.Set("pa55word")
								wg.bools["lan"].Value(false)
								wg.Invalidate()
							} else {
								wg.cx.Config.LAN.F()
								wg.cx.Config.DisableDNSSeed.F()
								// wg.cx.Config.MulticastPass.Set("pa55word")
								wg.cx.Config.AutoListen.T()
							}
							_ = wg.cx.Config.WriteToFile(wg.cx.Config.ConfigFile.V())
						},
					),
				).
					IconColor(checkColor).
					TextColor(textColor).
					Text("Solo (mine without peers)").
					Fn(gtx)
			},
		).
		Fn
}

func (wg *WalletGUI) cwfConfirmation() (out l.Widget) {
	return wg.CheckBox(
		wg.bools["ihaveread"].SetOnChange(
			func(b bool) {
				D.Ln("confirmed read", b)
				// if the password has been entered, we need to copy it to the variable
				if wg.createWalletPasswordsMatch() {
					wg.cx.Config.WalletPass.Set(wg.passwords["confirmPassEditor"].GetPassword())
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
		Fn
}

func (wg *WalletGUI) createWalletFormWidgets() (out []l.Widget) {
	out = append(
		out,
		wg.cwfLogoHeader(),
		wg.cwfHeader(),
		wg.cwfPasswordHeader(),
		wg.passwords["passEditor"].
			Fn,
		wg.passwords["confirmPassEditor"].
			Fn,
		wg.cwfSeedHeader(),
	)
	if wg.createVerifying {
		out = append(
			out, wg.cfwWordsVerify(),
		)
	} else if wg.restoring {
		out = append(out, wg.cfwRestore())
	} else {
		out = append(out, wg.cfwWords())
	}
	out = append(
		out,
		wg.cwfTestnetSettings(),
		wg.cwfConfirmation(),
	)
	return
}
