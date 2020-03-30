// SPDX-License-Identifier: Unlicense OR MIT

package gelook

import (
	"gioui.org/op"
	"image"
	"image/color"

	"github.com/p9c/pod/pkg/logi"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/p9c/pod/pkg/gel"
)

type Button struct {
	Text string
	// Color is the text color.
	Color        color.RGBA
	Font         text.Font
	TextSize     unit.Value
	Background   color.RGBA
	CornerRadius unit.Value
	Inset        layout.Inset
	shaper       text.Shaper
}

type ButtonLayout struct {
	Background   color.RGBA
	Color        color.RGBA
	CornerRadius unit.Value
	Inset        layout.Inset
}

type IconButton struct {
	Background color.RGBA
	Color      color.RGBA
	Icon       *DuoUIicon
	Size       unit.Value
	Padding    unit.Value
	Inset      layout.Inset
}

func (t *DuoUItheme) Button(txt string) Button {
	return Button{
		Text:         txt,
		Color:        HexARGB(t.Colors["ButtonText"]),
		CornerRadius: unit.Dp(4),
		Background:   HexARGB(t.Colors["ButtonBg"]),
		TextSize:     t.TextSize.Scale(14.0 / 16.0),
		Inset: layout.Inset{
			Top: unit.Dp(10), Bottom: unit.Dp(10),
			Left: unit.Dp(12), Right: unit.Dp(12),
		},
		shaper: t.Shaper,
	}
}

func (t *DuoUItheme) ButtonLayout() ButtonLayout {
	return ButtonLayout{
		Background:   HexARGB(t.Colors["ButtonBg"]),
		Color:        HexARGB(t.Colors["ButtonText"]),
		CornerRadius: unit.Dp(4),
		Inset:        layout.UniformInset(unit.Dp(12)),
	}
}

func (t *DuoUItheme) IconButton(icon *DuoUIicon) IconButton {
	return IconButton{
		Background: HexARGB(t.Colors["Primary"]),
		Color:      HexARGB(t.Colors["InvText"]),
		Icon:       icon,
		Size:       unit.Dp(56),
		Padding:    unit.Dp(16),
	}
}

func (b Button) Layout(gtx *layout.Context, button *gel.Button) {
	ButtonLayout{
		Background:   b.Background,
		CornerRadius: b.CornerRadius,
		Color:        b.Color,
		Inset:        b.Inset,
	}.Layout(gtx, button, func() {
		gel.Label{}.Layout(gtx, b.shaper, b.Font, b.TextSize, b.Text)
	})
}

func (b ButtonLayout) Layout(gtx *layout.Context, button *gel.Button, w layout.Widget) {
	hmin := gtx.Constraints.Width.Min
	vmin := gtx.Constraints.Height.Min
	layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func() {
			rr := float32(gtx.Px(b.CornerRadius))
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(gtx.Constraints.Width.Min),
					Y: float32(gtx.Constraints.Height.Min),
				}},
				NE: rr, NW: rr, SE: rr, SW: rr,
			}.Op(gtx.Ops).Add(gtx.Ops)
			fill(gtx, b.Background)
			for _, c := range button.History() {
				drawInk(gtx, c)
			}
		}),
		layout.Stacked(func() {
			layout.Center.Layout(gtx, func() {
				gtx.Constraints.Width.Min = hmin
				gtx.Constraints.Height.Min = vmin
				b.Inset.Layout(gtx, func() {
					paint.ColorOp{Color: b.Color}.Add(gtx.Ops)
					w()
				})
			})
			pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
			button.Layout(gtx)
		}),
	)
}

func (b IconButton) Layout(gtx *layout.Context, button *gel.Button) {
	layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func() {
			size := gtx.Constraints.Width.Min
			sizef := float32(size)
			rr := sizef * .5
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{X: sizef, Y: sizef}},
				NE:   rr, NW: rr, SE: rr, SW: rr,
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
			pointer.Ellipse(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
			button.Layout(gtx)
		}),
	)
}

var (
	buttonInsideLayoutList = &layout.List{
		Axis: layout.Vertical,
	}
)

type DuoUIbutton struct {
	Text string
	// Color is the text color.
	TxColor       color.RGBA
	BgColor       color.RGBA
	TxColorHover  color.RGBA
	BgColorHover  color.RGBA
	Font          text.Font
	TextSize      unit.Value
	Width         int
	Height        int
	CornerRadius  unit.Value
	Icon          *DuoUIicon
	IconSize      int
	IconColor     color.RGBA
	PaddingTop    unit.Value
	PaddingRight  unit.Value
	PaddingBottom unit.Value
	PaddingLeft   unit.Value
	shaper        text.Shaper
	hover         bool
}

func (t *DuoUItheme) DuoUIbutton(
	txtFont text.Typeface, txt,
	txtColor, bgColor,
	txtHoverColor, bgHoverColor,
	icon, iconColor string,
	textSize, iconSize,
	width, height,
	paddingTop,
	paddingRight,
	paddingBottom,
	paddingLeft int) DuoUIbutton {
	return DuoUIbutton{
		Text: txt,
		Font: text.Font{
			Typeface: txtFont,
		},
		TextSize:      unit.Dp(float32(textSize)),
		Width:         width,
		Height:        height,
		TxColor:       HexARGB(txtColor),
		BgColor:       HexARGB(bgColor),
		TxColorHover:  HexARGB(txtHoverColor),
		BgColorHover:  HexARGB(bgHoverColor),
		Icon:          t.Icons[icon],
		IconSize:      iconSize,
		IconColor:     HexARGB(iconColor),
		PaddingTop:    unit.Dp(float32(paddingTop)),
		PaddingRight:  unit.Dp(float32(paddingRight)),
		PaddingBottom: unit.Dp(float32(paddingBottom)),
		PaddingLeft:   unit.Dp(float32(paddingLeft)),
		shaper:        t.Shaper,
	}
}

