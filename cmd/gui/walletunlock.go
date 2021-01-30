package gui

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"
	
	"golang.org/x/exp/shiny/materialdesign/icons"
	"lukechampine.com/blake3"
	
	l "gioui.org/layout"
	"gioui.org/text"
	
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/gui"
	p9icons "github.com/p9c/pod/pkg/gui/ico/svg"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/util/logi"
)

func (wg *WalletGUI) unlockWallet(pass string) {
	Debug("entered password", pass)
	// unlock wallet
	wg.cx.Config.Lock()
	*wg.cx.Config.WalletPass = pass
	*wg.cx.Config.WalletOff = false
	wg.cx.Config.Unlock()
	// load config into a fresh variable
	cfg, _ := pod.EmptyConfig()
	var cfgFile []byte
	var err error
	if cfgFile, err = ioutil.ReadFile(*wg.cx.Config.ConfigFile); Check(err) {
		// this should not happen
		// TODO: panic-type conditions - for gui should have a notification maybe?
		panic("config file does not exist")
	}
	Debug("loaded config")
	if err = json.Unmarshal(cfgFile, &cfg); !Check(err) {
		Debug("unmarshaled config")
		bhb := blake3.Sum256([]byte(pass))
		bh := hex.EncodeToString(bhb[:])
		Debug(pass, bh, *cfg.WalletPass)
		if *cfg.WalletPass == bh {
			Debug("loading previously saved state")
			filename := filepath.Join(wg.cx.DataDir, "state.json")
			if logi.FileExists(filename) {
				Debug("#### loading state data...")
				if err = wg.State.Load(filename, wg.cx.Config.WalletPass); Check(err) {
					interrupt.Request()
				}
				Debug("#### loaded state data")
			}
			//
			// qrText := fmt.Sprintf("parallelcoin:%s?amount=%s&message=%s",
			// 	wg.State.currentReceivingAddress.Load().EncodeAddress(),
			// 	wg.inputs["receiveAmount"].GetText(),
			// 	wg.inputs["receiveMessage"].GetText(),
			// )
			// var qrc image.Image
			// if qrc, err = qrcode.Encode(qrText, 0, qrcode.ECLevelL, 4); !Check(err) {
			// 	iop := paint.NewImageOp(qrc)
			// 	wg.currentReceiveQRCode = &iop
			// 	wg.currentReceiveQR = wg.ButtonLayout(wg.currentReceiveCopyClickable.SetClick(func() {
			// 		Debug("clicked qr code copy clicker")
			// 		if err := clipboard.WriteAll(qrText); Check(err) {
			// 		}
			// 	})).
			// 		// CornerRadius(0.5).
			// 		// Corners(gui.NW | gui.SW | gui.NE).
			// 		Background("white").
			// 		Embed(
			// 			wg.Inset(0.125,
			// 				wg.Image().Src(*wg.currentReceiveQRCode).Scale(1).Fn,
			// 			).Fn,
			// 		).Fn
			// 	// *wg.currentReceiveQRCode = iop
			// }
			wg.stateLoaded.Store(true)
			
			wg.RecentTransactions(10, "recent")
			// wg.Invalidate()
			wg.RecentTransactions(-1, "history")
			// the entered password matches the stored hash
			Debug("now we can open the wallet")
			if err = wg.writeWalletCookie(); Check(err) {
			}
			*wg.cx.Config.NodeOff = false
			*wg.cx.Config.WalletOff = false
			save.Pod(wg.cx.Config)
			if !wg.node.Running() {
				wg.node.Start()
			} else {
				if wg.ChainClient != nil {
					wg.ChainClient.Disconnect()
					wg.ChainClient.Shutdown()
				}
				wg.node.Stop()
				wg.node.Start()
			}
			wg.wallet.Start()
			if err = wg.chainClient(); Check(err) {
			}
			if err = wg.walletClient(); Check(err) {
			}
			wg.unlockPassword.Wipe()
			wg.ready.Store(true)
			wg.Invalidate()
		}
	} else {
		Debug("failed to unlock the wallet")
	}
}

