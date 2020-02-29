// SPDX-License-Identifier: Unlicense OR MIT
package component

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
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

func SettingsTabs(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
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
				th.DuoUIbutton(th.Font.Primary, txt, th.Color.Light, th.Color.Info, "", th.Color.Dark, 16, 0, 80, 32, 4, 4).Layout(gtx, rc.Settings.Tabs.TabsList[txt])
			})
		})
	}
}

func DuoUIinputField(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, f *Field) func() {
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
					Editor(gtx, th, (rc.Settings.Daemon.Widgets[f.Field.Label]).(*controller.Editor), f.Field.Label, func(e controller.SubmitEvent) {})
				case "number":
					e := th.DuoUIeditor(f.Field.Label)
					e.Font.Typeface = th.Font.Primary
					e.Font.Style = text.Italic
					e.Layout(gtx, (rc.Settings.Daemon.Widgets[f.Field.Label]).(*controller.Editor))
				case "password":
					e := th.DuoUIeditor(f.Field.Label)
					e.Font.Typeface = th.Font.Primary
					e.Font.Style = text.Italic
					e.Layout(gtx, (rc.Settings.Daemon.Widgets[f.Field.Label]).(*controller.Editor))
				default:
				}
			case "switch":
				th.DuoUIcheckBox(f.Field.Label, th.Color.Dark, th.Color.Dark).Layout(gtx, (rc.Settings.Daemon.Widgets[f.Field.Label]).(*controller.CheckBox))
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

func SettingsFieldLabel(gtx *layout.Context, th *theme.DuoUItheme, f *Field) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			name := th.H6(fmt.Sprint(f.Field.Label))
			name.Font.Typeface = th.Font.Primary
			name.Layout(gtx)
		})
	}
}

func SettingsFieldDescription(gtx *layout.Context, th *theme.DuoUItheme, f *Field) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			desc := th.Body2(fmt.Sprint(f.Field.Description))
			desc.Font.Typeface = th.Font.Primary
			desc.Layout(gtx)
		})
	}
}
