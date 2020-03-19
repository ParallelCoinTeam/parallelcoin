package app

import (
	"os"

	"github.com/urfave/cli"

	"github.com/p9c/pod/cmd/kopach"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/chain/fork"

	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func KopachHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		L.Info("starting up kopach standalone miner for parallelcoin")
		Configure(cx, c.Command.Name)
		if cx.ActiveNet.Name == netparams.TestNet3Params.Name {
			fork.IsTestnet = true
		}
		quit := make(chan struct{})
		interrupt.AddHandler(func() {
			L.Debug("KopachHandle interrupt")
			close(quit)
			os.Exit(0)
		})
		err = kopach.KopachHandle(cx)(c)
		<-quit
		L.Debug("kopach main finished")
		return
	}
}
