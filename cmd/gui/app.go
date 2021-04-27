package gui

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	uberatomic "go.uber.org/atomic"
	"golang.org/x/exp/shiny/materialdesign/icons"

	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/text"

	"github.com/p9c/gel"
	"github.com/p9c/matrjoska/cmd/gui/cfg"
	"github.com/p9c/matrjoska/pkg/p9icons"
)

func (wg *WalletGUI) GetAppWidget() (a *gel.App) {
	a = wg.App(wg.Window.Width, uberatomic.NewString("home"), Break1).SetMainDirection(l.W)
	wg.MainApp = a
	wg.config = cfg.New(wg.Window, wg.quit)
	wg.configs = wg.config.Config()
	a.Pages(
		map[string]l.Widget{
			"home": wg.Page(
				"home", gel.Widgets{
					// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
					gel.WidgetSize{Widget: wg.OverviewPage()},
				},
			),
			"send": wg.Page(
				"send", gel.Widgets{
					// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
					gel.WidgetSize{Widget: wg.SendPage.Fn},
				},
			),
			"receive": wg.Page(
				"receive", gel.Widgets{
					// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
					gel.WidgetSize{Widget: wg.ReceivePage.Fn},
				},
			),
			"history": wg.Page(
				"history", gel.Widgets{
					// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
					gel.WidgetSize{Widget: wg.HistoryPage()},
				},
			),
			"settings": wg.Page(
				"settings", gel.Widgets{
					// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
					gel.WidgetSize{
						Widget: func(gtx l.Context) l.Dimensions {
							return wg.configs.Widget(wg.config)(gtx)
						},
					},
				},
			),
			"console": wg.Page(
				"console", gel.Widgets{
					// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
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
									wg.H4("are you sure?").Color(wg.MainApp.BodyColorGet()).Alignment(text.Middle).Fn,
								).
								Rigid(
									wg.Flex().
										// SpaceEvenly().
										Flexed(0.5, gel.EmptyMaxWidth()).
										Rigid(
											wg.Button(
												wg.clickables["quit"].SetClick(
													func() {
														// interrupt.Request()
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
					gel.WidgetSize{
						Widget: a.Placeholder("mining"),
					},
				},
			),
			"explorer": wg.Page(
				"explorer", gel.Widgets{
					gel.WidgetSize{
						Widget: a.Placeholder("explorer"),
					},
				},
			),
		},
	)
	a.SideBar(
		[]l.Widget{
			// wg.SideBarButton(" ", " ", 11),
			wg.SideBarButton("home", "home", 0),
			wg.SideBarButton("send", "send", 1),
			wg.SideBarButton("receive", "receive", 2),
			wg.SideBarButton("history", "history", 3),
			// wg.SideBarButton("explorer", "explorer", 6),
			// wg.SideBarButton("mining", "mining", 7),
			wg.SideBarButton("console", "console", 9),
			wg.SideBarButton("settings", "settings", 5),
			// wg.SideBarButton("log", "log", 10),
			wg.SideBarButton("help", "help", 8),
			// wg.SideBarButton(" ", " ", 11),
			// wg.SideBarButton("quit", "quit", 11),
		},
	)
	a.ButtonBar(
		[]l.Widget{
			
			// gel.EmptyMaxWidth(),
			// wg.PageTopBarButton(
			// 	"goroutines", 0, &icons.ActionBugReport, func(name string) {
			// 		wg.App.ActivePage(name)
			// 	}, a, "",
			// ),
			wg.PageTopBarButton(
				"help", 1, &icons.ActionHelp, func(name string) {
					wg.MainApp.ActivePage(name)
				}, a, "",
			),
			wg.PageTopBarButton(
				"home", 4, &icons.ActionLockOpen, func(name string) {
					wg.unlockPassword.Wipe()
					wg.unlockPassword.Focus()
					wg.WalletWatcher.Q()
					// if wg.WalletClient != nil {
					// 	wg.WalletClient.Disconnect()
					// 	wg.WalletClient = nil
					// }
					// wg.wallet.Stop()
					// wg.node.Stop()
					wg.State.SetActivePage("home")
					wg.unlockPage.ActivePage("home")
					wg.stateLoaded.Store(false)
					wg.ready.Store(false)
				}, a, "green",
			),
			// wg.Flex().Rigid(wg.Inset(0.5, gel.EmptySpace(0, 0)).Fn).Fn,
			// wg.PageTopBarButton(
			// 	"quit", 3, &icons.ActionExitToApp, func(name string) {
			// 		wg.MainApp.ActivePage(name)
			// 	}, a, "",
			// ),
		},
	)
	a.StatusBar(
		[]l.Widget{
			// func(gtx l.Context) l.Dimensions { return wg.RunStatusPanel(gtx) },
			// wg.Inset(0.5, gel.EmptySpace(0, 0)).Fn,
			// wg.Inset(0.5, gel.EmptySpace(0, 0)).Fn,
			wg.RunStatusPanel,
		},
		[]l.Widget{
			// gel.EmptyMaxWidth(),
			wg.StatusBarButton(
				"console", 3, &p9icons.Terminal, func(name string) {
					wg.MainApp.ActivePage(name)
				}, a,
			),
			wg.StatusBarButton(
				"log", 4, &icons.ActionList, func(name string) {
					D.Ln("click on button", name)
					if wg.MainApp.MenuOpen {
						wg.MainApp.MenuOpen = false
					}
					wg.MainApp.ActivePage(name)
				}, a,
			),
			wg.StatusBarButton(
				"settings", 5, &icons.ActionSettings, func(name string) {
					D.Ln("click on button", name)
					if wg.MainApp.MenuOpen {
						wg.MainApp.MenuOpen = false
					}
					wg.MainApp.ActivePage(name)
				}, a,
			),
			// wg.Inset(0.5, gel.EmptySpace(0, 0)).Fn,
		},
	)
	// a.PushOverlay(wg.toasts.DrawToasts())
	// a.PushOverlay(wg.dialog.DrawDialog())
	return
}

func (wg *WalletGUI) Page(title string, widget gel.Widgets) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		return wg.VFlex().
			// SpaceEvenly().
			Rigid(
				wg.Responsive(
					wg.Size.Load(), gel.Widgets{
						// p9.WidgetSize{
						// 	Widget: a.ButtonInset(0.25, a.H5(title).Color(wg.App.BodyColorGet()).Fn).Fn,
						// },
						gel.WidgetSize{
							// Size:   800,
							Widget: gel.EmptySpace(0, 0),
							// a.ButtonInset(0.25, a.Caption(title).Color(wg.BodyColorGet()).Fn).Fn,
						},
					},
				).Fn,
			).
			Flexed(
				1,
				wg.Inset(
					0.25,
					wg.Responsive(wg.Size.Load(), widget).Fn,
				).Fn,
			).Fn(gtx)
	}
}

func (wg *WalletGUI) SideBarButton(title, page string, index int) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		var scale float32
		scale = gel.Scales["H6"]
		var color string
		background := "Transparent"
		color = "DocText"
		var ins float32 = 0.5
		// var hl = false
		if wg.MainApp.ActivePageGet() == page || wg.MainApp.PreRendering {
			background = "PanelBg"
			scale = gel.Scales["H6"]
			color = "DocText"
			// ins = 0.5
			// hl = true
		}
		if title == " " {
			scale = gel.Scales["H6"] / 2
		}
		max := int(wg.MainApp.SideBarSize.V)
		if max > 0 {
			gtx.Constraints.Max.X = max
			gtx.Constraints.Min.X = max
		}
		// D.Ln("sideMAXXXXXX!!", max)
		return wg.Direction().E().Embed(
			wg.ButtonLayout(wg.sidebarButtons[index]).
				CornerRadius(scale).Corners(0).
				Background(background).
				Embed(
					wg.Inset(
						ins,
						func(gtx l.Context) l.Dimensions {
							return wg.H5(title).
								Color(color).
								Alignment(text.End).
								Fn(gtx)
						},
					).Fn,
				).
				SetClick(
					func() {
						if wg.MainApp.MenuOpen {
							wg.MainApp.MenuOpen = false
						}
						wg.MainApp.ActivePage(page)
					},
				).
				Fn,
		).
			Fn(gtx)
	}
}

