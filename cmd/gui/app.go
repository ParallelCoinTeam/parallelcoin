package gui

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	
	uberatomic "go.uber.org/atomic"
	"golang.org/x/exp/shiny/materialdesign/icons"
	
	l "gioui.org/layout"
	"gioui.org/text"
	
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/gui"
	"github.com/p9c/pod/pkg/gui/cfg"
	p9icons "github.com/p9c/pod/pkg/gui/ico/svg"
)

func (wg *WalletGUI) GetAppWidget() (a *gui.App) {
	a = wg.App(&wg.Window.Width, uberatomic.NewString("home"), wg.invalidate).SetMainDirection(l.W)
	wg.MainApp = a
	wg.MainApp.ThemeHook(
		func() {
			Debug("theme hook")
			// Debug(wg.bools)
			// wg.Colors.Lock()
			*wg.cx.Config.DarkTheme = *wg.Dark
			// a := wg.configs["config"]["DarkTheme"].Slot.(*bool)
			// *a = *wg.Dark
			// if wgb, ok := wg.config.Bools["DarkTheme"]; ok {
			// 	wgb.Value(*wg.Dark)
			// }
			// wg.Colors.Unlock()
			save.Pod(wg.cx.Config)
			wg.RecentTransactions(10, "recent")
			wg.RecentTransactions(-1, "history")
		},
	)
	wg.config = cfg.New(wg.cx, wg.Window)
	wg.configs = wg.config.Config()
	a.Pages(
		map[string]l.Widget{
			"home": wg.Page(
				"home", gui.Widgets{
					// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
					gui.WidgetSize{Widget: wg.OverviewPage()},
				},
			),
			"send": wg.Page(
				"send", gui.Widgets{
					// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
					gui.WidgetSize{Widget: wg.SendPage()},
				},
			),
			"receive": wg.Page(
				"receive", gui.Widgets{
					// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
					gui.WidgetSize{Widget: wg.ReceivePage()},
				},
			),
			"history": wg.Page(
				"history", gui.Widgets{
					// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
					gui.WidgetSize{Widget: wg.HistoryPage()},
				},
			),
			"settings": wg.Page(
				"settings", gui.Widgets{
					// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
					gui.WidgetSize{
						Widget: func(gtx l.Context) l.Dimensions {
							return wg.configs.Widget(wg.config)(gtx)
						},
					},
				},
			),
			"console": wg.Page(
				"console", gui.Widgets{
					// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
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
									wg.H4("are you sure?").Color(wg.MainApp.BodyColorGet()).Alignment(text.Middle).Fn,
								).
								Rigid(
									wg.Flex().
										// SpaceEvenly().
										Flexed(0.5, gui.EmptyMaxWidth()).
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
												"yes!!!").Fn,
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
					gui.WidgetSize{
						Widget: a.Placeholder("mining"),
					},
				},
			),
			"explorer": wg.Page(
				"explorer", gui.Widgets{
					gui.WidgetSize{
						Widget: a.Placeholder("explorer"),
					},
				},
			),
		},
	)
	a.SideBar(
		[]l.Widget{
			wg.SideBarButton(" ", " ", 11),
			wg.SideBarButton("home", "home", 0),
			wg.SideBarButton("send", "send", 1),
			wg.SideBarButton("receive", "receive", 2),
			wg.SideBarButton("history", "history", 3),
			wg.SideBarButton("explorer", "explorer", 6),
			wg.SideBarButton("mining", "mining", 7),
			wg.SideBarButton("console", "console", 9),
			wg.SideBarButton("settings", "settings", 5),
			wg.SideBarButton("log", "log", 10),
			wg.SideBarButton("help", "help", 8),
			wg.SideBarButton(" ", " ", 11),
			// wg.SideBarButton("quit", "quit", 11),
		},
	)
	a.ButtonBar(
		[]l.Widget{
			
			// gui.EmptyMaxWidth(),
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
					// wg.walletLocked.Store(true)
					wg.wallet.Stop()
					wg.unlockPage.ActivePage(name)
					wg.stateLoaded.Store(false)
					wg.txReady.Store(false)
				}, a, "green",
			),
			wg.Flex().Rigid(wg.Inset(0.5, gui.EmptySpace(0, 0)).Fn).Fn,
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
			// wg.Inset(0.5, gui.EmptySpace(0, 0)).Fn,
			// wg.Inset(0.5, gui.EmptySpace(0, 0)).Fn,
			wg.RunStatusPanel,
		},
		[]l.Widget{
			// gui.EmptyMaxWidth(),
			wg.StatusBarButton(
				"console", 2, &p9icons.Terminal, func(name string) {
					wg.MainApp.ActivePage(name)
				}, a,
			),
			wg.StatusBarButton(
				"log", 4, &icons.ActionList, func(name string) {
					Debug("click on button", name)
					if wg.MainApp.MenuOpen {
						wg.MainApp.MenuOpen = false
					}
					wg.MainApp.ActivePage(name)
				}, a,
			),
			wg.StatusBarButton(
				"settings", 5, &icons.ActionSettings, func(name string) {
					Debug("click on button", name)
					if wg.MainApp.MenuOpen {
						wg.MainApp.MenuOpen = false
					}
					wg.MainApp.ActivePage(name)
				}, a,
			),
			// wg.Inset(0.5, gui.EmptySpace(0, 0)).Fn,
		},
	)
	// a.PushOverlay(wg.toasts.DrawToasts())
	// a.PushOverlay(wg.dialog.DrawDialog())
	return
}

