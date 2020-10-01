package monitor

import (
	"os/exec"
	"strings"

	"gioui.org/layout"

	"github.com/p9c/pod/pkg/gui"
	"github.com/p9c/pod/pkg/util/logi"
)

// LogViewer renders the log view
func (s *State) LogViewer() layout.FlexChild {
	return gui.Flexed(1, func() {
		cs := s.Gtx.Constraints
		cs.Width.Min = cs.Width.Max / 2
		eb := s.EntryBuf
		if s.Config.FilterMode {
			eb = s.FilterBuf
		}
		s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
		s.Inset(4, func() {
			l := s.Lists["Log"]
			l.Axis = layout.Vertical
			l.ScrollToEnd = true
			l.Layout(s.Gtx, eb.Len(), func(i int) {
				if eb.Clicked == i {
					cs := s.Gtx.Constraints
					s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBgHilite")
				}
				b := eb.Get(i)
				if b == nil {
					return
				}
				color := "DocText"
				switch b.Level {
				case logi.Trace:
					color = "Secondary"
				case logi.Debug:
					color = "Info"
				case logi.Info:
					color = "Success"
				case logi.Warn:
					color = "Warning"
				case logi.Check:
					color = "Check"
				case logi.Error:
					color = "Danger"
				case logi.Fatal:
					color = "Fatal"
				}
				button := eb.GetButton(i)
				hider := eb.GetHider(i)
				ww := s.WindowWidth
				if s.Config.FilterOpen {
					ww -= 332
				}
				s.FlexHStart(
					gui.Flexed(1, func() {
						s.ButtonArea(func() {
							s.FlexHStart(
								gui.Rigid(func() {
									if ww > 480 {
										s.Inset(4, func() {
											s.Icon(logi.Tags[b.Level], color,
												"Transparent", 24)
										})
									}
								}),
								gui.Rigid(func() {
									if ww > 960 {
										s.FlexHStart(gui.Rigid(
											s.Text(b.Time.Format("15:04:05"),
												color, "Transparent",
												"Mono", "body2"),
										))
									}
								}),
								gui.Flexed(1, func() {
									tc := "DocText"
									if ww <= 480 {
										tc = color
									}
									s.FlexHStart(gui.Rigid(s.Text(b.Text, tc,
										"Transparent", "Mono", "body2"),
									))
								}),
								gui.Rigid(func() {
									if ww > 720 {
										s.FlexH(gui.Rigid(s.Text(b.Package,
											"PanelBg", "Transparent", "Primary",
											"body1"),
										))
									}
								}),
							)

						}, button)
						for button.Clicked(s.Gtx) {
							go func() {
								if s.Config.ClickCommand == "" {
									return
								}
								eb.Clicked = i
								split := strings.Split(b.CodeLocation, ":")
								v1 := split[0]
								v2 := split[1]
								c := strings.Replace(s.Config.ClickCommand, "$1", v1, 1)
								c = strings.Replace(c, "$2", v2, 1)
								Debug("running command", c)
								args := strings.Split(c, " ")
								cmd := exec.Command(args[0], args[1:]...)
								_ = cmd.Run()
								// s.Config.ClickCommand
							}()
						}
					}),
					gui.Rigid(func() {
						if ww > 640 {
							s.ButtonArea(func() {
								s.Inset(4, func() {
									s.Icon("HideItem", "PanelBg", "DocBg", 24)
								})
							}, hider)
							for hider.Clicked(s.Gtx) {
								s.Config.FilterNodes[eb.Get(i).
									Package].Hidden = true
							}
						}
					}),
				)
			})

		})
		s.W.Invalidate()
	})
}

func (s *State) RegenerateFilterBuf() {
	s.FilterBuf.Clear()
	// set all filters to not filter anything
	*s.Ctx.Config.LogLevel = logi.Trace
	for _, i := range s.Config.FilterNodes {
		i.Hidden = false
	}
	// regenerate the filter buffer
	for i := range s.EntryBuf.Buf {
		eb := s.EntryBuf.Buf[i]
		if s.FilterFunc(eb) {
			s.FilterBuf.Add(eb)
		}
	}
}
