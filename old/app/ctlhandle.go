package app

import (
	"github.com/gookit/color"
	"github.com/p9c/pod/pkg/logg"
	"github.com/p9c/pod/pkg/pod"
	"os"
	
	"github.com/p9c/pod/pkg/podconfig"
	
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/cmd/ctl"
)

const slash = string(os.PathSeparator)

func ctlHandleList(c *cli.Context) (e error) {
	// fmt.Println("Here are the available commands. Pausing a moment as it is a long list...")
	// time.Sleep(2 * time.Second)
	ctl.ListCommands()
	return nil
}

func ctlHandle(cx *pod.State) func(c *cli.Context) (e error) {
	return func(c *cli.Context) (e error) {
		logg.AppColorizer = color.Bit24(128, 128, 255, false).Sprint
		logg.App = "   ctl"
		*cx.Config.LogLevel = "off"
		podconfig.Configure(cx, "ctl", true)
		args := c.Args()
		if len(args) < 1 {
			return cli.ShowSubcommandHelp(c)
		}
		ctl.HelpPrint = func() {
			if e := cli.ShowSubcommandHelp(c); E.Chk(e) {
			}
		}
		ctl.Main(args, cx)
		return nil
	}
}
