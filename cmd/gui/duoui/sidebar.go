package duoui

import (
	"github.com/p9c/pod/pkg/gui/layout"
)

func (ui *DuoUI)DuoUIsidebar() func() {
	return func() {
		layout.Flex{Axis: layout.Vertical}.Layout(ui.ly.Context,
			layout.Rigid(ui.DuoUImenu()),
		)
	}
}
