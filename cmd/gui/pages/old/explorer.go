package pages

import (
	"fmt"
	"github.com/stalker-loki/pod/cmd/gui/component"
	"github.com/stalker-loki/pod/pkg/gui/gel"
	"time"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/stalker-loki/pod/cmd/gui/model"
	"github.com/stalker-loki/pod/cmd/gui/rcd"
	"github.com/stalker-loki/pod/pkg/gui/gelook"
)

var (
	explorerPanelElement = gel.NewPanel()
	txwidth              int
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

				explorerPanel := th.DuoUIpanel()
				explorerPanel.PanelObject = rc.Explorer.Blocks
				explorerPanel.ScrollBar = th.ScrollBar(16)
				explorerPanelElement.PanelObjectsNumber =
					len(rc.Explorer.Blocks)
				explorerPanel.Layout(gtx, explorerPanelElement,
					func(i int, in interface{}) {
						blocks := in.([]model.DuoUIblock)
						b := blocks[i]

						//blocksList.Layout(gtx, len(rc.Explorer.Blocks), func(i int) {
						//	b := rc.Explorer.Blocks[i]
						blockRow(rc, gtx, th, &b)
						//})
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
				th.DuoUIcounter(rc.GetBlocksExcerpts()).Layout(gtx,
					rc.Explorer.Page, "PAGE",
					fmt.Sprint(rc.Explorer.Page.Value))
			}),
			layout.Rigid(func() {
				th.DuoUIcounter(rc.GetBlocksExcerpts()).Layout(gtx,
					rc.Explorer.PerPage, "PER PAGE",
					fmt.Sprint(rc.Explorer.PerPage.Value))
			}),
		)
	}
}

func blockRowCellLabels(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) {
	cs := gtx.Constraints
	labels := th.DuoUIcontainer(0, th.Colors["Gray"])
	labels.FullWidth = true
	labels.Layout(gtx, layout.W, func() {
		// component.HorizontalLine(gtx, 1, th.Colors["Dark"])()
		gtx.Constraints.Width.Min = cs.Width.Max
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
					Top:   unit.Dp(8),
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
	for block.Link.Clicked(gtx) {
		rc.ShowPage = fmt.Sprintf("BLOCK %s", block.BlockHash)
		rc.GetSingleBlock(block.BlockHash)()
		component.SetPage(rc, blockPage(rc, gtx, th, block.BlockHash))
	}
	width := gtx.Constraints.Width.Max
	button := th.DuoUIbutton("", "",
		"", "",
		"", "",
		"", "",
		0, 0, 0, 0,
		0, 0, 0, 0)
	button.InsideLayout(gtx, block.Link, func() {
		layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func() {
				gtx.Constraints.Width.Min = width
				layout.Flex{
					Spacing: layout.SpaceBetween,
				}.Layout(gtx,
					layout.Rigid(func() {
						var linkButton gelook.DuoUIbutton
						linkButton = th.DuoUIbutton(th.Fonts["Mono"],
							fmt.Sprint(block.Height), th.Colors["Light"],
							th.Colors["Info"], th.Colors["Info"],
							th.Colors["Dark"], "", th.Colors["Dark"],
							14, 0, 60, 24,
							5, 8, 6, 8)
						linkButton.Layout(gtx, block.Link)
					}),
					layout.Rigid(func() {
						layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
							l := th.Body2(
								fmt.Sprint(time.Unix(block.Time, 0).
									Format("2006-01-02 15:04:05")))
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
						gelook.DuoUIcontainer{}.Layout(gtx, layout.Center, func() {
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
			}),
			layout.Rigid(func() {
				th.DuoUIline(gtx, 0, 0, 1,
					th.Colors["Gray"])()
			}))
	})
}
