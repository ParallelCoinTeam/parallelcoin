package p9

import (
	"image"
	"image/color"
	
	"gioui.org/f32"
	l "gioui.org/layout"
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
	return f.th.Stack().Stacked(f.w).Expanded(
		func(c l.Context) l.Dimensions {
			// gtx1 := CopyContextDimensionsWithMaxAxis(gtx, gtx.Constraints.Max, l.Vertical)
			// render the widgets onto a second context to get their dimensions
			gtx1 := CopyContextDimensionsWithMaxAxis(gtx, gtx.Constraints.Max, l.Horizontal)
			// generate the dimensions for all the list elements
			dd := GetDimensionList(
				gtx1, 1, func(gtx l.Context, index int) l.Dimensions {
					return f.w(gtx)
				},
			)
			dims := dd[0]
			if f.col != "" {
				fill(gtx, f.th.Colors.Get(f.col), dims.Size)
			}
			// // dims := f.w(gtx)
			// cs := gtx.Constraints
			// d := image.Point{X: cs.Max.X, Y: cs.Max.Y}
			// dr := f32.Rectangle{
			// 	Max: f32.Point{X: float32(dims.Size.X), Y: float32(dims.Size.Y)},
			// }
			// paint.ColorOp{Color: f.th.Colors.Get(f.col)}.Add(gtx.Ops)
			//
			// paint.PaintOp{}.Add(gtx.Ops)
			// gtx.Constraints.Constrain(d)
			gtx.Constraints.Constrain(dims.Size)
			dims = f.w(gtx)
			return dims
		},
	).Fn(gtx)
}

func fill(gtx l.Context, col color.NRGBA, bounds image.Point) {
	clip.UniformRRect(f32.Rectangle{
		Max: f32.Pt(float32(bounds.X), float32(bounds.Y)),
	}, 0).Add(gtx.Ops)
	paint.Fill(gtx.Ops, col)
}
