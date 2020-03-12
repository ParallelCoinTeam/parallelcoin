// SPDX-License-Identifier: Unlicense OR MIT
package component

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
	"github.com/p9c/pod/cmd/gui/rcd"
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
				i = groupsNumber - 1 - i
				t := rc.Settings.Daemon.Schema.Groups[i]
				txt := fmt.Sprint(t.Legend)
				for rc.Settings.Tabs.TabsList[txt].Clicked(gtx) {
					rc.Settings.Tabs.Current = txt
				}
				th.DuoUIbutton(th.Fonts["Primary"], txt, th.Colors["Light"], th.Colors["Info"], th.Colors["Info"], th.Colors["Light"], "", th.Colors["Dark"], 16, 0, 80, 32, 4, 4).Layout(gtx, rc.Settings.Tabs.TabsList[txt])
			})
		})
	}
}

func DuoUIinputField(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, f *Field) func() {
	return func() {
		layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(30), Left: unit.Dp(30), Right: unit.Dp(30)}.Layout(gtx, func() {
			switch f.Field.Type {
			case "array":
				switch f.Field.InputType {
				case "text":

				default:
				}
			case "input":
				switch f.Field.InputType {
				case "text":
					Editor(gtx, th, (rc.Settings.Daemon.Widgets[f.Field.Label]).(*gel.Editor), f.Field.Label, func(e gel.SubmitEvent) {})
				case "number":
					e := th.DuoUIeditor(f.Field.Label)
					e.Font.Typeface = th.Fonts["Primary"]
					e.Font.Style = text.Italic
					e.Layout(gtx, (rc.Settings.Daemon.Widgets[f.Field.Label]).(*gel.Editor))
				case "password":
					e := th.DuoUIeditor(f.Field.Label)
					e.Font.Typeface = th.Fonts["Primary"]
					e.Font.Style = text.Italic
					e.Layout(gtx, (rc.Settings.Daemon.Widgets[f.Field.Label]).(*gel.Editor))
				default:
				}
			case "switch":
				th.DuoUIcheckBox(f.Field.Label, th.Colors["Dark"], th.Colors["Dark"]).Layout(gtx, (rc.Settings.Daemon.Widgets[f.Field.Label]).(*gel.CheckBox))
			case "radio":
				//radioButtonsGroup := (duo.Configuration.Settings.Daemon.Widgets[fieldName]).(*widget.Enum)
				//layout.Flex{}.Layout(gtx,
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
				//duo.Theme.CheckBox("Checkbox").Layout(gtx, (duo.Configuration.Settings.Daemon.Widgets[fieldName]).(*widget.CheckBox))

			}
		})
	}
}

//
//var typeRegistry = make(map[string]reflect.Type)
//
//func makeInstance(name string) interface{} {
//	v := reflect.New(typeRegistry["cx.DuoUIconfigurationig."+name]).Elem()
//	// Maybe fill in fields here if necessary
//	return v.Interface()
//}

func SettingsFieldLabel(gtx *layout.Context, th *gelook.DuoUItheme, f *Field) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			name := th.H6(fmt.Sprint(f.Field.Label))
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
			desc.Layout(gtx)
		})
	}
}
