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
	dxn l.Direction
}

// Fill fills underneath a widget you can put over top of it, dxn sets which
// direction to place a smaller object, cardinal axes and center
func (th *Theme) Fill(col string, w l.Widget, dxn l.Direction) *Filler {
	return &Filler{th: th, col: col, w: w, dxn: dxn}
}

func (f *Filler) Embed(w l.Widget) *Filler {
	f.w = w
	return f
}

func (f *Filler) Fn(gtx l.Context) l.Dimensions {
	gtx1 := CopyContextDimensionsWithMaxAxis(gtx, gtx.Constraints.Max, l.Horizontal)
	// generate the dimensions for all the list elements
	dL := GetDimensionList(gtx1, 1, func(gtx l.Context, index int) l.Dimensions {
		return f.w(gtx)
	})
	fill(gtx, f.th.Colors.Get(f.col), dL[0].Size)
	return f.w(gtx)
}

func fill(gtx l.Context, col color.NRGBA, bounds image.Point) {
	clip.UniformRRect(f32.Rectangle{
		Max: f32.Pt(float32(bounds.X), float32(bounds.Y)),
	}, 0).Add(gtx.Ops)
	paint.Fill(gtx.Ops, col)
}
