// SPDX-License-Identifier: Unlicense OR MIT

package p9

import (
	"gioui.org/layout"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type _radioButton struct {
	_checkable
	th    *Theme
	key   string
	group *_enum
}

// RadioButton returns a RadioButton with a label. The key specifies the value for the _enum.
func (th *Theme) RadioButton(group *_enum, key, label string) *_radioButton {
	return &_radioButton{
		group: group,
		th:    th,
		_checkable: *th.Checkable().
			CheckedStateIcon(th.Icon().Src(icons.ToggleRadioButtonChecked)).
			UncheckedStateIcon(th.Icon().Src(icons.ToggleRadioButtonUnchecked)).
			Label(label),
		key: key,
	}
}

func (r *_radioButton) Key(key string) *_radioButton {
	r.key = key
	return r
}

func (r *_radioButton) Group(group *_enum) *_radioButton {
	r.group = group
	return r
}

// Fn updates enum and displays the radio button.
func (r _radioButton) Fn(gtx layout.Context) layout.Dimensions {
	dims := r._checkable.Fn(gtx, r.group.Value == r.key)
	gtx.Constraints.Min = dims.Size
	r.group.Layout(gtx, r.key)
	return dims
}
