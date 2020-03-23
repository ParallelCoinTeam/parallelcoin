// SPDX-License-Identifier: Unlicense OR MIT
package component

import (
	"fmt"
	"strconv"
	"strings"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
	"github.com/p9c/pod/pkg/pod"
)

var (
	groupsList = &layout.List{
		Axis:      layout.Horizontal,
		Alignment: layout.Start,
	}
)

type Field struct {
	Field *pod.Field
}

func SettingsTabs(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		groupsNumber := len(rc.Settings.Daemon.Schema.Groups)
		groupsList.Layout(gtx, groupsNumber, func(i int) {
			layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
				color := th.Colors["Light"]
				bgColor := th.Colors["Dark"]
				i = groupsNumber - 1 - i
				t := rc.Settings.Daemon.Schema.Groups[i]
				txt := fmt.Sprint(t.Legend)
				for rc.Settings.Tabs.TabsList[txt].Clicked(gtx) {
					rc.Settings.Tabs.Current = txt
				}
				if rc.Settings.Tabs.Current == txt {
					color = th.Colors["Dark"]
					bgColor = th.Colors["Light"]
				}
				th.DuoUIbutton(th.Fonts["Primary"],
					txt, color, bgColor, "", "", "", "",
					16, 0, 80, 32, 4, 4).Layout(gtx, rc.Settings.Tabs.TabsList[txt])
			})
		})
	}
}

func DuoUIinputField(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, f *Field) func() {
	return func() {
		switch f.Field.Type {
		case "stringSlice":
			switch f.Field.InputType {
			case "text":
				if f.Field.Model != "MiningAddrs" {
					StringsArrayEditor(gtx, th, (rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor),
						(rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor).Text(),
						func(e gel.EditorEvent) {
							rc.Settings.Daemon.Config[f.Field.Model] = strings.Fields((rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor).Text())
							L.Debug()
							if e != nil {
								rc.SaveDaemonCfg()
							}
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
						txt := rc.Settings.Daemon.Widgets[f.Field.Model].(*gel.Editor).Text()
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
				e := th.DuoUIeditor(f.Field.Label, "DocText", "DocBg", 32)
				e.Font.Typeface = th.Fonts["Primary"]
				e.Font.Style = text.Italic
				e.Layout(gtx, (rc.Settings.Daemon.Widgets[f.Field.Model]).(*gel.Editor))
			default:
			}
		case "switch":
			th.DuoUIcheckBox(f.Field.Label, th.Colors["Primary"],
				th.Colors["Primary"]).Layout(gtx,
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

//
// var typeRegistry = make(map[string]reflect.Type)
//
// func makeInstance(name string) interface{} {
//	v := reflect.New(typeRegistry["cx.DuoUIconfigurationig."+name]).Elem()
//	// Maybe fill in fields here if necessary
//	return v.Interface()
// }

func SettingsFieldLabel(gtx *layout.Context, th *gelook.DuoUItheme, f *Field) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			name := th.H6(fmt.Sprint(f.Field.Label))
			name.Color = th.Colors["Dark"]
			name.Font.Typeface = th.Fonts["Primary"]
			name.Layout(gtx)
		})
	}
}

func SettingsFieldDescription(gtx *layout.Context, th *gelook.DuoUItheme, f *Field) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			desc := th.Body2(fmt.Sprint(f.Field.Description))
			desc.Font.Typeface = th.Fonts["Primary"]
			desc.Color = th.Colors["Dark"]
			desc.Layout(gtx)
		})
	}
}
