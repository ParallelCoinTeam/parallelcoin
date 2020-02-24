// SPDX-License-Identifier: Unlicense OR MIT
package duoui

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/pkg/pod"
	"reflect"
)

type Field struct {
	Field *pod.Field
}

func (ui *DuoUI) InputField(f *Field) func() {
	return func() {
		layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(30), Left: unit.Dp(30), Right: unit.Dp(30)}.Layout(ui.ly.Context, func() {
			switch f.Field.Type {
			case "array":
				switch f.Field.InputType {
				case "text":

				default:
				}
			case "input":
				switch f.Field.InputType {
				case "text":
					e := ui.ly.Theme.DuoUIeditor(f.Field.Label)
					e.Font.Typeface = ui.ly.Theme.Font.Primary
					e.Font.Style = text.Italic
					e.Layout(ui.ly.Context, (ui.rc.Settings.Daemon.Widgets[f.Field.Label]).(*controller.Editor))
				case "number":
					e := ui.ly.Theme.DuoUIeditor(f.Field.Label)
					e.Font.Typeface = ui.ly.Theme.Font.Primary
					e.Font.Style = text.Italic
					e.Layout(ui.ly.Context, (ui.rc.Settings.Daemon.Widgets[f.Field.Label]).(*controller.Editor))
				case "password":
					e := ui.ly.Theme.DuoUIeditor(f.Field.Label)
					e.Font.Typeface = ui.ly.Theme.Font.Primary
					e.Font.Style = text.Italic
					e.Layout(ui.ly.Context, (ui.rc.Settings.Daemon.Widgets[f.Field.Label]).(*controller.Editor))
				default:
				}
			case "switch":
				ui.ly.Theme.DuoUIcheckBox(f.Field.Label, ui.ly.Theme.Color.Dark, ui.ly.Theme.Color.Dark).Layout(ui.ly.Context, (ui.rc.Settings.Daemon.Widgets[f.Field.Label]).(*controller.CheckBox))
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

var typeRegistry = make(map[string]reflect.Type)

func makeInstance(name string) interface{} {
	v := reflect.New(typeRegistry["cx.DuoUIconfigurationig."+name]).Elem()
	// Maybe fill in fields here if necessary
	return v.Interface()
}
