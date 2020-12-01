package gui

import (
	"golang.org/x/exp/shiny/materialdesign/icons"

	l "gioui.org/layout"
	"gioui.org/text"
	"github.com/p9c/pod/app/save"
	p9icons "github.com/p9c/pod/pkg/gui/ico/svg"
	"github.com/p9c/pod/pkg/gui/p9"
)

//
// func (wg *WalletGUI) WalletUnlockPage(gtx l.Context) l.Dimensions {
// 	return wg.th.Fill("PanelBg",
// 		wg.th.Inset(0.5,
// 			wg.th.H4("unlock wallet").Fn,
// 			// p9.EmptyMaxWidth(),
// 		).Fn,
// 	).Fn(gtx)
// }

func (wg *WalletGUI) getWalletUnlockAppWidget() (a *p9.App) {
	a = wg.th.App(wg.w["main"].Width)
	wg.unlockPage = a
	password := ""
	wg.unlockPassword = wg.th.Password("", &password, "Danger", "DocText", 26, func(pass string) {
		Debug("entered password", pass)
	})
	wg.unlockPage.ThemeHook(func() {
		Debug("theme hook")
		// Debug(wg.bools)
		*wg.cx.Config.DarkTheme = *wg.th.Dark
		a := wg.configs["config"]["DarkTheme"].Slot.(*bool)
		*a = *wg.th.Dark
		if wgb, ok := wg.config.Bools["DarkTheme"]; ok {
			wgb.Value(*wg.th.Dark)
		}
		save.Pod(wg.cx.Config)
	})
	wg.size = a.Size
	a.Pages(map[string]l.Widget{
		"main": wg.Page("overview", p9.Widgets{
			p9.WidgetSize{Widget:
			func(gtx l.Context) l.Dimensions {
				return wg.th.Flex().
					SpaceEvenly().
					AlignMiddle().
					Flexed(1,
						wg.th.VFlex().Flexed(0.5, p9.EmptyMaxHeight()).Rigid(
							wg.th.Flex().
								SpaceEvenly().
								AlignMiddle().
								Flexed(1,
									wg.th.Flex().AlignMiddle().
										Flexed(0.5, p9.EmptyMaxWidth()).
										Rigid(
											wg.th.VFlex().
												Rigid(
													wg.th.Flex().AlignEnd().
														Rigid(
															wg.th.Icon().Scale(p9.Scales["H1"]).Color("Danger").Src(&icons.ActionLock).Fn,
														).
														Rigid(
															wg.th.H2("locked").Color("Danger").Fn,
														).
														Fn,
												).
												Rigid(wg.th.Inset(0.5, p9.EmptySpace(0, 0)).Fn).
												Rigid(
													wg.unlockPassword.Fn,
												).
												// Rigid(
												// 	// wg.th.VFlex().
												// 	wg.th.Body1(fmt.Sprintf("lock automatically in %v",
												// 		time.Second*time.Duration(wg.intSliders["lockTimeout"].GetValue()))).Color("Danger").Fn,
												//
												// ).
												// Rigid(
												// 	func(gtx l.Context) l.Dimensions {
												// 		cs := gtx.Constraints
												// 		cs.Max.X = int(wg.th.TextSize.Scale(24).V)
												// 		// cs.Max.X = 200
												// 		cs.Min.X = cs.Max.X
												// 		cs.Constrain(image.Point{X: cs.Max.X, Y: cs.Max.Y})
												// 		return wg.intSliders["lockTimeout"].Fn(gtx)
												// 	},
												// ).
												Fn,
										).
										Flexed(0.5, p9.EmptyMaxWidth()).Fn,
								).
								Fn,
						).Flexed(0.5, p9.EmptyMaxHeight()).Fn,
					).
					Fn(gtx)
			},
			},
			// p9.WidgetSize{Widget: p9.EmptyMaxWidth()},
			// p9.WidgetSize{Widget: wg.OverviewPage()},
		}),
		// "send": wg.Page("send", p9.Widgets{
		// 	// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
		// 	p9.WidgetSize{Widget: wg.SendPage()},
		// }),
		// "receive": wg.Page("receive", p9.Widgets{
		// 	// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
		// 	p9.WidgetSize{Widget: wg.ReceivePage()},
		// }),
		// "history": wg.Page("history", p9.Widgets{
		// 	// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
		// 	p9.WidgetSize{Widget: wg.HistoryPage()},
		// }),
		"settings": wg.Page("settings", p9.Widgets{
			// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
			p9.WidgetSize{Widget: func(gtx l.Context) l.Dimensions {
				return wg.configs.Widget(wg.config)(gtx)
			}},
		}),
		"console": wg.Page("console", p9.Widgets{
			// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
			p9.WidgetSize{Widget: wg.console.Fn},
		}),
		"help": wg.Page("help", p9.Widgets{
			p9.WidgetSize{Widget: p9.EmptyMaxWidth()},
		}),
		"log": wg.Page("log", p9.Widgets{
			p9.WidgetSize{Widget: p9.EmptyMaxWidth()},
		}),
		"quit": wg.Page("quit", p9.Widgets{
			p9.WidgetSize{Widget: func(gtx l.Context) l.Dimensions {
				return wg.th.VFlex().
					SpaceEvenly().
					AlignMiddle().
					Rigid(
						wg.th.H4("are you sure?").Color(wg.unlockPage.BodyColorGet()).Alignment(text.Middle).Fn,
					).
					Rigid(
						wg.th.Flex().
							// SpaceEvenly().
							Flexed(0.5, p9.EmptyMaxWidth()).
							Rigid(
								wg.th.Button(wg.clickables["quit"].SetClick(func() {
									close(wg.quit)
								})).Color("Light").TextScale(2).Text("yes!!!").Fn,
							).
							Flexed(0.5, p9.EmptyMaxWidth()).
							Fn,
					).
					Fn(gtx)
			},
			},
		}),
		"goroutines": wg.Page("log", p9.Widgets{
			// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},

			p9.WidgetSize{Widget: func(gtx l.Context) l.Dimensions {
				le := func(gtx l.Context, index int) l.Dimensions {
					return wg.State.goroutines[index](gtx)
				}
				return func(gtx l.Context) l.Dimensions {
					return wg.th.Inset(0.25,
						wg.th.Fill("DocBg",
							wg.lists["recent"].
								Vertical().
								// Background("DocBg").Color("DocText").Active("Primary").
								Length(len(wg.State.goroutines)).
								ListElement(le).
								Fn,
						).Fn,
					).
						Fn(gtx)
				}(gtx)
				// wg.ShellRunCommandChan <- "stop"
				// consume.Kill(wg.Worker)
				// consume.Kill(wg.cx.StateCfg.Miner)
				// close(wg.cx.NodeKill)
				// close(wg.cx.KillAll)
				// time.Sleep(time.Second*3)
				// interrupt.Request()
				// os.Exit(0)
				// return l.Dimensions{}
			}},
		}),
		"mining": wg.Page("mining", p9.Widgets{
			p9.WidgetSize{Widget: func(gtx l.Context) l.Dimensions {
				return wg.th.VFlex().
					AlignMiddle().
					SpaceSides().
					Rigid(
						wg.th.Flex().
							Flexed(0.5, p9.EmptyMaxWidth()).
							Rigid(
								wg.th.H1("Mining").Fn,
							).
							Flexed(0.5, p9.EmptyMaxWidth()).
							Fn,
					).
					Fn(gtx)
			},
			},
		}),
		"explorer": wg.Page("explorer", p9.Widgets{
			p9.WidgetSize{Widget: func(gtx l.Context) l.Dimensions {
				return wg.th.VFlex().
					AlignMiddle().
					SpaceSides().
					Rigid(
						wg.th.Flex().
							Flexed(0.5, p9.EmptyMaxWidth()).
							Rigid(
								wg.th.H1("explorer").Fn,
							).
							Flexed(0.5, p9.EmptyMaxWidth()).
							Fn,
					).
					Fn(gtx)
			},
			},
		}),
	})
	// a.SideBar([]l.Widget{
	// 	wg.SideBarButton("overview", "main", 0),
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
	a.ButtonBar([]l.Widget{
		wg.PageTopBarButton("main", 4, &icons.ActionLock, func(name string) {
			wg.unlockPage.ActivePage(name)
		}, wg.unlockPage),
		wg.PageTopBarButton("console", 2, &p9icons.Terminal, func(name string) {
			wg.unlockPage.ActivePage(name)
		}, wg.unlockPage),
		wg.PageTopBarButton("goroutines", 0, &icons.ActionBugReport, func(name string) {
			wg.unlockPage.ActivePage(name)
		}, wg.unlockPage),
		wg.PageTopBarButton("help", 1, &icons.ActionHelp, func(name string) {
			wg.unlockPage.ActivePage(name)
		}, wg.unlockPage),
		wg.PageTopBarButton("quit", 3, &icons.ActionExitToApp, func(name string) {
			wg.unlockPage.ActivePage(name)
		}, wg.unlockPage),
	})
	a.StatusBar([]l.Widget{
		// func(gtx l.Context) l.Dimensions { return wg.RunStatusPanel(gtx) },
		wg.RunStatusPanel,
		wg.th.Flex().
			Flexed(1, p9.EmptyMaxWidth()).
			Rigid(
				wg.StatusBarButton("log", 4, &icons.ActionList, func(name string) {
					Debug("click on button", name)
					wg.unlockPage.ActivePage(name)
				}, wg.unlockPage),
			).
			Rigid(
				wg.StatusBarButton("settings", 5, &icons.ActionSettings, func(name string) {
					wg.unlockPage.ActivePage(name)
				}, wg.unlockPage),
			).
			Fn,
	})
	a.AddOverlay(wg.toasts.DrawToasts())
	a.AddOverlay(wg.dialog.DrawDialog())
	return
}
