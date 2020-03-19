package monitor

import (
	"fmt"

	"gioui.org/layout"
)

func (m *State) BottomBar() layout.FlexChild {
	return Rigid(func() {
		cs := m.Gtx.Constraints
		m.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg")
		m.FlexV(
			m.SettingsPage(),
			m.BuildPage(),
			m.StatusBar(),
		)
	})
}

func (m *State) StatusBar() layout.FlexChild {
	return Rigid(func() {
		cs := m.Gtx.Constraints
		m.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg")
		m.FlexH(
			m.RunControls(),
			m.RunmodeButtons(),
			m.BuildButtons(),
			m.SettingsButtons(),
		)
	})
}

func (m *State) RunControls() layout.FlexChild {
	return Rigid(func() {
		if !m.Running {
			m.IconButton("Run", "PanelBg", "PanelText", m.RunMenuButton)
			for m.RunMenuButton.Clicked(m.Gtx) {
				L.Debug("clicked run button")
				m.Running = true
			}
		}
		if m.Running {
			ic := "Pause"
			fg, bg := "PanelBg", "PanelText"
			if m.Pausing {
				ic = "Run"
				fg, bg = "PanelText", "PanelBg"
			}
			m.FlexH(Rigid(func() {
				m.IconButton("Stop", "PanelBg", "PanelText",
					m.StopMenuButton)
				for m.StopMenuButton.Clicked(m.Gtx) {
					L.Debug("clicked stop button")
					m.Running = false
					m.Pausing = false
				}
			}), Rigid(func() {
				m.IconButton(ic, fg, bg, m.PauseMenuButton)
				for m.PauseMenuButton.Clicked(m.Gtx) {
					if m.Pausing {
						L.Debug("clicked on resume button")
					} else {
						L.Debug("clicked pause button")
					}
					Toggle(&m.Pausing)
				}
			}), Rigid(func() {
				m.IconButton("Restart", "PanelBg", "PanelText",
					m.RestartMenuButton)
				for m.RestartMenuButton.Clicked(m.Gtx) {
					L.Debug("clicked restart button")
				}
			}),
			)
		}
	})
}

func (m *State) RunmodeButtons() layout.FlexChild {
	return Rigid(func() {
		m.FlexH(Rigid(func() {
			fg, bg := "ButtonText", "ButtonBg"
			if m.Config.RunModeOpen {
				fg, bg = "ButtonText", "ButtonBg"
			}
			m.TextButton(m.Config.RunMode, "Secondary",
				23, fg, bg,
				m.RunModeFoldButton)
			for m.RunModeFoldButton.Clicked(m.Gtx) {
				if !m.Running {
					Toggle(&m.Config.RunModeOpen)
					m.SaveConfig()
				}
			}
		}), Rigid(func() {
			if m.Config.RunModeOpen {
				modes := []string{
					"node", "wallet", "shell", "gui",
				}
				m.ModesList.Layout(m.Gtx, len(modes), func(i int) {
					if m.Config.RunMode != modes[i] {
						m.TextButton(modes[i], "Secondary",
							23, "ButtonText",
							"ButtonBg", m.ModesButtons[modes[i]])
					}
					for m.ModesButtons[modes[i]].Clicked(m.Gtx) {
						L.Debug(modes[i], "clicked")
						if m.Config.RunModeOpen {
							m.Config.RunMode = modes[i]
							m.Config.RunModeOpen = false
						}
						m.SaveConfig()
					}
				})
			}
		}),
		)
	})
}

func (m *State) BuildButtons() layout.FlexChild {
	return Rigid(func() {
		m.FlexH(Rigid(func() {
			bg, fg := "PanelBg", "PanelText"
			if m.Config.BuildOpen {
				bg, fg = "DocBg", "DocText"
			}
			m.TextButton("Build", "Secondary", 23,
				fg, bg, m.BuildFoldButton)
			for m.BuildFoldButton.Clicked(m.Gtx) {
				L.Debug("run mode folder clicked")
				switch {
				case !m.Config.BuildOpen:
					m.Config.BuildOpen = true
					m.Config.SettingsOpen = false
				case m.Config.BuildOpen:
					m.Config.BuildOpen = false
				}
				m.SaveConfig()
			}
		}),
		)
	})
}

func (m *State) BuildPage() layout.FlexChild {
	if !m.Config.BuildOpen {
		return Flexed(0, func() {})
	}
	var weight float32 = 0.5
	switch {
	case m.WindowHeight < 1024 && m.WindowWidth < 1024:
		weight = 1
	case m.WindowHeight < 600 && m.WindowWidth > 1024:
		weight = 1
	}
	return Flexed(weight, func() {
		cs := m.Gtx.Constraints
		m.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
		m.FlexV(Rigid(func() {
			m.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
			m.Inset(4, func() {})
		}), Rigid(func() {

			m.FlexH(Rigid(func() {
				m.TextButton("Build Configuration", "Secondary",
					23, "DocText", "DocBg",
					m.BuildTitleCloseButton)
				for m.BuildTitleCloseButton.Clicked(m.Gtx) {
					L.Debug("build configuration panel title close" +
						" button clicked")
					m.Config.BuildOpen = false
					m.SaveConfig()
				}
			}), Spacer(), Rigid(func() {
				m.IconButton("minimize", "DocText", "DocBg",
					m.BuildCloseButton)
				for m.BuildCloseButton.Clicked(m.Gtx) {
					L.Debug("settings panel close button clicked")
					m.Config.BuildOpen = false
					m.SaveConfig()
				}
			}),
			)
		}), Flexed(1, func() {
			cs := m.Gtx.Constraints
			m.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg")
			// m.Inset(8, func() {			})
		}), Rigid(func() {
			m.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
			m.Inset(4, func() {})
		}),
		)
	})
}

func (m *State) SettingsHeader() func() {
	return func() {
		layout.Flex{Spacing: layout.SpaceBetween}.Layout(m.Gtx,
			layout.Rigid(func() {
				m.Inset(0, m.SettingsTabs)
			}),
		)
	}
}
func (m *State) SettingsTabs() {
	groupsNumber := len(m.Rc.Settings.Daemon.Schema.Groups)
	m.GroupsList.Layout(m.Gtx, groupsNumber, func(i int) {
		color :=
			m.Theme.Colors["DocText"]
		bgColor :=
			m.Theme.Colors["DocBg"]
		i = groupsNumber - 1 - i
		t := m.Rc.Settings.Daemon.Schema.Groups[i]
		txt := fmt.Sprint(t.Legend)
		for m.Rc.Settings.Tabs.TabsList[txt].Clicked(m.Gtx) {
			m.Rc.Settings.Tabs.Current = txt
		}
		if m.Rc.Settings.Tabs.Current == txt {
			color =
				m.Theme.Colors["PanelText"]
			bgColor =
				m.Theme.Colors["PanelBg"]
		}
		m.Theme.DuoUIbutton(m.Theme.Fonts["Primary"],
			txt, color, bgColor, "", "", "", "",
			16, 0, 80, 32, 0, 0).Layout(m.Gtx,
			m.Rc.Settings.Tabs.TabsList[txt])
	})
}
