package gui

import (
	l "gioui.org/layout"
	"gioui.org/text"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/pod/pkg/gui/cfg"
	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func (ng *NodeGUI) GetAppWidget() (a *p9.App) {
	a = ng.th.App(ng.w.Width)
	ng.app = a
	// ng.size = ng.size
	ng.th.Colors.SetTheme(*ng.app.Dark)
	ng.config = cfg.New(ng.cx, ng.th)
	ng.configs = ng.config.Config()
	ng.app.ThemeHook(func() {})
	ng.app.Pages(p9.WidgetMap{
		"main": ng.Page("overview", p9.Widgets{
			p9.WidgetSize{
				Widget: func(gtx l.Context) l.Dimensions {
					return ng.th.VFlex().Rigid(
						ng.th.CardList(ng.lists["overview"], ng.app.CardBackgroundGet(),
							ng.app.CardContent("run settings", "Primary",
								ng.th.VFlex().Rigid(
									ng.th.Flex().
										Rigid(
											ng.th.Body1("run").Fn,
										).
										Rigid(
											ng.th.Switch(ng.bools["runstate"]).Fn,
										).
										Fn,
								).Rigid(
									ng.th.Flex().
										Rigid(
											ng.th.Body1("mode").Fn,
										).
										Rigid(
											ng.th.Responsive(gtx.Constraints.Max.X, p9.Widgets{
												{
													Widget: ng.th.VFlex().
														Rigid(
															ng.th.Inset(0.125,
																ng.th.RadioButton(ng.checkables["runmodenode"].
																	Color("DocText").
																	IconColor("DocText"), ng.enums["runmode"],
																	"node", "node").Fn,
															).Fn,
														).
														Rigid(
															ng.th.Inset(0.125,
																ng.th.RadioButton(ng.checkables["runmodewallet"].
																	Color("DocText").
																	IconColor("DocText"), ng.enums["runmode"],
																	"wallet", "wallet").Fn,
															).Fn,
														).
														Rigid(
															ng.th.Inset(0.125,
																ng.th.RadioButton(ng.checkables["runmodeshell"].
																	Color("DocText").
																	IconColor("DocText"), ng.enums["runmode"],
																	"shell", "shell").Fn,
															).Fn,
														).Fn,
												},
												{Size: 512,
													Widget: ng.th.Flex().
														Rigid(
															ng.th.Inset(0.125,
																ng.th.RadioButton(ng.checkables["runmodenode"].
																	Color("DocText").
																	IconColor("DocText"), ng.enums["runmode"],
																	"node", "node").Fn,
															).Fn,
														).
														Rigid(
															ng.th.Inset(0.125,
																ng.th.RadioButton(ng.checkables["runmodewallet"].
																	Color("DocText").
																	IconColor("DocText"), ng.enums["runmode"],
																	"wallet", "wallet").Fn,
															).Fn,
														).
														Rigid(
															ng.th.Inset(0.125,
																ng.th.RadioButton(ng.checkables["runmodeshell"].
																	Color("DocText").
																	IconColor("DocText"), ng.enums["runmode"],
																	"shell", "shell").Fn,
															).Fn,
														).Fn,
												},
											}).Fn,
										).Fn,
								).Fn,
							),
							ng.th.CardContent("mining info", ng.app.CardColorGet(),
								ng.th.Flex().
									Rigid(
										ng.th.Body1("I will show the current data about difficulty adjustment").
											Color(ng.app.CardColorGet()).
											Fn,
									).
									Fn,
							),
							ng.th.CardContent("network hashrate", ng.app.CardColorGet(),
								ng.th.Flex().
									Rigid(
										ng.th.Body1("i will show a graph of the hashrate on the lan").
											Color(ng.app.CardColorGet()).
											Fn,
									).
									Fn,
							),
							ng.th.CardContent("log", ng.app.CardColorGet(),
								ng.th.Flex().
									Flexed(1,
										ng.th.Body1("i will become a log viewer").
											Color(ng.app.CardColorGet()).
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
			p9.WidgetSize{Widget: func(gtx l.Context) l.Dimensions {
				return ng.configs.Widget(ng.config)(gtx)
			}},
		}),
		"help": ng.Page("help", p9.Widgets{
			p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
		}),
		"log": ng.Page("log", p9.Widgets{
			p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
		}),
		"quit": ng.Page("quit", p9.Widgets{
			p9.WidgetSize{Widget: func(gtx l.Context) l.Dimensions {
				return ng.th.VFlex().
					SpaceEvenly().
					// AlignMiddle().
					Rigid(
						ng.th.H4("are you sure?").Color(ng.app.BodyColorGet()).Alignment(text.Middle).Fn,
					).
					Rigid(
						ng.th.Flex().
							SpaceEvenly().
							Rigid(
								ng.th.Button(ng.clickables["quit"].SetClick(func() {
									interrupt.Request()
								})).Color(ng.app.TitleBarColorGet()).TextScale(2).Text("yes!!!").Fn,
							).Fn,
					).
					Fn(gtx)
			},
			},
		}),
	})
	ng.app.SideBar([]l.Widget{
		ng.SideBarButton("overview", "main", 0),
		ng.SideBarButton("settings", "settings", 5),
		// ng.SideBarButton("help", "help", 6),
		ng.SideBarButton("log", "log", 7),
		// ng.SideBarButton("quit", "quit", 8),
	})
	ng.app.ButtonBar([]l.Widget{
		ng.PageTopBarButton("help", 0, &icons.ActionHelp),
		// ng.PageTopBarButton("log", 1, icons.ActionList),
		// ng.PageTopBarButton("settings", 2, icons.ActionSettings),
		ng.PageTopBarButton("quit", 3, &icons.ActionExitToApp),
	})
	ng.app.StatusBar([]l.Widget{
		ng.RunStatusButton(),
		// ng.StatusBarButton("help", 0, icons.AVPlayArrow),
		ng.th.Flex().Rigid(
			ng.StatusBarButton("log", 1, &icons.ActionList),
		).Rigid(
			ng.StatusBarButton("settings", 2, &icons.ActionSettings),
		).Fn,
	})
	ng.app.Title("node")
	return
}

// Page renders a page. Note that the widgets you give it should be written wrapped in functions if
// the fluent declarations are used for values inside the ng parent type, as they are computed then at declaration
// and not at the time of execution.
func (ng *NodeGUI) Page(title string, widget p9.Widgets) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		return ng.th.Fill(ng.app.BodyBackgroundGet(),
			ng.th.VFlex().
				SpaceEvenly().
				Rigid(
					ng.th.Responsive(*ng.app.Size, p9.Widgets{
						{
							Widget: func(gtx l.Context) l.Dimensions {
								if ng.app.MenuOpen {
									return p9.EmptySpace(0, 0)(gtx)
								} else {
									return ng.th.Inset(0.25, ng.th.H6(title).Color(ng.app.BodyColorGet()).Fn).Fn(gtx)
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
					ng.th.Inset(0.25,
						ng.th.Responsive(*ng.app.Size, widget).Fn,
					).Fn,
				).Fn,
		).Fn(gtx)
	}
}

func (ng *NodeGUI) SideBarButton(title, page string, index int) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		gtx.Constraints.Max.X = int(ng.th.TextSize.Scale(12).V)
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return ng.th.ButtonLayout(ng.sidebarButtons[index]).Embed(
			func(gtx l.Context) l.Dimensions {
				background := "Transparent"
				color := "DocText"
				if ng.app.ActivePageGet() == page {
					background = "PanelBg"
					color = "PanelText"
				}
				var inPad, outPad float32 = 0.5, 0.25
				if *ng.app.Size >= 800 {
					inPad, outPad = 0.75, 0
				}
				return ng.th.Inset(outPad,
					ng.th.Fill(background,
						ng.th.Flex().
							Flexed(1,
								ng.th.Inset(inPad,
									ng.th.H6(title).
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
					if ng.app.MenuOpen {
						ng.app.MenuOpen = false
					}
					ng.app.ActivePage(page)
				}).
			Fn(gtx)
	}
}

func (ng *NodeGUI) PageTopBarButton(name string, index int, ico *[]byte) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		background := ng.app.TitleBarBackgroundGet()
		color := ng.app.MenuColorGet()
		if ng.app.ActivePageGet() == name {
			color = "PanelText"
			background = "PanelBg"
		}
		ic := ng.th.Icon().
			Scale(p9.Scales["H5"]).
			Color(color).
			Src(ico).
			Fn
		return ng.th.Flex().Rigid(
			// ng.Inset(0.25,
			ng.th.ButtonLayout(ng.buttonBarButtons[index]).
				CornerRadius(0).
				Embed(
					ng.th.Inset(0.4,
						ic,
					).Fn,
				).
				Background(background).
				SetClick(
					func() {
						if ng.app.MenuOpen {
							ng.app.MenuOpen = false
						}
						ng.app.ActivePage(name)
					}).
				Fn,
			// ).Fn,
		).Fn(gtx)
	}
}

func (ng *NodeGUI) StatusBarButton(name string, index int, ico *[]byte) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		background := ng.app.StatusBarBackgroundGet()
		color := ng.app.StatusBarColorGet()
		if ng.app.ActivePageGet() == name {
			background = ng.app.BodyBackgroundGet()
		}
		ic := ng.th.Icon().
			Scale(p9.Scales["H4"]).
			Color(color).
			Src(ico).
			Fn
		return ng.th.Flex().
			Rigid(
				ng.th.ButtonLayout(ng.statusBarButtons[index]).
					CornerRadius(0).
					Embed(
						ng.th.Inset(0.066, ic).Fn,
					).
					Background(background).
					SetClick(
						func() {
							if ng.app.MenuOpen {
								ng.app.MenuOpen = false
							}
							ng.app.ActivePage(name)
						}).
					Fn,
			).Fn(gtx)
	}
}

func (ng *NodeGUI) RunStatusButton() func(gtx l.Context) l.Dimensions {
	t, f := &icons.AVStop, &icons.AVPlayArrow
	return func(gtx l.Context) l.Dimensions {
		state := ng.bools["runstate"].GetValue()
		background := ng.app.StatusBarBackgroundGet()
		color := ng.app.StatusBarColorGet()
		var st bool
		if state {
			st = true
			background = "Primary"
		}
		var ico *[]byte
		if st {
			ico = t
		} else {
			ico = f
		}
		ic := ng.th.Icon().
			Scale(p9.Scales["H4"]).
			Color(color).
			Src(ico).
			Fn
		return ng.th.Flex().
			Rigid(
				ng.th.ButtonLayout(ng.statusBarButtons[0]).
					CornerRadius(0).
					Embed(
						ng.th.Inset(0.066, ic).Fn,
					).
					Background(background).
					SetClick(
						func() {
							ng.bools["runstate"].Value(!ng.bools["runstate"].GetValue())
						}).
					Fn,
			).
			Rigid(
				ng.th.Inset(0.33,
					p9.If(ng.bools["runstate"].GetValue(),
						ng.th.Indefinite().Scale(p9.Scales["H5"]).Fn,
						ng.th.Icon().
							Scale(p9.Scales["H5"]).
							Color("Primary").
							Src(&icons.ActionCheckCircle).
							Fn,
					),
				).Fn,
			).
			Rigid(
				ng.th.Inset(0.33,
					ng.th.H5("256789").Color(color).Fn,
				).Fn,
			).
			Fn(gtx)
	}
}
