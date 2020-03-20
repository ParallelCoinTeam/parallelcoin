package monitor

import (
	"gioui.org/layout"
)

func (s *State) BottomBar() layout.FlexChild {
	return Rigid(func() {
		cs := s.Gtx.Constraints
		s.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg")
		s.FlexV(
			s.SettingsPage(),
			s.BuildPage(),
			s.StatusBar(),
		)
	})
}

func (s *State) StatusBar() layout.FlexChild {
	return Rigid(func() {
		cs := s.Gtx.Constraints
		s.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg")
		s.FlexH(
			s.RunControls(),
			s.RunmodeButtons(),
			s.BuildButtons(),
			s.SettingsButtons(),
		)
	})
}

func (s *State) RunmodeButtons() layout.FlexChild {
	return Rigid(func() {
		s.FlexH(Rigid(func() {
			if !s.Config.RunModeOpen.Load() {
				fg, bg := "ButtonText", "ButtonBg"
				if s.Running.Load() {
					fg, bg = "DocBg", "DocText"
				}
				s.TextButton(s.Config.RunMode.Load(), "Secondary",
					23, fg, bg,
					s.RunModeFoldButton)
				for s.RunModeFoldButton.Clicked(s.Gtx) {
					if !s.Running.Load() {
						s.Config.RunModeOpen.Store(true)
						s.SaveConfig()
					}
				}
			} else {
				modes := []string{
					"node", "wallet", "shell", "gui", "monitor",
				}
				s.ModesList.Layout(s.Gtx, len(modes), func(i int) {
					// if s.Config.RunMode.Load() != modes[i] {
					s.TextButton(modes[i], "Secondary",
						23, "ButtonText",
						"ButtonBg", s.ModesButtons[modes[i]])
					// }
					for s.ModesButtons[modes[i]].Clicked(s.Gtx) {
						L.Debug(modes[i], "clicked")
						if s.Config.RunModeOpen.Load() {
							s.Config.RunMode.Store(modes[i])
							s.Config.RunModeOpen.Store(false)
						}
						s.SaveConfig()
					}
				})
			}
		}),
		)
	})
}
