package p9

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	l "gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/unit"

	"github.com/p9c/pod/pkg/gui/f32color"
)

type IconButton struct {
	th         *Theme
	background color.RGBA
	// Color is the icon color.
	color color.RGBA
	icon  *Icon
	// Size is the icon size.
	size   unit.Value
	inset  *Inset
	button *Clickable
}

// IconButton creates an icon with a circular background and an icon placed in the centre
func (th *Theme) IconButton(button *Clickable) *IconButton {
	return &IconButton{
		th:         th,
		background: th.Colors.Get("Primary"),
		color:      th.Colors.Get("DocBg"),
		size:       th.TextSize,
		inset:      th.Inset(0.33),
		button:     button,
		icon:       th.Icon(),
	}
}

// Background sets the color of the circular background
func (b *IconButton) Background(color string) *IconButton {
	b.background = b.th.Colors.Get(color)
	return b
}

// Color sets the color of the icon
func (b *IconButton) Color(color string) *IconButton {
	b.color = b.th.Colors.Get(color)
	return b
}

// Icon sets the icon to display
func (b *IconButton) Icon(data []byte) *IconButton {
	b.icon.color = b.color
	b.icon.Size(b.size)
	b.icon.Src(data)
	return b
}

// Scale changes the size of the icon as a ratio of the base font size
func (b *IconButton) Scale(scale float32) *IconButton {
	b.size = b.th.TextSize.Scale(scale*0.72)
	return b
}

// Inset sets the size of inset that goes in between the button background and the icon
func (b *IconButton) Inset(inset float32) *IconButton {
	b.inset = b.th.Inset(inset).Embed(b.button.Fn)
	return b
}

func (b *IconButton) SetClick(fn func()) *IconButton {
	b.button.SetClick(fn)
	return b
}

func (b *IconButton) SetPress(fn func()) *IconButton {
	b.button.SetPress(fn)
	return b
}

func (b *IconButton) SetCancel(fn func()) *IconButton {
	b.button.SetCancel(fn)
	return b
}

// Fn renders the icon button
func (b *IconButton) Fn(gtx l.Context) l.Dimensions {
	return b.th.Stack().Expanded(
		func(gtx l.Context) l.Dimensions {
			sizex, sizey := gtx.Constraints.Min.X, gtx.Constraints.Min.Y
			sizexf, sizeyf := float32(sizex), float32(sizey)
			rr := (sizexf + sizeyf) * .25
			clip.RRect{
				Rect: f32.Rectangle{Max: f32.Point{X: sizexf, Y: sizeyf}},
				NE:   rr, NW: rr, SE: rr, SW: rr,
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
		},
	).Stacked(
		b.inset.Embed(b.icon.Fn).Fn,
	).Expanded(func(gtx l.Context) l.Dimensions {
		pointer.Ellipse(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
		return b.button.Fn(gtx)
	}).Fn(gtx)
}
