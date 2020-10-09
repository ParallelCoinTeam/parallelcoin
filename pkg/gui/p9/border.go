package p9

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// _border lays out a widget and draws a border inside it.
type _border struct {
	th           *Theme
	color        color.RGBA
	cornerRadius unit.Value
	width        unit.Value
}

func (th *Theme) Border() *_border {
	b := &_border{th: th}
	return b
}

func (b *_border) Color(color color.RGBA) *_border {
	b.color = color
	return b
}
func (b *_border) CornerRadius(rad float32) *_border {
	b.cornerRadius = b.th.textSize.Scale(rad)
	return b
}
func (b *_border) Width(width float32) *_border {
	b.width = b.th.textSize.Scale(width)
	return b
}

func (b *_border) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	dims := w(gtx)
	sz := dims.Size
	rr := float32(gtx.Px(b.cornerRadius))
	st := op.Push(gtx.Ops)
	width := gtx.Px(b.width)
	clip.Border{
		Rect: f32.Rectangle{
			Max: layout.FPt(sz),
		},
		NE: rr, NW: rr, SE: rr, SW: rr,
		Width: float32(width),
	}.Add(gtx.Ops)
	dr := f32.Rectangle{
		Max: layout.FPt(sz),
	}
	paint.ColorOp{Color: b.color}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	st.Pop()
	return dims
}
