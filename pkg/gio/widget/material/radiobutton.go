// SPDX-License-Identifier: Unlicense OR MIT

package material

import (
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/text"
	"github.com/p9c/pod/pkg/gio/unit"
	"github.com/p9c/pod/pkg/gio/widget"
)

type RadioButton struct {
	checkable
	Key string
}

// RadioButton returns a RadioButton with a label. The key specifies
// the value for the Enum.
func (t *Theme) RadioButton(key, label string) RadioButton {
	return RadioButton{
		checkable: checkable{
			Label: label,

			Color:     t.Color.Text,
			IconColor: t.Color.Primary,
			Font: text.Font{
				Size: t.TextSize.Scale(14.0 / 16.0),
			},
			Size:               unit.Dp(26),
			shaper:             t.Shaper,
			checkedStateIcon:   t.radioCheckedIcon,
			uncheckedStateIcon: t.radioUncheckedIcon,
		},
		Key: key,
	}
}

func (r RadioButton) Layout(gtx *layout.Context, enum *widget.Enum) {
	r.layout(gtx, enum.Value(gtx) == r.Key)
	enum.Layout(gtx, r.Key)
}
