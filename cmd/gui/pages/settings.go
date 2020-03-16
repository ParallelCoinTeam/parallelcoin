package pages

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
)

var (
	fieldsList = &layout.List{
		Axis: layout.Vertical,
	}
	buttonSettingsSave = new(gel.Button)
)

func Settings(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) *gelook.DuoUIpage {
	return th.DuoUIpage("SETTINGS", 0, func() {}, component.ContentHeader(gtx, th, SettingsHeader(rc, gtx, th)), SettingsBody(rc, gtx, th), func() {
		//var msg string
		//if rc.Settings.Daemon.Config["DisableBanning"].(*bool) != true{
		//	msg = "ima"
		//}else{
		//	msg = "nema"
		////}
		//ttt := th.H6(fmt.Sprint(rc.Settings.Daemon.Config))
		//ttt.Color = gelook.HexARGB("ffcfcfcf")
		//ttt.Layout(gtx)
	})
}

func SettingsHeader(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(component.SettingsTabs(rc, gtx, th)),
			layout.Rigid(func() {
				var settingsSaveButton gelook.DuoUIbutton
				settingsSaveButton = th.DuoUIbutton(th.Fonts["Secondary"], "SAVE", th.Colors["Light"], th.Colors["Dark"], th.Colors["Dark"], th.Colors["Light"], "", th.Colors["Light"], 16, 0, 128, 48, 0, 0)
				for buttonSettingsSave.Clicked(gtx) {
					rc.SaveDaemonCfg()
				}
				settingsSaveButton.Layout(gtx, buttonSettingsSave)
			}),
		)
	}
}

func SettingsBody(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
			th.DuoUIitem(16, th.Colors["Dark"]).Layout(gtx, layout.N, func() {
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
								layout.Rigid(SettingsItemRow(rc, gtx, th, &tl)),
								layout.Rigid(component.HorizontalLine(gtx, 1, th.Colors["Dark"])))
						})
					}
				}
			})
		})
	}
}

func SettingsItemRow(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, f *component.Field) func() {
	return func() {
		layout.Flex{}.Layout(gtx,
			//layout.Rigid(func() {
			//	gelook.DuoUIdrawRectangle(gtx, 30, 3, th.Colors["Light"], [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			//}),
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
