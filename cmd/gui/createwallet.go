package gui

import (
	"encoding/hex"
	"os"
	"time"
	
	l "gioui.org/layout"
	
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/chain/mining/addresses"
	"github.com/p9c/pod/pkg/gui"
	"github.com/p9c/pod/pkg/util/hdkeychain"
	"github.com/p9c/pod/pkg/wallet"
)

const slash = string(os.PathSeparator)

func (wg *WalletGUI) CreateWalletPage(gtx l.Context) l.Dimensions {
	return wg.Fill("PanelBg", l.Center, 0, 0, wg.Inset(
		0.5,
		wg.Flex().
			SpaceAround().
			Flexed(0.5, gui.EmptyMaxHeight()).
			Rigid(
				func(gtx l.Context) l.Dimensions {
					return wg.VFlex().
						AlignMiddle().
						SpaceSides().
						Rigid(
							wg.H4("create new wallet").
								Color("PanelText").
								Fn,
						).
						Rigid(
							wg.Inset(
								0.25,
								wg.passwords["passEditor"].Fn,
							).
								Fn,
						).
						Rigid(
							wg.Inset(
								0.25,
								wg.passwords["confirmPassEditor"].Fn,
							).
								Fn,
						).
						Rigid(
							wg.Inset(
								0.25,
								wg.inputs["walletSeed"].Fn,
							).
								Fn,
						).
						Rigid(
							wg.Inset(
								0.25,
								func(gtx l.Context) l.Dimensions {
									// gtx.Constraints.Min.X = int(wg.TextSize.Scale(16).V)
									return wg.CheckBox(
										wg.bools["testnet"].SetOnChange(
											func(b bool) {
												go func() {
													Debug("testnet on?", b)
													// if the password has been entered, we need to copy it to the variable
													if wg.passwords["passEditor"].GetPassword() != "" ||
														wg.passwords["confirmPassEditor"].GetPassword() != "" ||
														len(wg.passwords["passEditor"].GetPassword()) >= 8 ||
														wg.passwords["passEditor"].GetPassword() ==
															wg.passwords["confirmPassEditor"].GetPassword() {
														*wg.cx.Config.WalletPass = wg.passwords["confirmPassEditor"].GetPassword()
														Debug("wallet pass", *wg.cx.Config.WalletPass)
													}
													if b {
														wg.cx.ActiveNet = &netparams.TestNet3Params
														fork.IsTestnet = true
													} else {
														wg.cx.ActiveNet = &netparams.MainNetParams
														fork.IsTestnet = false
													}
													Info("activenet:", wg.cx.ActiveNet.Name)
													*wg.cx.Config.Network = wg.cx.ActiveNet.Name
													if wg.cx.ActiveNet.Name == "testnet" {
														// TODO: obviously when we get to starting testnets this should not be done
														*wg.cx.Config.LAN = true  // mines without peer outside lan
														*wg.cx.Config.Solo = true // mines without peers
													}
													save.Pod(wg.cx.Config)
												}()
											},
										),
									).
										IconColor("Primary").
										TextColor("DocText").
										Text("Use testnet?").
										Fn(gtx)
								},
							).Fn,
						).
						Rigid(
							wg.Body1("your seed").
								Color("PanelText").
								Fn,
						).
						Rigid(
							func(gtx l.Context) l.Dimensions {
								gtx.Constraints.Max.X = int(wg.TextSize.Scale(22).V)
								return wg.Caption(wg.inputs["walletSeed"].GetText()).
									Font("go regular").
									TextScale(0.66).
									Fn(gtx)
							},
						).
						Rigid(
							wg.Inset(
								0.5,
								func(gtx l.Context) l.Dimensions {
									gtx.Constraints.Max.X = int(wg.TextSize.Scale(32).V)
									gtx.Constraints.Min.X = int(wg.TextSize.Scale(16).V)
									return wg.CheckBox(
										wg.bools["ihaveread"].SetOnChange(
											func(b bool) {
												Debug("confirmed read", b)
												// if the password has been entered, we need to copy it to the variable
												if wg.passwords["passEditor"].GetPassword() != "" ||
													wg.passwords["confirmPassEditor"].GetPassword() != "" ||
													len(wg.passwords["passEditor"].GetPassword()) >= 8 ||
													wg.passwords["passEditor"].GetPassword() ==
														wg.passwords["confirmPassEditor"].GetPassword() {
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
								return wg.Flex().
									Rigid(
										wg.Button(wg.clickables["createWallet"]).
											Background("Primary").
											Color("Light").
											SetClick(
												func() {
													go func() {
														// wg.NodeRunCommandChan <- "stop"
														Debug("clicked submit wallet")
														*wg.cx.Config.WalletFile = *wg.cx.Config.DataDir +
															string(os.PathSeparator) + wg.cx.ActiveNet.Name +
															string(os.PathSeparator) + wallet.DbName
														dbDir := *wg.cx.Config.WalletFile
														loader := wallet.NewLoader(wg.cx.ActiveNet, dbDir, 250)
														seed, _ := hex.DecodeString(wg.inputs["walletSeed"].GetText())
														pass := []byte(wg.passwords["passEditor"].GetPassword())
														*wg.cx.Config.WalletPass = string(pass)
														Debug("password", string(pass))
														save.Pod(wg.cx.Config)
														w, err := loader.CreateNewWallet(
															pass,
															pass,
															seed,
															time.Now(),
															false,
															wg.cx.Config,
															nil,
														)
														Debug("*** created wallet")
														if Check(err) {
															// return
														}
														Debug("refilling mining addresses")
														addresses.RefillMiningAddresses(
															w,
															wg.cx.Config,
															wg.cx.StateCfg,
														)
														Warn("done refilling mining addresses")
														Debug("starting main app")
														*wg.cx.Config.Generate = true
														*wg.cx.Config.GenThreads = 1
														*wg.cx.Config.NodeOff = false
														*wg.cx.Config.WalletOff = false
														save.Pod(wg.cx.Config)
														wg.miner.Start()
														*wg.noWallet = false
														wg.node.Start()
														if err = wg.writeWalletCookie(); Check(err) {
														}
														wg.wallet.Start()
													}()
												},
											).
											CornerRadius(0).
											Inset(0.5).
											Text("create wallet").
											Fn,
									).
									Fn(gtx)
							},
						).
						
						Fn(gtx)
				},
			).
			Flexed(0.5, gui.EmptyMaxWidth()).Fn,
	).Fn).Fn(gtx)
}
