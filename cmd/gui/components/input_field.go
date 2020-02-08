package components

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/f32"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/op/clip"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/pkg/gui/widget"
)

func DuoUIinputField(duo *models.DuoUI, cx *conte.Xt, fieldName, fieldModel string, lineEditor *widget.Editor) {
	//var btn material.Button
	//fmt.Println("daj sta das", makeInstance(fieldModel))
	bgcol := helpers.HexARGB("ffe4e4e4")
	brcol := helpers.HexARGB("ff303030")
	hmin := duo.DuoUIcontext.Constraints.Width.Min
	vmin := duo.DuoUIcontext.Constraints.Height.Min
	layout.Stack{Alignment: layout.Center}.Layout(duo.DuoUIcontext,
		layout.Expanded(func() {
			rr := float32(duo.DuoUIcontext.Px(unit.Dp(4)))
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(duo.DuoUIcontext.Constraints.Width.Min),
					Y: float32(duo.DuoUIcontext.Constraints.Height.Min),
				}},
				NE: rr, NW: rr, SE: rr, SW: rr,
			}.Op(duo.DuoUIcontext.Ops).Add(duo.DuoUIcontext.Ops)
			helpers.DuoUIfill(duo.DuoUIcontext, brcol)
		}),
		layout.Stacked(func() {
			duo.DuoUIcontext.Constraints.Width.Min = hmin
			duo.DuoUIcontext.Constraints.Height.Min = vmin
			layout.Align(layout.Center).Layout(duo.DuoUIcontext, func() {
				layout.Inset{Top: unit.Dp(1), Bottom: unit.Dp(1), Left: unit.Dp(1), Right: unit.Dp(1)}.Layout(duo.DuoUIcontext, func() {

					layout.Stack{Alignment: layout.Center}.Layout(duo.DuoUIcontext,
						layout.Expanded(func() {
							rr := float32(duo.DuoUIcontext.Px(unit.Dp(4)))
							clip.Rect{
								Rect: f32.Rectangle{Max: f32.Point{
									X: float32(duo.DuoUIcontext.Constraints.Width.Min),
									Y: float32(duo.DuoUIcontext.Constraints.Height.Min),
								}},
								NE: rr, NW: rr, SE: rr, SW: rr,
							}.Op(duo.DuoUIcontext.Ops).Add(duo.DuoUIcontext.Ops)
							helpers.DuoUIfill(duo.DuoUIcontext, bgcol)
						}),
						layout.Stacked(func() {
							duo.DuoUIcontext.Constraints.Width.Min = hmin
							duo.DuoUIcontext.Constraints.Height.Min = vmin
							layout.Align(layout.Center).Layout(duo.DuoUIcontext, func() {
								layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: unit.Dp(12), Right: unit.Dp(12)}.Layout(duo.DuoUIcontext, func() {
									//paint.ColorOp{Color: col}.Add(duo.DuoUIcontext.Ops)
									//widget.Label{}.Layout(duo.DuoUIcontext, btn.shaper, btn.Font, btn.Text)
									e := duo.DuoUItheme.DuoUIeditor(fieldName, fieldName)
									e.Font.Style = text.Italic

									e.Layout(duo.DuoUIcontext, lineEditor)
									for _, e := range lineEditor.Events(duo.DuoUIcontext) {
										if _, ok := e.(widget.SubmitEvent); ok {
											//topLabel = e.Text
											lineEditor.SetText("")
										}
									}
								})
							})
						}),
					)
				})
			})
		}),
	)
}
