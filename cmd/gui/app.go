package gui

import (
	"fmt"

	l "gioui.org/layout"
	"gioui.org/text"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/gui/cfg"
	p9icons "github.com/p9c/pod/pkg/gui/ico/svg"
	"github.com/p9c/pod/pkg/gui/p9"
)

func (wg *WalletGUI) GetAppWidget() (a *p9.App) {
	a = wg.th.App(*wg.size)
	wg.App = a
	wg.App.ThemeHook(func() {
		Debug("theme hook")
		Debug(wg.bools)
		*wg.cx.Config.DarkTheme = *wg.Dark
		a := wg.configs["config"]["DarkTheme"].Slot.(*bool)
		*a = *wg.Dark
		if wgb, ok := wg.config.Bools["DarkTheme"]; ok {
			wgb.Value(*wg.Dark)
		}
		save.Pod(wg.cx.Config)

	})
	wg.size = a.Size
	wg.config = cfg.New(wg.cx, wg.th)
	wg.configs = wg.config.Config()
	a.Pages(map[string]l.Widget{
		"main": wg.Page("overview", p9.Widgets{
			// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
			p9.WidgetSize{Widget: wg.OverviewPage()},
		}),
		"send": wg.Page("send", p9.Widgets{
			// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
			p9.WidgetSize{Widget: wg.SendPage()},
		}),
		"receive": wg.Page("receive", p9.Widgets{
			// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
			p9.WidgetSize{Widget: wg.ReceivePage()},
		}),
		"transactions": wg.Page("receive", p9.Widgets{
			// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
			p9.WidgetSize{Widget: wg.TransactionsPage()},
		}),
		"settings": wg.Page("settings", p9.Widgets{
			// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
			p9.WidgetSize{Widget: func(gtx l.Context) l.Dimensions {
				return wg.configs.Widget(wg.config)(gtx)
			}},
		}),
		"console": wg.Page("console", p9.Widgets{
			// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
			p9.WidgetSize{Widget: wg.ConsolePage()},
		}),
		"help": wg.Page("help", p9.Widgets{
			p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
		}),
		"log": wg.Page("log", p9.Widgets{
			p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
		}),
		"quit": wg.Page("quit", p9.Widgets{
			p9.WidgetSize{Widget: func(gtx l.Context) l.Dimensions {
				return wg.th.VFlex().
					SpaceEvenly().
					// AlignMiddle().
					Rigid(
						wg.th.H4("are you sure?").Color(wg.App.BodyColorGet()).Alignment(text.Middle).Fn,
					).
					Rigid(
						wg.th.Flex().
							SpaceEvenly().
							Rigid(
								wg.th.Button(wg.clickables["quit"].SetClick(func() {
									close(wg.quit)
								})).Color(wg.App.TitleBarColorGet()).TextScale(2).Text("yes!!!").Fn,
							).Fn,
					).
					Fn(gtx)
			},
			},
		}),
		"goroutines": wg.Page("log", p9.Widgets{
			p9.WidgetSize{Widget: func(gtx l.Context) l.Dimensions {
				le := func(gtx l.Context, index int) l.Dimensions {
					return wg.State.goroutines[index](gtx)
				}
				return func(gtx l.Context) l.Dimensions {
					return wg.Inset(0.25,
						wg.Fill("DocBg",
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
				// wg.RunCommandChan <- "stop"
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
			p9.WidgetSize{Widget: wg.th.VFlex().SpaceAround().AlignMiddle().Rigid(wg.th.H1("mining").Alignment(text.Middle).Fn).Fn},
		}),
	})
	a.SideBar([]l.Widget{
		wg.SideBarButton("overview", "main", 0),
		wg.SideBarButton("send", "send", 1),
		wg.SideBarButton("receive", "receive", 2),
		wg.SideBarButton("history", "transactions", 3),
		wg.SideBarButton("settings", "settings", 5),
		wg.SideBarButton("mining", "mining", 6),
		wg.SideBarButton("help", "help", 7),
		wg.SideBarButton("log", "log", 8),
		wg.SideBarButton("quit", "quit", 9),
	})
	a.ButtonBar([]l.Widget{
		wg.PageTopBarButton("goroutines", 0, &icons.ActionBugReport),
		wg.PageTopBarButton("help", 1, &icons.ActionHelp),
		wg.PageTopBarButton("console", 2, &icons.MapsLocalHotel),
		wg.PageTopBarButton("settings", 3, &icons.ActionSettings),
		//wg.PageTopBarButton("quit", 4, &icons.ActionExitToApp),
	})
	a.StatusBar([]l.Widget{
		func(gtx l.Context) l.Dimensions { return wg.RunStatusPanel(gtx) },
		wg.th.Flex().Rigid(
			wg.StatusBarButton("log", 1, &icons.ActionList),
		).Rigid(
			wg.StatusBarButton("settings", 2, &icons.ActionSettings),
		).Fn,
	})
	return
}

func (wg *WalletGUI) Page(title string, widget p9.Widgets) func(gtx l.Context) l.Dimensions {
	a := wg.th
	return func(gtx l.Context) l.Dimensions {
		return a.Fill(wg.BodyBackgroundGet(),
			a.VFlex().
				SpaceEvenly().
				Rigid(
					a.Responsive(*wg.Size, p9.Widgets{
						p9.WidgetSize{
							Widget: a.Inset(0.25, a.H5(title).Color(wg.BodyColorGet()).Fn).Fn,
						},
						p9.WidgetSize{
							Size:   800,
							Widget: p9.EmptySpace(0, 0),
							// a.Inset(0.25, a.Caption(title).Color(wg.BodyColorGet()).Fn).Fn,
						},
					}).Fn,
				).
				Flexed(1,
					a.Inset(0.25,
						a.Responsive(*wg.Size, widget).Fn,
					).Fn,
				).Fn,
		).Fn(gtx)
	}
}

func (wg *WalletGUI) SideBarButton(title, page string, index int) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		gtx.Constraints.Max.X = int(wg.TextSize.Scale(12).V)
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return wg.ButtonLayout(wg.sidebarButtons[index]).Embed(
			func(gtx l.Context) l.Dimensions {
				background := "Transparent"
				color := "DocText"
				if wg.ActivePageGet() == page {
					background = "PanelBg"
					color = "PanelText"
				}
				var inPad, outPad float32 = 0.25, 0.25
				if *wg.Size >= 800 {
					inPad, outPad = 0.5, 0
				}
				return wg.Inset(outPad,
					wg.Fill(background,
						wg.Flex().
							Flexed(1,
								wg.Inset(inPad,
									wg.H6(title).
										Color(color).
										Fn,
								).Fn,
							).Fn,
					).Fn,
				).Fn(gtx)
			},
		).
			Background("Transparent").
			SetClick(
				func() {
					if wg.MenuOpen {
						wg.MenuOpen = false
					}
					wg.ActivePage(page)
				}).
			Fn(gtx)
	}
}

func (wg *WalletGUI) PageTopBarButton(name string, index int, ico *[]byte) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		background := wg.TitleBarBackgroundGet()
		color := wg.MenuColorGet()
		if wg.ActivePageGet() == name {
			color = "PanelText"
			background = "PanelBg"
		}
		ic := wg.Icon().
			Scale(p9.Scales["H5"]).
			Color(color).
			Src(ico).
			Fn
		return wg.Flex().Rigid(
			// wg.Inset(0.25,
			wg.ButtonLayout(wg.buttonBarButtons[index]).
				CornerRadius(0).
				Embed(
					wg.Inset(0.375,
						ic,
					).Fn,
				).
				Background(background).
				SetClick(
					func() {
						if wg.MenuOpen {
							wg.MenuOpen = false
						}
						wg.ActivePage(name)
					}).
				Fn,
			// ).Fn,
		).Fn(gtx)
	}
}

func (wg *WalletGUI) StatusBarButton(name string, index int, ico *[]byte) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		background := wg.StatusBarBackgroundGet()
		color := wg.StatusBarColorGet()
		ic := wg.Icon().
			Scale(p9.Scales["H5"]).
			Color(color).
			Src(ico).
			Fn
		return wg.Flex().
			Rigid(
				wg.ButtonLayout(wg.statusBarButtons[index]).
					CornerRadius(0).
					Embed(
						wg.th.Inset(0.25, ic).Fn,
					).
					Background(background).
					SetClick(
						func() {
							if wg.MenuOpen {
								wg.MenuOpen = false
							}
							wg.ActivePage(name)
						}).
					Fn,
			).Fn(gtx)
	}
}

