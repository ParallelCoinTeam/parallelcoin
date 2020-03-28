package monitor

import (
	"gioui.org/layout"
	"github.com/p9c/pod/pkg/logi"
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
				color := "DocText"
				switch b.Level {
				case logi.Tags[logi.Trace]:
					color = "Hint"
				case logi.Tags[logi.Debug]:
					color = "Info"
				case logi.Tags[logi.Info]:
					color = "Success"
				case logi.Tags[logi.Warn]:
					color = "Warning"
				case logi.Tags[logi.Check], logi.Tags[logi.Error],
					logi.Tags[logi.Fatal]:
					color = "Danger"
					//case "FTL":
					//	color = "Danger"
				}
				s.FlexH(
					Rigid(
						s.Text(b.Level, color, "DocBg", "Mono", "body1"),
					),
					Rigid(
						s.Text(b.Text, "DocText", "DocBg", "Mono", "body1"),
					),
				)
			})
		})
		s.W.Invalidate()
	})
}
