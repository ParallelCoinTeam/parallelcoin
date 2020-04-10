package fui

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"image"
	"image/color"
)

type widget struct {
	ctx *layout.Context
	fn  []func(ctx *layout.Context)
}

// Widget returns a widget to chain paint/event ops into
func Widget() (out *widget) {
	out = &widget{}
	return
}

// Fill the widget's bounds with a given colour
func (w *widget) Fill(r, g, b, a byte) (out *widget) {
	col := color.RGBA{r, g, b, a}
	w.fn = append(w.fn, func(ctx *layout.Context) {
		cs := ctx.Constraints
		d := image.Point{
			X: cs.Width.Min,
			Y: cs.Height.Min,
		}
		dr := f32.Rectangle{
			Max: f32.Point{
				X: float32(d.X),
				Y: float32(d.Y),
			},
		}
		paint.ColorOp{Color: col}.Add(ctx.Ops)
		clip.Rect{Rect: dr}.Op(ctx.Ops).Add(ctx.Ops)
		paint.PaintOp{Rect: dr}.Add(ctx.Ops)
		ctx.Dimensions = layout.Dimensions{
			Size: d,
		}
	})
	return w
}

// Paint renders the widgets
func (w *widget) Paint(c *layout.Context) {
	for i := range w.fn {
		w.fn[i](c)
	}
}

// Prep the widgets to render later
func (w *widget) Prep(c *layout.Context) func() {
	return func() { w.Paint(c) }
}
