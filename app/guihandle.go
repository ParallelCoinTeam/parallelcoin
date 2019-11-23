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

var guiHandle = func(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) (err error) {
		var firstRun bool
		if !apputil.FileExists(*cx.Config.WalletFile) {
			firstRun = true
		}
		//utils.GetBiosMessage(view, "starting GUI")

		//err := gui.Services(cx)
		Configure(cx, c)

		var fs http.FileSystem = http.Dir("./pkg/gui/assets/filesystem")
		err = vfsgen.Generate(fs, vfsgen.Options{
			PackageName:  "guiFileSystem",
			BuildTags:    "dev",
			VariableName: "WalletGUI",
		})
		if err != nil {
			log.FATAL(err)
		}

		bios := &gui.Bios{
			Fs: &fs,
			IsFirstRun: firstRun,
		}


		gui.Loader(bios, cx)
		err = gui.Services(cx)
		if err != nil {
			log.ERROR(err)
		}
		// We open up wallet creation
		gui.GUI(bios, cx)

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
