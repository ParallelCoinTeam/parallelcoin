package app

import (
	"os"
	"os/exec"
	
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/pkg/conte"
	log "github.com/p9c/pod/pkg/logi"
)

var initHandle = func(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		log.L.Info("running configuration and wallet initialiser")
		Configure(cx, c)
		command := os.Args[0]
		args := append(os.Args[1:len(os.Args)-1], "wallet")
		log.L.Debug(args)
		firstWallet := exec.Command(command, args...)
		firstWallet.Stdin = os.Stdin
		firstWallet.Stdout = os.Stdout
		firstWallet.Stderr = os.Stderr
		err := firstWallet.Run()
		log.L.Debug("running it a second time for mining addresses")
		firstWallet = exec.Command(command, args...)
		firstWallet.Stdin = os.Stdin
		firstWallet.Stdout = os.Stdout
		firstWallet.Stderr = os.Stderr
		err = firstWallet.Run()
		log.L.Info("you should be ready to go to sync and mine on the network:", cx.ActiveNet.Name,
			"using datadir:", *cx.Config.DataDir)
		return err
	}
}
