// SPDX-License-Identifier: Unlicense OR MIT

package theme

import (
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
)

type DuoUIcheckBox struct {
	checkable
}

func (t *DuoUItheme) DuoUIcheckBox(label, color, iconColor string) DuoUIcheckBox {
	return DuoUIcheckBox{
		checkable{
			Font: text.Font{
				Typeface: t.Font.Primary,
			},
			Label:              label,
			Color:              HexARGB(color),
			IconColor:          HexARGB(iconColor),
			TextSize:           t.TextSize.Scale(14.0 / 16.0),
			Size:               unit.Dp(26),
			shaper:             t.Shaper,
			checkedStateIcon:   t.checkBoxCheckedIcon,
			uncheckedStateIcon: t.checkBoxUncheckedIcon,
		},
	}
}

func (c DuoUIcheckBox) Layout(gtx *layout.Context, checkBox *controller.CheckBox) {
	c.layout(gtx, checkBox.Checked(gtx))
	checkBox.Layout(gtx)
}
