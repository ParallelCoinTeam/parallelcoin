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

func (st *State) BuildButtons() layout.FlexChild {
	return Rigid(func() {
		st.FlexH(Rigid(func() {
			bg, fg := "PanelBg", "PanelText"
			if st.Config.BuildOpen {
				bg, fg = "DocBg", "DocText"
			}
			st.TextButton("Build", "Secondary", 23,
				fg, bg, st.BuildFoldButton)
			for st.BuildFoldButton.Clicked(st.Gtx) {
				L.Debug("run mode folder clicked")
				switch {
				case !st.Config.BuildOpen:
					st.Config.BuildOpen = true
					st.Config.SettingsOpen = false
				case st.Config.BuildOpen:
					st.Config.BuildOpen = false
				}
				st.SaveConfig()
			}
		}),
		)
	})
}

func (st *State) BuildPage() layout.FlexChild {
	if !st.Config.BuildOpen {
		return Flexed(0, func() {})
	}
	var weight float32 = 0.5
	switch {
	case st.WindowHeight < 1024 && st.WindowWidth < 1024:
		weight = 1
	case st.WindowHeight < 600 && st.WindowWidth > 1024:
		weight = 1
	}
	return Flexed(weight, func() {
		cs := st.Gtx.Constraints
		st.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
		st.FlexV(Rigid(func() {
			st.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
			st.Inset(4, func() {})
		}), Rigid(func() {

			st.FlexH(Rigid(func() {
				st.TextButton("Build Configuration", "Secondary",
					23, "DocText", "DocBg",
					st.BuildTitleCloseButton)
				for st.BuildTitleCloseButton.Clicked(st.Gtx) {
					L.Debug("build configuration panel title close" +
						" button clicked")
					st.Config.BuildOpen = false
					st.SaveConfig()
				}
			}), Spacer(), Rigid(func() {
				st.IconButton("minimize", "DocText", "DocBg",
					st.BuildCloseButton)
				for st.BuildCloseButton.Clicked(st.Gtx) {
					L.Debug("settings panel close button clicked")
					st.Config.BuildOpen = false
					st.SaveConfig()
				}
			}),
			)
		}), Flexed(1, func() {
			cs := st.Gtx.Constraints
			st.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg")
			// st.Inset(8, func() {			})
		}), Rigid(func() {
			st.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
			st.Inset(4, func() {})
		}),
		)
	})
}
