package gwallet

//
//func (g *GioWallet) getPages(d *mod.Dap) mod.Pages {
//	return mod.Pages{
//		"Welcome": mod.Page{
//			Title:  "Welcome",
//			Header: g.welcomeHeader(d.UI.Theme),
//			Body:   g.welcomeBody(),
//		},
//		"Overview": mod.Page{
//			Title:  "Overview",
//			Header: g.overviewHeader(),
//			Body:   g.overviewBody(d),
//		},
//		"Send": mod.Page{
//			Title:  "Send",
//			Header: g.sendHeader(d.UI.Theme),
//			Body:   g.sendBody(d.UI),
//		},
//		"Receive": mod.Page{
//			Title:  "Receive",
//			Header: g.receiveHeader(d.UI.Theme),
//			Body:   g.receiveBody(d.UI.Theme),
//		},
//		"Transactions": mod.Page{
//			Title:  "Transactions",
//			Header: g.transactionsHeader(d.UI.Theme),
//			Body:   g.transactionsBody(d.UI.Theme),
//		},
//		"Explore": mod.Page{
//			Title:  "Explore",
//			Header: g.exploreHeader(),
//			Body:   g.exploreBody(),
//		},
//		"Peers": mod.Page{
//			Title:  "Peers",
//			Header: g.peersHeader(d.UI.Theme),
//			Body:   g.peersBody(),
//		},
//		"Settings": mod.Page{
//			Title:  "Settings",
//			Header: g.settingsHeader(d.UI.Theme),
//			Body:   g.settingsBody(d.UI.Theme),
//		},
//	}
//}

//func Page(th *theme.Theme,page mod.Page) func(gtx C) D {
//	return ContainerLayout(th.Colors["White"], th.Colors["White"], th.Colors["White"], 0, 0, 0, func(gtx C) D {
//		return lyt.Format(gtx, "vflex(start,r(inset(0dp0dp0dp0dp,_)),f(1,inset(0dp0dp0dp0dp,_)))", page.Header, page.Body)
//	})
//}
//
//
//func (g *GioWallet) pageButton(d *mod.Dap, b *widget.Clickable, f func(), icon, page string) func(gtx C) D {
//	return func(gtx C) D {
//		btn := material.IconButton(d.UI.Theme.T, b, d.UI.Theme.Icons[icon])
//		btn.Inset = layout.Inset{unit.Dp(2), unit.Dp(2), unit.Dp(2), unit.Dp(2)}
//		btn.Size = unit.Dp(21)
//		btn.Background = helper.HexARGB(d.UI.Theme.Colors["Secondary"])
//		for b.Clicked() {
//			f()
//			d.UI.N.CurrentPage = page
//		}
//		return btn.Layout(gtx)
//	}
//}
