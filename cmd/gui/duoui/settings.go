package duoui

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
	"github.com/p9c/pod/pkg/log"
)

var (
	groupsList = &layout.List{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}
	fieldsList = &layout.List{
		Axis: layout.Vertical,
	}
)

func DuoUIsettings(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {

	layout.Flex{
		Axis: layout.Vertical,
	}.Layout(duo.Gc,
		layout.Rigid(func() {
			duo.Comp.Settings.Inset.Layout(duo.Gc, func() {
				helpers.DuoUIdrawRectangle(duo.Gc, duo.Cs.Width.Max, duo.Cs.Height.Max, helpers.HexARGB("ff30cfcf"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
				// Settings  <<<
				duo.Th.H5("settings :").Layout(duo.Gc)
				layout.Flex{
					Axis: layout.Vertical,
				}.Layout(duo.Gc,
					layout.Rigid(func() {

						duo.Th.H3(duo.Conf.Tabs.Current).Layout(duo.Gc)
					}),
					layout.Rigid(func() {
						cs := duo.Gc.Constraints
						helpers.DuoUIdrawRectangle(duo.Gc, cs.Width.Max, 64, helpers.HexARGB("ffcfcfcf"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
						groupsNumber := len(rc.Settings.Daemon.Schema.Groups)
						groupsList.Layout(duo.Gc, groupsNumber, func(i int) {
							in.Layout(duo.Gc, func() {
								i = groupsNumber - 1 - i
								t := rc.Settings.Daemon.Schema.Groups[i]
								txt := fmt.Sprint(t.Legend)
								for duo.Conf.Tabs.TabsList[txt].Clicked(duo.Gc) {
									duo.Conf.Tabs.Current = txt
									log.INFO(txt)
								}
								duo.Th.Button(txt).Layout(duo.Gc, duo.Conf.Tabs.TabsList[txt])
							})
						})
					}))

			})
		}),
		layout.Flexed(1, func() {
			//cs := duo.Gc.Constraints
			for _, fields := range rc.Settings.Daemon.Schema.Groups {
				if fmt.Sprint(fields.Legend) == duo.Conf.Tabs.Current {
					fieldsList.Layout(duo.Gc, len(fields.Fields), func(il int) {
						il = len(fields.Fields) - 1 - il
						tl := fields.Fields[il]
						layout.Flex{
							Axis: layout.Vertical,
						}.Layout(duo.Gc,
							layout.Rigid(func() {
								in.Layout(duo.Gc, func() {
									duo.Th.H6(fmt.Sprint(tl.Name)).Layout(duo.Gc)
								})
							}),
							layout.Rigid(func() {
								in.Layout(duo.Gc, func() {
									duo.Th.Body2(fmt.Sprint(tl.Description)).Layout(duo.Gc)
								})
							}),
						)
					})
				}
			}
		}),
	)

	// Overview >>>
}
