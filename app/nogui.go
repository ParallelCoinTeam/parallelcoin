// +build headless

package app

import (
	"fmt"
	"os"
	
	"github.com/urfave/cli"
	
	"github.com/parallelcointeam/parallelcoin/pkg/conte"
)

var guiHandle = func(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		fmt.Println("GUI was disabled for this build (server only version)")
		os.Exit(1)
		return nil
	}
}
