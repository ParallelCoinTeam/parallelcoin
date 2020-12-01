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
	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/util/hdkeychain"
	"github.com/p9c/pod/pkg/wallet"
)

//
// func (wg *WalletGUI) CreateWalletPageOld(gtx l.Context) l.Dimensions {
// 	return wg.th.
// 		Fill("PanelBg",
// 			wg.th.Flex().SpaceAround().AlignMiddle().
// 				Rigid(wg.th.
// 					VFlex().AlignMiddle().SpaceAround().
// 					Rigid(
// 						wg.th.VFlex().SpaceAround().AlignMiddle().
// 							Rigid(
// 								wg.th.H4("create new wallet").
// 									Color("PanelText").
// 									// Alignment(text.Middle).
// 									Fn,
// 							).
// 							Rigid(
// 								wg.th.Inset(0.25,
// 									wg.passwords["passEditor"].Fn,
// 								).Fn,
// 							).
// 							Rigid(
// 								wg.th.Inset(0.25,
// 									wg.passwords["confirmPassEditor"].Fn,
// 								).Fn,
// 							).
// 							Rigid(
// 								wg.th.Inset(0.25,
// 									wg.inputs["walletSeed"].Fn,
// 								).Fn,
// 							).
// 							Rigid(
// 								wg.th.Inset(0.25,
// 									wg.passwords["publicPassEditor"].Fn,
// 								).Fn,
// 							).
// 							Rigid(
// 								wg.th.Body1("your seed").
// 									Color("PanelText").
// 									Fn,
// 							).
// 							Rigid(
// 								wg.th.Inset(0.5,
// 									func(gtx l.Context) l.Dimensions {
// 										gtx.Constraints.Max.X = int(wg.th.TextSize.Scale(36).V)
// 										gtx.Constraints.Min.X = int(wg.th.TextSize.Scale(16).V)
// 										return wg.th.CheckBox(wg.bools["testnet"].SetOnChange(func(b bool) {
// 											Debug("confirmed read", b)
// 										})).
// 											IconColor("Primary").
// 											TextColor("DocText").
// 											Text("use testnet?").
// 											Fn(gtx)
// 									},
// 								).Fn,
// 							).
// 							// Rigid(
// 							// 	wg.th.Inset(0.25,
// 							// 		func(gtx l.Context) l.Dimensions {
// 							// 			// gtx.Constraints.Min.X = int(wg.th.TextSize.Scale(16).V)
// 							// 			return wg.th.CheckBox(wg.bools["testnet"].SetOnChange(func(b bool) {
// 							// 				Debug("testnet on?", b)
// 							// 				if b {
// 							// 					wg.cx.ActiveNet = &netparams.TestNet3Params
// 							// 					fork.IsTestnet = true
// 							// 				} else {
// 							// 					wg.cx.ActiveNet = &netparams.MainNetParams
// 							// 					fork.IsTestnet = false
// 							// 				}
// 							// 				Info("activenet:", wg.cx.ActiveNet.Name)
// 							// 				*wg.cx.Config.Network = wg.cx.ActiveNet.Name
// 							// 				Debug("")
// 							// 				if wg.cx.ActiveNet.Name == "testnet" {
// 							// 					// TODO: obviously when we get to starting testnets this should not be done
// 							// 					*wg.cx.Config.LAN = true  // mines without peer outside lan
// 							// 					*wg.cx.Config.Solo = true // mines without peers
// 							// 				}
// 							// 				save.Pod(wg.cx.Config)
// 							// 			})).
// 							// 				IconColor("Primary").
// 							// 				TextColor("DocText").
// 							// 				Text("Use testnet?").
// 							// 				Fn(gtx)
// 							// 		},
// 							// 	).Fn,
// 							// ).
// 							Rigid(
// 								func(gtx l.Context) l.Dimensions {
// 									gtx.Constraints.Max.X = int(wg.th.TextSize.Scale(22).V)
// 									return wg.th.Caption(wg.inputs["walletSeed"].GetText()).
// 										Font("go regular").
// 										TextScale(0.66).
// 										Fn(gtx)
// 								},
// 							).
// 							// Rigid(
// 							// 	wg.th.Inset(0.5,
// 							// 		func(gtx l.Context) l.Dimensions {
// 							// 			gtx.Constraints.Max.X = int(wg.th.TextSize.Scale(36).V)
// 							// 			gtx.Constraints.Min.X = int(wg.th.TextSize.Scale(16).V)
// 							// 			return wg.th.CheckBox(wg.bools["ihaveread"].SetOnChange(func(b bool) {
// 							// 				Debug("confirmed read", b)
// 							// 			})).
// 							// 				IconColor("Primary").
// 							// 				TextColor("DocText").
// 							// 				Text("I have stored the seed and password safely " +
// 							// 					"and understand it cannot be recovered").
// 							// 				Fn(gtx)
// 							// 		},
// 							// 	).Fn,
// 							// ).
// 							Rigid(
// 								func(gtx l.Context) l.Dimensions {
// 									var b []byte
// 									var err error
// 									seedValid := true
// 									if b, err = hex.DecodeString(wg.inputs["walletSeed"].GetText()); Check(err) {
// 										seedValid = false
// 									} else if len(b) != 0 && len(b) < hdkeychain.MinSeedBytes ||
// 										len(b) > hdkeychain.MaxSeedBytes {
// 										seedValid = false
// 									}
// 									if wg.passwords["passEditor"].GetPassword() == "" ||
// 										wg.passwords["confirmPassEditor"].GetPassword() == "" ||
// 										len(wg.passwords["passEditor"].GetPassword()) < 8 ||
// 										wg.passwords["passEditor"].GetPassword() !=
// 											wg.passwords["confirmPassEditor"].GetPassword() ||
// 										!seedValid ||
// 										!wg.bools["ihaveread"].GetValue() {
// 										gtx = gtx.Disabled()
// 									}
// 									return wg.th.Flex().
// 										Rigid(
// 											wg.th.Button(wg.clickables["createWallet"]).
// 												Background("Primary").
// 												Color("Light").
// 												SetClick(func() {
// 													// go func() {
// 													// wg.ShellRunCommandChan <- "stop"
// 													Debug("clicked submit wallet")
// 													*wg.cx.Config.WalletFile = *wg.cx.Config.DataDir +
// 														string(os.PathSeparator) + wg.cx.ActiveNet.Name +
// 														string(os.PathSeparator) + wallet.WalletDbName
// 													dbDir := *wg.cx.Config.WalletFile
// 													loader := wallet.NewLoader(wg.cx.ActiveNet, dbDir, 250)
// 													seed, _ := hex.DecodeString(wg.inputs["walletSeed"].GetText())
// 													w, err := loader.CreateNewWallet(
// 														[]byte(wg.passwords["publicPassEditor"].GetPassword()),
// 														[]byte(wg.passwords["passEditor"].GetPassword()),
// 														seed,
// 														time.Now(),
// 														false,
// 														wg.cx.Config,
// 													)
// 													if Check(err) {
// 														panic(err)
// 													}
// 													Warn("refilling mining addresses")
// 													addresses.RefillMiningAddresses(w, wg.cx.Config, wg.cx.StateCfg)
// 													Warn("done refilling mining addresses")
// 													w.Manager.Close()
// 													w.Stop()
// 													// Debug("starting up shell first time")
// 													// rand.Seed(time.Now().Unix())
// 													// nodeport := rand.Intn(60000) + 1024
// 													// walletport := rand.Intn(60000) + 1024
// 													// *wg.cx.Config.RPCListeners = []string{fmt.Sprintf("127.0.0.1:%d", nodeport)}
// 													// *wg.cx.Config.RPCConnect = fmt.Sprintf("127.0.0.1:%d", nodeport)
// 													// *wg.cx.Config.WalletRPCListeners = []string{fmt.Sprintf("127.0.0.1:%d", walletport)}
// 													// *wg.cx.Config.WalletServer = fmt.Sprintf("127.0.0.1:%d", walletport)
// 													// *wg.cx.Config.ServerTLS = false
// 													// *wg.cx.Config.TLS = false
// 													// *wg.cx.Config.GenThreads = 1 // probably want it to be max ultimately
// 													// wg.incdecs["generatethreads"].Current = 1
// 													// *wg.cx.Config.Generate = true // probably don't want on ultimately
// 													// save.Pod(wg.cx.Config)
//
// 													// Debug("opening wallet")
// 													// w, err = loader.OpenExistingWallet([]byte(*wg.cx.Config.WalletPass),
// 													// 	false, wg.cx.Config)
// 													// if err != nil {
// 													// 	panic(err)
// 													// }
// 													// args := []string{os.Args[0], "-D", *wg.cx.Config.DataDir,
// 													// 	"--pipelog", "wallet", "drophistory"}
// 													// runner := exec.Command(args[0], args[1:]...)
// 													// runner.Stderr = os.Stderr
// 													// runner.Stdout = os.Stderr
// 													// if err := runner.Start(); Check(err) {
// 													// }
// 													// time.Sleep(time.Second * 10)
// 													// wg.ShellRunCommandChan <- "stop"
// 													// wg.ShellRunCommandChan <- "run"
// 													// wg.ShellRunCommandChan <- "stop"
// 													// wg.ShellRunCommandChan <- "run"
// 													// time.Sleep(time.Second * 10)
// 													// time.Sleep(time.Second * 2)
// 													// interrupt.RequestRestart()
// 													// procAttr := new(os.ProcAttr)
// 													// procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
// 													// os.StartProcess(os.Args[0], os.Args[1:], procAttr)
// 													// *wg.App = *wg.GetAppWidget()
// 													Debug("starting main app")
// 													*wg.noWallet = false
// 													wg.running = false
// 													wg.mining = false
// 													if err = wg.Runner(); Check(err) {
// 													}
// 													wg.ShellRunCommandChan <- "run"
// 													wg.MinerRunCommandChan <- "run"
// 													// Exec()
// 													// }()
// 												}).
// 												CornerRadius(0).
// 												Inset(0.5).
// 												Text("create wallet").
// 												Fn,
// 										).
// 										Fn(gtx)
// 								},
// 							).
// 							Fn,
// 					).
// 					Fn,
// 				).
// 				Fn,
// 		).
// 		Fn(gtx)
// }

