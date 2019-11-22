package app

import (
	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/cmd/gui"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/shurcooL/vfsgen"
	"github.com/urfave/cli"
	"net/http"
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

		var fs http.FileSystem = http.Dir("./pkg/gui/assets")
		err = vfsgen.Generate(fs, vfsgen.Options{
			PackageName:  "guiLibs",
			BuildTags:    "!dev",
			VariableName: "WalletGUI",
		})
		if err != nil {
			log.FATAL(err)
		}
		cx.FileSystem = &fs

		if !apputil.FileExists(*cx.Config.WalletFile){
			// We can open wallet directly
			gui.Loader(cx)
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
