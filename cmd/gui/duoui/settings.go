package duoui

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
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

func (ui *DuoUI)DuoUIsettings() {
	layout.Flex{
		Axis: layout.Vertical,
	}.Layout(ui.ly.Context,
		layout.Rigid(func() {
			layout.UniformInset(unit.Dp(15)).Layout(ui.ly.Context, func() {
				layout.Flex{
					Axis: layout.Vertical,
				}.Layout(ui.ly.Context,
					layout.Rigid(func() {
						t := ui.ly.Theme.H3(ui.rc.Settings.Tabs.Current)
						t.Font.Typeface = ui.ly.Theme.Font.Primary
						t.Layout(ui.ly.Context)
					}),
					layout.Rigid(func() {
						cs := ui.ly.Context.Constraints
						theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 64, "ffcf44cf", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
						groupsNumber := len(ui.rc.Settings.Daemon.Schema.Groups)
						groupsList.Layout(ui.ly.Context, groupsNumber, func(i int) {
							layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
								i = groupsNumber - 1 - i
								t := ui.rc.Settings.Daemon.Schema.Groups[i]
								txt := fmt.Sprint(t.Legend)
								for ui.rc.Settings.Tabs.TabsList[txt].Clicked(ui.ly.Context) {
									ui.rc.Settings.Tabs.Current = txt
									log.INFO("unutra: ", txt)
								}
								ui.ly.Theme.DuoUIbutton(txt, "ff303030", "ff989898", "ff303030", 0, 125, 32, 4, 4, nil).Layout(ui.ly.Context, ui.rc.Settings.Tabs.TabsList[txt])
							})
						})
					}))
			})
		}),
		layout.Flexed(1, func() {
			for _, fields := range ui.rc.Settings.Daemon.Schema.Groups {
				if fmt.Sprint(fields.Legend) == ui.rc.Settings.Tabs.Current {
					fieldsList.Layout(ui.ly.Context, len(fields.Fields), func(il int) {
						il = len(fields.Fields) - 1 - il
						tl := Field{
							Field: &fields.Fields[il],
						}
						layout.Flex{
							Axis: layout.Vertical,
						}.Layout(ui.ly.Context,
							layout.Rigid(func() {
								layout.Flex{}.Layout(ui.ly.Context,
									layout.Rigid(func() {
										theme.DuoUIdrawRectangle(ui.ly.Context, 30, 3, ui.ly.Theme.Color.Dark, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
									}),
									layout.Flexed(0.62, func() {
										layout.Flex{
											Axis:    layout.Vertical,
											Spacing: 10,
										}.Layout(ui.ly.Context,
											layout.Rigid(func() {
												layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
													name := ui.ly.Theme.H6(fmt.Sprint(tl.Field.Name))
													name.Font.Typeface = ui.ly.Theme.Font.Primary
													name.Layout(ui.ly.Context)
												})
											}),
											layout.Rigid(func() {
												layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
													desc := ui.ly.Theme.Body2(fmt.Sprint(tl.Field.Description))
													desc.Font.Typeface = ui.ly.Theme.Font.Primary
													desc.Layout(ui.ly.Context)
												})
											}),
										)
									}),
									layout.Flexed(0.38, func() {
										layout.Align(layout.Start).Layout(ui.ly.Context, func() {
											layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(30), Left: unit.Dp(30), Right: unit.Dp(30)}.Layout(ui.ly.Context, func() {
												// TODO:
												// Input fileds must be set as theme part
												//tl.InputFields(ui.ly.Context, ui.ly.Theme, ui.rc.Settings)
											})
										})
									}),
								)
							}),
							layout.Rigid(func() {
								cs := ui.ly.Context.Constraints
								theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 1, ui.ly.Theme.Color.Dark, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
							}))
					})
				}
			}
		}),
	)
}
