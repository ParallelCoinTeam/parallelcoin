package p9

import (
	l "gioui.org/layout"
	"gioui.org/unit"
)

// App defines an application with a header, sidebar/menu, right side button bar, changeable body page widget and
// pop-over layers
type App struct {
	*Theme
	activePage         string
	bodyBackground     string
	bodyColor          string
	buttonBar          []l.Widget
	hideSideBar        bool
	hideTitleBar       bool
	layers             []l.Widget
	pages              map[string]l.Widget
	root               *Stack
	sideBar            []l.Widget
	sideBarSize        unit.Value
	title              string
	titleBarBackground string
	titleBarColor      string
	titleFont          string
	menuClickable      *Clickable
	menuButton         *IconButton
	menuIcon           []byte
	menuColor          string
	menuBackground     string
	responsive         *Responsive
}

func (th *Theme) App() *App {
	mc := th.Clickable()
	return &App{
		Theme:              th,
		activePage:         "",
		bodyBackground:     "DocBg",
		bodyColor:          "DocText",
		buttonBar:          nil,
		hideSideBar:        false,
		hideTitleBar:       false,
		layers:             nil,
		pages:              make(map[string]l.Widget),
		root:               th.Stack(),
		sideBar:            nil,
		sideBarSize:        th.TextSize.Scale(10),
		title:              "plan 9 from crypto space",
		titleBarBackground: "Primary",
		titleBarColor:      "DocBg",
		titleFont:          "plan9",
		menuClickable:      mc,
		menuButton:         th.IconButton(mc),
	}
}

// Fn renders the app widget
func (a *App) Fn(gtx l.Context) l.Dimensions {
	// barHeight := int(a.Theme.TextSize.Scale(3).V)
	return a.Flex().
		Rigid(
			a.Flex().Rigid(
				EmptySpace(int(a.sideBarSize.V), gtx.Constraints.Max.Y),
			).Fn,
		).
		Rigid(
			a.Flex().Vertical().
				Rigid(
					a.Flex().Flexed(1,
						a.Fill(a.titleBarBackground).Embed(
							a.Inset(0.5).Embed(
								a.H5(a.title).Fn,
							).Fn,
						).Fn,
					).Fn,
				).
				Flexed(1,
					a.Fill(a.bodyBackground).Embed(
						a.Flex().Flexed(1,
							a.Inset(0.5).Embed(
								a.Body1("body area").Fn,
							).Fn,
						).Fn,
					).Fn,
				).Fn,
		).
		Fn(gtx)
}

func (a *App) ActivePage(activePage string) {
	a.activePage = activePage
}
func (a *App) ActivePageGet() string {
	return a.activePage
}

func (a *App) BodyBackground(bodyBackground string) {
	a.bodyBackground = bodyBackground
}
func (a *App) BodyBackgroundGet() string {
	return a.bodyBackground
}

func (a *App) BodyColor(bodyColor string) {
	a.bodyColor = bodyColor
}
func (a *App) BodyColorGet() string {
	return a.bodyColor
}

func (a *App) ButtonBar(bar []l.Widget) {
	a.buttonBar = bar
}
func (a *App) ButtonBarGet() (bar []l.Widget) {
	return a.buttonBar
}

func (a *App) HideSideBar(hideSideBar bool) {
	a.hideSideBar = hideSideBar
}
func (a *App) HideSideBarGet() bool {
	return a.hideSideBar
}

func (a *App) HideTitleBar(hideTitleBar bool) {
	a.hideTitleBar = hideTitleBar
}
func (a *App) HideTitleBarGet() bool {
	return a.hideTitleBar
}

func (a *App) Layers(widgets []l.Widget) {
	a.layers = widgets
}
func (a *App) LayersGet() []l.Widget {
	return a.layers
}

func (a *App) MenuBackground(menuBackground string) {
	a.menuBackground = menuBackground
}
func (a *App) MenuBackgroundGet() string {
	return a.menuBackground
}

func (a *App) MenuColor(menuColor string) {
	a.menuColor = menuColor
}
func (a *App) MenuColorGet() string {
	return a.menuColor
}

func (a *App) MenuIcon(menuIcon []byte) {
	a.menuIcon = menuIcon
}
func (a *App) MenuIconGet() []byte {
	return a.menuIcon
}

func (a *App) Pages(widgets map[string]l.Widget) {
	a.pages = widgets
}
func (a *App) PagesGet() map[string]l.Widget {
	return a.pages
}

func (a *App) Root(root *Stack) {
	a.root = root
}
func (a *App) RootGet() *Stack {
	return a.root
}

func (a *App) SideBar(widgets []l.Widget) {
	a.sideBar = widgets
}
func (a *App) SideBarGet() []l.Widget {
	return a.sideBar
}

func (a *App) Title(title string) {
	a.title = title
}
func (a *App) TitleGet() string {
	return a.title
}

func (a *App) TitleBarBackground(TitleBarBackground string) {
	a.bodyBackground = TitleBarBackground
}
func (a *App) TitleBarBackgroundGet() string {
	return a.titleBarBackground
}

func (a *App) TitleBarColor(titleBarColor string) {
	a.titleBarColor = titleBarColor
}
func (a *App) TitleBarColorGet() string {
	return a.titleBarColor
}

func (a *App) TitleFont(font string) {
	a.titleFont = font
}
func (a *App) TitleFontGet() string {
	return a.titleFont
}
