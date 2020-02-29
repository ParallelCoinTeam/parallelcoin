package pages

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/component"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

var (
	previousBlockHashButton = new(controller.Button)
	nextBlockHashButton     = new(controller.Button)
)

func blockPage(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, block string) *theme.DuoUIpage {
	return th.DuoUIpage("BLOCK", 10, rc.GetSingleBlock(block), func() {}, singleBlockBody(rc, gtx, th, rc.SingleBlock), func() {})
}

func singleBlockBody(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, block btcjson.GetBlockVerboseResult) func() {
	return func() {
		widgets := []func(){
			component.ContentLabeledField(gtx, th, layout.Vertical, 16, 24, "Hash", fmt.Sprint(block.Hash)),
			component.TrioFields(gtx, th, 16, 32,
				"Height", fmt.Sprint(block.Height),
				"Confirmations", fmt.Sprint(block.Confirmations),
				"Time", fmt.Sprint(block.Time)),
			component.ContentLabeledField(gtx, th, layout.Vertical, 16, 12, "MerkleRoot", block.MerkleRoot),
			component.TrioFields(gtx, th, 18, 16,
				"PowAlgo", fmt.Sprint(block.PowAlgo),
				"Difficulty", fmt.Sprint(block.Difficulty),
				"Nonce", fmt.Sprint(block.Nonce)),
			component.ContentLabeledField(gtx, th, layout.Vertical, 16, 12, "PowHash", fmt.Sprint(block.PowHash)),
			component.TrioFields(gtx, th, 16, 16,
				"Size", fmt.Sprint(block.Size),
				"Weight", fmt.Sprint(block.Weight),
				"Bits", fmt.Sprint(block.Bits)),
			component.HorizontalLine(gtx, 1, th.Color.Dark),
			component.TrioFields(gtx, th, 16, 16,
				"TxNum", fmt.Sprint(block.TxNum),
				"StrippedSize", fmt.Sprint(block.StrippedSize),
				"Version", fmt.Sprint(block.Version)),
			component.ContentLabeledField(gtx, th, layout.Vertical, 16, 12, "Tx", fmt.Sprint(block.Tx)),
			component.ContentLabeledField(gtx, th, layout.Vertical, 14, 12, "RawTx", fmt.Sprint(block.RawTx)),
			component.PageNavButtons(rc, gtx, th, block.PreviousHash, block.NextHash, blockPage(rc, gtx, th, block.PreviousHash), blockPage(rc, gtx, th, block.NextHash)),
		}
		layautList.Layout(gtx, len(widgets), func(i int) {
			layout.UniformInset(unit.Dp(4)).Layout(gtx, widgets[i])
		})

	}
}
