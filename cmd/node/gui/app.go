package gui

import (
	l "gioui.org/layout"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/pod/pkg/gui/p9"
)

func (ng *NodeGUI) Page(title string, widget p9.Widgets) func(gtx l.Context) l.Dimensions {
	a := ng.th
	return func(gtx l.Context) l.Dimensions {
		// width := gtx.Constraints.Max.X
		// height := gtx.Constraints.Max.Y
		return a.VFlex().
			SpaceEvenly().
			Rigid(
				a.Responsive(*ng.Size, p9.Widgets{
					p9.WidgetSize{
						Widget: a.Inset(0.25, a.H5(title).Color(ng.BodyColorGet()).Fn).Fn,
					},
					p9.WidgetSize{
						Size:   800,
						Widget: a.Inset(0.25, a.Caption(title).Color(ng.BodyColorGet()).Fn).Fn,
					},
				}).Fn,
			).
			Flexed(1,
				a.Inset(0.25,
					a.Responsive(*ng.Size, widget).Fn,
				).Fn,
			).
			Fn(gtx)
	}
}

func (ng *NodeGUI) GetAppWidget() (a *p9.App) {
	a = ng.th.App(*ng.size)
	ng.App = a
	ng.size = a.Size
	a.Pages(map[string]l.Widget{
		"main": ng.Page("first", p9.Widgets{
			p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
		}),
		"second": ng.Page("second", p9.Widgets{
			p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
		}),
		"third": ng.Page("third", p9.Widgets{
			p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
		}),
		"fourth": ng.Page("fourth", p9.Widgets{
			p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
		}),
		"fifth": ng.Page("fifth", p9.Widgets{
			p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
		}),
		"settings": ng.Page("settings", p9.Widgets{
			p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
		}),
		"help": ng.Page("help", p9.Widgets{
			p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
		}),
	})
	a.SideBar([]l.Widget{
		a.
			ButtonLayout(
				ng.
					sidebarButtons[0]).
			Embed(
				func(gtx l.Context) l.Dimensions {
					background := "Transparent"
					color := "DocText"
					if ng.ActivePageGet() == "main" {
						background = "DocText"
						color = "DocBg"
					}
					return ng.Fill(background,
						ng.Flex().Flexed(1,
							ng.th.Inset(0.5,
								ng.th.H6("first").
									Color(color).
									Fn,
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
					a.ActivePage("main")
				}).
			Fn,
		ng.th.ButtonLayout(ng.sidebarButtons[1]).
			Embed(
				func(gtx l.Context) l.Dimensions {
					background := "Transparent"
					color := "DocText"
					if ng.ActivePageGet() == "second" {
						background = "DocText"
						color = "DocBg"
					}
					return ng.th.Fill(background,
						ng.th.Flex().Flexed(1,
							ng.th.Inset(0.5,
								ng.th.H6("second").
									Color(color).
									Fn,
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
					ng.ActivePage("second")
				}).
			Fn,
		ng.th.ButtonLayout(ng.sidebarButtons[2]).
			Embed(
				func(gtx l.Context) l.Dimensions {
					background := "Transparent"
					color := "DocText"
					if ng.ActivePageGet() == "third" {
						background = "DocText"
						color = "DocBg"
					}
					return ng.th.Fill(background,
						ng.th.Flex().Flexed(1,
							ng.th.Inset(0.5,
								ng.th.H6("third").
									Color(color).
									Fn,
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
					a.ActivePage("third")
				}).
			Fn,
		ng.th.ButtonLayout(ng.sidebarButtons[3]).
			Embed(
				func(gtx l.Context) l.Dimensions {
					background := "Transparent"
					color := "DocText"
					if a.ActivePageGet() == "fourth" {
						background = "DocText"
						color = "DocBg"
					}
					return ng.th.Fill(background,
						ng.th.Flex().Flexed(1,
							ng.th.Inset(0.5,
								ng.th.H6("fourth").
									Color(color).
									Fn,
							).Fn,
						).Fn,
					).Fn(gtx)
				},
			).
			Background("Transparent").
			SetClick(
				func() {
					if a.MenuOpen {
						a.MenuOpen = false
					}
					a.ActivePage("fourth")
				}).
			Fn,
		ng.th.ButtonLayout(ng.sidebarButtons[4]).
			Embed(
				func(gtx l.Context) l.Dimensions {
					background := "Transparent"
					color := "DocText"
					if a.ActivePageGet() == "fifth" {
						background = "DocText"
						color = "DocBg"
					}
					return ng.th.Fill(background,
						ng.th.Flex().Flexed(1,
							ng.th.Inset(0.5,
								ng.th.H6("fifth").
									Color(color).
									Fn,
							).Fn,
						).Fn,
					).Fn(gtx)
				},
			).
			Background("Transparent").
			SetClick(
				func() {
					if a.MenuOpen {
						a.MenuOpen = false
					}
					a.ActivePage("fifth")
				}).
			Fn,
		ng.th.ButtonLayout(ng.sidebarButtons[5]).
			Embed(
				func(gtx l.Context) l.Dimensions {
					background := "Transparent"
					color := "DocText"
					return ng.th.Fill(background,
						ng.th.Flex().Flexed(1,
							ng.th.Inset(0.5,
								ng.th.H6("invalid").
									Color(color).
									Fn,
							).Fn,
						).Fn,
					).Fn(gtx)
				},
			).
			Background("Transparent").
			SetClick(
				func() {
					if a.MenuOpen {
						a.MenuOpen = false
					}
					a.ActivePage("invalid")
				}).
			Fn,
	})
	a.ButtonBar([]l.Widget{
		func(gtx l.Context) l.Dimensions {
			background := a.TitleBarBackgroundGet()
			color := a.MenuColorGet()
			if a.ActivePageGet() == "help" {
				// background, color = color,background
				color = "DocText"
			}
			return a.Flex().Rigid(
				a.Inset(0.25,
					a.ButtonLayout(ng.buttonBarButtons[1]).
						CornerRadius(0).
						Embed(
							a.Inset(0.25,
								a.Icon().
									Scale(p9.Scales["H5"]).
									Color(color).
									Src(icons.ActionHelp).
									Fn,
							).Fn,
						).
						Background(background).
						SetClick(
							func() {
								a.ActivePage("help")
							}).
						Fn,
				).Fn,
			).Fn(gtx)
		},
		ng.PageTopBarButton("settings", icons.ActionSettings),
		// func(gtx l.Context) l.Dimensions {
		// 	background := a.TitleBarBackgroundGet()
		// 	color := a.MenuColorGet()
		// 	if a.ActivePageGet() == "settings" {
		// 		// background, color = color,background
		// 		color = "DocText"
		// 	}
		// 	return a.Flex().Rigid(
		// 		a.Inset(0.25,
		// 			a.ButtonLayout(ng.buttonBarButtons[0]).
		// 				CornerRadius(0).
		// 				Embed(
		// 					a.Inset(0.25,
		// 						a.Icon().
		// 							Scale(p9.Scales["H5"]).
		// 							Color(color).
		// 							Src(icons.ActionSettings).
		// 							Fn,
		// 					).Fn,
		// 				).
		// 				Background(background).
		// 				SetClick(
		// 					func() {
		// 						a.ActivePage("settings")
		// 					}).
		// 				Fn,
		// 		).Fn,
		// 	).Fn(gtx)
		// },
	})
	return
}

func (ng *NodeGUI) PageTopBarButton(name string, ico []byte) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		background := ng.TitleBarBackgroundGet()
		color := ng.MenuColorGet()
		if ng.ActivePageGet() == name {
			// background, color = color,background
			color = "DocText"
		}
		ic := ng.Icon().
			Scale(p9.Scales["H5"]).
			Color(color).
			Src(ico).
			Fn
		return ng.Flex().Rigid(
			ng.Inset(0.25,
				ng.ButtonLayout(ng.buttonBarButtons[0]).
					CornerRadius(0).
					Embed(
						ng.Inset(0.25, ic).Fn,
					).
					Background(background).
					SetClick(
						func() {
							ng.ActivePage(name)
						}).
					Fn,
			).Fn,
		).Fn(gtx)
	}
}
