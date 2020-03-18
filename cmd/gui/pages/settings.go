package pages

import (
	"fmt"

	"gioui.org/layout"

	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
)

var (
	fieldsList = &layout.List{
		Axis: layout.Vertical,
	}
	buttonSettingsRestart = new(gel.Button)
)

func Settings(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "SETTINGS",
		TxColor:       "",
		Command:       func() {},
		Border:        4,
		BorderColor:   th.Colors["Light"],
		Header:        SettingsHeader(rc, gtx, th),
		HeaderBgColor: "",
		HeaderPadding: 4,
		Body:          SettingsBody(rc, gtx, th),
		BodyBgColor:   th.Colors["Dark"],
		BodyPadding:   0,
		Footer:        func() {},
		FooterBgColor: "",
		FooterPadding: 0,
	}
	return th.DuoUIpage(page)

	//return th.DuoUIpage("SETTINGS", 0, func() {}, component.ContentHeader(gtx, th, SettingsHeader(rc, gtx, th)), SettingsBody(rc, gtx, th), func() {
	// var msg string
	// if rc.Settings.Daemon.Config["DisableBanning"].(*bool) != true{
	//	msg = "ima"
	// }else{
	//	msg = "nema"
	// //}
	// ttt := th.H6(fmt.Sprint(rc.Settings.Daemon.Config))
	// ttt.Color = gelook.HexARGB("ffcfcfcf")
	// ttt.Layout(gtx)
	//})
}

func SettingsHeader(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(component.SettingsTabs(rc, gtx, th)),
			layout.Rigid(func() {
				// var settingsRestartButton gelook.DuoUIbutton
				// settingsRestartButton = th.DuoUIbutton(th.Fonts["Secondary"],
				// 	"restart",
				// 	th.Colors["Light"],
				// 	th.Colors["Dark"],
				// 	th.Colors["Dark"],
				// 	th.Colors["Light"],
				// 	"",
				// 	th.Colors["Light"],
				// 	23, 0, 80, 48, 4, 4)
				// for buttonSettingsRestart.Clicked(gtx) {
				// 	rc.SaveDaemonCfg()
				// }
				// settingsRestartButton.Layout(gtx, buttonSettingsRestart)
			}),
		)
	}
}

func SettingsBody(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		// layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
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
							layout.Rigid(component.HorizontalLine(gtx, 1, th.Colors["Light"])))
					})
				}
			}
		})
		// })
	}
}

func SettingsItemRow(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, f *component.Field) func() {
	return func() {
		layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
		}.Layout(gtx,
			// layout.Rigid(func() {
			//	gelook.DuoUIdrawRectangle(gtx, 30, 3, th.Colors["Light"], [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			// }),
			layout.Flexed(0.62, func() {
				layout.Flex{
					Axis:    layout.Vertical,
					Spacing: 4,
				}.Layout(gtx,
					layout.Rigid(component.SettingsFieldLabel(gtx, th, f)),
					layout.Rigid(component.SettingsFieldDescription(gtx, th, f)),
				)
			}),
			layout.Flexed(1, component.DuoUIinputField(rc, gtx, th, f)),
		)
	}
}
