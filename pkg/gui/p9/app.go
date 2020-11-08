package p9

import (
	"fmt"

	l "gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/pod/pkg/gui/wallet/ico"
)

// App defines an application with a header, sidebar/menu, right side button bar, changeable body page widget and
// pop-over layers
type App struct {
	*Theme
	activePage          string
	bodyBackground      string
	bodyColor           string
	cardBackground      string
	cardColor           string
	buttonBar           []l.Widget
	hideSideBar         bool
	hideTitleBar        bool
	layers              []l.Widget
	logo                *[]byte
	logoClickable       *Clickable
	themeHook           func()
	menuBackground      string
	menuButton          *IconButton
	menuClickable       *Clickable
	menuColor           string
	menuIcon            *[]byte
	MenuOpen            bool
	pages               WidgetMap
	root                *Stack
	sideBar             []l.Widget
	sideBarBackground   string
	sideBarColor        string
	SideBarSize         unit.Value
	sideBarList         *List
	Size                *int
	statusBar           []l.Widget
	statusBarBackground string
	statusBarColor      string
	title               string
	titleBarBackground  string
	titleBarColor       string
	titleFont           string
	overlay             []func(gtx l.Context)
}

type WidgetMap map[string]l.Widget

func (th *Theme) App(size int) *App {
	mc := th.Clickable()
	return &App{
		Theme:               th,
		activePage:          "main",
		bodyBackground:      "PanelBg",
		bodyColor:           "PanelText",
		cardBackground:      "DocBg",
		cardColor:           "DocText",
		buttonBar:           nil,
		hideSideBar:         false,
		hideTitleBar:        false,
		layers:              nil,
		pages:               make(WidgetMap),
		root:                th.Stack(),
		SideBarSize:         th.TextSize.Scale(19),
		sideBarBackground:   "DocBg",
		sideBarColor:        "DocText",
		statusBarBackground: "DocBg",
		statusBarColor:      "DocText",
		sideBarList:         th.List(),
		logo:                &ico.ParallelCoin,
		logoClickable:       th.Clickable(),
		title:               "parallelcoin",
		titleBarBackground:  "Primary",
		titleBarColor:       "DocBg",
		titleFont:           "plan9",
		menuIcon:            &icons.NavigationMenu,
		menuClickable:       mc,
		menuButton:          th.IconButton(mc),
		menuColor:           "Light",
		MenuOpen:            false,
		Size:                &size,
	}
}

// Fn renders the app widget
func (a *App) Fn() func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		x := gtx.Constraints.Max.X
		a.Size = &x
		// TODO: put the root stack in here
		return a.Flex().Rigid(
			a.VFlex().
				Rigid(
					a.RenderHeader,
				).
				Flexed(1,
					a.MainFrame,
				).
				Rigid(
					a.RenderStatusBar,
				).
				Fn,
		).Fn(gtx)
	}
}
func (a *App) AddOverlay(overlay func(gtx l.Context)) *App {
	a.overlay = append(a.overlay, overlay)
	return a
}
func (a *App) Overlay() func(gtx l.Context) {
	return func(gtx l.Context) {
		for _, overlay := range a.overlay {
			overlay(gtx)
		}
	}
}

func (a *App) RenderStatusBar(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		bar := a.Flex().SpaceBetween()
		for x := range a.statusBar {
			i := x
			bar.Rigid(a.statusBar[i])
		}
		// out :=
		// a.Fill("PanelBg",
		// 	a.Inset(0.25,
		// 	bar.Fn
		// ).Fn
		// ).Fn
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		dims := a.Fill(a.statusBarBackground, bar.Fn).Fn(gtx)
		gtx.Constraints.Min = dims.Size
		gtx.Constraints.Max = dims.Size
		return dims
	}(gtx)
}

func (a *App) RenderHeader(gtx l.Context) l.Dimensions {
	return a.Flex().Flexed(1,
		a.Fill(a.titleBarBackground,
			a.Flex().
				Rigid(
					a.Responsive(*a.Size,
						Widgets{
							{Widget: a.MenuButton},
							{Size: 800, Widget: a.NoMenuButton}}).
						Fn,
				).
				Rigid(a.LogoAndTitle).
				Flexed(1,
					EmptyMinWidth(),
				).
				// Rigid(
				// 	a.DimensionCaption,
				// ).
				Rigid(
					a.RenderButtonBar,
				).
				Fn,
		).Fn,
	).Fn(gtx)
}

func (a *App) RenderButtonBar(gtx l.Context) l.Dimensions {
	out := a.Flex()
	for i := range a.buttonBar {
		out.Rigid(a.buttonBar[i])
	}
	dims := out.Fn(gtx)
	gtx.Constraints.Min = dims.Size
	gtx.Constraints.Max = dims.Size
	return dims
}

