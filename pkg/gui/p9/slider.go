package p9

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/p9c/pod/pkg/gui/f32color"
)

// Slider is for selecting a value in a range.
func (th *Theme) Slider(float *widget.Float, min, max float32) SliderStyle {
	return SliderStyle{
		Min:   min,
		Max:   max,
		Color: th.Color.Primary,
		Float: float,
	}
}

type SliderStyle struct {
	Min, Max float32
	Color    color.RGBA
	Float    *widget.Float
}

func (s SliderStyle) Fn(gtx layout.Context) layout.Dimensions {
	thumbRadiusInt := gtx.Px(unit.Dp(6))
	trackWidth := float32(gtx.Px(unit.Dp(2)))
	thumbRadius := float32(thumbRadiusInt)
	halfWidthInt := 2 * thumbRadiusInt
	halfWidth := float32(halfWidthInt)

	size := gtx.Constraints.Min
	// Keep a minimum length so that the track is always visible.
	minLength := halfWidthInt + 3*thumbRadiusInt + halfWidthInt
	if size.X < minLength {
		size.X = minLength
	}
	size.Y = 2 * halfWidthInt

	st := op.Push(gtx.Ops)
	op.Offset(f32.Pt(halfWidth, 0)).Add(gtx.Ops)
	gtx.Constraints.Min = image.Pt(size.X-2*halfWidthInt, size.Y)
	s.Float.Layout(gtx, halfWidthInt, s.Min, s.Max)
	thumbPos := halfWidth + s.Float.Pos()
	st.Pop()

	color := s.Color
	if gtx.Queue == nil {
		color = f32color.MulAlpha(color, 150)
	}

	// Draw track before thumb.
	st = op.Push(gtx.Ops)
	track := f32.Rectangle{
		Min: f32.Point{
			X: halfWidth,
			Y: halfWidth - trackWidth/2,
		},
		Max: f32.Point{
			X: thumbPos,
			Y: halfWidth + trackWidth/2,
		},
	}
	clip.RRect{Rect: track}.Add(gtx.Ops)
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{Rect: track}.Add(gtx.Ops)
	st.Pop()

	// Draw track after thumb.
	st = op.Push(gtx.Ops)
	track.Min.X = thumbPos
	track.Max.X = float32(size.X) - halfWidth
	clip.RRect{Rect: track}.Add(gtx.Ops)
	paint.ColorOp{Color: f32color.MulAlpha(color, 96)}.Add(gtx.Ops)
	paint.PaintOp{Rect: track}.Add(gtx.Ops)
	st.Pop()

	// Draw thumb.
	st = op.Push(gtx.Ops)
	thumb := f32.Rectangle{
		Min: f32.Point{
			X: thumbPos - thumbRadius,
			Y: halfWidth - thumbRadius,
		},
		Max: f32.Point{
			X: thumbPos + thumbRadius,
			Y: halfWidth + thumbRadius,
		},
	}
	rr := thumbRadius
	clip.RRect{
		Rect: thumb,
		NE:   rr, NW: rr, SE: rr, SW: rr,
	}.Add(gtx.Ops)
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{Rect: thumb}.Add(gtx.Ops)
	st.Pop()

	return layout.Dimensions{Size: size}
}
