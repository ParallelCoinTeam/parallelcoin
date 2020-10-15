package p9

import (
	"image"
	"image/color"

	"gioui.org/f32"
	l "gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
)

type (
	// Defining these as types gives flexibility later to create methods that modify them
	Fonts      map[string]text.Typeface
	Icons      map[string]*Icon
	Collection []text.FontFace
)

const Inf = 1e6

// Fill is a general fill function that covers the background of the current context space
func Fill(gtx l.Context, col color.RGBA) l.Dimensions {
	cs := gtx.Constraints
	d := cs.Min
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Constraints.Constrain(d)
	return l.Dimensions{Size: d}
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

func axisPoint(a l.Axis, main, cross int) image.Point {
	if a == l.Horizontal {
		return image.Point{X: main, Y: cross}
	} else {
		return image.Point{X: cross, Y: main}
	}
}

func axisMain(a l.Axis, sz image.Point) int {
	if a == l.Horizontal {
		return sz.X
	} else {
		return sz.Y
	}
}


func axisCross(a l.Axis, sz image.Point) int {
	if a == l.Horizontal {
		return sz.Y
	} else {
		return sz.X
	}
}

func axisMainConstraint(a l.Axis, cs l.Constraints) (int, int) {
	if a == l.Horizontal {
		return cs.Min.X, cs.Max.X
	} else {
		return cs.Min.Y, cs.Max.Y
	}
}

func axisCrossConstraint(a l.Axis, cs l.Constraints) (int, int) {
	if a == l.Horizontal {
		return cs.Min.Y, cs.Max.Y
	} else {
		return cs.Min.X, cs.Max.X
	}
}

func axisConstraints(a l.Axis, mainMin, mainMax, crossMin, crossMax int) l.Constraints {
	if a == l.Horizontal {
		return l.Constraints{Min: image.Pt(mainMin, crossMin), Max: image.Pt(mainMax, crossMax)}
	} else {
		return l.Constraints{Min: image.Pt(crossMin, mainMin), Max: image.Pt(crossMax, mainMax)}
	}
}
