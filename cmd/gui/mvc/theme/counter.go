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

//type button struct {
//	pressed      bool
//	Name         string
//	ColorBg      string
//	BorderRadius [4]float32
//	OperateValue interface{}
//	Font         text.Font
//	TxColor      color.RGBA
//	TextSize     unit.Value
//	shaper       text.Shaper
//}
//
//func (b *button) Layout(gtx *layout.Context, action func(interface{})) {
//	for _, e := range gtx.Events(b) {
//		if e, ok := e.(pointer.Event); ok {
//			switch e.Type {
//			case pointer.Press:
//				b.pressed = true
//				action(b.OperateValue)
//				log.INFO("Itwoikoos")
//
//			case pointer.Release:
//				b.pressed = false
//			}
//		}
//	}
//
//	cs := gtx.Constraints
//	if b.pressed {
//	}
//	pointer.Rect(
//		image.Rectangle{Max: image.Point{X: cs.Width.Min, Y: cs.Height.Min}},
//	).Add(gtx.Ops)
//	pointer.InputOp{Key: b}.Add(gtx.Ops)
//	paint.ColorOp{Color: b.TxColor}.Add(gtx.Ops)
//	controller.Label{
//		Alignment: text.Middle,
//	}.Layout(gtx, b.shaper, b.Font, b.TextSize, b.Name)
//}

type DuoUIcounter struct {
	increase *DuoUIbutton
	decrease *DuoUIbutton
	reset    *DuoUIbutton
	Font     text.Font
	TextSize unit.Value
	TxColor  color.RGBA
	BgColor  color.RGBA
	shaper   text.Shaper
}

func (t *DuoUItheme) DuoUIcounter() DuoUIcounter {
	return DuoUIcounter{
		increase: &DuoUIbutton{
			Text:         "increase",
			BgColor:      HexARGB(t.Color.Light),
			CornerRadius: unit.Value{},
			Font: text.Font{
				Typeface: t.Font.Primary,
			},
			TxColor:  HexARGB(t.Color.Primary),
			TextSize: unit.Dp(float32(16)),
			shaper:   t.Shaper,
		},
		decrease: &DuoUIbutton{
			Text:         "decrease",
			BgColor:      HexARGB(t.Color.Light),
			CornerRadius: unit.Value{},
			Font: text.Font{
				Typeface: t.Font.Primary,
			},
			TxColor:  HexARGB(t.Color.Primary),
			TextSize: unit.Dp(float32(16)),
			shaper:   t.Shaper,
		},
		reset: &DuoUIbutton{
			Text:         "reset",
			BgColor:      HexARGB(t.Color.Light),
			CornerRadius: unit.Value{},
			Font: text.Font{
				Typeface: t.Font.Primary,
			},
			TxColor:  HexARGB(t.Color.Primary),
			TextSize: unit.Dp(float32(16)),
			shaper:   t.Shaper,
		},
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
							//cs := gtx.Constraints

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
									for couterDecrease.Clicked(gtx) {
										cc.Decrease()
									}
									c.decrease.Layout(gtx, couterDecrease)
								}),
								layout.Flexed(0.4, func() {
									for couterReset.Clicked(gtx) {
										cc.Reset()
									}
									c.reset.Layout(gtx, couterReset)
								}),
							)
						}),
					)

				}),
			)
		}),
	)
}
