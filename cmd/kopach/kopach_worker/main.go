package kopach_worker

import (
	"net/rpc"
	"os"
	
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/kopach/worker"
	"github.com/p9c/pod/pkg/blockchain/chaincfg/netparams"
	"github.com/p9c/pod/pkg/blockchain/fork"
	log "github.com/p9c/pod/pkg/util/logi"
)

func KopachWorkerHandle(cx *conte.Xt) func(c *cli.Context) (e error) {
	return func(c *cli.Context) (e error) {
		if len(os.Args) > 3 {
			if os.Args[3] == netparams.TestNet3Params.Name {
				fork.IsTestnet = true
			}
		}
		if len(os.Args) > 4 {
			log.L.SetLevel(os.Args[4], true, "pod")
		}
		dbg.Ln("miner worker starting")
		w, conn := worker.New(os.Args[2], cx.KillAll, uint64(*cx.Config.UUID))
		// interrupt.AddHandler(
		// 	func() {
		// 		dbg.Ln("KopachWorkerHandle interrupt")
		// 		// if e := conn.Close(); err.Chk(e) {
		// 		// }
		// 		// quit.Q()
		// 	},
		// )
		e = rpc.Register(w)
		if e != nil  {
			dbg.Ln(e)
			return e
		}
		dbg.Ln("starting up worker IPC")
		rpc.ServeConn(conn)
		dbg.Ln("stopping worker IPC")
		// if e := conn.Close(); err.Chk(e) {
		// }
		// quit.Quit()
		dbg.Ln("finished")
		return nil
	}
}