func (a *App) MainFrame(gtx l.Context) l.Dimensions {
	return a.Flex().
		Rigid(
			a.Flex().
				Rigid(
					a.Fill(a.sideBarBackground,
						a.Responsive(*a.Size, Widgets{
							{
								Widget: func(gtx l.Context) l.Dimensions {
									return If(a.MenuOpen,
										// a.Fill(a.sideBarBackground,
										a.renderSideBar(),
										// ).Fn,
										EmptySpace(0, 0),
									)(gtx)
								},
							},
							{Size: 800,
								Widget:
								// a.Fill(a.sideBarBackground,
								a.renderSideBar(),
								// ).Fn,
							},
						},
						).Fn,
					).Fn,
				).Fn,
		).
		Flexed(1,
			a.RenderPage,
		).
		Fn(gtx)
}

func (a *App) MenuButton(gtx l.Context) l.Dimensions {
	bg := a.titleBarBackground
	color := a.menuColor
	if a.MenuOpen {
		color = "DocText"
		bg = a.sideBarBackground
	}
	return a.Flex().Rigid(
		// a.Inset(0.25,
		a.ButtonLayout(a.menuClickable).
			CornerRadius(0).
			Embed(
				a.Inset(0.4,
					a.Icon().
						Scale(Scales["H5"]).
						Color(color).
						Src(&icons.NavigationMenu).
						Fn,
				).Fn,
			).
			Background(bg).
			SetClick(
				func() {
					a.MenuOpen = !a.MenuOpen
				}).
			Fn,
		// ).Fn,
	).Fn(gtx)
}

func (a *App) NoMenuButton(gtx l.Context) l.Dimensions {
	a.MenuOpen = false
	return l.Dimensions{}
}

func (a *App) LogoAndTitle(gtx l.Context) l.Dimensions {
	return a.Flex().
		Rigid(
			a.Responsive(*a.Size, Widgets{
				{
					Widget: EmptySpace(0, 0),
				},
				{Size: 800,
					Widget: a.Inset(0.25,
						a.IconButton(
							a.logoClickable.
								SetClick(
									func() {
										Debug("clicked logo")
										*a.Dark = !*a.Dark
										a.Theme.Colors.SetTheme(*a.Dark)
										a.themeHook()
									},
								),
						).
							Icon(
								a.Icon().
									Scale(Scales["H6"]).
									Color("Light").
									Src(a.logo)).
							Background("Dark").Color("Light").
							Inset(0.25).
							Fn,
					).Fn,
				},
			},
			).Fn,
		).
		Rigid(
			a.Responsive(*a.Size, Widgets{
				{Size: 800,
					Widget: a.Inset(0.333,
						a.H5(a.title).Color("Light").Fn,
					).Fn,
				},
				{
					Widget: a.ButtonLayout(a.logoClickable).Embed(
						a.Inset(0.333,
							a.H5(a.title).Color("Light").Fn,
						).Fn,
					).Background("Transparent").Fn,
				},
			}).Fn,
		).Fn(gtx)
}

func (a *App) RenderPage(gtx l.Context) l.Dimensions {
	return a.Fill(a.bodyBackground,
		func(gtx l.Context) l.Dimensions {
			if page, ok := a.pages[a.activePage]; !ok {
				return a.Flex().
					Flexed(1,
						a.Inset(0.5,
							a.VFlex().SpaceEvenly().
								Rigid(
									a.H1("404").
										Alignment(text.Middle).
										Fn,
								).
								Rigid(
									a.Body1("page not found").
										Alignment(text.Middle).
										Fn,
								).
								Fn,
						).Fn,
					).Fn(gtx)
			} else {
				// _ = page
				// return EmptyMaxHeight()(gtx)
				return page(gtx)
			}
		},
	).Fn(gtx)
}

func (a *App) DimensionCaption(gtx l.Context) l.Dimensions {
	return a.Caption(fmt.Sprintf("%dx%d", gtx.Constraints.Max.X, gtx.Constraints.Max.Y)).Fn(gtx)
}

func (a *App) renderSideBar() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		out := a.sideBarList.
			Length(len(a.sideBar)).
			LeftSide(true).
			Vertical().
			// Background("DocBg").
			// Color("DocText").
			// Active("Primary").
			ListElement(func(gtx l.Context, index int) l.Dimensions {
				// gtx.Constraints.Max.X = gtx.Constraints.Min.X
				dims := a.sideBar[index](gtx)
				// Debug(dims)
				return dims
				// out := a.VFlex()
				// for i := range a.sideBar {
				// 	out.Rigid(a.sideBar[i])
				// }
				// return out.Fn(gtx)
			})
		// out.Rigid(EmptySpace(int(a.SideBarSize.V), 0))
		return out.Fn(gtx)
	}
}

