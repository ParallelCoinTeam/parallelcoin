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
	background color.RGBA
	// Color is the icon color.
	color color.RGBA
	icon  *widget.Icon
	// Size is the icon size.
	size   unit.Value
	inset  l.Inset
	button *w.Clickable
}

func (th *Theme) IconButton(button *w.Clickable, icon *widget.Icon) *_iconButton {
	return &_iconButton{
		background: th.Colors.Get("Primary"),
		color:      th.Colors.Get("InvText"),
		icon:       icon,
		size:       unit.Sp(24),
		inset:      l.UniformInset(unit.Sp(12)),
		button:     button,
	}
}

func (b *_iconButton) Fn(gtx l.Context) l.Dimensions {
	return l.Stack{Alignment: l.Center}.Layout(gtx,
		l.Expanded(func(gtx l.Context) l.Dimensions {
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
		}),
		l.Stacked(func(gtx l.Context) l.Dimensions {
			return b.inset.Layout(gtx, func(gtx l.Context) l.Dimensions {
				size := gtx.Px(b.size)
				if b.icon != nil {
					b.icon.Color = b.color
					b.icon.Layout(gtx, unit.Px(float32(size)))
				}
				return l.Dimensions{
					Size: image.Point{X: size, Y: size},
				}
			})
		}),
		l.Expanded(func(gtx l.Context) l.Dimensions {
			pointer.Ellipse(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
			return b.button.Fn(gtx)
		}),
	)
}
