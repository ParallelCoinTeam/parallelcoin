// +build !headless

package app

import (
	"github.com/p9c/pod/pkg/logg"
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/app/config"
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/gui"
)

func walletGUIHandle(cx *conte.Xt) func(c *cli.Context) (e error) {
	return func(c *cli.Context) (e error) {
		logg.App = c.Command.Name
		config.Configure(cx, c.Command.Name, true)
		// dbg.Ln(os.Args)
		// interrupt.AddHandler(func() {
		// 	dbg.Ln("wallet gui is shut down")
		// })
		if e = gui.Main(cx, c); err.Chk(e) {
		}
		dbg.Ln("pod gui finished")
		return
	}
}