func (wg *WalletGUI) PageTopBarButton(
	name string, index int, ico *[]byte, onClick func(string), app *gel.App,
	highlightColor string,
) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		background := "Transparent"
		// background := app.TitleBarBackgroundGet()
		color := app.MenuColorGet()
		
		if app.ActivePageGet() == name {
			color = "PanelText"
			// background = "scrim"
			background = "PanelBg"
		}
		// if name == "home" {
		// 	background = "scrim"
		// }
		if highlightColor != "" {
			color = highlightColor
		}
		ic := wg.Icon().
			Scale(gel.Scales["H5"]).
			Color(color).
			Src(ico).
			Fn
		return wg.Flex().Rigid(
			// wg.ButtonInset(0.25,
			wg.ButtonLayout(wg.buttonBarButtons[index]).
				CornerRadius(0).
				Embed(
					wg.Inset(
						0.375,
						ic,
					).Fn,
				).
				Background(background).
				SetClick(func() { onClick(name) }).
				Fn,
			// ).Fn,
		).Fn(gtx)
	}
}

func (wg *WalletGUI) StatusBarButton(
	name string,
	index int,
	ico *[]byte,
	onClick func(string),
	app *gel.App,
) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		background := app.StatusBarBackgroundGet()
		color := app.StatusBarColorGet()
		if app.ActivePageGet() == name {
			// background, color = color, background
			background = "PanelBg"
			// color = "Danger"
		}
		ic := wg.Icon().
			Scale(gel.Scales["H5"]).
			Color(color).
			Src(ico).
			Fn
		return wg.Flex().
			Rigid(
				wg.ButtonLayout(wg.statusBarButtons[index]).
					CornerRadius(0).
					Embed(
						wg.Inset(0.25, ic).Fn,
					).
					Background(background).
					SetClick(func() { onClick(name) }).
					Fn,
			).Fn(gtx)
	}
}

