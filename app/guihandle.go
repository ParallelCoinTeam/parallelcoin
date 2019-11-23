package app

import (
	"github.com/p9c/pod/cmd/gui"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/shurcooL/vfsgen"
	"github.com/urfave/cli"
	"net/http"
)

var guiHandle = func(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) (err error) {

		//utils.GetBiosMessage(view, "starting GUI")

		//err := gui.Services(cx)
		Configure(cx, c)

		var fs http.FileSystem = http.Dir("./pkg/gui/assets")
		err = vfsgen.Generate(fs, vfsgen.Options{
			PackageName:  "guiLibs",
			BuildTags:    "dev",
			VariableName: "WalletGUI",
		})
		if err != nil {
			log.FATAL(err)
		}

		bios := &gui.Bios{
			Fs: &fs,
			IsFirstRun: &cx.FirstRun,
		}


		gui.Loader(bios, cx)
		err = gui.Services(cx)
		if err != nil {
			log.ERROR(err)
		}
		// We open up wallet creation
		gui.GUI(bios)

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
