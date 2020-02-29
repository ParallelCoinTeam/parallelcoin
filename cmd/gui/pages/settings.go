package pages

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/component"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/log"
)

var (
	groupsList = &layout.List{
		Axis:      layout.Horizontal,
		Alignment: layout.Start,
	}
	fieldsList = &layout.List{
		Axis: layout.Vertical,
	}
	buttonSettingsSave = new(controller.Button)
)

func Settings(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) *theme.DuoUIpage {
	return th.DuoUIpage("SETTINGS", 0, func() {}, component.ContentHeader(gtx, th, headerSettings(rc, gtx, th)), settingsBody(rc, gtx, th), func() {})
}

func headerSettings(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func() {
				groupsNumber := len(rc.Settings.Daemon.Schema.Groups)
				groupsList.Layout(gtx, groupsNumber, func(i int) {
					layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
						i = groupsNumber - 1 - i
						t := rc.Settings.Daemon.Schema.Groups[i]
						txt := fmt.Sprint(t.Legend)
						for rc.Settings.Tabs.TabsList[txt].Clicked(gtx) {
							rc.Settings.Tabs.Current = txt
							log.INFO("unutra: ", txt)
						}
						th.DuoUIbutton(th.Font.Primary, txt, th.Color.Dark, "ff989898", "", th.Color.Dark, 16, 0, 80, 32, 4, 4).Layout(gtx, rc.Settings.Tabs.TabsList[txt])
					})
				})
			}),
			layout.Rigid(func() {
				var settingsSaveButton theme.DuoUIbutton
				settingsSaveButton = th.DuoUIbutton(th.Font.Secondary, "SAVE", th.Color.Light, th.Color.Dark, "", th.Color.Light, 16, 0, 128, 48, 0, 0)
				for buttonSettingsSave.Clicked(gtx) {
					//addressLineEditor.SetText(clipboard.Get())
				}
				settingsSaveButton.Layout(gtx, buttonSettingsSave)
			}),
		)
	}
}

func settingsBody(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		for _, fields := range rc.Settings.Daemon.Schema.Groups {
			if fmt.Sprint(fields.Legend) == rc.Settings.Tabs.Current {
				fieldsList.Layout(gtx, len(fields.Fields), func(il int) {
					il = len(fields.Fields) - 1 - il
					tl := component.Field{
						Field: &fields.Fields[il],
					}
					layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx,
						layout.Rigid(settingsItemRow(rc, gtx, th, &tl)),
						layout.Rigid(component.HorizontalLine(gtx, 1, th.Color.Dark)))
				})
			}
		}
	}
}

func settingsItemRow(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, f *component.Field) func() {
	return func() {
		layout.Flex{}.Layout(gtx,
			layout.Rigid(func() {
				theme.DuoUIdrawRectangle(gtx, 30, 3, th.Color.Dark, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			}),
			layout.Flexed(0.62, func() {
				layout.Flex{
					Axis:    layout.Vertical,
					Spacing: 10,
				}.Layout(gtx,
					layout.Rigid(component.SettingsFieldLabel(gtx, th, f)),
					layout.Rigid(component.SettingsFieldDescription(gtx, th, f)),
				)
			}),
			layout.Flexed(0.38, component.DuoUIinputField(rc, gtx, th, f)),
		)
	}
}
