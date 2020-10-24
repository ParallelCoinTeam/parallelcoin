package gui

import (
	l "gioui.org/layout"
	"gioui.org/text"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func (ng *WalletGUI) GetAppWidget() (a *p9.App) {
	a = ng.th.App(*ng.size)
	ng.App = a
	ng.size = a.Size
	a.Pages(map[string]l.Widget{
		"main": ng.Page("overview", p9.Widgets{
			p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
		}),
		"settings": ng.Page("settings", p9.Widgets{
			p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
		}),
		"help": ng.Page("help", p9.Widgets{
			p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
		}),
		"log": ng.Page("log", p9.Widgets{
			p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
		}),
		"quit": ng.Page("quit", p9.Widgets{
			p9.WidgetSize{Widget: a.VFlex().
				SpaceEvenly().
				// AlignMiddle().
				Rigid(
					a.H4("are you sure?").Color(ng.BodyColorGet()).Alignment(text.Middle).Fn,
				).
				Rigid(
					a.Flex().
						SpaceEvenly().
						Rigid(
							a.Button(ng.quitClickable.SetClick(func() {
								interrupt.Request()
							})).TextScale(2).Text("yes").Fn,
						).Fn,
				).
				Fn},
		}),
	})
	a.SideBar([]l.Widget{
		ng.SideBarButton("overview", "overview", 0),
		ng.SideBarButton("send", "send", 1),
		ng.SideBarButton("receive", "receive", 2),
		ng.SideBarButton("transactions", "transactions", 3),
		ng.SideBarButton("settings", "settings", 5),
		ng.SideBarButton("help", "help", 6),
		ng.SideBarButton("log", "log", 7),
		ng.SideBarButton("quit", "quit", 8),
	})
	a.ButtonBar([]l.Widget{
		ng.PageTopBarButton("help", 0, icons.ActionHelp),
		ng.PageTopBarButton("log", 1, icons.ActionList),
		ng.PageTopBarButton("settings", 2, icons.ActionSettings),
		ng.PageTopBarButton("quit", 3, icons.ActionExitToApp),
	})
	a.StatusBar([]l.Widget{
		ng.StatusBarButton("help", 0, icons.ActionHelp),
		ng.StatusBarButton("log", 1, icons.ActionList),
		ng.StatusBarButton("settings", 2, icons.ActionSettings),
	})
	return
}

func (ng *WalletGUI) Page(title string, widget p9.Widgets) func(gtx l.Context) l.Dimensions {
	a := ng.th
	return func(gtx l.Context) l.Dimensions {
		return a.Fill(ng.BodyBackgroundGet(),
			a.VFlex().
				SpaceEvenly().
				Rigid(
					a.Responsive(*ng.Size, p9.Widgets{
						p9.WidgetSize{
							Widget: a.Inset(0.25, a.H5(title).Color(ng.BodyColorGet()).Fn).Fn,
						},
						p9.WidgetSize{
							Size:   800,
							Widget: p9.EmptySpace(0, 0),
							// a.Inset(0.25, a.Caption(title).Color(ng.BodyColorGet()).Fn).Fn,
						},
					}).Fn,
				).
				Flexed(1,
					a.Inset(0.25,
						a.Responsive(*ng.Size, widget).Fn,
					).Fn,
				).Fn,
		).Fn(gtx)
	}
}

func (ng *WalletGUI) SideBarButton(title, page string, index int) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		gtx.Constraints.Max.X = int(ng.TextSize.Scale(12).V)
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return ng.ButtonLayout(ng.sidebarButtons[index]).Embed(
			func(gtx l.Context) l.Dimensions {
				background := "Transparent"
				color := "DocText"
				if ng.ActivePageGet() == page {
					background = "PanelBg"
					color = "PanelText"
				}
				var inPad, outPad float32 = 0.5, 0.25
				if *ng.Size >= 800 {
					inPad, outPad = 0.75, 0
				}
				return ng.Inset(outPad,
					ng.Fill(background,
						ng.Flex().
							Flexed(1,
								ng.Inset(inPad,
									ng.H6(title).
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
					if ng.MenuOpen {
						ng.MenuOpen = false
					}
					ng.ActivePage(page)
				}).
			Fn(gtx)
	}
}

func (ng *WalletGUI) PageTopBarButton(name string, index int, ico []byte) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		background := ng.TitleBarBackgroundGet()
		color := ng.MenuColorGet()
		if ng.ActivePageGet() == name {
			color = "PanelText"
			background = "PanelBg"
		}
		ic := ng.Icon().
			Scale(p9.Scales["H5"]).
			Color(color).
			Src(ico).
			Fn
		return ng.Flex().Rigid(
			// ng.Inset(0.25,
			ng.ButtonLayout(ng.buttonBarButtons[index]).
				CornerRadius(0).
				Embed(
					ng.Inset(0.375,
						ic,
					).Fn,
				).
				Background(background).
				SetClick(
					func() {
						if ng.MenuOpen {
							ng.MenuOpen = false
						}
						ng.ActivePage(name)
					}).
				Fn,
			// ).Fn,
		).Fn(gtx)
	}
}

func (ng *WalletGUI) StatusBarButton(name string, index int, ico []byte) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		background := ng.StatusBarBackgroundGet()
		color := ng.StatusBarColorGet()
		ic := ng.Icon().
			Scale(p9.Scales["H5"]).
			Color(color).
			Src(ico).
			Fn
		return ng.Flex().
			Rigid(
				ng.ButtonLayout(ng.statusBarButtons[index]).
					CornerRadius(0).
					Embed(
						ic,
					).
					Background(background).
					SetClick(
						func() {
							if ng.MenuOpen {
								ng.MenuOpen = false
							}
							ng.ActivePage(name)
						}).
					Fn,
			).Fn(gtx)
	}
}
