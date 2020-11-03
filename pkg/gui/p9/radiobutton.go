// SPDX-License-Identifier: Unlicense OR MIT

package p9

import (
	"gioui.org/layout"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type RadioButton struct {
	*Checkable
	th    *Theme
	key   string
	group *Enum
}

// RadioButton returns a RadioButton with a label. The key specifies the value for the Enum.
func (th *Theme) RadioButton(checkable *Checkable, group *Enum, key, label string) *RadioButton {
	// if checkable == nil {
	// 	debug.PrintStack()
	// 	os.Exit(0)
	// }
	return &RadioButton{
		group: group,
		th:    th,
		Checkable: checkable.
			CheckedStateIcon(icons.ToggleRadioButtonChecked).
			UncheckedStateIcon(icons.ToggleRadioButtonUnchecked).
			Label(label),
		key: key,
	}
}

// Key sets the key initially active on the radiobutton
func (r *RadioButton) Key(key string) *RadioButton {
	r.key = key
	return r
}

// Group sets the enum group of the radio button
func (r *RadioButton) Group(group *Enum) *RadioButton {
	r.group = group
	return r
}

// Fn updates enum and displays the radio button.
func (r RadioButton) Fn(gtx layout.Context) layout.Dimensions {
	dims := r.Checkable.Fn(gtx, r.group.Value() == r.key)
	gtx.Constraints.Min = dims.Size
	r.group.Fn(gtx, r.key)
	return dims
}
