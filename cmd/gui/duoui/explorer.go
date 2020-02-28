package duoui

import (
	"fmt"
	"gioui.org/layout"
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
						rc.GetSingleBlock(block.BlockHash)
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
				layout.Flex{
					Axis:    layout.Horizontal,
					Spacing: layout.SpaceAround,
				}.Layout(gtx,
					//layout.Flexed(0.6, func() {
					//	var linkButton theme.DuoUIbutton
					//	linkButton = th.DuoUIbutton(th.Font.Mono, fmt.Sprint(block.Height), th.Color.Light, th.Color.Dark, "", th.Color.Light, 16, 0, 60, 24, 0, 0)
					//	for block.Link.Clicked(gtx) {
					//		//clipboard.Set(b.BlockHash)
					//		rc.ShowPage = "BLOCK" + block.BlockHash
					//		setPage(rc, blockPage(rc, gtx, th, block.BlockHash))
					//	}
					//	linkButton.Layout(gtx, block.Link)
					//}),
					//layout.Rigid(func() {
					//	amount := th.H5(fmt.Sprintf("%0.8f", block.Amount))
					//	amount.Font.Typeface = th.Font.Primary
					//	amount.Color = theme.HexARGB(th.Color.Dark)
					//	amount.Alignment = text.End
					//	amount.Font.Variant = "Mono"
					//	amount.Font.Weight = text.Bold
					//	amount.Layout(gtx)
					//}),
					layout.Rigid(func() {
						sat := th.Body1(fmt.Sprint(block.TxNum))
						sat.Font.Typeface = th.Font.Primary
						sat.Color = theme.HexARGB(th.Color.Dark)
						sat.Layout(gtx)
					}),
					layout.Rigid(func() {
						sat := th.Body1(fmt.Sprint(block.Hash))
						sat.Font.Typeface = th.Font.Mono
						sat.Color = theme.HexARGB(th.Color.Dark)
						sat.Layout(gtx)
					}),
					layout.Rigid(func() {
						l := th.Body2(fmt.Sprint(block.Time))
						l.Font.Typeface = th.Font.Primary
						l.Color = theme.HexARGB(th.Color.Dark)
						l.Layout(gtx)
					}),
				)
			}),
			layout.Rigid(func() {
				sat := th.Body1(fmt.Sprintf("%0.8f", block.NextHash))
				sat.Color = theme.HexARGB(th.Color.Dark)
				sat.Layout(gtx)
			}))
	}
}