func (wg *WalletGUI) CreateWalletPage(gtx l.Context) l.Dimensions {
	return wg.th.Fill("PanelBg",
		wg.th.Inset(0.5,
			wg.th.Flex().
				SpaceAround().
				Flexed(0.5, p9.EmptyMaxHeight()).
				Rigid(
					func(gtx l.Context) l.Dimensions {
						return wg.th.VFlex().
							AlignMiddle().
							SpaceSides().
							Rigid(
								wg.th.H4("create new wallet").
									Color("PanelText").
									Fn,
							).
							Rigid(
								wg.th.Inset(0.25,
									wg.passwords["passEditor"].Fn,
								).
									Fn,
							).
							Rigid(
								wg.th.Inset(0.25,
									wg.passwords["confirmPassEditor"].Fn,
								).
									Fn,
							).
							Rigid(
								wg.th.Inset(0.25,
									wg.inputs["walletSeed"].Fn,
								).
									Fn,
							).
							Rigid(
								wg.th.Inset(0.25,
									func(gtx l.Context) l.Dimensions {
										// gtx.Constraints.Min.X = int(wg.th.TextSize.Scale(16).V)
										return wg.th.CheckBox(wg.bools["testnet"].SetOnChange(func(b bool) {
											Debug("testnet on?", b)
											if b {
												wg.cx.ActiveNet = &netparams.TestNet3Params
												fork.IsTestnet = true
											} else {
												wg.cx.ActiveNet = &netparams.MainNetParams
												fork.IsTestnet = false
											}
											Info("activenet:", wg.cx.ActiveNet.Name)
											*wg.cx.Config.Network = wg.cx.ActiveNet.Name
											Debug("")
											if wg.cx.ActiveNet.Name == "testnet" {
												// TODO: obviously when we get to starting testnets this should not be done
												*wg.cx.Config.LAN = true  // mines without peer outside lan
												*wg.cx.Config.Solo = true // mines without peers
											}
											save.Pod(wg.cx.Config)
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
									gtx.Constraints.Max.X = int(wg.th.TextSize.Scale(22).V)
									return wg.th.Caption(wg.inputs["walletSeed"].GetText()).
										Font("go regular").
										TextScale(0.66).
										Fn(gtx)
								},
							).

							Rigid(
								wg.th.Inset(0.5,
									func(gtx l.Context) l.Dimensions {
										gtx.Constraints.Max.X = int(wg.th.TextSize.Scale(32).V)
										gtx.Constraints.Min.X = int(wg.th.TextSize.Scale(16).V)
										return wg.th.CheckBox(wg.bools["ihaveread"].SetOnChange(func(b bool) {
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
													// go func() {
													// wg.ShellRunCommandChan <- "stop"
													Debug("clicked submit wallet")
													*wg.cx.Config.WalletFile = *wg.cx.Config.DataDir +
														string(os.PathSeparator) + wg.cx.ActiveNet.Name +
														string(os.PathSeparator) + wallet.WalletDbName
													dbDir := *wg.cx.Config.WalletFile
													loader := wallet.NewLoader(wg.cx.ActiveNet, dbDir, 250)
													seed, _ := hex.DecodeString(wg.inputs["walletSeed"].GetText())
													pass := []byte(wg.passwords["passEditor"].GetPassword())
													w, err := loader.CreateNewWallet(
														pass, pass,
														seed,
														time.Now(),
														false,
														wg.cx.Config,
													)
													if Check(err) {
														panic(err)
													}
													Warn("refilling mining addresses")
													addresses.RefillMiningAddresses(w, wg.cx.Config, wg.cx.StateCfg)
													Warn("done refilling mining addresses")
													w.Manager.Close()
													w.Stop()
													// Debug("starting up shell first time")
													// rand.Seed(time.Now().Unix())
													// nodeport := rand.Intn(60000) + 1024
													// walletport := rand.Intn(60000) + 1024
													// *wg.cx.Config.RPCListeners = []string{fmt.Sprintf("127.0.0.1:%d", nodeport)}
													// *wg.cx.Config.RPCConnect = fmt.Sprintf("127.0.0.1:%d", nodeport)
													// *wg.cx.Config.WalletRPCListeners = []string{fmt.Sprintf("127.0.0.1:%d", walletport)}
													// *wg.cx.Config.WalletServer = fmt.Sprintf("127.0.0.1:%d", walletport)
													// *wg.cx.Config.ServerTLS = false
													// *wg.cx.Config.TLS = false
													// *wg.cx.Config.GenThreads = 1 // probably want it to be max ultimately
													// wg.incdecs["generatethreads"].Current = 1
													// *wg.cx.Config.Generate = true // probably don't want on ultimately
													// save.Pod(wg.cx.Config)

													// Debug("opening wallet")
													// w, err = loader.OpenExistingWallet([]byte(*wg.cx.Config.WalletPass),
													// 	false, wg.cx.Config)
													// if err != nil {
													// 	panic(err)
													// }
													// args := []string{os.Args[0], "-D", *wg.cx.Config.DataDir,
													// 	"--pipelog", "wallet", "drophistory"}
													// runner := exec.Command(args[0], args[1:]...)
													// runner.Stderr = os.Stderr
													// runner.Stdout = os.Stderr
													// if err := runner.Start(); Check(err) {
													// }
													// time.Sleep(time.Second * 10)
													// wg.ShellRunCommandChan <- "stop"
													// wg.ShellRunCommandChan <- "run"
													// wg.ShellRunCommandChan <- "stop"
													// wg.ShellRunCommandChan <- "run"
													// time.Sleep(time.Second * 10)
													// time.Sleep(time.Second * 2)
													// interrupt.RequestRestart()
													// procAttr := new(os.ProcAttr)
													// procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
													// os.StartProcess(os.Args[0], os.Args[1:], procAttr)
													// *wg.App = *wg.GetAppWidget()
													Debug("starting main app")
													*wg.noWallet = false
													wg.running = false
													wg.mining = false
													if err = wg.Runner(); Check(err) {
													}
													wg.ShellRunCommandChan <- "run"
													wg.MinerRunCommandChan <- "run"
													// Exec()
													// }()
												}).
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
				Flexed(0.5, p9.EmptyMaxWidth()).Fn,
		).Fn,
	).Fn(gtx)
}
