package component

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/controller"
	"github.com/p9c/pod/pkg/gui/theme"
)

var (
	buttonHeader = new(controller.Button)
)

func ContentHeader(gtx *layout.Context, th *theme.DuoUItheme, b func()) func() {
	return func() {
		hmin := gtx.Constraints.Width.Min
		vmin := gtx.Constraints.Height.Min
		layout.Stack{Alignment: layout.Center}.Layout(gtx,
			layout.Expanded(func() {
				clip.Rect{
					Rect: f32.Rectangle{Max: f32.Point{
						X: float32(gtx.Constraints.Width.Min),
						Y: float32(gtx.Constraints.Height.Min),
					}},
				}.Op(gtx.Ops).Add(gtx.Ops)
				fill(gtx, theme.HexARGB(th.Colors["Primary"]))
			}),
			layout.Stacked(func() {
				gtx.Constraints.Width.Min = hmin
				gtx.Constraints.Height.Min = vmin
				layout.UniformInset(unit.Dp(8)).Layout(gtx, b)
			}),
		)
	}
}

func HeaderMenu(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, allPages *model.DuoUIpages) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			headerNav := []func(){
				headerMenuButton(rc, gtx, th, "", "CommunicationImportExport", buttonHeader),
				headerMenuButton(rc, gtx, th, "", "NotificationNetworkCheck", buttonHeader),
				headerMenuButton(rc, gtx, th, "", "NotificationSync", buttonHeader),
				headerMenuButton(rc, gtx, th, "", "NotificationSyncDisabled", buttonHeader),
				headerMenuButton(rc, gtx, th, "", "NotificationSyncProblem", buttonHeader),
				headerMenuButton(rc, gtx, th, "", "NotificationVPNLock", buttonHeader),
				headerMenuButton(rc, gtx, th, "", "NotificationWiFi", buttonHeader),
				headerMenuButton(rc, gtx, th, "", "MapsLayers", buttonHeader),
				headerMenuButton(rc, gtx, th, "", "MapsLayersClear", buttonHeader),
				headerMenuButton(rc, gtx, th, "", "ImageTimer", buttonHeader),
				headerMenuButton(rc, gtx, th, "", "ImageRemoveRedEye", buttonHeader),
				headerMenuButton(rc, gtx, th, "", "DeviceSignalCellular0Bar", buttonHeader),
				headerMenuButton(rc, gtx, th, "", "DeviceWidgets", buttonHeader),
				headerMenuButton(rc, gtx, th, "", "ActionTimeline", buttonHeader),
				headerMenuButton(rc, gtx, th, "", "HardwareWatch", buttonHeader),
				headerMenuButton(rc, gtx, th, "", "HardwareKeyboardHide", buttonHeader),
			}
			footerNav.Layout(gtx, len(headerNav), func(i int) {
				layout.UniformInset(unit.Dp(0)).Layout(gtx, headerNav[i])
			})
		})
	}
}

func headerMenuButton(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, text, icon string, headerButton *controller.Button) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			var footerMenuItem theme.DuoUIbutton
			footerMenuItem = th.DuoUIbutton("", "", "", "", "", th.Colors["Dark"], icon, CurrentCurrentPageColor(rc.ShowPage, text, navItemIconColor, th.Colors["Primary"]), footerMenuItemTextSize, footerMenuItemIconSize, footerMenuItemWidth, footerMenuItemHeight, footerMenuItemPaddingVertical, footerMenuItemPaddingHorizontal)
			for headerButton.Clicked(gtx) {
				rc.ShowPage = text
			}
			footerMenuItem.IconLayout(gtx, headerButton)
		})
	}
}
