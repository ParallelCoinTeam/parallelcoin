package monitor

import (
	"gioui.org/layout"
)

func (s *State) Sidebar() layout.FlexChild {
	return Rigid(func() {
		if !(s.Config.BuildOpen.Load() || s.Config.SettingsOpen.Load()) {
			s.Gtx.Constraints.Width.Max /= 2
		} else {
			s.Gtx.Constraints.Width.Max -= 360
		}
		if s.Gtx.Constraints.Width.Max > 360 {
			s.Gtx.Constraints.Width.Max = 360
		}
		//if s.Config.FilterOpen.Load() && (s.Config.BuildOpen.Load() ||
		//	s.Config.SettingsOpen.Load() && s.WindowWidth <= 800) {
		//	s.Config.FilterOpen.Store(false)
		//}
		if s.Config.FilterOpen.Load() {
			s.FlexV(Rigid(func() {
				s.Gtx.Constraints.Height.Max = 48
				s.Gtx.Constraints.Height.Min = 48
				s.FlexH(
					//Flexed(1, func() {
					//	cs := s.Gtx.Constraints
					//	s.Rectangle(cs.Width.Max, cs.Width.Max, "DocBg", "ff")
					//	if s.WindowWidth > 360 {
					//		//s.Gtx.Constraints.Width.Min = 32
					//		s.Icon("Filter", "DocText", "DocBg", 48)
					//	}
					//}),
					Flexed(1, func() {
						cs := s.Gtx.Constraints
						s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg", "ff")
						if s.WindowWidth > 640 {
							s.Label("Filter")
						}
					}), //Spacer(),
					Rigid(func() {
						s.IconButton("FilterAll", "DocText", "DocBg", s.FilterAllButton)
						for s.FilterAllButton.Clicked(s.Gtx) {
							L.Debug("filter all")
							//if !s.Config.FilterOpen.Load() {
							//	s.Config.BuildOpen.Store(false)
							//	s.Config.SettingsOpen.Store(false)
							//}
							//s.Config.FilterOpen.Toggle()
							//s.SaveConfig()
						}
					}), Rigid(func() {
						s.IconButton("FilterNone", "DocText", "DocBg", s.FilterNoneButton)
						for s.FilterNoneButton.Clicked(s.Gtx) {
							L.Debug("filter all")
							//if !s.Config.FilterOpen.Load() {
							//	s.Config.BuildOpen.Store(false)
							//	s.Config.SettingsOpen.Store(false)
							//}
							//s.Config.FilterOpen.Toggle()
							//s.SaveConfig()
						}
					}), Rigid(func() {
						s.IconButton("Filter", "DocText", "DocBg", s.FilterButton)
						for s.FilterButton.Clicked(s.Gtx) {
							L.Debug("filter header clicked")
							if !s.Config.FilterOpen.Load() {
								s.Config.BuildOpen.Store(false)
								s.Config.SettingsOpen.Store(false)
							}
							s.Config.FilterOpen.Toggle()
							s.SaveConfig()
						}
					}),
				)
			}), Flexed(1, func() {
				s.FlexH(
					Rigid(func() {
						s.Gtx.Constraints.Width.Min = 240
						cs := s.Gtx.Constraints
						s.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg", "ff")

					}),
				)
			}),
			)
		}
	})
}
