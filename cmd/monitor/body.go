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
				//fmt.Println("level", b.Level)
				switch b.Level {
				case logi.Trace:
					color = "Secondary"
				case logi.Debug:
					color = "Info"
				case logi.Info:
					color = "Success"
				case logi.Warn:
					color = "Warning"
				case logi.Check, logi.Error, logi.Fatal:
					color = "Danger"
					//case "FTL":
					//	color = "Danger"
				}
				s.FlexH(
					//Rigid(
					//	s.Text(fmt.Sprint(i), color, "DocBg", "Mono", "body1"),
					//),
					Rigid(
						//s.Text(b.Level, color, "DocBg", "Mono", "body1"),
						func() {
							s.Icon(logi.Tags[b.Level], color, "DocBg", 32)
						},
					),
					Rigid(
						s.Text(b.Time.Format("15:04:05"), color, "DocBg",
							"Mono",
							"body1"),
					),
					Flexed(1,
						s.Text(b.Text, "DocText", "DocBg", "Mono",
							"body1"),
					),
					Spacer(),
					Rigid(
						s.Text(b.Package, "PanelBg", "DocBg", "Primary",
							"h6"),
					),
				)
			})
		})
		s.W.Invalidate()
	})
}
