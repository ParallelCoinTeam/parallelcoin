package component

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/gel"
	"github.com/p9c/gelook"
	"github.com/p9c/util/interrupt"

	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
)

var (
	buttonQuit = new(gel.Button)
)

func QuitButton(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			var closeMeniItem gelook.DuoUIbutton
			closeMeniItem = th.DuoUIbutton("", "", "", th.Colors["Dark"], "", "", "closeIcon", CurrentCurrentPageColor(rc.ShowPage, "CLOSE", th.Colors["Light"], th.Colors["Primary"]), footerMenuItemTextSize, footerMenuItemIconSize, footerMenuItemWidth, footerMenuItemHeight, footerMenuItemPaddingVertical, footerMenuItemPaddingHorizontal)
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
					Cancel:      func() { rc.Dialog.Show = false },
					CustomField: func() {},
					Title:       "Are you sure?",
					Text:        "Confirm ParallelCoin close",
				}
			}
			closeMeniItem.IconLayout(gtx, buttonQuit)
		})
	}
}
