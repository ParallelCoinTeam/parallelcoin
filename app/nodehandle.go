package app

import (
	"github.com/p9c/pod/app/config"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/util/logi/serve"
	"github.com/urfave/cli"

	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/node"
)

func nodeHandle(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) (err error) {
		Trace("running node handler")
		serve.Log(cx.KillAll, save.Filters(*cx.Config.DataDir))
		config.Configure(cx, c.Command.Name)
		cx.NodeReady = make(chan struct{})
		cx.Node.Store(false)
		// serviceOptions defines the configuration options for the daemon as a service on Windows.
		type serviceOptions struct {
			ServiceCommand string `short:"s" long:"service" description:"Service command {install, remove, start, stop}"`
		}
		// runServiceCommand is only set to a real function on Windows.  It is used to parse and execute service
		// commands specified via the -s flag.
		var runServiceCommand func(string) error
		// Service options which are only added on Windows.
		serviceOpts := serviceOptions{}
		// Perform service command and exit if specified.  Invalid service commands show an appropriate error.
		// Only runs on Windows since the runServiceCommand function will be nil when not on Windows.
		if serviceOpts.ServiceCommand != "" && runServiceCommand != nil {
			err := runServiceCommand(serviceOpts.ServiceCommand)
			if err != nil {
				Error(err)
				return err
			}
			return nil
		}
		shutdownChan := make(chan struct{})
		go func() {
			err := node.Main(cx, shutdownChan)
			if err != nil {
				Error("error starting node ", err)
			}
		}()
		Debug("sending back node rpc server handler")
		cx.RPCServer = <-cx.NodeChan
		close(cx.NodeReady)
		cx.Node.Store(true)
		cx.WaitGroup.Wait()
		return nil
	}
}
