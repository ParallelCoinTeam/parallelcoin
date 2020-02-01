// SPDX-License-Identifier: Unlicense OR MIT

package parallel

import (
	"github.com/p9c/pod/pkg/gui/text"
	"image/color"
)

var (
	increase = &button{
		Name:         "increase",
		OperateValue: 1,
	}
	decrease = &button{
		Name:         "decrease",
		OperateValue: 1,
	}
	reset = &button{
		Name:         "reset",
		OperateValue: 0,
	}
	itemValue = item{
		i: 5,
	}
)

type button struct {
	pressed      bool
	Name         string
	Do           func(interface{})
	ColorBg      string
	BorderRadius [4]float32
	OperateValue interface{}
}
type item struct {
	i int
}

func (it *item) doIncrease(n int) {
	it.i = it.i + int(n)
}

func (it *item) doDecrease(n int) {
	it.i = it.i - int(n)
}
func (it *item) doReset() {
	it.i = 0
}

type DuoUIcounter struct {
	Font text.Font
	// Color is the text color.
	Color color.RGBA
	// Hint contains the text displayed when the editor is empty.
	Hint string
	// HintColor is the color of hint text.
	HintColor color.RGBA
	Text      string
	shaper    text.Shaper
}

//func (t *DuoUItheme) DuoUIcounter(hint, txt string) DuoUIeditor {
////	//
////
////	layout.Stack{}.Layout(duo.DuoUIcontext,
////		layout.Stacked(func() {
////			layout.Flex{}.Layout(duo.DuoUIcontext,
////				layout.Flexed(0.4, func() {
////					decrease.Do = func(n interface{}) {
////						itemValue.doDecrease(n.(int))
////					}
////					decrease.Layout(duo)
////				}),
////				layout.Flexed(0.2, func() {
////					layout.Flex{Axis: layout.Horizontal}.Layout(duo.DuoUIcontext,
////						layout.Rigid(func() {
////							//cs := duo.DuoUIcontext.Constraints
////							//helpers.DrawRectangle(duo.DuoUIcontext, cs.Width.Max, 120, "ff3030cf", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
////							in := layout.UniformInset(unit.Dp(0))
////							in.Layout(duo.DuoUIcontext, func() {
////								layout.Align(layout.Center).Layout(duo.DuoUIcontext, func() {
////									duo.DuoUItheme.Body2(fmt.Sprint(itemValue.i)).Layout(duo.DuoUIcontext)
////								})
////							})
////						}),
////						layout.Flexed(1, func() {
////							reset.Do = func(interface{}) {
////								itemValue.doReset()
////							}
////							reset.Layout(duo)
////						}),
////					)
////				}),
////				layout.Flexed(0.4, func() {
////					increase.Do = func(n interface{}) {
////						itemValue.doIncrease(n.(int))
////					}
////					increase.Layout(duo)
////				}),
////			)
////		}),
////	)
////}
//
//func (b *button) Layout(duo *models.DuoUI) {
//	for _, e := range duo.DuoUIcontext.Events(b) { // HLevent
//		if e, ok := e.(pointer.Event); ok { // HLevent
//			switch e.Type { // HLevent
//			case pointer.Press: // HLevent
//				b.pressed = true // HLevent
//				b.Do(b.OperateValue)
//			case pointer.Release: // HLevent
//				b.pressed = false // HLevent
//			}
//		}
//	}
//
//	cs := duo.DuoUIcontext.Constraints
//	colorBg := "ff30cfcf"
//	colorBorder := "ffcf3030"
//	border := unit.Dp(1)
//
//	if b.pressed {
//		colorBg = "ffcf30cf"
//		colorBorder = "ff303030"
//		border = unit.Dp(3)
//	}
//	pointer.Rect( // HLevent
//		image.Rectangle{Max: image.Point{X: cs.Width.Max, Y: cs.Height.Max}}, // HLevent
//	).Add(duo.DuoUIcontext.Ops)                       // HLevent
//	pointer.InputOp{Key: b}.Add(duo.DuoUIcontext.Ops) // HLevent
//	//helpers.DrawRectangle(gtx, cs.Width.Max, cs.Height.Max, colorBorder, b.BorderRadius, unit.Dp(0))
//	helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, 32, colorBorder, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
//
//	in := layout.UniformInset(border)
//	in.Layout(duo.DuoUIcontext, func() {
//		//helpers.DrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, colorBg, b.BorderRadius, unit.Dp(0))
//		helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, 30, colorBg, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
//		//cs := gtx.Constraints
//		label := duo.DuoUItheme.Caption(b.Name)
//		label.Alignment = text.Middle
//		label.Layout(duo.DuoUIcontext)
//	})
//}
