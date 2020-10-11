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
	w            layout.Widget
}

// Border creates a border with configurable color, width and corner radius.
func (th *Theme) Border() *_border {
	b := &_border{
		th: th,
	}
	b.CornerRadius(0.25).Color("Primary").Width(0.125)
	return b
}

// Color sets the color to render the border in
func (b *_border) Color(color string) *_border {
	b.color = b.th.Colors.Get(color)
	return b
}

// CornerRadius sets the radius of the curve on the corners
func (b *_border) CornerRadius(rad float32) *_border {
	b.cornerRadius = b.th.TextSize.Scale(rad)
	return b
}

// Width sets the width of the border line
func (b *_border) Width(width float32) *_border {
	b.width = b.th.TextSize.Scale(width)
	return b
}

func (b *_border) Embed(w layout.Widget) *_border {
	b.w = w
	return b
}

// Fn renders the border
func (b *_border) Fn(gtx layout.Context) layout.Dimensions {
	dims := b.w(gtx)
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
