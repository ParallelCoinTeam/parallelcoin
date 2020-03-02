package component

import (
	"fmt"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/cmd/gui/rcd"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
)

var (
	previousBlockHashButton = new(controller.Button)
	nextBlockHashButton     = new(controller.Button)
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
			layout.Flexed(0.99, field),
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

func TrioFields(gtx *layout.Context, th *theme.DuoUItheme, axis layout.Axis, labelTextSize, valueTextSize float32, unoLabel, unoValue, duoLabel, duoValue, treLabel, treValue string) func() {
	return func() {
		layout.Flex{
			Axis:    axis,
			Spacing: layout.SpaceAround,
		}.Layout(gtx,
			fieldAxis(axis, ContentLabeledField(gtx, th, layout.Vertical, labelTextSize, valueTextSize, unoLabel, fmt.Sprint(unoValue)), 0.3),
			fieldAxis(axis, ContentLabeledField(gtx, th, layout.Vertical, labelTextSize, valueTextSize, duoLabel, fmt.Sprint(duoValue)), 0.3),
			fieldAxis(axis, ContentLabeledField(gtx, th, layout.Vertical, labelTextSize, valueTextSize, treLabel, fmt.Sprint(treValue)), 0.3),
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

func ContentLabeledField(gtx *layout.Context, th *theme.DuoUItheme, axis layout.Axis, labelTextSize, valueTextSize float32, label, value string) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			layout.Flex{
				Axis: axis,
			}.Layout(gtx,
				layout.Rigid(contentField(gtx, th, label, th.Color.Light, th.Color.Dark, th.Font.Primary, labelTextSize)),
				layout.Rigid(contentField(gtx, th, value, th.Color.Light, th.Color.DarkGray, th.Font.Mono, valueTextSize)))
		})
	}
}

func PageNavButtons(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, previousBlockHash, nextBlockHash string, prevPage, nextPage *theme.DuoUIpage) func() {
	return func() {
		layout.Flex{}.Layout(gtx,
			layout.Flexed(0.495, func() {
				eh := chainhash.Hash{}
				if previousBlockHash != eh.String() {
					var previousBlockButton theme.DuoUIbutton
					previousBlockButton = th.DuoUIbutton(th.Font.Mono, "Previous Block "+previousBlockHash, th.Color.Light, th.Color.Info, th.Color.Info, th.Color.Light, "", th.Color.Light, 16, 0, 60, 24, 0, 0)
					for previousBlockHashButton.Clicked(gtx) {
						// clipboard.Set(b.BlockHash)
						rc.ShowPage = fmt.Sprintf("BLOCK %s", previousBlockHash)
						rc.GetSingleBlock(previousBlockHash)()
						SetPage(rc, prevPage)
					}
					previousBlockButton.Layout(gtx, previousBlockHashButton)
				}
			}),
			layout.Flexed(0.495, func() {
				if nextBlockHash != "" {
					var nextBlockButton theme.DuoUIbutton
					nextBlockButton = th.DuoUIbutton(th.Font.Mono, "Next Block "+nextBlockHash, th.Color.Light, th.Color.Info, th.Color.Info, th.Color.Light, "", th.Color.Light, 16, 0, 60, 24, 0, 0)
					for nextBlockHashButton.Clicked(gtx) {
						// clipboard.Set(b.BlockHash)
						rc.ShowPage = fmt.Sprintf("BLOCK %s", nextBlockHash)
						rc.GetSingleBlock(nextBlockHash)()
						SetPage(rc, nextPage)
					}
					nextBlockButton.Layout(gtx, nextBlockHashButton)
				}
			}))
	}
}

func contentField(gtx *layout.Context, th *theme.DuoUItheme, text, color, bgColor string, font text.Typeface, textSize float32) func() {
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
				fill(gtx, theme.HexARGB(bgColor))
			}),
			layout.Stacked(func() {
				gtx.Constraints.Width.Min = hmin
				gtx.Constraints.Height.Min = vmin
				layout.Center.Layout(gtx, func() {
					layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
						l := th.DuoUIlabel(unit.Dp(textSize), text)
						l.Font.Typeface = font
						l.Color = theme.HexARGB(color)
						l.Layout(gtx)
					})
				})
			}),
		)
	}
}
