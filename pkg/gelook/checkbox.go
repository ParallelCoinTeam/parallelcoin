// SPDX-License-Identifier: Unlicense OR MIT

package gelook

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/gel"
)

type DuoUIcheckBox struct {
	checkable
}

func (t *DuoUItheme) DuoUIcheckBox(label, color, iconColor string) DuoUIcheckBox {
	return DuoUIcheckBox{
		checkable{
			Font: text.Font{
				Typeface: t.Fonts["Primary"],
			},
			Label:              label,
			Color:              HexARGB(color),
			IconColor:          HexARGB(iconColor),
			TextSize:           t.TextSize.Scale(14.0 / 16.0),
			Size:               unit.Dp(26),
			shaper:             t.Shaper,
			checkedStateIcon:   t.Icons["Checked"],
			uncheckedStateIcon: t.Icons["Unchecked"],
		},
	}
}

func (c DuoUIcheckBox) Layout(gtx *layout.Context, checkBox *gel.CheckBox) {
	c.layout(gtx, checkBox.Checked(gtx))
	checkBox.Layout(gtx)
}
