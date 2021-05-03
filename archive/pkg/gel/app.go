package gel

import (
	"fmt"
	
	"go.uber.org/atomic"
	"golang.org/x/exp/shiny/materialdesign/icons"
	
	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/text"
	"github.com/p9c/gio/unit"
)

// App defines an application with a header, sidebar/menu, right side button bar, changeable body page widget and
// pop-over layers
type App struct {
	*Window
	activePage          *atomic.String
	invalidate          chan struct{}
	bodyBackground      string
	bodyColor           string
	pageBackground      string
	pageColor           string
	cardBackground      string
	cardColor           string
	buttonBar           []l.Widget
	hideSideBar         bool
	hideTitleBar        bool
	layers              []l.Widget
	Logo                *[]byte
	LogoClickable       *Clickable
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
	SideBarSize         *unit.Value
	sideBarList         *List
	Size                *atomic.Int32
	statusBar           []l.Widget
	statusBarRight      []l.Widget
	statusBarBackground string
	statusBarColor      string
	title               string
	titleBarBackground  string
	titleBarColor       string
	titleFont           string
	mainDirection       l.Direction
	PreRendering        bool
	Break1              float32
}

type WidgetMap map[string]l.Widget

func (w *Window) App(size *atomic.Int32, activePage *atomic.String, Break1 float32, ) *App {
	// mc := w.Clickable()
	a := &App{
		Window:              w,
		activePage:          activePage,
		bodyBackground:      "DocBg",
		bodyColor:           "DocText",
		pageBackground:      "PanelBg",
		pageColor:           "PanelText",
		cardBackground:      "DocBg",
		cardColor:           "DocText",
		buttonBar:           nil,
		hideSideBar:         false,
		hideTitleBar:        false,
		layers:              nil,
		pages:               make(WidgetMap),
		root:                w.Stack(),
		sideBarBackground:   "DocBg",
		sideBarColor:        "DocText",
		statusBarBackground: "DocBg",
		statusBarColor:      "DocText",
		Logo:                &icons.ActionSettingsApplications,
		title:               "gio elements application",
		titleBarBackground:  "Primary",
		titleBarColor:       "DocBg",
		titleFont:           "plan9",
		menuIcon:            &icons.NavigationMenu,
		menuColor:           "DocText",
		MenuOpen:            false,
		Size:                size,
		mainDirection:       l.Center + 1,
		Break1:              Break1,
		LogoClickable:       w.WidgetPool.GetClickable(),
		sideBarList:         w.WidgetPool.GetList(),
	}
	a.SideBarSize = &unit.Value{}
	return a
}

func (a *App) SetAppTitleText(title string) *App {
	a.title = title
	return a
}

func (a *App) AppTitleText() string {
	return a.title
}

func (a *App) SetLogo(logo *[]byte) *App {
	a.Logo = logo
	return a
}

func (a *App) GetLogo() string {
	return a.title
}

func (a *App) SetMainDirection(direction l.Direction) *App {
	a.mainDirection = direction
	return a
}

func (a *App) MainDirection() l.Direction {
	return a.mainDirection
}

// Fn renders the node widget
func (a *App) Fn() func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		return a.Fill(a.bodyBackground, l.Center, 0, 0,
			a.VFlex().
				Rigid(
					a.RenderHeader,
				).
				Flexed(
					1,
					a.MainFrame,
				).
				Rigid(
					a.RenderStatusBar,
				).
				Fn,
		).
			Fn(gtx)
	}
}

func (a *App) RenderStatusBar(gtx l.Context) l.Dimensions {
	bar := a.Flex()
	for x := range a.statusBar {
		i := x
		bar.Rigid(a.statusBar[i])
	}
	bar.Flexed(1, EmptyMaxWidth())
	for x := range a.statusBarRight {
		i := x
		bar.Rigid(a.statusBarRight[i])
	}
	return bar.Fn(gtx)
}

func (a *App) RenderHeader(gtx l.Context) l.Dimensions {
	return a.Flex().AlignMiddle().
		Rigid(
			a.Theme.Responsive(
				a.Size.Load(),
				Widgets{
					{Widget: If(len(a.sideBar) > 0, a.MenuButton, a.NoMenuButton)},
					{Size: a.Break1, Widget: a.NoMenuButton},
				},
			).
				Fn,
		).
		Flexed(
			1, If(
				float32(a.Width.Load()) >= a.TextSize.Scale(a.Break1).V,
				a.Direction().W().Embed(a.LogoAndTitle).Fn,
				a.Direction().Center().Embed(a.LogoAndTitle).Fn,
			),
		).
		Rigid(
			a.RenderButtonBar,
		).Fn(gtx)
}

func (a *App) RenderButtonBar(gtx l.Context) l.Dimensions {
	out := a.Theme.Flex()
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
			a.VFlex().
				Flexed(
					1,
					a.Responsive(
						a.Size.Load(), Widgets{
							{
								Widget: func(gtx l.Context) l.Dimensions {
									return If(
										a.MenuOpen,
										a.renderSideBar(),
										EmptySpace(0, 0),
									)(gtx)
								},
							},
							{
								Size: a.Break1,
								Widget:
								a.renderSideBar(),
							},
						},
					).Fn,
				).Fn,
		).
		Flexed(
			1,
			a.RenderPage,
		).
		Fn(gtx)
}

