package pages

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/controller"
	"github.com/p9c/pod/pkg/gui/theme"
)

var (
	blocksList = &layout.List{
		Axis: layout.Vertical,
	}
	blocksPanel = &controller.Panel{
		Name: "",
		PanelContentLayout: &layout.List{
			Axis:        layout.Vertical,
			ScrollToEnd: false,
		},
	}
)

func DuoUIexplorer(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) *theme.DuoUIpage {
	return th.DuoUIpage("EXPLORER", 0, rc.GetBlocksExcerpts(), component.ContentHeader(gtx, th, headerExplorer(rc, gtx, th)), bodyExplorer(rc, gtx, th), func() {})
}
func bodyExplorer(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		rc.GetBlocksExcerpts()
		th.DuoUIpanel(explorerContent(rc, gtx, th)).Layout(gtx, addressBookPanel)
	}
}

func explorerContent(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
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
			Spacing:   layout.SpaceBetween,
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Rigid(func() {
				th.DuoUIcounter(rc.GetBlocksExcerpts()).Layout(gtx, rc.Explorer.Page, "PAGE", fmt.Sprint(rc.Explorer.Page.Value))
			}),
			layout.Rigid(func() {
				th.DuoUIcounter(rc.GetBlocksExcerpts()).Layout(gtx, rc.Explorer.PerPage, "PER PAGE", fmt.Sprint(rc.Explorer.PerPage.Value))
			}),
		)
	}
}

func blockRow(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, block *model.DuoUIblock) {
	layout.UniformInset(unit.Dp(4)).Layout(gtx, func() {
		component.HorizontalLine(gtx, 1, th.Colors["Dark"])()
		layout.Flex{
			Spacing: layout.SpaceBetween,
		}.Layout(gtx,
			layout.Rigid(func() {
				var linkButton theme.DuoUIbutton
				linkButton = th.DuoUIbutton(th.Fonts["Mono"], fmt.Sprint(block.Height), th.Colors["Light"], th.Colors["Info"], th.Colors["Info"], th.Colors["Light"], "", th.Colors["Light"], 14, 0, 60, 24, 0, 0)
				for block.Link.Clicked(gtx) {
					rc.ShowPage = fmt.Sprintf("BLOCK %s", block.BlockHash)
					rc.GetSingleBlock(block.BlockHash)()
					component.SetPage(rc, blockPage(rc, gtx, th, block.BlockHash))
				}
				linkButton.Layout(gtx, block.Link)
			}),
			layout.Rigid(func() {
				l := th.Body2(block.Time)
				l.Font.Typeface = th.Fonts["Mono"]
				l.Color = theme.HexARGB(th.Colors["Dark"])
				l.Layout(gtx)
			}),
			layout.Rigid(func() {
				l := th.Body2(block.BlockHash)
				l.Font.Typeface = th.Fonts["Mono"]
				l.Color = theme.HexARGB(th.Colors["Dark"])
				l.Layout(gtx)
			}))
	})
}
