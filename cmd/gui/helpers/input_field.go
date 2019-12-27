package helpers

import (
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/pkg/gio/f32"
	"github.com/p9c/pod/pkg/gio/io/pointer"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/op"
	"github.com/p9c/pod/pkg/gio/op/clip"
	"github.com/p9c/pod/pkg/gio/op/paint"
	"github.com/p9c/pod/pkg/gio/unit"
	"github.com/p9c/pod/pkg/gio/widget"
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
button            = new(widget.Button)
)

func DuoUIinputField(duo *models.DuoUI, fieldName string, lineEditor *widget.Editor) {
	//var btn material.Button

	col := HexARGB("ff303055")
	bgcol := HexARGB("ff553030")
	hmin := duo.Gc.Constraints.Width.Min
	vmin := duo.Gc.Constraints.Height.Min
	layout.Stack{Alignment: layout.Center}.Layout(duo.Gc,
		layout.Expanded(func() {
			rr := float32(duo.Gc.Px(unit.Dp(4)))
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(duo.Gc.Constraints.Width.Min),
					Y: float32(duo.Gc.Constraints.Height.Min),
				}},
				NE: rr, NW: rr, SE: rr, SW: rr,
			}.Op(duo.Gc.Ops).Add(duo.Gc.Ops)
			fill(duo.Gc, bgcol)
			for _, c := range button.History() {
				drawInk(duo.Gc, c)
			}
		}),
		layout.Stacked(func() {
			duo.Gc.Constraints.Width.Min = hmin
			duo.Gc.Constraints.Height.Min = vmin
			layout.Align(layout.Center).Layout(duo.Gc, func() {
				layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: unit.Dp(12), Right: unit.Dp(12)}.Layout(duo.Gc, func() {
					paint.ColorOp{Color: col}.Add(duo.Gc.Ops)
					//widget.Label{}.Layout(duo.Gc, btn.shaper, btn.Font, btn.Text)
				})
			})
			pointer.Rect(image.Rectangle{Max: duo.Gc.Dimensions.Size}).Add(duo.Gc.Ops)
			button.Layout(duo.Gc)
		}),
	)
	//cs := duo.Gc.Constraints
	//DuoUIdrawRectangle(duo.Gc, cs.Width.Max, cs.Height.Max, HexARGB("ff305530"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
	//
	//e := duo.Th.Editor(fieldName)
	//e.Font.Style = text.Italic
	//
	//e.Layout(duo.Gc, lineEditor)
	//for _, e := range lineEditor.Events(duo.Gc) {
	//	if _, ok := e.(widget.SubmitEvent); ok {
	//		//topLabel = e.Text
	//		lineEditor.SetText("")
	//	}
	//}

	//kln.Layout(duo.Gc, func() {
	//	cs := duo.Gc.Constraints
	//	DuoUIdrawRectangle(duo.Gc, cs.Width.Max, cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, ne, nw, se, sw, inset)
	//	ln.Layout(duo.Gc, func() {
	//		cs := duo.Gc.Constraints
	//		DuoUIdrawRectangle(duo.Gc, cs.Width.Max, cs.Height.Max, color.RGBA{A: 0xff, R: 0xf4, G: 0xf4, B: 0xf4}, ne, nw, se, sw, inset)
	//		//in.Layout(duo.Gc, func() {
	//		//	e := duo.Th.Editor("Hint")
	//		//	e.Font.Style = text.Italic
	//		//	e.Font.Size = unit.Dp(24)
	//		//	e.Layout(duo.Gc, lineEditor)
	//		//	for _, e := range lineEditor.Events(duo.Gc) {
	//		//		if e, ok := e.(widget.SubmitEvent); ok {
	//		//			topLabel = e.Text
	//		//			lineEditor.SetText("")
	//		//		}
	//		//	}
	//		//})
	//	})
	//})
	//in := layout.UniformInset(inset)
	//in.Layout(duo.Gc, func() {
	//	square := f32.Rectangle{
	//		Max: f32.Point{
	//			X: float32(w),
	//			Y: float32(h),
	//		},
	//	}
	//	paint.ColorOp{Color: color}.Add(duo.Gc.Ops)
	//
	//	clip.Rect{Rect: square,
	//		NE: ne, NW: nw, SE: se, SW: sw}.Op(duo.Gc.Ops).Add(duo.Gc.Ops) // HLdraw
	//	paint.PaintOp{Rect: square}.Add(duo.Gc.Ops)
	//	duo.Gc.Dimensions = layout.Dimensions{Size: image.Point{X: w, Y: h}}
	//})
}

func drawInk(gtx *layout.Context, c widget.Click) {
	d := gtx.Now().Sub(c.Time)
	t := float32(d.Seconds())
	const duration = 0.5
	if t > duration {
		return
	}
	t = t / duration
	var stack op.StackOp
	stack.Push(gtx.Ops)
	size := float32(gtx.Px(unit.Dp(700))) * t
	rr := size * .5
	col := byte(0xaa * (1 - t*t))
	ink := paint.ColorOp{Color: color.RGBA{A: col, R: col, G: col, B: col}}
	ink.Add(gtx.Ops)
	op.TransformOp{}.Offset(c.Position).Offset(f32.Point{
		X: -rr,
		Y: -rr,
	}).Add(gtx.Ops)
	clip.Rect{
		Rect: f32.Rectangle{Max: f32.Point{
			X: float32(size),
			Y: float32(size),
		}},
		NE: rr, NW: rr, SE: rr, SW: rr,
	}.Op(gtx.Ops).Add(gtx.Ops)
	paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{X: float32(size), Y: float32(size)}}}.Add(gtx.Ops)
	stack.Pop()
	op.InvalidateOp{}.Add(gtx.Ops)
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
