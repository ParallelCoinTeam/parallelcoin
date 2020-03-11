package component

import (
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/gel"
	"github.com/p9c/gelook"
	"github.com/p9c/pod/cmd/gui/rcd"
	"image"
)

var (
	buttonDialogCancel = new(gel.Button)
	buttonDialogOK     = new(gel.Button)
	buttonDialogClose  = new(gel.Button)
)

func DuoUIdialog(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) {
	cs := gtx.Constraints
	layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func() {
			rr := float32(gtx.Px(unit.Dp(0)))
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(cs.Width.Max),
					Y: float32(cs.Height.Max),
				}},
				NE: rr, NW: rr, SE: rr, SW: rr,
			}.Op(gtx.Ops).Add(gtx.Ops)
			fill(gtx, gelook.HexARGB("ee000000"))
			pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
		}),
		layout.Stacked(func() {
			cs := gtx.Constraints
			layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Expanded(func() {
					rr := float32(gtx.Px(unit.Dp(0)))
					clip.Rect{
						Rect: f32.Rectangle{Max: f32.Point{
							X: float32(cs.Width.Max),
							Y: float32(cs.Height.Max),
						}},
						NE: rr, NW: rr, SE: rr, SW: rr,
					}.Op(gtx.Ops).Add(gtx.Ops)
					fill(gtx, gelook.HexARGB("ff888888"))
					pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
				}),
				layout.Stacked(func() {
					layout.Center.Layout(gtx, func() {
						layout.Inset{Top: unit.Dp(16), Bottom: unit.Dp(16), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func() {
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
											layout.Inset{Top: unit.Dp(0), Bottom: unit.Dp(8), Left: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, func() {
												cur := th.H4(rc.Dialog.Text)
												cur.Font.Typeface = th.Fonts["Primary"]
												cur.Color = gelook.HexARGB(th.Colors["Dark"])
												cur.Alignment = text.Start
												cur.Layout(gtx)
											})
										}),
									)
								}),
								layout.Rigid(func() {

									rc.Dialog.CustomField()

								}),
								layout.Rigid(func() {
									layout.Flex{
										Axis:      layout.Horizontal,
										Alignment: layout.Middle,
									}.Layout(gtx,
										layout.Rigid(dialogButon(gtx, th, rc.Dialog.Cancel, "CANCEL", "ffcf3030", "iconCancel", "ffcf8080", buttonDialogCancel)),
										layout.Rigid(dialogButon(gtx, th, rc.Dialog.Ok, "QUIT", "ff30cf30", "iconOK", "ff80cf80", buttonDialogOK)),
										layout.Rigid(dialogButon(gtx, th, rc.Dialog.Close, "RESTART", "ffcf8030", "iconClose", "ffcfa880", buttonDialogClose)),
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

func dialogButon(gtx *layout.Context, th *gelook.DuoUItheme, f func(), t, bgColor, icon, iconColor string, button *gel.Button) func() {
	var b gelook.DuoUIbutton
	return func() {
		layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func() {
			b = th.DuoUIbutton(th.Fonts["Primary"], t, th.Colors["Dark"], bgColor, th.Colors["Info"], bgColor, icon, iconColor, 16, 32, 120, 64, 0, 0)
			for button.Clicked(gtx) {
				f()
			}
			b.MenuLayout(gtx, button)
		})
	}
}
