package duoui

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

var (
	blocksList = &layout.List{
		Axis: layout.Vertical,
	}
	perPage = &controller.DuoUIcounter{
		Value:           20,
		OperateValue:    1,
		From:            0,
		To:              50,
		CounterIncrease: new(controller.Button),
		CounterDecrease: new(controller.Button),
		CounterReset:    new(controller.Button),
	}
	page = &controller.DuoUIcounter{
		Value:           0,
		OperateValue:    1,
		From:            0,
		To:              50,
		CounterIncrease: new(controller.Button),
		CounterDecrease: new(controller.Button),
		CounterReset:    new(controller.Button),
	}
)

func bodyExplorer(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		//rc.GetBlocksExcerpts(page.Value, perPage.Value)
		in := layout.UniformInset(unit.Dp(0))
		in.Layout(gtx, func() {
			blocksList.Layout(gtx, len(rc.Blocks), func(i int) {
				b := rc.Blocks[i]
				blockRow(rc, gtx, th, &b)
			})
		})
	}
}

func headerExplorer(gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		layout.Flex{
			Spacing: layout.SpaceBetween,
			Axis:    layout.Horizontal,
		}.Layout(gtx,
			//layout.Rigid(ui.txsFilter()),
			layout.Flexed(0.5, func() {
				th.DuoUIcounter().Layout(gtx, page)
			}),
			layout.Flexed(0.5, func() {
				th.DuoUIcounter().Layout(gtx, perPage)
			}),
		)
	}
}

func blockRow(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, block *model.DuoUIblock) {
	line(gtx, th.Color.Dark)()
	layout.Flex{
		Spacing: layout.SpaceBetween,
	}.Layout(gtx,
		layout.Rigid(func() {
			layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceAround,
			}.Layout(gtx,
				layout.Flexed(0.6, func() {
					var linkButton theme.DuoUIbutton
					linkButton = th.DuoUIbutton(th.Font.Mono, fmt.Sprint(block.Height), th.Color.Light, th.Color.Dark, "", th.Color.Light, 16, 0, 60, 24, 0, 0)
					for block.Link.Clicked(gtx) {
						//clipboard.Set(b.BlockHash)
						rc.ShowPage = "BLOCK" + block.BlockHash
						rc.GetSingleBlock(block.BlockHash)()
						setPage(rc, blockPage(rc, gtx, th, block.BlockHash))
					}
					linkButton.Layout(gtx, block.Link)
				}),
				layout.Rigid(func() {
					amount := th.H5(fmt.Sprintf("%0.8f", block.Amount))
					amount.Font.Typeface = th.Font.Primary
					amount.Color = theme.HexARGB(th.Color.Dark)
					amount.Alignment = text.End
					amount.Font.Variant = "Mono"
					amount.Font.Weight = text.Bold
					amount.Layout(gtx)
				}),
				layout.Rigid(func() {
					sat := th.Body1(fmt.Sprint(block.TxNum))
					sat.Font.Typeface = th.Font.Primary
					sat.Color = theme.HexARGB(th.Color.Dark)
					sat.Layout(gtx)
				}),
				layout.Rigid(func() {
					sat := th.Body1(fmt.Sprint(block.BlockHash))
					sat.Font.Typeface = th.Font.Mono
					sat.Color = theme.HexARGB(th.Color.Dark)
					sat.Layout(gtx)
				}),
				layout.Rigid(func() {
					l := th.Body2(block.Time)
					l.Font.Typeface = th.Font.Primary
					l.Color = theme.HexARGB(th.Color.Dark)
					l.Layout(gtx)
				}),
			)
		}),
		layout.Rigid(func() {
			sat := th.Body1(fmt.Sprintf("%0.8f", block.Amount))
			sat.Color = theme.HexARGB(th.Color.Dark)
			sat.Layout(gtx)
		}))
}

func blockPage(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, block string) *theme.DuoUIpage {
	return th.DuoUIpage("BLOCK", 0, rc.GetSingleBlock(block), func() {}, singleBlockBody(gtx, th, rc.SingleBlock), func() {})
}

