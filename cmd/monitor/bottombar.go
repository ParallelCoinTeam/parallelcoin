package monitor

import (
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/p9c/pod/cmd/gui/pages"
	"github.com/p9c/pod/pkg/gelook"
)

func BottomBar(m *State) layout.FlexChild {
	return layout.Rigid(func() {
		cs := m.Gtx.Constraints
		gelook.DuoUIdrawRectangle(m.Gtx, cs.Width.Max,
			cs.Height.Max, m.Theme.Colors["Primary"],
			[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		layout.Flex{
			Axis: layout.Vertical,
		}.Layout(m.Gtx,
			StatusBar(m),
			SettingsPage(m),
			BuildPage(m),
		)
	})
}

func StatusBar(m *State) layout.FlexChild {
	return layout.Rigid(func() {
		cs := m.Gtx.Constraints
		gelook.DuoUIdrawRectangle(m.Gtx, cs.Width.Max,
			cs.Height.Max, m.Theme.Colors["Primary"],
			[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		layout.Flex{
			Axis: layout.Horizontal,
		}.Layout(m.Gtx,
			RunControls(m),
			RunmodeButtons(m),
			BuildButtons(m),
			SettingsButtons(m),
		)
	})
}

func RunControls(m *State) layout.FlexChild {
	return layout.Rigid(func() {
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
			layout.Flex{Axis: layout.Horizontal}.Layout(m.Gtx,
				layout.Rigid(func() {
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
				layout.Rigid(func() {
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
				layout.Rigid(func() {
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

func RunmodeButtons(m *State) layout.FlexChild {
	return layout.Rigid(func() {
		layout.Flex{
			Axis: layout.Horizontal,
		}.Layout(m.Gtx,
			layout.Rigid(func() {
				bg := m.Theme.Colors["Primary"]
				if m.RunModeOpen {
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
					m.RunModeOpen = !m.RunModeOpen

				}
				s.Layout(m.Gtx, m.RunModeFoldButton)
			}),
			layout.Rigid(func() {
				if m.RunModeOpen {
					modes := []string{
						"node", "wallet", "shell", "gui",
					}
					m.ModesList.Layout(m.Gtx, len(modes), func(i int) {
						fg, bg := m.Theme.Colors["Light"],
							m.Theme.Colors["Dark"]
						if m.RunMode == modes[i] {
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
							if m.RunModeOpen {
								m.RunMode = modes[i]
								m.RunModeOpen = false
							}
						}
					})
				} else {
					layout.UniformInset(unit.Dp(8)).Layout(m.Gtx, func() {
						t := m.Theme.DuoUIlabel(unit.Dp(18), m.RunMode)
						t.Font.Typeface = m.Theme.Fonts["Primary"]
						t.Layout(m.Gtx)
					})
				}
			}),
		)
	})
}

func SettingsButtons(m *State) layout.FlexChild {
	return layout.Flexed(1, func() {
		layout.Flex{
			Axis: layout.Horizontal,
		}.Layout(m.Gtx,
			layout.Rigid(func() {
				bg := m.Theme.Colors["Primary"]
				if m.SettingsOpen {
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
					case !m.SettingsOpen:
						m.BuildOpen = false
						m.SettingsOpen = true
					case m.SettingsOpen:
						m.SettingsOpen = false
					}
				}
				// s.Layout(m.Gtx, settingsFoldButton)
			}),
			layout.Rigid(func() {
				if m.WindowWidth > 1024 && m.SettingsOpen {
					pages.SettingsHeader(m.Rc, m.Gtx, m.Theme)()
				}
			}),
		)
	})
}

func SettingsPage(m *State) layout.FlexChild {
	if !m.SettingsOpen {
		return layout.Flexed(0, func() {})
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
	return layout.Flexed(weight, func() {
		cs := m.Gtx.Constraints
		gelook.DuoUIdrawRectangle(m.Gtx, cs.Width.Max,
			cs.Height.Max, m.Theme.Colors["Primary"],
			[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		layout.Flex{Axis: layout.Vertical}.Layout(m.Gtx,
			layout.Rigid(func() {
				if m.WindowWidth < 1024 {
					pages.SettingsHeader(m.Rc, m.Gtx, m.Theme)()
				}
			}),
			layout.Rigid(func() {
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

func BuildButtons(m *State) layout.FlexChild {
	return layout.Rigid(func() {
		layout.Flex{
			Axis: layout.Horizontal,
		}.Layout(m.Gtx,
			layout.Rigid(func() {
				bg := m.Theme.Colors["Primary"]
				if m.BuildOpen {
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
					case !m.BuildOpen:
						m.BuildOpen = true
						m.SettingsOpen = false
					case m.BuildOpen:
						m.BuildOpen = false
					}

				}
				s.Layout(m.Gtx, m.BuildFoldButton)
			}),
		)
	})
}

func BuildPage(m *State) layout.FlexChild {
	if !m.BuildOpen {
		return layout.Flexed(0, func() {})
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
	return layout.Flexed(weight, func() {
		cs := m.Gtx.Constraints
		gelook.DuoUIdrawRectangle(m.Gtx, cs.Width.Max,
			cs.Height.Max, m.Theme.Colors["Dark"],
			[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		// layout.Flex{Axis: layout.Vertical}.Layout(m.Gtx,
		// 	layout.Rigid(func() {
		// 		if m.WindowWidth < 1024 {
		// 			pages.SettingsHeader(rc, m.Gtx, theme)()
		// 		}
		// 	}),
		// 	layout.Rigid(func() {
		// 		cs := m.Gtx.Constraints
		// 		gelook.DuoUIdrawRectangle(m.Gtx, cs.Width.Max,
		// 			cs.Height.Max, m.Theme.Colors["Dark"],
		// 			[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		// 		layout.UniformInset(unit.Dp(settingsInset)).Layout(m.Gtx,
		// 			pages.SettingsBody(rc, m.Gtx, theme))
		// 	}),
		// )
	})
}
