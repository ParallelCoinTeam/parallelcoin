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
	checkedStateIcon   *_icon
	uncheckedStateIcon *_icon
	shaper             text.Shaper
	checked            bool
}

// Checkable creates a checkbox type widget
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
		textSize:           th.TextSize.Scale(14.0 / 16.0),
		iconColor:          th.Colors.Get("Primary"),
		size:               th.TextSize.Scale(1),
		checkedStateIcon:   th.Icon().Src(icons.ToggleCheckBox).Color("Primary"),
		uncheckedStateIcon: th.Icon().Src(icons.ToggleCheckBoxOutlineBlank).Color("Primary"),
		shaper:             th.shaper,
	}
}

// _text sets the label on the checkbox
func (c *_checkable) Label(txt string) *_checkable {
	c.label = txt
	return c
}

// Color sets the color of the checkbox label
func (c *_checkable) Color(color string) *_checkable {
	c.color = c.th.Colors.Get(color)
	return c
}

// Font sets the font used on the label
func (c *_checkable) Font(font string) *_checkable {
	for i := range c.th.collection {
		if c.th.collection[i].Font.Typeface == text.Typeface(font) {
			c.font = c.th.collection[i].Font
			break
		}
	}
	return c
}

// TextScale sets the size of the font relative to the base text size
func (c *_checkable) TextScale(scale float32) *_checkable {
	c.textSize = c.th.TextSize.Scale(scale)
	return c
}

// IconColor sets the color of the icon
func (c *_checkable) IconColor(color string) *_checkable {
	c.iconColor = c.th.Colors.Get(color)
	return c
}

// Scale sets the size of the checkbox icon relative to the base font size
func (c *_checkable) Scale(size float32) *_checkable {
	c.size = c.th.TextSize.Scale(size)
	return c
}

// CheckedStateIcon loads the icon for the checked state
func (c *_checkable) CheckedStateIcon(ic *_icon) *_checkable {
	c.checkedStateIcon = ic
	return c
}

// UncheckedStateIcon loads the icon for the unchecked state
func (c *_checkable) UncheckedStateIcon(ic *_icon) *_checkable {
	c.uncheckedStateIcon = ic
	return c
}

// Fn renders the checkbox widget
func (c *_checkable) Fn(gtx l.Context, checked bool) l.Dimensions {
	var icon *_icon
	if checked {
		icon = c.checkedStateIcon.Scale(1.5)
	} else {
		icon = c.uncheckedStateIcon.Scale(1.5)
	}
	dims :=
		c.th.Flex().Rigid(
			c.th.Inset(0.25).Embed(
				func(gtx l.Context) l.Dimensions {
					size := gtx.Px(c.size)
					icon.color = c.iconColor
					if gtx.Queue == nil {
						icon.color = f32color.MulAlpha(icon.color, 150)
					}
					icon.Fn(gtx)
					return l.Dimensions{
						Size: image.Point{X: size, Y: size},
					}
				}).Fn,
		).Rigid(
			c.th.Inset(0.25).Embed(
				func(gtx l.Context) l.Dimensions {
					paint.ColorOp{Color: c.color}.Add(gtx.Ops)
					return widget.Label{}.Layout(gtx, c.shaper, c.font, c.textSize, c.label)
				},
			).Fn,
		).Fn(gtx)
	pointer.Rect(image.Rectangle{Max: dims.Size}).Add(gtx.Ops)
	return dims
}
