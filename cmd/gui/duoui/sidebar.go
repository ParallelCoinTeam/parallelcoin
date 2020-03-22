package duoui

import (
	"gioui.org/layout"
)

func (ui *DuoUI) DuoUIsidebar() func() {
	return func() {
		ui.ly.Theme.DuoUIitem(0, ui.ly.Theme.Colors["Dark"]).Layout(ui.ly.Context, layout.NW, func() {
			layout.Flex{Axis: layout.Vertical}.Layout(ui.ly.Context,
				layout.Rigid(ui.DuoUImenu()),
			)
		})
	}
}
