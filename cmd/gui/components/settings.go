package components

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/unit"
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
)

func DuoUIsettingsWidget(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) func() {
	return func() {
		layout.Flex{
			Axis: layout.Vertical,
		}.Layout(duo.DuoUIcontext,
			layout.Rigid(func() {
				layout.UniformInset(unit.Dp(15)).Layout(duo.DuoUIcontext, func() {
					layout.Flex{
						Axis: layout.Vertical,
					}.Layout(duo.DuoUIcontext,
						layout.Rigid(func() {
							duo.DuoUItheme.H3(duo.DuoUIconfiguration.Tabs.Current).Layout(duo.DuoUIcontext)
						}),
						layout.Rigid(func() {
							cs := duo.DuoUIcontext.Constraints
							helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, 64, helpers.HexARGB("ffcf44cf"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
							groupsNumber := len(rc.Settings.Daemon.Schema.Groups)
							groupsList.Layout(duo.DuoUIcontext, groupsNumber, func(i int) {
								layout.UniformInset(unit.Dp(0)).Layout(duo.DuoUIcontext, func() {
									i = groupsNumber - 1 - i
									t := rc.Settings.Daemon.Schema.Groups[i]
									txt := fmt.Sprint(t.Legend)
									for duo.DuoUIconfiguration.Tabs.TabsList[txt].Clicked(duo.DuoUIcontext) {
										duo.DuoUIconfiguration.Tabs.Current = txt
										log.INFO("unutra: ", txt)
									}
									duo.DuoUItheme.DuoUIbutton(txt, "ff303030", "ff989898", "ff303030", 0, 125, 32, 4, 4, nil).Layout(duo.DuoUIcontext, duo.DuoUIconfiguration.Tabs.TabsList[txt])
								})
							})
						}))
				})
			}),
			layout.Flexed(1, func() {
				for _, fields := range rc.Settings.Daemon.Schema.Groups {
					if fmt.Sprint(fields.Legend) == duo.DuoUIconfiguration.Tabs.Current {
						fieldsList.Layout(duo.DuoUIcontext, len(fields.Fields), func(il int) {
							il = len(fields.Fields) - 1 - il
							tl := Field{
								Field: &fields.Fields[il],
							}
							layout.Flex{
								Axis: layout.Vertical,
							}.Layout(duo.DuoUIcontext,
								layout.Rigid(func() {
									layout.Flex{}.Layout(duo.DuoUIcontext,
										layout.Rigid(func() {
											helpers.DuoUIdrawRectangle(duo.DuoUIcontext, 30, 3, duo.DuoUItheme.Color.Dark, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
										}),
										layout.Flexed(0.62, func() {
											layout.Flex{
												Axis:    layout.Vertical,
												Spacing: 10,
											}.Layout(duo.DuoUIcontext,
												layout.Rigid(func() {
													layout.UniformInset(unit.Dp(0)).Layout(duo.DuoUIcontext, func() {
														duo.DuoUItheme.H6(fmt.Sprint(tl.Field.Name)).Layout(duo.DuoUIcontext)
													})
												}),
												layout.Rigid(func() {
													layout.UniformInset(unit.Dp(0)).Layout(duo.DuoUIcontext, func() {
														duo.DuoUItheme.Body2(fmt.Sprint(tl.Field.Description)).Layout(duo.DuoUIcontext)
													})
												}),
											)
										}),
										layout.Flexed(0.38, func() {
											layout.Align(layout.Start).Layout(duo.DuoUIcontext, func() {
												layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(30), Left: unit.Dp(30), Right: unit.Dp(30)}.Layout(duo.DuoUIcontext, func() {
													tl.InputFields(duo, cx)
												})
											})
										}),
									)
								}),
								layout.Rigid(func() {
									cs := duo.DuoUIcontext.Constraints
									helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, 1, duo.DuoUItheme.Color.Dark, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
								}))
						})
					}
				}
			}),
		)
	}
}
