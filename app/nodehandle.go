package app

import (
	"sync"

	"github.com/urfave/cli"

	"github.com/parallelcointeam/parallelcoin/cmd/node"
	"github.com/parallelcointeam/parallelcoin/cmd/node/rpc"
	"github.com/parallelcointeam/parallelcoin/pkg/conte"
	"github.com/parallelcointeam/parallelcoin/pkg/util/cl"
)

func nodeHandle(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) (err error) {
		WARN("running node handler")
		var wg sync.WaitGroup
		Configure(cx)
		// serviceOptions defines the configuration options for the daemon as a service on Windows.
		type serviceOptions struct {
			ServiceCommand string `short:"s" long:"service" description:"Service command {install, remove, start, stop}"`
		}
		// runServiceCommand is only set to a real function on Windows.  It is
		// used to parse and execute service commands specified via the -s flag.
		var runServiceCommand func(string) error
		// Service options which are only added on Windows.
		//
		serviceOpts := serviceOptions{}
		// Perform service command and exit if specified.  Invalid service
		// commands show an appropriate error.
		// Only runs on Windows since the runServiceCommand function will be nil
		// when not on Windows.
		if serviceOpts.ServiceCommand != "" && runServiceCommand != nil {
			err := runServiceCommand(serviceOpts.ServiceCommand)
			if err != nil {
				cx.Log <- cl.Error{err}
				return err
			}
			return nil
		}
		shutdownChan := make(chan struct{})
		nodeChan := make(chan *rpc.Server)
		killswitch := make(chan struct{})
		go func() {
			err = node.Main(cx, shutdownChan, killswitch, nodeChan, &wg)
			if err != nil {
				ERROR("error starting node ", err)
			}
		}()
		cx.RPCServer = <-nodeChan
		wg.Wait()
		return nil
	}
}
