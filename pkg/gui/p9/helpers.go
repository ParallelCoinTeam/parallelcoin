package p9

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
)

type (
	// Defining these as types gives flexibility later to create methods that modify them
	Fonts      map[string]text.Typeface
	Icons      map[string]*Icon
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

func (th *Theme) GetFont(font string) *text.Font {
	for i := range th.collection {
		if th.collection[i].Font.Typeface == text.Typeface(font) {
			return &th.collection[i].Font
		}
	}
	return nil
}

func rgb(c uint32) color.RGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.RGBA {
	return color.RGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func toPointF(p image.Point) f32.Point {
	return f32.Point{X: float32(p.X), Y: float32(p.Y)}
}
