package app

import (
	"os"

	"github.com/p9c/pod/cmd/explorer"

	"github.com/urfave/cli"

	"github.com/p9c/pod/app/config"
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func explorerHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		Info("Starting up explorer GUI for ParallelCoin network")
		Debug(os.Args)
		config.Configure(cx, c.Command.Name, true)
		if cx.ActiveNet.Name == netparams.TestNet3Params.Name {
			fork.IsTestnet = true
		}
		interrupt.AddHandler(func() {
			Debug("Explorer gui is shut down")
			// os.Exit(0)
		})
		if err := explorer.Main(cx, c); Check(err) {
		}
		Debug("Explorer gui finished")
		return
	}
}
