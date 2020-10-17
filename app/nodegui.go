package app

import (
	"github.com/urfave/cli"

	"github.com/p9c/pod/app/conte"
)

func nodeGUIHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		Debug("running node gui")
		return
	}
}
