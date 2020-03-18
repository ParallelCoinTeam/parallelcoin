package monitor

import (
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/p9c/pod/cmd/gui/pages"
	"github.com/p9c/pod/pkg/gelook"
)

func (m *State) BottomBar() layout.FlexChild {
	return Rigid(func() {
		cs := m.Gtx.Constraints
		gelook.DuoUIdrawRectangle(m.Gtx, cs.Width.Max,
			cs.Height.Max, m.Theme.Colors["Primary"],
			[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		m.FlexVertical(
			m.SettingsPage(),
			m.BuildPage(),
			m.StatusBar(),
		)
	})
}

func (m *State) StatusBar() layout.FlexChild {
	return Rigid(func() {
		cs := m.Gtx.Constraints
		gelook.DuoUIdrawRectangle(m.Gtx, cs.Width.Max,
			cs.Height.Max, m.Theme.Colors["Primary"],
			[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		m.FlexHorizontal(
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
			m.Theme.DuoUIbutton("", "", "",
				m.Theme.Colors["Primary"], "",
				m.Theme.Colors["Dark"], "Run",
				m.Theme.Colors["Light"], 0, 41, 41, 41,
				0, 0).IconLayout(m.Gtx, m.RunMenuButton)
			for m.RunMenuButton.Clicked(m.Gtx) {
				L.Debug("clicked run button")
				m.Running = true
			}
		}
		if m.Running {
			m.FlexHorizontal(
				Rigid(func() {
					m.Theme.DuoUIbutton("", "", "",
						m.Theme.Colors["Dark"], "",
						m.Theme.Colors["Dark"], "Stop",
						m.Theme.Colors["Light"], 0, 41, 41, 41,
						0, 0).IconLayout(m.Gtx, m.StopMenuButton)
					for m.StopMenuButton.Clicked(m.Gtx) {
						L.Debug("clicked stop button")
						m.Running = false
						m.Pausing = false
					}
				}),
				Rigid(func() {
					ic := "Pause"
					rc := m.Theme.Colors["Dark"]
					if m.Pausing {
						ic = "Run"
						rc = m.Theme.Colors["Primary"]
					}
					m.Theme.DuoUIbutton("", "", "",
						rc, "",
						m.Theme.Colors["Dark"], ic,
						m.Theme.Colors["Light"], 0, 41, 41, 41,
						0, 0).IconLayout(m.Gtx, m.PauseMenuButton)
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
					m.Theme.DuoUIbutton("", "", "",
						m.Theme.Colors["Dark"], "",
						m.Theme.Colors["Dark"], "Restart",
						m.Theme.Colors["Light"], 0, 41, 41, 41,
						0, 0).IconLayout(m.Gtx, m.RestartMenuButton)
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
				bg := m.Theme.Colors["Primary"]
				if m.Config.RunModeOpen {
					bg = m.Theme.Colors["Dark"]
				}
				s := m.Theme.DuoUIbutton(
					m.Theme.Fonts["Secondary"],
					"mode",
					m.Theme.Colors["Light"],
					bg,
					m.Theme.Colors["Primary"],
					m.Theme.Colors["Light"],
					"settingsIcon",
					m.Theme.Colors["Light"],
					23, 23, 23, 23, 0, 0)
				for m.RunModeFoldButton.Clicked(m.Gtx) {
					L.Debug("run mode folder clicked")
					m.Config.RunModeOpen = !m.Config.RunModeOpen
					m.SaveConfig()
				}
				s.Layout(m.Gtx, m.RunModeFoldButton)
			}),
			Rigid(func() {
				if m.Config.RunModeOpen {
					modes := []string{
						"node", "wallet", "shell", "gui",
					}
					m.ModesList.Layout(m.Gtx, len(modes), func(i int) {
						fg, bg := m.Theme.Colors["Light"],
							m.Theme.Colors["Dark"]
						if m.Config.RunMode == modes[i] {
							fg, bg = m.Theme.Colors["Dark"],
								m.Theme.Colors["Light"]
						}
						m.Theme.DuoUIbutton(m.Theme.Fonts["Primary"],
							modes[i],
							fg,
							bg,
							"", "",
							"", "",
							16, 0, 80, 32, 4, 4).Layout(m.Gtx,
							m.ModesButtons[modes[i]])
						for m.ModesButtons[modes[i]].Clicked(m.Gtx) {
							L.Debug(modes[i], "clicked")
							if m.Config.RunModeOpen {
								m.Config.RunMode = modes[i]
								m.Config.RunModeOpen = false
							}
							m.SaveConfig()
						}
					})
				} else {
					layout.UniformInset(unit.Dp(8)).Layout(m.Gtx, func() {
						t := m.Theme.DuoUIlabel(unit.Dp(18), m.Config.RunMode)
						t.Font.Typeface = m.Theme.Fonts["Primary"]
						t.Layout(m.Gtx)
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
				bg := m.Theme.Colors["Primary"]
				if m.Config.SettingsOpen {
					bg = m.Theme.Colors["Dark"]
				}
				m.Theme.DuoUIbutton(m.Theme.Fonts["Primary"], "settings",
					m.Theme.Colors["Light"],
					bg, "",
					m.Theme.Colors["Dark"], "settingsIcon",
					m.Theme.Colors["Light"], 23, 32, 41, 41,
					0, 0).IconLayout(m.Gtx, m.SettingsFoldButton)
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
				// s.Layout(m.Gtx, settingsFoldButton)
			}),
			// Rigid(func() {
			// 	if m.WindowWidth > 1024 && m.Config.SettingsOpen {
			// 		pages.SettingsHeader(m.Rc, m.Gtx, m.Theme)()
			// 	}
			// }),
		)
	})
}

func (m *State) SettingsPage() layout.FlexChild {
	if !m.Config.SettingsOpen {
		return Flexed(0, func() {})
	}
	var weight float32 = 0.5
	var settingsInset float32 = 0
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
		gelook.DuoUIdrawRectangle(m.Gtx, cs.Width.Max,
			cs.Height.Max, m.Theme.Colors["Secondary"],
			[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		m.FlexVertical(
			Rigid(func() {
				m.FlexHorizontal(
					Rigid(func() {
						t := m.Theme.DuoUIlabel(unit.Dp(float32(
							24)), "Run Settings")
						t.Color = m.Theme.Colors["PanelText"]
						t.Font.Typeface = m.Theme.Fonts["Secondary"]
						layout.UniformInset(unit.Dp(8)).Layout(m.Gtx,
							func() {
								t.Layout(m.Gtx)
							},
						)
					}),
					Rigid(func() {
						// if m.WindowWidth < 1024 {
						pages.SettingsHeader(m.Rc, m.Gtx, m.Theme)()
						// }
					}),
					Flexed(1, func() {}),
					Rigid(func() {
						m.Theme.DuoUIbutton("", "settings",
							m.Theme.Colors["PanelText"],
							"", "",
							m.Theme.Colors["PanelBg"], "minimize",
							m.Theme.Colors["PanelText"],
							0, 32, 41, 41,
							0, 0).IconLayout(m.Gtx, m.SettingsCloseButton)
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
				gelook.DuoUIdrawRectangle(m.Gtx, cs.Width.Max,
					cs.Height.Max, m.Theme.Colors["Dark"],
					[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				layout.UniformInset(unit.Dp(settingsInset)).Layout(m.Gtx,
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
				bg := m.Theme.Colors["Primary"]
				if m.Config.BuildOpen {
					bg = m.Theme.Colors["Dark"]
				}
				s := m.Theme.DuoUIbutton(
					m.Theme.Fonts["Secondary"],
					"build",
					m.Theme.Colors["Light"],
					bg,
					m.Theme.Colors["Primary"],
					m.Theme.Colors["Light"],
					"settingsIcon",
					m.Theme.Colors["Light"],
					23, 23, 23, 23, 0, 0)
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
				s.Layout(m.Gtx, m.BuildFoldButton)
			}),
		)
	})
}

func (m *State) BuildPage() layout.FlexChild {
	if !m.Config.BuildOpen {
		return Flexed(0, func() {})
	}
	var weight float32 = 0.5
	// var settingsInset float32 = 0
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
		gelook.DuoUIdrawRectangle(m.Gtx, cs.Width.Max,
			cs.Height.Max, m.Theme.Colors["Secondary"],
			[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		m.FlexVertical(
			Rigid(func() {
				m.FlexHorizontal(
					Rigid(func() {
						t := m.Theme.DuoUIlabel(unit.Dp(float32(
							24)), "Build Configuration")
						t.Color = m.Theme.Colors["PanelText"]
						t.Font.Typeface = m.Theme.Fonts["Secondary"]
						layout.UniformInset(unit.Dp(8)).Layout(m.Gtx,
							func() {
								t.Layout(m.Gtx)
							},
						)
					}),
					Flexed(1, func() {}),
					Rigid(func() {
						m.Theme.DuoUIbutton("", "settings",
							m.Theme.Colors["PanelText"],
							"", "",
							m.Theme.Colors["PanelBg"], "minimize",
							m.Theme.Colors["PanelText"],
							0, 32, 41, 41,
							0, 0).IconLayout(m.Gtx, m.BuildCloseButton)
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
				gelook.DuoUIdrawRectangle(m.Gtx, cs.Width.Max,
					cs.Height.Max, m.Theme.Colors["PanelBg"],
					[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			}),
		)
	})
}
