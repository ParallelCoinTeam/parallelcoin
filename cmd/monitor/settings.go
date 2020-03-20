package monitor

import (
	"fmt"
	"strconv"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/p9c/pod/cmd/gui/rcd"
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
	var settingsInset = 0
	switch {
	case s.WindowWidth < 1024 && s.WindowHeight > 1024:
		// weight = 0.333
	case s.WindowHeight < 1024 && s.WindowWidth < 1024:
		weight = 1
	case s.WindowHeight < 600 && s.WindowWidth > 1024:
		weight = 1
	}
	return Flexed(weight, func() {
		s.Inset(settingsInset, func() {
			cs := s.Gtx.Constraints
			s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
			s.FlexV(
				Rigid(func() {
					cs := s.Gtx.Constraints
					s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
					s.Inset(4, func() {})
				}),
				Rigid(func() {
					s.FlexH(
						Rigid(func() {
							s.TextButton("Settings", "Secondary",
								23, "DocText", "DocBg",
								s.SettingsTitleCloseButton)
							for s.SettingsTitleCloseButton.Clicked(s.Gtx) {
								L.Debug("settings panel title close button clicked")
								s.Config.SettingsOpen.Store(false)
								s.SaveConfig()
							}
						}),
						Rigid(func() {
							if s.WindowWidth > 640 {
								s.SettingsTabs()
							}
						}),
						Spacer(),
						Rigid(func() {
							s.IconButton("minimize", "DocText", "DocBg",
								s.SettingsCloseButton)
							for s.SettingsCloseButton.Clicked(s.Gtx) {
								L.Debug("settings panel close button clicked")
								s.Config.SettingsOpen.Store(false)
								s.SaveConfig()
							}
						}),
					)
				}),
				Rigid(func() {
					if s.WindowWidth < 640 {
						cs := s.Gtx.Constraints
						s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
						s.SettingsTabs()
					}
				}),
				Flexed(1, func() {
					s.Inset(settingsInset, func() {
						cs := s.Gtx.Constraints
						s.Rectangle(cs.Width.Max, cs.Height.Max,
							"PanelBg")
						s.SettingsBody()
					})
				}),
				Rigid(func() {
					cs := s.Gtx.Constraints
					s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
					s.Inset(4, func() {})
				}),
			)
		})
	})
}

func (s *State) SettingsTabs() {
	groupsNumber := len(s.Rc.Settings.Daemon.Schema.Groups)
	s.GroupsList.Layout(s.Gtx, groupsNumber, func(i int) {
		color :=
			"DocText"
		bgColor :=
			"DocBg"
		i = groupsNumber - 1 - i
		txt := s.Rc.Settings.Daemon.Schema.Groups[i].Legend
		for s.Rc.Settings.Tabs.TabsList[txt].Clicked(s.Gtx) {
			s.Rc.Settings.Tabs.Current = txt
		}
		if s.Rc.Settings.Tabs.Current == txt {
			color =
				"PanelText"
			bgColor =
				"PanelBg"
		}
		s.TextButton(txt, "Primary", 16,
			color, bgColor, s.Rc.Settings.Tabs.TabsList[txt])
	})
}

func (s *State) SettingsBody() {
	s.FlexH(
		Rigid(func() {
			s.Theme.DuoUIitem(4, s.Theme.Colors["PanelBg"]).
				Layout(s.Gtx, layout.N, func() {
					for _, fields := range s.Rc.Settings.Daemon.Schema.Groups {
						if fmt.Sprint(fields.Legend) == s.Rc.Settings.Tabs.Current {
							s.SettingsFields.Layout(s.Gtx, len(fields.Fields),
								func(il int) {
									il = len(fields.Fields) - 1 - il
									tl := &Field{
										Field: &fields.Fields[il],
									}
									s.FlexH(
										s.SettingsItemLabel(tl),
										s.SettingsItemInput(tl),
									)
								},
							)
						}
					}
				})
		}),
	)
}

func (s *State) SettingsItemLabel(f *Field) layout.FlexChild {
	return Flexed(0.5, func() {
		s.Inset(10, func() {
			s.FlexV(
				Rigid(
					s.SettingsFieldLabel(s.Gtx, s.Theme, f),
				),
				Rigid(
					s.SettingsFieldDescription(s.Gtx, s.Theme, f),
				),
			)
		})
	})
}

