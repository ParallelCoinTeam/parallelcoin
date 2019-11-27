// +build headless

package app

import (
	"os"

	"github.com/p9c/cli"

	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
)

var guiHandle = func(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		log.WARN("GUI was disabled for this build (server only version)")
		os.Exit(1)
		return nil
	}
}
