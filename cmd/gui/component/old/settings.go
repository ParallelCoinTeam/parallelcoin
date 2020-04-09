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
	"github.com/p9c/pod/pkg/gui/gel"
	"github.com/p9c/pod/pkg/gui/gelook"
	"github.com/p9c/pod/pkg/pod"
)

var groupsList = &layout.List{
	Axis:      layout.Horizontal,
	Alignment: layout.Start,
}

type Field struct {
	Field *pod.Field
}

func SettingsTabs(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	rcs := rc.Settings
	return func() {
		groupsNumber := len(rcs.Daemon.Schema.Groups)
		groupsList.Layout(gtx, groupsNumber, func(i int) {
			layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
				color := th.Colors["Light"]
				bgColor := th.Colors["Dark"]
				i = groupsNumber - 1 - i
				t := rcs.Daemon.Schema.Groups[i]
				txt := fmt.Sprint(t.Legend)
				for rcs.Tabs.TabsList[txt].Clicked(gtx) {
					rcs.Tabs.Current = txt
				}
				if rcs.Tabs.Current == txt {
					color = th.Colors["Dark"]
					bgColor = th.Colors["Light"]
				}
				th.DuoUIbutton(th.Fonts["Primary"],
					txt, color, bgColor, "", "", "", "",
					16, 0, 80, 32,
					4, 4, 4, 4).Layout(gtx,
					rcs.Tabs.TabsList[txt])
			})
		})
	}
}

func DuoUIinputField(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, f *Field) func() {
	rcs := rc.Settings
	rdw := rcs.Daemon.Widgets
	rdCfg := rcs.Daemon.Config
	return func() {
		switch f.Field.Type {
		case "stringSlice":
			switch f.Field.InputType {
			case "text":
				if f.Field.Model != "MiningAddrs" {
					StringsArrayEditor(gtx, th, (rdw[f.Field.
						Model]).(*gel.Editor),
						(rdw[f.Field.Model]).(*gel.Editor).Text(),
						func(e gel.EditorEvent) {
							rdCfg[f.Field.Model] = strings.Fields((rdw[f.Field.Model]).(*gel.
							Editor).Text())
							Debug()
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
				Editor(gtx, th, (rdw[f.Field.Model]).(*gel.Editor),
					(rdw[f.Field.Model]).(*gel.Editor).Text(),
					func(e gel.EditorEvent) {
						txt := rdw[f.Field.Model].(*gel.Editor).
							Text()
						rdCfg[f.Field.Model] = txt
						if e != nil {
							rc.SaveDaemonCfg()
						}
					})()
			case "number":
				Editor(gtx, th, (rdw[f.Field.Model]).(*gel.
				Editor), (rdw[f.Field.Model]).(*gel.Editor).
					Text(),
					func(e gel.EditorEvent) {
						number, err := strconv.Atoi((rdw[f.
							Field.Model]).(*gel.Editor).Text())
						if err == nil {
						}
						rdCfg[f.Field.Model] = number
						if e != nil {
							rc.SaveDaemonCfg()
						}
					})()
			case "decimal":
				Editor(gtx, th, (rdw[f.Field.Model]).(*gel.
				Editor), (rdw[f.Field.Model]).(*gel.Editor).
					Text(),
					func(e gel.EditorEvent) {
						decimal, err := strconv.ParseFloat((rcs.Daemon.
							Widgets[f.Field.Model]).(*gel.Editor).Text(), 64)
						if err != nil {
						}
						rdCfg[f.Field.Model] = decimal
						if e != nil {
							rc.SaveDaemonCfg()
						}
					})()
			case "password":
				e := th.DuoUIeditor(f.Field.Label, "DocText", "DocBg", 32)
				e.Font.Typeface = th.Fonts["Primary"]
				e.Font.Style = text.Italic
				e.Layout(gtx, (rdw[f.Field.Model]).(*gel.Editor))
			default:
			}
		case "switch":
			sw := th.DuoUIcheckBox(f.Field.Label, th.Colors["Primary"],
				th.Colors["Primary"])
			sw.PillColor = th.Colors["LightGray"]
			sw.PillColorChecked = th.Colors["LightGrayI"]
			sw.CircleColor = th.Colors["LightGrayII"]
			sw.CircleColorChecked = th.Colors["Primary"]
			sw.DrawLayout(gtx,
				(rdw[f.Field.Model]).(*gel.CheckBox))
			if (rdw[f.Field.Model]).(*gel.CheckBox).Checked(gtx) {
				if !*rdCfg[f.Field.Model].(*bool) {
					tt := true
					rdCfg[f.Field.Model] = &tt
					rc.SaveDaemonCfg()
				}
			} else {
				if *rdCfg[f.Field.Model].(*bool) {
					ff := false
					rdCfg[f.Field.Model] = &ff
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
