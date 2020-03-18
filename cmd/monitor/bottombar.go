package monitor

import (
	"gioui.org/layout"

	"github.com/p9c/pod/cmd/gui/pages"
)

func (m *State) BottomBar() layout.FlexChild {
	return Rigid(func() {
		cs := m.Gtx.Constraints
		m.Rectangle(cs.Width.Max, cs.Height.Max, "Primary")
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
		m.Rectangle(cs.Width.Max, cs.Height.Max, "Primary")
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
			m.IconButton("Run", "Primary", m.RunMenuButton)
			for m.RunMenuButton.Clicked(m.Gtx) {
				L.Debug("clicked run button")
				m.Running = true
			}
		}
		if m.Running {
			m.FlexHorizontal(
				Rigid(func() {
					m.IconButton("Stop", "Dark", m.StopMenuButton)
					for m.StopMenuButton.Clicked(m.Gtx) {
						L.Debug("clicked stop button")
						m.Running = false
						m.Pausing = false
					}
				}),
				Rigid(func() {
					ic := "Pause"
					rc := "Dark"
					if m.Pausing {
						ic = "Run"
						rc = "Primary"
					}
					m.IconButton(ic, rc, m.PauseMenuButton)
					for m.PauseMenuButton.Clicked(m.Gtx) {
						if m.Pausing {
							L.Debug("clicked on resume button")
						} else {
							L.Debug("clicked pause button")
						}
						m.Pausing = !m.Pausing
					}
				}),
				Rigid(func() {
					m.IconButton("Restart", "Dark", m.RestartMenuButton)
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
		m.FlexHorizontal(
			Rigid(func() {
				fg, bg := "Light", "Primary"
				if m.Config.RunModeOpen {
					fg, bg = "Light", "Dark"
				}
				m.TextButton(m.Config.RunMode, "Secondary", 23, fg, bg,
					m.RunModeFoldButton)
				for m.RunModeFoldButton.Clicked(m.Gtx) {
					m.Config.RunModeOpen = !m.Config.RunModeOpen
					m.SaveConfig()
				}
			}),
			Rigid(func() {
				if m.Config.RunModeOpen {
					modes := []string{
						"node", "wallet", "shell", "gui",
					}
					m.ModesList.Layout(m.Gtx, len(modes), func(i int) {
						if m.Config.RunMode != modes[i] {
							m.TextButton(modes[i], "Primary", 16, "Light",
								"Dark", m.ModesButtons[modes[i]])
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

func (m *State) SettingsButtons() layout.FlexChild {
	return Flexed(1, func() {
		m.FlexHorizontal(
			Rigid(func() {
				bg := "Primary"
				if m.Config.SettingsOpen {
					bg = "Dark"
				}
				m.IconButton("settingsIcon", bg, m.SettingsFoldButton)
				for m.SettingsFoldButton.Clicked(m.Gtx) {
					L.Debug("settings folder clicked")
					switch {
					case !m.Config.SettingsOpen:
						m.Config.BuildOpen = false
						m.Config.SettingsOpen = true
					case m.Config.SettingsOpen:
						m.Config.SettingsOpen = false
					}
					m.SaveConfig()
				}
			}),
		)
	})
}

func (m *State) SettingsPage() layout.FlexChild {
	if !m.Config.SettingsOpen {
		return Flexed(0, func() {})
	}
	var weight float32 = 0.5
	var settingsInset = 0
	switch {
	case m.WindowWidth < 1024 && m.WindowHeight > 1024:
		// weight = 0.333
	case m.WindowHeight < 1024 && m.WindowWidth < 1024:
		weight = 1
	case m.WindowHeight < 600 && m.WindowWidth > 1024:
		weight = 1
	}
	return Flexed(weight, func() {
		cs := m.Gtx.Constraints
		m.Rectangle(cs.Width.Max, cs.Height.Max, "Secondary")
		m.FlexV(
			Rigid(func() {
				m.FlexHorizontal(
					Rigid(func() {
						m.TextButton("Run Settings", "Secondary",
							23, "PanelText", "Secondary",
							m.SettingsTitleCloseButton)
						for m.SettingsTitleCloseButton.Clicked(m.Gtx) {
							L.Debug("settings panel title close button clicked")
							m.Config.SettingsOpen = false
							m.SaveConfig()
						}
						// t := m.Theme.DuoUIlabel(unit.Dp(float32(
						// 	24)), "Run Settings")
						// t.Color = m.Theme.Colors["PanelText"]
						// t.Font.Typeface = m.Theme.Fonts["Secondary"]
						// layout.UniformInset(unit.Dp(8)).Layout(m.Gtx,
						// 	func() {
						// 		t.Layout(m.Gtx)
						// 	},
						// )
					}),
					Rigid(func() {
						// if m.WindowWidth < 1024 {
						pages.SettingsHeader(m.Rc, m.Gtx, m.Theme)()
						// }
					}),
					Flexed(1, func() {}),
					Rigid(func() {
						m.IconButton("minimize", "Secondary",
							m.SettingsCloseButton)
						for m.SettingsCloseButton.Clicked(m.Gtx) {
							L.Debug("settings panel close button clicked")
							m.Config.SettingsOpen = false
							m.SaveConfig()
						}
					}),
				)
			}),
			Rigid(func() {
				cs := m.Gtx.Constraints
				m.Rectangle(cs.Width.Max, cs.Height.Max, "Dark")
				m.Inset(settingsInset,
					pages.SettingsBody(m.Rc, m.Gtx, m.Theme))
			}),
		)
	})
}

func (m *State) BuildButtons() layout.FlexChild {
	return Rigid(func() {
		layout.Flex{
			Axis: layout.Horizontal,
		}.Layout(m.Gtx,
			Rigid(func() {
				bg := "Primary"
				if m.Config.BuildOpen {
					bg = "Dark"
				}
				m.TextButton("Build", "Secondary", 23,
					"PanelText", bg, m.BuildFoldButton)
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
		m.Rectangle(cs.Width.Max, cs.Height.Max, "Secondary")
		m.FlexV(
			Rigid(func() {
				m.FlexHorizontal(
					Rigid(func() {
						m.TextButton("Build Configuration", "Secondary",
							23, "PanelText", "Secondary",
							m.BuildTitleCloseButton)
						for m.BuildTitleCloseButton.Clicked(m.Gtx) {
							L.Debug("build configuration panel title close" +
								" button clicked")
							m.Config.BuildOpen = false
							m.SaveConfig()
						}
					}),
					Flexed(1, func() {}),
					Rigid(func() {
						m.IconButton("minimize", "Secondary",
							m.BuildCloseButton)
						for m.BuildCloseButton.Clicked(m.Gtx) {
							L.Debug("settings panel close button clicked")
							m.Config.BuildOpen = false
							m.SaveConfig()
						}
					}),
				)
			}),
			Flexed(1, func() {
				cs := m.Gtx.Constraints
				m.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg")
			}),
		)
	})
}
