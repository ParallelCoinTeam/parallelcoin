// +build headless

package app

import (
	"github.com/stalker-loki/app/slog"
	"os"

	"github.com/urfave/cli"

	"github.com/p9c/pod/app/conte"
)

var guiHandle = func(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		slog.Warn("GUI was disabled for this build (server only version)")
		os.Exit(1)
		return nil
	}
}

var monitorHandle = func(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		slog.Warn("GUI was disabled for this build (server only version)")
		os.Exit(1)
		return nil
	}
}
