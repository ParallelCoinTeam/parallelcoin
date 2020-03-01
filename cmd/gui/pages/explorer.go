package pages

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/component"
	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/cmd/gui/rcd"
)

var (
	blocksList = &layout.List{
		Axis: layout.Vertical,
	}
)

func Explorer(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) *theme.DuoUIpage {
	return th.DuoUIpage("EXPLORER", 0, rc.GetBlocksExcerpts(), component.ContentHeader(gtx, th, headerExplorer(rc, gtx, th)), bodyExplorer(rc, gtx, th), func() {})
}
func bodyExplorer(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		rc.GetBlocksExcerpts()
		in := layout.UniformInset(unit.Dp(0))
		in.Layout(gtx, func() {
			blocksList.Layout(gtx, len(rc.Explorer.Blocks), func(i int) {
				b := rc.Explorer.Blocks[i]
				blockRow(rc, gtx, th, &b)
			})
		})
	}
}

func headerExplorer(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		layout.Flex{
			Spacing: layout.SpaceBetween,
			Axis:    layout.Horizontal,
		}.Layout(gtx,
			//layout.Rigid(ui.txsFilter()),
			layout.Rigid(component.Label(gtx, th, th.Font.Primary, 12, th.Color.Light, "Block count: "+fmt.Sprint(rc.Status.Node.BlockCount))),
			layout.Rigid(component.Label(gtx, th, th.Font.Primary, 12, th.Color.Light, "Pages: "+fmt.Sprint(rc.Explorer.Page.To))),
			//layout.Rigid(component.Label(gtx, th, th.Font.Primary, 12, th.Color.Light, "Block count: "+fmt.Sprint(rc.Status.Node.BlockCount))),
			layout.Flexed(0.5, func() {
				th.DuoUIcounter(rc.GetBlocksExcerpts()).Layout(gtx, rc.Explorer.Page)
			}),
			layout.Flexed(0.5, func() {
				th.DuoUIcounter(rc.GetBlocksExcerpts()).Layout(gtx, rc.Explorer.PerPage)
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
				linkButton = th.DuoUIbutton(th.Font.Mono, fmt.Sprint(block.Height), th.Color.Light, th.Color.Info, th.Color.Info, th.Color.Light, "", th.Color.Light, 14, 0, 60, 24, 0, 0)
				for block.Link.Clicked(gtx) {
					rc.ShowPage = fmt.Sprintf("BLOCK %s", block.BlockHash)
					rc.GetSingleBlock(block.BlockHash)()
					component.SetPage(rc, blockPage(rc, gtx, th, block.BlockHash))
				}
				linkButton.Layout(gtx, block.Link)
			}),
			layout.Rigid(func() {
				l := th.Body2(block.Time)
				l.Font.Typeface = th.Font.Mono
				l.Color = theme.HexARGB(th.Color.Dark)
				l.Layout(gtx)
			}),
			layout.Rigid(func() {
				l := th.Body2(block.BlockHash)
				l.Font.Typeface = th.Font.Mono
				l.Color = theme.HexARGB(th.Color.Dark)
				l.Layout(gtx)
			}))
	})
}
