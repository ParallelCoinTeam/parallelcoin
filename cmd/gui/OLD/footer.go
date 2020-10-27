package gui

import (
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/gioapp/gel/helper"
	"github.com/p9c/pod/pkg/gui/wallet/btn"
	"github.com/p9c/pod/pkg/gui/wallet/dap/box"
	"github.com/p9c/pod/pkg/gui/wallet/lyt"
	"github.com/p9c/pod/pkg/gui/wallet/theme"
)

var (
	syncBtn           = new(widget.Clickable)
	unitBtn           = new(widget.Clickable)
	connectionsBtn    = new(widget.Clickable)
	footerSearchInput = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
)

func (g *GuiAppModel) footer() func(gtx C) D {
	//return ContainerLayout(d.UI.Theme.Colors["Info"], d.UI.Theme.Colors["Info"], d.UI.Theme.Colors["Info"], 0, 0, 0, func(gtx C) D {
	return box.BoxBase(g.ui.Theme.Colors["FooterBg"], func(gtx C) D {
		return lyt.Format(gtx, "hflex(middle,f(1,inset(16dp0dp16dp6dp,_)),r(inset(0dp0dp0dp0dp,_)))",
			g.footerStatus(),
			g.footerMenu(),
		)
	})
}

func (g *GuiAppModel) footerMenu() func(gtx C) D {
	return func(gtx C) D {
		return lyt.Format(gtx, "hflex(middle,r(_),r(_),r(_))",
			footerButton(g.ui.Theme, sendBtn, "Send", func() {}),
			footerButton(g.ui.Theme, connectionsBtn, "Close", func() {}),
			footerButton(g.ui.Theme, sendBtn, "CounterPlus", func() {}),
		)
	}

}

func (g *GuiAppModel) footerStatus() func(gtx C) D {
	return func(gtx C) D {
		return lyt.Format(gtx, "hflex(middle,r(inset(0dp0dp0dp0dp,_)),f(1,_))",
			//ContainerLayout(d.UI.Theme.Colors["White"], d.UI.Theme.Colors["Dark"], d.UI.Theme.Colors["White"], 1, 1, 1, func(gtx C) D {
			func(gtx C) D {
				title := theme.Body(g.ui.Theme, "Connectiong to peers...")
				title.Alignment = text.Start
				return title.Layout(gtx)
			},
			box.BoxBase(g.ui.Theme.Colors["Primary"], func(gtx C) D {
				title := theme.Body(g.ui.Theme, "100 years and 12 weeks behind.")
				title.Alignment = text.Start
				title.Color = helper.HexARGB(g.ui.Theme.Colors["Light"])
				return title.Layout(gtx)
			}),
		)
	}
}

func footerButton(th *theme.Theme, c *widget.Clickable, icon string, onClick func()) func(gtx C) D {
	return func(gtx C) D {
		b := btn.IconTextBtn(th, c, "max(inset(0dp0dp0dp0dp,_))", true, "")
		b.Icon = th.Icons[icon]
		b.IconSize = unit.Dp(16)
		b.CornerRadius = unit.Dp(0)
		b.Background = th.Colors["NavBg"]
		b.TextColor = th.Colors["NavItem"]
		for c.Clicked() {
			onClick()
		}
		return b.Layout(gtx)
	}
}
