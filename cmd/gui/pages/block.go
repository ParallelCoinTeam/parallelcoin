package pages

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/cmd/gui/controller"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/cmd/gui/theme"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

var (
	previousBlockHashButton = new(controller.Button)
	nextBlockHashButton     = new(controller.Button)
)

func blockPage(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, block string) *theme.DuoUIpage {
	return th.DuoUIpage("BLOCK", 10, rc.GetSingleBlock(block), func() {}, singleBlockBody(rc, gtx, th, rc.Explorer.SingleBlock), func() {})
}

func singleBlockBody(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, block btcjson.GetBlockVerboseResult) func() {
	return func() {

		duo := layout.Horizontal
		if gtx.Constraints.Width.Max < 1280 {
			duo = layout.Vertical
		}
		trio := layout.Horizontal
		if gtx.Constraints.Width.Max < 780 {
			trio = layout.Vertical
		}

		widgets := []func(){

			component.UnoField(gtx, component.ContentLabeledField(gtx, th, layout.Vertical, 16, 24, "Hash", fmt.Sprint(block.Hash))),
			component.DuoFields(gtx, duo,
				component.TrioFields(gtx, th, trio, 16, 32,
					"Height", fmt.Sprint(block.Height),
					"Confirmations", fmt.Sprint(block.Confirmations),
					"Time", fmt.Sprint(block.Time)),
				component.TrioFields(gtx, th, trio, 18, 16,
					"PowAlgo", fmt.Sprint(block.PowAlgo),
					"Difficulty", fmt.Sprint(block.Difficulty),
					"Nonce", fmt.Sprint(block.Nonce)),
			),
			component.DuoFields(gtx, duo,
				component.ContentLabeledField(gtx, th, layout.Vertical, 16, 12, "MerkleRoot", block.MerkleRoot),
				component.ContentLabeledField(gtx, th, layout.Vertical, 16, 12, "PowHash", fmt.Sprint(block.PowHash)),
			),
			component.HorizontalLine(gtx, 1, th.Color.Dark),
			component.DuoFields(gtx, duo,
				component.TrioFields(gtx, th, trio, 16, 16,
					"Size", fmt.Sprint(block.Size),
					"Weight", fmt.Sprint(block.Weight),
					"Bits", fmt.Sprint(block.Bits)),
				component.TrioFields(gtx, th, trio, 16, 16,
					"TxNum", fmt.Sprint(block.TxNum),
					"StrippedSize", fmt.Sprint(block.StrippedSize),
					"Version", fmt.Sprint(block.Version)),
			),
			component.UnoField(gtx, component.ContentLabeledField(gtx, th, layout.Vertical, 16, 12, "Tx", fmt.Sprint(block.Tx))),
			component.UnoField(gtx, component.ContentLabeledField(gtx, th, layout.Vertical, 14, 12, "RawTx", fmt.Sprint(block.RawTx))),
			component.PageNavButtons(rc, gtx, th, block.PreviousHash, block.NextHash, blockPage(rc, gtx, th, block.PreviousHash), blockPage(rc, gtx, th, block.NextHash)),
		}
		layautList.Layout(gtx, len(widgets), func(i int) {
			layout.UniformInset(unit.Dp(4)).Layout(gtx, widgets[i])
		})

	}
}
