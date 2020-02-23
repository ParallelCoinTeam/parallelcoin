package duoui

import (
	"gioui.org/layout"
)

func (ui *DuoUI) DuoUIsidebar() func() {
	return func() {
		layout.Flex{Axis: layout.Vertical}.Layout(ui.ly.Context,
			layout.Rigid(ui.DuoUImenu()),
		)
	}
}
