package monitor

import (
	"gioui.org/layout"
)

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
			st.FlexV(Flexed(1, func() {
				st.Inset(8, func() {
					// cs := st.Gtx.Constraints
					// st.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
					st.BuildConfigPage()
				})
			}))
		}), Rigid(func() {
			st.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
			st.Inset(4, func() {})
		}),
		)
	})
}

func (st *State) BuildConfigPage() {
	st.FlexV(Rigid(func() {
		st.Inset(4, func() {
			st.FlexH(Rigid(func() {
				st.Inset(8,
					st.Text("Run in", "PanelText", "Primary", "h6"),
				)
			}), Rigid(func() {
				if st.RunningInRepo {
					fg, bg := "DocText", "DocBg"
					if st.Config.RunInRepo {
						fg, bg = "ButtonText", "ButtonBg"
					}
					st.TextButton("repo", "Primary", 16,
						fg, bg, st.RunningInRepoButton)
					for st.RunningInRepoButton.Clicked(st.Gtx) {
						st.Config.RunInRepo = true
						st.CannotRun = false
					}
				}
			}), Rigid(func() {
				fg, bg := "DocText", "DocBg"
				if !st.Config.RunInRepo {
					fg, bg = "ButtonText", "ButtonBg"
				}
				st.TextButton("profile", "Primary", 16,
					fg, bg, st.RunFromProfileButton)
				for st.RunFromProfileButton.Clicked(st.Gtx) {
					st.Config.RunInRepo = false
					st.CannotRun = false
				}
			}), Rigid(func() {
				txt := "run pod in its repository"
				if !st.Config.RunInRepo {
					txt = "not implemented"
					st.CannotRun = true
				}
				st.Inset(8,
					st.Text(txt, "PanelText", "Primary", "h6"),
				)
			}),
			)
		})
	}), Rigid(func() {
		st.Inset(4, func() {
			st.FlexH(Rigid(func() {
				st.Inset(8,
					st.Text("Use Go version", "PanelText", "Primary", "h6"),
				)
			}), Rigid(func() {
				if st.HasGo {
					fg, bg := "DocText", "DocBg"
					if st.Config.UseBuiltinGo {
						fg, bg = "ButtonText", "ButtonBg"
					}
					st.TextButton("builtin", "Primary", 16,
						fg, bg, st.UseBuiltinGoButton)
					for st.UseBuiltinGoButton.Clicked(st.Gtx) {
						st.Config.UseBuiltinGo = true
						st.CannotRun = false
						if !st.HasGo {
							st.CannotRun = true
						}
					}
				}
			}), Rigid(func() {
				fg, bg := "DocText", "DocBg"
				if !st.Config.UseBuiltinGo {
					fg, bg = "ButtonText", "ButtonBg"
				}
				st.TextButton("install new", "Primary", 16,
					fg, bg, st.InstallNewGoButton)
				for st.InstallNewGoButton.Clicked(st.Gtx) {
					st.Config.UseBuiltinGo = false
					st.CannotRun = false
					if !st.HasOtherGo {
						st.CannotRun = true
					}
				}
			}), Rigid(func() {
				txt := "build using built in go"
				if !st.Config.UseBuiltinGo {
					txt = "not implemented"
					st.CannotRun = true
				}
				st.Inset(8,
					st.Text(txt, "PanelText", "Primary", "h6"),
				)
			}),
			)
		})
	}),
	)
}
