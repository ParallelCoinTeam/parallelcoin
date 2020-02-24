// SPDX-License-Identifier: Unlicense OR MIT

package theme

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
)

var ()

type DuoUIpage struct {
	Title   string
	TxColor string
	Font    text.Font
	BgColor string
	//Icon              *DuoUIicon
	//IconSize          int
	//IconColor         color.RGBA
	//PaddingVertical   unit.Value
	//PaddingHorizontal unit.Value
	shaper text.Shaper
	layout func()
}

func (t *DuoUItheme) DuoUIpage(txt string, paddingVertical, paddingHorizontal float32, f func()) *DuoUIpage {
	return &DuoUIpage{
		Title: txt,
		Font:  text.Font{
			//Size: t.TextSize.Scale(14.0 / 16.0),
		},
		TxColor: t.Color.Dark,
		BgColor: t.Color.Light,
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
					DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, p.BgColor, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
					p.layout()
				})
				// Overview >>>
			})
		}),
	)
}
