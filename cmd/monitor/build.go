package monitor

import (
	"gioui.org/layout"
)

func (s *State) BuildButtons() layout.FlexChild {
	return Rigid(func() {
		s.FlexH(Rigid(func() {
			bg, fg := "PanelBg", "PanelText"
			if s.Config.BuildOpen.Load() {
				bg, fg = "DocBg", "DocText"
			}
			s.TextButton("Build", "Secondary", 23,
				fg, bg, s.BuildFoldButton)
			for s.BuildFoldButton.Clicked(s.Gtx) {
				L.Debug("run mode folder clicked")
				switch {
				case !s.Config.BuildOpen.Load():
					s.Config.BuildOpen.Store(true)
					s.Config.SettingsOpen.Store(false)
				case s.Config.BuildOpen.Load():
					s.Config.BuildOpen.Store(false)
				}
				s.SaveConfig()
			}
		}),
		)
	})
}

func (s *State) BuildPage() layout.FlexChild {
	if !s.Config.BuildOpen.Load() {
		return Flexed(0, func() {})
	}
	var weight float32 = 0.5
	switch {
	case s.Config.BuildZoomed.Load():
		weight = 1
	case s.WindowHeight < 1024 && s.WindowWidth < 1024:
		weight = 1
	case s.WindowHeight < 600 && s.WindowWidth > 1024:
		weight = 1
	}
	return Flexed(weight, func() {
		cs := s.Gtx.Constraints
		s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
		s.FlexV(Rigid(func() {
			s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
			s.Inset(4, func() {})
		}), Rigid(func() {
			s.FlexH(Rigid(func() {
				s.TextButton("Build Configuration", "Secondary",
					23, "DocText", "DocBg",
					s.BuildTitleCloseButton)
				for s.BuildTitleCloseButton.Clicked(s.Gtx) {
					L.Debug("build configuration panel title close" +
						" button clicked")
					s.Config.BuildOpen.Store(false)
					s.SaveConfig()
				}
			}), Spacer(), Rigid(func() {
				if !(s.WindowHeight < 1024 && s.WindowWidth < 1024 ||
					s.WindowHeight < 600 && s.WindowWidth > 1024) {
					ic := "zoom"
					if s.Config.BuildZoomed.Load() {
						ic = "minimize"
					}
					s.IconButton(ic, "DocText", "DocBg",
						s.BuildZoomButton)
					for s.BuildZoomButton.Clicked(s.Gtx) {
						L.Debug("settings panel fold button clicked")
						s.Config.BuildZoomed.Toggle()
						s.SaveConfig()
					}
				}
			}), Spacer(), Rigid(func() {
				s.IconButton("foldIn", "DocText", "DocBg",
					s.BuildCloseButton)
				for s.BuildCloseButton.Clicked(s.Gtx) {
					L.Debug("settings panel close button clicked")
					s.Config.BuildOpen.Store(false)
					s.SaveConfig()
				}
			}),
			)
		}), Flexed(1, func() {
			cs := s.Gtx.Constraints
			s.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg")
			s.FlexV(Flexed(1, func() {
				s.Inset(8, func() {
					// cs := s.Gtx.Constraints
					// s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
					s.BuildConfigPage()
				})
			}))
		}), Rigid(func() {
			s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
			s.Inset(4, func() {})
		}),
		)
	})
}

func (s *State) BuildConfigPage() {
	s.FlexV(Rigid(func() {
		s.Inset(4, func() {
			s.FlexH(Rigid(func() {
				s.Inset(8,
					s.Text("Run in", "PanelText", "Primary", "h6"),
				)
			}), Rigid(func() {
				if s.RunningInRepo {
					fg, bg := "DocText", "DocBg"
					if s.Config.RunInRepo.Load() {
						fg, bg = "ButtonText", "ButtonBg"
					}
					s.TextButton("repo", "Primary", 16,
						fg, bg, s.RunningInRepoButton)
					for s.RunningInRepoButton.Clicked(s.Gtx) {
						if !s.Config.Running.Load() {
							s.Config.RunInRepo.Store(true)
							s.CannotRun = false
							s.SaveConfig()
						}
					}
				}
			}), Rigid(func() {
				fg, bg := "DocText", "DocBg"
				if !s.Config.RunInRepo.Load() {
					fg, bg = "ButtonText", "ButtonBg"
				}
				s.TextButton("profile", "Primary", 16,
					fg, bg, s.RunFromProfileButton)
				for s.RunFromProfileButton.Clicked(s.Gtx) {
					if !s.Config.Running.Load() {
						s.Config.RunInRepo.Store(false)
						s.CannotRun = false
						s.SaveConfig()
					}
				}
			}), Rigid(func() {
				txt := "run pod in its repository"
				if !s.Config.RunInRepo.Load() {
					txt = "not implemented"
					s.CannotRun = true
				}
				s.Inset(8,
					s.Text(txt, "PanelText", "Primary", "h6"),
				)
			}),
			)
		})
	}), Rigid(func() {
		s.Inset(4, func() {
			s.FlexH(Rigid(func() {
				s.Inset(8,
					s.Text("Use Go version", "PanelText", "Primary", "h6"),
				)
			}), Rigid(func() {
				if s.HasGo {
					fg, bg := "DocText", "DocBg"
					if s.Config.UseBuiltinGo.Load() {
						fg, bg = "ButtonText", "ButtonBg"
					}
					s.TextButton("builtin", "Primary", 16,
						fg, bg, s.UseBuiltinGoButton)
					for s.UseBuiltinGoButton.Clicked(s.Gtx) {
						if !s.Config.RunInRepo.Load() {
							s.Config.UseBuiltinGo.Store(true)
							s.CannotRun = false
							if !s.HasGo {
								s.CannotRun = true
							}
						}
					}
				}
			}), Rigid(func() {
				fg, bg := "DocText", "DocBg"
				if !s.Config.UseBuiltinGo.Load() {
					fg, bg = "ButtonText", "ButtonBg"
				}
				s.TextButton("install new", "Primary", 16,
					fg, bg, s.InstallNewGoButton)
				for s.InstallNewGoButton.Clicked(s.Gtx) {
					if !s.Config.RunInRepo.Load() {
						s.Config.UseBuiltinGo.Store(false)
						s.CannotRun = false
						if !s.HasOtherGo {
							s.CannotRun = true
						}
					}
				}
			}), Rigid(func() {
				txt := "build using built in go"
				if !s.Config.UseBuiltinGo.Load() {
					txt = "not implemented"
					s.CannotRun = true
				}
				s.Inset(8,
					s.Text(txt, "PanelText", "Primary", "h6"),
				)
			}),
			)
		})
	}),
	)
}
