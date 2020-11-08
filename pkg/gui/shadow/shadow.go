package shadow

import (
	"gioui.org/f32"
	l "gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/egonelbre/expgio/surface/f32color"
	"github.com/gioapp/gel/helper"
	"github.com/p9c/pod/pkg/gui/p9"
	"image/color"
)

type shadow struct {
	theme        *p9.Theme
	cornerRadius unit.Value
	elevation    unit.Value
	color        color.RGBA
}

func Shadow(th *p9.Theme) *shadow {
	return &shadow{
		theme:        th,
		cornerRadius: unit.Dp(5),
		elevation:    unit.Dp(5),
		color:        helper.HexARGB(th.Colors["PanelBg"]),
	}
}

func (s *shadow) Drop(w l.Dimensions) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		sz := gtx.Constraints.Min
		rr := float32(gtx.Px(s.cornerRadius))
		r := f32.Rect(0, 0, float32(sz.X), float32(sz.Y))
		s.layout(gtx, r, rr)
		clip.UniformRRect(r, rr).Add(gtx.Ops)
		paint.Fill(gtx.Ops, s.color)
		return w
	}
}

func (s *shadow) layout(gtx l.Context, r f32.Rectangle, rr float32) {
	if s.elevation.V <= 0 {
		return
	}
	offset := pxf(gtx.Metric, s.elevation)
	d := int(offset + 1)
	if d > 4 {
		d = 4
	}
	a := float32(s.color.A) / 0xff
	background := (f32color.RGBA{A: a * 0.4 / float32(d*d)}).SRGB()
	for x := 0; x <= d; x++ {
		for y := 0; y <= d; y++ {
			px, py := float32(x)/float32(d)-0.5, float32(y)/float32(d)-0.15
			stack := op.Push(gtx.Ops)
			op.Offset(f32.Pt(px*offset, py*offset)).Add(gtx.Ops)
			clip.UniformRRect(r, rr).Add(gtx.Ops)
			paint.Fill(gtx.Ops, background)
			stack.Pop()
		}
	}
}

func outset(r f32.Rectangle, y, s float32) f32.Rectangle {
	r.Min.X += s
	r.Min.Y += s + y
	r.Max.X += -s
	r.Max.Y += -s + y
	return r
}

func pxf(c unit.Metric, v unit.Value) float32 {
	switch v.U {
	case unit.UnitPx:
		return v.V
	case unit.UnitDp:
		s := c.PxPerDp
		if s == 0 {
			s = 1
		}
		return s * v.V
	case unit.UnitSp:
		s := c.PxPerSp
		if s == 0 {
			s = 1
		}
		return s * v.V
	default:
		panic("unknown unit")
	}
}
