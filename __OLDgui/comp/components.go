package comp

import (
	"github.com/p9c/pod/pkg/duos/db"
	"github.com/p9c/pod/pkg/duos/mod"
	"github.com/p9c/pod/pkg/duos/sys"
	"github.com/p9c/pod/pkg/svelte/__OLDvue/pnl"
)

func Apps(d db.DuOSdb) (c []mod.DuOScomp) {
	c = append(c, sys.Boot())
	//c = append(c, sys.Dev())
	c = append(c, sys.Display())

	//c = append(c, serv.Blocks())
	c = append(c, pnl.Send())
	c = append(c, pnl.AddressBook())
	//c = append(c, pnl.Blocks())
	c = append(c, pnl.Peers())
	c = append(c, pnl.Settings())
	c = append(c, pnl.Transactions())
	c = append(c, pnl.TransactionsExcerpts())
	c = append(c, pnl.TimeBalance())
	//c = append(c, sys.Nav())
	//c = append(c, pnl.Manager())
	//c = append(c, pnl.LayoutConfig())
	//c = append(c, pnl.Test())
	//c = append(c, pnl.ChartA())
	//c = append(c, pnl.ChartB())
	c = append(c, pnl.WalletStatus())
	c = append(c, pnl.Status())
	c = append(c, pnl.LocalHashRate())
	c = append(c, pnl.NetworkHashRate())
	return c
}

//func Components(d db.DuOSdb) (c []mod.DuOScomp) {
//	c = append(c, pnl.Address())
//	c = append(c, lib.Logo())
//	c = append(c, html.Header())
//	c = append(c, lib.Sidebar())
//	c = append(c, lib.Screen())
//	c = append(c, lib.Menu())
//
//	return c
//}

var GetAppHtml string = `<!DOCTYPE html><html lang="en" ><head><meta charset="UTF-8"><title>ParallelCoin Wallet - True Story</title></head><body><header is="boot" id="boot"></header><display :is="display" id="display" class=" lightTheme"></display><footer is="dev" id="dev"></footer></body></html>`
