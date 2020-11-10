package p9

import (
	"image"

	l "gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type Filler struct {
	th  *Theme
	col string
	w   l.Widget
}

// Fill fills underneath a widget you can put over top of it
func (th *Theme) Fill(col string, w l.Widget) *Filler {
	return &Filler{th: th, col: col, w: w}
}

func (f *Filler) Embed(w l.Widget) *Filler {
	f.w = w
	return f
}

func (f *Filler) Fn(gtx l.Context) l.Dimensions {
	return f.th.Stack().
		Expanded(
			func(c l.Context) l.Dimensions {
				d := f.w(gtx).Size
				// cs := gtx.Constraints
				// d := cs.Min
				clipRect := image.Rectangle{
					Max: image.Point{X: d.X, Y: d.Y},
				}
				st := op.Push(gtx.Ops)
				clip.Rect(clipRect).Add(gtx.Ops)
				paint.ColorOp{Color: f.th.Colors.Get(f.col)}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				st.Pop()
				gtx.Constraints.Constrain(d)
				return l.Dimensions{
					Size: image.Point{X: d.X, Y: d.Y},
				}
			},
		).
		Stacked(f.w).
		Fn(gtx)
}
