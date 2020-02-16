// SPDX-License-Identifier: Unlicense OR MIT
package duoui

import (
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/view"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/pod"
	"reflect"
)

//
//var (
//	editor     = new(widget.Editor)
//	lineEditor = &widget.Editor{
//		SingleLine: true,
//		Submit:     true,
//	}
//	button            = new(widget.Button)
//	greenButton       = new(widget.Button)
//	iconButton        = new(widget.Button)
//	radioButtonsGroup = new(widget.Enum)
//	list              = &layout.List{
//		Axis: layout.Vertical,
//	}
//	green    = true
//	topLabel = "Hello, Gio"
//	icon     *material.Icon
//	checkbox = new(widget.CheckBox)
//)

type Field struct {
	Field *pod.Field
}

func (f *Field) InputFields(ui *DuoUI) {

	switch f.Field.Type {
	case "array":
		switch f.Field.InputType {
		case "text":

		default:
		}
	case "input":
		switch f.Field.InputType {
		case "text":
			view.DuoUIinputField(ui.ly, f.Field.Name, f.Field.Model, (ui.rc.Settings.Daemon.Widgets[f.Field.Name]).(*controller.Editor))
		case "number":
			e := ui.ly.Theme.DuoUIeditor(f.Field.Name)
			e.Font.Typeface = ui.ly.Theme.Font.Primary
			e.Font.Style = text.Italic
			lineEditor := (ui.rc.Settings.Daemon.Widgets[f.Field.Name]).(*controller.Editor)
			e.Layout(ui.ly.Context, lineEditor)
			for _, e := range lineEditor.Events(ui.ly.Context) {
				if _, ok := e.(controller.SubmitEvent); ok {
					//topLabel = e.Text
					lineEditor.SetText("")
				}
			}
		default:
		}
	case "switch":
		ui.ly.Theme.DuoUIcheckBox(f.Field.Name).Layout(ui.ly.Context, (ui.rc.Settings.Daemon.Widgets[f.Field.Name]).(*controller.CheckBox))
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
}

var typeRegistry = make(map[string]reflect.Type)

func makeInstance(name string) interface{} {
	v := reflect.New(typeRegistry["cx.DuoUIconfigurationig."+name]).Elem()
	// Maybe fill in fields here if necessary
	return v.Interface()
}
