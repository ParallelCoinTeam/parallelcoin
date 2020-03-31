package monitor

import (
	"gioui.org/layout"
	"github.com/p9c/pod/pkg/logi"
	"os/exec"
	"strings"
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
				button := s.EntryBuf.GetButton(i)
				s.ButtonArea(func() {
					s.FlexH(
						//Rigid(
						//	s.Text(fmt.Sprint(i), color, "DocBg", "Mono", "body1"),
						//),
						Rigid(
							//s.Text(b.Level, color, "DocBg", "Mono", "body1"),
							func() {
								s.Icon(logi.Tags[b.Level], color, "Transparent",
									24)
							},
						),
						Rigid(
							s.Text(b.Time.Format("15:04:05"), color, "Transparent",
								"Mono",
								"body2"),
						),
						Flexed(1,
							s.Text(b.Text, "DocText", "Transparent", "Mono",
								"body2"),
						),
						Spacer(),
						Rigid(
							s.Text(b.Package, "PanelBg", "Transparent", "Primary",
								"h6"),
						),
					)
				}, button)
				for button.Clicked(s.Gtx) {
					go func() {
						if s.Config.ClickCommand == "" {
							return
						}
						split := strings.Split(b.CodeLocation, ":")
						v1 := split[0]
						v2 := split[1]
						c := strings.Replace(s.Config.ClickCommand, "$1", v1, 1)
						c = strings.Replace(c, "$2", v2, 1)
						Debug("running command", c)
						args := strings.Split(c, " ")
						cmd := exec.Command(args[0], args[1:]...)
						_ = cmd.Run()
						//s.Config.ClickCommand
					}()
				}
			})
		})
		s.W.Invalidate()
	})
}
