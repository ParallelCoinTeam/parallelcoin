package monitor

import (
	"gioui.org/layout"
	"github.com/p9c/pod/pkg/gui"
	"github.com/p9c/pod/pkg/gui/gel"
)

func (s *State) BuildButtons(headless bool) layout.FlexChild {
	gtx := s.Gtx
	if headless {
		gtx = s.Htx
	}
	return gui.Rigid(func() {
		if s.WindowWidth >= 360 || !s.Config.FilterOpen {
			s.FlexH(gui.Rigid(func() {
				bg, fg := "PanelBg", "PanelText"
				if s.Config.BuildOpen {
					bg, fg = "DocBg", "DocText"
				}
				b := s.Buttons["BuildFold"]
				//s.IconButton("Build", fg, bg, b)
				s.ButtonArea(func() {
					gtx.Constraints.Width.Max = 48
					gtx.Constraints.Height.Max = 48
					cs := gtx.Constraints
					s.Rectangle(cs.Width.Max, cs.Height.Max, bg, "ff")
					s.Inset(8, func() {
						s.Icon("Build", fg, bg, 32)
					})
				}, b)
				for b.Clicked(s.Gtx) {
					Debug("run mode folder clicked")
					if !s.Config.BuildOpen {
						s.Config.FilterOpen = false
						s.Config.SettingsOpen = false
					}
					s.Config.BuildOpen = !s.Config.BuildOpen
					s.SaveConfig()
				}
			}),
			)
		}
	})
}

func (s *State) BuildPage(headless bool) layout.FlexChild {
	gtx := s.Gtx
	if headless {
		gtx = s.Htx
	}
	if !s.Config.BuildOpen {
		return gui.Flexed(0, func() {})
	}
	var weight float32 = 0.5
	switch {
	case s.Config.BuildZoomed:
		weight = 1
	case s.WindowHeight <= 800 && s.WindowWidth <= 800:
		weight = 1
	case s.WindowHeight <= 600 && s.WindowWidth > 800:
		weight = 1
	}
	return gui.Flexed(weight, func() {
		cs := gtx.Constraints
		s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg", "ff")
		s.FlexV(gui.Rigid(func() {
			s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg", "ff")
			s.Inset(4, func() {})
		}), gui.Rigid(func() {
			s.FlexH(gui.Rigid(func() {
				s.Label("Monitor Configuration", "DocText", "DocBg")
			}), s.Spacer("DocBg"), gui.Rigid(func() {
				if !(s.WindowHeight <= 800 && s.WindowWidth <= 800 ||
					s.WindowHeight <= 600 && s.WindowWidth > 800) {
					ic := "zoom"
					if s.Config.BuildZoomed {
						ic = "minimize"
					}
					b := s.Buttons["BuildZoom"]
					s.IconButton(ic, "DocText", "DocBg", b)
					for b.Clicked(s.Gtx) {
						Debug("settings panel fold button clicked")
						s.Config.BuildZoomed = !s.Config.BuildZoomed
						s.SaveConfig()
					}
				}
			}), s.Spacer("DocBg"), gui.Rigid(func() {
				b := s.Buttons["BuildClose"]
				s.IconButton("foldIn", "DocText", "DocBg", b)
				for b.Clicked(s.Gtx) {
					Debug("settings panel close button clicked")
					s.Config.BuildOpen = false
					s.SaveConfig()
				}
			}),
			)
		}), gui.Flexed(1, func() {
			cs := gtx.Constraints
			s.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg", "ff")
			s.FlexV(gui.Flexed(1, func() {
				s.Inset(8, func() {
					// cs := gtx.Constraints
					// s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
					//if s.Config.BuildOpen {
					s.BuildConfigPage(headless)
					//}
				})
			}))
		}), gui.Rigid(func() {
			s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg", "ff")
			s.Inset(4, func() {})
		}),
		)
	})
}

