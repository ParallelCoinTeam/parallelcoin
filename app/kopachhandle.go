package app

import (
	"github.com/p9c/pod/cmd/kopach"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/log"
	"github.com/urfave/cli"
	"os"

	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func KopachHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		log.INFO("starting up kopach standalone miner for parallelcoin")
		Configure(cx, c)
		if cx.ActiveNet.Name == netparams.TestNet3Params.Name {
			fork.IsTestnet = true
		}
		quit := make(chan struct{})
		interrupt.AddHandler(func() {
			close(quit)
			os.Exit(0)
		})
		err = kopach.KopachHandle(cx)(c)
		<-quit
		log.DEBUG("kopach main finished")
		return
	}
}
//
// func kopachGUIHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
// 	return func(c *cli.Context) (err error) {
// 		log.INFO("starting up kopach standalone miner with GUI for parallelcoin")
// 		Configure(cx, c)
// 		quit := make(chan struct{})
// 		interrupt.AddHandler(func() {
// 			close(quit)
// 			os.Exit(0)
// 		})
// 		err = kopach.KopachHandle(cx)(c)
// 		<-quit
// 		log.DEBUG("kopach main finished")
// 		return
// 	}
// }

