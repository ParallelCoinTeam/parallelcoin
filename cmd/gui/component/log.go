package component

import (
	"gioui.org/layout"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/theme"
)

func iconButton(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, page *theme.DuoUIpage) func() {
	return func() {
		var logMenuItem theme.DuoUIbutton
		logMenuItem = th.DuoUIbutton("", "", "", th.Colors["Dark"], "", "", "traceIcon", CurrentCurrentPageColor(rc.ShowPage, "LOG", th.Colors["Light"], th.Colors["Primary"]), footerMenuItemTextSize, footerMenuItemIconSize, footerMenuItemWidth, footerMenuItemHeight, footerMenuItemPaddingVertical, footerMenuItemPaddingHorizontal)
		for buttonLog.Clicked(gtx) {
			SetPage(rc, page)
			rc.ShowPage = "LOG"
		}
		logMenuItem.IconLayout(gtx, buttonLog)
	}
}
