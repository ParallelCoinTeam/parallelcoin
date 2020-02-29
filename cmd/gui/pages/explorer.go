package pages

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/component"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"image"
	"image/color"
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
	previousBlockHashButton = new(controller.Button)
	nextBlockHashButton     = new(controller.Button)
)

func Explorer(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) *theme.DuoUIpage {
	return th.DuoUIpage("EXPLORER", 0, rc.GetBlocksExcerpts(page.Value, perPage.Value), component.ContentHeader(gtx, th, headerExplorer(gtx, th)), bodyExplorer(rc, gtx, th), func() {})
}
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
	layout.UniformInset(unit.Dp(4)).Layout(gtx, func() {

		component.HorizontalLine(gtx, 1, th.Color.Dark)()
		layout.Flex{
			Spacing: layout.SpaceBetween,
		}.Layout(gtx,
			layout.Rigid(func() {
				var linkButton theme.DuoUIbutton
				linkButton = th.DuoUIbutton(th.Font.Mono, fmt.Sprint(block.Height), th.Color.Light, th.Color.Info, "", th.Color.Light, 14, 0, 60, 24, 0, 0)
				for block.Link.Clicked(gtx) {
					//clipboard.Set(b.BlockHash)
					rc.ShowPage = "BLOCK" + block.BlockHash
					rc.GetSingleBlock(block.BlockHash)()
					component.SetPage(rc, blockPage(rc, gtx, th, block.BlockHash))
				}
				linkButton.Layout(gtx, block.Link)
			}),
			layout.Rigid(func() {
				l := th.Body2(block.BlockHash)
				l.Font.Typeface = th.Font.Mono
				l.Color = theme.HexARGB(th.Color.Dark)
				l.Layout(gtx)
			}))
	})
}

func blockPage(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, block string) *theme.DuoUIpage {
	return th.DuoUIpage("BLOCK", 10, rc.GetSingleBlock(block), func() {}, singleBlockBody(rc, gtx, th, rc.SingleBlock), func() {})
}

func singleBlockBody(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, block btcjson.GetBlockVerboseResult) func() {
	return func() {

		widgets := []func(){
			blockField(gtx, th, layout.Vertical, 16, 24, "Hash", fmt.Sprint(block.Hash)),
			trioFields(gtx, th, 16, 32,
				"Height", fmt.Sprint(block.Height),
				"Confirmations", fmt.Sprint(block.Confirmations),
				"Time", fmt.Sprint(block.Time)),
			blockField(gtx, th, layout.Vertical, 16, 12, "MerkleRoot", block.MerkleRoot),
			trioFields(gtx, th, 18, 16,
				"PowAlgo", fmt.Sprint(block.PowAlgo),
				"Difficulty", fmt.Sprint(block.Difficulty),
				"Nonce", fmt.Sprint(block.Nonce)),
			blockField(gtx, th, layout.Vertical, 16, 12, "PowHash", fmt.Sprint(block.PowHash)),
			trioFields(gtx, th, 16, 16,
				"Size", fmt.Sprint(block.Size),
				"Weight", fmt.Sprint(block.Weight),
				"Bits", fmt.Sprint(block.Bits)),
			component.HorizontalLine(gtx, 1, th.Color.Dark),
			trioFields(gtx, th, 16, 16,
				"TxNum", fmt.Sprint(block.TxNum),
				"StrippedSize", fmt.Sprint(block.StrippedSize),
				"Version", fmt.Sprint(block.Version)),
			blockField(gtx, th, layout.Vertical, 16, 12, "Tx", fmt.Sprint(block.Tx)),
			blockField(gtx, th, layout.Vertical, 14, 12, "RawTx", fmt.Sprint(block.RawTx)),
			blockNavButtons(rc, gtx, th, block.PreviousHash, block.NextHash),
		}
		layautList.Layout(gtx, len(widgets), func(i int) {
			layout.UniformInset(unit.Dp(4)).Layout(gtx, widgets[i])
		})

	}
}

func trioFields(gtx *layout.Context, th *theme.DuoUItheme, labelTextSize, valueTextSize float32, unoLabel, unoValue, duoLabel, duoValue, treLabel, treValue string) func() {

	return func() {
		layout.Flex{
			Axis:    layout.Horizontal,
			Spacing: layout.SpaceBetween,
		}.Layout(gtx,
			layout.Flexed(0.3, blockField(gtx, th, layout.Vertical, labelTextSize, valueTextSize, unoLabel, fmt.Sprint(unoValue))),
			layout.Flexed(0.3, blockField(gtx, th, layout.Vertical, labelTextSize, valueTextSize, duoLabel, fmt.Sprint(duoValue))),
			layout.Flexed(0.3, blockField(gtx, th, layout.Vertical, labelTextSize, valueTextSize, treLabel, fmt.Sprint(treValue))),
		)

	}
}

func blockField(gtx *layout.Context, th *theme.DuoUItheme, axis layout.Axis, labelTextSize, valueTextSize float32, label, value string) func() {
	return func() {
		layout.Flex{
			Axis: axis,
		}.Layout(gtx,
			layout.Rigid(blockFieldValue(gtx, th, label, th.Color.Light, th.Color.Dark, th.Font.Primary, labelTextSize)),
			layout.Rigid(blockFieldValue(gtx, th, value, th.Color.Light, th.Color.DarkGray, th.Font.Mono, valueTextSize)))
	}
}

func blockNavButtons(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, previousBlockHash, nextBlockHash string) func() {
	return func() {
		layout.Flex{}.Layout(gtx,
			layout.Flexed(0.5, func() {
				var previousBlockButton theme.DuoUIbutton
				previousBlockButton = th.DuoUIbutton(th.Font.Mono, "Previous Block "+previousBlockHash, th.Color.Light, th.Color.Info, "", th.Color.Light, 16, 0, 60, 24, 0, 0)
				for previousBlockHashButton.Clicked(gtx) {
					//clipboard.Set(b.BlockHash)
					rc.ShowPage = "BLOCK" + previousBlockHash
					rc.GetSingleBlock(previousBlockHash)()
					component.SetPage(rc, blockPage(rc, gtx, th, previousBlockHash))
				}
				previousBlockButton.Layout(gtx, previousBlockHashButton)
			}),
			layout.Flexed(0.5, func() {
				var nextBlockButton theme.DuoUIbutton
				nextBlockButton = th.DuoUIbutton(th.Font.Mono, "Next Block "+nextBlockHash, th.Color.Light, th.Color.Info, "", th.Color.Light, 16, 0, 60, 24, 0, 0)
				for nextBlockHashButton.Clicked(gtx) {
					//clipboard.Set(b.BlockHash)
					rc.ShowPage = "BLOCK" + nextBlockHash
					rc.GetSingleBlock(nextBlockHash)()
					component.SetPage(rc, blockPage(rc, gtx, th, nextBlockHash))
				}
				nextBlockButton.Layout(gtx, nextBlockHashButton)
			}))
	}
}

func blockFieldValue(gtx *layout.Context, th *theme.DuoUItheme, text, color, bgColor string, font text.Typeface, textSize float32) func() {
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

func fill(gtx *layout.Context, col color.RGBA) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d}
}
