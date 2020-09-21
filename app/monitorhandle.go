// +build !headless

package app

import (
	"github.com/urfave/cli"

	"github.com/stalker-loki/pod/app/config"

	"github.com/stalker-loki/pod/app/conte"
	"github.com/stalker-loki/pod/cmd/gui/rcd"
	"github.com/stalker-loki/pod/cmd/monitor"
)

var monitorHandle = func(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		config.Configure(cx, c.Command.Name, true)
		rc := rcd.RcInit(cx)
		Warn("starting monitor GUI")
		return monitor.Run(cx, rc)
	}
}
