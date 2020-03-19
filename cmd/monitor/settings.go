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

func (st *State) SettingsButtons() layout.FlexChild {
	return Flexed(1, func() {
		st.FlexH(Rigid(func() {
			bg, fg := "PanelBg", "PanelText"
			if st.Config.SettingsOpen {
				bg, fg = "DocBg", "DocText"
			}
			st.TextButton("Settings", "Secondary",
				23, fg, bg, st.SettingsFoldButton)
			for st.SettingsFoldButton.Clicked(st.Gtx) {
				L.Debug("settings folder clicked")
				switch {
				case !st.Config.SettingsOpen:
					st.Config.BuildOpen = false
					st.Config.SettingsOpen = true
				case st.Config.SettingsOpen:
					st.Config.SettingsOpen = false
				}
				st.SaveConfig()
			}
		}),
		)
	})
}

func (st *State) SettingsPage() layout.FlexChild {
	if !st.Config.SettingsOpen {
		return Flexed(0, func() {})
	}
	var weight float32 = 0.5
	var settingsInset = 0
	switch {
	case st.WindowWidth < 1024 && st.WindowHeight > 1024:
		// weight = 0.333
	case st.WindowHeight < 1024 && st.WindowWidth < 1024:
		weight = 1
	case st.WindowHeight < 600 && st.WindowWidth > 1024:
		weight = 1
	}
	return Flexed(weight, func() {
		st.Inset(settingsInset, func() {
			cs := st.Gtx.Constraints
			st.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
			st.FlexV(
				Rigid(func() {
					cs := st.Gtx.Constraints
					st.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
					st.Inset(4, func() {})
				}),
				Rigid(func() {
					st.FlexH(
						Rigid(func() {
							st.TextButton("Settings", "Secondary",
								23, "DocText", "DocBg",
								st.SettingsTitleCloseButton)
							for st.SettingsTitleCloseButton.Clicked(st.Gtx) {
								L.Debug("settings panel title close button clicked")
								st.Config.SettingsOpen = false
								st.SaveConfig()
							}
						}),
						Rigid(func() {
							if st.WindowWidth > 640 {
								st.SettingsTabs()
							}
						}),
						Spacer(),
						Rigid(func() {
							st.IconButton("minimize", "DocText", "DocBg",
								st.SettingsCloseButton)
							for st.SettingsCloseButton.Clicked(st.Gtx) {
								L.Debug("settings panel close button clicked")
								st.Config.SettingsOpen = false
								st.SaveConfig()
							}
						}),
					)
				}),
				Rigid(func() {
					if st.WindowWidth < 640 {
						cs := st.Gtx.Constraints
						st.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
						st.SettingsTabs()
					}
				}),
				Flexed(1, func() {
					st.Inset(settingsInset, func() {
						cs := st.Gtx.Constraints
						st.Rectangle(cs.Width.Max, cs.Height.Max,
							"PanelBg")
						st.SettingsBody()
					})
				}),
				Rigid(func() {
					cs := st.Gtx.Constraints
					st.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
					st.Inset(4, func() {})
				}),
			)
		})
	})
}

func (st *State) SettingsTabs() {
	groupsNumber := len(st.Rc.Settings.Daemon.Schema.Groups)
	st.GroupsList.Layout(st.Gtx, groupsNumber, func(i int) {
		color :=
			"DocText"
		bgColor :=
			"DocBg"
		i = groupsNumber - 1 - i
		txt := st.Rc.Settings.Daemon.Schema.Groups[i].Legend
		for st.Rc.Settings.Tabs.TabsList[txt].Clicked(st.Gtx) {
			st.Rc.Settings.Tabs.Current = txt
		}
		if st.Rc.Settings.Tabs.Current == txt {
			color =
				"PanelText"
			bgColor =
				"PanelBg"
		}
		st.TextButton(txt, "Primary", 16,
			color, bgColor, st.Rc.Settings.Tabs.TabsList[txt])
	})
}

func (st *State) SettingsBody() {
	st.FlexH(
		Rigid(func() {
			st.Theme.DuoUIitem(4, st.Theme.Colors["PanelBg"]).
				Layout(st.Gtx, layout.N, func() {
					for _, fields := range st.Rc.Settings.Daemon.Schema.Groups {
						if fmt.Sprint(fields.Legend) == st.Rc.Settings.Tabs.Current {
							st.SettingsFields.Layout(st.Gtx, len(fields.Fields),
								func(il int) {
									il = len(fields.Fields) - 1 - il
									tl := &Field{
										Field: &fields.Fields[il],
									}
									st.FlexH(
										st.SettingsItemLabel(tl),
										st.SettingsItemInput(tl),
									)
								},
							)
						}
					}
				})
		}),
	)
}

func (st *State) SettingsItemLabel(f *Field) layout.FlexChild {
	return Flexed(0.5, func() {
		st.Inset(10, func() {
			st.FlexV(
				Rigid(
					st.SettingsFieldLabel(st.Gtx, st.Theme, f),
				),
				Rigid(
					st.SettingsFieldDescription(st.Gtx, st.Theme, f),
				),
			)
		})
	})
}

func (st *State) SettingsItemInput(f *Field) layout.FlexChild {
	return Flexed(0.5, func() {
		st.Inset(10,
			DuoUIinputField(st.Rc, st.Gtx, st.Theme, &Field{Field: f.Field}),
		)
	})
}

func (st *State) SettingsFieldLabel(gtx *layout.Context, th *gelook.DuoUItheme, f *Field) func() {
	return func() {
		name := th.H6(fmt.Sprint(f.Field.Label))
		name.Color = th.Colors["PanelText"]
		name.Font.Typeface = th.Fonts["Primary"]
		name.Layout(gtx)
	}
}

func (st *State) SettingsFieldDescription(gtx *layout.Context, th *gelook.DuoUItheme, f *Field) func() {
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
				if f.Field.Model != "MinerPass" {
					StringsArrayEditor(gtx, th, (rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor),
						(rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor).Text(),
						func(e gel.EditorEvent) {
							rc.Settings.Daemon.Config[f.Field.Model] = (rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor).Text()
						})()
				}
			default:

			}
		case "input":
			switch f.Field.InputType {
			case "text":
				Editor(gtx, th, (rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor),
					(rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor).Text(),
					func(e gel.EditorEvent) {
						rc.Settings.Daemon.Config[f.Field.Model] = (rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor).Text()
					})()
			case "number":
				Editor(gtx, th, (rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor), (rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor).Text(),
					func(e gel.EditorEvent) {
						number, err := strconv.Atoi((rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor).Text())
						if err == nil {
						}
						rc.Settings.Daemon.Config[f.Field.Model] = number
					})()
			case "decimal":
				Editor(gtx, th, (rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor), (rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor).Text(),
					func(e gel.EditorEvent) {
						decimal, err := strconv.ParseFloat((rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor).Text(), 64)
						if err != nil {
						}
						rc.Settings.Daemon.Config[f.Field.Model] = decimal
					})()
			case "password":
				e := th.DuoUIeditor(f.Field.Label)
				e.Font.Typeface = th.Fonts["Primary"]
				e.Font.Style = text.Italic
				e.Layout(gtx, (rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor))
			default:
			}
		case "switch":
			th.DuoUIcheckBox(f.Field.Label, th.Colors["Light"], th.Colors["Light"]).Layout(gtx, (rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.CheckBox))
			if (rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.CheckBox).Checked(gtx) {
				rc.Settings.Daemon.Config[f.Field.Model] = true
			} else {
				rc.Settings.Daemon.Config[f.Field.Model] = false
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
