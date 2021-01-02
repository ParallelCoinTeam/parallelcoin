package gui

import (
	"image"
	"image/color"
	
	"gioui.org/f32"
	l "gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type Filler struct {
	*Window
	col string
	w   l.Widget
	dxn l.Direction
}

// Fill fills underneath a widget you can put over top of it, dxn sets which
// direction to place a smaller object, cardinal axes and center
func (w *Window) Fill(col string, embed l.Widget, dxn l.Direction) *Filler {
	return &Filler{Window: w, col: col, w: embed, dxn: dxn}
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
	fill(gtx, f.Theme.Colors.GetNRGBAFromName(f.col), dL[0].Size)
	return f.dxn.Layout(gtx, f.w)
}

func fill(gtx l.Context, col color.NRGBA, bounds image.Point) {
	clip.UniformRRect(f32.Rectangle{
		Max: f32.Pt(float32(bounds.X), float32(bounds.Y)),
	}, 0).Add(gtx.Ops)
	paint.Fill(gtx.Ops, col)
}
