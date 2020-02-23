package theme

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/op/paint"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
	"image/color"
)

var (
	couterIncrease = new(controller.Button)
	couterDecrease = new(controller.Button)
	couterReset    = new(controller.Button)
)

type DuoUIcounter struct {
	increase DuoUIbutton
	decrease DuoUIbutton
	reset    DuoUIbutton
	Font     text.Font
	TextSize unit.Value
	TxColor  color.RGBA
	BgColor  color.RGBA
	shaper   text.Shaper
}

func (t *DuoUItheme) DuoUIcounter() DuoUIcounter {
	return DuoUIcounter{
		increase: t.DuoUIbutton(t.Font.Secondary, "INCREASE", t.Color.Primary, t.Color.Light, "", "", 16, 0, 128, 48, 0, 0),
		decrease: t.DuoUIbutton(t.Font.Secondary, "DECREASE", t.Color.Primary, t.Color.Light, "", "", 16, 0, 128, 48, 0, 0),
		reset:    t.DuoUIbutton(t.Font.Secondary, "RESET", t.Color.Primary, t.Color.Light, "", "", 16, 0, 128, 48, 0, 0),
		Font: text.Font{
			Typeface: t.Font.Primary,
		},
		TxColor:  HexARGB(t.Color.Light),
		TextSize: unit.Dp(float32(16)),
		shaper:   t.Shaper,
	}
}

func (c DuoUIcounter) Layout(gtx *layout.Context, cc *controller.DuoUIcounter) {
	layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceSides,
	}.Layout(gtx,
		layout.Flexed(0.5, func() {
			layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceSides,
			}.Layout(gtx,
				layout.Flexed(0.5, func() {

					layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func() {
							in := layout.UniformInset(unit.Dp(0))
							in.Layout(gtx, func() {
								paint.ColorOp{Color: c.TxColor}.Add(gtx.Ops)
								controller.Label{
									Alignment: text.Middle,
								}.Layout(gtx, c.shaper, c.Font, c.TextSize, fmt.Sprint(cc.Value))
							})
						}),
						layout.Flexed(1, func() {
							layout.Flex{}.Layout(gtx,
								layout.Flexed(0.4, func() {
									for couterIncrease.Clicked(gtx) {
										cc.Increase()
									}
									c.increase.Layout(gtx, couterIncrease)
								}),
								layout.Flexed(0.2, func() {
									for couterReset.Clicked(gtx) {
										cc.Reset()
									}
									c.reset.Layout(gtx, couterReset)
								}),
								layout.Flexed(0.4, func() {
									for couterDecrease.Clicked(gtx) {
										cc.Decrease()
									}
									c.decrease.Layout(gtx, couterDecrease)
								}),
							)
						}),
					)

				}),
			)
		}),
	)
}
