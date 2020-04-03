package monitor

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"strconv"
	"strings"
	"time"

	"github.com/p9c/pod/pkg/gui/gel"
	"github.com/p9c/pod/pkg/gui/gelook"
	"github.com/p9c/pod/pkg/pod"
)

type Field struct {
	Field *pod.Field
}

func (s *State) SettingsButtons() layout.FlexChild {
	return Rigid(func() {
		if s.WindowWidth >= 360 || !s.Config.FilterOpen {
			bg, fg := "PanelBg", "PanelText"
			if s.Config.SettingsOpen {
				bg, fg = "DocBg", "DocText"
			}
			b := s.Buttons["SettingsFold"]
			s.ButtonArea(func() {
				s.Gtx.Constraints.Width.Max = 48
				s.Gtx.Constraints.Height.Max = 48
				cs := s.Gtx.Constraints
				s.Rectangle(cs.Width.Max, cs.Height.Max, bg, "ff")
				s.Inset(8, func() {
					s.Icon("settingsIcon", fg, bg, 32)
				})
			}, b)
			//s.IconButton("settingsIcon", fg, bg, b)
			for b.Clicked(s.Gtx) {
				Debug("settings folder clicked")
				if !s.Config.SettingsOpen {
					s.Config.FilterOpen = false
					s.Config.BuildOpen = false
				}
				s.Config.SettingsOpen = !s.Config.SettingsOpen
				s.SaveConfig()
			}
		}
	})
}

const settingsTabBreak = 960
const settingsTabBreakMedium = 640
const settingsTabBreakSmall = 512

func (s *State) SettingsPage() layout.FlexChild {
	if !s.Config.SettingsOpen {
		return Flexed(0, func() {})
	}
	var weight float32 = 0.5
	switch {
	case s.Config.SettingsZoomed:
		weight = 1
	//case s.WindowWidth < 1024 && s.WindowHeight > 1024:
	// weight = 0.333
	case s.WindowHeight <= 960 && s.WindowWidth <= 960:
		weight = 1
	case s.WindowHeight <= 600 && s.WindowWidth > 960:
		weight = 1
	}
	return Flexed(weight, func() {
		cs := s.Gtx.Constraints
		s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg", "ff")
		s.FlexV(
			Rigid(func() {
				cs := s.Gtx.Constraints
				s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg", "ff")
				s.Inset(4, func() {})
			}),
			Rigid(func() {
				s.FlexH(Rigid(func() {
					s.Label("Pod Settings")
				}), Flexed(1, func() {
					if s.WindowWidth > settingsTabBreak {
						s.SettingsTabs(27)
					}
				}), Rigid(func() {
					if !(s.WindowHeight <= 800 && s.WindowWidth <= 800 ||
						s.WindowHeight <= 600 && s.WindowWidth > 800) {
						ic := "zoom"
						if s.Config.SettingsZoomed {
							ic = "minimize"
						}
						b := s.Buttons["SettingsZoom"]
						//s.IconButton(ic, "DocText", "DocBg", b)
						s.ButtonArea(func() {
							s.Inset(8, func() {
								s.Icon(ic, "DocText", "DocBg", 32)
							})
						}, b)
						for b.Clicked(s.Gtx) {
							Debug("settings panel close button clicked")
							s.Config.SettingsZoomed = !s.Config.SettingsZoomed
							s.SaveConfig()
						}
					}
				}), Rigid(func() {
					b := s.Buttons["SettingsClose"]
					s.ButtonArea(func() {
						s.Inset(8, func() {
							s.Icon("foldIn", "DocText", "DocBg", 32)
						})
					}, b)
					//s.IconButton("foldIn", "DocText", "DocBg", b)
					for b.Clicked(s.Gtx) {
						Debug("settings panel close button clicked")
						s.Config.SettingsOpen = false
						s.SaveConfig()
					}
				}),
				)
			}), Rigid(func() {
				if s.WindowWidth <= settingsTabBreak {
					cs := s.Gtx.Constraints
					s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg", "ff")
					si := 17
					if s.WindowWidth >= settingsTabBreakSmall {
						si = 20
					}
					if s.WindowWidth >= settingsTabBreakMedium {
						si = 24
					}
					s.SettingsTabs(si)
				}
			}), Flexed(1, func() {
				cs := s.Gtx.Constraints
				s.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg", "ff")
				s.Inset(8, func() { s.SettingsBody() })
			}), Rigid(func() {
				cs := s.Gtx.Constraints
				s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg", "ff")
				s.Inset(4, func() {})
			}),
		)
	})
}

