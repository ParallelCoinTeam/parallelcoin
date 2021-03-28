// +podbuild !headless

package app

import (
	"github.com/gookit/color"
	"github.com/p9c/log"
	"github.com/p9c/pod/pkg/pod"
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/cmd/gui"
	"github.com/p9c/pod/pkg/podconfig"
)

func walletGUIHandle(cx *pod.State) func(c *cli.Context) (e error) {
	return func(c *cli.Context) (e error) {
		log.AppColorizer = color.Bit24(128, 255, 255, false).Sprint
		log.App = "   gui"
		D.Ln("starting up parallelcoin pod gui...")
		// fork.ForkCalc()
		podconfig.Configure(cx, "gui", true)
		// D.Ln(os.Args)
		// interrupt.AddHandler(func() {
		// 	D.Ln("wallet gui is shut down")
		// })
		if e = gui.Main(cx, c); E.Chk(e) {
		}
		D.Ln("pod gui finished")
		return
	}
}
