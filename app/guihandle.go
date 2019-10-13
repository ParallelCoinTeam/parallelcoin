// +build !headless

package app

import (
	"fmt"
	"github.com/p9c/pod/pkg/duos/core"
	"os"
	"sync"
	"sync/atomic"

	"github.com/urfave/cli"

	"github.com/p9c/pod/cmd/node"
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/gui"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/wallet"
)

var guiHandle = func(d *core.DuOS) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		log.WARN("starting gui")
		Configure(d.Cx)
		shutdownChan := make(chan struct{})
		walletChan := make(chan *wallet.Wallet)
		nodeChan := make(chan *rpc.Server)
		d.Cx.WalletKill = make(chan struct{})
		d.Cx.NodeKill = make(chan struct{})
		d.Cx.Wallet = &atomic.Value{}
		d.Cx.Wallet.Store(false)
		d.Cx.Node = &atomic.Value{}
		d.Cx.Node.Store(false)
		var err error
		var wg sync.WaitGroup
		if !*d.Cx.Config.NodeOff {
			go func() {
				log.INFO("starting node")
				err = node.Main(d.Cx, shutdownChan, d.Cx.NodeKill, nodeChan, &wg)
				if err != nil {
					fmt.Println("error running node:", err)
					os.Exit(1)
				}
			}()
			log.DEBUG("waiting for nodeChan")
			d.Cx.RPCServer = <-nodeChan
			log.DEBUG("nodeChan sent")
			d.Cx.Node.Store(true)
		}
		if !*d.Cx.Config.WalletOff {
			go func() {
				log.INFO("starting wallet")
				err = walletmain.Main(d.Cx.Config, d.Cx.StateCfg,
					d.Cx.ActiveNet, walletChan, d.Cx.WalletKill, &wg)
				if err != nil {
					fmt.Println("error running wallet:", err)
					os.Exit(1)
				}
			}()
			log.DEBUG("waiting for walletChan")
			d.Cx.WalletServer = <-walletChan
			log.DEBUG("walletChan sent")
			d.Cx.Wallet.Store(true)
		}
		interrupt.AddHandler(func() {
			log.WARN("interrupt received, " +
				"shutting down shell modules")
			close(d.Cx.WalletKill)
			close(d.Cx.NodeKill)
		})
		gui.Main(d)
		if !d.Cx.Node.Load().(bool) {
			close(d.Cx.WalletKill)
		}
		if !d.Cx.Wallet.Load().(bool) {
			close(d.Cx.NodeKill)
		}
		return err
	}
}
