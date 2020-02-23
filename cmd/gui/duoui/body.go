package duoui

import (
	"gioui.org/layout"
)

func (ui *DuoUI) DuoUIbody() func() {
	return func() {
		layout.Flex{Axis: layout.Horizontal}.Layout(ui.ly.Context,
			layout.Rigid(ui.DuoUIsidebar()),
			layout.Flexed(1, ui.DuoUIcontent()),
		)
	}
}
