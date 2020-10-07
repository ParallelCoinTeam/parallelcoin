package p9

import (
	l "gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

type _radioButton struct {
	_checkable
	key   string
	group *widget.Enum
}

// RadioButton returns a RadioButton with a label. The key specifies the value for the Enum.
func (th *Theme) RadioButton(group *widget.Enum, key, label string) *_radioButton {
	return &_radioButton{
		group: group,
		_checkable: _checkable{
			label:              label,
			color:              th.Colors.Get("Text"),
			iconColor:          th.Colors.Get("Primary"),
			textSize:           th.textSize.Scale(14.0 / 16.0),
			size:               unit.Dp(26),
			shaper:             th.shaper,
			checkedStateIcon:   th.icons["Checked"],
			uncheckedStateIcon: th.icons["Unchecked"],
		},
		key: key,
	}
}

// Fn updates enum and displays the radio _button.
func (r *_radioButton) Fn(gtx l.Context) l.Dimensions {
	dims := r._checkable.Fn(gtx, r.group.Value == r.key)
	gtx.Constraints.Min = dims.Size
	r.group.Layout(gtx, r.key)
	return dims
}