func (wg *WalletGUI) SetNodeRunState(b bool) {
	go func() {
		D.Ln("node run state is now", b)
		if b {
			wg.node.Start()
		} else {
			wg.node.Stop()
		}
	}()
}

func (wg *WalletGUI) SetWalletRunState(b bool) {
	go func() {
		D.Ln("node run state is now", b)
		if b {
			wg.wallet.Start()
		} else {
			wg.wallet.Stop()
		}
	}()
}

func (wg *WalletGUI) RunStatusPanel(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		t, f := &p9icons.Link, &p9icons.LinkOff
		var runningIcon *[]byte
		if wg.node.Running() {
			runningIcon = t
		} else {
			runningIcon = f
		}
		miningIcon := &p9icons.Mine
		if !wg.miner.Running() {
			miningIcon = &p9icons.NoMine
		}
		controllerIcon := &icons.NotificationSyncDisabled
		if wg.cx.Config.Controller.True() {
			controllerIcon = &icons.NotificationSync
		}
		discoverColor :=
			"DocText"
		discoverIcon :=
			&icons.DeviceWiFiTethering
		if wg.cx.Config.Discovery.False() {
			discoverIcon =
				&icons.CommunicationPortableWiFiOff
			discoverColor =
				"scrim"
		}
		clr := "scrim"
		if wg.cx.Config.Controller.True() {
			clr = "DocText"
		}
		clr2 := "DocText"
		if wg.cx.Config.GenThreads.V() == 0 {
			clr2 = "scrim"
		}
		// background := wg.App.StatusBarBackgroundGet()
		color := wg.MainApp.StatusBarColorGet()
		ic := wg.Icon().
			Scale(gel.Scales["H5"]).
			Color(color).
			Src(&icons.NavigationRefresh).
			Fn
		return wg.Flex().AlignMiddle().
			Rigid(
				wg.ButtonLayout(wg.statusBarButtons[0]).
					CornerRadius(0).
					Embed(
						wg.Inset(
							0.25,
							wg.Icon().
								Scale(gel.Scales["H5"]).
								Color("DocText").
								Src(runningIcon).
								Fn,
						).Fn,
					).
					Background(wg.MainApp.StatusBarBackgroundGet()).
					SetClick(
						func() {
							go func() {
								D.Ln("clicked node run control button", wg.node.Running())
								// wg.toggleNode()
								wg.unlockPassword.Wipe()
								wg.unlockPassword.Focus()
								if wg.node.Running() {
									if wg.wallet.Running() {
										wg.WalletWatcher.Q()
									}
									wg.node.Stop()
									wg.ready.Store(false)
									wg.stateLoaded.Store(false)
									wg.State.SetActivePage("home")
								} else {
									wg.node.Start()
									// wg.ready.Store(true)
									// wg.stateLoaded.Store(true)
								}
							}()
						},
					).
					Fn,
			).
			Rigid(
				wg.Inset(
					0.33,
					wg.Body1(fmt.Sprintf("%d", wg.State.bestBlockHeight.Load())).
						Font("go regular").TextScale(gel.Scales["Caption"]).
						Color("DocText").
						Fn,
				).Fn,
			).
			Rigid(
				wg.ButtonLayout(wg.statusBarButtons[6]).
					CornerRadius(0).
					Embed(
						wg.Inset(
							0.25,
							wg.Icon().
								Scale(gel.Scales["H5"]).
								Color(discoverColor).
								Src(discoverIcon).
								Fn,
						).Fn,
					).
					Background(wg.MainApp.StatusBarBackgroundGet()).
					SetClick(
						func() {
							go func() {
								wg.cx.Config.Discovery.Flip()
								_ = wg.cx.Config.WriteToFile(wg.cx.Config.ConfigFile.V())
								I.Ln("discover enabled:",
									wg.cx.Config.Discovery.True())
							}()
						},
					).
					Fn,
			).
			Rigid(
				wg.Inset(
					0.33,
					wg.Caption(fmt.Sprintf("%d LAN %d", len(wg.otherNodes), wg.peerCount.Load())).
						Font("go regular").
						Color("DocText").
						Fn,
				).Fn,
			).
			Rigid(
				wg.ButtonLayout(wg.statusBarButtons[7]).
					CornerRadius(0).
					Embed(
						wg.
							Inset(
								0.25, wg.
									Icon().
									Scale(gel.Scales["H5"]).
									Color(clr).
									Src(controllerIcon).Fn,
							).Fn,
					).
					Background(wg.MainApp.StatusBarBackgroundGet()).
					SetClick(
						func() {
							if wg.ChainClient != nil && !wg.ChainClient.Disconnected() {
								wg.cx.Config.Controller.Flip()
								I.Ln("controller running:",
									wg.cx.Config.Controller.True())
								var e error
								if e = wg.ChainClient.SetGenerate(
									wg.cx.Config.Controller.True(),
									wg.cx.Config.GenThreads.V(),
								); !E.Chk(e) {
								}
							}
							// // wg.toggleMiner()
							// go func() {
							// 	if wg.miner.Running() {
							// 		*wg.cx.Config.Generate = false
							// 		wg.miner.Stop()
							// 	} else {
							// 		wg.miner.Start()
							// 		*wg.cx.Config.Generate = true
							// 	}
							// 	save.Save(wg.cx.Config)
							// }()
						},
					).
					Fn,
			).
			Rigid(
				wg.ButtonLayout(wg.statusBarButtons[1]).
					CornerRadius(0).
					Embed(
						wg.Inset(
							0.25, wg.
								Icon().
								Scale(gel.Scales["H5"]).
								Color(clr2).
								Src(miningIcon).Fn,
						).Fn,
					).
					Background(wg.MainApp.StatusBarBackgroundGet()).
					SetClick(
						func() {
							// wg.toggleMiner()
							go func() {
								if wg.cx.Config.GenThreads.V() != 0 {
									if wg.miner.Running() {
										wg.cx.Config.Generate.F()
										wg.miner.Stop()
									} else {
										wg.miner.Start()
										wg.cx.Config.Generate.T()
									}
									_ = wg.cx.Config.WriteToFile(wg.cx.Config.ConfigFile.V())
								}
							}()
						},
					).
					Fn,
			).
			Rigid(
				func(gtx l.Context) l.Dimensions {
					return wg.incdecs["generatethreads"].
						// Color("DocText").
						// Background(wg.MainApp.StatusBarBackgroundGet()).
						Fn(gtx)
				},
			).
			Rigid(
				func(gtx l.Context) l.Dimensions {
					if !wg.wallet.Running() {
						return l.Dimensions{}
					}
					
					return wg.Flex().
						Rigid(
							wg.ButtonLayout(wg.statusBarButtons[2]).
								CornerRadius(0).
								Embed(
									wg.Inset(0.25, ic).Fn,
								).
								Background(wg.MainApp.StatusBarBackgroundGet()).
								SetClick(
									func() {
										D.Ln("clicked reset wallet button")
										go func() {
											var e error
											wasRunning := wg.wallet.Running()
											D.Ln("was running", wasRunning)
											if wasRunning {
												wg.wallet.Stop()
											}
											args := []string{
												os.Args[0],
												 "DD"+
												wg.cx.Config.DataDir.V(),
												"pipelog",
												"walletpass"+
												wg.cx.Config.WalletPass.V(),
												"wallet",
												"drophistory",
											}
											runner := exec.Command(args[0], args[1:]...)
											runner.Stderr = os.Stderr
											runner.Stdout = os.Stderr
											if e = wg.writeWalletCookie(); E.Chk(e) {
											}
											if e = runner.Run(); E.Chk(e) {
											}
											if wasRunning {
												wg.wallet.Start()
											}
										}()
									},
								).
								Fn,
						).Fn(gtx)
				},
			).
			Fn(gtx)
	}(gtx)
}

