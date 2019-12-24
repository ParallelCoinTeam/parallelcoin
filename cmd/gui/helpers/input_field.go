package helpers

import (
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/pkg/gio/f32"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/op/clip"
	"github.com/p9c/pod/pkg/gio/op/paint"
	"github.com/p9c/pod/pkg/gio/unit"
	"github.com/p9c/pod/pkg/gio/widget"
	"image"
	"image/color"
)

var (
	topLabel   = "testtopLabel"
	lineEditor = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	list = &layout.List{
		Axis: layout.Vertical,
	}
	kln = layout.UniformInset(unit.Dp(0))
	ln  = layout.UniformInset(unit.Dp(1))
	in  = layout.UniformInset(unit.Dp(16))
)

func DuoUIinputField(duo *models.DuoUI, w, h int, color color.RGBA, ne, nw, se, sw float32, inset unit.Value) {

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
	in := layout.UniformInset(inset)
	in.Layout(duo.Gc, func() {
		square := f32.Rectangle{
			Max: f32.Point{
				X: float32(w),
				Y: float32(h),
			},
		}
		paint.ColorOp{Color: color}.Add(duo.Gc.Ops)

		clip.Rect{Rect: square,
			NE: ne, NW: nw, SE: se, SW: sw}.Op(duo.Gc.Ops).Add(duo.Gc.Ops) // HLdraw
		paint.PaintOp{Rect: square}.Add(duo.Gc.Ops)
		duo.Gc.Dimensions = layout.Dimensions{Size: image.Point{X: w, Y: h}}
	})
}
