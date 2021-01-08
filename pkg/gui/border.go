package gui

import (
	"image/color"
	
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// Border lays out a widget and draws a border inside it.
type Border struct {
	*Window
	color        color.NRGBA
	cornerRadius unit.Value
	width        unit.Value
	w            layout.Widget
	corners      int
}

// Border creates a border with configurable color, width and corner radius.
func (w *Window) Border() *Border {
	b := &Border{Window: w}
	b.CornerRadius(0.25).Color("Primary").Width(0.125)
	return b
}

// Corners sets the corners that are rounded
func (b *Border) Corners(corners int) *Border {
	b.corners = corners
	return b
}

// Color sets the color to render the border in
func (b *Border) Color(color string) *Border {
	b.color = b.Theme.Colors.GetNRGBAFromName(color)
	return b
}

// CornerRadius sets the radius of the curve on the corners
func (b *Border) CornerRadius(rad float32) *Border {
	b.cornerRadius = b.Theme.TextSize.Scale(rad)
	return b
}

// Width sets the width of the border line
func (b *Border) Width(width float32) *Border {
	b.width = b.Theme.TextSize.Scale(width)
	return b
}

func (b *Border) Embed(w layout.Widget) *Border {
	b.w = w
	return b
}

// Fn renders the border
func (b *Border) Fn(gtx layout.Context) layout.Dimensions {
	dims := b.w(gtx)
	sz := dims.Size
	rr := float32(gtx.Px(b.cornerRadius))
	st := op.Push(gtx.Ops)
	width := gtx.Px(b.width)
	clip.Border{
		Rect:  f32.Rectangle{Max: layout.FPt(sz)},
		NE:    ifDir(rr, b.corners&NE),
		NW:    ifDir(rr, b.corners&NW),
		SE:    ifDir(rr, b.corners&SE),
		SW:    ifDir(rr, b.corners&SW),
		Width: float32(width),
	}.Add(gtx.Ops)
	paint.ColorOp{Color: b.color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	st.Pop()
	return dims
}
