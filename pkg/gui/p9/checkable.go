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
	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/pod/pkg/gui/f32color"
)

type _checkable struct {
	th                 *Theme
	label              string
	color              color.RGBA
	font               text.Font
	textSize           unit.Value
	iconColor          color.RGBA
	size               unit.Value
	checkedStateIcon   *Ico
	uncheckedStateIcon *Ico
	shaper             text.Shaper
	checked            bool
}

func (th *Theme) Checkable() *_checkable {
	font := "bariol regular"
	var f text.Font
	for i := range th.collection {
		if th.collection[i].Font.Typeface == text.Typeface(font) {
			f = th.collection[i].Font
			break
		}
	}
	return &_checkable{
		th:                 th,
		label:              "checkable",
		color:              th.Colors.Get("DocText"),
		font:               f,
		textSize:           th.textSize.Scale(14.0 / 16.0),
		iconColor:          th.Colors.Get("Primary"),
		size:               th.textSize.Scale(1),
		checkedStateIcon:   th.Icon().Src(icons.ToggleCheckBox).Color("Primary"),
		uncheckedStateIcon: th.Icon().Src(icons.ToggleCheckBoxOutlineBlank).Color("Primary"),
		shaper:             th.shaper,
	}
}

func (c *_checkable) Label(txt string) *_checkable {
	c.label = txt
	return c
}

func (c *_checkable) Color(color string) *_checkable {
	c.color = c.th.Colors.Get(color)
	return c
}

func (c *_checkable) Font(font string) *_checkable {
	for i := range c.th.collection {
		if c.th.collection[i].Font.Typeface == text.Typeface(font) {
			c.font = c.th.collection[i].Font
			break
		}
	}
	return c
}

func (c *_checkable) TextScale(scale float32) *_checkable {
	c.textSize = c.th.textSize.Scale(scale)
	return c
}

func (c *_checkable) IconColor(color string) *_checkable {
	c.iconColor = c.th.Colors.Get(color)
	return c
}

func (c *_checkable) Scale(size float32) *_checkable {
	c.size = c.th.textSize.Scale(size)
	return c
}

func (c *_checkable) CheckedStateIcon(ic *Ico) *_checkable {
	c.checkedStateIcon = ic
	return c
}

func (c *_checkable) UncheckedStateIcon(ic *Ico) *_checkable {
	c.uncheckedStateIcon = ic
	return c
}

func (c *_checkable) Fn(gtx l.Context, checked bool) l.Dimensions {
	var icon *Ico
	if checked {
		icon = c.checkedStateIcon.Scale(2)
	} else {
		icon = c.uncheckedStateIcon.Scale(2)
	}
	min := gtx.Constraints.Min
	dims := l.Flex{Alignment: l.Middle}.Layout(gtx,
		l.Rigid(func(gtx l.Context) l.Dimensions {
			return l.Center.Layout(gtx, func(gtx l.Context) l.Dimensions {
				return l.UniformInset(c.th.textSize.Scale(0.25)).Layout(gtx, func(gtx l.Context) l.Dimensions {
					size := gtx.Px(c.size)
					icon.color = c.iconColor
					if gtx.Queue == nil {
						icon.color = f32color.MulAlpha(icon.color, 150)
					}
					icon.Fn(gtx)
					return l.Dimensions{
						Size: image.Point{X: size, Y: size},
					}
				})
			})
		}),
		l.Rigid(func(gtx l.Context) l.Dimensions {
			gtx.Constraints.Min = min
			return l.W.Layout(gtx, func(gtx l.Context) l.Dimensions {
				return l.UniformInset(c.th.textSize.Scale(0.25)).Layout(gtx, func(gtx l.Context) l.Dimensions {
					paint.ColorOp{Color: c.color}.Add(gtx.Ops)
					return widget.Label{}.Layout(gtx, c.shaper, c.font, c.textSize, c.label)
				})
			})
		}),
	)
	pointer.Rect(image.Rectangle{Max: dims.Size}).Add(gtx.Ops)
	return dims
}