func (wg *WalletGUI) writeWalletCookie() (e error) {
	// for security with apps launching the wallet, the public password can be set with a file that is deleted after
	walletPassPath := filepath.Join(wg.cx.Config.DataDir.V(), wg.cx.ActiveNet.Name, "wp.txt")
	D.Ln("runner", walletPassPath)
	b := wg.cx.Config.WalletPass.Bytes()
	if e = ioutil.WriteFile(walletPassPath, b, 0700); E.Chk(e) {
	}
	D.Ln("created password cookie")
	return
}

//
// func (wg *WalletGUI) toggleNode() {
// 	if wg.node.Running() {
// 		wg.node.Stop()
// 		*wg.cx.Config.NodeOff = true
// 	} else {
// 		wg.node.Start()
// 		*wg.cx.Config.NodeOff = false
// 	}
// 	save.Save(wg.cx.Config)
// }
//
// func (wg *WalletGUI) startNode() {
// 	if !wg.node.Running() {
// 		wg.node.Start()
// 	}
// 	D.Ln("startNode")
// }
//
// func (wg *WalletGUI) stopNode() {
// 	if wg.wallet.Running() {
// 		wg.stopWallet()
// 		wg.unlockPassword.Wipe()
// 		// wg.walletLocked.Store(true)
// 	}
// 	if wg.node.Running() {
// 		wg.node.Stop()
// 	}
// 	D.Ln("stopNode")
// }
//
// func (wg *WalletGUI) toggleMiner() {
// 	if wg.miner.Running() {
// 		wg.miner.Stop()
// 		*wg.cx.Config.Generate = false
// 	}
// 	if !wg.miner.Running() && *wg.cx.Config.GenThreads > 0 {
// 		wg.miner.Start()
// 		*wg.cx.Config.Generate = true
// 	}
// 	save.Save(wg.cx.Config)
// }
//
// func (wg *WalletGUI) startMiner() {
// 	if *wg.cx.Config.GenThreads == 0 && wg.miner.Running() {
// 		wg.stopMiner()
// 		D.Ln("was zero threads")
// 	} else {
// 		wg.miner.Start()
// 		D.Ln("startMiner")
// 	}
// }
//
// func (wg *WalletGUI) stopMiner() {
// 	if wg.miner.Running() {
// 		wg.miner.Stop()
// 	}
// 	D.Ln("stopMiner")
// }
//
// func (wg *WalletGUI) toggleWallet() {
// 	if wg.wallet.Running() {
// 		wg.stopWallet()
// 		*wg.cx.Config.WalletOff = true
// 	} else {
// 		wg.startWallet()
// 		*wg.cx.Config.WalletOff = false
// 	}
// 	save.Save(wg.cx.Config)
// }
//
// func (wg *WalletGUI) startWallet() {
// 	if !wg.node.Running() {
// 		wg.startNode()
// 	}
// 	if !wg.wallet.Running() {
// 		wg.wallet.Start()
// 		wg.unlockPassword.Wipe()
// 		// wg.walletLocked.Store(false)
// 	}
// 	D.Ln("startWallet")
// }
//
// func (wg *WalletGUI) stopWallet() {
// 	if wg.wallet.Running() {
// 		wg.wallet.Stop()
// 		// wg.unlockPassword.Wipe()
// 		// wg.walletLocked.Store(true)
// 	}
// 	wg.unlockPassword.Wipe()
// 	D.Ln("stopWallet")
// }
