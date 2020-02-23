package duoui

import (
	"gioui.org/layout"
)

func (ui *DuoUI) DuoUIoverview() func() {
	return func() {
		viewport := layout.Flex{Axis: layout.Horizontal}
		if ui.ly.Context.Constraints.Width.Max < 780 {
			viewport = layout.Flex{Axis: layout.Vertical}
		}
		viewport.Layout(ui.ly.Context,
			layout.Flexed(0.5, ui.DuoUIbalance()),
			layout.Flexed(0.5, ui.DuoUIlatestTransactions()),
		)
	}
}
