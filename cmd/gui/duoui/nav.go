package duoui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
)

var (
	navButtonOverview    = new(controller.Button)
	navButtonSend        = new(controller.Button)
	navButtonReceive     = new(controller.Button)
	navButtonAddressBook = new(controller.Button)
	navButtonHistory     = new(controller.Button)
	mainNav              = &layout.List{
		Axis: layout.Vertical,
	}

	navItemWidth             int = 96
	navItemHeight            int = 72
	navItemTextSize          int = 48
	navItemTconSize          int = 36
	navItemPaddingVertical   int = 8
	navItemPaddingHorizontal int = 0
	navItemIconColor             = "ffacacac"
)

func (ui *DuoUI) DuoUImenu() func() {
	return func() {

		//overviewButton :=
		//historyButton :=

		in := layout.UniformInset(unit.Dp(0))

		layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
			Spacing:   layout.SpaceEvenly}.Layout(ui.ly.Context,
			layout.Rigid(func() {
				layout.Flex{}.Layout(ui.ly.Context,
					layout.Rigid(func() {
						in.Layout(ui.ly.Context, func() {

							navButtons := []func(){
								ui.navMenuButton("OVERVIEW", "overviewIcon", navButtonOverview),
								ui.navMenuLine(),
								ui.navMenuButton("SEND", "sendIcon", navButtonSend),
								ui.navMenuLine(),
								ui.navMenuButton("RECEIVE", "receiveIcon", navButtonReceive),
								ui.navMenuLine(),
								ui.navMenuButton("ADDRESSBOOK", "addressBookIcon", navButtonAddressBook),
								ui.navMenuLine(),
								ui.navMenuButton("HISTORY", "historyIcon", navButtonHistory),
							}
							mainNav.Layout(ui.ly.Context, len(navButtons), func(i int) {
								layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, navButtons[i])
							})
						})
					}),
				)
			}),
		)
	}
}

func currentCurrentPageColor(showPage, page, color, currentPageColor string) (c string) {
	if showPage == page {
		c = currentPageColor
	} else {
		c = color
	}
	return
}

func (ui *DuoUI) navMenuButton(title, icon string, navButton *controller.Button) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
			var menuItem theme.DuoUIbutton
			menuItem = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Secondary, title, ui.ly.Theme.Color.Dark, ui.ly.Theme.Color.LightGrayII, icon, currentCurrentPageColor(ui.rc.ShowPage, title, navItemIconColor, ui.ly.Theme.Color.Primary), navItemTextSize, navItemTconSize, navItemWidth, navItemHeight, navItemPaddingVertical, navItemPaddingHorizontal)
			for navButton.Clicked(ui.ly.Context) {
				ui.rc.ShowPage = title
			}
			menuItem.MenuLayout(ui.ly.Context, navButton)
		})
	}
}

func (ui *DuoUI) navMenuLine() func() {
	return func() {
		theme.DuoUIdrawRectangle(ui.ly.Context, int(navItemWidth), 1, ui.ly.Theme.Color.LightGrayIII, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
	}
}
