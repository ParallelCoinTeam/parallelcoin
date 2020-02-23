package duoui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/ico"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"image"
)

// Main wallet screen
func (ui *DuoUI) DuoUIsplashScreen() {
	cs := ui.ly.Context.Constraints
	theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, cs.Height.Max, ui.ly.Theme.Color.Dark, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
	logo, _ := theme.NewDuoUIicon(ico.ParallelCoin)
	layout.Flex{Axis: layout.Vertical}.Layout(ui.ly.Context,
		layout.Rigid(func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(ui.ly.Context,
				layout.Rigid(func() {
					layout.UniformInset(unit.Dp(8)).Layout(ui.ly.Context, func() {
						size := ui.ly.Context.Px(unit.Dp(256)) - 2*ui.ly.Context.Px(unit.Dp(8))
						if logo != nil {
							logo.Color = theme.HexARGB(ui.ly.Theme.Color.Dark)
							logo.Color = theme.HexARGB(ui.ly.Theme.Color.Light)
							logo.Layout(ui.ly.Context, unit.Px(float32(size)))
						}
						ui.ly.Context.Dimensions = layout.Dimensions{
							Size: image.Point{X: size, Y: size},
						}
					})
				}),
				layout.Flexed(1, func() {
					layout.UniformInset(unit.Dp(60)).Layout(ui.ly.Context, func() {
						txt := ui.ly.Theme.H1("PLAN NINE FROM FAR, FAR AWAY SPACE")
						txt.Color = theme.HexARGB(ui.ly.Theme.Color.Light)
						txt.Layout(ui.ly.Context)
					})
				}),
			)
		}),
		layout.Flexed(1, ui.DuoUIlogger()),
	)
}

// Main wallet screen
func (ui *DuoUI) DuoUImainScreen() {
	cs := ui.ly.Context.Constraints
	theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, cs.Height.Max, ui.ly.Theme.Color.Dark, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
	layout.Flex{Axis: layout.Vertical}.Layout(ui.ly.Context,
		layout.Rigid(ui.DuoUIheader()),
		layout.Flexed(1, ui.DuoUIbody()),
		layout.Rigid(ui.DuoUIfooter()),
	)
}
