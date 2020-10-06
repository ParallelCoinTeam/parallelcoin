package p9

import (
	"image/color"

	l "gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

type _checkbox struct {
	_checkable
	checkBox *widget.Bool
}

func (th *Theme) CheckBox(checkBox *widget.Bool, color color.RGBA, label string) *_checkbox {
	return &_checkbox{
		checkBox: checkBox,
		_checkable: _checkable{
			label:              label,
			color:              color,
			iconColor:          th.Colors.Get("Primary"),
			textSize:           th.textSize.Scale(14.0 / 16.0),
			size:               unit.Dp(26),
			shaper:             th.shaper,
			checkedStateIcon:   th.icons["Checked"],
			uncheckedStateIcon: th.icons["Unchecked"],
		},
	}
}

// Layout updates the checkBox and displays it.
func (c *_checkbox) Layout(gtx l.Context) l.Dimensions {
	dims := c.fn(gtx, c.checkBox.Value)
	gtx.Constraints.Min = dims.Size
	c.checkBox.Layout(gtx)
	return dims
}
