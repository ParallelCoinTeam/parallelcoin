package p9

import (
	"image"

	"gioui.org/io/pointer"
	l "gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Checkable struct {
	th                 *Theme
	label              string
	color              string
	font               text.Font
	textSize           unit.Value
	iconColor          string
	size               unit.Value
	checkedStateIcon   *[]byte
	uncheckedStateIcon *[]byte
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
		color:              "Primary",
		font:               f,
		textSize:           th.TextSize.Scale(14.0 / 16.0),
		iconColor:          "Primary",
		size:               th.TextSize.Scale(1.5),
		checkedStateIcon:   &icons.ToggleCheckBox,
		uncheckedStateIcon: &icons.ToggleCheckBoxOutlineBlank,
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
	c.color = color
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
	c.iconColor = color
	return c
}

// Scale sets the size of the checkbox icon relative to the base font size
func (c *Checkable) Scale(size float32) *Checkable {
	c.size = c.th.TextSize.Scale(size)
	return c
}

// CheckedStateIcon loads the icon for the checked state
func (c *Checkable) CheckedStateIcon(ic *[]byte) *Checkable {
	c.checkedStateIcon = ic
	return c
}

// UncheckedStateIcon loads the icon for the unchecked state
func (c *Checkable) UncheckedStateIcon(ic *[]byte) *Checkable {
	c.uncheckedStateIcon = ic
	return c
}

// Fn renders the checkbox widget
func (c *Checkable) Fn(gtx l.Context, checked bool) l.Dimensions {
	var icon *Icon
	if checked {
		icon = c.th.Icon().
			Color(c.color).
			Src(c.checkedStateIcon)
	} else {
		icon = c.th.Icon().
			Color(c.color).
			Src(c.uncheckedStateIcon)
	}
	icon.size = c.size
	// Debugs(icon)
	dims :=
		c.th.Flex().Rigid(
			// c.th.Inset(0.25,
			func(gtx l.Context) l.Dimensions {
				size := gtx.Px(c.size)
				// icon.color = c.iconColor
				// TODO: maybe make a special code for raw colors to do this kind of alpha
				//  or add a parameter to apply it
				// if gtx.Queue == nil {
				// 	icon.color = f32color.MulAlpha(c.th.Colors.Get(icon.color), 150)
				// }
				icon.Fn(gtx)
				return l.Dimensions{
					Size: image.Point{X: size, Y: size},
				}
			},
			// ).Fn,
		).Rigid(
			// c.th.Inset(0.25,
			func(gtx l.Context) l.Dimensions {
				paint.ColorOp{Color: c.th.Colors.Get(c.color)}.Add(gtx.Ops)
				return c.th.Caption(c.label).Color(c.color).Fn(gtx)
				// return widget.Label{}.Layout(gtx, c.shaper, c.font, c.textSize, c.label)
			},
			// ).Fn,
		).Fn(gtx)
	pointer.Rect(image.Rectangle{Max: dims.Size}).Add(gtx.Ops)
	return dims
}
