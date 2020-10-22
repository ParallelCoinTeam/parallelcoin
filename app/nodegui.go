package app

import (
	"os"

	"github.com/urfave/cli"

	"github.com/p9c/pod/app/config"
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/node/gui"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func nodeGUIHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		Info("starting up node gui for parallelcoin")
		Debug(os.Args)
		config.Configure(cx, c.Command.Name, true)
		if cx.ActiveNet.Name == netparams.TestNet3Params.Name {
			fork.IsTestnet = true
		}
		interrupt.AddHandler(func() {
			Debug("node gui is shut down")
			os.Exit(0)
		})
		if err := gui.Main(cx, c); Check(err) {
		}
		Debug("node gui finished")
		return
	}
}
