package duoui

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
)

var (
	logoButton = new(controller.Button)
)

func (ui *DuoUI) DuoUIheader() func() {
	return func() {
		layout.Flex{Axis: layout.Horizontal}.Layout(ui.ly.Context,
			layout.Rigid(func() {
				layout.Align(layout.Center).Layout(ui.ly.Context, func() {
					var logoMeniItem theme.DuoUIbutton
					logoMeniItem = ui.ly.Theme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", 64, 96, 96, 8, 8, ui.ly.Theme.Icons["logo"])
					for logoButton.Clicked(ui.ly.Context) {
						//d.mod.CurrentPage = "NETWORK"
					}
					logoMeniItem.Layout(ui.ly.Context, logoButton)
				})
			}),
			layout.Flexed(1, func() {
				layout.Align(layout.Start).Layout(ui.ly.Context, func() {
					layout.Inset{Top: unit.Dp(24), Bottom: unit.Dp(8), Left: unit.Dp(0), Right: unit.Dp(4)}.Layout(ui.ly.Context, func() {
						currentPage := ui.ly.Theme.H4(ui.rc.ShowPage)
						currentPage.Color = ui.ly.Theme.Color.Light
						currentPage.Alignment = text.Start
						currentPage.Layout(ui.ly.Context)
					})
				})
			}),
			layout.Rigid(func() {
				layout.Align(layout.Center).Layout(ui.ly.Context, func() {
					layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(16), Left: unit.Dp(16), Right: unit.Dp(4)}.Layout(ui.ly.Context, func() {
						balance := ui.ly.Theme.Body2(ui.rc.Status.Wallet.Balance +
							" " + ui.rc.Settings.Abbrevation)
						balance.Color = ui.ly.Theme.Color.Light
						balance.Font.Typeface = "bariol"
						balance.Alignment = text.End
						balance.Layout(ui.ly.Context)
					})
				})
			}),
			layout.Rigid(func() {
				layout.Align(layout.Center).Layout(ui.ly.Context, func() {
					layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(16), Left: unit.Dp(16), Right: unit.Dp(4)}.Layout(ui.ly.Context, func() {
						balance := ui.ly.Theme.Body2("dimenzion: " + fmt.Sprint(ui.ly.Context.Constraints.Width.Max))
						balance.Color = ui.ly.Theme.Color.Light
						balance.Alignment = text.End
						balance.Font.Typeface = "bariol"
						balance.Layout(ui.ly.Context)
					})
				})
			}))
	}
}
