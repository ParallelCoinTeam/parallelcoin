// +build !headless

package app

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/urfave/cli"

	"github.com/parallelcointeam/parallelcoin/cmd/gui"
	"github.com/parallelcointeam/parallelcoin/cmd/node"
	"github.com/parallelcointeam/parallelcoin/cmd/node/rpc"
	"github.com/parallelcointeam/parallelcoin/cmd/walletmain"
	"github.com/parallelcointeam/parallelcoin/pkg/conte"
	"github.com/parallelcointeam/parallelcoin/pkg/log"
	"github.com/parallelcointeam/parallelcoin/pkg/util/interrupt"
	"github.com/parallelcointeam/parallelcoin/pkg/wallet"
)

var guiHandle = func(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		log.WARN("starting gui")
		Configure(cx)
		shutdownChan := make(chan struct{})
		walletChan := make(chan *wallet.Wallet)
		nodeChan := make(chan *rpc.Server)
		cx.WalletKill = make(chan struct{})
		cx.NodeKill = make(chan struct{})
		cx.Wallet = &atomic.Value{}
		cx.Wallet.Store(false)
		cx.Node = &atomic.Value{}
		cx.Node.Store(false)
		var err error
		var wg sync.WaitGroup
		if !*cx.Config.NodeOff {
			go func() {
				log.INFO("starting node")
				err = node.Main(cx, shutdownChan, cx.NodeKill, nodeChan, &wg)
				if err != nil {
					fmt.Println("error running node:", err)
					os.Exit(1)
				}
			}()
			log.DEBUG("waiting for nodeChan")
			cx.RPCServer = <-nodeChan
			log.DEBUG("nodeChan sent")
			cx.Node.Store(true)
		}
		if !*cx.Config.WalletOff {
			go func() {
				log.INFO("starting wallet")
				err = walletmain.Main(cx.Config, cx.StateCfg,
					cx.ActiveNet, walletChan, cx.WalletKill, &wg)
				if err != nil {
					fmt.Println("error running wallet:", err)
					os.Exit(1)
				}
			}()
			log.DEBUG("waiting for walletChan")
			cx.WalletServer = <-walletChan
			log.DEBUG("walletChan sent")
			cx.Wallet.Store(true)
		}
		interrupt.AddHandler(func() {
			log.WARN("interrupt received, " +
				"shutting down shell modules")
			close(cx.WalletKill)
			close(cx.NodeKill)
		})
		gui.Main(cx, &wg)
		if !cx.Node.Load().(bool) {
			close(cx.WalletKill)
		}
		if !cx.Wallet.Load().(bool) {
			close(cx.NodeKill)
		}
		return err
	}
}