func (a *App) MenuButton(gtx l.Context) l.Dimensions {
	bg := "Transparent"
	color := a.menuColor
	if a.MenuOpen {
		color = "DocText"
		bg = a.sideBarBackground
	}
	return a.Theme.Flex().SpaceEvenly().AlignEnd().
		Rigid(
			a.ButtonLayout(a.menuClickable).
				CornerRadius(0).
				Embed(
					a.Inset(
						0.375,
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
					},
				).
				Fn,
		).Fn(gtx)
}

func (a *App) NoMenuButton(_ l.Context) l.Dimensions {
	a.MenuOpen = false
	return l.Dimensions{}
}

func (a *App) LogoAndTitle(gtx l.Context) l.Dimensions {
	return a.Theme.Responsive(
		a.Size.Load(), Widgets{
			{
				Widget: a.Theme.Flex().AlignMiddle().
					Rigid(
						a.
							Inset(
								0.25, a.
									IconButton(
										a.LogoClickable.
											SetClick(
												func() {
													D.Ln("clicked logo")
													a.Theme.Dark.Flip()
													a.Theme.Colors.SetDarkTheme(a.Theme.Dark.True())
												},
											),
									).
									Icon(
										a.Icon().
											Scale(Scales["H6"]).
											Color("DocText").
											Src(a.Logo),
									).
									Background("Transparent").
									Color("DocText").
									ButtonInset(0.25).
									Corners(0).
									Fn,
							).
							Fn,
					).
					Rigid(
						a.H5(a.ActivePageGet()).
							Color("DocText").Fn,
					).
					Fn,
			},
			{
				Size: a.Break1,
				Widget: a.Theme.Flex().AlignMiddle().
					Rigid(
						a.
							Inset(
								0.25, a.
									IconButton(
										a.LogoClickable.
											SetClick(
												func() {
													D.Ln("clicked logo")
													a.Theme.Dark.Flip()
													a.Theme.Colors.SetDarkTheme(a.Theme.Dark.True())
												},
											),
									).
									Icon(
										a.Icon().
											Scale(Scales["H6"]).
											Color("DocText").
											Src(a.Logo),
									).
									Background("Transparent").Color("DocText").
									ButtonInset(0.25).
									Corners(0).
									Fn,
							).
							Fn,
					).
					Rigid(
						a.H5(a.title).Color("DocText").Fn,
					).
					Fn,
			},
		},
	).Fn(gtx)
}

func (a *App) RenderPage(gtx l.Context) l.Dimensions {
	return a.Fill(
		a.pageBackground, l.Center, 0, 0, a.Inset(
			0.25,
			func(gtx l.Context) l.Dimensions {
				if page, ok := a.pages[a.activePage.Load()]; !ok {
					return a.Flex().
						Flexed(
							1,
							a.VFlex().SpaceEvenly().
								Rigid(
									a.H1("404").
										Alignment(text.Middle).
										Fn,
								).
								Rigid(
									a.Body1("page "+a.activePage.Load()+" not found").
										Alignment(text.Middle).
										Fn,
								).
								Fn,
						).Fn(gtx)
				} else {
					return page(gtx)
				}
			},
		).Fn,
	).Fn(gtx)
}

func (a *App) DimensionCaption(gtx l.Context) l.Dimensions {
	return a.Caption(fmt.Sprintf("%dx%d", gtx.Constraints.Max.X, gtx.Constraints.Max.Y)).Fn(gtx)
}

func (a *App) renderSideBar() l.Widget {
	if len(a.sideBar) > 0 {
		le := func(gtx l.Context, index int) l.Dimensions {
			dims := a.sideBar[index](gtx)
			return dims
		}
		return func(gtx l.Context) l.Dimensions {
			a.PreRendering = true
			gtx1 := CopyContextDimensionsWithMaxAxis(gtx, l.Horizontal)
			// generate the dimensions for all the list elements
			allDims := GetDimensionList(gtx1, len(a.sideBar), le)
			a.PreRendering = false
			max := 0
			for _, i := range allDims {
				if i.Size.X > max {
					max = i.Size.X
				}
			}
			a.SideBarSize.V = float32(max)
			gtx.Constraints.Max.X = max
			gtx.Constraints.Min.X = max
			out := a.VFlex().
				Rigid(
					a.sideBarList.
						Length(len(a.sideBar)).
						LeftSide(true).
						Vertical().
						ListElement(le).
						Fn,
				)
			return out.Fn(gtx)
		}
	} else {
		return EmptySpace(0, 0)
	}
}

func (a *App) ActivePage(activePage string) *App {
	a.Invalidate()
	a.activePage.Store(activePage)
	return a
}
func (a *App) ActivePageGet() string {
	return a.activePage.Load()
}
func (a *App) ActivePageGetAtomic() *atomic.String {
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

func (a *App) StatusBar(bar, barR []l.Widget) *App {
	a.statusBar = bar
	a.statusBarRight = barR
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

func (a *App) Placeholder(title string) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		return a.VFlex().
			AlignMiddle().
			SpaceSides().
			Rigid(
				a.Flex().
					Flexed(0.5, EmptyMaxWidth()).
					Rigid(
						a.H1(title).Fn,
					).
					Flexed(0.5, EmptyMaxWidth()).
					Fn,
			).
			Fn(gtx)
	}
}
