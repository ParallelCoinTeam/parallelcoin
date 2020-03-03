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
		Axis:      layout.Vertical,
		Alignment: layout.Start,
	}
	headerListItems = []*theme.DuoUIicon{
		//""
	}
)

func (ui *DuoUI) DuoUIheader() func() {
	return func() {
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
			layout.Rigid(component.Label(ui.ly.Context, ui.ly.Theme, ui.ly.Theme.Font.Primary, 12, ui.ly.Theme.Color.Light, ui.rc.Status.Wallet.Balance+" "+ui.rc.Settings.Abbrevation)),
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
