package app

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/p9c/log"
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

func ctlHandle(ifc interface{}) (e error) {
	var cx *pod.State
	var ok bool
	if cx, ok = ifc.(*pod.State); !ok {
		return fmt.Errorf("cannot run without a state")
	}
	log.AppColorizer = color.Bit24(128, 128, 255, false).Sprint
	log.App = "   ctl"
	cx.Config.LogLevel.Set("off")
	podconfig.Configure(cx, true)
	args := os.Args
	ctl.Main(args, cx)
	return nil
}
