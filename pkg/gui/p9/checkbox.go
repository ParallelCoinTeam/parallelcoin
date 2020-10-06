package p9

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

type CheckBoxStyle struct {
	checkable
	CheckBox *widget.Bool
}

func CheckBox(th *Theme, checkBox *widget.Bool, color color.RGBA, label string) CheckBoxStyle {
	return CheckBoxStyle{
		CheckBox: checkBox,
		checkable: checkable{
			Label:              label,
			Color:              color,
			IconColor:          th.Colors.Get("Primary"),
			TextSize:           th.TextSize.Scale(14.0 / 16.0),
			Size:               unit.Dp(26),
			shaper:             th.Shaper,
			checkedStateIcon:   th.Icons["Checked"],
			uncheckedStateIcon: th.Icons["Unchecked"],
		},
	}
}

// Layout updates the checkBox and displays it.
func (c CheckBoxStyle) Layout(gtx layout.Context) layout.Dimensions {
	dims := c.layout(gtx, c.CheckBox.Value)
	gtx.Constraints.Min = dims.Size
	c.CheckBox.Layout(gtx)
	return dims
}
