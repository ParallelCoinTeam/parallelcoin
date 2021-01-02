package gui

import (
	"image"
	
	"gioui.org/f32"
	"gioui.org/io/pointer"
	l "gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"golang.org/x/exp/shiny/materialdesign/icons"
	
	"github.com/p9c/pod/pkg/gui/f32color"
)

type IconButton struct {
	*Window
	background string
	// Color is the icon color.
	color string
	icon  *Icon
	// Size is the icon size.
	size   unit.Value
	inset  *Inset
	button *Clickable
}

// IconButton creates an icon with a circular background and an icon placed in the centre
func (w *Window) IconButton(button *Clickable) *IconButton {
	return &IconButton{
		Window:      w,
		background: "Primary",
		color:      "DocBg",
		size:       w.TextSize,
		inset:      w.Inset(0.33, nil),
		button:     button,
		icon:       w.Icon().Src(&icons.AlertError),
	}
}

// Background sets the color of the circular background
func (b *IconButton) Background(color string) *IconButton {
	b.background = color
	return b
}

// Color sets the color of the icon
func (b *IconButton) Color(color string) *IconButton {
	b.color = color
	return b
}

// Icon sets the icon to display
func (b *IconButton) Icon(ic *Icon) *IconButton {
	b.icon = ic
	return b
}

// Scale changes the size of the icon as a ratio of the base font size
func (b *IconButton) Scale(scale float32) *IconButton {
	b.size = b.Theme.TextSize.Scale(scale * 0.72)
	return b
}

// Inset sets the size of inset that goes in between the button background and the icon
func (b *IconButton) ButtonInset(inset float32) *IconButton {
	b.inset = b.Inset(inset, b.button.Fn)
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
	return b.Stack().Expanded(
		func(gtx l.Context) l.Dimensions {
			sizex, sizey := gtx.Constraints.Min.X, gtx.Constraints.Min.Y
			sizexf, sizeyf := float32(sizex), float32(sizey)
			rr := (sizexf + sizeyf) * .25
			clip.RRect{
				Rect: f32.Rectangle{Max: f32.Point{X: sizexf, Y: sizeyf}},
				NE:   rr, NW: rr, SE: rr, SW: rr,
			}.Add(gtx.Ops)
			background := b.Theme.Colors.Get(b.background)
			if gtx.Queue == nil {
				background = f32color.MulAlpha(background, 150)
			}
			var dims l.Dimensions
			if b.background != "" {
				dims = Fill(gtx, background)
			}
			for _, c := range b.button.History() {
				drawInk(gtx, c)
			}
			return dims
		},
	).Stacked(
		b.inset.Embed(b.icon.Fn).Fn,
	).Expanded(
		func(gtx l.Context) l.Dimensions {
			pointer.Ellipse(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
			return b.button.Fn(gtx)
		},
	).Fn(gtx)
}
