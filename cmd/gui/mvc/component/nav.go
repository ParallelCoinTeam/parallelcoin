package component

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/cmd/gui/rcd"
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
)

func MainNavigation(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, allPages *model.DuoUIpages) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			mainNav.Layout(gtx, len(navButtons(rc, gtx, th, allPages)), func(i int) {
				layout.UniformInset(unit.Dp(0)).Layout(gtx, navButtons(rc, gtx, th, allPages)[i])
			})
		})
	}
}

func navButtons(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, allPages *model.DuoUIpages) []func() {
	return []func(){
		navMenuButton(rc, gtx, th, allPages.Theme["OVERVIEW"], "OVERVIEW", "overviewIcon", navButtonOverview),
		navMenuLine(gtx, th),
		navMenuButton(rc, gtx, th, allPages.Theme["SEND"], "SEND", "sendIcon", navButtonSend),
		navMenuLine(gtx, th),
		navMenuButton(rc, gtx, th, allPages.Theme["RECEIVE"], "RECEIVE", "receiveIcon", navButtonReceive),
		navMenuLine(gtx, th),
		navMenuButton(rc, gtx, th, allPages.Theme["ADDRESSBOOK"], "ADDRESSBOOK", "addressBookIcon", navButtonAddressBook),
		navMenuLine(gtx, th),
		navMenuButton(rc, gtx, th, allPages.Theme["HISTORY"], "HISTORY", "historyIcon", navButtonHistory),
	}
}

func navMenuButton(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, page *theme.DuoUIpage, title, icon string, navButton *controller.Button) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			var menuItem theme.DuoUIbutton
			menuItem = th.DuoUIbutton(th.Font.Secondary, title, th.Color.Dark, th.Color.LightGrayII, th.Color.LightGrayII, th.Color.Dark, icon, CurrentCurrentPageColor(rc.ShowPage, title, navItemIconColor, th.Color.Primary), navItemTextSize, navItemTconSize, navItemWidth, navItemHeight, navItemPaddingVertical, navItemPaddingHorizontal)
			for navButton.Clicked(gtx) {
				rc.ShowPage = title
				SetPage(rc, page)
			}
			menuItem.MenuLayout(gtx, navButton)
		})
	}
}

func navMenuLine(gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		theme.DuoUIdrawRectangle(gtx, int(navItemWidth), 1, th.Color.LightGrayIII, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
	}
}