func (b DuoUIbutton) Layout(gtx *layout.Context, button *gel.Button) {
	hmin := gtx.Constraints.Width.Min
	vmin := gtx.Constraints.Height.Min
	if b.Height > 0 {
		//vmin = b.Height
	}
	txColor := b.TxColor
	bgColor := b.BgColor
	if button.Hover(gtx) {
		txColor = b.TxColorHover
		bgColor = b.BgColorHover
		logi.L.Info("")
		logi.L.Info("oce")
		logi.L.Info("")
	}
	layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func() {
			rr := float32(gtx.Px(unit.Dp(0)))
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(gtx.Constraints.Width.Min),
					Y: float32(gtx.Constraints.Height.Min),
				}},
				NE: rr, NW: rr, SE: rr, SW: rr,
			}.Op(gtx.Ops).Add(gtx.Ops)
			fill(gtx, bgColor)
			for _, c := range button.History() {
				drawInk(gtx, c)
			}
		}),
		layout.Stacked(func() {
			gtx.Constraints.Width.Min = hmin
			gtx.Constraints.Height.Min = vmin
			layout.Center.Layout(gtx, func() {
				layout.Inset{Top: b.PaddingTop, Bottom: b.PaddingBottom, Left: b.PaddingLeft, Right: b.PaddingRight}.Layout(gtx, func() {

					paint.ColorOp{Color: txColor}.Add(gtx.Ops)
					gel.Label{
						Alignment: text.Middle,
					}.Layout(gtx, b.shaper, b.Font, b.TextSize, b.Text)
				})
			})
			pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
			button.Layout(gtx)
		}),
	)
}

func (b DuoUIbutton) IconLayout(gtx *layout.Context, button *gel.Button) {
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
			fill(gtx, b.BgColor)
			for _, c := range button.History() {
				drawInk(gtx, c)
			}
		}),
		layout.Stacked(func() {
			gtx.Constraints.Width.Min = b.Width
			gtx.Constraints.Height.Min = b.Height
			layout.Center.Layout(gtx, func() {
				layout.Inset{Top: b.PaddingTop, Bottom: b.PaddingBottom, Left: b.PaddingLeft, Right: b.PaddingRight}.Layout(gtx, func() {
					b.Icon.Color = b.IconColor
					b.Icon.Layout(gtx, unit.Dp(float32(b.IconSize)-b.PaddingTop.V-b.PaddingBottom.V))
				})
				gtx.Dimensions = layout.Dimensions{
					Size: image.Point{X: b.IconSize, Y: b.IconSize},
				}
			})
			pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
			button.Layout(gtx)
		}),
	)
}

func (b DuoUIbutton) MenuLayout(gtx *layout.Context, button *gel.Button) {
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
			fill(gtx, b.BgColor)
			for _, c := range button.History() {
				drawInk(gtx, c)
			}
		}),
		layout.Stacked(func() {
			gtx.Constraints.Width.Min = b.Width
			gtx.Constraints.Height.Min = b.Height
			layout.Center.Layout(gtx, func() {
				layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func() {
						layout.Center.Layout(gtx, func() {
							layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
								b.Icon.Color = b.IconColor
								b.Icon.Layout(gtx, unit.Dp(float32(b.IconSize)))
							})
							gtx.Dimensions = layout.Dimensions{
								Size: image.Point{X: b.IconSize, Y: b.IconSize},
							}
						})
					}),
					layout.Rigid(func() {
						layout.Center.Layout(gtx, func() {
							paint.ColorOp{Color: b.TxColor}.Add(gtx.Ops)
							gel.Label{
								Alignment: text.Middle,
							}.Layout(gtx, b.shaper, b.Font, unit.Dp(12), b.Text)
						})
					}))
			})
			pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
			button.Layout(gtx)
		}),
	)
}

func (b DuoUIbutton) InsideLayout(gtx *layout.Context, button *gel.Button, inside func()) {
	hmin := gtx.Constraints.Width.Min
	vmin := gtx.Constraints.Height.Min
	bgColor := b.BgColor
	layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func() {
			rr := float32(gtx.Px(unit.Dp(0)))
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(gtx.Constraints.Width.Min),
					Y: float32(gtx.Constraints.Height.Min),
				}},
				NE: rr, NW: rr, SE: rr, SW: rr,
			}.Op(gtx.Ops).Add(gtx.Ops)
			fill(gtx, bgColor)
			for _, c := range button.History() {
				drawInk(gtx, c)
			}
		}),
		layout.Stacked(func() {
			gtx.Constraints.Width.Min = hmin
			gtx.Constraints.Height.Min = vmin
			layout.Center.Layout(gtx, func() {
				layout.Inset{Top: b.PaddingTop, Bottom: b.PaddingBottom, Left: b.PaddingLeft, Right: b.PaddingRight}.Layout(gtx, func() {
					inside()
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

func drawInk(gtx *layout.Context, c gel.Click) {
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
