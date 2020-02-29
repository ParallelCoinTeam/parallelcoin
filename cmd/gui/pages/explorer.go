package pages

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/component"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/cmd/gui/rcd"
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

func Explorer(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) *theme.DuoUIpage {
	return th.DuoUIpage("EXPLORER", 8, rc.GetBlocksExcerpts(page.Value, perPage.Value), component.ContentHeader(gtx, th, headerExplorer(gtx, th)), bodyExplorer(rc, gtx, th), func() {})
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
					rc.ShowPage = fmt.Sprintf("BLOCK %s", block.BlockHash)
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