func (s *State) BuildConfigPage(headless bool) {
	gtx := s.Gtx
	if headless {
		gtx = s.Htx
	}
	s.FlexV(
		//s.FlexH(
		gui.Rigid(func() {
			s.Inset(4, func() {
				s.FlexH(
					gui.Rigid(func() {
						s.Inset(8,
							s.Text("Run in", "PanelText", "PanelBg", "Primary", "h6"),
						)
					}), gui.Rigid(func() {
						if s.RunningInRepo {
							fg, bg := "DocText", "DocBg"
							if s.Config.RunInRepo {
								fg, bg = "ButtonText", "ButtonBg"
							}
							b := s.Buttons["RunningInRepo"]
							s.TextButton("repo", "Primary", 16,
								fg, bg, b)
							for b.Clicked(s.Gtx) {
								if !s.Config.Running {
									s.Config.RunInRepo = true
									s.CannotRun = false
									s.SaveConfig()
								}
							}
						}
					}), gui.Rigid(func() {
						fg, bg := "DocText", "DocBg"
						if !s.Config.RunInRepo {
							fg, bg = "ButtonText", "ButtonBg"
						}
						b := s.Buttons["RunFromProfile"]
						s.TextButton("profile", "Primary", 16,
							fg, bg, b)
						for b.Clicked(s.Gtx) {
							if !s.Config.Running {
								s.Config.RunInRepo = false
								s.CannotRun = false
								s.SaveConfig()
							}
						}
					}), gui.Rigid(func() {
						txt := "run pod in its repository"
						if !s.Config.RunInRepo {
							txt = "not implemented"
							s.CannotRun = true
						}
						s.Inset(8,
							s.Text(txt, "PanelText", "PanelBg", "Primary", "h6"),
						)
					}),
				)
			})
		}),
		gui.Rigid(func() {
			s.Inset(4, func() {
				s.FlexH(gui.Rigid(func() {
					s.Inset(8,
						s.Text("Use Go version", "PanelText", "PanelBg", "Primary", "h6"),
					)
				}), gui.Rigid(func() {
					if s.HasGo {
						fg, bg := "DocText", "DocBg"
						if s.Config.UseBuiltinGo {
							fg, bg = "ButtonText", "ButtonBg"
						}
						b := s.Buttons["UseBuiltinGo"]
						s.TextButton("builtin", "Primary", 16,
							fg, bg, b)
						for b.Clicked(s.Gtx) {
							if !s.Config.RunInRepo {
								s.Config.UseBuiltinGo = true
								s.CannotRun = false
								if !s.HasGo {
									s.CannotRun = true
								}
							}
						}
					}
				}), gui.Rigid(func() {
					fg, bg := "DocText", "DocBg"
					if !s.Config.UseBuiltinGo {
						fg, bg = "ButtonText", "ButtonBg"
					}
					b := s.Buttons["InstallNewGo"]
					s.TextButton("install new", "Primary", 16,
						fg, bg, b)
					for b.Clicked(s.Gtx) {
						if !s.Config.RunInRepo {
							s.Config.UseBuiltinGo = false
							s.CannotRun = false
							if !s.HasOtherGo {
								s.CannotRun = true
							}
						}
					}
				}), gui.Rigid(func() {
					txt := "build using built in go"
					if !s.Config.UseBuiltinGo {
						txt = "not implemented"
						s.CannotRun = true
					}
					s.Inset(8,
						s.Text(txt, "PanelText", "PanelBg", "Primary", "h6"),
					)
				}),
				)
			})
		}), gui.Rigid(func() {
			s.Inset(4, func() {
				s.FlexH(
					gui.Rigid(func() {
						s.Inset(8,
							s.Text("Log entry click command", "PanelText",
								"PanelBg",
								"Primary", "h6"),
						)
					}), gui.Rigid(func() {
						ww := len(s.Config.ClickCommand)
						//if ww < 12 {
						//	ww = 12
						//}
						gtx.Constraints.Width.Max = ww*10 + 30
						gtx.Constraints.Width.Min = ww*10 + 30
						s.Editor(&s.CommandEditor, ww, func(e gel.EditorEvent) {
							if e != nil {
								txt := s.CommandEditor.Text()
								if s.Config.ClickCommand == txt {
									return
								}
								s.Config.ClickCommand = txt
								Debug(s.Config.ClickCommand)
								s.SaveConfig()
							}
						}, headless)()
					}), gui.Rigid(func() {
						s.Inset(8,
							s.Text("When a log entry is clicked run this"+
								" command with variables substituted for"+
								" values from the log entry:\n\n"+
								"$1 is the source code file location\n"+
								"$2 is the line number", "PanelText",
								"PanelBg",
								"Primary", "h6"),
						)
					}),
				)
			})
		}),
	)
}
