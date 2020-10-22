package gui

import (
	l "gioui.org/layout"
	"gioui.org/text"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func (ng *NodeGUI) GetAppWidget() (a *p9.App) {
	a = ng.th.App(*ng.size)
	ng.App = a
	ng.size = ng.Size
	ng.Theme.Colors.SetTheme(ng.Dark)
	ng.Pages(p9.WidgetMap{
		"main": ng.Page("overview", p9.Widgets{
			p9.WidgetSize{
				Widget:
				func(gtx l.Context) l.Dimensions {
					return ng.VFlex().Rigid(
						ng.CardList(ng.lists["overview"], ng.CardBackgroundGet(),
							ng.CardContent("run settings", "Primary",
								ng.VFlex().Rigid(
									ng.Flex().
										Rigid(
											ng.Body1("run").Fn,
										).
										Rigid(
											ng.Switch(ng.bools["runstate"]).Fn,
										).
										Fn,
								).Rigid(
									ng.Flex().
										Rigid(
											ng.Body1("mode").Fn,
										).
										Rigid(
											ng.Responsive(gtx.Constraints.Max.X, p9.Widgets{
												{
													Widget: ng.VFlex().
														Rigid(
															ng.RadioButton(ng.enums["runmode"], "node", "node").Fn,
														).
														Rigid(
															ng.RadioButton(ng.enums["runmode"], "wallet", "wallet").Fn,
														).
														Rigid(
															ng.RadioButton(ng.enums["runmode"], "shell", "shell").Fn,
														).Fn,
												},
												{Size: 512,
													Widget: ng.Flex().
														Rigid(
															ng.RadioButton(ng.enums["runmode"], "node", "node").Fn,
														).
														Rigid(
															ng.RadioButton(ng.enums["runmode"], "wallet", "wallet").Fn,
														).
														Rigid(
															ng.RadioButton(ng.enums["runmode"], "shell", "shell").Fn,
														).Fn,
												},
											}).Fn,
										).Fn,
								).Fn,
							),
							ng.CardContent("mining info", ng.CardColorGet(),
								ng.Flex().
									Rigid(
										ng.Body1("I will show the current data about difficulty adjustment").
											Color(ng.CardColorGet()).
											Fn,
									).
									Fn,
							),
							ng.CardContent("network hashrate", ng.CardColorGet(),
								ng.Flex().
									Rigid(
										ng.Body1("i will show a graph of the hashrate on the lan").
											Color(ng.CardColorGet()).
											Fn,
									).
									Fn,
							),
							ng.CardContent("log", ng.CardColorGet(),
								ng.Flex().
									Flexed(1,
										ng.Body1("i will become a log viewer").
											Color(ng.CardColorGet()).
											Fn,
									).
									Fn,
							),
						),
					).Fn(gtx)
				},
			},
		}),
		"settings": ng.Page("settings", p9.Widgets{
			p9.WidgetSize{Widget: ng.Config()},
		}),
		"help": ng.Page("help", p9.Widgets{
			p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
		}),
		"log": ng.Page("log", p9.Widgets{
			p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
		}),
		"quit": ng.Page("quit", p9.Widgets{
			p9.WidgetSize{Widget:
			func(gtx l.Context) l.Dimensions {
				return ng.VFlex().
					SpaceEvenly().
					// AlignMiddle().
					Rigid(
						ng.H4("are you sure?").Color(ng.BodyColorGet()).Alignment(text.Middle).Fn,
					).
					Rigid(
						ng.Flex().
							SpaceEvenly().
							Rigid(
								ng.Button(ng.quitClickable.SetClick(func() {
									interrupt.Request()
								})).Color(ng.TitleBarColorGet()).TextScale(2).Text("yes!!!").Fn,
							).Fn,
					).
					Fn(gtx)
			},
			},
		}),
	})
	ng.SideBar([]l.Widget{
		ng.SideBarButton("overview", "main", 0),
		ng.SideBarButton("settings", "settings", 5),
		// ng.SideBarButton("help", "help", 6),
		ng.SideBarButton("log", "log", 7),
		// ng.SideBarButton("quit", "quit", 8),
	})
	ng.ButtonBar([]l.Widget{
		ng.PageTopBarButton("help", 0, icons.ActionHelp),
		// ng.PageTopBarButton("log", 1, icons.ActionList),
		// ng.PageTopBarButton("settings", 2, icons.ActionSettings),
		ng.PageTopBarButton("quit", 3, icons.ActionExitToApp),
	})
	ng.StatusBar([]l.Widget{
		ng.RunStatusButton(),
		// ng.StatusBarButton("help", 0, icons.AVPlayArrow),
		ng.Flex().Rigid(
			ng.StatusBarButton("log", 1, icons.ActionList),
		).Rigid(
			ng.StatusBarButton("settings", 2, icons.ActionSettings),
		).Fn,
	})
	ng.Title("node")
	return
}

// Page renders a page. Note that the widgets you give it should be written wrapped in functions if
// the fluent declarations are used for values inside the ng parent type, as they are computed then at declaration
// and not at the time of execution.
func (ng *NodeGUI) Page(title string, widget p9.Widgets) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		return ng.Fill(ng.BodyBackgroundGet(),
			ng.VFlex().
				SpaceEvenly().
				Rigid(
					ng.Responsive(*ng.Size, p9.Widgets{
						{
							Widget: func(gtx l.Context) l.Dimensions {
								if ng.MenuOpen {
									return p9.EmptySpace(0, 0)(gtx)
								} else {
									return ng.Inset(0.25, ng.H6(title).Color(ng.BodyColorGet()).Fn).Fn(gtx)
								}
							},
						},
						{
							Size:   800,
							Widget: p9.EmptySpace(0, 0),
						},
					}).Fn,
				).
				Flexed(1,
					ng.Inset(0.25,
						ng.Responsive(*ng.Size, widget).Fn,
					).Fn,
				).Fn,
		).Fn(gtx)
	}
}

