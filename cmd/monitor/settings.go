package monitor

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/pkg/gui"
	"github.com/p9c/pod/pkg/gui/gel"
	"github.com/p9c/pod/pkg/gui/gelook"
	"github.com/p9c/pod/pkg/pod"
	"strconv"
	"strings"
	"time"
)

type Field struct {
	Field *pod.Field
}

func (s *State) SettingsButtons() layout.FlexChild {
	return gui.Rigid(func() {
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
				s.Rectangle(cs.Width.Max, cs.Height.Max, bg)
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
		return gui.Flexed(0, func() {})
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
	return gui.Flexed(weight, func() {
		cs := s.Gtx.Constraints
		s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
		s.FlexV(
			gui.Rigid(func() {
				cs := s.Gtx.Constraints
				s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
				s.Inset(4, func() {})
			}),
			gui.Rigid(func() {
				s.FlexH(gui.Rigid(func() {
					s.Label("Pod Settings", "DocText", "DocBg")
				}), gui.Flexed(1, func() {
					if s.WindowWidth > settingsTabBreak {
						s.SettingsTabs(27)
					}
				}), gui.Rigid(func() {
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
				}), gui.Rigid(func() {
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
			}), gui.Rigid(func() {
				if s.WindowWidth <= settingsTabBreak {
					cs := s.Gtx.Constraints
					s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
					si := 17
					if s.WindowWidth >= settingsTabBreakSmall {
						si = 20
					}
					if s.WindowWidth >= settingsTabBreakMedium {
						si = 24
					}
					s.SettingsTabs(si)
				}
			}), gui.Flexed(1, func() {
				cs := s.Gtx.Constraints
				s.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg")
				s.Inset(8, func() { s.SettingsBody() })
			}), gui.Rigid(func() {
				cs := s.Gtx.Constraints
				s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
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
func (s *State) SettingsItem(fields pod.Group) func(il int) {
	return func(il int) {
		//il = len(fields.Fields) - 1 - il
		tl := &Field{
			Field: &fields.Fields[il],
		}
		s.FlexH(gui.Flexed(1, func() {
			s.Inset(8, func() {
				s.FlexV(
					gui.Rigid(s.SettingsFieldLabel(tl)),
					gui.Rigid(func() {
						s.FlexH(
							gui.Rigid(s.SettingsItemInput(tl)),
							gui.Rigid(s.SettingsFieldDescription(s.Gtx, s.Theme, tl)),
						)
					}),
				)
			})
		}),
		)
	}
}

func (s *State) SettingsBody() {
	s.FlexH(gui.Rigid(func() {
		s.Theme.DuoUIcontainer(4, s.Theme.Colors["PanelBg"]).Layout(s.Gtx, layout.N, func() {
			for _, fields := range s.Rc.Settings.Daemon.Schema.Groups {
				if fmt.Sprint(fields.Legend) == s.Rc.Settings.Tabs.Current {
					s.Lists["SettingsFields"].Layout(s.Gtx,
						len(fields.Fields), s.SettingsItem(fields))
				}
			}
		})
	}))
}

func (s *State) SettingsItemLabel(f *Field) func() {
	return func() {
		//s.Gtx.Constraints.Width.Max = 32 * 10
		s.Gtx.Constraints.Width.Min = 32 * 10
		s.Inset(8, func() {
			s.FlexV(gui.Rigid(s.SettingsFieldLabel(f)))
		})
	}
}

func (s *State) SettingsItemInput(f *Field) func() {
	return func() {
		s.Inset(4, s.InputField(&Field{Field: f.Field}))
	}
}

func (s *State) SettingsFieldLabel(f *Field) func() {
	return func() {
		s.Inset(4, func() {
			layout.W.Layout(s.Gtx, func() {
				name := s.Theme.H6(fmt.Sprint(f.Field.Label))
				name.Color = s.Theme.Colors["DocText"]
				name.Font.Typeface = s.Theme.Fonts["Secondary"]
				name.Layout(s.Gtx)
			})
		})
	}
}

func (s *State) SettingsFieldDescription(gtx *layout.Context, th *gelook.DuoUItheme, f *Field) func() {
	return func() {
		s.Inset(4, func() {
			layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceAround,
			}.Layout(s.Gtx, gui.Rigid(func() {
				desc := th.Body1(fmt.Sprint(f.Field.Description))
				desc.Font.Typeface = th.Fonts["Primary"]
				desc.Color = th.Colors["DocText"]
				desc.Layout(gtx)
			}),
			)
		})
	}
}

func (s *State) InputField(f *Field) func() {
	rcs := s.Rc.Settings
	rdw, rdCfg := rcs.Daemon.Widgets, rcs.Daemon.Config
	return func() {
		switch f.Field.Type {
		case "stringSlice":
			switch f.Field.InputType {
			case "text":
				if f.Field.Model != "MiningAddrs" {
					s.StringsArrayEditor(
						(rdw[f.Field.Model]).(*gel.Editor),
						(rdw[f.Field.Model]).(*gel.Editor).Text(),
						func(e gel.EditorEvent) {
							rdCfg[f.Field.Model] =
								strings.Fields((rdw[f.Field.Model]).(*gel.Editor).Text())
							Debug()
							if e != nil {
								s.Rc.SaveDaemonCfg()
							}
						},
					)()
				}
			default:
			}
		case "input":
			switch f.Field.InputType {
			case "text":
				s.Editor(
					(rdw[f.Field.Model]).(*gel.Editor),
					(rdw[f.Field.Model]).(*gel.Editor).Text(),
					func(e gel.EditorEvent) {
						txt := rdw[f.Field.Model].(*gel.Editor).Text()
						rdCfg[f.Field.Model] = txt
						if e != nil {
							s.Rc.SaveDaemonCfg()
						}
					},
				)()
			case "number":
				s.Editor(
					(rdw[f.Field.Model]).(*gel.Editor),
					(rdw[f.Field.Model]).(*gel.Editor).Text(),
					func(e gel.EditorEvent) {
						number, err :=
							strconv.Atoi((rdw[f.Field.Model]).(*gel.Editor).Text())
						if err == nil {
						}
						rdCfg[f.Field.Model] = number
						if e != nil {
							s.Rc.SaveDaemonCfg()
						}
					},
				)()
			case "decimal":
				s.Editor(
					(rdw[f.Field.Model]).(*gel.Editor),
					(rdw[f.Field.Model]).(*gel.Editor).Text(),
					func(e gel.EditorEvent) {
						decimal, err :=
							strconv.ParseFloat(
								(rdw[f.Field.Model]).(*gel.Editor).Text(), 64)
						if err != nil {
						}
						rdCfg[f.Field.Model] = decimal
						if e != nil {
							s.Rc.SaveDaemonCfg()
						}
					})()
			case "password":
				s.PasswordEditor((rdw[f.Field.Model]).(*gel.Editor), func(e gel.EditorEvent) {
					txt := rdw[f.Field.Model].(*gel.Editor).Text()
					rdCfg[f.Field.Model] = txt
					if e != nil {
						s.Rc.SaveDaemonCfg()
					}
				})()
			case "time":
				//Debug("rendering duration")
				s.Editor(
					(rdw[f.Field.Model]).(*gel.Editor),
					(rdw[f.Field.Model]).(*gel.Editor).Text(),
					func(e gel.EditorEvent) {
						txt := rdw[f.Field.Model].(*gel.Editor).Text()
						var err error
						if rdCfg[f.Field.Model], err = time.ParseDuration(
							txt); Check(err) {
						}
						if e != nil {
							s.Rc.SaveDaemonCfg()
						}
					},
				)()
			default:
			}
		case "switch":
			bg, fg := s.Theme.Colors["DocBg"], s.Theme.Colors["DocText"]
			sw := s.Theme.DuoUIcheckBox("",
				//f.Field.Label,
				s.Theme.Colors["Primary"],
				s.Theme.Colors["Primary"])
			sw.PillColor = bg
			sw.CircleColor = fg
			sw.PillColorChecked = s.Theme.Colors["PrimaryDim"]
			sw.CircleColorChecked = s.Theme.Colors["Primary"]
			sw.DrawLayout(s.Gtx, rdw[f.Field.Model].(*gel.CheckBox))
			if (rdw[f.Field.Model]).(*gel.CheckBox).Checked(s.Gtx) {
				if !*rdCfg[f.Field.Model].(*bool) {
					tt := true
					rdCfg[f.Field.Model] = &tt
					s.Rc.SaveDaemonCfg()
				}
			} else {
				if *rdCfg[f.Field.Model].(*bool) {
					ff := false
					rdCfg[f.Field.Model] = &ff
					s.Rc.SaveDaemonCfg()
				}
			}
		case "radio":
		// radioButtonsGroup := (duo.Configuration.Settings.Daemon.Widgets[fieldName]).(*widget.Enum)
		// layout.Flex{}.Layout(g,
		//	layout.Rigid(func() {
		//		duo.Theme.RadioButton("r1", "RadioButton1").Layout(g,
		//		radioButtonsGroup)
		//
		//	}),
		//	layout.Rigid(func() {
		//		duo.Theme.RadioButton("r2", "RadioButton2").Layout(g,
		//		radioButtonsGroup)
		//
		//	}),
		//	layout.Rigid(func() {
		//		duo.Theme.RadioButton("r3", "RadioButton3").Layout(g,
		//		radioButtonsGroup)
		//
		//	}))
		default:
			// duo.Theme.CheckBox("Checkbox").Layout(g,
			//(duo.Configuration.Settings.Daemon.Widgets[fieldName]).(*widget.CheckBox))
		}
	}
}

const textWidth = 10

func (s *State) Editor(editorController *gel.Editor, label string, handler func(gel.EditorEvent)) func() {
	t, g := s.Theme, s.Gtx
	bg, fg := t.Colors["DocBg"], t.Colors["DocText"]
	return func() {
		t.DuoUIcontainer(8, bg).Layout(g, layout.NW, func() {
			width := g.Constraints.Width.Max / 2
			if width > 320 {
				width = 320
			}
			e := t.DuoUIeditor(label, fg, bg, width)
			e.Font.Typeface = t.Fonts["Mono"]
			e.TextSize = unit.Dp(16)
			layout.UniformInset(unit.Dp(4)).Layout(g, func() {
				e.Layout(g, editorController)
			})
			for _, e := range editorController.Events(g) {
				switch e.(type) {
				case gel.EditorEvent:
					//case gel.ChangeEvent:
					handler(e)
				}
			}
		})
	}
}

func (s *State) PasswordEditor(editorController *gel.Editor, handler func(gel.EditorEvent)) func() {
	t, g := s.Theme, s.Gtx
	bg, fg := t.Colors["DocBg"], t.Colors["DocText"]
	if !editorController.Focused() {
		fg = t.Colors["Transparent"]
	}
	return func() {
		t.DuoUIcontainer(8, bg).Layout(g, layout.NW, func() {
			width := g.Constraints.Width.Max / 2
			if width > 320 {
				width = 320
			}
			e := t.DuoUIeditor("", fg, bg, width)
			e.Font.Typeface = t.Fonts["Mono"]
			e.TextSize = unit.Dp(16)
			layout.UniformInset(unit.Dp(4)).Layout(g, func() {
				e.Layout(g, editorController)
			})
			for _, e := range editorController.Events(g) {
				switch e.(type) {
				//case gel.ChangeEvent:
				case gel.EditorEvent:
					handler(e)
				}
			}
		})
	}
}

func (s *State) StringsArrayEditor(editorController *gel.Editor, label string, handler func(gel.EditorEvent)) func() {
	t, g := s.Theme, s.Gtx
	bg, fg := t.Colors["DocBg"], t.Colors["DocText"]
	return func() {
		t.DuoUIcontainer(8, bg).Layout(g, layout.NW,
			func() {
				width := g.Constraints.Width.Max / 2
				if width > 320 {
					width = 320
				}
				e := t.DuoUIeditor(label, fg, bg, width)
				e.Font.Typeface = t.Fonts["Mono"]
				layout.UniformInset(unit.Dp(4)).Layout(g, func() {
					e.Layout(g, editorController)
				})
				for _, e := range editorController.Events(g) {
					switch e.(type) {
					case gel.ChangeEvent:
						handler(e)
					}
				}
			})
	}
}
