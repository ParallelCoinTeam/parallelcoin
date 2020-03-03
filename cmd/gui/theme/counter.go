package theme

import (
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/controller"
	"github.com/p9c/pod/pkg/log"
	"image"
)

type DuoUIcounter struct {
	increase     DuoUIbutton
	decrease     DuoUIbutton
	reset        DuoUIbutton
	pageFunction func()
	Font         text.Font
	Text         string
	TextSize     unit.Value
	TxColor      string
	BgColor      string
	shaper       text.Shaper
}

func (t *DuoUItheme) DuoUIcounter(pageFunction func()) DuoUIcounter {
	return DuoUIcounter{
		//ToDo Replace theme's buttons with counter exclusive buttons, set icons for increase/decrease
		increase: t.DuoUIbutton("", "", "", t.Color.Primary, "", t.Color.Light, "counterPlusIcon", t.Color.Light, 0, 24, 24, 24, 0, 0),
		decrease: t.DuoUIbutton("", "", "", t.Color.Primary, "", t.Color.Light, "counterMinusIcon", t.Color.Light, 0, 24, 24, 24, 0, 0),
		//reset:        t.DuoUIbutton(t.Font.Secondary, "RESET", t.Color.Primary, t.Color.Light, t.Color.Light, t.Color.Primary, "", "", 12, 0, 0, 48, 48, 0),
		pageFunction: pageFunction,
		Font: text.Font{
			Typeface: t.Font.Primary,
		},
		TxColor:  t.Color.Light,
		BgColor:  t.Color.Dark,
		TextSize: unit.Dp(float32(18)),
		shaper:   t.Shaper,
	}
}

func (c DuoUIcounter) Layout(gtx *layout.Context, cc *controller.DuoUIcounter, label string) {

	hmin := gtx.Constraints.Width.Min
	vmin := gtx.Constraints.Height.Min
	txColor := c.TxColor
	bgColor := c.BgColor
	layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func() {
			rr := float32(gtx.Px(unit.Dp(0)))
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(gtx.Constraints.Width.Min),
					Y: float32(gtx.Constraints.Height.Min),
				}},
				NE: rr, NW: rr, SE: rr, SW: rr,
			}.Op(gtx.Ops).Add(gtx.Ops)
			fill(gtx, HexARGB(bgColor))
			//for _, c := range button.History() {
			//	drawInk(gtx, c)
			//}
		}),
		layout.Stacked(func() {
			gtx.Constraints.Width.Min = hmin
			gtx.Constraints.Height.Min = vmin
			layout.Center.Layout(gtx, func() {
				layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: unit.Dp(12), Right: unit.Dp(12)}.Layout(gtx, func() {

					layout.Flex{}.Layout(gtx,
						layout.Rigid(func() {
							for cc.CounterDecrease.Clicked(gtx) {
								cc.Decrease()
								c.pageFunction()
							}
							c.decrease.IconLayout(gtx, cc.CounterDecrease)
						}),
						//layout.Flexed(0.2, func() {
						//	//for cc.CounterReset.Clicked(gtx) {
						//	//	cc.Reset()
						//	//	c.pageFunction()
						//	//}
						//	//c.reset.Layout(gtx, cc.CounterReset)
						//}),
						layout.Rigid(func() {
							for cc.CounterIncrease.Clicked(gtx) {
								cc.Increase()
								c.pageFunction()
							}
							c.increase.IconLayout(gtx, cc.CounterIncrease)
						}),
						layout.Rigid(func() {
							layout.Center.Layout(gtx, func() {
								paint.ColorOp{Color: HexARGB(c.TxColor)}.Add(gtx.Ops)
								controller.Label{
									Alignment: text.Middle,
								}.Layout(gtx, c.shaper, c.Font, unit.Dp(12), c.Text)
							})
						}))

					)
					paint.ColorOp{Color: HexARGB(txColor)}.Add(gtx.Ops)
					controller.Label{
						Alignment: text.Middle,
					}.Layout(gtx, c.shaper, c.Font, c.TextSize, c.Text)

				})
			})
			pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)

		}),
	)

}
