package explorer

import (
	"strconv"
	
	"golang.org/x/exp/shiny/materialdesign/icons"
	
	l "gioui.org/layout"
	"gioui.org/text"
	
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/gui"
	"github.com/p9c/pod/pkg/gui/cfg"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func (ex *Explorer) GetAppWidget() (a *gui.App) {
	a = ex.th.App(&ex.w.Width, nil, nil)
	ex.App = a
	ex.App.ThemeHook(func() {
		Debug("theme hook")
		Debug(ex.bools)
		*ex.cx.Config.DarkTheme = *ex.th.Dark
		a := ex.configs["config"]["DarkTheme"].Slot.(*bool)
		*a = *ex.th.Dark
		if wgb, ok := ex.config.Bools["DarkTheme"]; ok {
			wgb.Value(*ex.th.Dark)
		}
		save.Pod(ex.cx.Config)

	})
	ex.size = a.Size
	ex.config = cfg.New(ex.cx, ex.th)
	ex.configs = ex.config.Config()
	a.Pages(map[string]l.Widget{
		"home": ex.Page("home", gui.Widgets{
			gui.WidgetSize{Widget: ex.Blocks()},
		}),
		"help": ex.Page("help", gui.Widgets{
			gui.WidgetSize{Widget: gui.EmptyMaxHeight()},
		}),
		"quit": ex.Page("quit", gui.Widgets{
			gui.WidgetSize{Widget: func(gtx l.Context) l.Dimensions {
				return ex.th.VFlex().
					SpaceEvenly().
					// AlignMiddle().
					Rigid(
						ex.th.H4("are you sure?").Color(ex.App.BodyColorGet()).Alignment(text.Middle).Fn,
					).
					Rigid(
						ex.th.Flex().
							SpaceEvenly().
							Rigid(
								ex.th.Button(ex.clickables["quit"].SetClick(func() {
									interrupt.Request()
								})).Color(ex.App.TitleBarColorGet()).TextScale(2).Text("yes!!!").Fn,
							).Fn,
					).
					Fn(gtx)
			},
			},
		}),
	})
	a.ButtonBar([]l.Widget{
		ex.PageTopBarButton("help", 0, &icons.ActionHelp),
		ex.PageTopBarButton("quit", 3, &icons.ActionExitToApp),
	})
	a.StatusBar([]l.Widget{
		ex.RunStatusButton(),
		ex.th.Flex().Rigid(
			ex.StatusBarButton("log", 1, &icons.ActionList),
		).Fn,
	})
	return
}

func (ex *Explorer) Page(title string, widget gui.Widgets) func(gtx l.Context) l.Dimensions {
	a := ex.th
	return func(gtx l.Context) l.Dimensions {
		return a.Fill(ex.BodyBackgroundGet(), a.VFlex().
			SpaceEvenly().
			Rigid(
				a.Responsive(*ex.Size, gui.Widgets{
					gui.WidgetSize{
						Widget: a.Inset(0.25, a.H5(title).Color(ex.BodyColorGet()).Fn).Fn,
					},
					gui.WidgetSize{
						Size:   800,
						Widget: gui.EmptySpace(0, 0),
						// a.Inset(0.25, a.Caption(title).Color(ex.BodyColorGet()).Fn).Fn,
					},
				}).Fn,
			).
			Flexed(1,
				a.Inset(0.25,
					a.Responsive(*ex.Size, widget).Fn,
				).Fn,
			).Fn, l.Center).Fn(gtx)
	}
}

func (ex *Explorer) PageTopBarButton(name string, index int, ico *[]byte) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		background := ex.TitleBarBackgroundGet()
		color := ex.MenuColorGet()
		if ex.ActivePageGet() == name {
			color = "PanelText"
			background = "PanelBg"
		}
		ic := ex.th.Icon().
			Scale(gui.Scales["H5"]).
			Color(color).
			Src(ico).
			Fn
		return ex.th.Flex().Rigid(
			// ex.Inset(0.25,
			ex.th.ButtonLayout(ex.buttonBarButtons[index]).
				CornerRadius(0).
				Embed(
					ex.th.Inset(0.375,
						ic,
					).Fn,
				).
				Background(background).
				SetClick(
					func() {
						if ex.MenuOpen {
							ex.MenuOpen = false
						}
						ex.ActivePage(name)
					}).
				Fn,
			// ).Fn,
		).Fn(gtx)
	}
}

func (ex *Explorer) StatusBarButton(name string, index int, ico *[]byte) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		background := ex.StatusBarBackgroundGet()
		color := ex.StatusBarColorGet()
		ic := ex.th.Icon().
			Scale(gui.Scales["H5"]).
			Color(color).
			Src(ico).
			Fn
		return ex.th.Flex().
			Rigid(
				ex.th.ButtonLayout(ex.statusBarButtons[index]).
					CornerRadius(0).
					Embed(
						ic,
					).
					Background(background).
					SetClick(
						func() {
							if ex.MenuOpen {
								ex.MenuOpen = false
							}
							ex.ActivePage(name)
						}).
					Fn,
			).Fn(gtx)
	}
}

func (ex *Explorer) SetRunState(b bool) {
	go func() {
		Debug("run state is now", b)
		if b {
			ex.RunCommandChan <- "run"
			// ex.running = b
		} else {
			ex.RunCommandChan <- "stop"
			// ex.running = b
		}
	}()
}

func (ex *Explorer) RunStatusButton() func(gtx l.Context) l.Dimensions {
	t, f := &icons.AVStop, &icons.AVPlayArrow
	return func(gtx l.Context) l.Dimensions {
		background := ex.App.StatusBarBackgroundGet()
		color := ex.App.StatusBarColorGet()
		var ico *[]byte
		if ex.running {
			ico = t
		} else {
			ico = f
		}
		ic := ex.th.Icon().
			Scale(gui.Scales["H4"]).
			Color(color).
			Src(ico).
			Fn
		return ex.th.Flex().
			Rigid(
				ex.th.ButtonLayout(ex.statusBarButtons[0]).
					CornerRadius(0).
					Embed(
						ex.th.Inset(0.066, ic).Fn,
					).
					Background(background).
					SetClick(
						func() {
							ex.SetRunState(!ex.running)
						}).
					Fn,
			).
			Rigid(
				ex.th.Inset(0.33,
					gui.If(ex.running,
						ex.th.Indefinite().Scale(gui.Scales["H5"]).Fn,
						ex.th.Icon().
							Scale(gui.Scales["H5"]).
							Color("Primary").
							Src(&icons.ActionCheckCircle).
							Fn,
					),
				).Fn,
			).
			Rigid(
				ex.th.Inset(0.33,
					ex.th.H5(strconv.FormatInt(int64(ex.State.bestBlockHeight), 10)).Color(color).Fn,
				).Fn,
			).
			Fn(gtx)
	}
}
