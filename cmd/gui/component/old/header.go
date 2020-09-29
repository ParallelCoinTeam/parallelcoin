package component

import (
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/gel"
	"github.com/p9c/pod/pkg/gui/gelook"
)

var buttonHeader = new(gel.Button)

// func ContentHeader(gtx *layout.Context, th *gelook.DuoUItheme, b func()) func() {
//	return func() {
//		hmin := gtx.Constraints.Width.Min
//		vmin := gtx.Constraints.Height.Min
//		layout.Stack{Alignment: layout.Center}.Layout(gtx,
//			layout.Expanded(func() {
//				clip.Rect{
//					Rect: f32.Rectangle{Max: f32.Point{
//						X: float32(gtx.Constraints.Width.Min),
//						Y: float32(gtx.Constraints.Height.Min),
//					}},
//				}.Op(gtx.Ops).Add(gtx.Ops)
//				fill(gtx, gelook.HexARGB(th.Colors["Primary"]))
//			}),
//			layout.Stacked(func() {
//				gtx.Constraints.Width.Min = hmin
//				gtx.Constraints.Height.Min = vmin
//				layout.UniformInset(unit.Dp(0)).Layout(gtx, b)
//			}),
//		)
//	}
// }

func HeaderMenu(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, allPages *model.DuoUIpages) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			headerNav := []func(){
				headerMenuButton(rc, gtx, th, "",
					"CommunicationImportExport", buttonHeader),
				headerMenuButton(rc, gtx, th, "",
					"NotificationNetworkCheck", buttonHeader),
				headerMenuButton(rc, gtx, th, "",
					"NotificationSync", buttonHeader),
				headerMenuButton(rc, gtx, th, "",
					"NotificationSyncDisabled", buttonHeader),
				headerMenuButton(rc, gtx, th, "",
					"NotificationSyncProblem", buttonHeader),
				headerMenuButton(rc, gtx, th, "",
					"NotificationVPNLock", buttonHeader),
				headerMenuButton(rc, gtx, th, "",
					"MapsLayers", buttonHeader),
				headerMenuButton(rc, gtx, th, "",
					"MapsLayersClear", buttonHeader),
				headerMenuButton(rc, gtx, th, "",
					"ImageTimer", buttonHeader),
				headerMenuButton(rc, gtx, th, "",
					"ImageRemoveRedEye", buttonHeader),
				headerMenuButton(rc, gtx, th, "",
					"DeviceSignalCellular0Bar", buttonHeader),
				headerMenuButton(rc, gtx, th, "",
					"ActionTimeline", buttonHeader),
				headerMenuButton(rc, gtx, th, "",
					"HardwareWatch", buttonHeader),
			}
			footerNav.Layout(gtx, len(headerNav), func(i int) {
				layout.UniformInset(unit.Dp(0)).Layout(gtx, headerNav[i])
			})
		})
	}
}

func headerMenuButton(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, text, icon string, headerButton *gel.Button) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			var footerMenuItem gelook.DuoUIbutton
			footerMenuItem = th.DuoUIbutton("", "",
				"", "", "",
				th.Colors["Dark"], icon,
				CurrentCurrentPageColor(rc.ShowPage, text, navItemIconColor,
					th.Colors["Primary"]),
				footerMenuItemTextSize, footerMenuItemIconSize,
				footerMenuItemWidth, footerMenuItemHeight,
				footerMenuItemPaddingVertical, footerMenuItemPaddingHorizontal,
				footerMenuItemPaddingVertical, footerMenuItemPaddingHorizontal)
			for headerButton.Clicked(gtx) {
				rc.ShowPage = text
			}
			footerMenuItem.IconLayout(gtx, headerButton)
		})
	}
}
