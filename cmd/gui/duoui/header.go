package duoui

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
)

var (
	logoButton = new(controller.Button)
)

func (ui *DuoUI) DuoUIheader() func() {
	return func() {
		width := ui.ly.Context.Constraints.Width.Max
		layout.Flex{Axis: layout.Horizontal}.Layout(ui.ly.Context,
			layout.Rigid(func() {
				var logoMeniItem theme.DuoUIbutton
				logoMeniItem = ui.ly.Theme.DuoUIbutton("", "", "", ui.ly.Theme.Color.Dark, "logo", ui.ly.Theme.Color.Light, 16, 64, 96, 96, 8, 8)
				for logoButton.Clicked(ui.ly.Context) {
					ui.changeLightDark()
				}
				logoMeniItem.IconLayout(ui.ly.Context, logoButton)
			}),
			layout.Flexed(1, func() {
				layout.Inset{Top: unit.Dp(24), Bottom: unit.Dp(8), Left: unit.Dp(0), Right: unit.Dp(4)}.Layout(ui.ly.Context, func() {
					currentPage := ui.ly.Theme.H4(ui.rc.ShowPage)
					currentPage.Color = theme.HexARGB(ui.ly.Theme.Color.Light)
					currentPage.Alignment = text.Start
					currentPage.Layout(ui.ly.Context)
				})
			}),
			layout.Rigid(func() {
				layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(16), Left: unit.Dp(16), Right: unit.Dp(4)}.Layout(ui.ly.Context, func() {
					balance := ui.ly.Theme.Body2(ui.rc.Status.Wallet.Balance +
						" " + ui.rc.Settings.Abbrevation)
					balance.Color = theme.HexARGB(ui.ly.Theme.Color.Light)
					balance.Font.Typeface = "bariol"
					balance.Alignment = text.End
					balance.Layout(ui.ly.Context)
				})
			}),
			layout.Rigid(func() {
				layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(16), Left: unit.Dp(16), Right: unit.Dp(4)}.Layout(ui.ly.Context, func() {
					balance := ui.ly.Theme.Body2("dimenzion: " + fmt.Sprint(width))
					balance.Color = theme.HexARGB(ui.ly.Theme.Color.Light)
					balance.Alignment = text.End
					balance.Font.Typeface = "bariol"
					balance.Layout(ui.ly.Context)
				})
			}))
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
