// SPDX-License-Identifier: Unlicense OR MIT

package gelook

import (
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"github.com/p9c/pod/pkg/gui/gel"
	"image"
	"image/color"

	"gioui.org/text"
	"gioui.org/unit"
)

type checkable struct {
	Label              string
	Color              color.RGBA
	Font               text.Font
	TextSize           unit.Value
	IconColor          color.RGBA
	Size               unit.Value
	shaper             text.Shaper
	checkedStateIcon   *DuoUIIcon
	uncheckedStateIcon *DuoUIIcon
	PillColor          string
	PillColorChecked   string
	CircleColor        string
	CircleColorChecked string
}

func (c *checkable) layout(gtx *layout.Context, checked bool) {

	var icon *DuoUIIcon
	if checked {
		icon = c.checkedStateIcon
	} else {
		icon = c.uncheckedStateIcon
	}

	hMin := gtx.Constraints.Width.Min
	vMin := gtx.Constraints.Height.Min
	layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func() {
			layout.Center.Layout(gtx, func() {
				layout.UniformInset(unit.Dp(2)).Layout(gtx, func() {
					size := gtx.Px(c.Size)
					icon.Color = c.IconColor
					icon.Layout(gtx, unit.Px(float32(size)))
					gtx.Dimensions = layout.Dimensions{
						Size: image.Point{X: size, Y: size},
					}
				})
			})
		}),

		layout.Rigid(func() {
			gtx.Constraints.Width.Min = hMin
			gtx.Constraints.Height.Min = vMin
			layout.W.Layout(gtx, func() {
				layout.UniformInset(unit.Dp(2)).Layout(gtx, func() {
					paint.ColorOp{Color: c.Color}.Add(gtx.Ops)
					gel.Label{}.Layout(gtx, c.shaper, c.Font, c.TextSize, c.Label)
				})
			})
		}),
	)
	pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
}

func (c *checkable) drawLayout(gtx *layout.Context, checked bool) {
	state := layout.W
	pillColor := c.PillColor
	circleColor := c.CircleColor
	if checked {
		state = layout.E
		pillColor = c.PillColorChecked
		circleColor = c.CircleColorChecked
	}
	hMin := gtx.Constraints.Width.Min
	vMin := gtx.Constraints.Height.Min
	layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func() {
			layout.Center.Layout(gtx, func() {
				layout.UniformInset(unit.Dp(2)).Layout(gtx, func() {
					layout.Center.Layout(gtx, func() {
						gtx.Constraints.Width.Min = 64
						gtx.Constraints.Height.Min = 32
						//DuoUIDrawRectangle(gtx, 64, 32, "ff888888", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
						layout.Center.Layout(gtx, func() {
							DuoUIDrawRectangle(gtx, 48, 16, pillColor, [4]float32{8, 8, 8, 8}, [4]float32{0, 0, 0, 0})
						})
						state.Layout(gtx, func() {
							DuoUIDrawRectangle(gtx, 24, 24, circleColor, [4]float32{12, 12, 12, 12}, [4]float32{0, 0, 0, 0})
						})
						pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
					})
					gtx.Dimensions = layout.Dimensions{
						Size: image.Point{X: 64, Y: 32},
					}
				})
			})
		}),

		layout.Rigid(func() {
			gtx.Constraints.Width.Min = hMin
			gtx.Constraints.Height.Min = vMin
			layout.E.Layout(gtx, func() {
				layout.UniformInset(unit.Dp(2)).Layout(gtx, func() {
					paint.ColorOp{Color: c.Color}.Add(gtx.Ops)
					gel.Label{}.Layout(gtx, c.shaper, c.Font, c.TextSize, c.Label)
				})
			})
		}),
	)
	pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
}
