// SPDX-License-Identifier: Unlicense OR MIT

package parallel

import (
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/pkg/gui/widget"
)

type DuoUIcheckBox struct {
	checkable
}

func (t *DuoUItheme) DuoUIcheckBox(label string) DuoUIcheckBox {
	return DuoUIcheckBox{
		checkable{
			Label:     label,
			Color:     t.Color.Text,
			IconColor: t.Color.Primary,
			Font: text.Font{
				Typeface: t.Font.Primary,
				Size: t.TextSize.Scale(14.0 / 16.0),
			},
			Size:               unit.Dp(26),
			shaper:             t.Shaper,
			checkedStateIcon:   t.checkBoxCheckedIcon,
			uncheckedStateIcon: t.checkBoxUncheckedIcon,
		},
	}
}

func (c DuoUIcheckBox) Layout(gtx *layout.Context, checkBox *widget.CheckBox) {
	c.layout(gtx, checkBox.Checked(gtx))
	checkBox.Layout(gtx)
}
