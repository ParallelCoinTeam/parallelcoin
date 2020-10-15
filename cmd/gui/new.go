package gwallet

import (
	"fmt"
	"gioui.org/widget"
	"github.com/p9c/pod/pkg/gui/wallet/dap/mod"
	"github.com/p9c/pod/pkg/gui/wallet/dap/page"
	"github.com/p9c/pod/pkg/gui/wallet/nav"
	rpcclient "github.com/p9c/pod/pkg/rpc/client"
)

func NewGioWallet(d *mod.Dap) mod.Sap {
	g := &GioWallet{
		ui: d.UI,
	}

	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:18444",
		User:         "bitnode",
		Pass:         "44444",
		HTTPPostMode: true,  // Bitcoin core only supports HTTP POST mode
		TLS:          false, // Bitcoin core does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {

	}

	g.rpc = client

	defer g.rpc.Shutdown()

	//g.UI.pages = g.getPages()
	//g.menuItems = g.getMenuItems()

	g.Status.bal = &Balances{}
	g.GetOverview()
	g.ui.N.MenuItems = append(g.getMenuItems(g.ui))
	g.ui.N.CurrentPage = page.Page{
		Title:  "Overview",
		Header: g.overviewHeader(),
		Body:   g.overviewBody(),
		Footer: noReturn,
	}
	d.UI.F = g.footer()

	CreateSendAddressItem()()

	balances = dummyBalances
	latestTransactions = dummyLatestTxs

	transactions = dummyTxs

	return mod.Sap{
		Title: "ParallelCoin",
		App:   g,
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
}

func (g *GioWallet) getMenuItem(hide bool, title string, header, body, footer func(gtx C) D) nav.Item {
	return nav.Item{
		Title: title,
		Icon:  g.ui.Theme.Icons[title],
		Btn:   new(widget.Clickable),
		Page: page.Page{
			Title:  title,
			Header: header,
			Body:   body,
			Footer: footer,
		},
		HideOnMob: false,
	}
}

func (g *GioWallet) getMenuItems(ui *mod.UserInterface) []nav.Item {
	return []nav.Item{
		g.getMenuItem(false, "Overview", g.overviewHeader(), g.overviewBody(), noReturn),
		g.getMenuItem(false, "Send", g.sendHeader(), g.sendBody(), g.sendFooter()),
		g.getMenuItem(false, "Receive", g.receiveHeader(), g.receiveBody(), noReturn),
		g.getMenuItem(false, "Transactions", g.transactionsHeader(), g.transactionsBody(), noReturn),
		//g.getMenuItem(true, "Explore", g.exploreHeader(), g.exploreBody(), noReturn),
		//g.getMenuItem(true, "Peers", g.peersHeader(), g.peersBody(), noReturn),
		//g.getMenuItem(true, "Settings", g.settingsHeader(), g.settingsBody(), noReturn),
	}
}
