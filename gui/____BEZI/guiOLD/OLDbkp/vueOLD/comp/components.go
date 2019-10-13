package comp

import (
	"github.com/p9c/pod/cmd/gui/OLDbkp/vue/comp/comp"
	"github.com/p9c/pod/cmd/gui/OLDbkp/vue/comp/panel"
	"github.com/p9c/pod/cmd/gui/OLDbkp/vue/comp/sys"
	"github.com/p9c/pod/cmd/gui/OLDbkp/vue/db"
	"github.com/p9c/pod/cmd/gui/OLDbkp/vue/mod"
)

func Apps(d db.DuoVUEdb) (c []mod.DuoVUEcomp) {
	c = append(c, sys.Boot())
	//c = append(c, sys.Dev())
	c = append(c, sys.Display())

	//c = append(c, serv.Blocks())
	c = append(c, panel.Send())
	c = append(c, panel.AddressBook())
	//c = append(c, panel.Blocks())
	c = append(c, panel.Peers())
	c = append(c, panel.Settings())
	c = append(c, panel.Transactions())
	c = append(c, panel.TransactionsExcerpts())
	c = append(c, panel.TimeBalance())
	//c = append(c, sys.Nav())
	//c = append(c, panel.Manager())
	//c = append(c, panel.LayoutConfig())
	//c = append(c, panel.Test())
	//c = append(c, panel.ChartA())
	//c = append(c, panel.ChartB())
	c = append(c, panel.WalletStatus())
	c = append(c, panel.Status())
	c = append(c, panel.LocalHashRate())
	c = append(c, panel.NetworkHashRate())
	return c
}
func Components(d db.DuoVUEdb) (c []mod.DuoVUEcomp) {
	c = append(c, panel.Address())
	c = append(c, comp.Logo())
	c = append(c, comp.Header())
	c = append(c, comp.Sidebar())
	c = append(c, comp.Screen())
	c = append(c, comp.Menu())

	return c
}

var GetAppHtml string = `<!DOCTYPE html><html lang="en" ><head><meta charset="UTF-8"><title>ParallelCoin Wallet - True Story</title></head><body><header is="boot" id="boot"></header><display :is="display" id="display" class=" lightTheme"></display><footer is="dev" id="dev"></footer></body></html>`
