package comp

import (
	"git.parallelcoin.io/dev/pod/cmd/gui/vue/comp/panel"
	"git.parallelcoin.io/dev/pod/cmd/gui/vue/comp/part"
	"git.parallelcoin.io/dev/pod/cmd/gui/vue/comp/serv"
	"git.parallelcoin.io/dev/pod/cmd/gui/vue/comp/sys"
	"git.parallelcoin.io/dev/pod/cmd/gui/vue/db"
	"git.parallelcoin.io/dev/pod/cmd/gui/vue/mod"
)

func Components(d db.DuoVUEdb) (c []mod.DuoVUEcomp) {
	c = append(c, sys.Boot())
	//c = append(c, serv())
	//c = append(c, sys.Dev())
	c = append(c, sys.Display())
	c = append(c, serv.SrvNode())
	//c = append(c, serv.SrvWallet())
	c = append(c, sys.Screen())
	//c = append(c, core.alerts())
	//c = append(c, serv.Blocks())
	c = append(c, part.PartAddress())
	c = append(c, panel.Send())
	//c = append(c, panel.Receive())
	c = append(c, panel.AddressBook())
	//c = append(c, panel.Blocks())
	c = append(c, panel.Peers())
	//c = append(c, panel.Settings())
	c = append(c, panel.Transactions())
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
