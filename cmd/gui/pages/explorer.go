package pages

import (
	"fmt"
	"time"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gelook"
)

var (
	blocksList = &layout.List{
		Axis: layout.Vertical,
	}
	txwidth int
)

func DuoUIexplorer(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "EXPLORER",
		TxColor:       "",
		Command:       rc.GetBlocksExcerpts(),
		Border:        4,
		BorderColor:   th.Colors["Light"],
		Header:        explorerHeader(rc, gtx, th),
		HeaderBgColor: "",
		HeaderPadding: 4,
		Body:          bodyExplorer(rc, gtx, th),
		BodyBgColor:   th.Colors["Light"],
		BodyPadding:   4,
		Footer:        func() {},
		FooterBgColor: "",
		FooterPadding: 0,
	}
	return th.DuoUIpage(page)
}
func bodyExplorer(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		rc.GetBlocksExcerpts()
		layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func() {
				blockRowCellLabels(rc, gtx, th)
			}),
			layout.Flexed(1, func() {
				blocksList.Layout(gtx, len(rc.Explorer.Blocks), func(i int) {
					b := rc.Explorer.Blocks[i]
					blockRow(rc, gtx, th, &b)
				})
			}),
		)
	}
}

func explorerHeader(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
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

func blockRowCellLabels(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) {
	layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
		// component.HorizontalLine(gtx, 1, th.Colors["Dark"])()
		layout.Flex{
			Spacing: layout.SpaceBetween,
		}.Layout(gtx,
			layout.Rigid(func() {
				layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
					l := th.Body2("Height")
					l.Font.Typeface = th.Fonts["Mono"]
					l.Alignment = text.Middle
					l.Color = th.Colors["Dark"]
					l.Layout(gtx)
				})
			}),
			layout.Rigid(func() {
				layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
					l := th.Body2("Time")
					l.Font.Typeface = th.Fonts["Mono"]
					l.Alignment = text.Middle
					l.Color = th.Colors["Dark"]
					l.Layout(gtx)
				})
			}),
			layout.Rigid(func() {
				layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
					l := th.Body2("Confirmations")
					l.Font.Typeface = th.Fonts["Mono"]
					l.Color = th.Colors["Dark"]
					l.Layout(gtx)
				})
			}),
			layout.Rigid(func() {
				layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
					l := th.Body2("TxNum")
					l.Font.Typeface = th.Fonts["Mono"]
					l.Color = th.Colors["Dark"]
					l.Layout(gtx)
				})
			}),
			layout.Rigid(func() {
				layout.Inset{
					Right: unit.Dp(float32(txwidth - 64)),
				}.Layout(gtx, func() {
					l := th.Body2("BlockHash")
					l.Font.Typeface = th.Fonts["Mono"]
					l.Color = th.Colors["Dark"]
					l.Layout(gtx)
				})
			}))
	})
}

func blockRow(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, block *model.DuoUIblock) {
	layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
		gelook.DuoUIdrawRectangle(gtx, gtx.Constraints.Width.Max, gtx.Constraints.Height.Max, th.Colors["Light"], [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		th.DuoUIline(gtx, 0, 0, 1, th.Colors["Gray"])()
		layout.Flex{
			Spacing: layout.SpaceBetween,
		}.Layout(gtx,
			layout.Rigid(func() {
				var linkButton gelook.DuoUIbutton
				linkButton = th.DuoUIbutton(th.Fonts["Mono"], fmt.Sprint(block.Height), th.Colors["Light"], th.Colors["Info"], th.Colors["Info"], th.Colors["Dark"], "", th.Colors["Dark"], 14, 0, 60, 24, 0, 0)
				for block.Link.Clicked(gtx) {
					rc.ShowPage = fmt.Sprintf("BLOCK %s", block.BlockHash)
					rc.GetSingleBlock(block.BlockHash)()
					component.SetPage(rc, blockPage(rc, gtx, th, block.BlockHash))
				}
				linkButton.Layout(gtx, block.Link)
			}),
			layout.Rigid(func() {
				layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
					l := th.Body2(fmt.Sprint(time.Unix(block.Time, 0).Format("2006-01-02 15:04:05")))
					l.Font.Typeface = th.Fonts["Mono"]
					l.Alignment = text.Middle
					l.Color = th.Colors["Dark"]
					l.Layout(gtx)
				})
			}),
			layout.Rigid(func() {
				layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
					l := th.Body2(fmt.Sprint(block.Confirmations))
					l.Font.Typeface = th.Fonts["Mono"]
					l.Color = th.Colors["Dark"]
					l.Layout(gtx)
				})
			}),
			layout.Rigid(func() {
				layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
					l := th.Body2(fmt.Sprint(block.TxNum))
					l.Font.Typeface = th.Fonts["Mono"]
					l.Color = th.Colors["Dark"]
					l.Layout(gtx)
				})
			}),
			layout.Rigid(func() {
				layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
					l := th.Body2(block.BlockHash)
					l.Font.Typeface = th.Fonts["Mono"]
					l.Color = th.Colors["Dark"]
					l.Layout(gtx)
					txwidth = gtx.Dimensions.Size.X
				})
			}))
	})
}
