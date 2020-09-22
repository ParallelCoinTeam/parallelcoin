// Package monitor is a log viewer and filter and configuration interface
//
// +build !headless

package pkg

import (
	"gioui.org/layout"
	"github.com/p9c/pod/pkg/gui"
	"github.com/p9c/pod/pkg/gui/gel"
	"github.com/stalker-loki/app/slog"
)

func (s *State) BuildButtons() layout.FlexChild {
	return gui.Rigid(func() {
		if s.WindowWidth >= 360 || !s.Config.FilterOpen {
			s.FlexH(gui.Rigid(func() {
				bg, fg := "PanelBg", "PanelText"
				if s.Config.BuildOpen {
					bg, fg = "DocBg", "DocText"
				}
				b := s.Buttons["BuildFold"]
				// s.IconButton("Build", fg, bg, b)
				s.ButtonArea(func() {
					s.Gtx.Constraints.Width.Max = 48
					s.Gtx.Constraints.Height.Max = 48
					cs := s.Gtx.Constraints
					s.Rectangle(cs.Width.Max, cs.Height.Max, bg)
					s.Inset(8, func() {
						s.Icon("Build", fg, bg, 32)
					})
				}, b)
				for b.Clicked(s.Gtx) {
					slog.Debug("run mode folder clicked")
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

func (s *State) BuildPage() layout.FlexChild {
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
		cs := s.Gtx.Constraints
		s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
		s.FlexV(
			gui.Rigid(func() {
				s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
				s.Inset(4, func() {})
			}),
			gui.Rigid(func() {
				s.FlexHStart(
					gui.Rigid(func() {
						s.Label("Monitor Configuration", "DocText", "DocBg")
					}),
					s.Spacer("DocBg"),
					gui.Rigid(func() {
						if !(s.WindowHeight <= 800 && s.WindowWidth <= 800 ||
							s.WindowHeight <= 600 && s.WindowWidth > 800) {
							ic := "zoom"
							if s.Config.BuildZoomed {
								ic = "minimize"
							}
							b := s.Buttons["BuildZoom"]
							s.IconButton(ic, "DocText", "DocBg", b)
							for b.Clicked(s.Gtx) {
								slog.Debug("settings panel fold button clicked")
								s.Config.BuildZoomed = !s.Config.BuildZoomed
								s.SaveConfig()
							}
						}
					}),
					gui.Rigid(func() {
						b := s.Buttons["BuildClose"]
						s.IconButton("foldIn", "DocText", "DocBg", b)
						for b.Clicked(s.Gtx) {
							slog.Debug("settings panel close button clicked")
							s.Config.BuildOpen = false
							s.SaveConfig()
						}
					}),
				)
			}),
			gui.Flexed(1, func() {
				cs := s.Gtx.Constraints
				s.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg")
				s.FlexV(
					gui.Flexed(1, func() {
						s.Inset(8, func() {
							s.BuildConfigPage()
						})
					}))
			}),
			gui.Rigid(func() {
				s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
				s.Inset(4, func() {})
			}),
		)
	})
}

func (s *State) BuildConfigPage() {
	s.FlexV(
		gui.Rigid(func() {
			s.Inset(4, func() {
				s.FlexHStart(
					gui.Rigid(
						s.Text("Run in", "PanelText", "PanelBg", "Primary", "h6"),
					),
					gui.Rigid(func() {
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
					}),
					gui.Rigid(func() {
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
					}),
					gui.Rigid(func() {
						txt := "run pod in its repository"
						if !s.Config.RunInRepo {
							txt = "not implemented"
							s.CannotRun = true
						}
						s.Text(txt, "PanelText", "PanelBg", "Primary", "h6")()
					}),
				)
			})
		}),
		gui.Rigid(func() {
			s.Inset(4, func() {
				s.FlexHStart(
					gui.Rigid(
						s.Text("Use Go version", "PanelText", "PanelBg",
							"Primary", "h6"),
					),
					gui.Rigid(func() {
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
					}),
					gui.Rigid(func() {
						fg, bg := "DocText", "DocBg"
						if !s.Config.UseBuiltinGo {
							fg, bg = "ButtonText", "ButtonBg"
						}
						b := s.Buttons["InstallNewGo"]
						s.TextButton("install new", "Primary", 16, fg, bg, b)
						for b.Clicked(s.Gtx) {
							if !s.Config.RunInRepo {
								s.Config.UseBuiltinGo = false
								s.CannotRun = false
								if !s.HasOtherGo {
									s.CannotRun = true
								}
							}
						}
					}),
					gui.Rigid(func() {
						txt := "build using built in go"
						if !s.Config.UseBuiltinGo {
							txt = "not implemented"
							s.CannotRun = true
						}
						s.Text(txt, "PanelText", "PanelBg", "Primary", "h6")()
					}),
				)
			})
		}),
		gui.Rigid(func() {
			s.Inset(4, func() {
				s.FlexHStart(
					gui.Rigid(
						s.Text("Beta tester/Filter mode", "PanelText",
							"PanelBg", "Primary", "h6"),
					),
					gui.Rigid(func() {
						bg, fg := s.Theme.Colors["DocBg"], s.Theme.Colors["DocText"]
						cb := s.Theme.DuoUIcheckBox("", fg, bg)
						cb.PillColor = bg
						cb.CircleColor = fg
						cb.PillColorChecked = bg
						cb.CircleColorChecked = s.Theme.Colors["Primary"]
						isChecked := s.CheckBoxes["FilterMode"].Checked(s.Gtx)
						if isChecked != s.Config.FilterMode {
							s.Config.FilterMode = !s.Config.FilterMode
							s.SaveConfig()
							if s.Config.FilterMode {
								go s.RegenerateFilterBuf()
								s.FilterFunc = s.FilterOn
							} else {
								s.FilterFunc = NoFilter
							}
						}
						cb.DrawLayout(s.Gtx, s.CheckBoxes["FilterMode"])
					}),
					gui.Rigid(s.Text("Sets logger to trace but dynamically"+
						" filters log view according to filter settings on"+
						" sidebar", "PanelText", "PanelBg", "Primary", "h6"),
					),
				)
			})
		}),
		gui.Rigid(func() {
			s.Inset(4, func() {
				s.FlexHStart(
					gui.Rigid(
						s.Text("Log entry click command", "PanelText",
							"PanelBg", "Primary", "h6"),
					),
					gui.Rigid(func() {
						ww := len(s.Config.ClickCommand)
						s.Gtx.Constraints.Width.Max = ww*10 + 30
						s.Gtx.Constraints.Width.Min = ww*10 + 30
						s.Editor(&s.CommandEditor, "", func(e gel.EditorEvent) {
							if e != nil {
								txt := s.CommandEditor.Text()
								if s.Config.ClickCommand == txt {
									return
								}
								s.Config.ClickCommand = txt
								slog.Debug(s.Config.ClickCommand)
								s.SaveConfig()
							}
						})()
					}),
					gui.Rigid(s.Text("When a log entry is clicked run this"+
						" command with variables substituted for"+
						" values from the log entry:\n\n"+
						"$1 is the source code file location\n"+
						"$2 is the line number", "PanelText",
						"PanelBg",
						"Primary", "h6"),
					),
				)
			})
		}),
	)
}
