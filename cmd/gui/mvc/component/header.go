package component

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
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
				fill(gtx, theme.HexARGB(th.Color.Primary))
			}),
			layout.Stacked(func() {
				gtx.Constraints.Width.Min = hmin
				gtx.Constraints.Height.Min = vmin
				layout.UniformInset(unit.Dp(8)).Layout(gtx, b)
			}),
		)
	}
}
