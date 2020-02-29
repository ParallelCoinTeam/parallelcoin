// SPDX-License-Identifier: Unlicense OR MIT
package component

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/pod"
	"reflect"
)

type Field struct {
	Field *pod.Field
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
					e := th.DuoUIeditor(f.Field.Label)
					e.Font.Typeface = th.Font.Primary
					e.Font.Style = text.Italic
					e.Layout(gtx, (rc.Settings.Daemon.Widgets[f.Field.Label]).(*controller.Editor))
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

var typeRegistry = make(map[string]reflect.Type)

func makeInstance(name string) interface{} {
	v := reflect.New(typeRegistry["cx.DuoUIconfigurationig."+name]).Elem()
	// Maybe fill in fields here if necessary
	return v.Interface()
}