func (s *State) SettingsItemInput(f *Field) layout.FlexChild {
	return Flexed(0.5, func() {
		s.Inset(10,
			DuoUIinputField(s.Rc, s.Gtx, s.Theme, &Field{Field: f.Field}),
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

func StringsArrayEditor(gtx *layout.Context, th *gelook.DuoUItheme, editorControler *gel.Editor, label string, handler func(gel.EditorEvent)) func() {
	return func() {
		layout.UniformInset(unit.Dp(4)).Layout(gtx, func() {
			cs := gtx.Constraints
			gelook.DuoUIdrawRectangle(gtx, cs.Width.Max, 32,
				th.Colors["Light"], [4]float32{0, 0, 0, 0}, [4]float32{0, 0,
					0, 0})
			e := th.DuoUIeditor(label)
			e.Font.Typeface = th.Fonts["Mono"]
			// e.Font.Style = text.Italic
			layout.UniformInset(unit.Dp(4)).Layout(gtx, func() {
				e.Layout(gtx, editorControler)
			})
			for _, e := range editorControler.Events(gtx) {
				switch e.(type) {
				case gel.ChangeEvent:
					handler(e)
				}
			}
		})
	}
}

func DuoUIinputField(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, f *Field) func() {
	return func() {
		switch f.Field.Type {
		case "stringSlice":
			switch f.Field.InputType {
			case "text":
				//if f.Field.Model != "MinerPass" {
				StringsArrayEditor(gtx, th, (rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor),
					(rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor).Text(),
					func(e gel.EditorEvent) {
						rc.Settings.Daemon.Config[f.Field.Model] = (rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor).Text()
						L.Debug()
					})()
				//}
			default:

			}
		case "input":
			switch f.Field.InputType {
			case "text":
				Editor(gtx, th, (rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor),
					(rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor).Text(),
					func(e gel.EditorEvent) {
						txt := rc.Settings.Daemon.Widgets[f.Field.Model].(
						*gel.Editor).Text()
						rc.Settings.Daemon.Config[f.Field.Model] = txt
						if e != nil {
							rc.SaveDaemonCfg()
						}
					})()
			case "number":
				Editor(gtx, th, (rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor), (rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor).Text(),
					func(e gel.EditorEvent) {
						number, err := strconv.Atoi((rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor).Text())
						if err == nil {
						}
						rc.Settings.Daemon.Config[f.Field.Model] = number
						if e != nil {
							rc.SaveDaemonCfg()
						}
					})()
			case "decimal":
				Editor(gtx, th, (rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor), (rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor).Text(),
					func(e gel.EditorEvent) {
						decimal, err := strconv.ParseFloat((rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor).Text(), 64)
						if err != nil {
						}
						rc.Settings.Daemon.Config[f.Field.Model] = decimal
						if e != nil {
							rc.SaveDaemonCfg()
						}
					})()
			case "password":
				e := th.DuoUIeditor(f.Field.Label)
				e.Font.Typeface = th.Fonts["Primary"]
				e.Font.Style = text.Italic
				e.Layout(gtx, (rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor))
			default:
			}
		case "switch":
			th.DuoUIcheckBox(f.Field.Label, th.Colors["PanelText"],
				th.Colors["PanelText"]).Layout(gtx,
				(rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.CheckBox))
			if (rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.CheckBox).Checked(gtx) {
				if !*rc.Settings.Daemon.Config[f.Field.Model].(*bool) {
					tt := true
					rc.Settings.Daemon.Config[f.Field.Model] = &tt
					rc.SaveDaemonCfg()
				}
			} else {
				if *rc.Settings.Daemon.Config[f.Field.Model].(*bool) {
					ff := false
					rc.Settings.Daemon.Config[f.Field.Model] = &ff
					rc.SaveDaemonCfg()
				}
			}
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

func Editor(gtx *layout.Context, th *gelook.DuoUItheme, editorControler *gel.Editor, label string, handler func(gel.EditorEvent)) func() {
	return func() {
		layout.UniformInset(unit.Dp(4)).Layout(gtx, func() {
			cs := gtx.Constraints
			gelook.DuoUIdrawRectangle(gtx, cs.Width.Max, 32,
				th.Colors["Gray"],
				[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
				gelook.DuoUIdrawRectangle(gtx, cs.Width.Max, 32,
					th.Colors["Light"], [4]float32{0, 0, 0, 0}, [4]float32{0, 0,
						0, 0})
				e := th.DuoUIeditor(label)
				e.Font.Typeface = th.Fonts["Mono"]
				// e.Font.Style = text.Italic
				layout.UniformInset(unit.Dp(4)).Layout(gtx, func() {
					e.Layout(gtx, editorControler)
				})
				for _, e := range editorControler.Events(gtx) {
					switch e.(type) {
					case gel.ChangeEvent:
						handler(e)
					}
				}
			})
		})
	}
}