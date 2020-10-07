package p9

import (
	"image"
	"image/color"

	"gioui.org/f32"
	l "gioui.org/layout"
	"gioui.org/op/paint"
)

type _filler struct {
	th  *Theme
	col string
	w   l.Widget
}

// Fill fills underneath a widget you can put over top of it
func (th *Theme) Fill(col string) *_filler {
	return &_filler{th: th, col: col}
}

func (f *_filler) Widget(w l.Widget) *_filler {
	f.w = w
	return f
}

func (f *_filler) Fn(gtx l.Context) l.Dimensions {
	return f.th.Stack().Stacked(f.w).Expanded(
		func(c l.Context) l.Dimensions {
			dims := f.w(gtx)
			cs := gtx.Constraints
			d := image.Point{X: cs.Max.X, Y: cs.Max.Y}
			dr := f32.Rectangle{
				Max: f32.Point{X: float32(dims.Size.X), Y: float32(dims.Size.Y)},
			}
			paint.ColorOp{Color: f.th.Colors.Get(f.col)}.Add(gtx.Ops)
			paint.PaintOp{Rect: dr}.Add(gtx.Ops)
			gtx.Constraints.Constrain(d)
			f.w(gtx)
			gtx.Constraints.Constrain(dims.Size)
			return dims
		},
	).Fn(gtx)
}

func bounds(gtx l.Context) f32.Rectangle {
	cs := gtx.Constraints
	d := cs.Max
	return f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
}

func (th *Theme) fill(gtx l.Context, col color.RGBA) {
	dr := bounds(gtx)
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
}

func fill(gtx *l.Context, col string) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Max.X, Y: cs.Max.Y}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: HexARGB(col)}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Constraints.Constrain(d)
}
