package component

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/controller"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/cmd/gui/theme"
	"github.com/p9c/pod/pkg/util/interrupt"
)

var (
	buttonQuit = new(controller.Button)
)

func QuitButton(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			var closeMeniItem theme.DuoUIbutton
			closeMeniItem = th.DuoUIbutton("", "", "", th.Color.Dark, "", "", "closeIcon", CurrentCurrentPageColor(rc.ShowPage, "CLOSE", th.Color.Light, th.Color.Primary), footerMenuItemTextSize, footerMenuItemIconSize, footerMenuItemWidth, footerMenuItemHeight, footerMenuItemPaddingVertical, footerMenuItemPaddingHorizontal)
			for buttonQuit.Clicked(gtx) {
				rc.Dialog.Show = true
				rc.Dialog = &model.DuoUIdialog{
					Show: true,
					Ok: func() {
						interrupt.Request()
						// TODO make this close the window or at least switch to a shutdown screen
						rc.Dialog.Show = false
					},
					Close: func() {
						interrupt.RequestRestart()
						// TODO make this close the window or at least switch to a shutdown screen
						rc.Dialog.Show = false
					},
					Cancel: func() { rc.Dialog.Show = false },
					Title:  "Are you sure?",
					Text:   "Confirm ParallelCoin close",
				}
			}
			closeMeniItem.IconLayout(gtx, buttonQuit)
		})
	}
}
