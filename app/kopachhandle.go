package app

import (
	"github.com/p9c/pod/app/config"
	"github.com/p9c/pkg/app/slog"
	"os"

	"github.com/urfave/cli"

	"github.com/p9c/pod/cmd/kopach"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/chain/fork"

	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func KopachHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		slog.Info("starting up kopach standalone miner for parallelcoin")
		config.Configure(cx, c.Command.Name, true)
		if cx.ActiveNet.Name == netparams.TestNet3Params.Name {
			fork.IsTestnet = true
		}
		quit := make(chan struct{})
		interrupt.AddHandler(func() {
			slog.Debug("KopachHandle interrupt")
			close(quit)
			os.Exit(0)
		})
		err = kopach.Handle(cx)(c)
		<-quit
		slog.Debug("kopach main finished")
		return
	}
}