func (wg *WalletGUI) getWalletUnlockAppWidget() (a *gui.App) {
	a = wg.App(&wg.Window.Width, wg.State.activePage, wg.invalidate, Break1).SetMainDirection(l.Center + 1)
	wg.unlockPage = a
	password := ""
	exitButton := wg.WidgetPool.GetClickable()
	unlockButton := wg.WidgetPool.GetClickable()
	wg.unlockPassword = wg.Password(
		"enter password", &password, "DocText",
		"DocBg", "PanelBg", func(pass string) {
			go wg.unlockWallet(pass)
		},
	)
	wg.unlockPage.ThemeHook(
		func() {
			Debug("theme hook")
			// Debug(wg.bools)
			*wg.cx.Config.DarkTheme = *wg.Dark
			a := wg.configs["config"]["DarkTheme"].Slot.(*bool)
			*a = *wg.Dark
			if wgb, ok := wg.config.Bools["DarkTheme"]; ok {
				wgb.Value(*wg.Dark)
			}
			save.Pod(wg.cx.Config)
		},
	)
	a.Pages(
		map[string]l.Widget{
			"home": wg.Page(
				"home", gui.Widgets{
					gui.WidgetSize{
						Widget:
						func(gtx l.Context) l.Dimensions {
							var dims l.Dimensions
							return wg.Flex().
								AlignMiddle().
								Flexed(
									1,
									wg.VFlex().
										Flexed(0.5, gui.EmptyMaxHeight()).
										Rigid(
											wg.Flex().
												SpaceEvenly().
												AlignMiddle().
												Rigid(
													wg.Fill(
														"DocBg", l.Center, wg.TextSize.V, 0,
														wg.Inset(
															0.5,
															wg.Flex().
																AlignMiddle().
																Rigid(
																	wg.VFlex().
																		AlignMiddle().
																		Rigid(
																			func(gtx l.Context) l.Dimensions {
																				dims = wg.Flex().
																					AlignBaseline().
																					Rigid(
																						wg.Fill(
																							"Fatal",
																							l.Center,
																							wg.TextSize.V/2,
																							0,
																							wg.Inset(
																								0.5,
																								wg.Icon().
																									Scale(gui.Scales["H3"]).
																									Color("DocBg").
																									Src(&icons.ActionLock).Fn,
																							).Fn,
																						).Fn,
																					).
																					Rigid(
																						wg.Inset(
																							0.5,
																							gui.EmptySpace(0, 0),
																						).Fn,
																					).
																					Rigid(
																						wg.H2("locked").Color("DocText").Fn,
																					).
																					Fn(gtx)
																				return dims
																			},
																		).
																		Rigid(wg.Inset(0.5, gui.EmptySpace(0, 0)).Fn).
																		Rigid(
																			func(gtx l.Context) l.
																			Dimensions {
																				gtx.Constraints.Max.
																					X = dims.Size.X
																				return wg.
																					unlockPassword.
																					Fn(gtx)
																			},
																		).
																		Rigid(wg.Inset(0.5, gui.EmptySpace(0, 0)).Fn).
																		Rigid(
																			wg.Body1(
																				fmt.Sprintf(
																					"%v idle timeout",
																					time.Duration(wg.incdecs["idleTimeout"].GetCurrent())*time.Second,
																				),
																			).
																				Color("DocText").
																				Font("bariol bold").
																				Fn,
																		).
																		Rigid(
																			wg.Flex().
																				Rigid(
																					wg.Body1("Idle timeout in seconds:").Color(
																						"DocText",
																					).Fn,
																				).
																				Rigid(
																					wg.incdecs["idleTimeout"].
																						Color("DocText").
																						Background("DocBg").
																						Scale(gui.Scales["Caption"]).
																						Fn,
																				).
																				Fn,
																		).
																		Rigid(
																			wg.Flex().
																				Rigid(
																					wg.Inset(
																						0.25,
																						wg.ButtonLayout(
																							exitButton.SetClick(
																								func() {
																									interrupt.Request()
																								},
																							),
																						).
																							CornerRadius(0.5).
																							Corners(0).
																							Background("PanelBg").
																							Embed(
																								// wg.Fill("DocText",
																								wg.Inset(
																									0.25,
																									wg.Flex().AlignMiddle().
																										Rigid(
																											wg.Icon().
																												Scale(
																													gui.
																														Scales["H4"],
																												).
																												Color("DocText").
																												Src(
																													&icons.
																														MapsDirectionsRun,
																												).Fn,
																										).
																										Rigid(
																											wg.Inset(
																												0.5,
																												gui.EmptySpace(
																													0,
																													0,
																												),
																											).Fn,
																										).
																										Rigid(
																											wg.H6("exit").Color("DocText").Fn,
																										).
																										Rigid(
																											wg.Inset(
																												0.5,
																												gui.EmptySpace(
																													0,
																													0,
																												),
																											).Fn,
																										).
																										Fn,
																								).Fn,
																								// l.Center,
																								// wg.TextSize.V/2).Fn,
																							).Fn,
																					).Fn,
																				).
																				Rigid(
																					wg.Inset(
																						0.25,
																						wg.ButtonLayout(
																							unlockButton.SetClick(
																								func() {
																									// pass := wg.unlockPassword.Editor().Text()
																									pass := wg.unlockPassword.GetPassword()
																									Debug(
																										">>>>>>>>>>> unlock password",
																										pass,
																									)
																									wg.unlockWallet(pass)
																									
																								},
																							),
																						).Background("Primary").
																							CornerRadius(0.5).
																							Corners(0).
																							Embed(
																								wg.Inset(
																									0.25,
																									wg.Flex().AlignMiddle().
																										Rigid(
																											wg.Icon().
																												Scale(gui.Scales["H4"]).
																												Color("Light").
																												Src(&icons.ActionLockOpen).Fn,
																										).
																										Rigid(
																											wg.Inset(
																												0.5,
																												gui.EmptySpace(
																													0,
																													0,
																												),
																											).Fn,
																										).
																										Rigid(
																											wg.H6("unlock").Color("Light").Fn,
																										).
																										Rigid(
																											wg.Inset(
																												0.5,
																												gui.EmptySpace(
																													0,
																													0,
																												),
																											).Fn,
																										).
																										Fn,
																								).Fn,
																							).Fn,
																					).Fn,
																				).
																				Fn,
																		).
																		Fn,
																).
																Fn,
														).Fn,
													).Fn,
												).
												Fn,
										).Flexed(0.5, gui.EmptyMaxHeight()).Fn,
								).
								Fn(gtx)
						},
					},
				},
			),
			"settings": wg.Page(
				"settings", gui.Widgets{
					gui.WidgetSize{
						Widget: func(gtx l.Context) l.Dimensions {
							return wg.configs.Widget(wg.config)(gtx)
						},
					},
				},
			),
			"console": wg.Page(
				"console", gui.Widgets{
					gui.WidgetSize{Widget: wg.console.Fn},
				},
			),
			"help": wg.Page(
				"help", gui.Widgets{
					gui.WidgetSize{Widget: wg.HelpPage()},
				},
			),
			"log": wg.Page(
				"log", gui.Widgets{
					gui.WidgetSize{Widget: a.Placeholder("log")},
				},
			),
			"quit": wg.Page(
				"quit", gui.Widgets{
					gui.WidgetSize{
						Widget: func(gtx l.Context) l.Dimensions {
							return wg.VFlex().
								SpaceEvenly().
								AlignMiddle().
								Rigid(
									wg.H4("are you sure?").Color(wg.unlockPage.BodyColorGet()).Alignment(text.Middle).Fn,
								).
								Rigid(
									wg.Flex().
										// SpaceEvenly().
										Flexed(0.5, gui.EmptyMaxWidth()).
										Rigid(
											wg.Button(
												wg.clickables["quit"].SetClick(
													func() {
														wg.gracefulShutdown()
														// close(wg.quit)
													},
												),
											).Color("Light").TextScale(5).Text(
												"yes!!!",
											).Fn,
										).
										Flexed(0.5, gui.EmptyMaxWidth()).
										Fn,
								).
								Fn(gtx)
						},
					},
				},
			),
			// "goroutines": wg.Page(
			// 	"log", p9.Widgets{
			// 		// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
			//
			// 		p9.WidgetSize{
			// 			Widget: func(gtx l.Context) l.Dimensions {
			// 				le := func(gtx l.Context, index int) l.Dimensions {
			// 					return wg.State.goroutines[index](gtx)
			// 				}
			// 				return func(gtx l.Context) l.Dimensions {
			// 					return wg.ButtonInset(
			// 						0.25,
			// 						wg.Fill(
			// 							"DocBg",
			// 							wg.lists["recent"].
			// 								Vertical().
			// 								// Background("DocBg").Color("DocText").Active("Primary").
			// 								Length(len(wg.State.goroutines)).
			// 								ListElement(le).
			// 								Fn,
			// 						).Fn,
			// 					).
			// 						Fn(gtx)
			// 				}(gtx)
			// 				// wg.NodeRunCommandChan <- "stop"
			// 				// consume.Kill(wg.Worker)
			// 				// consume.Kill(wg.cx.StateCfg.Miner)
			// 				// close(wg.cx.NodeKill)
			// 				// close(wg.cx.KillAll)
			// 				// time.Sleep(time.Second*3)
			// 				// interrupt.Request()
			// 				// os.Exit(0)
			// 				// return l.Dimensions{}
			// 			},
			// 		},
			// 	},
			// ),
			"mining": wg.Page(
				"mining", gui.Widgets{
					gui.WidgetSize{Widget: a.Placeholder("mining")},
				},
			),
			"explorer": wg.Page(
				"explorer", gui.Widgets{
					gui.WidgetSize{Widget: a.Placeholder("explorer")},
				},
			),
		},
	)
	// a.SideBar([]l.Widget{
	// 	wg.SideBarButton("overview", "overview", 0),
	// 	wg.SideBarButton("send", "send", 1),
	// 	wg.SideBarButton("receive", "receive", 2),
	// 	wg.SideBarButton("history", "history", 3),
	// 	wg.SideBarButton("explorer", "explorer", 6),
	// 	wg.SideBarButton("mining", "mining", 7),
	// 	wg.SideBarButton("console", "console", 9),
	// 	wg.SideBarButton("settings", "settings", 5),
	// 	wg.SideBarButton("log", "log", 10),
	// 	wg.SideBarButton("help", "help", 8),
	// 	wg.SideBarButton("quit", "quit", 11),
	// })
	a.ButtonBar(
		[]l.Widget{
			// wg.PageTopBarButton(
			// 	"goroutines", 0, &icons.ActionBugReport, func(name string) {
			// 		wg.unlockPage.ActivePage(name)
			// 	}, wg.unlockPage, "",
			// ),
			wg.PageTopBarButton(
				"help", 1, &icons.ActionHelp, func(name string) {
					wg.unlockPage.ActivePage(name)
				}, wg.unlockPage, "",
			),
			wg.PageTopBarButton(
				"home", 4, &icons.ActionLock, func(name string) {
					wg.unlockPage.ActivePage(name)
				}, wg.unlockPage, "Danger",
			),
			// wg.Flex().Rigid(wg.Inset(0.5, gui.EmptySpace(0, 0)).Fn).Fn,
			// wg.PageTopBarButton(
			// 	"quit", 3, &icons.ActionExitToApp, func(name string) {
			// 		wg.unlockPage.ActivePage(name)
			// 	}, wg.unlockPage, "",
			// ),
		},
	)
	a.StatusBar(
		[]l.Widget{
			// wg.Inset(0.5, gui.EmptySpace(0, 0)).Fn,
			wg.RunStatusPanel,
		},
		[]l.Widget{
			wg.StatusBarButton(
				"console", 3, &p9icons.Terminal, func(name string) {
					wg.MainApp.ActivePage(name)
				}, a,
			),
			wg.StatusBarButton(
				"log", 4, &icons.ActionList, func(name string) {
					Debug("click on button", name)
					wg.unlockPage.ActivePage(name)
				}, wg.unlockPage,
			),
			wg.StatusBarButton(
				"settings", 5, &icons.ActionSettings, func(name string) {
					wg.unlockPage.ActivePage(name)
				}, wg.unlockPage,
			),
			// wg.Inset(0.5, gui.EmptySpace(0, 0)).Fn,
		},
	)
	// a.PushOverlay(wg.toasts.DrawToasts())
	// a.PushOverlay(wg.dialog.DrawDialog())
	return
}
