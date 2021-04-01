package kopach_worker

import (
	"github.com/gookit/color"
	"github.com/p9c/pod/pkg/chaincfg"
	"github.com/p9c/pod/pkg/fork"
	"github.com/p9c/log"
	"github.com/p9c/pod/pkg/pod"
	"net/rpc"
	"os"
	
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/cmd/kopach/worker"
)

func KopachWorkerHandle(cx *pod.State) func(c *cli.Context) (e error) {
	return func(c *cli.Context) (e error) {
		log.AppColorizer = color.Bit24(255, 128, 128, false).Sprint
		log.App = "worker"
		if len(os.Args) > 3 {
			if os.Args[3] == chaincfg.TestNet3Params.Name {
				fork.IsTestnet = true
			}
		}
		if len(os.Args) > 4 {
			log.SetLogLevel(os.Args[4])
		}
		D.Ln("miner worker starting")
		w, conn := worker.New(os.Args[2], cx.KillAll, uint64(cx.Config.UUID.V()))
		e = rpc.Register(w)
		if e != nil {
			D.Ln(e)
			return e
		}
		D.Ln("starting up worker IPC")
		rpc.ServeConn(conn)
		D.Ln("stopping worker IPC")
		D.Ln("finished")
		return nil
	}
}
