package duoui

import (
	"gioui.org/layout"
	"github.com/p9c/pod/cmd/gui/mvc/component"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
)

var (
	logoButton = new(controller.Button)
)

func (ui *DuoUI) DuoUIheader() func() {
	return func() {
		layout.Flex{Axis: layout.Horizontal}.Layout(ui.ly.Context,
			layout.Rigid(func() {
				var logoMeniItem theme.DuoUIbutton
				logoMeniItem = ui.ly.Theme.DuoUIbutton("", "", "", ui.ly.Theme.Color.Dark, "", "", "logo", ui.ly.Theme.Color.Light, 16, 64, 96, 96, 8, 8)
				for logoButton.Clicked(ui.ly.Context) {
					ui.changeLightDark()
				}
				logoMeniItem.IconLayout(ui.ly.Context, logoButton)
			}),
			layout.Rigid(component.Label(ui.ly.Context, ui.ly.Theme, ui.ly.Theme.Font.Primary, 12, ui.ly.Theme.Color.Dark, ui.rc.Status.Wallet.Balance+" "+ui.rc.Settings.Abbrevation)),
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
