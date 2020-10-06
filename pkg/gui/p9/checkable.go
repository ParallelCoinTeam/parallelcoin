package p9

import (
	"image"
	"image/color"

	"gioui.org/io/pointer"
	l "gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/p9c/pod/pkg/gui/f32color"
)

type _checkable struct {
	label              string
	color              color.RGBA
	font               text.Font
	textSize           unit.Value
	iconColor          color.RGBA
	size               unit.Value
	shaper             text.Shaper
	checkedStateIcon   *_icon
	uncheckedStateIcon *_icon
}

func (th *Theme) Checkable() *_checkable {
	return &_checkable{}
}

func (c *_checkable) fn(gtx l.Context, checked bool) l.Dimensions {
	var icon *_icon
	if checked {
		icon = c.checkedStateIcon
	} else {
		icon = c.uncheckedStateIcon
	}
	min := gtx.Constraints.Min
	dims := l.Flex{Alignment: l.Middle}.Layout(gtx,
		l.Rigid(func(gtx l.Context) l.Dimensions {
			return l.Center.Layout(gtx, func(gtx l.Context) l.Dimensions {
				return l.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx l.Context) l.Dimensions {
					size := gtx.Px(c.size)
					icon.color = c.iconColor
					if gtx.Queue == nil {
						icon.color = f32color.MulAlpha(icon.color, 150)
					}
					icon.Fn(&gtx, unit.Px(float32(size)))
					return l.Dimensions{
						Size: image.Point{X: size, Y: size},
					}
				})
			})
		}),
		l.Rigid(func(gtx l.Context) l.Dimensions {
			gtx.Constraints.Min = min
			return l.W.Layout(gtx, func(gtx l.Context) l.Dimensions {
				return l.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx l.Context) l.Dimensions {
					paint.ColorOp{Color: c.color}.Add(gtx.Ops)
					return widget.Label{}.Layout(gtx, c.shaper, c.font, c.textSize, c.label)
				})
			})
		}),
	)
	pointer.Rect(image.Rectangle{Max: dims.Size}).Add(gtx.Ops)
	return dims
}
