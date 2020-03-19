// +build !headless

package app

import (
	"github.com/p9c/pod/app/config"
	"github.com/urfave/cli"

	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/cmd/monitor"
	"github.com/p9c/pod/pkg/conte"
)

var monitorHandle = func(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		config.Configure(cx, c.Command.Name)
		rc := rcd.RcInit(cx)
		L.Warn("starting monitor GUI")
		return monitor.Run(cx, rc)
	}
}
