package monitor

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"strconv"
	"strings"
	"time"

	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
	"github.com/p9c/pod/pkg/pod"
)

type Field struct {
	Field *pod.Field
}

func (s *State) SettingsButtons() layout.FlexChild {
	return Flexed(1, func() {
		s.FlexH(Rigid(func() {
			bg, fg := "PanelBg", "PanelText"
			if s.Config.SettingsOpen.Load() {
				bg, fg = "DocBg", "DocText"
			}
			s.TextButton("Settings", "Secondary",
				23, fg, bg, s.SettingsFoldButton)
			for s.SettingsFoldButton.Clicked(s.Gtx) {
				L.Debug("settings folder clicked")
				switch {
				case !s.Config.SettingsOpen.Load():
					s.Config.BuildOpen.Store(false)
					s.Config.SettingsOpen.Store(true)
				case s.Config.SettingsOpen.Load():
					s.Config.SettingsOpen.Store(false)
				}
				s.SaveConfig()
			}
		}),
		)
	})
}

func (s *State) SettingsPage() layout.FlexChild {
	if !s.Config.SettingsOpen.Load() {
		return Flexed(0, func() {})
	}
	var weight float32 = 0.5
	switch {
	case s.Config.SettingsZoomed.Load():
		weight = 1
	//case s.WindowWidth < 1024 && s.WindowHeight > 1024:
	// weight = 0.333
	case s.WindowHeight < 1024 && s.WindowWidth < 1024:
		weight = 1
	case s.WindowHeight < 600 && s.WindowWidth > 1024:
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
					s.TextButton("Settings", "Secondary",
						23, "DocText", "DocBg",
						s.SettingsTitleCloseButton)
					for s.SettingsTitleCloseButton.Clicked(s.Gtx) {
						L.Debug("settings panel title close button clicked")
						s.Config.SettingsOpen.Store(false)
						s.SaveConfig()
					}
				}), Rigid(func() {
					if s.WindowWidth > 640 {
						s.SettingsTabs()
					}
				}), Spacer(), Rigid(func() {
					if !(s.WindowHeight < 1024 && s.WindowWidth < 1024 ||
						s.WindowHeight < 600 && s.WindowWidth > 1024) {
						ic := "zoom"
						if s.Config.SettingsZoomed.Load() {
							ic = "minimize"
						}
						s.IconButton(ic, "DocText", "DocBg",
							s.SettingsZoomButton)
						for s.SettingsZoomButton.Clicked(s.Gtx) {
							L.Debug("settings panel close button clicked")
							s.Config.SettingsZoomed.Toggle()
							s.SaveConfig()
						}
					}
				}), Rigid(func() {
					s.IconButton("foldIn", "DocText", "DocBg",
						s.SettingsCloseButton)
					for s.SettingsCloseButton.Clicked(s.Gtx) {
						L.Debug("settings panel close button clicked")
						s.Config.SettingsOpen.Store(false)
						s.SaveConfig()
					}
				}),
				)
			}), Rigid(func() {
				if s.WindowWidth < 640 {
					cs := s.Gtx.Constraints
					s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg", "ff")
					s.SettingsTabs()
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

func (s *State) SettingsTabs() {
	groupsNumber := len(s.Rc.Settings.Daemon.Schema.Groups)

	s.GroupsList.Layout(s.Gtx, groupsNumber, func(i int) {
		color := "DocText"
		bgColor := "DocBg"
		i = groupsNumber - 1 - i
		txt := s.Rc.Settings.Daemon.Schema.Groups[i].Legend
		for s.Rc.Settings.Tabs.TabsList[txt].Clicked(s.Gtx) {
			s.Rc.Settings.Tabs.Current = txt
			s.Config.SettingsTab.Store(txt)
		}
		if s.Rc.Settings.Tabs.Current == txt {
			color = "PanelText"
			bgColor = "PanelBg"
		}
		s.TextButton(txt, "Primary", 16,
			color, bgColor, s.Rc.Settings.Tabs.TabsList[txt])
	})
}

func (s *State) SettingsBody() {
	s.FlexH(Flexed(1, func() {
		s.Theme.DuoUIitem(4, s.Theme.Colors["PanelBg"]).Layout(s.Gtx, layout.N, func() {
			for _, fields := range s.Rc.Settings.Daemon.Schema.Groups {
				if fmt.Sprint(fields.Legend) == s.Rc.Settings.Tabs.Current {
					s.SettingsFields.Layout(s.Gtx, len(fields.Fields), func(il int) {
						il = len(fields.Fields) - 1 - il
						tl := &Field{
							Field: &fields.Fields[il],
						}
						if tl.Field.Type == "switch" {
							s.FlexH(Flexed(1, func() {
								s.FlexH(
									s.SettingsItemLabel(tl),
									s.SettingsItemInput(tl),
								)
							}))
						} else {
							s.FlexH(Flexed(1, func() {
								s.FlexH(
									s.SettingsItemLabel(tl),
									s.SettingsItemInput(tl),
								)
							}))
						}
					},
					)
				}
			}
		})
	}),
	)
}

func (s *State) SettingsItemLabel(f *Field) layout.FlexChild {
	return Rigid(func() {
		s.Gtx.Constraints.Width.Max = 16*16 + 8
		s.Gtx.Constraints.Width.Min = 16*16 + 8
		s.Inset(10, func() {
			s.FlexV(
				Rigid(s.SettingsFieldLabel(s.Gtx, s.Theme, f)),
				Rigid(s.SettingsFieldDescription(s.Gtx, s.Theme, f)),
			)
		})
	})
}

func (s *State) SettingsItemInput(f *Field) layout.FlexChild {
	return Rigid(func() {
		s.Inset(10,
			s.InputField(&Field{Field: f.Field}),
		)
	})
}

func (s *State) SettingsFieldLabel(gtx *layout.Context, th *gelook.DuoUItheme, f *Field) func() {
	return func() {
		name := th.H6(fmt.Sprint(f.Field.Label))
		name.Color = th.Colors["PanelText"]
		name.Font.Typeface = th.Fonts["Primary"]
		name.Layout(gtx)
	}
}

func (s *State) SettingsFieldDescription(gtx *layout.Context, th *gelook.DuoUItheme, f *Field) func() {
	return func() {
		desc := th.Body2(fmt.Sprint(f.Field.Description))
		desc.Font.Typeface = th.Fonts["Primary"]
		desc.Color = th.Colors["PanelText"]
		desc.Layout(gtx)
	}
}

func (s *State) InputField(f *Field) func() {
	return func() {
		//gtx.Constraints.Width.Max = 8 + 32*16
		s.Gtx.Constraints.Width.Min = 8 + 32*16
		rsd := s.Rc.Settings.Daemon
		fld := f.Field
		fm := fld.Model
		rwe, ok := rsd.Widgets[fm].(*gel.Editor)
		var rwc *gel.CheckBox
		if !ok {
			rwc, ok = rsd.Widgets[fm].(*gel.CheckBox)
			if !ok {
				return
			}
		}
		_ = rwc
		switch fld.Type {
		case "stringSlice":
			switch fld.InputType {
			case "text":
				if fm != "MiningAddrs" {
					s.StringsArrayEditor(rsd.Widgets[fm].(*gel.
					Editor), (rsd.Widgets[fm]).(*gel.Editor).Text(), 42,
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
				s.Editor(rwe, 32, func(e gel.EditorEvent) {
					txt := rwe.Text()
					rsd.Config[fm] = txt
					if e != nil {
						s.Rc.SaveDaemonCfg()
					}
				})()
			case "number":
				s.Editor(rwe, 15, func(e gel.EditorEvent) {
					number, err := strconv.Atoi(rwe.Text())
					if err == nil {
						rsd.Config[fm] = number
					}
					if e != nil {
						s.Rc.SaveDaemonCfg()
					}
				})()
			case "time":
				s.Editor(rwe, 10, func(e gel.EditorEvent) {
					duration, err := time.ParseDuration(rwe.Text())
					if err == nil {
						rsd.Config[fm] = duration
					}
					if e != nil {
						s.Rc.SaveDaemonCfg()
					}
				})()
			case "decimal":
				s.Editor(rwe, 15, func(e gel.EditorEvent) {
					decimal, err := strconv.ParseFloat(rwe.Text(), 64)
					if err != nil {
						rsd.Config[fm] = decimal
					}
					if e != nil {
						s.Rc.SaveDaemonCfg()
					}
				})()
			case "password":
				s.PasswordEditor(rwe, 32, func(e gel.EditorEvent) {
					txt := rwe.Text()
					rsd.Config[fm] = txt
					if e != nil {
						s.Rc.SaveDaemonCfg()
					}
				})()
			default:
			}
		case "switch":
			s.Gtx.Constraints.Width.Max = 32 //16*15 + 8
			s.Gtx.Constraints.Width.Min = 16*15 + 8
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

func (s *State) Editor(editorControler *gel.Editor, width int,
	handler func(gel.EditorEvent)) func() {
	return func() {
		layout.UniformInset(unit.Dp(4)).Layout(s.Gtx, func() {
			outerColor := "DocBg"
			innerColor := "PanelBg"
			textColor := "PanelText"
			if editorControler.Focused() {
				outerColor = "DocText"
				//innerColor = "DocBg"
				//textColor = "PanelBg"
			}
			s.Rectangle(width*16+6, 38, outerColor, "bb", 4)
			s.Inset(3, func() {
				s.Rectangle(width*16, 32, innerColor, "ff", 2)
				e := s.Theme.DuoUIeditor(editorControler.Text(),
					s.Theme.Colors[textColor], s.Theme.Colors[innerColor], width)
				e.Font.Typeface = s.Theme.Fonts["Mono"]
				s.Inset(4, func() {
					s.FlexH(Rigid(func() {
						e.Layout(s.Gtx, editorControler)
					}),
					)
				})
				for _, e := range editorControler.Events(s.Gtx) {
					switch e.(type) {
					case gel.ChangeEvent:
						handler(e)
					}
				}
			})
		})
	}
}

func (s *State) PasswordEditor(editorControler *gel.Editor, width int,
	handler func(gel.EditorEvent)) func() {
	return func() {
		layout.UniformInset(unit.Dp(4)).Layout(s.Gtx, func() {
			outerColor := "DocBg"
			innerColor := "PanelBg"
			textColor := "PanelBg"
			if editorControler.Focused() {
				outerColor = "DocText"
				innerColor = "DocBg"
				textColor = "PanelBg"
			}
			s.Rectangle(width*16+6, 38, outerColor, "bb", 4)
			s.Inset(3, func() {
				s.Rectangle(width*16, 32, innerColor, "ff", 2)
				e := s.Theme.DuoUIeditor(editorControler.Text(),
					s.Theme.Colors[textColor], s.Theme.Colors[innerColor], width)
				e.Font.Typeface = s.Theme.Fonts["Mono"]
				s.Inset(4, func() {
					s.FlexH(Rigid(func() {
						e.Layout(s.Gtx, editorControler)
					}),
					)
				})
				for _, e := range editorControler.Events(s.Gtx) {
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
		//if len(split[len(split)-1]) < 1 && len(split) > 1 {
		//	split = split[:len(split)-1]
		//}
		height := 19*len(split) + 6
		//L.Debug(len(split), height, split)
		s.Theme.DuoUIitem(8, s.Theme.Colors["PanelBg"]).Layout(s.Gtx, layout.NW, func() {
			outerColor := "DocBg"
			innerColor := "PanelBg"
			textColor := "PanelText"
			if editorController.Focused() {
				outerColor = "DocText"
				innerColor = "PanelBg"
				textColor = "PanelText"
			}
			s.Rectangle(width*16+12, height+12, outerColor, "ff", 4)
			s.Inset(3, func() {
				s.Rectangle(width*16+6, height+6, innerColor, "ff", 2)
				s.Inset(6, func() {
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
