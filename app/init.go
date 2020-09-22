package app

import (
	"github.com/p9c/pod/app/config"
	"github.com/stalker-loki/app/slog"
	"os"
	"os/exec"

	"github.com/urfave/cli"

	"github.com/p9c/pod/app/conte"
)

var initHandle = func(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		slog.Info("running configuration and wallet initialiser")
		config.Configure(cx, c.Command.Name, true)
		command := os.Args[0]
		args := append(os.Args[1:len(os.Args)-1], "wallet")
		slog.Debug(args)
		firstWallet := exec.Command(command, args...)
		firstWallet.Stdin = os.Stdin
		firstWallet.Stdout = os.Stdout
		firstWallet.Stderr = os.Stderr
		err = firstWallet.Run()
		slog.Debug("running it a second time for mining addresses")
		firstWallet = exec.Command(command, args...)
		firstWallet.Stdin = os.Stdin
		firstWallet.Stdout = os.Stdout
		firstWallet.Stderr = os.Stderr
		err = firstWallet.Run()
		slog.Info("you should be ready to go to sync and mine on the network:", cx.ActiveNet.Name,
			"using datadir:", *cx.Config.DataDir)
		return err
	}
}
