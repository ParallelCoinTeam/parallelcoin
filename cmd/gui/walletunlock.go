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

	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/text"

	"github.com/p9c/interrupt"

	"github.com/p9c/gel"
	"github.com/p9c/matrjoska/pkg/p9icons"
)

func (wg *WalletGUI) unlockWallet(pass string) {
	D.Ln("entered password", pass)
	// unlock wallet
	// wg.cx.Config.Lock()
	wg.cx.Config.WalletPass.Set(pass)
	wg.cx.Config.WalletOff.F()
	// wg.cx.Config.Unlock()
	// load config into a fresh variable
	// cfg := podcfgs.GetDefaultConfig()
	var cfgFile []byte
	var e error
	if cfgFile, e = ioutil.ReadFile(wg.cx.Config.ConfigFile.V()); E.Chk(e) {
		// this should not happen
		// TODO: panic-type conditions - for gel should have a notification maybe?
		panic("config file does not exist")
	}
	cfg := wg.cx.Config
	D.Ln("loaded config")
	if e = json.Unmarshal(cfgFile, &cfg); !E.Chk(e) {
		D.Ln("unmarshaled config")
		bhb := blake3.Sum256([]byte(pass))
		bh := hex.EncodeToString(bhb[:])
		I.Ln(pass, bh, cfg.WalletPass.V())
		if cfg.WalletPass.V() == bh {
			D.Ln("loading previously saved state")
			filename := filepath.Join(wg.cx.Config.DataDir.V(), "state.json")
			// if log.FileExists(filename) {
			I.Ln("#### loading state data...")
			if e = wg.State.Load(filename, wg.cx.Config.WalletPass.Bytes()); !E.Chk(e) {
				D.Ln("#### loaded state data")
			}
			// it is as though it is loaded if it didn't exist
			wg.stateLoaded.Store(true)
			// the entered password matches the stored hash
			wg.cx.Config.NodeOff.F()
			wg.cx.Config.WalletOff.F()
			wg.cx.Config.WalletPass.Set(pass)
			if e = wg.cx.Config.WriteToFile(wg.cx.Config.ConfigFile.V()); E.Chk(e) {
			}
			wg.cx.Config.WalletPass.Set(pass)
			wg.WalletWatcher = wg.Watcher()
			// }
			//
			// qrText := fmt.Sprintf("parallelcoin:%s?amount=%s&message=%s",
			// 	wg.State.currentReceivingAddress.Load().EncodeAddress(),
			// 	wg.inputs["receiveAmount"].GetText(),
			// 	wg.inputs["receiveMessage"].GetText(),
			// )
			// var qrc image.Image
			// if qrc, e = qrcode.Encode(qrText, 0, qrcode.ECLevelL, 4); !E.Chk(e) {
			// 	iop := paint.NewImageOp(qrc)
			// 	wg.currentReceiveQRCode = &iop
			// 	wg.currentReceiveQR = wg.ButtonLayout(wg.currentReceiveCopyClickable.SetClick(func() {
			// 		D.Ln("clicked qr code copy clicker")
			// 		if e := clipboard.WriteAll(qrText); E.Chk(e) {
			// 		}
			// 	})).
			// 		// CornerRadius(0.5).
			// 		// Corners(gel.NW | gel.SW | gel.NE).
			// 		Background("white").
			// 		Embed(
			// 			wg.Inset(0.125,
			// 				wg.Image().Src(*wg.currentReceiveQRCode).Scale(1).Fn,
			// 			).Fn,
			// 		).Fn
			// 	// *wg.currentReceiveQRCode = iop
			// }
		}
	} else {
		D.Ln("failed to unlock the wallet")
	}
}

