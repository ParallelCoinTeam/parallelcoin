package app

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"

	"github.com/p9c/pod/cmd/ctl"
	"github.com/p9c/pod/pkg/conte"
)

const slash = string(os.PathSeparator)

func ctlHandleList(c *cli.Context) error {
	fmt.Println("Here are the available commands. Pausing a moment as it is a long list...")
	time.Sleep(2 * time.Second)
	ctl.ListCommands()
	return nil
}

func ctlHandle(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		Configure(cx, c)
		args := c.Args()
		if len(args) < 1 {
			return cli.ShowSubcommandHelp(c)
		}
		ctl.HelpPrint = func() {
			err := cli.ShowSubcommandHelp(c)
			if err != nil {
				L.Error(err)
			}
		}
		ctl.Main(args, cx)
		return nil
	}
}

func ctlGUIHandle(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		Configure(cx, c)

		return nil
	}
}
