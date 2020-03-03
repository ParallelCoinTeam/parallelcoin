package duoui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/cmd/gui/controller"
	"github.com/p9c/pod/cmd/gui/theme"
	"image/color"
)

var (
	logoButton = new(controller.Button)
	headerList = &layout.List{
		Axis:      layout.Horizontal,
		Alignment: layout.Start,
	}
)

func (ui *DuoUI) DuoUIheader() func() {
	return func() {

		headerListItems := []*theme.DuoUIicon{
			ui.ly.Theme.Icons["CommunicationImportExport"],
			ui.ly.Theme.Icons["NotificationNetworkCheck"],
			ui.ly.Theme.Icons["NotificationSync"],
			ui.ly.Theme.Icons["NotificationSyncDisabled"],
			ui.ly.Theme.Icons["NotificationSyncProblem"],
			ui.ly.Theme.Icons["NotificationVPNLock"],
			ui.ly.Theme.Icons["NotificationWiFi"],
			ui.ly.Theme.Icons["MapsLayers"],
			ui.ly.Theme.Icons["MapsLayersClear"],
			ui.ly.Theme.Icons["ImageTimer"],
			ui.ly.Theme.Icons["ImageRemoveRedEye"],
			ui.ly.Theme.Icons["DeviceSignalCellular0Bar"],
			ui.ly.Theme.Icons["DeviceWidgets"],
			ui.ly.Theme.Icons["ActionTimeline"],
			ui.ly.Theme.Icons["HardwareWatch"],
			ui.ly.Theme.Icons["HardwareKeyboardHide"],
		}

		layout.Flex{
			Axis:    layout.Horizontal,
			Spacing: layout.SpaceBetween,
		}.Layout(ui.ly.Context,
			layout.Rigid(func() {
				var logoMeniItem theme.DuoUIbutton
				logoMeniItem = ui.ly.Theme.DuoUIbutton("", "", "", ui.ly.Theme.Color.Dark, "", "", "logo", ui.ly.Theme.Color.Light, 16, 64, 96, 96, 8, 8)
				for logoButton.Clicked(ui.ly.Context) {
					ui.changeLightDark()
				}
				logoMeniItem.IconLayout(ui.ly.Context, logoButton)
			}),
			layout.Rigid(component.Label(ui.ly.Context, ui.ly.Theme, ui.ly.Theme.Font.Primary, 12, ui.ly.Theme.Color.Dark, ui.rc.Status.Wallet.Balance.Load()+" "+ui.rc.Settings.Abbrevation)),
			layout.Flexed(1, func() {
				headerList.Layout(ui.ly.Context, len(headerListItems), func(i int) {
					layout.UniformInset(unit.Dp(16)).Layout(ui.ly.Context,
						func() {
							headerListItems[i].Color = color.RGBA{A: 0xff, R: 0xcf, G: 0x30, B: 0x30}
							headerListItems[i].Layout(ui.ly.Context, unit.Dp(float32(32)))
						},
					)
				})
			}),
		)
	}
}

func (ui *DuoUI) changeLightDark() {
	light := ui.ly.Theme.Color.Light
	dark := ui.ly.Theme.Color.Dark
	lightGray := ui.ly.Theme.Color.LightGrayIII
	darkGray := ui.ly.Theme.Color.DarkGrayII
	ui.ly.Theme.Color.Light = dark
	ui.ly.Theme.Color.Dark = light
	ui.ly.Theme.Color.LightGrayIII = darkGray
	ui.ly.Theme.Color.DarkGrayII = lightGray
}
