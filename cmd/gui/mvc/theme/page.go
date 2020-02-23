// SPDX-License-Identifier: Unlicense OR MIT

package theme

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
)

var ()

type DuoUIpage struct {
	Text    string
	TxColor color.RGBA
	Font    text.Font
	BgColor color.RGBA
	//Icon              *DuoUIicon
	//IconSize          int
	//IconColor         color.RGBA
	//PaddingVertical   unit.Value
	//PaddingHorizontal unit.Value
	shaper text.Shaper
	layout func()
}

func (t *DuoUItheme) DuoUIpage(txt, txtColor, bgColor string, paddingVertical, paddingHorizontal float32, f func()) *DuoUIpage {
	return &DuoUIpage{
		Text: txt,
		Font: text.Font{
			//Size: t.TextSize.Scale(14.0 / 16.0),
		},
		TxColor: HexARGB(txtColor),
		BgColor: HexARGB(bgColor),
		//PaddingVertical:   unit.Dp(paddingVertical),
		//PaddingHorizontal: unit.Dp(paddingHorizontal),
		shaper: t.Shaper,
		layout: f,
	}
}

func (p DuoUIpage) Layout(gtx *layout.Context) {
	layout.Flex{}.Layout(gtx,
		layout.Flexed(1, func() {
			in := layout.UniformInset(unit.Dp(0))
			in.Layout(gtx, func() {
				cs := gtx.Constraints
				DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, "ffacacac", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				// Overview <<<
				layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
					cs := gtx.Constraints
					DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, "ffcfcfcf", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
					p.layout()
				})
				// Overview >>>
			})
		}),
	)
}
