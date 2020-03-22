package monitor

import (
	"gioui.org/layout"
)

func (s *State) BottomBar() layout.FlexChild {
	return Rigid(func() {
		cs := s.Gtx.Constraints
		s.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg", "ff")
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
		s.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg", "ff")
		s.FlexH(
			s.RunControls(),
			s.RunmodeButtons(),
			s.BuildButtons(),
			s.SettingsButtons(),
			Spacer(),
			s.Filter(),
		)
	})
}

func (s *State) RunmodeButtons() layout.FlexChild {
	return Rigid(func() {
		s.FlexH(Rigid(func() {
			if !s.Config.RunModeOpen.Load() {
				fg, bg := "ButtonText", "ButtonBg"
				if s.Config.Running.Load() {
					fg, bg = "DocBg", "DocText"
				}
				txt := s.Config.RunMode.Load()
				cs := s.Gtx.Constraints
				if cs.Width.Max <= 240 {
					txt = txt[:1]
				}
				s.TextButton(txt, "Secondary",
					32, fg, bg,
					s.RunModeFoldButton)
				for s.RunModeFoldButton.Clicked(s.Gtx) {
					if !s.Config.Running.Load() {
						s.Config.RunModeOpen.Store(true)
						s.SaveConfig()
					}
				}
			} else {
				modes := []string{
					"node", "wallet", "shell", "gui", "mon",
				}
				s.ModesList.Layout(s.Gtx, len(modes), func(i int) {
					mm := modes[i]
					txt := mm
					if s.WindowWidth <= 720 && s.Config.FilterOpen.Load() {
						txt = txt[:1]
					}
					cs := s.Gtx.Constraints
					s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg", "ff")
					s.TextButton(txt, "Secondary",
						32, "DocText",
						"DocBg", s.ModesButtons[mm])
					for s.ModesButtons[mm].Clicked(s.Gtx) {
						L.Debug(mm, "clicked")
						if s.Config.RunModeOpen.Load() {
							s.Config.RunMode.Store(mm)
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

func (s *State) Filter() layout.FlexChild {
	return Rigid(func() {
		fg, bg := "PanelText", "PanelBg"
		if s.Config.FilterOpen.Load() {
			fg, bg = "DocText", "DocBg"
		}
		if !(s.Config.FilterOpen.Load() && s.WindowWidth <= 720) ||
			(!s.Config.FilterOpen.Load() && s.WindowWidth > 480) {
			s.IconButton("Filter", fg, bg, s.FilterButton)
			for s.FilterButton.Clicked(s.Gtx) {
				L.Debug("clicked filter button")
				if !s.Config.FilterOpen.Load() {
					if s.WindowWidth < 1024 {
						s.Config.SettingsOpen.Store(false)
						s.Config.BuildOpen.Store(false)
					}
				}
				s.Config.FilterOpen.Toggle()
				s.SaveConfig()
			}
		}
	})
}
