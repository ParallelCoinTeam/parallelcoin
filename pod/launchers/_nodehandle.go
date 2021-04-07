package launchers

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/p9c/monorepo/node/node"
	"github.com/p9c/monorepo/pkg/log"
	"github.com/p9c/monorepo/pkg/pod"
	
	"github.com/p9c/monorepo/cmd/walletmain"
	"github.com/p9c/monorepo/pkg/apputil"
	"github.com/p9c/monorepo/pkg/qu"
	"github.com/p9c/monorepo/pod/podconfig"
)

func nodeHandle(ifc interface{}) (e error) {
	var cx *pod.State
	var ok bool
	if cx, ok = ifc.(*pod.State); !ok {
		return fmt.Errorf("cannot run without a state")
	}
	log.AppColorizer = color.Bit24(128, 128, 255, false).Sprint
	log.App = "  node"
	F.Ln("running node handler")
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
		if e = runServiceCommand(serviceOpts.ServiceCommand); E.Chk(e) {
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
			if _, e = walletmain.GenerateRPCKeyPair(cx.Config, true); E.Chk(e) {
			}
		}
	}
	// if cx.Config.NodeOff.False() {
	// go func() {
	if e := node.Main(cx); E.Chk(e) {
		E.Ln("error starting node ", e)
	}
	// }()
	I.Ln("starting node")
	if cx.Config.DisableRPC.False() {
		cx.RPCServer = <-cx.NodeChan
		cx.NodeReady.Q()
		cx.Node.Store(true)
		I.Ln("node started")
	}
	// }
	cx.WaitWait()
	I.Ln("node is now fully shut down")
	cx.WaitGroup.Wait()
	<-cx.KillAll
	return nil
}
