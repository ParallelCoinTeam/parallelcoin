package kopach_worker

import (
	"github.com/gookit/color"
	"github.com/p9c/pod/pkg/blockchain/chaincfg"
	"github.com/p9c/pod/pkg/logg"
	"net/rpc"
	"os"
	
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/kopach/worker"
	"github.com/p9c/pod/pkg/blockchain/fork"
)

func KopachWorkerHandle(cx *conte.Xt) func(c *cli.Context) (e error) {
	return func(c *cli.Context) (e error) {
		logg.AppColorizer = color.Bit24(255, 128, 128, false).Sprint
		logg.App = "worker"
		if len(os.Args) > 3 {
			if os.Args[3] == chaincfg.TestNet3Params.Name {
				fork.IsTestnet = true
			}
		}
		if len(os.Args) > 4 {
			logg.SetLogLevel(os.Args[4])
		}
		D.Ln("miner worker starting")
		w, conn := worker.New(os.Args[2], cx.KillAll, uint64(*cx.Config.UUID))
		// interrupt.AddHandler(
		// 	func() {
		// 		D.Ln("KopachWorkerHandle interrupt")
		// 		// if e := conn.Close(); E.Chk(e) {
		// 		// }
		// 		// quit.Q()
		// 	},
		// )
		e = rpc.Register(w)
		if e != nil {
			D.Ln(e)
			return e
		}
		D.Ln("starting up worker IPC")
		rpc.ServeConn(conn)
		D.Ln("stopping worker IPC")
		// if e := conn.Close(); E.Chk(e) {
		// }
		// quit.Quit()
		D.Ln("finished")
		return nil
	}
}
