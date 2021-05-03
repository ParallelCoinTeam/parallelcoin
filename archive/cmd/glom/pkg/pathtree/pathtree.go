package pathtree

import (
	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/text"
	uberatomic "go.uber.org/atomic"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/opts/binary"
	"github.com/p9c/opts/meta"

	"github.com/p9c/gel"
)

type Widget struct {
	*gel.Window
	*gel.App
	activePage       *uberatomic.String
	sidebarButtons   []*gel.Clickable
	statusBarButtons []*gel.Clickable
	buttonBarButtons []*gel.Clickable
	Size             *uberatomic.Int32
}

func New(w *gel.Window) (wg *Widget) {
	activePage := uberatomic.NewString("home")
	w.Dark = binary.New(meta.Data{}, false, func(b bool) error { return nil })
	w.Colors.SetDarkTheme(false)
	// I.S(w.Colors)
	app := w.App(w.Width, activePage, 48)
	wg = &Widget{
		Window:     w,
		App:        app,
		activePage: uberatomic.NewString("home"),
		Size:       w.Width,
	}
	wg.GetButtons()
	app.Pages(
		map[string]l.Widget{
			"home": wg.Page(
				"home", gel.Widgets{
					// p9.WidgetSize{Widget: p9.EmptyMaxHeight()},
					gel.WidgetSize{
						Widget: wg.Flex().Flexed(1, wg.H3("glom").Fn).Fn,
					},
				},
			),
		},
	)
	app.SideBar(
		[]l.Widget{
			// wg.SideBarButton(" ", " ", 11),
			wg.SideBarButton("home", "home", 0),
		},
	)
	app.ButtonBar(
		[]l.Widget{
			wg.PageTopBarButton(
				"help", 0, &icons.ActionHelp, func(name string) {
				}, app, "",
			),
			wg.PageTopBarButton(
				"home", 1, &icons.ActionLockOpen, func(name string) {
					wg.App.ActivePage(name)
				}, app, "green",
			),
			// wg.Flex().Rigid(wg.Inset(0.5, gel.EmptySpace(0, 0)).Fn).Fn,
			// wg.PageTopBarButton(
			// 	"quit", 3, &icons.ActionExitToApp, func(name string) {
			// 		wg.MainApp.ActivePage(name)
			// 	}, a, "",
			// ),
		},
	)
	app.StatusBar(
		[]l.Widget{
			wg.StatusBarButton(
				"log", 0, &icons.ActionList, func(name string) {
					D.Ln("click on button", name)
				}, app,
			),
		},
		[]l.Widget{
			wg.StatusBarButton(
				"settings", 1, &icons.ActionSettings, func(name string) {
					D.Ln("click on button", name)
				}, app,
			),
		},
	)
	return
}

func (w *Widget) Fn(gtx l.Context) l.Dimensions {
	return w.App.Fn()(gtx)
}

func (w *Widget) Page(title string, widget gel.Widgets) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		return w.VFlex().
			// SpaceEvenly().
			Rigid(
				w.Responsive(
					w.Size.Load(), gel.Widgets{
						// p9.WidgetSize{
						// 	Widget: a.ButtonInset(0.25, a.H5(title).Color(wg.App.BodyColorGet()).Fn).Fn,
						// },
						gel.WidgetSize{
							// Size:   800,
							Widget: gel.EmptySpace(0, 0),
							// a.ButtonInset(0.25, a.Caption(title).Color(wg.BodyColorGet()).Fn).Fn,
						},
					},
				).Fn,
			).
			Flexed(
				1,
				w.Inset(
					0.25,
					w.Responsive(w.Size.Load(), widget).Fn,
				).Fn,
			).Fn(gtx)
	}
}

func (wg *Widget) GetButtons() {
	wg.sidebarButtons = make([]*gel.Clickable, 2)
	// wg.walletLocked.Store(true)
	for i := range wg.sidebarButtons {
		wg.sidebarButtons[i] = wg.Clickable()
	}
	wg.buttonBarButtons = make([]*gel.Clickable, 2)
	for i := range wg.buttonBarButtons {
		wg.buttonBarButtons[i] = wg.Clickable()
	}
	wg.statusBarButtons = make([]*gel.Clickable, 2)
	for i := range wg.statusBarButtons {
		wg.statusBarButtons[i] = wg.Clickable()
	}
}

func (w *Widget) SideBarButton(title, page string, index int) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		var scale float32
		scale = gel.Scales["H6"]
		var color string
		background := "Transparent"
		color = "DocText"
		var ins float32 = 0.5
		// var hl = false
		if w.App.ActivePageGet() == page || w.App.PreRendering {
			background = "PanelBg"
			scale = gel.Scales["H6"]
			color = "DocText"
			// ins = 0.5
			// hl = true
		}
		if title == " " {
			scale = gel.Scales["H6"] / 2
		}
		max := int(w.App.SideBarSize.V)
		if max > 0 {
			gtx.Constraints.Max.X = max
			gtx.Constraints.Min.X = max
		}
		// D.Ln("sideMAXXXXXX!!", max)
		return w.Direction().E().Embed(
			w.ButtonLayout(w.sidebarButtons[index]).
				CornerRadius(scale).Corners(0).
				Background(background).
				Embed(
					w.Inset(
						ins,
						func(gtx l.Context) l.Dimensions {
							return w.H5(title).
								Color(color).
								Alignment(text.End).
								Fn(gtx)
						},
					).Fn,
				).
				SetClick(
					func() {
						if w.App.MenuOpen {
							w.App.MenuOpen = false
						}
						w.App.ActivePage(page)
					},
				).
				Fn,
		).
			Fn(gtx)
	}
}

func (w *Widget) PageTopBarButton(
	name string, index int, ico *[]byte, onClick func(string), app *gel.App,
	highlightColor string,
) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		background := "Transparent"
		// background := node.TitleBarBackgroundGet()
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
		ic := w.Icon().
			Scale(gel.Scales["H5"]).
			Color(color).
			Src(ico).
			Fn
		return w.Flex().Rigid(
			// wg.ButtonInset(0.25,
			w.ButtonLayout(w.buttonBarButtons[index]).
				CornerRadius(0).
				Embed(
					w.Inset(
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

func (w *Widget) StatusBarButton(
	name string,
	index int,
	ico *[]byte,
	onClick func(string),
	app *gel.App,
) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		background := app.StatusBarBackgroundGet()
		color := app.StatusBarColorGet()
		if app.ActivePageGet() == name {
			// background, color = color, background
			background = "PanelBg"
			// color = "Danger"
		}
		ic := w.Icon().
			Scale(gel.Scales["H5"]).
			Color(color).
			Src(ico).
			Fn
		return w.Flex().
			Rigid(
				w.ButtonLayout(w.statusBarButtons[index]).
					CornerRadius(0).
					Embed(
						w.Inset(0.25, ic).Fn,
					).
					Background(background).
					SetClick(func() { onClick(name) }).
					Fn,
			).Fn(gtx)
	}
}
