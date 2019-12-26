// SPDX-License-Identifier: Unlicense OR MIT
package duoui

import (
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/pkg/gio/text"
	"github.com/p9c/pod/pkg/gio/widget"
	"github.com/p9c/pod/pkg/pod"
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

type Field struct{

	field *pod.Field
}

func (f *Field)inputFields(duo *models.DuoUI) {

	switch f.field.Type {
	case "array":
		switch f.field.InputType {
		case "text":

		default:
		}
	case "input":
		switch f.field.InputType {
		case "text":
			e := duo.Th.Editor(f.field.Name)
			e.Font.Style = text.Italic
			lineEditor := (duo.Conf.Settings.Daemon.Widgets[f.field.Name]).(*widget.Editor)
			e.Layout(duo.Gc, lineEditor)
			for _, e := range lineEditor.Events(duo.Gc) {
				if _, ok := e.(widget.SubmitEvent); ok {
					//topLabel = e.Text
					lineEditor.SetText("")
				}
			}
		case "number":
			e := duo.Th.Editor(f.field.Name)
			e.Font.Style = text.Italic
			lineEditor := (duo.Conf.Settings.Daemon.Widgets[f.field.Name]).(*widget.Editor)
			e.Layout(duo.Gc, lineEditor)
			for _, e := range lineEditor.Events(duo.Gc) {
				if _, ok := e.(widget.SubmitEvent); ok {
					//topLabel = e.Text
					lineEditor.SetText("")
				}
			}
		default:
		}
	case "switch":
		duo.Th.CheckBox(f.field.Name).Layout(duo.Gc, (duo.Conf.Settings.Daemon.Widgets[f.field.Name]).(*widget.CheckBox))
	case "radio":
		//radioButtonsGroup := (duo.Conf.Settings.Daemon.Widgets[fieldName]).(*widget.Enum)
		//layout.Flex{}.Layout(duo.Gc,
		//	layout.Rigid(func() {
		//		duo.Th.RadioButton("r1", "RadioButton1").Layout(duo.Gc,  radioButtonsGroup)
		//
		//	}),
		//	layout.Rigid(func() {
		//		duo.Th.RadioButton("r2", "RadioButton2").Layout(duo.Gc, radioButtonsGroup)
		//
		//	}),
		//	layout.Rigid(func() {
		//		duo.Th.RadioButton("r3", "RadioButton3").Layout(duo.Gc, radioButtonsGroup)
		//
		//	}))
	default:
		//duo.Th.CheckBox("Checkbox").Layout(duo.Gc, (duo.Conf.Settings.Daemon.Widgets[fieldName]).(*widget.CheckBox))

	}
}