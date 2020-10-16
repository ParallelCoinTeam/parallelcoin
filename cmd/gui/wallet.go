package gwallet

import (
	"gioui.org/text"
	"github.com/p9c/pod/pkg/gui/wallet/lyt"
	"github.com/p9c/pod/pkg/gui/wallet/theme"
)

//func Main(d *mod.Dap) W {
//	return func(gtx C) D {
//		return lyt.Format(gtx, d.UI.R.Mod["Main"].(string),
//			func(gtx C) D {
//				return d.UI.N.Nav(d.UI.Theme, gtx)
//			},
//			func(gtx C) D {
//				return lyt.Format(gtx, d.UI.R.Mod["Content"].(string),
//					noReturn,
//				d.Page(d.UI.Theme,d.UI.P[d.UI.N.CurrentPage]),
//					d.footer(d),
//				)
//			})
//	}
//}
//
//func (g *GioWallet) BeforeMain(gtx C) {
//	//g.Resposnsivity(gtx)
//	//g.UI.N.mode = g.UI.Res.Mode
//	//g.UI.N.navLayout = g.UI.Res.Mod["Nav"].(string)
//	//g.UI.N.itemLayout = g.UI.Res.Mod["NavIconAndLabel"].(string)
//	//g.UI.N.axis = g.UI.Res.Mod["NavItemsAxis"].(layout.Axis)
//	//g.UI.N.size = g.UI.Res.Mod["NavSize"].(int)
//	//g.UI.N.noContent = g.UI.N.wide
//	//g.UI.N.logo = g.logo()
//}
//
//func (g *GioWallet) AfterMain() {
//
//}
//
//func (g *GioWallet) Tik() func() {
//	return func() {
//		g.GetBalances()
//		g.GetLatestTransactions()
//	}
//}

func labeledRow(th *theme.Theme, label string, content func(gtx C) D) func(gtx C) D {
	return func(gtx C) D {
		return lyt.Format(gtx, "hflex(start,r(inset(8dp0dp0dp8dp,_)),f(1,_))",
			func(gtx C) D {
				gtx.Constraints.Min.X = 80
				gtx.Constraints.Min.X = 80
				title := theme.Body(th, label)
				title.Alignment = text.Start
				return title.Layout(gtx)
			},
			content,
		)
	}
}
