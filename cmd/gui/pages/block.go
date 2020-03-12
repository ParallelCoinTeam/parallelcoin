package pages

import (
	"encoding/json"
	"fmt"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

var (
	previousBlockHashButton                                = new(gel.Button)
	nextBlockHashButton                                    = new(gel.Button)
	algoHeadColor, algoHeadBgColor, algoColor, algoBgColor string
)

func blockPage(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, block string) *gelook.DuoUIpage {
	return th.DuoUIpage("BLOCK", 10, rc.GetSingleBlock(block), func() {}, singleBlockBody(rc, gtx, th, rc.Explorer.SingleBlock), func() {})
}

func singleBlockBody(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, block btcjson.GetBlockVerboseResult) func() {
	return func() {

		switch block.PowAlgo {
		case "argon2i":
			algoHeadColor = "Light"
			algoHeadBgColor = "Dark"
			algoColor = "Light"
			algoBgColor = "DarkGray"
		case "blake2b":
			algoHeadColor = "Light"
			algoHeadBgColor = "Dark"
			algoColor = "Dark"
			algoBgColor = "LightGrayI"
		case "x11":
			algoHeadColor = "Light"
			algoHeadBgColor = "Dark"
			algoColor = "Light"
			algoBgColor = "DarkGrayI"
		case "keccak":
			algoHeadColor = "Light"
			algoHeadBgColor = "Dark"
			algoColor = "Light"
			algoBgColor = "Secondary"
		case "blake3":
			algoHeadColor = "Light"
			algoHeadBgColor = "Dark"
			algoColor = "Dark"
			algoBgColor = "Success"
		case "scrypt":
			algoHeadColor = "Light"
			algoHeadBgColor = "Dark"
			algoColor = "DarkGrayI"
			algoBgColor = "Warning"
		case "sha256d":
			algoHeadColor = "Light"
			algoHeadBgColor = "Dark"
			algoColor = "LightGrayI"
			algoBgColor = "Info"
		case "skein":
			algoHeadColor = "Light"
			algoHeadBgColor = "Dark"
			algoColor = "Light"
			algoBgColor = "DarkGrayII"
		case "stribog":
			algoHeadColor = "Light"
			algoHeadBgColor = "Dark"
			algoColor = "Light"
			algoBgColor = "Danger"
		}

		duo := layout.Horizontal
		if gtx.Constraints.Width.Max < 1280 {
			duo = layout.Vertical
		}
		trio := layout.Horizontal
		if gtx.Constraints.Width.Max < 780 {
			trio = layout.Vertical
		}

		blockJSON, _ := json.MarshalIndent(block, "", "  ")
		blockText := string(blockJSON)
		widgets := []func(){

			component.UnoField(gtx, component.ContentLabeledField(gtx, th, layout.Vertical, 12, 18, "Hash", "Light", "Dark", "Light", "DarkGray", fmt.Sprint(block.Hash))),
			component.DuoFields(gtx, duo,
				component.TrioFields(gtx, th, trio, 12, 16,
					"Height", fmt.Sprint(block.Height), "Light", "Dark", "Light", "DarkGray",
					"Confirmations", fmt.Sprint(block.Confirmations), "Light", "Dark", "Light", "DarkGray",
					"Time", fmt.Sprint(time.Unix(block.Time, 0).Format("2006-01-02 15:04:05")), "Light", "Dark", "Light", "DarkGray",
				),
				component.TrioFields(gtx, th, trio, 12, 16,
					"PowAlgo", fmt.Sprint(block.PowAlgo), algoHeadColor, algoHeadBgColor, algoColor, algoBgColor,
					"Difficulty", fmt.Sprint(block.Difficulty), "Light", "Dark", "Light", "DarkGray",
					"Nonce", fmt.Sprint(block.Nonce), "Light", "Dark", "Light", "DarkGray",
				),
			),
			component.DuoFields(gtx, duo,
				component.ContentLabeledField(gtx, th, layout.Vertical, 12, 12, "MerkleRoot", "Light", "Dark", "Light", "DarkGray", block.MerkleRoot),
				component.ContentLabeledField(gtx, th, layout.Vertical, 12, 12, "PowHash", "Light", "Dark", "Light", "DarkGray", fmt.Sprint(block.PowHash)),
			),
			component.HorizontalLine(gtx, 1, th.Colors["Dark"]),
			component.DuoFields(gtx, duo,
				component.TrioFields(gtx, th, trio, 12, 16,
					"Size", fmt.Sprint(block.Size), "Light", "Dark", "Light", "DarkGray",
					"Weight", fmt.Sprint(block.Weight), "Light", "Dark", "Light", "DarkGray",
					"Bits", fmt.Sprint(block.Bits), "Light", "Dark", "Light", "DarkGray",
				),
				component.TrioFields(gtx, th, trio, 12, 16,
					"TxNum", fmt.Sprint(block.TxNum), "Light", "Dark", "Light", "DarkGray",
					"StrippedSize", fmt.Sprint(block.StrippedSize), "Light", "Dark", "Light", "DarkGray",
					"Version", fmt.Sprint(block.Version), "Light", "Dark", "Light", "DarkGray",
				),
			),
			component.UnoField(gtx, component.ContentLabeledField(gtx, th, layout.Vertical, 12, 12, "Tx", "Light", "Dark", "Light", "DarkGray", fmt.Sprint(block.Tx))),
			component.UnoField(gtx, component.ContentLabeledField(gtx, th, layout.Vertical, 12, 12, "RawTx", "Light", "Dark", "Light", "DarkGray", fmt.Sprint(blockText))),
			component.PageNavButtons(rc, gtx, th, block.PreviousHash, block.NextHash, blockPage(rc, gtx, th, block.PreviousHash), blockPage(rc, gtx, th, block.NextHash)),
		}
		layautList.Layout(gtx, len(widgets), func(i int) {
			layout.UniformInset(unit.Dp(4)).Layout(gtx, widgets[i])
		})

	}
}
