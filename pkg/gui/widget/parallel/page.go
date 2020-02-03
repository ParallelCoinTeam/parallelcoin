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
	Width        float32
	Height       float32
	BgColor      color.RGBA
	CornerRadius unit.Value
	//Icon              *DuoUIicon
	//IconSize          int
	//IconColor         color.RGBA
	//PaddingVertical   unit.Value
	//PaddingHorizontal unit.Value
	shaper text.Shaper
}

func (t *DuoUItheme) DuoUIpage(txtColor, bgColor string, width, height, paddingVertical, paddingHorizontal float32) DuoUIpage {
	return DuoUIpage{
		Font: text.Font{
			Size: t.TextSize.Scale(14.0 / 16.0),
		},
		Width:   width,
		Height:  height,
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
			cs := gtx.Constraints
			in := layout.UniformInset(unit.Dp(0))
			in.Layout(gtx, func() {

				DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, HexARGB("ffacacac"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				// Overview <<<
				in := layout.UniformInset(unit.Dp(1))
				in.Layout(gtx, func() {
					DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, HexARGB("ffcfcfcf"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
					content()
				})
				// Overview >>>
			})
		}),
	)
}
