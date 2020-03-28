package monitor

import (
	"gioui.org/layout"
)

func (s *State) Body() layout.FlexChild {
	return Flexed(1, func() {
		cs := s.Gtx.Constraints
		cs.Width.Min = cs.Width.Max / 2
		s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg", "ff")
		s.Inset(8, func() {
			s.LogList.Axis = layout.Vertical
			s.LogList.ScrollToEnd = true
			s.LogList.Layout(s.Gtx, s.EntryBuf.Len(), func(i int) {
				b := s.EntryBuf.Get(i)
				//L.Debugs(b)
				s.FlexH(Rigid(s.Text(b.Text, "DocText", "DocBg",
					"Mono",
					"body1")))
			})
		})
		s.W.Invalidate()
	})
}
