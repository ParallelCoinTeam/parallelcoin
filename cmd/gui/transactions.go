package gwallet

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/gioapp/gel/helper"
	"github.com/p9c/pod/pkg/gui/wallet/dap/box"
	"github.com/p9c/pod/pkg/gui/wallet/lyt"
	"github.com/p9c/pod/pkg/gui/wallet/theme"
	"image"
)

var (
	transactions     []Tx
	transactionsList = &layout.List{
		Axis: layout.Vertical,
	}
)

func (g *GioWallet) GetTransactions() {

}

func (g *GioWallet) transactionsHeader() func(gtx C) D {
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

func (g *GioWallet) transactionsBody() func(gtx C) D {
	return box.BoxPanel(g.ui.Theme, func(gtx C) D {
		return transactionsList.Layout(gtx, len(transactions), func(gtx C, i int) D {
			tx := transactions[i]
			return material.Clickable(gtx, tx.Btn, func(gtx C) D {
				return lyt.Format(gtx, "vflex(r(_),r(_))",
					func(gtx C) D {
						return lyt.Format(gtx, "hflex(r(_),r(inset(4dp4dp4dp4dp,_)),r(inset(4dp4dp4dp4dp,_)),f(1,inset(4dp4dp4dp4dp,_)),r(inset(4dp4dp4dp4dp,_)))",
							func(gtx C) D {
								var d D
								size := gtx.Px(unit.Dp(24))
								l := g.ui.Theme.Icons["Explore"]
								l.Color = helper.HexARGB(g.ui.Theme.Colors["Primary"])
								l.Layout(gtx, unit.Px(float32(size)))
								d = D{
									Size: image.Point{X: size, Y: size},
								}
								return d
							},
							func(gtx C) D {
								title := theme.Body(g.ui.Theme, tx.Time)
								title.Alignment = text.Start
								return title.Layout(gtx)
							},
							func(gtx C) D {
								title := theme.Body(g.ui.Theme, tx.Type)
								title.Alignment = text.End
								return title.Layout(gtx)
							},
							func(gtx C) D {
								title := theme.Body(g.ui.Theme, tx.Address)
								title.Alignment = text.Start
								return title.Layout(gtx)
							},
							func(gtx C) D {
								title := theme.Body(g.ui.Theme, tx.Amount)
								title.Alignment = text.Start
								return title.Layout(gtx)
							},
						)
					},
					helper.DuoUIline(false, 1, 0, 1, g.ui.Theme.Colors["Border"]))
			})
		})
	})
}
