package app

import (
	"github.com/p9c/pod/pkg/logg"
	"os"
	"os/exec"
	
	"github.com/p9c/pod/app/config"
	
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/app/conte"
)

var initHandle = func(cx *conte.Xt) func(c *cli.Context) (e error) {
	return func(c *cli.Context) (e error) {
		logg.App = "  init"
		inf.Ln("running configuration and wallet initialiser")
		config.Configure(cx, c.Command.Name, true)
		args := append(os.Args[1:len(os.Args)-1], "wallet")
		dbg.Ln(args)
		var command []string
		command = append(command, os.Args[0])
		command = append(command, args...)
		// command = apputil.PrependForWindows(command)
		firstWallet := exec.Command(command[0], command[1:]...)
		firstWallet.Stdin = os.Stdin
		firstWallet.Stdout = os.Stdout
		firstWallet.Stderr = os.Stderr
		e = firstWallet.Run()
		dbg.Ln("running it a second time for mining addresses")
		secondWallet := exec.Command(command[0], command[1:]...)
		secondWallet.Stdin = os.Stdin
		secondWallet.Stdout = os.Stdout
		secondWallet.Stderr = os.Stderr
		e = firstWallet.Run()
		inf.Ln("you should be ready to go to sync and mine on the network:", cx.ActiveNet.Name,
			"using datadir:", *cx.Config.DataDir)
		return e
	}
}
