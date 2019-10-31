package app

import (
	"github.com/urfave/cli"

	"github.com/p9c/pod/pkg/conte"
)

var guiHandle = func(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		//L.Info("GUI was disabled for this build (server only version)")
		//os.Exit(1)
		return nil
	}
}