func singleBlockBody(gtx *layout.Context, th *theme.DuoUItheme, block btcjson.GetBlockVerboseResult) func() {
	return func() {
		line(gtx, th.Color.Dark)()
		layout.Flex{
			Spacing: layout.SpaceBetween,
		}.Layout(gtx,
			layout.Rigid(func() {

				widgets := []func(){
					blockField(gtx, th, "Confirmations:", fmt.Sprint(block.Confirmations)),
					blockField(gtx, th, "TxNum:", fmt.Sprint(block.TxNum)),
					blockField(gtx, th, "Hash:", fmt.Sprint(block.Hash)),
					blockField(gtx, th, "Time:", fmt.Sprint(block.Time)),
					blockField(gtx, th, "StrippedSize:", fmt.Sprint(block.StrippedSize)),
					blockField(gtx, th, "Size:", fmt.Sprint(block.Size)),
					blockField(gtx, th, "Weight:", fmt.Sprint(block.Weight)),
					blockField(gtx, th, "Height:", fmt.Sprint(block.Height)),
					blockField(gtx, th, "Version:", fmt.Sprint(block.Version)),
					blockField(gtx, th, "VersionHex:", fmt.Sprint(block.VersionHex)),
					blockField(gtx, th, "PowAlgoID:", fmt.Sprint(block.PowAlgoID)),
					blockField(gtx, th, "PowAlgo:", fmt.Sprint(block.PowAlgo)),
					blockField(gtx, th, "PowHash:", fmt.Sprint(block.PowHash)),
					blockField(gtx, th, "MerkleRoot:", block.MerkleRoot),
					blockField(gtx, th, "Tx:", fmt.Sprint(block.Tx)),
					blockField(gtx, th, "RawTx:", fmt.Sprint(block.RawTx)),
					blockField(gtx, th, "Time:", fmt.Sprint(block.Time)),
					blockField(gtx, th, "Nonce:", fmt.Sprint(block.Nonce)),
					blockField(gtx, th, "Bits:", block.Bits),
					blockField(gtx, th, "Difficulty:", fmt.Sprint(block.Difficulty)),
					blockField(gtx, th, "PreviousHash:", block.PreviousHash),
					blockField(gtx, th, "NextHash:", block.NextHash),
				}
				layautList.Layout(gtx, len(widgets), func(i int) {
					layout.UniformInset(unit.Dp(8)).Layout(gtx, widgets[i])
				})
			}),
			layout.Rigid(func() {
				//sat := th.Body1(fmt.Sprintf("%0.8f", block.NextHash))
				sat := th.Body1(block.NextHash)
				sat.Color = theme.HexARGB(th.Color.Dark)
				sat.Font.Typeface = th.Font.Mono
				sat.Layout(gtx)
			}))
	}
}

func blockField(gtx *layout.Context, th *theme.DuoUItheme, label, value string) func() {
	return func() {
		layout.Flex{}.Layout(gtx,
			layout.Rigid(stackedField(gtx, th, label, th.Color.Success, th.Color.Dark, th.Font.Primary)),
			layout.Flexed(1, blockFieldValue(gtx, th, value)))
	}
}

func stackedField(gtx *layout.Context, th *theme.DuoUItheme, text, color, bgColor string, font text.Typeface) func() {
	return func() {
		width := 120
		height := 40
		layout.Stack{Alignment: layout.Center}.Layout(gtx,
			layout.Expanded(func() {
				rr := float32(gtx.Px(unit.Dp(0)))
				clip.Rect{
					Rect: f32.Rectangle{Max: f32.Point{
						X: float32(width),
						Y: float32(height),
					}},
					NE: rr, NW: rr, SE: rr, SW: rr,
				}.Op(gtx.Ops).Add(gtx.Ops)
				fill(gtx, theme.HexARGB(bgColor))
			}),
			layout.Stacked(func() {
				gtx.Constraints.Width.Min = width
				gtx.Constraints.Height.Min = height
				layout.Center.Layout(gtx, func() {
					//paint.ColorOp{Color: b.TxColor}.Add(gtx.Ops)
					//controller.Label{
					//	Alignment: text.Middle,
					//}.Layout(gtx, b.shaper, b.Font, unit.Dp(12), b.Text)
					l := th.Body2(text)
					l.Font.Typeface = font
					l.Color = theme.HexARGB(color)
					l.Layout(gtx)
				})
			}),
		)
	}

}
func blockFieldValue(gtx *layout.Context, th *theme.DuoUItheme, value string) func() {
	return func() {
		l := th.Body2(value)
		l.Font.Typeface = th.Font.Mono
		l.Color = theme.HexARGB(th.Color.Dark)
		l.Layout(gtx)
	}

}
