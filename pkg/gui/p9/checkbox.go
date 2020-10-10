package p9

import (
	l "gioui.org/layout"
)

// CheckBox creates a checkbox with a text label
func (th *Theme) CheckBox(checkBox *_bool) *_checkbox {
	var (
		color     = "DocText"
		textColor = "Primary"
		label     = "this is a label"
	)
	chk := th.Checkable()
	chk.Font("bariol regular").Color(textColor)
	return &_checkbox{
		color:      color,
		textColor:  textColor,
		label:      label,
		checkBox:   checkBox,
		_checkable: chk,
		action:     func(b bool) {},
	}
}

type _checkbox struct {
	*_checkable
	checkBox                *_bool
	color, textColor, label string
	action                  func(b bool)
}

// IconColor sets the color of the icon in the checkbox
func (c *_checkbox) IconColor(color string) *_checkbox {
	c._checkable.iconColor = c.th.Colors.Get(color)
	return c
}

// TextColor sets the color of the text label
func (c *_checkbox) TextColor(color string) *_checkbox {
	c._checkable.color = c.th.Colors.Get(color)
	return c
}

// TextScale sets the scale relative to the base font size for the text label
func (c *_checkbox) TextScale(scale float32) *_checkbox {
	c.textSize = c.th.textSize.Scale(scale)
	return c
}

// Text sets the text to be rendered on the checkbox
func (c *_checkbox) Text(label string) *_checkbox {
	c._checkable.label = label
	return c
}

// IconScale sets the scaling of the check icon
func (c *_checkbox) IconScale(scale float32) *_checkbox {
	c.size = c.th.textSize.Scale(scale)
	return c
}

// SetOnChange sets the callback when a state change event occurs
func (c *_checkbox) SetOnChange(fn func(b bool)) *_checkbox {
	c.action = fn
	return c
}

// Fn renders the checkbox
func (c *_checkbox) Fn(gtx l.Context) l.Dimensions {
	dims := c._checkable.Fn(gtx, c.checkBox.GetValue())
	gtx.Constraints.Min = dims.Size
	c.checkBox.Fn(gtx)
	return dims
}
