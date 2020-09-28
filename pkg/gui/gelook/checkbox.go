// SPDX-License-Identifier: Unlicense OR MIT

package gelook

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/p9c/pod/pkg/gui/gel"
)

type DuoUICheckBox struct {
	checkable
}

func (t *DuoUITheme) DuoUICheckBox(label, color, iconColor string) DuoUICheckBox {
	return DuoUICheckBox{
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

func (c DuoUICheckBox) Layout(gtx *layout.Context, checkBox *gel.CheckBox) {
	c.layout(gtx, checkBox.Checked(gtx))
	checkBox.Layout(gtx)
}

func (c DuoUICheckBox) DrawLayout(gtx *layout.Context, checkBox *gel.CheckBox) {
	c.drawLayout(gtx, checkBox.Checked(gtx))
	checkBox.Layout(gtx)
}
