package app

import (
	"github.com/gookit/color"
	"github.com/p9c/pod/pkg/logg"
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/app/config"
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/util/qu"
	
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/node"
)

func rpcNodeHandle(cx *conte.Xt) func(c *cli.Context) (e error) {
	*cx.Config.DisableController = true
	return nodeHandle(cx)
}

func nodeHandle(cx *conte.Xt) func(c *cli.Context) (e error) {
	return func(c *cli.Context) (e error) {
		logg.App=color.Bit24(0,0,255,false).Sprint("  node")
		trc.Ln("running node handler")
		config.Configure(cx, c.Command.Name, true)
		cx.NodeReady = qu.T()
		cx.Node.Store(false)
		// serviceOptions defines the configuration options for the daemon as a service on Windows.
		type serviceOptions struct {
			ServiceCommand string `short:"s" long:"service" description:"Service command {install, remove, start, stop}"`
		}
		// runServiceCommand is only set to a real function on Windows. It is used to parse and execute service commands
		// specified via the -s flag.
		runServiceCommand := func(string) (e error) { return nil }
		// Service options which are only added on Windows.
		serviceOpts := serviceOptions{}
		// Perform service command and exit if specified. Invalid service commands show an appropriate error. Only runs
		// on Windows since the runServiceCommand function will be nil when not on Windows.
		if serviceOpts.ServiceCommand != "" && runServiceCommand != nil {
			if e = runServiceCommand(serviceOpts.ServiceCommand); err.Chk(e) {
				return e
			}
			return nil
		}
		config.Configure(cx, c.Command.Name, true)
		dbg.Ln("starting shell")
		if *cx.Config.TLS || *cx.Config.ServerTLS {
			// generate the tls certificate if configured
			if apputil.FileExists(*cx.Config.RPCCert) &&
				apputil.FileExists(*cx.Config.RPCKey) &&
				apputil.FileExists(*cx.Config.CAFile) {
			} else {
				if _, e = walletmain.GenerateRPCKeyPair(cx.Config, true); err.Chk(e) {
				}
			}
		}
		if !*cx.Config.NodeOff {
			go func() {
				if e := node.Main(cx); err.Chk(e) {
					err.Ln("error starting node ", e)
				}
			}()
			inf.Ln("starting node")
			if !*cx.Config.DisableRPC {
				cx.RPCServer = <-cx.NodeChan
				cx.NodeReady.Q()
				cx.Node.Store(true)
				inf.Ln("node started")
			}
		}
		cx.WaitWait()
		inf.Ln("node is now fully shut down")
		cx.WaitGroup.Wait()
		<-cx.KillAll
		return nil
	}
}
