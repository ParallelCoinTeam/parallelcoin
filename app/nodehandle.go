package app

import (
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/app/config"
	"github.com/p9c/pod/cmd/walletmain"
	qu "github.com/p9c/pod/pkg/util/quit"
	
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/node"
)

func nodeHandle(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) (err error) {
		Trace("running node handler")
		config.Configure(cx, c.Command.Name, true)
		cx.NodeReady = qu.T()
		cx.Node.Store(false)
		// serviceOptions defines the configuration options for the daemon as a service on Windows.
		type serviceOptions struct {
			ServiceCommand string `short:"s" long:"service" description:"Service command {install, remove, start, stop}"`
		}
		// runServiceCommand is only set to a real function on Windows. It is used to parse and execute service commands
		// specified via the -s flag.
		runServiceCommand := func(string) error { return nil }
		// Service options which are only added on Windows.
		serviceOpts := serviceOptions{}
		// Perform service command and exit if specified. Invalid service commands show an appropriate error. Only runs
		// on Windows since the runServiceCommand function will be nil when not on Windows.
		if serviceOpts.ServiceCommand != "" && runServiceCommand != nil {
			err := runServiceCommand(serviceOpts.ServiceCommand)
			if err != nil {
				Error(err)
				return err
			}
			return nil
		}
		config.Configure(cx, c.Command.Name, true)
		Debug("starting shell")
		if *cx.Config.TLS || *cx.Config.ServerTLS {
			// generate the tls certificate if configured
			if apputil.FileExists(*cx.Config.RPCCert) && apputil.FileExists(*cx.Config.RPCKey) &&
				apputil.FileExists(*cx.Config.CAFile) {
				
			} else {
				_, _ = walletmain.GenerateRPCKeyPair(cx.Config, true)
			}
		}
		if !*cx.Config.NodeOff {
			go func() {
				err := node.Main(cx)
				if err != nil {
					Error("error starting node ", err)
				}
			}()
			Info("starting node")
			if !*cx.Config.DisableRPC {
				cx.RPCServer = <-cx.NodeChan
				cx.NodeReady.Q()
				cx.Node.Store(true)
				Info("node started")
			}
		}
		cx.WaitWait()
		Info("node is now fully shut down")
		cx.WaitGroup.Wait()
		<-cx.KillAll
		return nil
	}
}