func (wg *WalletGUI) Page(title string, widget gui.Widgets) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		return wg.VFlex().
			// SpaceEvenly().
			Rigid(
				wg.Responsive(
					*wg.Size, gui.Widgets{
						// p9.WidgetSize{
						// 	Widget: a.ButtonInset(0.25, a.H5(title).Color(wg.App.BodyColorGet()).Fn).Fn,
						// },
						gui.WidgetSize{
							// Size:   800,
							Widget: gui.EmptySpace(0, 0),
							// a.ButtonInset(0.25, a.Caption(title).Color(wg.BodyColorGet()).Fn).Fn,
						},
					},
				).Fn,
			).
			Flexed(
				1,
				wg.Inset(
					0.25,
					wg.Responsive(*wg.Size, widget).Fn,
				).Fn,
			).Fn(gtx)
	}
}

func (wg *WalletGUI) SideBarButton(title, page string, index int) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		var scale float32
		scale = gui.Scales["H6"]
		var color string
		background := "Transparent"
		color = "DocText"
		var font string
		font = "plan9"
		var ins float32 = 0.5
		// var hl = false
		if wg.MainApp.ActivePageGet() == page || wg.MainApp.PreRendering {
			background = "PanelBg"
			scale = gui.Scales["H6"]
			color = "DocText"
			font = "plan9"
			ins = 1
			// hl = true
		}
		if title == " " {
			scale = gui.Scales["H6"] / 2
		}
		max := int(wg.MainApp.SideBarSize.V)
		if max > 0 {
			gtx.Constraints.Max.X = max
			gtx.Constraints.Min.X = max
		}
		// Debug("sideMAXXXXXX!!", max)
		return wg.Direction().E().Embed(
			wg.ButtonLayout(wg.sidebarButtons[index]).
				CornerRadius(scale).Corners(gui.NW | gui.SW | gui.NE).
				Background(background).
				Embed(
					wg.Inset(ins,
						func(gtx l.Context) l.Dimensions {
							return wg.Label().
								Font(font).
								Text(title).
								TextScale(scale).
								Color(color).Alignment(text.End).
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
	name string, index int, ico *[]byte, onClick func(string), app *gui.App,
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
			Scale(gui.Scales["H5"]).
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
	app *gui.App,
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
			Scale(gui.Scales["H5"]).
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
		Debug("node run state is now", b)
		if b {
			wg.node.Start()
		} else {
			wg.node.Stop()
		}
	}()
}

func (wg *WalletGUI) SetWalletRunState(b bool) {
	go func() {
		Debug("node run state is now", b)
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
		return wg.Flex().AlignMiddle().
			Rigid(
				wg.ButtonLayout(wg.statusBarButtons[0]).
					CornerRadius(0).
					Embed(
						wg.Inset(
							0.25,
							wg.Icon().
								Scale(gui.Scales["H5"]).
								Color("DocText").
								Src(runningIcon).
								Fn,
						).Fn,
					).
					Background(wg.MainApp.StatusBarBackgroundGet()).
					SetClick(
						func() {
							go func() {
								Debug("clicked node run control button", wg.node.Running())
								// wg.toggleNode()
								if wg.node.Running() {
									if wg.wallet.Running() {
										wg.wallet.Stop()
									}
									wg.node.Stop()
								} else {
									wg.node.Start()
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
						Font("go regular").TextScale(gui.Scales["Caption"]).
						Color("DocText").
						Fn,
				).Fn,
			).
			Rigid(
				wg.ButtonLayout(wg.statusBarButtons[1]).
					CornerRadius(0).
					Embed(
						func(gtx l.Context) l.Dimensions {
							clr := "DocText"
							if *wg.cx.Config.GenThreads == 0 {
								clr = "scrim"
							}
							return wg.
								Inset(
									0.25, wg.
										Icon().
										Scale(gui.Scales["H5"]).
										Color(clr).
										Src(miningIcon).Fn,
								).Fn(gtx)
						},
					).
					Background(wg.MainApp.StatusBarBackgroundGet()).
					SetClick(
						func() {
							// wg.toggleMiner()
							go func() {
								if wg.miner.Running() {
									*wg.cx.Config.Generate = false
									wg.miner.Stop()
								} else {
									wg.miner.Start()
									*wg.cx.Config.Generate = true
								}
								save.Pod(wg.cx.Config)
							}()
						},
					).
					Fn,
			).
			Rigid(
				wg.incdecs["generatethreads"].
					Color("DocText").
					Background(wg.MainApp.StatusBarBackgroundGet()).
					Fn,
			).
			Rigid(
				func(gtx l.Context) l.Dimensions {
					if !wg.wallet.Running() {
						return l.Dimensions{}
					}
					// background := wg.App.StatusBarBackgroundGet()
					color := wg.MainApp.StatusBarColorGet()
					ic := wg.Icon().
						Scale(gui.Scales["H5"]).
						Color(color).
						Src(&icons.NavigationRefresh).
						Fn
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
										Debug("clicked reset wallet button")
										go func() {
											var err error
											wasRunning := wg.wallet.Running()
											Debug("was running", wasRunning)
											if wasRunning {
												wg.wallet.Stop()
											}
											args := []string{
												os.Args[0],
												"-D",
												*wg.cx.Config.DataDir,
												"--pipelog",
												"--walletpass",
												*wg.cx.Config.WalletPass,
												"wallet",
												"drophistory",
											}
											runner := exec.Command(args[0], args[1:]...)
											runner.Stderr = os.Stderr
											runner.Stdout = os.Stderr
											if err = wg.writeWalletCookie(); Check(err) {
											}
											if err = runner.Run(); Check(err) {
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

func (wg *WalletGUI) writeWalletCookie() (err error) {
	// for security with apps launching the wallet, the public password can be set with a file that is deleted after
	walletPassPath := *wg.cx.Config.DataDir + slash + wg.cx.ActiveNet.Params.Name + slash + "wp.txt"
	Debug("runner", walletPassPath)
	wp := *wg.cx.Config.WalletPass
	b := []byte(wp)
	if err = ioutil.WriteFile(walletPassPath, b, 0700); Check(err) {
	}
	Debug("created password cookie")
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
// 	save.Pod(wg.cx.Config)
// }
//
// func (wg *WalletGUI) startNode() {
// 	if !wg.node.Running() {
// 		wg.node.Start()
// 	}
// 	Debug("startNode")
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
// 	Debug("stopNode")
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
// 	save.Pod(wg.cx.Config)
// }
//
// func (wg *WalletGUI) startMiner() {
// 	if *wg.cx.Config.GenThreads == 0 && wg.miner.Running() {
// 		wg.stopMiner()
// 		Debug("was zero threads")
// 	} else {
// 		wg.miner.Start()
// 		Debug("startMiner")
// 	}
// }
//
// func (wg *WalletGUI) stopMiner() {
// 	if wg.miner.Running() {
// 		wg.miner.Stop()
// 	}
// 	Debug("stopMiner")
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
// 	save.Pod(wg.cx.Config)
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
// 	Debug("startWallet")
// }
//
// func (wg *WalletGUI) stopWallet() {
// 	if wg.wallet.Running() {
// 		wg.wallet.Stop()
// 		// wg.unlockPassword.Wipe()
// 		// wg.walletLocked.Store(true)
// 	}
// 	wg.unlockPassword.Wipe()
// 	Debug("stopWallet")
// }
