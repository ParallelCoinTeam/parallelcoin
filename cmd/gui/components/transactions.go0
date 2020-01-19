// SPDX-License-Identifier: Unlicense OR MIT

package components

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"image/color"

	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/text"
	"github.com/p9c/pod/pkg/gio/unit"
)

var (
	list = &layout.List{
		Axis: layout.Vertical,
	}
)

type DuoUItransactions struct {
	Text string
	// Color is the text color.
	TxColor      color.RGBA
	Font         text.Font
	Width        float32
	Height       float32
	BgColor      color.RGBA
	CornerRadius unit.Value
	//Icon              *DuoUIicon
	//IconSize          int
	//IconColor         color.RGBA
	//PaddingVertical   unit.Value
	//PaddingHorizontal unit.Value
	shaper *text.Shaper
	//transactions models.DuoUItransactionsExcerpts
}

func (t *DuoUItheme) DuoUItransactions(bgColor string, width, height, paddingVertical, paddingHorizontal float32) DuoUItransactions {
	return DuoUItransactions{
		Font: text.Font{
			Size: t.TextSize.Scale(14.0 / 16.0),
		},
		Width:  width,
		Height: height,
		//TxColor: helpers.HexARGB(txtColor),
		BgColor: helpers.HexARGB(bgColor),
		//PaddingVertical:   unit.Dp(paddingVertical),
		//PaddingHorizontal: unit.Dp(paddingHorizontal),
		shaper: t.Shaper,
		//transactions: transactions,
	}
}

func (txs DuoUItransactions) Layout(gtx *layout.Context, content func()) {

	in := layout.UniformInset(unit.Dp(30))
	in.Layout(gtx, func() {
		layout.Flex{
			//Axis:      0,
			//Spacing:   0,
			//Alignment: 0,
		}.Layout(gtx,
			// Balance status item
			layout.Rigid(func() {

				//const n = 5
				//list.Layout(gtx, n, func(i int) {
				//	txt := fmt.Sprintf("List element #%d", i)
				//
				//	duo.DuoUItheme.H3(txt).Layout(gtx)
				//})
				//transList := &layout.List{
				//	Axis: layout.Vertical,
				//}
				//transList.Layout(gtx, len(txs.transactions.Txs), func(i int) {
				//	// Invert list
				//	//i = len(txs.Txs) - 1 - i
				//	//t := table[i]
				//	//a := 1.0
				//	//const duration = 5
				//
				//	widgets := []func(){
				//		func() {
				//			//layout.Rigid(func() {
				//			//	tim := duo.DuoUItheme.Body1(t.TxID)
				//			//	tim.Color = helpers.Alpha(a, tim.Color)
				//			//	tim.Layout(gtx)
				//			//})
				//		},
				//		func() {
				//			//layout.Rigid(func() {
				//			//amount := duo.DuoUItheme.H5(fmt.Sprintf("%0.8f", t.Amount))
				//			//amount.Color = helpers.RGB(0x003300)
				//			//amount.Alignment = text.End
				//			//amount.Font.Variant = "Mono"
				//			//amount.Font.Weight = text.Bold
				//			//amount.Layout(gtx)
				//			//}),
				//		},
				//		func() {
				//			//layout.Rigid(func() {
				//			//sat := duo.DuoUItheme.Body1(t.Category)
				//			//sat.Color = helpers.Alpha(a, sat.Color)
				//			//sat.Layout(gtx)
				//			//}),
				//		},
				//		func() {
				//			//layout.Rigid(func() {
				//			//sat := duo.DuoUItheme.Body1(fmt.Sprintf("%0.8f", t.Amount))
				//			//sat.Color = helpers.Alpha(a, sat.Color)
				//			//sat.Layout(gtx)
				//			//}),
				//		},
				//		func() {
				//			//layout.Rigid(func() {
				//			//l := duo.DuoUItheme.Body2(helpers.FormatTime(t.Time))
				//			//l.Color = duo.DuoUItheme.Color.Hint
				//			//l.Color = helpers.Alpha(a, l.Color)
				//			//l.Layout(gtx)
				//			//})
				//		},
				//		func() {
				//
				//		},
				//	}
				//	list.Layout(gtx, len(widgets), func(i int) {
				//		layout.UniformInset(unit.Dp(16)).Layout(gtx, widgets[i])
				//	})

			}))
	})
}
