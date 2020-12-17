package kopach_worker

import (
	"net/rpc"
	"os"
	
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/kopach/worker"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/chain/fork"
	log "github.com/p9c/pod/pkg/util/logi"
)

func KopachWorkerHandle(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		// we take one parameter, name of the network, as this does not change during the lifecycle of the miner worker
		// and is required to get the correct hash functions due to differing hard fork heights. A misconfigured miner
		// will use the wrong hash functions so the controller will log an error and this should be part of any miner
		// control or GUI interface built with pod. Since mainnet is over 200k at writing, mining set to testnet will be
		// correct for mainnet anyway, it is only the other way around that there could be problems with testnet
		// probably never as high as this and hard fork activates early for testing as pre-hardfork doesn't need testing
		// or CPU mining.
		if len(os.Args) > 3 {
			if os.Args[3] == netparams.TestNet3Params.Name {
				fork.IsTestnet = true
			}
		}
		if len(os.Args) > 4 {
			log.L.SetLevel(os.Args[4], true, "pod")
		}
		Debug("miner worker starting")
		w, conn := worker.New(os.Args[2], cx.KillAll)
		// interrupt.AddHandler(
		// 	func() {
		// 		Debug("KopachWorkerHandle interrupt")
		// 		// if err := conn.Close(); Check(err) {
		// 		// }
		// 		// quit.Q()
		// 	},
		// )
		err := rpc.Register(w)
		if err != nil {
			Debug(err)
			return err
		}
		Debug("starting up worker IPC")
		rpc.ServeConn(conn)
		Debug("stopping worker IPC")
		// if err := conn.Close(); Check(err) {
		// }
		// quit.Quit()
		Debug("finished")
		return nil
	}
}
