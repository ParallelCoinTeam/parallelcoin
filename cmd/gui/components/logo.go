// SPDX-License-Identifier: Unlicense OR MIT

package components

import (
	"image"
	"image/color"

	"github.com/p9c/pod/pkg/gio/f32"
	"github.com/p9c/pod/pkg/gio/io/pointer"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/op/clip"
	"github.com/p9c/pod/pkg/gio/unit"
	"github.com/p9c/pod/pkg/gio/widget"
)



type DuoUIlogo struct {
	Background color.RGBA
	Color      color.RGBA
	Icon       *DuoUIicon
	Size       unit.Value
	Padding    unit.Value
}

func (t *DuoUItheme) DuoUIiconLogo(icon *DuoUIicon) DuoUIlogo {
	return DuoUIlogo{
		Background: t.Color.Primary,
		Color:      t.Color.InvText,
		Icon:       icon,
		Size:       unit.Dp(64),
		Padding:    unit.Dp(16),
	}
}

func (b DuoUIlogo) Layout(gtx *layout.Context, button *widget.Button) {
	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			size := float32(gtx.Constraints.Width.Min)
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{X: size, Y: size}},
			}.Op(gtx.Ops).Add(gtx.Ops)
			fill(gtx, b.Background)
			for _, c := range button.History() {
				drawInk(gtx, c)
			}
		}),
		layout.Stacked(func() {
			layout.UniformInset(b.Padding).Layout(gtx, func() {
				size := gtx.Px(b.Size) - 2*gtx.Px(b.Padding)
				if b.Icon != nil {
					b.Icon.Color = b.Color
					b.Icon.Layout(gtx, unit.Px(float32(size)))
				}
				gtx.Dimensions = layout.Dimensions{
					Size: image.Point{X: size, Y: size},
				}
			})
			//pointer.Ellipse(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
			pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
			button.Layout(gtx)
		}),
	)
}
