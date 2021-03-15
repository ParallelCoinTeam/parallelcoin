package gui

import (
	l "gioui.org/layout"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/gui"
	p9icons "github.com/p9c/pod/pkg/gui/ico/svg"
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
		wg.H6("wallet seed").
			Color("PanelText").
			Fn,
		// wg.Caption("(hex)").
		// 	Color("PanelText").
		// 	Fn,
		wg.inputs["walletSeed"].
			Fn,
		func(gtx l.Context) l.Dimensions {
			// gtx.Constraints.Max.X = int(wg.TextSize.Scale(22).V)
			return wg.Inset(
				0.25,
				wg.Caption(wg.inputs["walletSeed"].GetText()).
					Font("go regular").
					// TextScale(0.66).
					Fn,
			).Fn(gtx)
		},
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
