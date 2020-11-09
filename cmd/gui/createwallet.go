package gui

import (
	"encoding/hex"

	l "gioui.org/layout"

	"github.com/p9c/pod/pkg/util/hdkeychain"
)

func (wg *WalletGUI) WalletPage(gtx l.Context) l.Dimensions {
	return wg.th.
		Fill("PanelBg",
			wg.th.Flex().SpaceAround().AlignMiddle().
				Rigid(
					wg.th.
						VFlex().AlignMiddle().SpaceAround().
						Rigid(
							wg.th.
								VFlex().SpaceAround().AlignMiddle().
								Rigid(
									wg.th.H4("create new wallet").
										Color("PanelText").
										// Alignment(text.Middle).
										Fn,
								).
								Rigid(
									wg.th.Inset(0.5,
										wg.passwords["passEditor"].Fn,
									).Fn,
								).
								Rigid(
									wg.th.Inset(0.5,
										wg.passwords["confirmPassEditor"].Fn,
									).Fn,
								).
								Rigid(
									wg.th.Inset(0.5,
										wg.inputs["walletSeed"].Fn,
									).Fn,
								).
								Rigid(
									wg.th.Inset(0.5,
										wg.passwords["publicPassEditor"].Fn,
									).Fn,
								).
								Rigid(
									wg.th.Inset(0.5,
										func(gtx l.Context) l.Dimensions {
											gtx.Constraints.Min.X = int(wg.th.TextSize.Scale(16).V)
											return wg.CheckBox(wg.bools["testnet"].SetOnChange(func(b bool) {
												Debug("testnet on?", b)
											})).
												IconColor("Primary").
												TextColor("DocText").
												// IconScale(0.1).
												Text("Use testnet?").
												Fn(gtx)
										},
									).Fn,
								).
								Rigid(
									wg.th.Inset(0.5,
										func(gtx l.Context) l.Dimensions {
											gtx.Constraints.Min.X = int(wg.th.TextSize.Scale(16).V)
											return wg.CheckBox(wg.bools["ihaveread"].SetOnChange(func(b bool) {
												Debug("confirmed read", b)
											})).
												IconColor("Primary").
												TextColor("DocText").
												// IconScale(0.1).
												Text("I have stored the seed and password safely " +
													"and understand it cannot be recovered").
												Fn(gtx)
										},
									).Fn,
								).
								Rigid(
									func(gtx l.Context) l.Dimensions {
										var b []byte
										var err error
										seedValid := true
										if b, err = hex.DecodeString(wg.inputs["walletSeed"].GetText()); Check(err) {
											seedValid = false
										} else if len(b) != 0 && len(b) < hdkeychain.MinSeedBytes ||
											len(b) > hdkeychain.MaxSeedBytes {
											seedValid = false
										}
										if wg.passwords["passEditor"].GetPassword() == "" ||
											wg.passwords["confirmPassEditor"].GetPassword() == "" ||
											wg.passwords["passEditor"].GetPassword() != wg.passwords["confirmPassEditor"].GetPassword() ||
											!seedValid ||
											!wg.bools["ihaveread"].GetValue() {
											gtx = gtx.Disabled()
										}
										return wg.th.Flex().
											Rigid(
												wg.th.Button(wg.clickables["createWallet"]).
													Background("Primary").
													Color("Light").
													SetClick(func() {
														Debug("clicked submit wallet")
													}).
													CornerRadius(0).
													Inset(0.5).
													Text("create wallet").
													Fn,
											).
											Fn(gtx)
									},
								).
								Fn,
						).
						Fn,
				).
				Fn,
		).
		Fn(gtx)
}
