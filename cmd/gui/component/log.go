package component

import (
	"gioui.org/layout"
	"github.com/p9c/gelook"

	"github.com/p9c/pod/cmd/gui/rcd"
)

func iconButton(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, page *gelook.DuoUIpage) func() {
	return func() {
		var logMenuItem gelook.DuoUIbutton
		logMenuItem = th.DuoUIbutton("", "", "", th.Colors["Dark"], "", "", "traceIcon", CurrentCurrentPageColor(rc.ShowPage, "LOG", th.Colors["Light"], th.Colors["Primary"]), footerMenuItemTextSize, footerMenuItemIconSize, footerMenuItemWidth, footerMenuItemHeight, footerMenuItemPaddingVertical, footerMenuItemPaddingHorizontal)
		for buttonLog.Clicked(gtx) {
			SetPage(rc, page)
			rc.ShowPage = "LOG"
		}
		logMenuItem.IconLayout(gtx, buttonLog)
	}
}
