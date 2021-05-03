package gel

import (
	"image"
	"image/color"
	
	"github.com/p9c/gio/f32"
	"github.com/p9c/gio/io/pointer"
	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/op"
	"github.com/p9c/gio/op/clip"
	"github.com/p9c/gio/op/paint"
	"github.com/p9c/gio/unit"
	
	"github.com/p9c/gel/f32color"
)

type Switch struct {
	*Window
	color struct {
		enabled  color.NRGBA
		disabled color.NRGBA
	}
	swtch *Bool
}

// Switch creates a boolean switch widget (basically a checkbox but looks like a switch)
func (w *Window) Switch(swtch *Bool) *Switch {
	sw := &Switch{
		Window: w,
		swtch:  swtch,
	}
	sw.color.enabled = w.Colors.GetNRGBAFromName("Primary")
	sw.color.disabled = w.Colors.GetNRGBAFromName("scrim")
	return sw
}

// EnabledColor sets the color to draw for the enabled state
func (s *Switch) EnabledColor(color string) *Switch {
	s.color.enabled = s.Theme.Colors.GetNRGBAFromName(color)
	return s
}

// DisabledColor sets the color to draw for the disabled state
func (s *Switch) DisabledColor(color string) *Switch {
	s.color.disabled = s.Theme.Colors.GetNRGBAFromName(color)
	return s
}

func (s *Switch) SetHook(fn func(b bool)) *Switch {
	s.swtch.SetOnChange(fn)
	return s
}

// Fn updates the switch and displays it.
func (s *Switch) Fn(gtx l.Context) l.Dimensions {
	return s.Inset(0.25, func(gtx l.Context) l.Dimensions {
		trackWidth := gtx.Px(s.Theme.TextSize.Scale(2.5))
		trackHeight := gtx.Px(unit.Dp(16))
		thumbSize := gtx.Px(unit.Dp(20))
		trackOff := float32(thumbSize-trackHeight) * .5
		
		// Draw track.
		stack := op.Save(gtx.Ops)
		trackCorner := float32(trackHeight) / 2
		trackRect := f32.Rectangle{Max: f32.Point{
			X: float32(trackWidth),
			Y: float32(trackHeight),
		}}
		col := s.color.disabled
		if s.swtch.value {
			col = s.color.enabled
		}
		if gtx.Queue == nil {
			col = f32color.MulAlpha(col, 150)
		}
		trackColor := f32color.MulAlpha(col, 200)
		op.Offset(f32.Point{Y: trackOff}).Add(gtx.Ops)
		clip.RRect{
			Rect: trackRect,
			NE:   trackCorner, NW: trackCorner, SE: trackCorner, SW: trackCorner,
		}.Add(gtx.Ops)
		paint.ColorOp{Color: trackColor}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		stack.Load()
		
		// Draw thumb ink.
		stack = op.Save(gtx.Ops)
		inkSize := gtx.Px(unit.Dp(44))
		rr := float32(inkSize) * .5
		inkOff := f32.Point{
			X: float32(trackWidth)*.5 - rr,
			Y: -rr + float32(trackHeight)*.5 + trackOff,
		}
		op.Offset(inkOff).Add(gtx.Ops)
		gtx.Constraints.Min = image.Pt(inkSize, inkSize)
		clip.RRect{
			Rect: f32.Rectangle{
				Max: l.FPt(gtx.Constraints.Min),
			},
			NE: rr, NW: rr, SE: rr, SW: rr,
		}.Add(gtx.Ops)
		for _, p := range s.swtch.History() {
			drawInk(gtx, p)
		}
		stack.Load()
		
		// Compute thumb offset and color.
		stack = op.Save(gtx.Ops)
		if s.swtch.value {
			off := trackWidth - thumbSize
			op.Offset(f32.Point{X: float32(off)}).Add(gtx.Ops)
		}
		
		// Draw thumb shadow, a translucent disc slightly larger than the
		// thumb itself.
		shadowStack := op.Save(gtx.Ops)
		shadowSize := float32(2)
		// Center shadow horizontally and slightly adjust its Y.
		op.Offset(f32.Point{X: -shadowSize / 2, Y: -.75}).Add(gtx.Ops)
		drawDisc(gtx.Ops, float32(thumbSize)+shadowSize, color.NRGBA(argb(0x55000000)))
		shadowStack.Load()
		
		// Draw thumb.
		drawDisc(gtx.Ops, float32(thumbSize), col)
		stack.Load()
		
		// Set up click area.
		stack = op.Save(gtx.Ops)
		clickSize := gtx.Px(unit.Dp(40))
		clickOff := f32.Point{
			X: (float32(trackWidth) - float32(clickSize)) * .5,
			Y: (float32(trackHeight)-float32(clickSize))*.5 + trackOff,
		}
		op.Offset(clickOff).Add(gtx.Ops)
		sz := image.Pt(clickSize, clickSize)
		pointer.Ellipse(image.Rectangle{Max: sz}).Add(gtx.Ops)
		gtx.Constraints.Min = sz
		s.swtch.Fn(gtx)
		stack.Load()
		
		dims := image.Point{X: trackWidth, Y: thumbSize}
		return l.Dimensions{Size: dims}
	}).Fn(gtx)
}

func drawDisc(ops *op.Ops, sz float32, col color.NRGBA) {
	defer op.Save(ops).Load()
	rr := sz / 2
	r := f32.Rectangle{Max: f32.Point{X: sz, Y: sz}}
	clip.RRect{
		Rect: r,
		NE:   rr, NW: rr, SE: rr, SW: rr,
	}.Add(ops)
	paint.ColorOp{Color: col}.Add(ops)
	paint.PaintOp{}.Add(ops)
}
