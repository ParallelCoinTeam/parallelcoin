// +build !headless

package app

import (
	"github.com/stalker-loki/app/slog"
	"github.com/urfave/cli"

	"github.com/p9c/pod/app/config"

	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/cmd/monitor"
)

var monitorHandle = func(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		config.Configure(cx, c.Command.Name, true)
		rc := rcd.RcInit(cx)
		slog.Warn("starting monitor GUI")
		return monitor.Run(cx, rc)
	}
}
