package gel

import (
	l "github.com/p9c/gio/layout"
)

// CheckBox creates a checkbox with a text label
func (w *Window) CheckBox(checkBox *Bool) *Checkbox {
	var (
		color     = "DocText"
		textColor = "Primary"
		label     = "this is a label"
	)
	chk := w.Checkable()
	chk.Font("bariol regular").Color(textColor)
	return &Checkbox{
		color:     color,
		textColor: textColor,
		label:     label,
		checkBox:  checkBox,
		Checkable: chk,
		action:    func(b bool) {},
	}
}

type Checkbox struct {
	*Checkable
	checkBox                *Bool
	color, textColor, label string
	action                  func(b bool)
}

// IconColor sets the color of the icon in the checkbox
func (c *Checkbox) IconColor(color string) *Checkbox {
	c.Checkable.iconColor = color
	return c
}

// TextColor sets the color of the text label
func (c *Checkbox) TextColor(color string) *Checkbox {
	c.Checkable.color = color
	return c
}

// TextScale sets the scale relative to the base font size for the text label
func (c *Checkbox) TextScale(scale float32) *Checkbox {
	c.textSize = c.TextSize.Scale(scale)
	return c
}

// Text sets the text to be rendered on the checkbox
func (c *Checkbox) Text(label string) *Checkbox {
	c.Checkable.label = label
	return c
}

// IconScale sets the scaling of the check icon
func (c *Checkbox) IconScale(scale float32) *Checkbox {
	c.size = c.TextSize.Scale(scale)
	return c
}

// SetOnChange sets the callback when a state change event occurs
func (c *Checkbox) SetOnChange(fn func(b bool)) *Checkbox {
	c.action = fn
	return c
}

// Fn renders the checkbox
func (c *Checkbox) Fn(gtx l.Context) l.Dimensions {
	dims := c.Checkable.Fn(gtx, c.checkBox.GetValue())
	gtx.Constraints.Min = dims.Size
	c.checkBox.Fn(gtx)
	return dims
}
