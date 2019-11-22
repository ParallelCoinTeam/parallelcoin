package app

import (
	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/cmd/gui"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/urfave/cli"
)

type Bios struct {
	Theme      bool   `json:"theme"`
	IsBoot     bool   `json:"boot"`
	IsBootMenu bool   `json:"menu"`
	IsBootLogo bool   `json:"logo"`
	IsLoading  bool   `json:"loading"`
	IsDev      bool   `json:"dev"`
	IsScreen   string `json:"screen"`
}

var guiHandle = func(cx *conte.Xt) func(c *cli.Context) error{
	return func(c *cli.Context) (err error) {

		//utils.GetBiosMessage(view, "starting GUI")

		Configure(cx, c)
		//err := gui.Services(cx)


		if !apputil.FileExists(*cx.Config.WalletFile){
			// We can open wallet directly
			gui.FirstRun(cx)
		}

		err = gui.Services(cx)
			if err != nil{
				log.ERROR(err)
			}
			// We open up wallet creation
			gui.GUI(cx)

		//b.IsBootLogo = false
		//b.IsBoot = false

		if !cx.Node.Load().(bool) {
			close(cx.WalletKill)
		}
		if !cx.Wallet.Load().(bool) {
			close(cx.NodeKill)
		}
		return
	}
}
