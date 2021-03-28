package app

import (
	"github.com/gookit/color"
	"github.com/p9c/log"
	"github.com/p9c/pod/pkg/chaincfg"
	"github.com/p9c/pod/pkg/fork"
	"github.com/p9c/pod/pkg/pod"
	"os"
	
	"github.com/p9c/pod/pkg/util/interrupt"
	
	"github.com/p9c/pod/pkg/podconfig"
	
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/cmd/kopach"
)

func KopachHandle(cx *pod.State) func(c *cli.Context) (e error) {
	return func(c *cli.Context) (e error) {
		log.AppColorizer = color.Bit24(255, 128, 128, false).Sprint
		log.App = "kopach"
		I.Ln("starting up kopach standalone miner for parallelcoin")
		D.Ln(os.Args)
		podconfig.Configure(cx, "kopach", true)
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