func (s *State) SettingsTabs(size int) {
	groupsNumber := len(s.Rc.Settings.Daemon.Schema.Groups)
	s.Lists["Groups"].Layout(s.Gtx, groupsNumber, func(i int) {
		color := "DocText"
		bgColor := "DocBg"
		i = groupsNumber - 1 - i
		txt := s.Rc.Settings.Daemon.Schema.Groups[i].Legend
		for s.Rc.Settings.Tabs.TabsList[txt].Clicked(s.Gtx) {
			s.Rc.Settings.Tabs.Current = txt
			s.Config.SettingsTab = txt
		}
		if s.Rc.Settings.Tabs.Current == txt {
			color = "PanelText"
			bgColor = "PanelBg"
		}
		s.TextButton(txt, "Primary", size,
			color, bgColor, s.Rc.Settings.Tabs.TabsList[txt])
	})
}

func (s *State) SettingsBody() {
	s.FlexH(
		Rigid(func() {
			s.Theme.DuoUIcontainer(4, s.Theme.Colors["PanelBg"]).
				Layout(s.Gtx, layout.N, func() {
					for _, fields := range s.Rc.Settings.Daemon.Schema.Groups {
						if fmt.Sprint(fields.Legend) == s.Rc.Settings.Tabs.Current {
							s.Lists["SettingsFields"].Layout(s.Gtx,
								len(fields.Fields), func(il int) {
									//il = len(fields.Fields) - 1 - il
									tl := &Field{
										Field: &fields.Fields[il],
									}
									s.FlexH(Flexed(1, func() {
										s.Inset(8, func() {
											s.FlexV(
												//Flexed(0.2, func() {}),
												Rigid(s.SettingsFieldLabel(tl)),
												Rigid(func() {
													s.FlexH(
														//Rigid(func() {}),
														Rigid(s.SettingsItemInput(tl)),
														Rigid(s.SettingsFieldDescription(s.Gtx, s.Theme, tl)),
													)
												}),
											)
										})
									}))
								})
						}
					}
				})
		}))
}

func (s *State) SettingsItemLabel(f *Field) func() {
	return func() {
		//s.Gtx.Constraints.Width.Max = 32 * 10
		s.Gtx.Constraints.Width.Min = 32 * 10
		s.Inset(4, func() {
			s.FlexV(
				Rigid(s.SettingsFieldLabel(f)),
			)
		})
	}
}

func (s *State) SettingsItemInput(f *Field) func() {
	return func() {
		s.Inset(4,
			s.InputField(&Field{Field: f.Field}),
		)
	}
}

func (s *State) SettingsFieldLabel(f *Field) func() {
	return func() {
		layout.W.Layout(s.Gtx, func() {
			name := s.Theme.H6(fmt.Sprint(f.Field.Label))
			name.Color = s.Theme.Colors["DocText"]
			name.Font.Typeface = s.Theme.Fonts["Secondary"]
			name.Layout(s.Gtx)
		})
	}
}

func (s *State) SettingsFieldDescription(gtx *layout.Context, th *gelook.DuoUItheme, f *Field) func() {
	return func() {
		layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceAround}.Layout(s.Gtx, Rigid(func() {
			desc := th.Body1(fmt.Sprint(f.Field.Description))
			desc.Font.Typeface = th.Fonts["Primary"]
			desc.Color = th.Colors["DocText"]
			desc.Layout(gtx)
		}),
		)
	}
}

