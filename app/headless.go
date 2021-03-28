// +podbuild headless

package app

import (
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/log"
	"os"
	
	"github.com/urfave/cli"
)

var walletGUIHandle = func(cx *pod.State) func(c *cli.Context) (e error) {
	return func(c *cli.Context) (e error) {
		log.App = c.Command.Name
		W.Ln("GUI was disabled for this podbuild (server only version)")
		os.Exit(1)
		return nil
	}
}