func (wg *WalletGUI) getWalletUnlockAppWidget() (a *gel.App) {
	a = wg.App(wg.Window.Width, wg.State.activePage, Break1).SetMainDirection(l.Center + 1)
	wg.unlockPage = a
	password := wg.cx.Config.WalletPass
	exitButton := wg.WidgetPool.GetClickable()
	unlockButton := wg.WidgetPool.GetClickable()
	wg.unlockPassword = wg.Password(
		"enter password", password, "DocText",
		"DocBg", "PanelBg", func(pass string) {
			I.Ln("wallet unlock initiated", pass)
			wg.unlockWallet(pass)
		},
	)
	// wg.unlockPage.SetThemeHook(
	// 	func() {
	// 		D.Ln("theme hook")
	// 		// D.Ln(wg.bools)
	// 		wg.cx.Config.DarkTheme.Set(*wg.Dark)
	// 		b := wg.configs["config"]["DarkTheme"].Slot.(*bool)
	// 		*b = *wg.Dark
	// 		if wgb, ok := wg.config.Bools["DarkTheme"]; ok {
	// 			wgb.Value(*wg.Dark)
	// 		}
	// 		var e error
	// 		if e = wg.cx.Config.WriteToFile(wg.cx.Config.ConfigFile.V()); E.Chk(e) {
	// 		}
	// 	},
	// )
	a.Pages(
		map[string]l.Widget{
			"home": wg.Page(
				"home", gel.Widgets{
					gel.WidgetSize{
						Widget:
						func(gtx l.Context) l.Dimensions {
							var dims l.Dimensions
							return wg.Flex().
								AlignMiddle().
								Flexed(
									1,
									wg.VFlex().
										Flexed(0.5, gel.EmptyMaxHeight()).
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
																									Scale(gel.Scales["H3"]).
																									Color("DocBg").
																									Src(&icons.ActionLock).Fn,
																							).Fn,
																						).Fn,
																					).
																					Rigid(
																						wg.Inset(
																							0.5,
																							gel.EmptySpace(0, 0),
																						).Fn,
																					).
																					Rigid(
																						wg.H2("locked").Color("DocText").Fn,
																					).
																					Fn(gtx)
																				return dims
																			},
																		).
																		Rigid(wg.Inset(0.5, gel.EmptySpace(0, 0)).Fn).
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
																		Rigid(wg.Inset(0.5, gel.EmptySpace(0, 0)).Fn).
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
																						Scale(gel.Scales["Caption"]).
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
																													gel.Scales["H4"],
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
																												gel.EmptySpace(
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
																												gel.EmptySpace(
																													0,
																													0,
																												),
																											).Fn,
																										).
																										Fn,
																								).Fn,
																								// l.Center,
																								// wg.TextSize.True/2).Fn,
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
																									D.Ln(
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
																												Scale(gel.Scales["H4"]).
																												Color("Light").
																												Src(&icons.ActionLockOpen).Fn,
																										).
																										Rigid(
																											wg.Inset(
																												0.5,
																												gel.EmptySpace(
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
																												gel.EmptySpace(
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
										).Flexed(0.5, gel.EmptyMaxHeight()).Fn,
								).
								Fn(gtx)
						},
					},
				},
			),
			"settings": wg.Page(
				"settings", gel.Widgets{
					gel.WidgetSize{
						Widget: func(gtx l.Context) l.Dimensions {
							return wg.configs.Widget(wg.config)(gtx)
						},
					},
				},
			),
			"console": wg.Page(
				"console", gel.Widgets{
					gel.WidgetSize{Widget: wg.console.Fn},
				},
			),
			"help": wg.Page(
				"help", gel.Widgets{
					gel.WidgetSize{Widget: wg.HelpPage()},
				},
			),
			"log": wg.Page(
				"log", gel.Widgets{
					gel.WidgetSize{Widget: a.Placeholder("log")},
				},
			),
			"quit": wg.Page(
				"quit", gel.Widgets{
					gel.WidgetSize{
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
										Flexed(0.5, gel.EmptyMaxWidth()).
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
										Flexed(0.5, gel.EmptyMaxWidth()).
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
				"mining", gel.Widgets{
					gel.WidgetSize{Widget: a.Placeholder("mining")},
				},
			),
			"explorer": wg.Page(
				"explorer", gel.Widgets{
					gel.WidgetSize{Widget: a.Placeholder("explorer")},
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
			// wg.Flex().Rigid(wg.Inset(0.5, gel.EmptySpace(0, 0)).Fn).Fn,
			// wg.PageTopBarButton(
			// 	"quit", 3, &icons.ActionExitToApp, func(name string) {
			// 		wg.unlockPage.ActivePage(name)
			// 	}, wg.unlockPage, "",
			// ),
		},
	)
	a.StatusBar(
		[]l.Widget{
			// wg.Inset(0.5, gel.EmptySpace(0, 0)).Fn,
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
					D.Ln("click on button", name)
					wg.unlockPage.ActivePage(name)
				}, wg.unlockPage,
			),
			wg.StatusBarButton(
				"settings", 5, &icons.ActionSettings, func(name string) {
					wg.unlockPage.ActivePage(name)
				}, wg.unlockPage,
			),
			// wg.Inset(0.5, gel.EmptySpace(0, 0)).Fn,
		},
	)
	// a.PushOverlay(wg.toasts.DrawToasts())
	// a.PushOverlay(wg.dialog.DrawDialog())
	return
}
