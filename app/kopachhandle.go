package app

import (
	"github.com/gookit/color"
	"github.com/p9c/pod/pkg/chaincfg"
	"github.com/p9c/pod/pkg/logg"
	"os"
	
	"github.com/p9c/pod/pkg/util/interrupt"
	
	"github.com/p9c/pod/app/config"
	
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/cmd/kopach"
	"github.com/p9c/pod/pkg/fork"
	
	"github.com/p9c/pod/app/conte"
)

func KopachHandle(cx *conte.Xt) func(c *cli.Context) (e error) {
	return func(c *cli.Context) (e error) {
		logg.AppColorizer = color.Bit24(255, 128, 128, false).Sprint
		logg.App = "kopach"
		I.Ln("starting up kopach standalone miner for parallelcoin")
		D.Ln(os.Args)
		config.Configure(cx, "kopach", true)
		if cx.ActiveNet.Name == chaincfg.TestNet3Params.Name {
			fork.IsTestnet = true
		}
		defer cx.KillAll.Q()
		e = kopach.Handle(cx)(c)
		<-interrupt.HandlersDone
		D.Ln("kopach main finished")
		return
	}
}
