package gel

import (
	"image"
	"image/color"
	
	"github.com/p9c/gio/f32"
	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/op"
	"github.com/p9c/gio/op/clip"
	"github.com/p9c/gio/op/paint"
	"github.com/p9c/gio/unit"
	
	"github.com/p9c/gel/f32color"
)

type Slider struct {
	*Window
	min, max float32
	color    color.NRGBA
	float    *Float
}

// Slider is for selecting a value in a range.
func (w *Window) Slider() *Slider {
	return &Slider{
		Window: w,
		color:  w.Colors.GetNRGBAFromName("Primary"),
	}
}

// Min sets the value at the left hand side
func (s *Slider) Min(min float32) *Slider {
	s.min = min
	return s
}

// Max sets the value at the right hand side
func (s *Slider) Max(max float32) *Slider {
	s.max = max
	return s
}

// Color sets the color to draw the slider in
func (s *Slider) Color(color string) *Slider {
	s.color = s.Theme.Colors.GetNRGBAFromName(color)
	return s
}

// Float sets the initial value
func (s *Slider) Float(f *Float) *Slider {
	s.float = f
	return s
}

// Fn renders the slider
func (s *Slider) Fn(gtx l.Context) l.Dimensions {
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
	
	st := op.Save(gtx.Ops)
	op.Offset(f32.Pt(halfWidth, 0)).Add(gtx.Ops)
	gtx.Constraints.Min = image.Pt(size.X-2*halfWidthInt, size.Y)
	s.float.Fn(gtx, halfWidthInt, s.min, s.max)
	thumbPos := halfWidth + s.float.Pos()
	st.Load()
	
	col := s.color
	if gtx.Queue == nil {
		col = f32color.MulAlpha(col, 150)
	}
	
	// Draw track before thumb.
	st = op.Save(gtx.Ops)
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
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	st.Load()
	
	// Draw track after thumb.
	st = op.Save(gtx.Ops)
	track.Min.X = thumbPos
	track.Max.X = float32(size.X) - halfWidth
	clip.RRect{Rect: track}.Add(gtx.Ops)
	paint.ColorOp{Color: f32color.MulAlpha(col, 96)}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	st.Load()
	
	// Draw thumb.
	st = op.Save(gtx.Ops)
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
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	st.Load()
	
	return l.Dimensions{Size: size}
}