func (wg *WalletGUI) SetRunState(b bool) {
	go func() {
		Debug("run state is now", b)
		if b {
			wg.RunCommandChan <- "run"
			// wg.running = b
		} else {
			wg.RunCommandChan <- "stop"
			// wg.running = b
		}
	}()
}

func (wg *WalletGUI) RunStatusPanel(gtx l.Context) l.Dimensions {
	t, f := &p9icons.Link, &p9icons.LinkOff
	var runningIcon *[]byte
	if wg.running {
		runningIcon = t
	} else {
		runningIcon = f
	}
	miningIcon := &p9icons.Mine
	if !wg.mining {
		miningIcon = &p9icons.NoMine
	}
	return wg.th.Flex().
		Rigid(
			wg.th.ButtonLayout(wg.statusBarButtons[0]).
				CornerRadius(0).
				Embed(
					wg.th.Inset(0.25,
						wg.th.Icon().
							Scale(p9.Scales["H5"]).
							Color("DocText").
							Src(runningIcon).
							Fn,
					).Fn,
				).
				Background("DocBg").
				SetClick(
					func() {
						wg.SetRunState(!wg.running)
					}).
				Fn,
		).
		Rigid(
			wg.th.Inset(0.25,
				p9.If(wg.running,
					wg.th.Indefinite().Scale(p9.Scales["H5"]).Fn,
					wg.th.Icon().
						Scale(p9.Scales["H5"]).
						Color("Primary").
						Src(&icons.ActionCheckCircle).
						Fn,
				),
			).Fn,
		).
		Rigid(wg.th.
			Inset(0.25,
				wg.Icon().
					Scale(p9.Scales["H5"]).
					Color("DocText").
					Src(&icons.DeviceWidgets).
					Fn,
			).
			Fn,
		).
		Rigid(
			wg.th.Inset(0.33,
				wg.th.Body1(fmt.Sprintf("%-8d", wg.State.bestBlockHeight)).
					Font("go regular").
					Color("DocText").
					Fn,
			).Fn,
		).
		Rigid(
			wg.th.ButtonLayout(wg.statusBarButtons[3]).
				CornerRadius(0).
				Embed(wg.th.
					Inset(0.25, wg.th.
						Icon().
						Scale(p9.Scales["H5"]).
						Color("DocText").
						Src(miningIcon).Fn,
					).Fn,
				).
				Background("DocBg").
				SetClick(
					func() {
						Debug("clicked miner control stop/start button")
						wg.mining = !wg.mining
						// wg.SetRunState(!wg.running)
					}).
				Fn,
		).
		Rigid(
			wg.incdecs["generatethreads"].
				SetColor("DocText").
				SetBackground("DocBg").
				Fn,
		).
		Fn(gtx)
}
