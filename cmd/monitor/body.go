package monitor

import (
	"gioui.org/layout"
)

func (s *State) Body() layout.FlexChild {
	return Flexed(1, func() {
		cs := s.Gtx.Constraints
		s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg", "ff")
		s.Inset(8, func() {

		})
	})
}
