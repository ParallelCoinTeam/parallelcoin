// SPDX-License-Identifier: Unlicense OR MIT

package p9

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type RadioButtonStyle struct {
	_checkable
	th    *Theme
	Key   string
	Group *widget.Enum
}

// RadioButton returns a RadioButton with a label. The key specifies the value for the Enum.
func (th *Theme) RadioButton(group *widget.Enum, key, label string) RadioButtonStyle {
	return RadioButtonStyle{
		Group: group,
		th: th,
		_checkable: *th.Checkable().
			CheckedStateIcon(th.Icon().Src(icons.ToggleRadioButtonChecked)).
			UncheckedStateIcon(th.Icon().Src(icons.ToggleRadioButtonUnchecked)).
			Label(label),
		Key: key,
	}
}

// Fn updates enum and displays the radio button.
func (r RadioButtonStyle) Fn(gtx layout.Context) layout.Dimensions {
	dims := r._checkable.Fn(gtx, r.Group.Value == r.Key)
	gtx.Constraints.Min = dims.Size
	r.Group.Layout(gtx, r.Key)
	return dims
}
