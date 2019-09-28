package app

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"

	"github.com/p9c/pod/cmd/ctl"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
)

const slash = string(os.PathSeparator)

func
ctlHandleList(c *cli.Context) error {
	fmt.Println("Here are the available commands. Pausing a moment as it is a long list...")
	time.Sleep(2 * time.Second)
	ctl.ListCommands()
	return nil
}

func
ctlHandle(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		Configure(cx)
		args := c.Args()
		if len(args) < 1 {
			return cli.ShowSubcommandHelp(c)
		}
		ctl.HelpPrint = func() {
			err := cli.ShowSubcommandHelp(c)
			if err != nil {
				fmt.Println(err)
			}
		}
		ctl.Main(args, cx)
		return nil
	}
}

func
kopachHandle(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		log.WARN("not implemented")
		// Configure(cx)
		// kopachQuit := make(chan struct{})
		// interrupt.AddHandler(func() { close(kopachQuit) })
		// kopach.NewWorker(*cx.Config.MinerListener, *cx.Config.MinerPass, *cx.Config.KopachListener, *cx.Config.DataDir,
		// 	*cx.Config.KopachBias, *cx.Config.GenThreads, kopachQuit)
		return nil
	}
}
