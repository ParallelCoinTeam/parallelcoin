package monitor

import (
	"gioui.org/layout"
	"github.com/p9c/pod/pkg/gui"
)

func (s *State) BottomBar() layout.FlexChild {
	return gui.Rigid(func() {
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
	return gui.Rigid(func() {
		cs := s.Gtx.Constraints
		s.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg", "ff")
		s.FlexH(
			s.RunControls(),
			s.RunmodeButtons(),
			//s.Spacer("PanelBg"),
			gui.Flexed(1, func() {
				s.Gtx.Constraints.Height.Max = 48
				cs := s.Gtx.Constraints
				s.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg", "FF")
			}),
			s.BuildButtons(),
			s.SettingsButtons(),
			s.Filter(),
		)
	})
}

func (s *State) RunmodeButtons() layout.FlexChild {
	fg, bg := "ButtonText", "ButtonBg"
	return gui.Rigid(func() {
		s.FlexH(gui.Rigid(func() {
			if !s.Config.RunModeOpen {
				txt := s.Config.RunMode
				if s.Config.Running {
					cs := s.Gtx.Constraints
					bg, fg = "DocBg", "DocText"
					s.Rectangle(cs.Width.Min, 48, bg, "ff")
					s.Label(txt, fg, bg)
				} else {
					b := s.Buttons["RunModeFold"]
					s.TextButton(txt, "Secondary", 34, fg, bg, b)
					for b.Clicked(s.Gtx) {
						if !s.Config.Running {
							s.Config.RunModeOpen = true
							s.SaveConfig()
						}
					}
				}
			} else {
				modes := []string{
					"node", "wallet", "shell", "gui", "mon",
				}
				s.Lists["Modes"].Layout(s.Gtx, len(modes), func(i int) {
					mm := modes[i]
					fg := "DocBg"
					if modes[i] == s.Config.RunMode {
						fg = "DocText"
					}
					txt := mm
					if s.WindowWidth <= 880 && s.Config.FilterOpen ||
						s.WindowWidth <= 640 && !s.Config.FilterOpen {
						txt = txt[:1]
					}
					cs := s.Gtx.Constraints
					s.Rectangle(cs.Width.Max, cs.Height.Max, "ButtonBg",
						"ff")
					s.TextButton(txt, "Secondary", 34, fg,
						"ButtonBg", s.ModesButtons[modes[i]])
					for s.ModesButtons[modes[i]].Clicked(s.Gtx) {
						Debug(mm, "clicked")
						if s.Config.RunModeOpen {
							s.Config.RunMode = modes[i]
							s.Config.RunModeOpen = false
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
	fg, bg := "PanelText", "PanelBg"
	if s.Config.FilterOpen {
		fg, bg = "DocText", "DocBg"
	}
	return gui.Rigid(func() {
		b := s.Buttons["Filter"]
		s.ButtonArea(func() {
			cs := s.Gtx.Constraints
			s.Rectangle(cs.Width.Max, cs.Height.Max, bg, "ff")
			s.Inset(8, func() {
				s.Icon("Filter", fg, bg, 32)
			})
		}, b)
		for b.Clicked(s.Gtx) {
			Debug("clicked filter button")
			if !s.Config.FilterOpen {
				s.Config.SettingsOpen = false
				s.Config.BuildOpen = false
			}
			s.Config.FilterOpen = !s.Config.FilterOpen
			s.SaveConfig()
		}
		//}
	})
}
