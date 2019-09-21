package comp

import (
	"github.com/parallelcointeam/parallelcoin/cmd/gui/vue/comp/panel"
	"github.com/parallelcointeam/parallelcoin/cmd/gui/vue/comp/serv"
	"github.com/parallelcointeam/parallelcoin/cmd/gui/vue/comp/sys"
	"github.com/parallelcointeam/parallelcoin/cmd/gui/vue/db"
	"github.com/parallelcointeam/parallelcoin/cmd/gui/vue/mod"
)

func Apps(d db.DuoVUEdb) (c []mod.DuoVUEcomp) {
	//c = append(c, sys.Boot())
	//c = append(c, sys.Dev())
	c = append(c, sys.Display())
	c = append(c, serv.Alert())
	//c = append(c, serv.Blocks())
	c = append(c, panel.Send())
	c = append(c, panel.AddressBook())
	//c = append(c, panel.Blocks())
	//c = append(c, panel.Peers())
	//c = append(c, panel.Settings())
	//c = append(c, panel.Transactions())
	//c = append(c, sys.Nav())
	//c = append(c, panel.Manager())
	//c = append(c, panel.LayoutConfig())
	//c = append(c, panel.Test())
	//c = append(c, panel.ChartA())
	//c = append(c, panel.ChartB())
	c = append(c, panel.WalletStatus())
	c = append(c, panel.Status())
	//c = append(c, panel.NetworkHashRate())
	return c
}
func Components(d db.DuoVUEdb) (c []mod.DuoVUEcomp) {
	c = append(c, panel.Address())
	return c
}
