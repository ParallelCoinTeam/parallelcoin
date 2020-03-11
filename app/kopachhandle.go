package app

import (
	"github.com/p9c/pod/cmd/kopach"
	"github.com/p9c/chaincfg/netparams"
	"github.com/p9c/fork"
	log "github.com/p9c/logi"
	"github.com/urfave/cli"
	"os"

	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/util/interrupt"
)

func KopachHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		log.L.Info("starting up kopach standalone miner for parallelcoin")
		Configure(cx, c)
		if cx.ActiveNet.Name == netparams.TestNet3Params.Name {
			fork.IsTestnet = true
		}
		quit := make(chan struct{})
		interrupt.AddHandler(func() {
			log.L.Debug("KopachHandle interrupt")
			close(quit)
			os.Exit(0)
		})
		err = kopach.KopachHandle(cx)(c)
		<-quit
		log.L.Debug("kopach main finished")
		return
	}
}
