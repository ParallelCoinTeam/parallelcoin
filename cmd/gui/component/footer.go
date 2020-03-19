package component

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
)

var (
	buttonLog        = new(gel.Button)
	buttonSettings   = new(gel.Button)
	buttonNetwork    = new(gel.Button)
	buttonBlocks     = new(gel.Button)
	buttonConsole    = new(gel.Button)
	buttonHelp       = new(gel.Button)
	navItemIconColor = "ffacacac"
	cornerNav        = &layout.List{
		Axis: layout.Horizontal,
	}
	footerNav = &layout.List{
		Axis: layout.Horizontal,
	}
	footerMenuItemWidth             = 48
	footerMenuItemHeight            = 48
	footerMenuItemTextSize          = 16
	footerMenuItemIconSize          = 32
	footerMenuItemPaddingVertical   = 8
	footerMenuItemPaddingHorizontal = 8
)

func footerMenuButton(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, page *gelook.DuoUIpage, text, icon string, footerButton *gel.Button) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			var footerMenuItem gelook.DuoUIbutton
			if icon != "" {
				footerMenuItem = th.DuoUIbutton("", "", "", "", "", th.Colors["Dark"], icon, CurrentCurrentPageColor(rc.ShowPage, page.Title, navItemIconColor, th.Colors["Primary"]), footerMenuItemTextSize, footerMenuItemIconSize, footerMenuItemWidth, footerMenuItemHeight, footerMenuItemPaddingVertical, footerMenuItemPaddingHorizontal)
				for footerButton.Clicked(gtx) {
					rc.ShowPage = page.Title
					page.Command()
					SetPage(rc, page)
				}
				footerMenuItem.IconLayout(gtx, footerButton)
			} else {
				footerMenuItem = th.DuoUIbutton(th.Fonts["Primary"], text, CurrentCurrentPageColor(rc.ShowPage, page.Title, th.Colors["Light"], th.Colors["Primary"]), "", "", "", "", "", footerMenuItemTextSize, footerMenuItemIconSize, 0, footerMenuItemHeight, footerMenuItemPaddingVertical, 0)
				for footerButton.Clicked(gtx) {
					rc.ShowPage = page.Title
					page.Command()
					SetPage(rc, page)
				}
				footerMenuItem.Layout(gtx, footerButton)
			}
		})
	}
}

func FooterLeftMenu(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, allPages *model.DuoUIpages) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			cornerButtons := []func(){
				QuitButton(rc, gtx, th),
				// footerMenuButton(rc, gtx, th, allPages.Theme["EXPLORER"], "BLOCKS: "+fmt.Sprint(rc.Status.Node.BlockCount), "", buttonBlocks),
				footerMenuButton(rc, gtx, th, allPages.Theme["LOG"], "LOG", "traceIcon", buttonLog),
			}
			cornerNav.Layout(gtx, len(cornerButtons), func(i int) {
				layout.UniformInset(unit.Dp(0)).Layout(gtx, cornerButtons[i])
			})
		})
	}
}

func FooterRightMenu(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, allPages *model.DuoUIpages) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			navButtons := []func(){
				// footerMenuButton(rc, gtx, th, allPages.Theme["NETWORK"], "", "networkIcon", buttonNetwork),
				// footerMenuButton(rc, gtx, th, allPages.Theme["NETWORK"], "CONNECTIONS: "+fmt.Sprint(rc.Status.Node.ConnectionCount), "", buttonNetwork),
				// footerMenuButton(rc, gtx, th, allPages.Theme["EXPLORER"], "", "DeviceWidgets", buttonBlocks),
				footerMenuButton(rc, gtx, th, allPages.Theme["NETWORK"], "", "network", buttonNetwork),
				footerMenuButton(rc, gtx, th, allPages.Theme["NETWORK"], "CONNECTIONS: "+fmt.Sprint(rc.Status.Node.ConnectionCount.Load()), "", buttonNetwork),
				footerMenuButton(rc, gtx, th, allPages.Theme["EXPLORER"], "", "DeviceWidgets", buttonBlocks),
				footerMenuButton(rc, gtx, th, allPages.Theme["EXPLORER"], "BLOCKS: "+fmt.Sprint(rc.Status.Node.BlockCount.Load()), "", buttonBlocks),
				footerMenuButton(rc, gtx, th, allPages.Theme["MINER"], "", "helpIcon", buttonHelp),
				footerMenuButton(rc, gtx, th, allPages.Theme["CONSOLE"], "", "consoleIcon", buttonConsole),
				footerMenuButton(rc, gtx, th, allPages.Theme["SETTINGS"], "", "settingsIcon", buttonSettings),
			}
			footerNav.Layout(gtx, len(navButtons), func(i int) {
				layout.UniformInset(unit.Dp(0)).Layout(gtx, navButtons[i])
			})
		})
	}
}
