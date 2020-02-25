package duoui

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
)

func  (ui *DuoUI)DuoUIcontent() func() {
	return func() {
		ui.ly.Pages[ui.rc.ShowPage].Layout(ui.ly.Context)
	}
}


func (ui *DuoUI) contentHeader(b func()) func() {
	return func() {
		hmin := ui.ly.Context.Constraints.Width.Min
		vmin := ui.ly.Context.Constraints.Height.Min
		layout.Stack{Alignment: layout.Center}.Layout(ui.ly.Context,
			layout.Expanded(func() {
				clip.Rect{
					Rect: f32.Rectangle{Max: f32.Point{
						X: float32(ui.ly.Context.Constraints.Width.Min),
						Y: float32(ui.ly.Context.Constraints.Height.Min),
					}},
				}.Op(ui.ly.Context.Ops).Add(ui.ly.Context.Ops)
				fill(ui.ly.Context, theme.HexARGB(ui.ly.Theme.Color.Primary))
			}),
			layout.Stacked(func() {
				ui.ly.Context.Constraints.Width.Min = hmin
				ui.ly.Context.Constraints.Height.Min = vmin
				layout.UniformInset(unit.Dp(8)).Layout(ui.ly.Context, b)
			}),
		)
	}
}
