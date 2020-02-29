package component

import (
	"gioui.org/layout"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/cmd/gui/rcd"
)

func LogButton(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		var logMenuItem theme.DuoUIbutton
		logMenuItem = th.DuoUIbutton("", "", "", th.Color.Dark, "traceIcon", CurrentCurrentPageColor(rc.ShowPage, "LOG", th.Color.Light, th.Color.Primary), footerMenuItemTextSize, footerMenuItemIconSize, footerMenuItemWidth, footerMenuItemHeight, footerMenuItemPaddingVertical, footerMenuItemPaddingHorizontal)

		for buttonLog.Clicked(gtx) {
			rc.ShowPage = "LOG"
		}
		logMenuItem.IconLayout(gtx, buttonLog)
	}
}
