package duoui

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
	"image/color"
)

func DuoUIsettings(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	//groupsList := &layout.List{
	//	Axis:      layout.Horizontal,
	//	Alignment: layout.Middle,
	//}
	fieldsList := &layout.List{
		Axis: layout.Vertical,
	}
	selectedTab := "wallet"

	//tabsList := make(map[int]*widget.Button)

	layout.Flex{
		Axis: layout.Vertical,
	}.Layout(duo.Gc,
		layout.Rigid(func() {
			duo.Comp.Settings.Inset.Layout(duo.Gc, func() {
				helpers.DuoUIdrawRectangle(duo.Gc, duo.Cs.Width.Max, duo.Cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0xcf, B: 0xcf}, 0, 0, 0, 0, unit.Dp(0))
				// Settings  <<<
				//duo.Th.H5("settings :").Layout(duo.Gc)
				//layout.Flex{
				//	Axis: layout.Vertical,
				//}.Layout(duo.Gc,

				//layout.Rigid(func() {
				//	cs := duo.Gc.Constraints
				//	helpers.DuoUIdrawRectangle(duo.Gc, cs.Width.Max, 64, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}, 0, 0, 0, 0, unit.Dp(0))
				//	widgets := []func(){
				//		func() {
				//			duo.Th.H3(topLabel).Layout(gtx)
				//		},
				//		func() {
				//			duo.Gc.Constraints.Height.Max = duo.Gc.Px(unit.Dp(200))
				//			duo.Th.Editor("Hint").Layout(duo.Gc, editor)
				//		},
				//	}
				//	groupsList.Layout(duo.Gc, len(widgets), func(i int) {
				//		layout.UniformInset(unit.Dp(16)).Layout(duo.Gc, widgets[i])
				//	})
				//}))

				//groupsList.Layout(duo.Gc, len(rc.Settings.Daemon.Schema.Groups), func(i int) {
				//			in.Layout(duo.Gc, func() {
				//				i = len(rc.Settings.Daemon.Schema.Groups) - 1 - i
				//				t := rc.Settings.Daemon.Schema.Groups[i]
				//				txt := fmt.Sprint(t.Legend)
				//				tabsList[i] = new(widget.Button)
				//
				//				for tabsList[i].Clicked(duo.Gc) {
				//					selectedTab = txt
				//					log.INFO(txt)
				//				}
				//				//th.Button("Click me!").Layout(gtx, button)
				//				var btn material.Button
				//				btn = duo.Th.Button(txt)
				//				btn.Layout(duo.Gc, tabsList[i])
				//			})
				//
				//		})
				//	}))
			})
		}),
		layout.Flexed(1, func() {
			//cs := duo.Gc.Constraints

			for _, fields := range rc.Settings.Daemon.Schema.Groups {

				if fmt.Sprint(fields.Legend) == selectedTab {

					fieldsList.Layout(duo.Gc, len(fields.Fields), func(il int) {
						il = len(fields.Fields) - 1 - il
						tl := fields.Fields[il]
						txtc := fmt.Sprint(tl.Name)
						duo.Th.H6(txtc).Layout(duo.Gc)
					})
				}
			}
		}),
	)

	// Overview >>>
}
