package p9

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	l "gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/p9c/pod/pkg/gui/f32color"
	w "github.com/p9c/pod/pkg/gui/widget"
)

type _iconButton struct {
	th         *Theme
	background color.RGBA
	// Color is the icon color.
	color color.RGBA
	icon  *Ico
	// Size is the icon size.
	size   unit.Value
	inset  *_inset
	button *w.Clickable
}

// IconButton creates an icon with a circular background and an icon placed in the centre
func (th *Theme) IconButton(button *w.Clickable) *_iconButton {
	return &_iconButton{
		th:         th,
		background: th.Colors.Get("Primary"),
		color:      th.Colors.Get("DocBg"),
		size:       th.textSize,
		inset:      th.Inset(0.5),
		button:     button,
	}
}

// Background sets the color of the circular background
func (b *_iconButton) Background(color string) *_iconButton {
	b.background = b.th.Colors.Get(color)
	return b
}

// Color sets the color of the icon
func (b *_iconButton) Color(color string) *_iconButton {
	b.color = b.th.Colors.Get(color)
	return b
}

// Icon sets the icon to display
func (b *_iconButton) Icon(ico *Ico) *_iconButton {
	b.icon = ico
	return b
}

// Scale changes the size of the icon as a ratio of the base font size
func (b *_iconButton) Scale(scale float32) *_iconButton {
	b.size = b.th.textSize.Scale(scale)
	return b
}

// Inset sets the size of inset that goes in between the button background and the icon
func (b *_iconButton) Inset(inset float32) *_iconButton {
	b.inset = b.th.Inset(inset).Widget(b.button.Fn)
	return b
}

// Fn renders the icon button
func (b *_iconButton) Fn(gtx l.Context) l.Dimensions {
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
				drawInk(gtx, widget.Press(c))
			}
			return dims
		},
	).Stacked(
		b.inset.Widget(b.icon.Fn).Fn,
	).Expanded(func(gtx l.Context) l.Dimensions {
		pointer.Ellipse(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
		return b.button.Fn(gtx)
	}).Fn(gtx)
}
