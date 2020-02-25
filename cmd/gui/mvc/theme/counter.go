package theme

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"image/color"
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
		increase: t.DuoUIbutton(t.Font.Secondary, "INCREASE", t.Color.Primary, t.Color.Light, "", "", 12, 0, 64, 44, 0, 0),
		decrease: t.DuoUIbutton(t.Font.Secondary, "DECREASE", t.Color.Primary, t.Color.Light, "", "", 12, 0, 64, 44, 0, 0),
		reset:    t.DuoUIbutton(t.Font.Secondary, "RESET", t.Color.Primary, t.Color.Light, "", "", 12, 0, 64, 44, 0, 0),
		Font: text.Font{
			Typeface: t.Font.Primary,
		},
		TxColor:  HexARGB(t.Color.Light),
		TextSize: unit.Dp(float32(24)),
		shaper:   t.Shaper,
	}
}

func (c DuoUIcounter) Layout(gtx *layout.Context, cc *controller.DuoUIcounter) {

	layout.Flex{}.Layout(gtx,
		layout.Flexed(0.2, func() {
			paint.ColorOp{Color: c.TxColor}.Add(gtx.Ops)
			controller.Label{
				Alignment: text.Middle,
			}.Layout(gtx, c.shaper, c.Font, c.TextSize, fmt.Sprint(cc.Value))
		}),
		layout.Flexed(0.3, func() {
			for cc.CounterIncrease.Clicked(gtx) {
				cc.Increase()
			}
			c.increase.Layout(gtx, cc.CounterIncrease)
		}),
		layout.Flexed(0.2, func() {
			for cc.CounterReset.Clicked(gtx) {
				cc.Reset()
			}
			c.reset.Layout(gtx, cc.CounterReset)
		}),
		layout.Flexed(0.3, func() {
			for cc.CounterDecrease.Clicked(gtx) {
				cc.Decrease()
			}
			c.decrease.Layout(gtx, cc.CounterDecrease)
		}),
	)
}