func (a *App) ActivePage(activePage string) *App {
	a.activePage = activePage
	return a
}
func (a *App) ActivePageGet() string {
	return a.activePage
}

func (a *App) BodyBackground(bodyBackground string) *App {
	a.bodyBackground = bodyBackground
	return a
}
func (a *App) BodyBackgroundGet() string {
	return a.bodyBackground
}

func (a *App) BodyColor(bodyColor string) *App {
	a.bodyColor = bodyColor
	return a
}
func (a *App) BodyColorGet() string {
	return a.bodyColor
}

func (a *App) CardBackground(cardBackground string) *App {
	a.cardBackground = cardBackground
	return a
}
func (a *App) CardBackgroundGet() string {
	return a.cardBackground
}

func (a *App) CardColor(cardColor string) *App {
	a.cardColor = cardColor
	return a
}
func (a *App) CardColorGet() string {
	return a.cardColor
}

func (a *App) ButtonBar(bar []l.Widget) *App {
	a.buttonBar = bar
	return a
}
func (a *App) ButtonBarGet() (bar []l.Widget) {
	return a.buttonBar
}

func (a *App) HideSideBar(hideSideBar bool) *App {
	a.hideSideBar = hideSideBar
	return a
}
func (a *App) HideSideBarGet() bool {
	return a.hideSideBar
}

func (a *App) HideTitleBar(hideTitleBar bool) *App {
	a.hideTitleBar = hideTitleBar
	return a
}
func (a *App) HideTitleBarGet() bool {
	return a.hideTitleBar
}

func (a *App) Layers(widgets []l.Widget) *App {
	a.layers = widgets
	return a
}
func (a *App) LayersGet() []l.Widget {
	return a.layers
}

func (a *App) MenuBackground(menuBackground string) *App {
	a.menuBackground = menuBackground
	return a
}
func (a *App) MenuBackgroundGet() string {
	return a.menuBackground
}

func (a *App) MenuColor(menuColor string) *App {
	a.menuColor = menuColor
	return a
}
func (a *App) MenuColorGet() string {
	return a.menuColor
}

func (a *App) MenuIcon(menuIcon *[]byte) *App {
	a.menuIcon = menuIcon
	return a
}
func (a *App) MenuIconGet() *[]byte {
	return a.menuIcon
}

func (a *App) Pages(widgets WidgetMap) *App {
	a.pages = widgets
	return a
}
func (a *App) PagesGet() WidgetMap {
	return a.pages
}

func (a *App) Root(root *Stack) *App {
	a.root = root
	return a
}
func (a *App) RootGet() *Stack {
	return a.root
}

func (a *App) SideBar(widgets []l.Widget) *App {
	a.sideBar = widgets
	return a
}
func (a *App) SideBarBackground(sideBarBackground string) *App {
	a.sideBarBackground = sideBarBackground
	return a
}
func (a *App) SideBarBackgroundGet() string {
	return a.sideBarBackground
}

func (a *App) SideBarColor(sideBarColor string) *App {
	a.sideBarColor = sideBarColor
	return a
}
func (a *App) SideBarColorGet() string {
	return a.sideBarColor
}

func (a *App) SideBarGet() []l.Widget {
	return a.sideBar
}

func (a *App) StatusBar(bar []l.Widget) *App {
	a.statusBar = bar
	return a
}
func (a *App) StatusBarBackground(statusBarBackground string) *App {
	a.statusBarBackground = statusBarBackground
	return a
}
func (a *App) StatusBarBackgroundGet() string {
	return a.statusBarBackground
}

func (a *App) StatusBarColor(statusBarColor string) *App {
	a.statusBarColor = statusBarColor
	return a
}
func (a *App) StatusBarColorGet() string {
	return a.statusBarColor
}

func (a *App) StatusBarGet() (bar []l.Widget) {
	return a.statusBar
}
func (a *App) Title(title string) *App {
	a.title = title
	return a
}
func (a *App) TitleBarBackground(TitleBarBackground string) *App {
	a.bodyBackground = TitleBarBackground
	return a
}
func (a *App) TitleBarBackgroundGet() string {
	return a.titleBarBackground
}

func (a *App) TitleBarColor(titleBarColor string) *App {
	a.titleBarColor = titleBarColor
	return a
}
func (a *App) TitleBarColorGet() string {
	return a.titleBarColor
}

func (a *App) TitleFont(font string) *App {
	a.titleFont = font
	return a
}
func (a *App) TitleFontGet() string {
	return a.titleFont
}
func (a *App) TitleGet() string {
	return a.title
}

func (a *App) ThemeHook(f func()) *App {
	a.themeHook = f
	return a
}
