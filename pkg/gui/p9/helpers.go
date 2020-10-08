package p9

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
)

type (
	// Defining these as types gives flexibility later to create methods that modify them
	Fonts      map[string]text.Typeface
	Icons      map[string]*Ico
	Collection []text.FontFace
)

// Fill is a general fill function that covers the background of the current context space
func Fill(gtx layout.Context, col color.RGBA) layout.Dimensions {
	cs := gtx.Constraints
	d := cs.Min
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Constraints.Constrain(d)
	return layout.Dimensions{Size: d}
}
