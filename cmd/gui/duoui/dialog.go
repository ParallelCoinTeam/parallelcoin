package duoui

import (
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/controller"
	"github.com/p9c/pod/cmd/gui/theme"
	"image"
	"image/color"
)

var (
	buttonDialogCancel = new(controller.Button)
	buttonDialogOK     = new(controller.Button)
	buttonDialogClose  = new(controller.Button)

	list = &layout.List{
		Axis: layout.Vertical,
	}
)

func fill(gtx *layout.Context, col color.RGBA) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d}
}

func (ui *DuoUI) DuoUIdialog() {
	cs := ui.ly.Context.Constraints
	layout.Stack{Alignment: layout.Center}.Layout(ui.ly.Context,
		layout.Expanded(func() {
			rr := float32(ui.ly.Context.Px(unit.Dp(0)))
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(cs.Width.Max),
					Y: float32(cs.Height.Max),
				}},
				NE: rr, NW: rr, SE: rr, SW: rr,
			}.Op(ui.ly.Context.Ops).Add(ui.ly.Context.Ops)
			fill(ui.ly.Context, theme.HexARGB("ee000000"))
			pointer.Rect(image.Rectangle{Max: ui.ly.Context.Dimensions.Size}).Add(ui.ly.Context.Ops)

		}),
		layout.Stacked(func() {
			cs := ui.ly.Context.Constraints

			layout.Stack{Alignment: layout.Center}.Layout(ui.ly.Context,
				layout.Expanded(func() {
					rr := float32(ui.ly.Context.Px(unit.Dp(0)))
					clip.Rect{
						Rect: f32.Rectangle{Max: f32.Point{
							X: float32(cs.Width.Max),
							Y: float32(cs.Height.Max),
						}},
						NE: rr, NW: rr, SE: rr, SW: rr,
					}.Op(ui.ly.Context.Ops).Add(ui.ly.Context.Ops)
					fill(ui.ly.Context, theme.HexARGB("ff888888"))
					pointer.Rect(image.Rectangle{Max: ui.ly.Context.Dimensions.Size}).Add(ui.ly.Context.Ops)

				}),
				layout.Stacked(func() {
					layout.Center.Layout(ui.ly.Context, func() {
						layout.Inset{Top: unit.Dp(16), Bottom: unit.Dp(16), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(ui.ly.Context, func() {

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
											layout.Inset{Top: unit.Dp(0), Bottom: unit.Dp(8), Left: unit.Dp(4), Right: unit.Dp(4)}.Layout(ui.ly.Context, func() {
												cur := ui.ly.Theme.H4(ui.rc.Dialog.Text)
												cur.Font.Typeface = ui.ly.Theme.Font.Primary
												cur.Color = theme.HexARGB(ui.ly.Theme.Color.Dark)
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
										layout.Rigid(ui.dialogButon(ui.rc.Dialog.Cancel, "CANCEL", "ffcf3030", "iconCancel", "ffcf8080", buttonDialogCancel)),
										layout.Rigid(ui.dialogButon(ui.rc.Dialog.Ok, "QUIT", "ff30cf30", "iconOK", "ff80cf80", buttonDialogOK)),
										layout.Rigid(ui.dialogButon(ui.rc.Dialog.Close, "RESTART", "ffcf8030", "iconClose", "ffcfa880", buttonDialogClose)),
									)
								}),
							)
						})
					})
				}),
			)

		}),
	)
}

func (ui *DuoUI) dialogButon(f func(), t, bgColor, icon, iconColor string, button *controller.Button) func() {
	var b theme.DuoUIbutton
	return func() {
		layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(ui.ly.Context, func() {
			b = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Primary, t, ui.ly.Theme.Color.Dark, bgColor, ui.ly.Theme.Color.Info, bgColor, icon, iconColor, 16, 48, 120, 60, 0, 0)
			for button.Clicked(ui.ly.Context) {
				f()
			}
			b.MenuLayout(ui.ly.Context, button)
		})
	}
}