func (s *State) InputField(f *Field) func() {
	return func() {
		rsd := s.Rc.Settings.Daemon
		fld := f.Field
		fm := fld.Model
		//var rwc *gel.CheckBox
		rwe, ok := rsd.Widgets[fm].(*gel.Editor)
		if !ok || rwe == nil {
			return
		}
		//rwc, ok = rsd.Widgets[fm].(*gel.CheckBox)
		//if !ok {
		//	return
		//}
		w := 0
		w = len(rwe.Text())
		switch fld.Type {
		case "stringSlice":
			switch fld.InputType {
			case "text":
				//s.Gtx.Constraints.Width.Min = (len(rwe.Text())-3)*10
				ww := len(rwe.Text())*10 + 40
				//if ww < 12 {
				//	ww = 12
				//}
				s.Gtx.Constraints.Width.Max = ww
				if fm != "MiningAddrs" {
					w := len((rsd.Widgets[fm]).(*gel.Editor).Text())
					s.StringsArrayEditor(rsd.Widgets[fm].(*gel.Editor), (rsd.Widgets[fm]).(*gel.Editor).Text(), w,
						func(e gel.EditorEvent) {
							rsd.Config[fm] = strings.Fields(rwe.Text())
							if e != nil {
								s.Rc.SaveDaemonCfg()
							}
						})()
				}
			default:
			}
		case "input":
			switch fld.InputType {
			case "text":
				ww := len(rwe.Text())
				//if ww < 12 {
				//	ww = 12
				//}
				s.Gtx.Constraints.Width.Max = ww*10 + 30
				s.Gtx.Constraints.Width.Min = ww*10 + 30
				s.Editor(rwe, w, func(e gel.EditorEvent) {
					txt := rwe.Text()
					rsd.Config[fm] = txt
					if e != nil {
						s.Rc.SaveDaemonCfg()
					}
				})()
			case "number":
				ww := len(rwe.Text())
				//if ww < 12 {
				//	ww = 12
				//}
				s.Gtx.Constraints.Width.Max = ww*10 + 30
				s.Gtx.Constraints.Width.Min = ww*10 + 30
				s.Editor(rwe, w, func(e gel.EditorEvent) {
					number, err := strconv.Atoi(rwe.Text())
					if err == nil {
						rsd.Config[fm] = number
					}
					if e != nil {
						s.Rc.SaveDaemonCfg()
					}
				})()
			case "time":
				ww := len(rwe.Text())
				//if ww < 12 {
				//	ww = 12
				//}
				s.Gtx.Constraints.Width.Max = ww*10 + 30
				s.Gtx.Constraints.Width.Min = ww*10 + 30
				s.Editor(rwe, w, func(e gel.EditorEvent) {
					duration, err := time.ParseDuration(rwe.Text())
					if err == nil {
						rsd.Config[fm] = duration
					}
					if e != nil {
						s.Rc.SaveDaemonCfg()
					}
				})()
			case "decimal":
				ww := len(rwe.Text())
				//if ww < 12 {
				//	ww = 12
				//}
				s.Gtx.Constraints.Width.Max = ww*10 + 30
				s.Gtx.Constraints.Width.Min = ww*10 + 30
				s.Editor(rwe, w, func(e gel.EditorEvent) {
					decimal, err := strconv.ParseFloat(rwe.Text(), 64)
					if err != nil {
						rsd.Config[fm] = decimal
					}
					if e != nil {
						s.Rc.SaveDaemonCfg()
					}
				})()
			case "password":
				ww := len(rwe.Text())
				//if ww < 12 {
				//	ww = 12
				//}
				s.Gtx.Constraints.Width.Max = ww*10 + 30
				s.Gtx.Constraints.Width.Min = ww*10 + 30
				s.PasswordEditor(rwe, w, func(e gel.EditorEvent) {
					txt := rwe.Text()
					rsd.Config[fm] = txt
					if e != nil {
						s.Rc.SaveDaemonCfg()
					}
				})()
			default:
			}
		case "switch":
			ww := 3
			s.Gtx.Constraints.Width.Max = ww * 10
			s.Gtx.Constraints.Width.Min = ww * 10
			layout.W.Layout(s.Gtx, func() {
				//s.Rectangle(32, 32, "DocBg", "88")
				color := "DocBg"
				if *rsd.Config[fm].(*bool) {
					color = "DocText"
				}
				s.Theme.DuoUIcheckBox("",
					//fld.Label,
					s.Theme.Colors[color],
					s.Theme.Colors[color]).Layout(s.Gtx,
					(rsd.Widgets[fm]).(*gel.CheckBox))
				if (rsd.Widgets[fm]).(*gel.CheckBox).Checked(s.Gtx) {
					if !*rsd.Config[fm].(*bool) {
						t := true
						rsd.Config[fm] = &t
						s.Rc.SaveDaemonCfg()
					}
				} else {
					if *rsd.Config[fm].(*bool) {
						f := false
						rsd.Config[fm] = &f
						s.Rc.SaveDaemonCfg()
					}
				}
			})
		case "radio":
			// radioButtonsGroup := (duo.Configuration.Settings.Daemon.Widgets[fieldName]).(*widget.Enum)
			// layout.Flex{}.Layout(gtx,
			//	layout.Rigid(func() {
			//		duo.Theme.RadioButton("r1", "RadioButton1").Layout(gtx,  radioButtonsGroup)
			//
			//	}),
			//	layout.Rigid(func() {
			//		duo.Theme.RadioButton("r2", "RadioButton2").Layout(gtx, radioButtonsGroup)
			//
			//	}),
			//	layout.Rigid(func() {
			//		duo.Theme.RadioButton("r3", "RadioButton3").Layout(gtx, radioButtonsGroup)
			//
			//	}))
		default:
			// duo.Theme.CheckBox("Checkbox").Layout(gtx, (duo.Configuration.Settings.Daemon.Widgets[fieldName]).(*widget.CheckBox))

		}
	}
}

