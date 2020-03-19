package monitor

import (
	"gioui.org/layout"
)

func (st *State) BottomBar() layout.FlexChild {
	return Rigid(func() {
		cs := st.Gtx.Constraints
		st.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg")
		st.FlexV(
			st.SettingsPage(),
			st.BuildPage(),
			st.StatusBar(),
		)
	})
}

func (st *State) StatusBar() layout.FlexChild {
	return Rigid(func() {
		cs := st.Gtx.Constraints
		st.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg")
		st.FlexH(
			st.RunControls(),
			st.RunmodeButtons(),
			st.BuildButtons(),
			st.SettingsButtons(),
		)
	})
}

func (st *State) RunControls() layout.FlexChild {
	return Rigid(func() {
		if st.CannotRun {
			return
		}
		if !st.Running {
			st.IconButton("Run", "PanelBg", "PanelText", st.RunMenuButton)
			for st.RunMenuButton.Clicked(st.Gtx) {
				L.Debug("clicked run button")
				if !st.Config.RunModeOpen {
					st.Running = true
				}
			}
		}
		if st.Running {
			ic := "Pause"
			fg, bg := "PanelBg", "PanelText"
			if st.Pausing {
				ic = "Run"
				fg, bg = "PanelText", "PanelBg"
			}
			st.FlexH(Rigid(func() {
				st.IconButton("Stop", "PanelBg", "PanelText",
					st.StopMenuButton)
				for st.StopMenuButton.Clicked(st.Gtx) {
					L.Debug("clicked stop button")
					st.Running = false
					st.Pausing = false
				}
			}), Rigid(func() {
				st.IconButton(ic, fg, bg, st.PauseMenuButton)
				for st.PauseMenuButton.Clicked(st.Gtx) {
					if st.Pausing {
						L.Debug("clicked on resume button")
					} else {
						L.Debug("clicked pause button")
					}
					Toggle(&st.Pausing)
				}
			}), Rigid(func() {
				st.IconButton("Restart", "PanelBg", "PanelText",
					st.RestartMenuButton)
				for st.RestartMenuButton.Clicked(st.Gtx) {
					L.Debug("clicked restart button")
				}
			}),
			)
		}
	})
}

func (st *State) RunmodeButtons() layout.FlexChild {
	return Rigid(func() {
		st.FlexH(Rigid(func() {
			fg, bg := "ButtonText", "ButtonBg"
			if st.Running {
				fg, bg = "DocBg", "DocText"
			}
			st.TextButton(st.Config.RunMode, "Secondary",
				23, fg, bg,
				st.RunModeFoldButton)
			for st.RunModeFoldButton.Clicked(st.Gtx) {
				if !st.Running {
					Toggle(&st.Config.RunModeOpen)
					st.SaveConfig()
				}
			}
		}), Rigid(func() {
			if st.Config.RunModeOpen {
				modes := []string{
					"node", "wallet", "shell", "gui",
				}
				st.ModesList.Layout(st.Gtx, len(modes), func(i int) {
					if st.Config.RunMode != modes[i] {
						st.TextButton(modes[i], "Secondary",
							23, "ButtonText",
							"ButtonBg", st.ModesButtons[modes[i]])
					}
					for st.ModesButtons[modes[i]].Clicked(st.Gtx) {
						L.Debug(modes[i], "clicked")
						if st.Config.RunModeOpen {
							st.Config.RunMode = modes[i]
							st.Config.RunModeOpen = false
						}
						st.SaveConfig()
					}
				})
			}
		}),
		)
	})
}
