// SPDX-License-Identifier: Unlicense OR MIT
package duoui

import (
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/widget"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/text"
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
	field *pod.Field
}

func (f *Field) inputFields(duo *models.DuoUI, cx *conte.Xt) {

	switch f.field.Type {
	case "array":
		switch f.field.InputType {
		case "text":

		default:
		}
	case "input":
		switch f.field.InputType {
		case "text":
			//helpers.DuoUIinputField(duo, cx, f.field.Name, f.field.Model, (duo.DuoUIconfiguration.Settings.Daemon.Widgets[f.field.Name]).(*widget.Editor) )
		case "number":
			e := duo.DuoUItheme.DuoUIeditor(f.field.Name, f.field.Name)
			e.Font.Style = text.Italic
			lineEditor := (duo.DuoUIconfiguration.Settings.Daemon.Widgets[f.field.Name]).(*widget.DuoUIeditor)
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
		duo.DuoUItheme.DuoUIcheckBox(f.field.Name).Layout(duo.DuoUIcontext, (duo.DuoUIconfiguration.Settings.Daemon.Widgets[f.field.Name]).(*widget.CheckBox))
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
