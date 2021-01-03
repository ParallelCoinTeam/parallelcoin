package shadow

import (
	"image/color"
	
	"gioui.org/f32"
	l "gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	
	"github.com/p9c/pod/pkg/gui/f32color"
)

func Shadow(gtx l.Context, cornerRadius, elevation unit.Value, shadowColor color.NRGBA, content func(gtx l.Context) l.Dimensions) l.Dimensions {
	sz := content(gtx).Size
	rr := float32(gtx.Px(cornerRadius))
	r := f32.Rect(0, 0, float32(sz.X), float32(sz.Y))
	layoutShadow(gtx, r, elevation, rr, shadowColor)
	clip.UniformRRect(r, rr).Add(gtx.Ops)
	return content(gtx)
}

// TODO: Shadow directions
func layoutShadow(gtx l.Context, r f32.Rectangle, elevation unit.Value, rr float32, shadowColor color.NRGBA) {
	if elevation.V <= 0 {
		return
	}
	offset := pxf(gtx.Metric, elevation)
	d := int(offset + 1)
	if d > 4 {
		d = 4
	}
	a := float32(shadowColor.A) / 0xff
	background := (f32color.RGBA{A: a * 0.4 / float32(d*d)}).SRGB()
	for x := 0; x <= d; x++ {
		for y := 0; y <= d; y++ {
			px, py := float32(x)/float32(d)-0.5, float32(y)/float32(d)-0.15
			stack := op.Push(gtx.Ops)
			op.Offset(f32.Pt(px*offset, py*offset)).Add(gtx.Ops)
			clip.UniformRRect(r, rr).Add(gtx.Ops)
			paint.Fill(gtx.Ops, color.NRGBA(background))
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
