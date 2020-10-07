package p9

import (
	l "gioui.org/layout"

	w "github.com/p9c/pod/pkg/gui/widget"
)

type _checkbox struct {
	*_checkable
	checkBox                *w.Bool
	color, textColor, label string
	action                  func(b, cs bool)
}

func (th *Theme) CheckBox(checkBox *w.Bool) *_checkbox {
	var (
		color     = "DocText"
		textColor = "Primary"
		label     = "this is a label"
	)
	chk := &_checkable{
		th:                 th,
		label:              label,
		color:              th.Colors.Get(textColor),
		iconColor:          th.Colors.Get(color),
		textSize:           th.textSize.Scale(14.0 / 16.0),
		size:               th.textSize.Scale(14.0 / 16.0 * 2),
		shaper:             th.shaper,
		checkedStateIcon:   th.icons["Checked"],
		uncheckedStateIcon: th.icons["Unchecked"],
	}
	chk = chk.Font("bariol regular").Color(textColor)
	return &_checkbox{
		color:      color,
		textColor:  textColor,
		label:      label,
		checkBox:   checkBox,
		_checkable: chk,
	}
}

func (c *_checkbox) IconColor(color string) *_checkbox {
	c._checkable.iconColor = c.th.Colors.Get(color)
	return c
}

func (c *_checkbox) TextColor(color string) *_checkbox {
	c._checkable.color = c.th.Colors.Get(color)
	return c
}

func (c *_checkbox) Label(label string) *_checkbox {
	c._checkable.label = label
	return c
}

func (c *_checkbox) Action(fn func(b bool)) *_checkbox {
	c.checkBox.SetHook(fn)
	return c
}

// Fn updates the checkBox and displays it.
func (c *_checkbox) Fn(gtx l.Context) l.Dimensions {
	dims := c._checkable.Fn(gtx, *c.checkBox.GetValue())
	gtx.Constraints.Min = dims.Size
	c.checkBox.Fn(gtx)
	return dims
}
