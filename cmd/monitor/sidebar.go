package monitor

import (
	"gioui.org/layout"
)

func (s *State) Sidebar() layout.FlexChild {
	return Rigid(func() {
		if !(s.Config.BuildOpen || s.Config.SettingsOpen) {
			s.Gtx.Constraints.Width.Max /= 2
		} else {
			s.Gtx.Constraints.Width.Max -= 360
		}
		if s.Gtx.Constraints.Width.Max > 360 {
			s.Gtx.Constraints.Width.Max = 360
		}
		if s.Config.FilterOpen {
			s.FlexV(Rigid(func() {

				s.Gtx.Constraints.Height.Max = 48
				s.Gtx.Constraints.Height.Min = 48
				s.FlexH(
					Flexed(1, func() {
						cs := s.Gtx.Constraints
						s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg", "ff")
						if s.WindowWidth > 640 {
							s.Label("Filter")
						}
					}), Rigid(func() {
						s.IconButton("HideAll", "DocText", "DocBg", s.FilterHideButton)
						for s.FilterHideButton.Clicked(s.Gtx) {
							L.Debug("hide all")
							s.Loggers.CloseAllItems(s)
							s.SaveConfig()
						}
					}), Rigid(func() {
						s.IconButton("ShowAll", "DocText", "DocBg", s.FilterShowButton)
						for s.FilterShowButton.Clicked(s.Gtx) {
							L.Debug("show all")
							s.Loggers.OpenAllItems(s)
							s.SaveConfig()
						}
					}), Rigid(func() {
						s.IconButton("ShowItem", "DocText", "DocBg", s.FilterAllButton)
						for s.FilterAllButton.Clicked(s.Gtx) {
							L.Debug("filter all")
							s.Loggers.ShowAllItems(s)
							s.SaveConfig()
						}
					}), Rigid(func() {
						s.IconButton("HideItem", "DocText", "DocBg", s.FilterNoneButton)
						for s.FilterNoneButton.Clicked(s.Gtx) {
							L.Debug("filter none")
							s.Loggers.HideAllItems(s)
							s.SaveConfig()
						}
					}), Rigid(func() {
						s.IconButton("Filter", "DocText", "DocBg", s.FilterButton)
						for s.FilterButton.Clicked(s.Gtx) {
							L.Debug("filter header clicked")
							if !s.Config.FilterOpen {
								s.Config.BuildOpen = false
								s.Config.SettingsOpen = false
							}
							s.Config.FilterOpen = !s.Config.FilterOpen
							s.SaveConfig()
						}
					}),
				)
			}), Flexed(1, func() {
				s.Inset(16, func() {
					s.FlexV(
						Flexed(1, func() {
							s.Gtx.Constraints.Width.Min = 240
							cs := s.Gtx.Constraints
							s.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg", "ff")
							s.FlexV(Flexed(1, func() {
								s.Loggers.GetWidget(s)
							}))
						}),
					)
				})
			}),
			)
		}
	})
}
