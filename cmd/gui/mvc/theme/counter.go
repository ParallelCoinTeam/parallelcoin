// SPDX-License-Identifier: Unlicense OR MIT

package theme

import (
	"github.com/p9c/pod/pkg/gui/io/pointer"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
	"image"
	"image/color"
)

var (
	buttonLayoutList = &layout.List{
		Axis: layout.Vertical,
	}
)

type button struct {
	pressed      bool
	Name         string
	Do           func(interface{})
	Font         text.Font
	TxColor      color.RGBA
	ColorBg      string
	BorderRadius [4]float32
	OperateValue interface{}
	shaper       text.Shaper
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
	increase *button
	decrease *button
	reset    *button
	item     *item
	Font     text.Font
	// Color is the text color.
	Color color.RGBA
	// Hint contains the text displayed when the editor is empty.
	Hint string
	// HintColor is the color of hint text.
	HintColor color.RGBA
	// Color is the text color.
	TxColor           color.RGBA
	Width             float32
	Height            float32
	BgColor           color.RGBA
	CornerRadius      unit.Value
	PaddingVertical   unit.Value
	PaddingHorizontal unit.Value
	shaper            text.Shaper
	hover             bool
}

func (t *DuoUItheme) DuoUIcounter(txtFont text.Typeface, txtColor, bgColor string, width, height, paddingVertical, paddingHorizontal float32) DuoUIcounter {
	itemValue := new(item)
	return DuoUIcounter{
		increase: &button{
			Name: "increase",
			Do: func(n interface{}) {
				itemValue.doIncrease(n.(int))
			},
			OperateValue: 1,
			TxColor:      HexARGB(txtColor),
			shaper:       t.Shaper,
		},
		decrease: &button{
			Name: "decrease",
			Do: func(n interface{}) {
				itemValue.doDecrease(n.(int))
			},
			OperateValue: 1,
			TxColor:      HexARGB(txtColor),
			shaper:       t.Shaper,
		},
		reset: &button{
			Name: "reset",
			Do: func(interface{}) {
				itemValue.doReset()
			},
			OperateValue: 0,
			TxColor:      HexARGB(txtColor),
			shaper:       t.Shaper,
		},
		item: itemValue,
		Font: text.Font{
			Typeface: txtFont,
			//Size:     t.TextSize.Scale(8.0 / 10.0),
		},
		Width:             width,
		Height:            height,
		TxColor:           HexARGB(txtColor),
		BgColor:           HexARGB(bgColor),
		PaddingVertical:   unit.Dp(paddingVertical),
		PaddingHorizontal: unit.Dp(paddingHorizontal),
		shaper:            t.Shaper,
	}
}

func (c DuoUIcounter) Layout(gtx *layout.Context, th *DuoUItheme) {
	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func() {
			cs := gtx.Constraints
			DuoUIdrawRectangle(gtx, cs.Width.Max, 20, "ff3030cf", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			//theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, cs.Height.Max, ui.ly.Theme.Color.Dark, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			in := layout.UniformInset(unit.Dp(0))
			in.Layout(gtx, func() {
				//layout.Align(layout.Center).Layout(gtx, func() {
				//	//th.H3(fmt.Sprint(c.itemValue.i)).Layout(gtx)
				//
				//	paint.ColorOp{Color: c.TxColor}.Add(gtx.Ops)
				//	controller.Label{
				//		Alignment: text.Middle,
				//	}.Layout(gtx, c.shaper, c.Font, fmt.Sprint(c.item.i))
				//
				//})
			})
		}),
		layout.Rigid(func() {
			layout.Flex{}.Layout(gtx,
				layout.Flexed(0.4, func() {
					c.increase.Layout(gtx)
				}),
				layout.Flexed(0.2, func() {
					c.reset.Layout(gtx)
				}),
				layout.Flexed(0.4, func() {
					c.decrease.Layout(gtx)
				}),
			)
		}),
	)
}

func (b *button) Layout(gtx *layout.Context) {
	for _, e := range gtx.Events(b) { // HLevent
		if e, ok := e.(pointer.Event); ok { // HLevent
			switch e.Type { // HLevent
			case pointer.Press: // HLevent
				b.pressed = true // HLevent
				b.Do(b.OperateValue)
			case pointer.Release: // HLevent
				b.pressed = false // HLevent
			}
		}
	}
	cs := gtx.Constraints
	colorBg := "ff303030"
	colorBorder := "ffcf3030"
	border := unit.Dp(1)

	if b.pressed {
		colorBg = "ffcf30cf"
		colorBorder = "ff303030"
		border = unit.Dp(3)
	}
	pointer.Rect( // HLevent
		image.Rectangle{Max: image.Point{X: cs.Width.Max, Y: cs.Height.Max}}, // HLevent
	).Add(gtx.Ops)                       // HLevent
	pointer.InputOp{Key: b}.Add(gtx.Ops) // HLevent
	DuoUIdrawRectangle(gtx, cs.Width.Max, 22, colorBorder, b.BorderRadius, [4]float32{0, 0, 0, 0})
	in := layout.UniformInset(border)
	in.Layout(gtx, func() {
		DuoUIdrawRectangle(gtx, cs.Width.Max, 20, colorBg, b.BorderRadius, [4]float32{0, 0, 0, 0})
		//cs := gtx.Constraints
		//layout.Align(layout.Center).Layout(gtx, func() {
		//	paint.ColorOp{Color: b.TxColor}.Add(gtx.Ops)
		//	controller.Label{
		//		Alignment: text.Middle,
		//	}.Layout(gtx, b.shaper, b.Font, b.Name)
		//})
	})
}
