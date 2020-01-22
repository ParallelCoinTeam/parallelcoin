package helpers

import (
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/f32"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/op/clip"
	"github.com/p9c/pod/pkg/gui/op/paint"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/pkg/gui/widget"
	"image"
	"image/color"
)

var (
//topLabel   = "testtopLabel"
//lineEditor = &widget.Editor{
//	SingleLine: true,
//	Submit:     true,
//}
//list = &layout.List{
//	Axis: layout.Vertical,
//}
//kln = layout.UniformInset(unit.Dp(0))
//ln  = layout.UniformInset(unit.Dp(1))
//in  = layout.UniformInset(unit.Dp(16))
//button            = new(widget.Button)
)

func DuoUIinputField(duo *models.DuoUI, cx *conte.Xt, fieldName, fieldModel string, lineEditor *widget.Editor) {
	//var btn material.Button
	//fmt.Println("daj sta das", makeInstance(fieldModel))
	bgcol := hexARGB("ffe4e4e4")
	brcol := hexARGB("ff303030")
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
			fill(duo.DuoUIcontext, brcol)
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
							fill(duo.DuoUIcontext, bgcol)
						}),
						layout.Stacked(func() {
							duo.DuoUIcontext.Constraints.Width.Min = hmin
							duo.DuoUIcontext.Constraints.Height.Min = vmin
							layout.Align(layout.Center).Layout(duo.DuoUIcontext, func() {
								layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: unit.Dp(12), Right: unit.Dp(12)}.Layout(duo.DuoUIcontext, func() {
									//paint.ColorOp{Color: col}.Add(duo.DuoUIcontext.Ops)
									//widget.Label{}.Layout(duo.DuoUIcontext, btn.shaper, btn.Font, btn.Text)
									e := duo.DuoUItheme.DuoUIeditor(fieldName,fieldName)
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
	//cs := duo.DuoUIcontext.Constraints
	//DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, HexARGB("ff305530"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
	//
	//e := duo.DuoUItheme.DuoUIeditor(fieldName)
	//e.Font.Style = text.Italic
	//
	//e.Layout(duo.DuoUIcontext, lineEditor)
	//for _, e := range lineEditor.Events(duo.DuoUIcontext) {
	//	if _, ok := e.(widget.SubmitEvent); ok {
	//		//topLabel = e.Text
	//		lineEditor.SetText("")
	//	}
	//}

	//kln.Layout(duo.DuoUIcontext, func() {
	//	cs := duo.DuoUIcontext.Constraints
	//	DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, ne, nw, se, sw, inset)
	//	ln.Layout(duo.DuoUIcontext, func() {
	//		cs := duo.DuoUIcontext.Constraints
	//		DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, color.RGBA{A: 0xff, R: 0xf4, G: 0xf4, B: 0xf4}, ne, nw, se, sw, inset)
	//		//in.Layout(duo.DuoUIcontext, func() {
	//		//	e := duo.DuoUItheme.DuoUIeditor("Hint")
	//		//	e.Font.Style = text.Italic
	//		//	e.Font.Size = unit.Dp(24)
	//		//	e.Layout(duo.DuoUIcontext, lineEditor)
	//		//	for _, e := range lineEditor.Events(duo.DuoUIcontext) {
	//		//		if e, ok := e.(widget.SubmitEvent); ok {
	//		//			topLabel = e.Text
	//		//			lineEditor.SetText("")
	//		//		}
	//		//	}
	//		//})
	//	})
	//})
	//in := layout.UniformInset(inset)
	//in.Layout(duo.DuoUIcontext, func() {
	//	square := f32.Rectangle{
	//		Max: f32.Point{
	//			X: float32(w),
	//			Y: float32(h),
	//		},
	//	}
	//	paint.ColorOp{Color: color}.Add(duo.DuoUIcontext.Ops)
	//
	//	clip.Rect{Rect: square,
	//		NE: ne, NW: nw, SE: se, SW: sw}.Op(duo.DuoUIcontext.Ops).Add(duo.DuoUIcontext.Ops) // HLdraw
	//	paint.PaintOp{Rect: square}.Add(duo.DuoUIcontext.Ops)
	//	duo.DuoUIcontext.Dimensions = layout.Dimensions{Size: image.Point{X: w, Y: h}}
	//})
}

func fill(gtx *layout.Context, col color.RGBA) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d}
}
