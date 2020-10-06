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
	Label              string
	Color              color.RGBA
	Font               text.Font
	TextSize           unit.Value
	IconColor          color.RGBA
	Size               unit.Value
	shaper             text.Shaper
	checkedStateIcon   *_icon
	uncheckedStateIcon *_icon
}

func (c *_checkable) layout(gtx l.Context, checked bool) l.Dimensions {
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
					size := gtx.Px(c.Size)
					icon.Color = c.IconColor
					if gtx.Queue == nil {
						icon.Color = f32color.MulAlpha(icon.Color, 150)
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
					paint.ColorOp{Color: c.Color}.Add(gtx.Ops)
					return widget.Label{}.Layout(gtx, c.shaper, c.Font, c.TextSize, c.Label)
				})
			})
		}),
	)
	pointer.Rect(image.Rectangle{Max: dims.Size}).Add(gtx.Ops)
	return dims
}