func (ng *NodeGUI) SideBarButton(title, page string, index int) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		return ng.ButtonLayout(ng.sidebarButtons[index]).Embed(
			func(gtx l.Context) l.Dimensions {
				// gtx.Constraints.Max.X = int(ng.TextSize.Scale(12).V)
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

func (ng *NodeGUI) PageTopBarButton(name string, index int, ico []byte) func(gtx l.Context) l.Dimensions {
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
					ng.Inset(0.4,
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

func (ng *NodeGUI) StatusBarButton(name string, index int, ico []byte) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		background := ng.StatusBarBackgroundGet()
		color := ng.StatusBarColorGet()
		if ng.ActivePageGet() == name {
			background = ng.BodyBackgroundGet()
		}
		ic := ng.Icon().
			Scale(p9.Scales["H4"]).
			Color(color).
			Src(ico).
			Fn
		return ng.Flex().
			Rigid(
				ng.ButtonLayout(ng.statusBarButtons[index]).
					CornerRadius(0).
					Embed(
						ng.Inset(0.066, ic).Fn,
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

func (ng *NodeGUI) RunStatusButton() func(gtx l.Context) l.Dimensions {
	t, f := icons.AVStop, icons.AVPlayArrow
	return func(gtx l.Context) l.Dimensions {
		state := ng.bools["runstate"].GetValue()
		background := ng.StatusBarBackgroundGet()
		color := ng.StatusBarColorGet()
		var st bool
		if state {
			st = true
			background = "Primary"
		}
		var ico []byte
		if st {
			ico = t
		} else {
			ico = f
		}
		ic := ng.Icon().
			Scale(p9.Scales["H4"]).
			Color(color).
			Src(ico).
			Fn
		return ng.Flex().
			Rigid(
				ng.ButtonLayout(ng.statusBarButtons[0]).
					CornerRadius(0).
					Embed(
						ng.Inset(0.066, ic).Fn,
					).
					Background(background).
					SetClick(
						func() {
							ng.bools["runstate"].Value(!ng.bools["runstate"].GetValue())
						}).
					Fn,
			).
			Rigid(
				ng.Inset(0.33,
					p9.If(ng.bools["runstate"].GetValue(),
						ng.Indefinite().Scale(p9.Scales["H5"]).Fn,
						ng.Icon().
							Scale(p9.Scales["H5"]).
							Color("Primary").
							Src(icons.ActionCheckCircle).
							Fn,
					),
				).Fn,
			).
			Rigid(
				ng.Inset(0.33,
					ng.H5("256789").Color(color).Fn,
				).Fn,
			).
			Fn(gtx)
	}
}
