package app

import (
	"github.com/urfave/cli"

	"github.com/p9c/pod/app/config"
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/gui"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func walletGUIHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		Info("starting up node gui for parallelcoin")
		// Debug(os.Args)
		config.Configure(cx, c.Command.Name, true)
		interrupt.AddHandler(func() {
			Debug("wallet gui is shut down")
			// os.Exit(0)
		})
		if err := gui.Main(cx, c); Check(err) {
		}
		Debug("pod gui finished")
		return
	}
}
