// SPDX-License-Identifier: Unlicense OR MIT

package parallel

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/pkg/gui/io/pointer"
	"image"
	"image/color"

	"github.com/p9c/pod/pkg/gui/widget"
	"github.com/p9c/pod/pkg/gui/f32"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/op"
	"github.com/p9c/pod/pkg/gui/op/clip"
	"github.com/p9c/pod/pkg/gui/op/paint"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
)

var (
	buttonLayoutList = &layout.List{
		Axis: layout.Vertical,
	}
)

type DuoUIbutton struct {
	Text string
	// Color is the text color.
	TxColor           color.RGBA
	Font              text.Font
	Width             float32
	Height            float32
	BgColor           color.RGBA
	CornerRadius      unit.Value
	Icon              *DuoUIicon
	IconSize          int
	IconColor         color.RGBA
	PaddingVertical   unit.Value
	PaddingHorizontal unit.Value
	shaper            text.Shaper
	hover             bool
}

func (t *DuoUItheme) DuoUIbutton(txt, txtColor, bgColor, iconColor string, iconSize int, width, height, paddingVertical, paddingHorizontal float32, icon *DuoUIicon) DuoUIbutton {
	return DuoUIbutton{
		Text: txt,
		Font: text.Font{
			Size: t.TextSize.Scale(8.0 / 10.0),
		},
		Width:             width,
		Height:            height,
		TxColor:           helpers.HexARGB(txtColor),
		BgColor:           helpers.HexARGB(bgColor),
		Icon:              icon,
		IconSize:          iconSize,
		IconColor:             helpers.HexARGB(iconColor),
		PaddingVertical:   unit.Dp(paddingVertical),
		PaddingHorizontal: unit.Dp(paddingHorizontal),
		shaper:            t.Shaper,
	}
}

func (b DuoUIbutton) Layout(gtx *layout.Context, button *widget.Button) {
	col := b.TxColor
	bgcol := b.BgColor
	layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func() {
			rr := float32(gtx.Px(unit.Dp(0)))
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(b.Width),
					Y: float32(b.Height),
				}},
				NE: rr, NW: rr, SE: rr, SW: rr,
			}.Op(gtx.Ops).Add(gtx.Ops)
			fill(gtx, bgcol)
			for _, c := range button.History() {
				drawInk(gtx, c)
			}
		}),
		layout.Stacked(func() {
			gtx.Constraints.Width.Min = int(b.Width)
			gtx.Constraints.Height.Min = int(b.Height)
			layout.Align(layout.Center).Layout(gtx, func() {
				buttonLayout := []func(){
					func() {
						if b.Icon != nil {
							layout.Inset{Top: b.PaddingVertical, Bottom: b.PaddingVertical, Left: b.PaddingHorizontal, Right: b.PaddingHorizontal}.Layout(gtx, func() {
								if b.Icon != nil {
									b.Icon.Color = b.IconColor
									b.Icon.Layout(gtx, unit.Px(float32(b.IconSize)))
								}
								gtx.Dimensions = layout.Dimensions{
									Size: image.Point{X: b.IconSize, Y: b.IconSize},
								}
							})
						}
					},
					func() {
						if b.Text != "" {
							//layout.Inset{Top: unit.Dp(0), Bottom: unit.Dp(0), Left: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, func() {
								paint.ColorOp{Color: col}.Add(gtx.Ops)
								widget.Label{
									Alignment:text.Middle,
								}.Layout(gtx, b.shaper, b.Font, b.Text)
							//})
						}
					},
				}
				buttonLayoutList.Layout(gtx, len(buttonLayout), func(i int) {
					layout.UniformInset(unit.Dp(0)).Layout(gtx, buttonLayout[i])
				})
			})
			pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
			button.Layout(gtx)
		}),
	)
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
