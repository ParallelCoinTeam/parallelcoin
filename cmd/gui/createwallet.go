package gui

import (
	"encoding/hex"
	"time"

	l "gioui.org/layout"

	"github.com/p9c/pod/pkg/util/hdkeychain"
	"github.com/p9c/pod/pkg/wallet"
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
									wg.th.Inset(0.25,
										wg.passwords["passEditor"].Fn,
									).Fn,
								).
								Rigid(
									wg.th.Inset(0.25,
										wg.passwords["confirmPassEditor"].Fn,
									).Fn,
								).
								Rigid(
									wg.th.Inset(0.25,
										wg.inputs["walletSeed"].Fn,
									).Fn,
								).
								Rigid(
									wg.th.Inset(0.25,
										wg.passwords["publicPassEditor"].Fn,
									).Fn,
								).
								Rigid(
									wg.th.Inset(0.25,
										func(gtx l.Context) l.Dimensions {
											gtx.Constraints.Min.X = int(wg.th.TextSize.Scale(16).V)
											return wg.CheckBox(wg.bools["testnet"].SetOnChange(func(b bool) {
												Debug("testnet on?", b)
											})).
												IconColor("Primary").
												TextColor("DocText").
												Text("Use testnet?").
												Fn(gtx)
										},
									).Fn,
								).
								Rigid(
									wg.th.Body1("your seed").
										Color("PanelText").
										Fn,
								).
								Rigid(
									func(gtx l.Context) l.Dimensions {
										gtx.Constraints.Max.X = int(wg.TextSize.Scale(22).V)
										return wg.th.Caption(wg.inputs["walletSeed"].GetText()).
											Font("go regular").
											TextScale(0.66).
											Fn(gtx)
									},
								).
								Rigid(
									wg.th.Inset(0.5,
										func(gtx l.Context) l.Dimensions {
											gtx.Constraints.Max.X = int(wg.th.TextSize.Scale(36).V)
											gtx.Constraints.Min.X = int(wg.th.TextSize.Scale(16).V)
											return wg.CheckBox(wg.bools["ihaveread"].SetOnChange(func(b bool) {
												Debug("confirmed read", b)
											})).
												IconColor("Primary").
												TextColor("DocText").
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
											len(wg.passwords["passEditor"].GetPassword()) < 8 ||
											wg.passwords["passEditor"].GetPassword() !=
												wg.passwords["confirmPassEditor"].GetPassword() ||
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
														go func() {
															dbDir := *wg.cx.Config.WalletFile
															loader := wallet.NewLoader(wg.cx.ActiveNet, dbDir, 250)
															seed, _ := hex.DecodeString(wg.inputs["walletSeed"].GetText())
															w, err := loader.CreateNewWallet(
																[]byte(wg.passwords["publicPassEditor"].GetPassword()),
																[]byte(wg.passwords["passEditor"].GetPassword()),
																seed,
																time.Now(),
																false,
																wg.cx.Config,
															)
															if Check(err) {
																panic(err)
															}
															w.Manager.Close()
															wg.noWallet = false
															Debug("starting up shell")
															wg.running = false
															wg.ShellRunCommandChan <- "run"
														}()
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
