// SPDX-License-Identifier: Unlicense OR MIT
package components

import (
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/pkg/gui/widget"
	"github.com/p9c/pod/pkg/conte"
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

func (f *Field) InputFields(duo *models.DuoUI, cx *conte.Xt) {

	switch f.Field.Type {
	case "array":
		switch f.Field.InputType {
		case "text":

		default:
		}
	case "input":
		switch f.Field.InputType {
		case "text":
			DuoUIinputField(duo, cx, f.Field.Name, f.Field.Model, (duo.DuoUIconfiguration.Settings.Daemon.Widgets[f.Field.Name]).(*widget.Editor))
		case "number":
			e := duo.DuoUItheme.DuoUIeditor(f.Field.Name, f.Field.Name)
			e.Font.Style = text.Italic
			lineEditor := (duo.DuoUIconfiguration.Settings.Daemon.Widgets[f.Field.Name]).(*widget.Editor)
			e.Layout(duo.DuoUIcontext, lineEditor)
			for _, e := range lineEditor.Events(duo.DuoUIcontext) {
				if _, ok := e.(widget.SubmitEvent); ok {
					//topLabel = e.Text
					lineEditor.SetText("")
				}
			}
		default:
		}
	case "switch":
		duo.DuoUItheme.DuoUIcheckBox(f.Field.Name).Layout(duo.DuoUIcontext, (duo.DuoUIconfiguration.Settings.Daemon.Widgets[f.Field.Name]).(*widget.CheckBox))
	case "radio":
		//radioButtonsGroup := (duo.DuoUIconfiguration.Settings.Daemon.Widgets[fieldName]).(*widget.Enum)
		//layout.Flex{}.Layout(duo.DuoUIcontext,
		//	layout.Rigid(func() {
		//		duo.DuoUItheme.RadioButton("r1", "RadioButton1").Layout(duo.DuoUIcontext,  radioButtonsGroup)
		//
		//	}),
		//	layout.Rigid(func() {
		//		duo.DuoUItheme.RadioButton("r2", "RadioButton2").Layout(duo.DuoUIcontext, radioButtonsGroup)
		//
		//	}),
		//	layout.Rigid(func() {
		//		duo.DuoUItheme.RadioButton("r3", "RadioButton3").Layout(duo.DuoUIcontext, radioButtonsGroup)
		//
		//	}))
	default:
		//duo.DuoUItheme.CheckBox("Checkbox").Layout(duo.DuoUIcontext, (duo.DuoUIconfiguration.Settings.Daemon.Widgets[fieldName]).(*widget.CheckBox))

	}
}

var typeRegistry = make(map[string]reflect.Type)

func makeInstance(name string) interface{} {
	v := reflect.New(typeRegistry["cx.DuoUIconfigurationig."+name]).Elem()
	// Maybe fill in fields here if necessary
	return v.Interface()
}
