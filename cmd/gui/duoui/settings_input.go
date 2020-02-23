// SPDX-License-Identifier: Unlicense OR MIT
package duoui

import (
	"fmt"
	"gioui.org/text"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/pkg/log"
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
			//ui.DuoUIinputField(f)
			e := ui.ly.Theme.DuoUIeditor(f.Field.Name)
			e.Font.Typeface = ui.ly.Theme.Font.Primary
			e.Font.Style = text.Italic
			e.Layout(ui.ly.Context, (ui.rc.Settings.Daemon.Widgets[f.Field.Name]).(*controller.Editor))
			(ui.rc.Settings.Daemon.Widgets[f.Field.Name]).(*controller.Editor).SetText(fmt.Sprint(*f.Field.Value.(*string)))
			log.INFO(f.Field.Value)
		case "number":
			//ui.DuoUIinputField(f)
			e := ui.ly.Theme.DuoUIeditor(f.Field.Name)
			e.Font.Typeface = ui.ly.Theme.Font.Primary
			e.Font.Style = text.Italic
			e.Layout(ui.ly.Context, (ui.rc.Settings.Daemon.Widgets[f.Field.Name]).(*controller.Editor))
			//(ui.rc.Settings.Daemon.Widgets[f.Field.Name]).(*controller.Editor).SetText(fmt.Sprint(*f.Field.Value.(*string)))
			log.INFO(f.Field.Value)
		case "password":
			//ui.DuoUIinputField(f)
			e := ui.ly.Theme.DuoUIeditor(f.Field.Name)
			e.Font.Typeface = ui.ly.Theme.Font.Primary
			e.Font.Style = text.Italic
			e.Layout(ui.ly.Context, (ui.rc.Settings.Daemon.Widgets[f.Field.Name]).(*controller.Editor))
			//(ui.rc.Settings.Daemon.Widgets[f.Field.Name]).(*controller.Editor).SetText(fmt.Sprint(*f.Field.Value.(*string)))
			log.INFO(f.Field.Value)
		default:
		}
	case "switch":
		ui.ly.Theme.DuoUIcheckBox(f.Field.Name, "ff303030", "ff303030").Layout(ui.ly.Context, (ui.rc.Settings.Daemon.Widgets[f.Field.Name]).(*controller.CheckBox))
		log.INFO("")
		log.INFO(*f.Field.Value.(*bool))
		log.INFO("testLohg")

		(ui.rc.Settings.Daemon.Widgets[f.Field.Name]).(*controller.CheckBox).SetChecked(*f.Field.Value.(*bool))

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
