package gui

import (
	l "gioui.org/layout"
	"gioui.org/text"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/pod/pkg/gui/cfg"
	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func (wg *WalletGUI) GetAppWidget() (a *p9.App) {
	a = wg.th.App(*wg.size)
	wg.App = a
	wg.size = a.Size
	wg.config = cfg.New(wg.cx)
	wg.configs = wg.config.Config()
	a.Pages(map[string]l.Widget{
		"main": wg.Page("overview", p9.Widgets{
			p9.WidgetSize{Widget: wg.OverviewPage()},
		}),
		"send": wg.Page("send", p9.Widgets{
			p9.WidgetSize{Widget: wg.SendPage()},
		}),
		"settings": wg.Page("settings", p9.Widgets{
			p9.WidgetSize{Widget: func(gtx l.Context) l.Dimensions {
				return wg.configs.Widget(wg.config)(gtx)
			}},
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
									interrupt.Request()
								})).Color(wg.App.TitleBarColorGet()).TextScale(2).Text("yes!!!").Fn,
							).Fn,
					).
					Fn(gtx)
			},
			},
		}),
	})
	a.SideBar([]l.Widget{
		wg.SideBarButton("overview", "overview", 0),
		wg.SideBarButton("send", "send", 1),
		wg.SideBarButton("receive", "receive", 2),
		wg.SideBarButton("transactions", "transactions", 3),
		wg.SideBarButton("settings", "settings", 5),
		wg.SideBarButton("help", "help", 6),
		wg.SideBarButton("log", "log", 7),
		wg.SideBarButton("quit", "quit", 8),
	})
	a.ButtonBar([]l.Widget{
		wg.PageTopBarButton("help", 0, icons.ActionHelp),
		wg.PageTopBarButton("log", 1, icons.ActionList),
		wg.PageTopBarButton("settings", 2, icons.ActionSettings),
		wg.PageTopBarButton("quit", 3, icons.ActionExitToApp),
	})
	a.StatusBar([]l.Widget{
		wg.StatusBarButton("help", 0, icons.ActionHelp),
		wg.StatusBarButton("log", 1, icons.ActionList),
		wg.StatusBarButton("settings", 2, icons.ActionSettings),
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
				var inPad, outPad float32 = 0.5, 0.25
				if *wg.Size >= 800 {
					inPad, outPad = 0.75, 0
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

func (wg *WalletGUI) PageTopBarButton(name string, index int, ico []byte) func(gtx l.Context) l.Dimensions {
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

func (wg *WalletGUI) StatusBarButton(name string, index int, ico []byte) func(gtx l.Context) l.Dimensions {
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
						ic,
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
