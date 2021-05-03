package gel

import (
	"errors"
	"image"
	"image/color"
	"time"
	
	"github.com/p9c/gio/f32"
	"github.com/p9c/gio/io/system"
	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/op"
	"github.com/p9c/gio/text"
)

// Defining these as types gives flexibility later to create methods that modify them

type Fonts map[string]text.Typeface
type Icons map[string]*Icon
type Collection []text.FontFace

func (c Collection) Font(font string) (out text.Font, e error) {
	for i := range c {
		if c[i].Font.Typeface == text.Typeface(font) {
			out = c[i].Font
			return
		}
	}
	return out, errors.New("font " + font + " not found")
	
}

const Inf = 1e6

// Fill is a general fill function that covers the background of the current context space
func Fill(gtx l.Context, col color.NRGBA) l.Dimensions {
	cs := gtx.Constraints
	d := cs.Min
	// dr := f32.Rectangle{
	// 	Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	// }
	// paint.ColorOp{Color: col}.Add(gtx.Ops)
	// paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	fill(gtx, col, d, 0, 0)
	return l.Dimensions{Size: d}
}

// func (th *Theme) GetFont(font string) *text.Font {
// 	for i := range th.collection {
// 		if th.collection[i].Font.Typeface == text.Typeface(font) {
// 			return &th.collection[i].Font
// 		}
// 	}
// 	return nil
// }

// func rgb(c uint32) color.RGBA {
// 	return argb(0xff000000 | c)
// }

func argb(c uint32) color.RGBA {
	return color.RGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

// FPt converts an point to a f32.Point.
func Fpt(p image.Point) f32.Point {
	return f32.Point{
		X: float32(p.X), Y: float32(p.Y),
	}
}

func axisPoint(a l.Axis, main, cross int) image.Point {
	if a == l.Horizontal {
		return image.Point{X: main, Y: cross}
	} else {
		return image.Point{X: cross, Y: main}
	}
}

// axisConvert a point in (x, y) coordinates to (main, cross) coordinates,
// or vice versa. Specifically, Convert((x, y)) returns (x, y) unchanged
// for the horizontal axis, or (y, x) for the vertical axis.
func axisConvert(a l.Axis, pt image.Point) image.Point {
	if a == l.Horizontal {
		return pt
	}
	return image.Pt(pt.Y, pt.X)
}

func axisMain(a l.Axis, sz image.Point) int {
	if a == l.Horizontal {
		return sz.X
	}
		return sz.Y
}

func axisCross(a l.Axis, sz image.Point) int {
	if a == l.Horizontal {
		return sz.Y
	}
		return sz.X
}

func axisMainConstraint(a l.Axis, cs l.Constraints) (int, int) {
	if a == l.Horizontal {
		return cs.Min.X, cs.Max.X
	}
	return cs.Min.Y, cs.Max.Y
}

func axisCrossConstraint(a l.Axis, cs l.Constraints) (int, int) {
	if a == l.Horizontal {
		return cs.Min.Y, cs.Max.Y
	}
	return cs.Min.X, cs.Max.X
}

func axisConstraints(a l.Axis, mainMin, mainMax, crossMin, crossMax int) l.Constraints {
	if a == l.Horizontal {
		return l.Constraints{Min: image.Pt(mainMin, crossMin), Max: image.Pt(mainMax, crossMax)}
	}
	return l.Constraints{Min: image.Pt(crossMin, mainMin), Max: image.Pt(crossMax, mainMax)}
}

func EmptySpace(x, y int) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		return l.Dimensions{
			Size: image.Point{
				X: x,
				Y: y,
			},
		}
	}
}

func EmptyFromSize(size image.Point) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		return l.Dimensions{
			Size: size,
		}
	}
}

func EmptyMaxWidth() func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		return l.Dimensions{
			Size:     image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Min.Y},
			Baseline: 0,
		}
	}
}
func EmptyMaxHeight() func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		return l.Dimensions{Size: image.Point{X: gtx.Constraints.Min.X, Y: gtx.Constraints.Min.Y}}
	}
}

func EmptyMinWidth() func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		return l.Dimensions{
			Size:     image.Point{X: gtx.Constraints.Min.X, Y: gtx.Constraints.Min.Y},
			Baseline: 0,
		}
	}
}
func EmptyMinHeight() func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		return l.Dimensions{Size: image.Point{Y: gtx.Constraints.Min.Y}}
	}
}

// CopyContextDimensionsWithMaxAxis copies the dimensions out with the max set by an image.Point along the axis
func CopyContextDimensionsWithMaxAxis(gtx l.Context, axis l.Axis) l.Context {
	ip := image.Point{}
	if axis == l.Horizontal {
		ip.Y = gtx.Constraints.Max.Y
		ip.X = gtx.Constraints.Max.X
	} else {
		ip.Y = gtx.Constraints.Max.Y
		ip.X = gtx.Constraints.Max.X
	}
	var ops op.Ops
	gtx1 := l.NewContext(
		&ops, system.FrameEvent{
			Now:    time.Now(),
			Metric: gtx.Metric,
			Size:   ip,
		},
	)
	if axis == l.Horizontal {
		gtx1.Constraints.Min.X = 0
		gtx1.Constraints.Min.Y = gtx.Constraints.Min.Y
	} else {
		gtx1.Constraints.Min.X = gtx.Constraints.Min.X
		gtx1.Constraints.Min.Y = 0
	}
	return gtx1
}

// GetInfContext creates a context with infinite max constraints
func GetInfContext(gtx l.Context) l.Context {
	ip := image.Point{}
	ip.Y = Inf
	ip.X = Inf
	var ops op.Ops
	gtx1 := l.NewContext(
		&ops, system.FrameEvent{
			Now:    time.Now(),
			Metric: gtx.Metric,
			Size:   ip,
		},
	)
	return gtx1
}

func If(value bool, t, f l.Widget) l.Widget {
	if value {
		return t
	} else {
		return f
	}
}

func (th *Theme) SliceToWidget(w []l.Widget, axis l.Axis) l.Widget {
	var out *Flex
	if axis == l.Horizontal {
		out = th.Flex().AlignStart()
	} else {
		out = th.VFlex().AlignStart()
	}
	for i := range w {
		out.Rigid(w[i])
	}
	return out.Fn
}
