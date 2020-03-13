package component

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/p9c/pod/cmd/gui/rcd"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
)

var (
	previousBlockHashButton = new(gel.Button)
	nextBlockHashButton     = new(gel.Button)
	list                    = &layout.List{
		Axis: layout.Vertical,
	}
)

func UnoField(gtx *layout.Context, field func()) func() {
	return func() {
		layout.Flex{
			Axis:    layout.Horizontal,
			Spacing: layout.SpaceAround,
		}.Layout(gtx,
			layout.Flexed(1, field),
		)

	}
}
func DuoFields(gtx *layout.Context, axis layout.Axis, left, right func()) func() {
	return func() {
		layout.Flex{
			Axis:    axis,
			Spacing: layout.SpaceAround,
		}.Layout(gtx,
			fieldAxis(axis, left, 0.5),
			fieldAxis(axis, right, 0.5),
		)
	}
}

func TrioFields(gtx *layout.Context, th *gelook.DuoUItheme, axis layout.Axis, labelTextSize, valueTextSize float32, unoLabel, unoValue, unoHeadcolor, unoHeadbgColor, unoColor, unoBgColor, duoLabel, duoValue, duoHeadcolor, duoHeadbgColor, duoColor, duoBgColor, treLabel, treValue, treHeadcolor, treHeadbgColor, treColor, treBgColor string) func() {
	return func() {
		layout.Flex{
			Axis:    axis,
			Spacing: layout.SpaceAround,
		}.Layout(gtx,
			fieldAxis(axis, ContentLabeledField(gtx, th, layout.Vertical, 4, labelTextSize, valueTextSize, unoLabel, unoHeadcolor, unoHeadbgColor, unoColor, unoBgColor, fmt.Sprint(unoValue)), 0.3333),
			fieldAxis(axis, ContentLabeledField(gtx, th, layout.Vertical, 4, labelTextSize, valueTextSize, duoLabel, duoHeadcolor, duoHeadbgColor, duoColor, duoBgColor, fmt.Sprint(duoValue)), 0.3333),
			fieldAxis(axis, ContentLabeledField(gtx, th, layout.Vertical, 4, labelTextSize, valueTextSize, treLabel, treHeadbgColor, treHeadcolor, treColor, treBgColor, fmt.Sprint(treValue)), 0.3333),
		)
	}
}

func fieldAxis(axis layout.Axis, field func(), size float32) layout.FlexChild {
	var f layout.FlexChild
	switch axis {
	case layout.Horizontal:
		f = layout.Flexed(size, field)
	case layout.Vertical:
		f = layout.Rigid(field)
	}
	return f
}

func ContentLabeledField(gtx *layout.Context, th *gelook.DuoUItheme, axis layout.Axis, margin, labelTextSize, valueTextSize float32, label, headcolor, headbgColor, color, bgColor, value string) func() {
	return func() {
		layout.UniformInset(unit.Dp(margin)).Layout(gtx, func() {
			layout.Flex{
				Axis: axis,
			}.Layout(gtx,
				layout.Rigid(contentField(gtx, th, label, th.Colors[headcolor], th.Colors[headbgColor], th.Fonts["Primary"], 4, labelTextSize)),
				layout.Rigid(contentField(gtx, th, value, th.Colors[color], th.Colors[bgColor], th.Fonts["Mono"], 4, valueTextSize)))
		})
	}
}

func PageNavButtons(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, previousBlockHash, nextBlockHash string, prevPage, nextPage *gelook.DuoUIpage) func() {
	return func() {
		layout.Flex{}.Layout(gtx,
			layout.Flexed(0.5, func() {
				eh := chainhash.Hash{}
				if previousBlockHash != eh.String() {
					pageNavButton(rc, gtx, th, nextPage, previousBlockHashButton, "Previous Block", previousBlockHash)
				}
			}),
			layout.Flexed(0.5, func() {
				if nextBlockHash != "" {
					pageNavButton(rc, gtx, th, nextPage, nextBlockHashButton, "Next Block", nextBlockHash)
				}
			}))
	}
}

func pageNavButton(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, page *gelook.DuoUIpage, b *gel.Button, label, hash string) {
	layout.UniformInset(unit.Dp(4)).Layout(gtx, func() {
		var blockButton gelook.DuoUIbutton
		blockButton = th.DuoUIbutton(th.Fonts["Mono"], label+" "+hash, th.Colors["Light"], th.Colors["Info"], th.Colors["Info"], th.Colors["Light"], "", th.Colors["Light"], 16, 0, 60, 24, 0, 0)
		for b.Clicked(gtx) {
			rc.ShowPage = fmt.Sprintf("BLOCK %s", hash)
			rc.GetSingleBlock(hash)()
			SetPage(rc, page)
		}
		blockButton.Layout(gtx, b)
	})
}

func contentField(gtx *layout.Context, th *gelook.DuoUItheme, text, color, bgColor string, font text.Typeface, padding, textSize float32) func() {
	return func() {
		hmin := gtx.Constraints.Width.Min
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
				fill(gtx, gelook.HexARGB(bgColor))
			}),
			layout.Stacked(func() {
				gtx.Constraints.Width.Min = hmin
				gtx.Constraints.Height.Min = vmin
				layout.Center.Layout(gtx, func() {
					layout.UniformInset(unit.Dp(padding)).Layout(gtx, func() {
						l := th.DuoUIlabel(unit.Dp(textSize), text)
						l.Font.Typeface = font
						l.Color = gelook.HexARGB(color)
						l.Layout(gtx)
					})
				})
			}),
		)
	}
}
