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

type DuoUIpage struct {
	Title   string
	TxColor string
	// Font          text.Font
	shaper  text.Shaper
	Command func()

	Header        func()
	HeaderBgColor string
	HeaderPadding float32
	// header
	// header
	Border      float32
	BorderColor string
	Body        func()
	BodyBgColor string
	BodyPadding float32
	// body
	// body
	Footer        func()
	FooterBgColor string
	FooterPadding float32
	// footer
	// footer
}

func (t *DuoUItheme) DuoUIpage(p DuoUIpage) *DuoUIpage {
	return &DuoUIpage{
		Title: p.Title,
		// Font: text.Font{
		// Size: t.TextSize.Scale(14.0 / 16.0),
		// },
		TxColor:       t.Colors["Dark"],
		shaper:        t.Shaper,
		Command:       p.Command,
		Header:        p.Header,
		HeaderBgColor: t.Colors["Primary"],
		HeaderPadding: p.HeaderPadding,
		Border:        p.Border,
		BorderColor:   p.BorderColor,
		Body:          p.Body,
		BodyBgColor:   p.BodyBgColor,
		BodyPadding:   p.BodyPadding,
		Footer:        p.Footer,
		FooterBgColor: t.Colors["Secondary"],
		FooterPadding: p.FooterPadding,
	}
}

func (p DuoUIpage) Layout(gtx *layout.Context) {
	layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(pageElementLayout(gtx, layout.N, p.HeaderBgColor, p.HeaderPadding, p.Header)),
		layout.Flexed(1, func() {
			DuoUIfill(gtx, p.BorderColor)
			layout.UniformInset(unit.Dp(p.Border)).Layout(gtx, pageElementLayout(gtx, layout.N, p.BodyBgColor, p.BodyPadding, p.Body))
		}),
		layout.Rigid(pageElementLayout(gtx, layout.N, p.FooterBgColor, p.FooterPadding, p.Footer)),
	)
}

func pageElementLayout(gtx *layout.Context, direction layout.Direction, background string, padding float32, elementContent func()) func() {
	return func() {
		hmin := gtx.Constraints.Width.Max
		vmin := gtx.Constraints.Height.Min
		layout.Stack{Alignment: layout.W}.Layout(gtx,
			layout.Expanded(func() {
				rr := float32(gtx.Px(unit.Dp(0)))
				clip.Rect{
					Rect: f32.Rectangle{Max: f32.Point{
						X: float32(gtx.Constraints.Width.Min),
						Y: float32(gtx.Constraints.Height.Min),
					}},
					NE: rr, NW: rr, SE: rr, SW: rr,
				}.Op(gtx.Ops).Add(gtx.Ops)
				fill(gtx, HexARGB(background))
				pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
			}),
			layout.Stacked(func() {
				gtx.Constraints.Width.Min = hmin
				gtx.Constraints.Height.Min = vmin
				direction.Layout(gtx, func() {
					layout.Flex{}.Layout(gtx,
						layout.Flexed(1, func() {
							layout.Inset{
								Top:    unit.Dp(padding),
								Right:  unit.Dp(padding),
								Bottom: unit.Dp(padding),
								Left:   unit.Dp(padding),
							}.Layout(gtx, elementContent)
						}))
				})
			}),
		)
	}
}
