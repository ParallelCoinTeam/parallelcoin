package gel

import (
	"image/color"
	
	"github.com/p9c/gio/f32"
	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/op/clip"
	"github.com/p9c/gio/unit"
	
	"github.com/p9c/gel/f32color"
)

type ButtonLayout struct {
	*Window
	background   color.NRGBA
	cornerRadius unit.Value
	button       *Clickable
	w            l.Widget
	corners      int
}

// ButtonLayout creates a button with a background and another widget over top
func (w *Window) ButtonLayout(button *Clickable) *ButtonLayout {
	return &ButtonLayout{
		Window:       w,
		button:       button,
		background:   w.Colors.GetNRGBAFromName("ButtonBg"),
		cornerRadius: w.TextSize.Scale(0.25),
	}
}

// Corners sets which corners have the radius of rounding
func (b *ButtonLayout) Corners(corners int) *ButtonLayout {
	b.corners = corners
	return b
}

// Background sets the background color of the button
func (b *ButtonLayout) Background(color string) *ButtonLayout {
	b.background = b.Theme.Colors.GetNRGBAFromName(color)
	return b
}

// CornerRadius sets the radius of the corners of the button
func (b *ButtonLayout) CornerRadius(radius float32) *ButtonLayout {
	b.cornerRadius = b.Theme.TextSize.Scale(radius)
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
	return b.Stack().Alignment(l.Center).
		Expanded(
			func(gtx l.Context) l.Dimensions {
				rr := float32(gtx.Px(b.cornerRadius))
				clip.RRect{
					Rect: f32.Rectangle{Max: f32.Point{
						X: float32(gtx.Constraints.Min.X),
						Y: float32(gtx.Constraints.Min.Y),
					}},
					NW: ifDir(rr, b.corners&NW),
					NE: ifDir(rr, b.corners&NE),
					SW: ifDir(rr, b.corners&SW),
					SE: ifDir(rr, b.corners&SE),
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