const textWidth = 10

func (s *State) Editor(editorController *gel.Editor, width int,
	handler func(gel.EditorEvent)) func() {
	return func() {
		layout.UniformInset(unit.Dp(4)).Layout(s.Gtx, func() {
			outerColor := "DocBg"
			innerColor := "PanelBg"
			textColor := "PanelText"
			if editorController.Focused() {
				outerColor = "DocText"
				//innerColor = "DocBg"
				//textColor = "PanelBg"
			}
			width++
			s.Rectangle(width*textWidth+16, 40, outerColor, "ff", 4)
			s.Inset(3, func() {
				s.Rectangle(width*textWidth+10, 34, innerColor, "ff", 2)
				e := s.Theme.DuoUIeditor(editorController.Text(),
					s.Theme.Colors[textColor], s.Theme.Colors[innerColor], width)
				e.Font.Typeface = s.Theme.Fonts["Mono"]
				s.Inset(5, func() {
					s.FlexH(Rigid(func() {
						e.Layout(s.Gtx, editorController)
					}),
					)
				})
				for _, e := range editorController.Events(s.Gtx) {
					switch e.(type) {
					case gel.ChangeEvent:
						handler(e)
					}
				}
			})
		})
	}
}

func (s *State) PasswordEditor(editorController *gel.Editor, width int,
	handler func(gel.EditorEvent)) func() {
	return func() {
		layout.UniformInset(unit.Dp(4)).Layout(s.Gtx, func() {
			outerColor := "DocBg"
			innerColor := "PanelBg"
			textColor := "PanelBg"
			if editorController.Focused() {
				outerColor = "DocText"
				innerColor = "DocBg"
				textColor = "PanelBg"
			}
			width++
			s.Rectangle(width*textWidth+16, 40, outerColor, "ff", 4)
			s.Inset(3, func() {
				s.Rectangle(width*textWidth+10, 34, innerColor, "ff", 2)
				e := s.Theme.DuoUIeditor(editorController.Text(),
					s.Theme.Colors[textColor], s.Theme.Colors[innerColor], width)
				e.Font.Typeface = s.Theme.Fonts["Mono"]
				s.Inset(5, func() {
					s.FlexH(Rigid(func() {
						e.Layout(s.Gtx, editorController)
					}),
					)
				})
				for _, e := range editorController.Events(s.Gtx) {
					switch e.(type) {
					case gel.ChangeEvent:
						handler(e)
					}
				}
			})
		})
	}
}
func (s *State) StringsArrayEditor(editorController *gel.Editor, label string, width int, handler func(gel.EditorEvent)) func() {
	return func() {
		split := strings.Split(label, "\n")
		maxLen := 0
		for i := range split {
			if len(split[i]) > maxLen {
				maxLen = len(split[i])
			}
		}
		//if len(split[len(split)-1]) < 1 && len(split) > 2 {
		//	split = split[:len(split)-1]
		//}
		if maxLen < 9 {
			maxLen = 9
		}
		s.Gtx.Constraints.Width.Max = maxLen*textWidth + 30
		s.Gtx.Constraints.Width.Min = maxLen*textWidth + 30
		width = maxLen
		height := 19 * len(split)
		//Debug(len(split), height, split)
		s.Theme.DuoUIcontainer(0, s.Theme.Colors["PanelBg"]).Layout(s.Gtx, layout.N, func() {
			outerColor := "DocBg"
			innerColor := "PanelBg"
			textColor := "PanelText"
			if editorController.Focused() {
				outerColor = "DocText"
				innerColor = "PanelBg"
				textColor = "PanelText"
			}
			if width < 9 {
				width = 9
			}
			s.Rectangle(width*textWidth+16, height+16, outerColor, "ff", 4)
			s.Inset(3, func() {
				s.Rectangle(width*textWidth+10, height+10, innerColor, "ff", 2)
				s.Inset(5, func() {
					e := s.Theme.DuoUIeditor(label,
						s.Theme.Colors[textColor], s.Theme.Colors[innerColor], width)
					e.Font.Typeface = s.Theme.Fonts["Mono"]
					e.Layout(s.Gtx, editorController)
					for _, e := range editorController.Events(s.Gtx) {
						switch e.(type) {
						case gel.ChangeEvent:
							handler(e)
						}
					}
				})
			})
		})
	}
}
