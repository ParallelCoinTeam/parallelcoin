package monitor

import (
	"gioui.org/layout"
)

func (m *State) Body() layout.FlexChild {
	return Flexed(1, func() {
		cs := m.Gtx.Constraints
		m.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
		m.Inset(8, func(){

		})
	})
}
