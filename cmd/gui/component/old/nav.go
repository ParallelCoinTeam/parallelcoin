package component

import (
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/gel"
	"github.com/p9c/pod/pkg/gui/gelook"
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
)

func MainNavigation(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, allPages *model.DuoUIpages, nav *model.DuoUInav) func() {
	return func() {
		navButtons := navButtons(rc, gtx, th, allPages, nav)
		gtx.Constraints.Width.Max = nav.Width
		th.DuoUIcontainer(0, th.Colors["Dark"]).Layout(gtx, layout.NW, func() {
			mainNav.Layout(gtx, len(navButtons), func(i int) {
				layout.UniformInset(unit.Dp(0)).Layout(gtx, navButtons[i])
			})
		})
	}
}

func navButtons(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, allPages *model.DuoUIpages, nav *model.DuoUInav) []func() {
	return []func(){
		navMenuButton(rc, gtx, th, allPages.Theme["OVERVIEW"], nav,
			"OVERVIEW", "overviewIcon", navButtonOverview),
		th.DuoUIline(gtx, 0, 0, 1,
			th.Colors["LightGrayIII"]),
		navMenuButton(rc, gtx, th, allPages.Theme["SEND"], nav,
			"SEND", "sendIcon", navButtonSend),
		// navMenuLine(gtx, th),
		// navMenuButton(rc, gtx, th, allPages.Theme["RECEIVE"], "RECEIVE",
		//"receiveIcon", navButtonReceive),
		th.DuoUIline(gtx, 0, 0, 1,
			th.Colors["LightGrayIII"]),
		navMenuButton(rc, gtx, th, allPages.Theme["ADDRESSBOOK"], nav,
			"ADDRESSBOOK", "addressBookIcon", navButtonAddressBook),
		th.DuoUIline(gtx, 0, 0, 1,
			th.Colors["LightGrayIII"]),
		navMenuButton(rc, gtx, th, allPages.Theme["HISTORY"], nav,
			"HISTORY", "historyIcon", navButtonHistory),
	}
}

func navMenuButton(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, page *gelook.DuoUIpage, nav *model.DuoUInav, title, icon string, navButton *gel.Button) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			var menuItem gelook.DuoUIbutton
			menuItem = th.DuoUIbutton(th.Fonts["Secondary"],
				title, th.Colors["Dark"],
				th.Colors["LightGrayII"],
				th.Colors["LightGrayII"],
				th.Colors["Dark"], icon,
				CurrentCurrentPageColor(rc.ShowPage, title, navItemIconColor,
					th.Colors["Primary"]), nav.TextSize, nav.IconSize,
				nav.Width, nav.Height, nav.PaddingVertical,
				nav.PaddingHorizontal, nav.PaddingVertical, nav.PaddingHorizontal)
			for navButton.Clicked(gtx) {
				rc.ShowPage = title
				SetPage(rc, page)
			}
			menuItem.MenuLayout(gtx, navButton)
		})
	}
}

//func navMenuLine(gtx *layout.Context, th *gelook.DuoUItheme) func() {
//	return func() {
//		gelook.DuoUIdrawRectangle(gtx, int(navItemWidth), 1,
//		th.Colors["LightGrayIII"], [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
//	}
//}
