package duoui

import (
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
)

var (
	buttonDialogCancel = new(controller.Button)
	buttonDialogOK     = new(controller.Button)
	buttonDialogClose  = new(controller.Button)

	list               = &layout.List{
		Axis: layout.Vertical,
	}
)

// Main wallet screen
func (ui *DuoUI) DuoUIdialog() {
	cs := ui.ly.Context.Constraints
	theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, cs.Height.Max, "ee000000", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		//cs := ui.ly.Context.Constraints
		theme.DuoUIdrawRectangle(ui.ly.Context, 408, 150, ui.ly.Theme.Color.Primary, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

		layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		}.Layout(ui.ly.Context,
			layout.Rigid(func() {
				layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(ui.ly.Context,
					layout.Rigid(func() {

							layout.Inset{Top: unit.Dp(24), Bottom: unit.Dp(8), Left: unit.Dp(0), Right: unit.Dp(4)}.Layout(ui.ly.Context, func() {
								cur := ui.ly.Theme.H4(ui.rc.Dialog.Text)
								cur.Font.Typeface = "bariol"
								cur.Color = theme.HexARGB(ui.ly.Theme.Color.Light)
								cur.Alignment = text.Start
								cur.Layout(ui.ly.Context)
							})
					}),
				)
			}),
			layout.Rigid(func() {
				layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(ui.ly.Context,
					layout.Rigid(ui.dialogButon(func(){ui.rc.Dialog.Cancel()},"CANCEL", "ffcf3030", "iconCancel", "ffcf8080", buttonDialogCancel)),
					layout.Rigid(ui.dialogButon(func(){ui.rc.Dialog.Ok()},"QUIT",  "ff30cf30", "iconOK", "ff80cf80", buttonDialogOK)),
					layout.Rigid(ui.dialogButon(func(){ui.rc.Dialog.Close()},"RESTART", "ffcf8030", "iconClose", "ffcfa880", buttonDialogClose)),
				)
			}),
		)

}

func (ui *DuoUI)dialogButon(f func(),t, bgColor, icon, iconColor string, button *controller.Button) func() {
	var b theme.DuoUIbutton
	return func() {
		layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(ui.ly.Context, func() {
			b = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Primary, t, "ff303030", bgColor, icon, iconColor, 48, 120, 60, 0, 0)
			for button.Clicked(ui.ly.Context) {
				f()
			}
			b.Layout(ui.ly.Context, button)
		})
	}
}
