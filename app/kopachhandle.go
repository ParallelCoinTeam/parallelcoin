package app

import (
	"os"

	"github.com/p9c/pod/app/config"

	"github.com/urfave/cli"

	"github.com/p9c/pod/cmd/kopach"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/chain/fork"

	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func KopachHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		Info("starting up kopach standalone miner for parallelcoin")
		Debug(os.Args)
		config.Configure(cx, c.Command.Name, true)
		if cx.ActiveNet.Name == netparams.TestNet3Params.Name {
			fork.IsTestnet = true
		}
		// quit := make(qu.C)
		interrupt.AddHandler(func() {
			Debug("Handle interrupt")
			close(cx.KillAll)
			// os.Exit(0)
		})
		err = kopach.Handle(cx)(c)
		<-cx.KillAll
		Debug("kopach main finished")
		return
	}
}
