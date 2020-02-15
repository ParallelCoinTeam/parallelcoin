// SPDX-License-Identifier: Unlicense OR MIT

package theme

import (
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
)

type DuoUIcheckBox struct {
	checkable
}

func (t *DuoUItheme) DuoUIcheckBox(label string) DuoUIcheckBox {
	return DuoUIcheckBox{
		checkable{
			Label:     label,
			Color:     HexARGB(t.Color.Light),
			IconColor: HexARGB(t.Color.Light),
			Font: text.Font{
				Typeface: t.Font.Primary,
				Size: t.TextSize.Scale(14.0 / 16.0),
			},
			Size:               unit.Dp(26),
			shaper:             t.Shaper,
			checkedStateIcon:   t.Icons["checkBoxCheckedIcon"],
			uncheckedStateIcon: t.Icons["checkBoxUncheckedIcon"],
		},
	}
}

func (c DuoUIcheckBox) Layout(gtx *layout.Context, checkBox *controller.CheckBox) {
	c.layout(gtx, checkBox.Checked(gtx))
	checkBox.Layout(gtx)
}
