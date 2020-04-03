package component

import (
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/gel"
	"github.com/p9c/pod/pkg/gui/gelook"
	"github.com/p9c/pod/pkg/util/interrupt"
)

var (
	buttonQuit = new(gel.Button)
)

func QuitButton(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			var closeMeniItem gelook.DuoUIbutton
			closeMeniItem = th.DuoUIbutton("", "", "", th.Colors["Dark"], "", "", "closeIcon", CurrentCurrentPageColor(rc.ShowPage, "CLOSE", th.Colors["Light"], th.Colors["Primary"]), footerMenuItemTextSize, footerMenuItemIconSize, footerMenuItemWidth, footerMenuItemHeight, 0, 0, 0, 0)
			for buttonQuit.Clicked(gtx) {
				rc.Dialog.Show = true
				rc.Dialog = &model.DuoUIdialog{
					Show: true,
					Green: func() {
						interrupt.Request()
						// TODO make this close the window or at least switch to a shutdown screen
						rc.Dialog.Show = false
					},
					GreenLabel: "QUIT",
					Orange: func() {
						interrupt.RequestRestart()
						// TODO make this close the window or at least switch to a shutdown screen
						rc.Dialog.Show = false
					},
					OrangeLabel: "RESTART",
					Red:         func() { rc.Dialog.Show = false },
					RedLabel:    "CANCEL",
					CustomField: func() {},
					Title:       "Are you sure?",
					Text:        "Confirm ParallelCoin close",
				}
			}
			closeMeniItem.IconLayout(gtx, buttonQuit)
		})
	}
}
