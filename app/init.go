package app

import (
	"os"
	"os/exec"

	"github.com/p9c/pod/app/config"

	"github.com/urfave/cli"

	"github.com/p9c/pod/app/conte"
)

var initHandle = func(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		Info("running configuration and wallet initialiser")
		config.Configure(cx, c.Command.Name, true)
		args := append(os.Args[1:len(os.Args)-1], "wallet")
		Debug(args)
		var command []string
		command = append(command, os.Args[0])
		command = append(command, args...)
		firstWallet := exec.Command(command[0], command[1:]...)
		firstWallet.Stdin = os.Stdin
		firstWallet.Stdout = os.Stdout
		firstWallet.Stderr = os.Stderr
		err := firstWallet.Run()
		Debug("running it a second time for mining addresses")
		secondWallet := exec.Command(command[0], command[1:]...)
		secondWallet.Stdin = os.Stdin
		secondWallet.Stdout = os.Stdout
		secondWallet.Stderr = os.Stderr
		err = firstWallet.Run()
		Info("you should be ready to go to sync and mine on the network:", cx.ActiveNet.Name,
			"using datadir:", *cx.Config.DataDir)
		return err
	}
}
