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

func (ng *NodeGUI) SideBarButton(title, page string, index int) func(gtx l.Context) l.Dimensions {
	return ng.ButtonLayout(ng.sidebarButtons[index]).Embed(
		func(gtx l.Context) l.Dimensions {
			background := "Transparent"
			color := "DocText"
			if ng.ActivePageGet() == page {
				background = "DocText"
				color = "DocBg"
			}
			return ng.Fill(background,
				ng.Flex().Flexed(1,
					ng.th.Inset(0.5,
						ng.th.H6(title).
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
				ng.ActivePage(page)
			}).
		Fn
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
		"log": ng.Page("log", p9.Widgets{
			p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
		}),
	})
	a.SideBar([]l.Widget{
		ng.SideBarButton("first", "main", 0),
		ng.SideBarButton("second", "second", 1),
		ng.SideBarButton("third", "third", 2),
		ng.SideBarButton("fourth", "fourth", 3),
		ng.SideBarButton("fifth", "fifth", 4),
		ng.SideBarButton("settings", "settings", 5),
		ng.SideBarButton("help", "help", 6),
		ng.SideBarButton("log", "log", 7),
		ng.SideBarButton("invalid", "invalid", 8),
	})
	a.ButtonBar([]l.Widget{
		ng.PageTopBarButton("help", 0, icons.ActionHelp),
		ng.PageTopBarButton("log", 1, icons.ActionList),
		ng.PageTopBarButton("settings", 2, icons.ActionSettings),
	})
	return
}

func (ng *NodeGUI) PageTopBarButton(name string, index int, ico []byte) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		background := ng.TitleBarBackgroundGet()
		color := ng.MenuColorGet()
		if ng.ActivePageGet() == name {
			// background, color = color,background
			color = "Dark"
		}
		ic := ng.Icon().
			Scale(p9.Scales["H5"]).
			Color(color).
			Src(ico).
			Fn
		return ng.Flex().Rigid(
			ng.Inset(0.25,
				ng.ButtonLayout(ng.buttonBarButtons[index]).
					CornerRadius(0).
					Embed(
						ng.Inset(0.25, ic).Fn,
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
			).Fn,
		).Fn(gtx)
	}
}
