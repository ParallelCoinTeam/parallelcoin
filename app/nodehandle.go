package app

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/p9c/log"
	"github.com/p9c/pod/pkg/opts"
	"github.com/p9c/pod/pkg/pod"
	
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/apputil"
	"github.com/p9c/pod/pkg/podconfig"
	"github.com/p9c/qu"
	
	"github.com/p9c/pod/cmd/node"
)

func nodeHandle(ifc interface{}) (e error) {
	var cx *pod.State
	var ok bool
	if cx, ok = ifc.(*pod.State); !ok {
		return fmt.Errorf("cannot run without a state")
	}
	log.AppColorizer = color.Bit24(128, 128, 255, false).Sprint
	log.App = "  node"
	opts.F.Ln("running node handler")
	podconfig.Configure(cx, true)
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
		if e = runServiceCommand(serviceOpts.ServiceCommand); opts.E.Chk(e) {
			return e
		}
		return nil
	}
	// config.Configure(cx, c.Command.Name, true)
	// D.Ln("starting shell")
	if cx.Config.ClientTLS.True() || cx.Config.ServerTLS.True() {
		// generate the tls certificate if configured
		if apputil.FileExists(cx.Config.RPCCert.V()) &&
			apputil.FileExists(cx.Config.RPCKey.V()) &&
			apputil.FileExists(cx.Config.CAFile.V()) {
		} else {
			if _, e = walletmain.GenerateRPCKeyPair(cx.Config, true); opts.E.Chk(e) {
			}
		}
	}
	// if cx.Config.NodeOff.False() {
	// go func() {
	if e := node.Main(cx); opts.E.Chk(e) {
		opts.E.Ln("error starting node ", e)
	}
	// }()
	opts.I.Ln("starting node")
	if cx.Config.DisableRPC.False() {
		cx.RPCServer = <-cx.NodeChan
		cx.NodeReady.Q()
		cx.Node.Store(true)
		opts.I.Ln("node started")
	}
	// }
	cx.WaitWait()
	opts.I.Ln("node is now fully shut down")
	cx.WaitGroup.Wait()
	<-cx.KillAll
	return nil
}
