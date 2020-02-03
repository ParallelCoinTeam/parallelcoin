package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/pkg/gui/widget"
	"github.com/p9c/pod/pkg/gui/widget/parallel"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
)

var (
	buttonDialogCancel = new(widget.Button)
	buttonDialogOK     = new(widget.Button)
	buttonDialogClose  = new(widget.Button)
)

// Main wallet screen
func DuoUIdialog(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	// START View <<<
	//cs := duo.DuoUIcontext.Constraints
	//	helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, "ee303030", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
	//layout.Flexed(1, func() {
	//
	//	layout.Align(layout.Center).Layout(duo.DuoUIcontext, func() {
	//		layout.Inset{Top: unit.Dp(24), Bottom: unit.Dp(8), Left: unit.Dp(0), Right: unit.Dp(4)}.Layout(duo.DuoUIcontext, func() {
	//			cur := duo.DuoUItheme.H4("dddddddddddddddddddddddddddd")
	//			cur.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
	//			cur.Alignment = text.Start
	//			cur.Layout(duo.DuoUIcontext)
	//		})
	//	})
	//
	//})
	iconCancel, _ := parallel.NewDuoUIicon(icons.NavigationCancel)
	iconOK, _ := parallel.NewDuoUIicon(icons.NavigationCheck)
	iconClose, _ := parallel.NewDuoUIicon(icons.NavigationClose)

	cs := duo.DuoUIcontext.Constraints
	helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, "ee000000", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
	layout.Align(layout.Center).Layout(duo.DuoUIcontext, func() {
		cs := duo.DuoUIcontext.Constraints
		helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Min, cs.Height.Min, "ff000555", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

		layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		}.Layout(duo.DuoUIcontext,
			layout.Rigid(func() {
				layout.Flex{
					Axis: layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(duo.DuoUIcontext,
					layout.Rigid(func(){
							layout.Align(layout.Center).Layout(duo.DuoUIcontext, func() {
								layout.Inset{Top: unit.Dp(24), Bottom: unit.Dp(8), Left: unit.Dp(0), Right: unit.Dp(4)}.Layout(duo.DuoUIcontext, func() {
									cur := duo.DuoUItheme.H4("Dialog box!")
									cur.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
									cur.Alignment = text.Start
									cur.Layout(duo.DuoUIcontext)
								})
							})
					}),
				)
			}),
			layout.Rigid(func() {
				layout.Flex{
					Axis: layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(duo.DuoUIcontext,
					layout.Rigid(dialogButon("Cancel", duo, rc, buttonDialogCancel, iconCancel)),
					layout.Rigid(dialogButon("OK", duo, rc, buttonDialogOK, iconOK)),
					layout.Rigid(dialogButon("Close", duo, rc, buttonDialogClose, iconClose)),

				)
			}),
		)
	})
}

func dialogButon(text string, duo *models.DuoUI, rc *rcd.RcVar, button *widget.Button, icon *parallel.DuoUIicon) func() {
	var b parallel.DuoUIbutton
	return func() {
		layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(duo.DuoUIcontext, func() {
			b = duo.DuoUItheme.DuoUIbutton(text, "ffcf30cf", "ff3030cf", "ff30cfcf", 24, 120, 60, 0, 0, icon)
			for button.Clicked(duo.DuoUIcontext) {
				rc.IsNotificationRun = false
			}
			b.Layout(duo.DuoUIcontext, button)
		})
	}
}
