package gui

import (
	"image/color"

	"gioui.org/f32"
	l "gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/unit"
	
	"github.com/p9c/pod/pkg/gui/f32color"
)

type ButtonLayout struct {
	th           *Theme
	background   color.NRGBA
	cornerRadius unit.Value
	button       *Clickable
	w            l.Widget
}

// ButtonLayout creates a button with a background and another widget over top
func (th *Theme) ButtonLayout(button *Clickable) *ButtonLayout {
	return &ButtonLayout{
		th:           th,
		button:       button,
		background:   th.Colors.Get("ButtonBg"),
		cornerRadius: th.TextSize.Scale(0.125),
	}
}

// Background sets the background color of the button
func (b *ButtonLayout) Background(color string) *ButtonLayout {
	b.background = b.th.Colors.Get(color)
	return b
}

// CornerRadius sets the radius of the corners of the button
func (b *ButtonLayout) CornerRadius(radius float32) *ButtonLayout {
	b.cornerRadius = b.th.TextSize.Scale(radius)
	return b
}

// Embed a widget in the button
func (b *ButtonLayout) Embed(w l.Widget) *ButtonLayout {
	b.w = w
	return b
}

func (b *ButtonLayout) SetClick(fn func()) *ButtonLayout {
	b.button.SetClick(fn)
	return b
}

func (b *ButtonLayout) SetCancel(fn func()) *ButtonLayout {
	b.button.SetCancel(fn)
	return b
}

func (b *ButtonLayout) SetPress(fn func()) *ButtonLayout {
	b.button.SetPress(fn)
	return b
}

// Fn is the function that draws the button and its child widget
func (b *ButtonLayout) Fn(gtx l.Context) l.Dimensions {
	min := gtx.Constraints.Min
	return b.th.Stack().Alignment(l.Center).
		Expanded(
			func(gtx l.Context) l.Dimensions {
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
			}).
		Stacked(
			func(gtx l.Context) l.Dimensions {
				gtx.Constraints.Min = min
				return l.Center.Layout(gtx, b.w)
			}).
		Expanded(b.button.Fn).
		Fn(gtx)
}
