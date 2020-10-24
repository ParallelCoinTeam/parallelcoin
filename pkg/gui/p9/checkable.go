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

type Checkable struct {
	th                 *Theme
	label              string
	color              color.RGBA
	font               text.Font
	textSize           unit.Value
	iconColor          color.RGBA
	size               unit.Value
	checkedStateIcon   *Icon
	uncheckedStateIcon *Icon
	shaper             text.Shaper
	checked            bool
}

// Checkable creates a checkbox type widget
func (th *Theme) Checkable() *Checkable {
	font := "bariol regular"
	var f text.Font
	for i := range th.collection {
		if th.collection[i].Font.Typeface == text.Typeface(font) {
			f = th.collection[i].Font
			break
		}
	}
	return &Checkable{
		th:                 th,
		label:              "checkable",
		color:              th.Colors.Get("Primary"),
		font:               f,
		textSize:           th.TextSize.Scale(14.0 / 16.0),
		iconColor:          th.Colors.Get("Primary"),
		size:               th.TextSize.Scale(1),
		checkedStateIcon:   th.Icon().Src(icons.ToggleCheckBox).Color("Primary"),
		uncheckedStateIcon: th.Icon().Src(icons.ToggleCheckBoxOutlineBlank).Color("Primary"),
		shaper:             th.shaper,
	}
}

// Text sets the label on the checkbox
func (c *Checkable) Label(txt string) *Checkable {
	c.label = txt
	return c
}

// Color sets the color of the checkbox label
func (c *Checkable) Color(color string) *Checkable {
	c.color = c.th.Colors.Get(color)
	return c
}

// Font sets the font used on the label
func (c *Checkable) Font(font string) *Checkable {
	for i := range c.th.collection {
		if c.th.collection[i].Font.Typeface == text.Typeface(font) {
			c.font = c.th.collection[i].Font
			break
		}
	}
	return c
}

// TextScale sets the size of the font relative to the base text size
func (c *Checkable) TextScale(scale float32) *Checkable {
	c.textSize = c.th.TextSize.Scale(scale)
	return c
}

// IconColor sets the color of the icon
func (c *Checkable) IconColor(color string) *Checkable {
	c.iconColor = c.th.Colors.Get(color)
	return c
}

// Scale sets the size of the checkbox icon relative to the base font size
func (c *Checkable) Scale(size float32) *Checkable {
	c.size = c.th.TextSize.Scale(size)
	return c
}

// CheckedStateIcon loads the icon for the checked state
func (c *Checkable) CheckedStateIcon(ic *Icon) *Checkable {
	c.checkedStateIcon = ic
	return c
}

// UncheckedStateIcon loads the icon for the unchecked state
func (c *Checkable) UncheckedStateIcon(ic *Icon) *Checkable {
	c.uncheckedStateIcon = ic
	return c
}

// Fn renders the checkbox widget
func (c *Checkable) Fn(gtx l.Context, checked bool) l.Dimensions {
	var icon *Icon
	if checked {
		icon = c.checkedStateIcon.Scale(1.5)
	} else {
		icon = c.uncheckedStateIcon.Scale(1.5)
	}
	dims :=
		c.th.Flex().Rigid(
			c.th.Inset(0.25,
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
			c.th.Inset(0.25,
				func(gtx l.Context) l.Dimensions {
					paint.ColorOp{Color: c.color}.Add(gtx.Ops)
					return widget.Label{}.Layout(gtx, c.shaper, c.font, c.textSize, c.label)
				},
			).Fn,
		).Fn(gtx)
	pointer.Rect(image.Rectangle{Max: dims.Size}).Add(gtx.Ops)
	return dims
}
