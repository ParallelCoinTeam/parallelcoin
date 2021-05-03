package gel

import (
	"image"
	"image/color"
	
	"github.com/p9c/gio/f32"
	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/op/clip"
	"github.com/p9c/gio/op/paint"
)

// Filler fills the background of a widget with a specified color and corner
// radius
type Filler struct {
	*Window
	col          string
	w            l.Widget
	dxn          l.Direction
	cornerRadius float32
	corners      int
}

const (
	NW = 1 << iota
	NE
	SW
	SE
)

// Fill fills underneath a widget you can put over top of it, dxn sets which
// direction to place a smaller object, cardinal axes and center
func (w *Window) Fill(col string, dxn l.Direction, radius float32, corners int, embed l.Widget) *Filler {
	return &Filler{Window: w, col: col, w: embed, dxn: dxn, cornerRadius: radius, corners: corners}
}

// Fn renders the fill with the widget inside
func (f *Filler) Fn(gtx l.Context) l.Dimensions {
	gtx1 := CopyContextDimensionsWithMaxAxis(gtx, l.Horizontal)
	// generate the dimensions for all the list elements
	dL := GetDimensionList(gtx1, 1, func(gtx l.Context, index int) l.Dimensions {
		return f.w(gtx)
	})
	fill(gtx, f.Colors.GetNRGBAFromName(f.col), dL[0].Size, f.cornerRadius, f.corners)
	return f.dxn.Layout(gtx, f.w)
}

func ifDir(radius float32, dir int) float32 {
	if dir != 0 {
		return radius
	}
	return 0
}

func fill(gtx l.Context, col color.NRGBA, bounds image.Point, radius float32, cnrs int) {
	rect := f32.Rectangle{
		Max: f32.Pt(float32(bounds.X), float32(bounds.Y)),
	}
	clip.RRect{
		Rect: rect,
		NW:   ifDir(radius, cnrs&NW),
		NE:   ifDir(radius, cnrs&NE),
		SW:   ifDir(radius, cnrs&SW),
		SE:   ifDir(radius, cnrs&SE),
	}.Add(gtx.Ops)
	paint.Fill(gtx.Ops, col)
}
