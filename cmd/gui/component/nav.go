package component

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/gel"
	"github.com/p9c/gelook"

	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
)

var (
	navButtonOverview    = new(gel.Button)
	navButtonSend        = new(gel.Button)
	navButtonReceive     = new(gel.Button)
	navButtonAddressBook = new(gel.Button)
	navButtonHistory     = new(gel.Button)
	mainNav              = &layout.List{
		Axis: layout.Vertical,
	}

	navItemWidth             = 96
	navItemHeight            = 72
	navItemTextSize          = 48
	navItemTconSize          = 36
	navItemPaddingVertical   = 8
	navItemPaddingHorizontal = 0
)

func MainNavigation(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, allPages *model.DuoUIpages) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			mainNav.Layout(gtx, len(navButtons(rc, gtx, th, allPages)), func(i int) {
				layout.UniformInset(unit.Dp(0)).Layout(gtx, navButtons(rc, gtx, th, allPages)[i])
			})
		})
	}
}

func navButtons(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, allPages *model.DuoUIpages) []func() {
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

func navMenuButton(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, page *gelook.DuoUIpage, title, icon string, navButton *gel.Button) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			var menuItem gelook.DuoUIbutton
			menuItem = th.DuoUIbutton(th.Fonts["Secondary"], title, th.Colors["Dark"], th.Colors["LightGrayII"], th.Colors["LightGrayII"], th.Colors["Dark"], icon, CurrentCurrentPageColor(rc.ShowPage, title, navItemIconColor, th.Colors["Primary"]), navItemTextSize, navItemTconSize, navItemWidth, navItemHeight, navItemPaddingVertical, navItemPaddingHorizontal)
			for navButton.Clicked(gtx) {
				rc.ShowPage = title
				page.Command()
				SetPage(rc, page)
			}
			menuItem.MenuLayout(gtx, navButton)
		})
	}
}

func navMenuLine(gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		gelook.DuoUIdrawRectangle(gtx, int(navItemWidth), 1, th.Colors["LightGrayIII"], [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
	}
}
