package gui

import (
	"fmt"
	"time"

	l "gioui.org/layout"
	"gioui.org/text"
	"github.com/gioapp/gel/helper"

	"github.com/p9c/pod/pkg/gui/wallet/dap/box"
	"github.com/p9c/pod/pkg/gui/wallet/lyt"
	"github.com/p9c/pod/pkg/gui/wallet/theme"
)

func (g *GuiAppModel) GetTransactions() {

}

func (g *GuiAppModel) transactionsHeader() func(gtx C) D {
	return box.BoxBase(g.ui.Theme.Colors["PanelBg"], func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		helper.Fill(gtx, helper.HexARGB(g.ui.Theme.Colors["PanelBg"]))
		return lyt.Format(gtx, "vflex(middle,r(inset(5dp0dp5dp0dp,_))))",
			func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				title := theme.H6(g.ui.Theme, "Transactions Header")
				title.Alignment = text.Start
				return title.Layout(gtx)
			})
	})
}

func (g *GuiAppModel) sssssStransactionsBody() func(gtx C) D {
	return box.BoxPanel(g.ui.Theme, noReturn)
	// return box.BoxPanel(g.ui.Theme, func(gtx C) D {
	// return transactionsList.Layout(gtx, len(transactions), func(gtx C, i int) D {
	//	tx := transactions[i]
	//	return material.Clickable(gtx, tx.Btn, func(gtx C) D {
	//		return lyt.Format(gtx, "vflex(r(_),r(_))",
	//			func(gtx C) D {
	//				return lyt.Format(gtx, "hflex(r(_),r(inset(4dp4dp4dp4dp,_)),r(inset(4dp4dp4dp4dp,_)),f(1,inset(4dp4dp4dp4dp,_)),r(inset(4dp4dp4dp4dp,_)))",
	//					func(gtx C) D {
	//						var d D
	//						size := gtx.Px(unit.Dp(24))
	//						l := g.ui.Theme.Icons["Explore"]
	//						l.Color = helper.HexARGB(g.ui.Theme.Colors["Primary"])
	//						l.Layout(gtx, unit.Px(float32(size)))
	//						d = D{
	//							Size: image.Point{X: size, Y: size},
	//						}
	//						return d
	//					},
	//					func(gtx C) D {
	//						title := theme.Body(g.ui.Theme, tx.Time)
	//						title.Alignment = text.Start
	//						return title.Layout(gtx)
	//					},
	//					func(gtx C) D {
	//						title := theme.Body(g.ui.Theme, tx.Type)
	//						title.Alignment = text.End
	//						return title.Layout(gtx)
	//					},
	//					func(gtx C) D {
	//						title := theme.Body(g.ui.Theme, tx.Address)
	//						title.Alignment = text.Start
	//						return title.Layout(gtx)
	//					},
	//					func(gtx C) D {
	//						title := theme.Body(g.ui.Theme, tx.Amount)
	//						title.Alignment = text.Start
	//						return title.Layout(gtx)
	//					},
	//				)
	//			},
	//			helper.DuoUIline(false, 1, 0, 1, g.ui.Theme.Colors["Border"]))
	//	})
	// })
	// })
}

func (g *GuiAppModel) transactionsBody() func(gtx C) D {
	return func(gtx C) D {
		return g.Inset(0.25,
			g.Flex().Flexed(1, func(gtx l.Context) l.Dimensions {
				return g.lists["latestTransactions"].End().ScrollToEnd().Length(g.worker.solutionCount).ListElement(
					func(gtx l.Context, i int) l.Dimensions {
						return g.Flex().Rigid(
							g.Button(g.solButtons[i].SetClick(func() {
								currentBlock = g.worker.solutions[i]
								Debug("clicked for block", currentBlock.height)
								g.modalWidget = g.BlockDetails
								g.modalOn = true
							})).Text(fmt.Sprint(g.worker.solutions[i].height)).Inset(0.5).Fn,
						).Flexed(1,
							g.Inset(0.25,
								g.Flex().Vertical().Rigid(
									g.Flex().Rigid(
										g.Body1(g.worker.solutions[i].algo).Font("plan9").Fn,
									).Flexed(1,
										g.Flex().Vertical().Rigid(
											g.Body1(g.worker.solutions[i].hash).
												Font("go regular").
												TextScale(0.75).
												Alignment(text.End).
												Fn,
										).Rigid(
											g.Caption(fmt.Sprint(
												g.worker.solutions[i].time.Format(time.RFC3339))).
												Alignment(text.End).
												Fn,
										).Fn,
									).Fn,
								).Fn,
							).Fn,
						).Fn(gtx)
					}).Fn(gtx)
			}).Fn,
		).Fn(gtx)
	}
}
