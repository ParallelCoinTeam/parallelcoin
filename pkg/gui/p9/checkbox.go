package p9

import (
	l "gioui.org/layout"

	w "github.com/p9c/pod/pkg/gui/widget"
)

type _checkbox struct {
	*_checkable
	checkBox *w.Bool
}

func (th *Theme) CheckBox(checkBox *w.Bool, color, textColor string, label string) *_checkbox {
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
		checkBox:   checkBox,
		_checkable: chk,
	}
}

// Fn updates the checkBox and displays it.
func (c *_checkbox) Fn(gtx l.Context) l.Dimensions {
	dims := c._checkable.Fn(gtx, *c.checkBox.GetValue())
	gtx.Constraints.Min = dims.Size
	c.checkBox.Fn(gtx)
	return dims
}
