package gui

import (
	"gioui.org/layout"

	"github.com/p9c/pod/pkg/gui/wallet/dap/box"
	"github.com/p9c/pod/pkg/gui/wallet/dap/res"
	"github.com/p9c/pod/pkg/gui/wallet/lyt"
)

func (g *GuiAppModel) BeforeMain(gtx C) {
	g.ui.R = res.Resposnsivity(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
	g.ui.N.Mode = g.ui.R.Mode
	g.ui.N.NavLayout = g.ui.R.Mod["Nav"].(string)
	g.ui.N.ItemLayout = g.ui.R.Mod["NavIconAndLabel"].(string)
	g.ui.N.Axis = g.ui.R.Mod["NavItemsAxis"].(layout.Axis)
	g.ui.N.Size = g.ui.R.Mod["NavSize"].(int)
	g.ui.N.NoContent = g.ui.N.Wide
	g.ui.N.LogoWidget = g.ui.N.LogoLayout(g.ui.Theme)
}

func (g *GuiAppModel) AfterMain(gtx C) {
	// pop.Popup(gtx, g.ui.Theme, func(gtx C)D{return material.H3(g.ui.Theme.T,"tetstette").Layout(gtx)})
	// return lyt.Format(gtx, "hflex(middle,f(1,inset(8dp8dp8dp8dp,_)))",
	// pop.Popup(g.ui.Theme, func(gtx C) D {
	//	title := theme.Body(g.ui.Theme, "Requested payments history")
	//	title.Alignment = text.Start
	//	return title.Layout(gtx)
	//	}),
	// })

}

func (g *GuiAppModel) Main(gtx C) D {
	return lyt.Format(gtx, "max(inset(0dp0dp0dp0dp,_))", func(gtx C) D {
		return lyt.Format(gtx, g.ui.R.Mod["Container"].(string),
			box.BoxBase(g.ui.Theme.Colors["NavBg"], g.ui.N.Nav(g.ui.Theme, gtx)),
			func(gtx C) D {
				return lyt.Format(gtx, g.ui.R.Mod["Main"].(string),
					g.ui.N.CurrentPage.P(g.ui.Theme, g.ui.R.Mod["Page"].(string)),
					g.ui.F,
				)
			})
	})
}
