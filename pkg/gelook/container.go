// SPDX-License-Identifier: Unlicense OR MIT

package gelook

import (
	"image"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/text"
	"gioui.org/unit"
)

type DuoUIcontainer struct {
	// Color is the text color.
	Color         string
	Font          text.Font
	TextSize      unit.Value
	Background    string
	TxColor       string
	BgColor       string
	TxColorHover  string
	BgColorHover  string
	FullWidth     bool
	Width         int
	Height        int
	CornerRadius  int
	PaddingTop    int
	PaddingRight  int
	PaddingBottom int
	PaddingLeft   int
	shaper        text.Shaper
	link          bool
	hover         bool
}

func (t *DuoUItheme) DuoUIcontainer(padding int, background string) DuoUIcontainer {
	return DuoUIcontainer{
		Font: text.Font{
			Typeface: t.Fonts["Primary"],
		},
		// Color:      rgb(0xffffff),
		PaddingTop:    padding,
		PaddingRight:  padding,
		PaddingBottom: padding,
		PaddingLeft:   padding,
		Background:    background,
		TextSize:      t.TextSize.Scale(14.0 / 16.0),
		shaper:        t.Shaper,
	}
}

func (d DuoUIcontainer) Layout(gtx *layout.Context, direction layout.Direction, itemContent func()) {
	hmin := gtx.Constraints.Width.Min
	vmin := gtx.Constraints.Height.Min
	if d.FullWidth {
		hmin = gtx.Constraints.Width.Max
	}
	layout.Stack{Alignment: layout.W}.Layout(gtx,
		layout.Expanded(func() {
			rr := float32(gtx.Px(unit.Dp(float32(d.CornerRadius))))
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(gtx.Constraints.Width.Min),
					Y: float32(gtx.Constraints.Height.Min),
				}},
				NE: rr, NW: rr, SE: rr, SW: rr,
			}.Op(gtx.Ops).Add(gtx.Ops)
			fill(gtx, HexARGB(d.Background))
			pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
		}),
		layout.Stacked(func() {
			gtx.Constraints.Width.Min = hmin
			gtx.Constraints.Height.Min = vmin
			direction.Layout(gtx, func() {
				layout.Inset{
					Top:    unit.Dp(float32(d.PaddingTop)),
					Right:  unit.Dp(float32(d.PaddingRight)),
					Bottom: unit.Dp(float32(d.PaddingBottom)),
					Left:   unit.Dp(float32(d.PaddingLeft)),
				}.Layout(gtx, itemContent)
			})
		}),
	)
}
