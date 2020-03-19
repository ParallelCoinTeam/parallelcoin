package monitor

import (
	"gioui.org/layout"
)

func (st *State) Body() layout.FlexChild {
	return Flexed(1, func() {
		cs := st.Gtx.Constraints
		st.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
		st.Inset(8, func(){

		})
	})
}
