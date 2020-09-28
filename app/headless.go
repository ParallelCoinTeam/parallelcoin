// +build headless

package app

import (
	"github.com/stalker-loki/app/slog"
	"os"

	"github.com/urfave/cli"

	"github.com/p9c/pod/app/conte"
)

var guiHandle = func(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		slog.Warn("GUI was disabled for this build (server only version)")
		defer os.Exit(1)
		return
	}
}

var monitorHandle = func(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		slog.Warn("GUI was disabled for this build (server only version)")
		defer os.Exit(1)
		return
	}
}
