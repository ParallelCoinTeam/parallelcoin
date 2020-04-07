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
	checkedStateIcon   *DuoUIicon
	uncheckedStateIcon *DuoUIicon
	PillColor          string
	PillColorChecked   string
	CircleColor        string
	CircleColorChecked string
}

func (c *checkable) layout(gtx *layout.Context, checked bool) {

	var icon *DuoUIicon
	if checked {
		icon = c.checkedStateIcon
	} else {
		icon = c.uncheckedStateIcon
	}

	hmin := gtx.Constraints.Width.Min
	vmin := gtx.Constraints.Height.Min
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
			gtx.Constraints.Width.Min = hmin
			gtx.Constraints.Height.Min = vmin
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
	state := layout.E
	pillColor := c.PillColor
	circleColor := c.CircleColor
	if checked {
		state = layout.W
		pillColor = c.PillColorChecked
		circleColor = c.CircleColorChecked
	}
	hmin := gtx.Constraints.Width.Min
	vmin := gtx.Constraints.Height.Min
	layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func() {
			layout.Center.Layout(gtx, func() {
				layout.UniformInset(unit.Dp(2)).Layout(gtx, func() {
					layout.Center.Layout(gtx, func() {
						gtx.Constraints.Width.Min = 64
						gtx.Constraints.Height.Min = 32
						//DuoUIdrawRectangle(gtx, 64, 32, "ff888888", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
						layout.Center.Layout(gtx, func() {
							DuoUIdrawRectangle(gtx, 48, 16, pillColor, [4]float32{8, 8, 8, 8}, [4]float32{0, 0, 0, 0})
						})
						state.Layout(gtx, func() {
							DuoUIdrawRectangle(gtx, 24, 24, circleColor, [4]float32{12, 12, 12, 12}, [4]float32{0, 0, 0, 0})
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
			gtx.Constraints.Width.Min = hmin
			gtx.Constraints.Height.Min = vmin
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
