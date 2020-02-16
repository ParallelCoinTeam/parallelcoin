package view

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/pkg/gui/f32"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/op/clip"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
)

func DuoUIinputField(duo *model.DuoUI, fieldName, fieldModel string, lineEditor *controller.Editor) func() {

	return func() {
		//var btn material.Button
		//fmt.Println("daj sta das", makeInstance(fieldModel))
		bgcol := theme.HexARGB("ffe4e4e4")
		brcol := theme.HexARGB("ff303030")
		hmin := duo.Context.Constraints.Width.Min
		vmin := duo.Context.Constraints.Height.Min
		layout.Stack{Alignment: layout.Center}.Layout(duo.Context,
			layout.Expanded(func() {
				rr := float32(duo.Context.Px(unit.Dp(4)))
				clip.Rect{
					Rect: f32.Rectangle{Max: f32.Point{
						X: float32(duo.Context.Constraints.Width.Min),
						Y: float32(duo.Context.Constraints.Height.Min),
					}},
					NE: rr, NW: rr, SE: rr, SW: rr,
				}.Op(duo.Context.Ops).Add(duo.Context.Ops)
				helpers.DuoUIfill(duo.Context, brcol)
			}),
			layout.Stacked(func() {
				duo.Context.Constraints.Width.Min = hmin
				duo.Context.Constraints.Height.Min = vmin
					layout.Inset{Top: unit.Dp(1), Bottom: unit.Dp(1), Left: unit.Dp(1), Right: unit.Dp(1)}.Layout(duo.Context, func() {

						layout.Stack{Alignment: layout.Center}.Layout(duo.Context,
							layout.Expanded(func() {
								rr := float32(duo.Context.Px(unit.Dp(4)))
								clip.Rect{
									Rect: f32.Rectangle{Max: f32.Point{
										X: float32(duo.Context.Constraints.Width.Min),
										Y: float32(duo.Context.Constraints.Height.Min),
									}},
									NE: rr, NW: rr, SE: rr, SW: rr,
								}.Op(duo.Context.Ops).Add(duo.Context.Ops)
								helpers.DuoUIfill(duo.Context, bgcol)
							}),
							layout.Stacked(func() {
								duo.Context.Constraints.Width.Min = hmin
								duo.Context.Constraints.Height.Min = vmin
									layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: unit.Dp(12), Right: unit.Dp(12)}.Layout(duo.Context, func() {
										//paint.ColorOp{Color: col}.Add(duo.Context.Ops)
										//widget.Label{}.Layout(duo.Context, btn.shaper, btn.Font, btn.Text)
										e := duo.Theme.DuoUIeditor(fieldName)
										e.Font.Style = text.Italic

										e.Layout(duo.Context, lineEditor)
										for _, e := range lineEditor.Events(duo.Context) {
											if _, ok := e.(controller.SubmitEvent); ok {
												//topLabel = e.Text
												lineEditor.SetText("")
											}
										}
									})
								}),
						)
					})
			}),
		)
	}
}