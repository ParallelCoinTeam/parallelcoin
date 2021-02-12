// +build !headless

package app

import (
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/app/config"
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/gui"
)

func walletGUIHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		// Debug(os.Args)
		config.Configure(cx, c.Command.Name, true)
		// interrupt.AddHandler(func() {
		// 	Debug("wallet gui is shut down")
		// })
		if err = gui.Main(cx, c); Check(err) {
		}
		Debug("pod gui finished")
		return
	}
}
