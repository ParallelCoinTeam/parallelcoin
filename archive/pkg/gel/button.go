package gel

import (
	"image/color"
	"math"
	"strings"
	
	"github.com/p9c/gio/f32"
	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/op"
	"github.com/p9c/gio/op/clip"
	"github.com/p9c/gio/op/paint"
	"github.com/p9c/gio/text"
	"github.com/p9c/gio/unit"
	
	"github.com/p9c/gel/f32color"
)

// Button is a material text label icon with options to change all features
type Button struct {
	*Window
	background   color.NRGBA
	color        color.NRGBA
	cornerRadius unit.Value
	font         text.Font
	inset        *l.Inset
	text         string
	textSize     unit.Value
	button       *Clickable
	shaper       text.Shaper
}

// Button is a regular material text button where all the dimensions, colors, corners and font can be changed
func (w *Window) Button(btn *Clickable) *Button {
	var font text.Font
	var e error
	if font, e = w.collection.Font("plan9"); E.Chk(e) {
	}
	return &Button{
		Window: w,
		text:   strings.ToUpper("text unset"),
		// default sets
		font:         font,
		color:        w.Colors.GetNRGBAFromName("DocBg"),
		cornerRadius: w.TextSize.Scale(0.25),
		background:   w.Colors.GetNRGBAFromName("Primary"),
		textSize:     w.TextSize,
		inset: &l.Inset{
			Top:    w.TextSize.Scale(0.5),
			Bottom: w.TextSize.Scale(0.5),
			Left:   w.TextSize.Scale(0.5),
			Right:  w.TextSize.Scale(0.5),
		},
		button: btn,
		shaper: w.shaper,
	}
}

// Background sets the background color
func (b *Button) Background(background string) *Button {
	b.background = b.Theme.Colors.GetNRGBAFromName(background)
	return b
}

// Color sets the text color
func (b *Button) Color(color string) *Button {
	b.color = b.Theme.Colors.GetNRGBAFromName(color)
	return b
}

// CornerRadius sets the corner radius (all measurements are scaled from the base text size)
func (b *Button) CornerRadius(cornerRadius float32) *Button {
	b.cornerRadius = b.TextSize.Scale(cornerRadius)
	return b
}

// Font sets the font style
func (b *Button) Font(font string) *Button {
	var fon text.Font
	var e error
	if fon, e = b.collection.Font(font); !E.Chk(e) {
		b.font = fon
	} else {
		panic(e)
	}
	return b
}

// Inset sets the inset between the button border and the text
func (b *Button) Inset(scale float32) *Button {
	b.inset = &l.Inset{
		Top:    b.TextSize.Scale(scale),
		Right:  b.TextSize.Scale(scale),
		Bottom: b.TextSize.Scale(scale),
		Left:   b.TextSize.Scale(scale),
	}
	return b
}

// Text sets the text on the button
func (b *Button) Text(text string) *Button {
	b.text = text
	return b
}

// TextScale sets the dimensions of the text as a fraction of the base text size
func (b *Button) TextScale(scale float32) *Button {
	b.textSize = b.Theme.TextSize.Scale(scale)
	return b
}

// SetClick defines the callback to run on a click (mouse up) event
func (b *Button) SetClick(fn func()) *Button {
	b.button.SetClick(fn)
	return b
}

// SetCancel sets the callback to run when the user presses down over the button
// but then moves out of its hitbox before release (click)
func (b *Button) SetCancel(fn func()) *Button {
	b.button.SetCancel(fn)
	return b
}

func (b *Button) SetPress(fn func()) *Button {
	b.button.SetPress(fn)
	return b
}

// Fn renders the button
func (b *Button) Fn(gtx l.Context) l.Dimensions {
	bl := &ButtonLayout{
		background:   b.background,
		cornerRadius: b.cornerRadius,
		button:       b.button,
	}
	fn := func(gtx l.Context) l.Dimensions {
		return b.inset.Layout(
			gtx, func(gtx l.Context) l.Dimensions {
				// paint.ColorOp{Color: b.color}.Add(gtx.Ops)
				return b.Flex().Rigid(b.Label().Text(b.text).TextScale(b.textSize.V / b.TextSize.V).Fn).Fn(gtx)
				// b.Window.Text().
				// Alignment(text.Middle).
				// Fn(gtx, b.shaper, b.font, b.textSize, b.text)
			},
		)
	}
	bl.Embed(fn)
	return bl.Fn(gtx)
}

func drawInk(c l.Context, p press) {
	// duration is the number of seconds for the completed animation: expand while
	// fading in, then out.
	const (
		expandDuration = float32(0.5)
		fadeDuration   = float32(0.9)
	)
	now := c.Now
	t := float32(now.Sub(p.Start).Seconds())
	end := p.End
	if end.IsZero() {
		// If the press hasn't ended, don't fade-out.
		end = now
	}
	endt := float32(end.Sub(p.Start).Seconds())
	// Compute the fade-in/out position in [0;1].
	var alphat float32
	{
		var haste float32
		if p.Cancelled {
			// If the press was cancelled before the inkwell was fully faded in, fast
			// forward the animation to match the fade-out.
			if h := 0.5 - endt/fadeDuration; h > 0 {
				haste = h
			}
		}
		// Fade in.
		half1 := t/fadeDuration + haste
		if half1 > 0.5 {
			half1 = 0.5
		}
		// Fade out.
		half2 := float32(now.Sub(end).Seconds())
		half2 /= fadeDuration
		half2 += haste
		if half2 > 0.5 {
			// Too old.
			return
		}
		
		alphat = half1 + half2
	}
	// Compute the expand position in [0;1].
	sizet := t
	if p.Cancelled {
		// Freeze expansion of cancelled presses.
		sizet = endt
	}
	sizet /= expandDuration
	// Animate only ended presses, and presses that are fading in.
	if !p.End.IsZero() || sizet <= 1.0 {
		op.InvalidateOp{}.Add(c.Ops)
	}
	if sizet > 1.0 {
		sizet = 1.0
	}
	if alphat > .5 {
		// Start fadeout after half the animation.
		alphat = 1.0 - alphat
	}
	// Twice the speed to attain fully faded in at 0.5.
	t2 := alphat * 2
	// Beziér ease-in curve.
	alphaBezier := t2 * t2 * (3.0 - 2.0*t2)
	sizeBezier := sizet * sizet * (3.0 - 2.0*sizet)
	size := float32(c.Constraints.Min.X)
	if h := float32(c.Constraints.Min.Y); h > size {
		size = h
	}
	// Cover the entire constraints min rectangle.
	size *= 2 * float32(math.Sqrt(2))
	// Apply curve values to size and color.
	size *= sizeBezier
	alpha := 0.7 * alphaBezier
	const col = 0.8
	ba, bc := byte(alpha*0xff), byte(col*0xff)
	defer op.Save(c.Ops).Load()
	rgba := f32color.MulAlpha(color.NRGBA{A: 0xff, R: bc, G: bc, B: bc}, ba)
	ink := paint.ColorOp{Color: rgba}
	ink.Add(c.Ops)
	rr := size * .5
	op.Offset(
		p.Position.Add(
			f32.Point{
				X: -rr,
				Y: -rr,
			},
		),
	).Add(c.Ops)
	clip.RRect{
		Rect: f32.Rectangle{
			Max: f32.Point{
				X: size,
				Y: size,
			},
		},
		NE: rr, NW: rr, SE: rr, SW: rr,
	}.Add(c.Ops)
	paint.PaintOp{}.Add(c.Ops)
}
