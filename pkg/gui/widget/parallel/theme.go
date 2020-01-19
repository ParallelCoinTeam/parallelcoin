// SPDX-License-Identifier: Unlicense OR MIT

package parallel

import (
	"image"
	"image/color"

	"github.com/p9c/pod/pkg/gui/f32"
	"github.com/p9c/pod/pkg/gui/font"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/op/paint"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type DuoUItheme struct {
	Shaper text.Shaper
	Color  struct {
		Primary color.RGBA
		Text    color.RGBA
		Hint    color.RGBA
		InvText color.RGBA
	}
	TextSize              unit.Value
	checkBoxCheckedIcon   *DuoUIicon
	checkBoxUncheckedIcon *DuoUIicon
	radioCheckedIcon      *DuoUIicon
	radioUncheckedIcon    *DuoUIicon
}

func NewDuoUItheme() *DuoUItheme {
	t := &DuoUItheme{
		Shaper: font.Default(),
	}
	t.Color.Primary = rgb(0x3f51b5)
	t.Color.Text = rgb(0x303030)
	t.Color.Hint = rgb(0x888888)
	t.Color.InvText = rgb(0xcfcfcf)
	t.TextSize = unit.Sp(16)

	t.checkBoxCheckedIcon = mustIcon(NewDuoUIicon(icons.ToggleCheckBox))
	t.checkBoxUncheckedIcon = mustIcon(NewDuoUIicon(icons.ToggleCheckBoxOutlineBlank))
	t.radioCheckedIcon = mustIcon(NewDuoUIicon(icons.ToggleRadioButtonChecked))
	t.radioUncheckedIcon = mustIcon(NewDuoUIicon(icons.ToggleRadioButtonUnchecked))

	return t
}

func mustIcon(ic *DuoUIicon, err error) *DuoUIicon {
	if err != nil {
		panic(err)
	}
	return ic
}

func rgb(c uint32) color.RGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.RGBA {
	return color.RGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func fill(gtx *layout.Context, col color.RGBA) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d}
}
