package component

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/controller"
	"github.com/p9c/pod/pkg/gui/theme"
)

var (
	buttonLog        = new(controller.Button)
	buttonSettings   = new(controller.Button)
	buttonNetwork    = new(controller.Button)
	buttonBlocks     = new(controller.Button)
	buttonConsole    = new(controller.Button)
	buttonHelp       = new(controller.Button)
	navItemIconColor = "ffacacac"
	cornerNav        = &layout.List{
		Axis: layout.Horizontal,
	}
	footerNav = &layout.List{
		Axis: layout.Horizontal,
	}
	footerMenuItemWidth             int = 48
	footerMenuItemHeight            int = 48
	footerMenuItemTextSize          int = 16
	footerMenuItemIconSize          int = 32
	footerMenuItemPaddingVertical   int = 8
	footerMenuItemPaddingHorizontal int = 8
)

func footerMenuButton(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, page *theme.DuoUIpage, text, icon string, footerButton *controller.Button) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			var footerMenuItem theme.DuoUIbutton
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

func FooterLeftMenu(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, allPages *model.DuoUIpages) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			cornerButtons := []func(){
				QuitButton(rc, gtx, th),
				//footerMenuButton(rc, gtx, th, allPages.Theme["EXPLORER"], "BLOCKS: "+fmt.Sprint(rc.Status.Node.BlockCount), "", buttonBlocks),
				footerMenuButton(rc, gtx, th, allPages.Theme["LOG"], "LOG", "traceIcon", buttonLog),
			}
			cornerNav.Layout(gtx, len(cornerButtons), func(i int) {
				layout.UniformInset(unit.Dp(0)).Layout(gtx, cornerButtons[i])
			})
		})
	}
}

func FooterRightMenu(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, allPages *model.DuoUIpages) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			navButtons := []func(){
<<<<<<< HEAD
				//footerMenuButton(rc, gtx, th, allPages.Theme["NETWORK"], "", "networkIcon", buttonNetwork),
				footerMenuButton(rc, gtx, th, allPages.Theme["NETWORK"], "CONNECTIONS: "+fmt.Sprint(rc.Status.Node.ConnectionCount), "", buttonNetwork),
				//footerMenuButton(rc, gtx, th, allPages.Theme["EXPLORER"], "", "DeviceWidgets", buttonBlocks),
=======
				footerMenuButton(rc, gtx, th, allPages.Theme["NETWORK"], "", "settingsIcon", buttonSettings),
				footerMenuButton(rc, gtx, th, allPages.Theme["NETWORK"],
					"CONNECTIONS: "+fmt.Sprint(rc.Status.Node.ConnectionCount.Load()), "", buttonNetwork),
				footerMenuButton(rc, gtx, th, allPages.Theme["EXPLORER"], "", "DeviceWidgets", buttonSettings),
>>>>>>> 12f1f46da1ab69d1b64d41693cefff7e0254953a
				footerMenuButton(rc, gtx, th, allPages.Theme["EXPLORER"], "BLOCKS: "+fmt.Sprint(rc.Status.Node.BlockCount), "", buttonBlocks),
				footerMenuButton(rc, gtx, th, allPages.Theme["MINER"], "", "helpIcon", buttonHelp),
				footerMenuButton(rc, gtx, th, allPages.Theme["CONSOLE"], "", "HardwareKeyboardHide", buttonConsole),
				footerMenuButton(rc, gtx, th, allPages.Theme["SETTINGS"], "", "settingsIcon", buttonSettings),
			}
			footerNav.Layout(gtx, len(navButtons), func(i int) {
				layout.UniformInset(unit.Dp(0)).Layout(gtx, navButtons[i])
			})
		})
	}
}
