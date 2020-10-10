package p9

import (
	"image/color"

	"gioui.org/f32"
	l "gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/unit"

	"github.com/p9c/pod/pkg/gui/f32color"
)

type _buttonLayout struct {
	th           *Theme
	background   color.RGBA
	cornerRadius unit.Value
	button       *_clickable
	w            l.Widget
}

// ButtonLayout creates a button with a background and another widget over top
func (th *Theme) ButtonLayout(button *_clickable) *_buttonLayout {
	return &_buttonLayout{
		th:           th,
		button:       button,
		background:   th.Colors.Get("ButtonBg"),
		cornerRadius: th.textSize.Scale(0.125),
	}
}

// Background sets the background color of the button
func (b *_buttonLayout) Background(color string) *_buttonLayout {
	b.background = b.th.Colors.Get(color)
	return b
}

// CornerRadius sets the radius of the corners of the button
func (b *_buttonLayout) CornerRadius(radius float32) *_buttonLayout {
	b.cornerRadius = b.th.textSize.Scale(radius)
	return b
}

// Embed a widget in the button
func (b *_buttonLayout) Embed(w l.Widget) *_buttonLayout {
	b.w = w
	return b
}

// Fn is the function that draws the button and its child widget
func (b *_buttonLayout) Fn(gtx l.Context) l.Dimensions {
	min := gtx.Constraints.Min
	return l.Stack{Alignment: l.Center}.Layout(gtx,
		l.Expanded(func(gtx l.Context) l.Dimensions {
			rr := float32(gtx.Px(b.cornerRadius))
			clip.RRect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(gtx.Constraints.Min.X),
					Y: float32(gtx.Constraints.Min.Y),
				}},
				NE: rr, NW: rr, SE: rr, SW: rr,
			}.Add(gtx.Ops)
			background := b.background
			if gtx.Queue == nil {
				background = f32color.MulAlpha(b.background, 150)
			}
			dims := Fill(gtx, background)
			for _, c := range b.button.History() {
				drawInk(gtx, c)
			}
			return dims
		}),
		l.Stacked(func(gtx l.Context) l.Dimensions {
			gtx.Constraints.Min = min
			return l.Center.Layout(gtx, b.w)
		}),
		l.Expanded(b.button.Fn),
	)
}
