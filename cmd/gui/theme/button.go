// SPDX-License-Identifier: Unlicense OR MIT

package theme

import (
	"github.com/p9c/pod/pkg/gio/io/pointer"
	"image"
	"image/color"

	"github.com/p9c/pod/cmd/gui/widget"
	"github.com/p9c/pod/pkg/gio/f32"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/op"
	"github.com/p9c/pod/pkg/gio/op/clip"
	"github.com/p9c/pod/pkg/gio/op/paint"
	"github.com/p9c/pod/pkg/gio/text"
	"github.com/p9c/pod/pkg/gio/unit"
)

type DuoUIbutton struct {
	Text string
	// Color is the text color.
	TxColor      color.RGBA
	Font         text.Font
	Width        float32
	Height       float32
	BgColor      color.RGBA
	CornerRadius unit.Value
	Icon         *DuoUIicon
	IconSize     unit.Value
	IconColor    color.RGBA
	Padding      unit.Value
	shaper       *text.Shaper
	hover        bool
}

func (t *DuoUItheme) DuoUIbutton(txt, txcolor, bgcolor string, width, height, iconSize, padding float32, icon *DuoUIicon) DuoUIbutton {
	return DuoUIbutton{
		Text: txt,
		Font: text.Font{
			Size: t.TextSize.Scale(14.0 / 16.0),
		},
		Width:    width,
		Height:   height,
		TxColor:  t.Color.Primary,
		BgColor:  t.Color.InvText,
		Icon:     icon,
		IconSize: unit.Dp(iconSize),
		Padding:  unit.Dp(padding),
		shaper:   t.Shaper,
	}
}

//
//func (b DuoUIbutton) Layout(gtx *layout.Context, button *widget.Button) {
//	layout.Stack{Alignment: layout.Center}.Layout(gtx,
//		layout.Expanded(func() {
//			rr := float32(gtx.Px(unit.Dp(0)))
//			clip.Rect{
//				Rect: f32.Rectangle{Max: f32.Point{
//					X: float32(b.Width),
//					Y: float32(b.Height),
//				}},
//				NE: rr, NW: rr, SE: rr, SW: rr,
//			}.Op(gtx.Ops).Add(gtx.Ops)
//			fill(gtx, b.BgColor)
//			for _, c := range button.History() {
//				drawInk(gtx, c)
//			}
//		}),
//		layout.Stacked(func() {
//			layout.Flex{
//				Axis:      layout.Vertical,
//				Alignment: layout.Middle,
//			}.Layout(gtx,
//				layout.Rigid(func() {
//					layout.UniformInset(b.Padding).Layout(gtx, func() {
//						size := gtx.Px(b.IconSize) - 2*gtx.Px(b.Padding)
//						if b.Icon != nil {
//							b.Icon.Color = b.TxColor
//							b.Icon.Layout(gtx, unit.Px(float32(size)))
//						}
//						gtx.Dimensions = layout.Dimensions{
//							Size: image.Point{X: size, Y: size},
//						}
//					})
//				}),
//				layout.Rigid(func() {
//					layout.Align(layout.Center).Layout(gtx, func() {
//						layout.Inset{Top: unit.Dp(0), Bottom: unit.Dp(0), Left: unit.Dp(0), Right: unit.Dp(0)}.Layout(gtx, func() {
//							paint.ColorOp{Color: b.TxColor}.Add(gtx.Ops)
//							widget.Label{}.Layout(gtx, b.shaper, b.Font, b.Text)
//						})
//					})
//				}),
//			)
//			pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
//			button.Layout(gtx)
//		}))
//}

func (b DuoUIbutton) Layout(gtx *layout.Context, button *widget.Button) {
	st := layout.Stack{Alignment: layout.Center}
	ico := layout.Stacked(func() {

		layout.UniformInset(b.Padding).Layout(gtx, func() {
			size := gtx.Px(b.IconSize) - 2*gtx.Px(b.Padding)

			layout.Flex{
				Axis:      layout.Vertical,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(func() {
					if b.Icon != nil {
						b.Icon.Color = b.IconColor
						b.Icon.Layout(gtx, unit.Px(float32(size)))
					}
					gtx.Dimensions = layout.Dimensions{
						Size: image.Point{X: size, Y: size},
					}
				}),
				layout.Rigid(func() {

				}),
			)
			pointer.Ellipse(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
			button.Layout(gtx)
		})
	})
	bg := layout.Expanded(func() {

	})
	st.Layout(gtx, bg, ico)
}

func toPointF(p image.Point) f32.Point {
	return f32.Point{X: float32(p.X), Y: float32(p.Y)}
}

func toRectF(r image.Rectangle) f32.Rectangle {
	return f32.Rectangle{
		Min: toPointF(r.Min),
		Max: toPointF(r.Max),
	}
}

func drawInk(gtx *layout.Context, c widget.Click) {
	d := gtx.Now().Sub(c.Time)
	t := float32(d.Seconds())
	const duration = 0.5
	if t > duration {
		return
	}
	t = t / duration
	var stack op.StackOp
	stack.Push(gtx.Ops)
	size := float32(gtx.Px(unit.Dp(700))) * t
	rr := size * .5
	col := byte(0xaa * (1 - t*t))
	ink := paint.ColorOp{Color: color.RGBA{A: col, R: col, G: col, B: col}}
	ink.Add(gtx.Ops)
	op.TransformOp{}.Offset(c.Position).Offset(f32.Point{
		X: -rr,
		Y: -rr,
	}).Add(gtx.Ops)
	clip.Rect{
		Rect: f32.Rectangle{Max: f32.Point{
			X: float32(size),
			Y: float32(size),
		}},
		NE: rr, NW: rr, SE: rr, SW: rr,
	}.Op(gtx.Ops).Add(gtx.Ops)
	paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{X: float32(size), Y: float32(size)}}}.Add(gtx.Ops)
	stack.Pop()
	op.InvalidateOp{}.Add(gtx.Ops)
}
