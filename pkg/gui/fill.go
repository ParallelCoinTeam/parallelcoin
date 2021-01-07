package gui

import (
	"image"
	"image/color"
	
	"gioui.org/f32"
	l "gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// Filler fills the background of a widget with a specified color and corner
// radius
type Filler struct {
	*Window
	col            string
	w              l.Widget
	dxn            l.Direction
	cornerRadius   float32
	partialRounded l.Direction
}

// Fill fills underneath a widget you can put over top of it, dxn sets which
// direction to place a smaller object, cardinal axes and center
func (w *Window) Fill(col string, dxn l.Direction, radius float32, partialRounded l.Direction, embed l.Widget) *Filler {
	return &Filler{Window: w, col: col, w: embed, dxn: dxn, cornerRadius: radius, partialRounded: partialRounded}
}

// Fn renders the fill with the widget inside
func (f *Filler) Fn(gtx l.Context) l.Dimensions {
	gtx1 := CopyContextDimensionsWithMaxAxis(gtx, gtx.Constraints.Max, l.Horizontal)
	// generate the dimensions for all the list elements
	dL := GetDimensionList(gtx1, 1, func(gtx l.Context, index int) l.Dimensions {
		return f.w(gtx)
	})
	fill(gtx, f.Theme.Colors.GetNRGBAFromName(f.col), dL[0].Size, f.cornerRadius, f.partialRounded)
	return f.dxn.Layout(gtx, f.w)
}

func fill(gtx l.Context, col color.NRGBA, bounds image.Point, radius float32, partialRounded l.Direction) {
	rect := f32.Rectangle{
		Max: f32.Pt(float32(bounds.X), float32(bounds.Y)),
	}
	var dSE, dSW, dNE, dNW float32
	switch partialRounded {
	case l.N:
		dSE = 0
		dSW = 0
		dNE = radius
		dNW = radius
	case l.NE:
		dSE = radius
		dSW = 0
		dNE = radius
		dNW = radius
	case l.E:
		dSE = radius
		dSW = 0
		dNE = radius
		dNW = 0
	case l.SE:
		dSE = radius
		dSW = radius
		dNE = radius
		dNW = 0
	case l.S:
		dSE = radius
		dSW = radius
		dNE = 0
		dNW = 0
	case l.SW:
		dSE = radius
		dSW = radius
		dNE = 0
		dNW = radius
	case l.W:
		dSE = 0
		dSW = radius
		dNE = 0
		dNW = radius
	case l.NW:
		dSE = 0
		dSW = radius
		dNE = radius
		dNW = radius
	case l.Center:
		dSE = radius
		dSW = radius
		dNE = radius
		dNW = radius
	default:
		dSE = 0
		dSW = 0
		dNE = 0
		dNW = 0
	}
	
	clip.RRect{
		Rect: rect,
		SE:   dSE,
		SW:   dSW,
		NE:   dNE,
		NW:   dNW,
	}.Add(gtx.Ops)
	paint.Fill(gtx.Ops, col)
}
