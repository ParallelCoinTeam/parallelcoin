// SPDX-License-Identifier: Unlicense OR MIT

package parallel

import (
	"image/color"

	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
)

var (


)

type DuoUIpage struct {
	Text string
	TxColor      color.RGBA
	Font         text.Font
	BgColor      color.RGBA
	//Icon              *DuoUIicon
	//IconSize          int
	//IconColor         color.RGBA
	//PaddingVertical   unit.Value
	//PaddingHorizontal unit.Value
	shaper text.Shaper
}

func (t *DuoUItheme) DuoUIpage(txtColor, bgColor string, paddingVertical, paddingHorizontal float32) DuoUIpage {
	return DuoUIpage{
		Font: text.Font{
			Size: t.TextSize.Scale(14.0 / 16.0),
		},
		TxColor: HexARGB(txtColor),
		BgColor: HexARGB(bgColor),
		//PaddingVertical:   unit.Dp(paddingVertical),
		//PaddingHorizontal: unit.Dp(paddingHorizontal),
		shaper: t.Shaper,
	}
}

func (b DuoUIpage) Layout(gtx *layout.Context, content func()) {
	layout.Flex{}.Layout(gtx,
		layout.Flexed(1, func() {
			in := layout.UniformInset(unit.Dp(0))
			in.Layout(gtx, func() {
				cs := gtx.Constraints
				DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, HexARGB("ffacacac"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				// Overview <<<
				layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
					cs := gtx.Constraints
					DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, HexARGB("ffcfcfcf"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
					content()
				})
				// Overview >>>
			})
		}),
	)
}
