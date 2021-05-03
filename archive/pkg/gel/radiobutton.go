// SPDX-License-Identifier: Unlicense OR MIT

package gel

import (
	"golang.org/x/exp/shiny/materialdesign/icons"

	l "github.com/p9c/gio/layout"
)

type RadioButton struct {
	*Checkable
	*Window
	key   string
	group *Enum
}

// RadioButton returns a RadioButton with a label. The key specifies the value for the Enum.
func (w *Window) RadioButton(checkable *Checkable, group *Enum, key,
	label string) *RadioButton {
	// if checkable == nil {
	// 	debug.PrintStack()
	// 	os.Exit(0)
	// }
	return &RadioButton{
		group:  group,
		Window: w,
		Checkable: checkable.
			CheckedStateIcon(&icons.ToggleRadioButtonChecked). // Color("Primary").
			UncheckedStateIcon(&icons.ToggleRadioButtonUnchecked). // Color("PanelBg").
			Label(label), // .Color("DocText").IconColor("PanelBg"),
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
func (r RadioButton) Fn(gtx l.Context) l.Dimensions {
	dims := r.Checkable.Fn(gtx, r.group.Value() == r.key)
	gtx.Constraints.Min = dims.Size
	r.group.Fn(gtx, r.key)
	return dims
}
