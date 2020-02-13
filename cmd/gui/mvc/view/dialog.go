package view

import (
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
	"image/color"
)

var (
	buttonDialogCancel = new(controller.Button)
	buttonDialogOK     = new(controller.Button)
	buttonDialogClose  = new(controller.Button)
)

type DuoUIdialog struct {
	*model.DuOScomponent
}

func DuoCOMdialog() *DuoUIdialog {
	dialog := *new(DuoUIdialog)
	cd := &model.DuOScomponent{
		Name:    "logger",
		Version: "0.1",
		//Model:      ,
		//Controller: c,
	}
	*dialog.DuOScomponent = *cd
	return &dialog
}

func (b *DuoUIdialog) View(gtx *layout.Context, th *theme.DuoUItheme, dl *model.DuoUIdialog) func() {
	return func() {
		cs := gtx.Constraints
		theme.DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, "ee000000", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		layout.Align(layout.Center).Layout(gtx, func() {
			//cs := gtx.Constraints
			theme.DuoUIdrawRectangle(gtx, 408, 150, th.Color.Primary, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

			layout.Flex{
				Axis:      layout.Vertical,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(func() {
					layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Rigid(func() {
							layout.Align(layout.Center).Layout(gtx, func() {
								layout.Inset{Top: unit.Dp(24), Bottom: unit.Dp(8), Left: unit.Dp(0), Right: unit.Dp(4)}.Layout(gtx, func() {
									cur := th.H4(dl.Text)
									cur.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
									cur.Alignment = text.Start
									cur.Layout(gtx)
								})
							})
						}),
					)
				}),
				layout.Rigid(func() {
					layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Rigid(dialogButon(func() { dl.Cancel() }, gtx, th, "CANCEL", "ffcfcfcf", "ffcf3030", "ffcfcfcf", buttonDialogCancel, th.Icons["iconCancel"])),
						layout.Rigid(dialogButon(func() { dl.Ok() }, gtx, th, "OK", "ffcfcfcf", "ff308030", "ffcfcfcf", buttonDialogOK, th.Icons["iconOK"])),
						layout.Rigid(dialogButon(func() { dl.Show = false }, gtx, th, "CLOSE", "ffcfcfcf", "ffcf8030", "ffcfcfcf", buttonDialogClose, th.Icons["iconClose"])),
					)
				}),
			)
		})
	}
}
func dialogButon(f func(), gtx *layout.Context, th *theme.DuoUItheme, t, txtColor, bgColor, iconColor string, button *controller.Button, icon *theme.DuoUIicon) func() {

	var b theme.DuoUIbutton
	return func() {
		layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func() {
			b = th.DuoUIbutton(t, txtColor, bgColor, iconColor, 24, 120, 60, 0, 0, icon)
			for button.Clicked(gtx) {
				f()
			}
			b.Layout(gtx, button)
		})
	}
}
