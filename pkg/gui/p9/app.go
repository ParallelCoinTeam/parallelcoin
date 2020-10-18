package p9

import (
	l "gioui.org/layout"
)

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
		title:              "plan 9 from crypto space",
		titleBarBackground: "Primary",
		titleBarColor:      "DocBg",
		titleFont:          "plan9",
		menuClickable:      mc,
		menuButton:         th.IconButton(mc),
	}
}

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
	title              string
	titleBarBackground string
	titleBarColor      string
	titleFont          string
	menuClickable      *Clickable
	menuButton         *IconButton
	menuIcon           []byte
	menuColor          string
	menuBackground     string
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

func (a *App) Root(root l.Stack) {
	a.root = root
}
func (a *App) RootGet() l.Stack {
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
