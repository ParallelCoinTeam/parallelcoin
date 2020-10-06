package p9

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

type RadioButtonStyle struct {
	checkable
	Key   string
	Group *widget.Enum
}

// RadioButton returns a RadioButton with a label. The key specifies the value for the Enum.
func (th *Theme) RadioButton(group *widget.Enum, key, label string) RadioButtonStyle {
	return RadioButtonStyle{
		Group: group,
		checkable: checkable{
			Label: label,

			Color:              th.Colors.Get("Text"),
			IconColor:          th.Colors.Get("Primary"),
			TextSize:           th.TextSize.Scale(14.0 / 16.0),
			Size:               unit.Dp(26),
			shaper:             th.Shaper,
			checkedStateIcon:   th.Icons["Checked"],
			uncheckedStateIcon: th.Icons["Unchecked"],
		},
		Key: key,
	}
}

// Fn updates enum and displays the radio button.
func (r RadioButtonStyle) Fn(gtx layout.Context) layout.Dimensions {
	dims := r.layout(gtx, r.Group.Value == r.Key)
	gtx.Constraints.Min = dims.Size
	r.Group.Layout(gtx, r.Key)
	return dims
}
