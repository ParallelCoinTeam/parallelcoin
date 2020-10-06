package p9

import (
	"image/color"

	"gioui.org/f32"
	l "gioui.org/layout"
	"gioui.org/op/paint"
)

type _filler struct {
	col color.RGBA
	w   l.Widget
}

// Fill fills underneath a widget you can put over top of it
func (th *Theme) Fill(col string) *_filler {
	return &_filler{col: th.Colors.Get(col)}
}

func (f *_filler) Widget(w l.Widget) *_filler {
	f.w = w
	return f
}

func (f *_filler) Fn(gtx l.Context) l.Dimensions {
	return l.Stack{Alignment: l.Center}.Layout(gtx,
		l.Expanded(func(gtx l.Context) l.Dimensions {
			cs := gtx.Constraints
			d := cs.Max
			dr := f32.Rectangle{
				Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
			}
			paint.ColorOp{Color: f.col}.Add(gtx.Ops)
			paint.PaintOp{Rect: dr}.Add(gtx.Ops)
			return l.Dimensions{Size: d}
		}),
		l.Expanded(f.w),
	)
}
