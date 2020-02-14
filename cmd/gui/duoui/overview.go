package duoui

import (
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/pkg/gui/layout"
)

func (ui *DuoUI) DuoUIoverview() func() {
	return func() {
		viewport := layout.Flex{Axis: layout.Horizontal}
		if ui.ly.Context.Constraints.Width.Max < 780 {
			viewport = layout.Flex{Axis: layout.Vertical}
		}
		viewport.Layout(ui.ly.Context,
			layout.Flexed(0.5, func() {
				cs := ui.ly.Context.Constraints
				theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, cs.Height.Max, ui.ly.Theme.Color.Light, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				ui.DuoUIbalance()
			}),
			layout.Flexed(0.5, func() {
				cs := ui.ly.Context.Constraints
				theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, cs.Height.Max, "ff424242", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				ui.DuoUIlatestTransactions()
			}),
		)
	}
}