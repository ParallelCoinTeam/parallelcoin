package monitor

import "gioui.org/layout"

func (s *State) Sidebar() layout.FlexChild {
	return Rigid(func() {
		if !(s.Config.BuildOpen.Load() || s.Config.SettingsOpen.Load()) {
			s.Gtx.Constraints.Width.Max /= 2
			//s.Gtx.Constraints.Width.Max -= 64
			//if s.Gtx.Constraints.Width.Max > 240 {
			//	s.Gtx.Constraints.Width.Max = 240
			//}
		} else {
			s.Gtx.Constraints.Width.Max -= 480
		}
		if s.Gtx.Constraints.Width.Max > 480 {
			s.Gtx.Constraints.Width.Max = 480
		}
		if s.Config.FilterOpen.Load() && (s.Config.BuildOpen.Load() ||
			s.Config.SettingsOpen.Load() && s.WindowWidth <= 1024) {
			s.Config.FilterOpen.Store(false)
		}
		if s.Config.FilterOpen.Load() {
			s.FlexV(
				Rigid(func() {
					//s.Gtx.Constraints.Width.Min = 256
					s.Gtx.Constraints.Height.Max = 64
					s.Gtx.Constraints.Height.Min = 64
					cs := s.Gtx.Constraints
					s.Rectangle(cs.Width.Max, cs.Width.Max, "DocBg", "ff")
					s.FlexH(
						Flexed(1, func() {
							s.TextButton("Filter", "Secondary",
								16, "DocText", "DocBg", s.FilterHeaderButton)
							for s.FilterHeaderButton.Clicked(s.Gtx) {
								L.Debug("filter header clicked")
								if !s.Config.FilterOpen.Load() {
									s.Config.BuildOpen.Store(false)
									s.Config.SettingsOpen.Store(false)
								}
								s.Config.FilterOpen.Toggle()
								s.SaveConfig()
							}
						}),
						//Rigid(func(){
						//	s.IconButton("foldIn", "DocText", "DocBg",
						//		s.FilterClearButton)
						//}),
					)
				}),
				Flexed(1, func() {
					s.FlexH(
						//Rigid(func() {
						//	s.Gtx.Constraints.Width.Max = 8
						//	s.Gtx.Constraints.Width.Min = 8
						//	cs := s.Gtx.Constraints
						//	s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg", "ff")
						//}),
						Rigid(func() {
							//s.Gtx.Constraints.Width.Max = 240
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
