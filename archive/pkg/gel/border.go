package gel

import (
	"image/color"
	
	"github.com/p9c/gio/f32"
	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/op/clip"
	"github.com/p9c/gio/op/paint"
	"github.com/p9c/gio/unit"
)

// Border lays out a widget and draws a border inside it.
type Border struct {
	*Window
	color        color.NRGBA
	cornerRadius unit.Value
	width        unit.Value
	w            l.Widget
}

// Border creates a border with configurable color, width and corner radius.
func (w *Window) Border() *Border {
	b := &Border{Window: w}
	b.CornerRadius(0.25).Color("Primary").Width(0.125)
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

func (b *Border) Embed(w l.Widget) *Border {
	b.w = w
	return b
}

// Fn renders the border
func (b *Border) Fn(gtx l.Context) l.Dimensions {
	dims := b.w(gtx)
	sz := l.FPt(dims.Size)
	
	rr := float32(gtx.Px(b.cornerRadius))
	width := float32(gtx.Px(b.width))
	sz.X -= width
	sz.Y -= width
	
	r := f32.Rectangle{Max: sz}
	r = r.Add(f32.Point{X: width * 0.5, Y: width * 0.5})
	
	paint.FillShape(gtx.Ops,
		b.color,
		clip.Stroke{
			Path:  clip.UniformRRect(r, rr).Path(gtx.Ops),
			Style: clip.StrokeStyle{Width: width},
		}.Op(),
	)
	
	return dims
}
