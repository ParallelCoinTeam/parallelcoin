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
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/wallet"
)

var guiHandle = func(d *core.DuOS) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		log.WARN("starting __OLDgui")
		Configure(d.CtX)
		shutdownChan := make(chan struct{})
		walletChan := make(chan *wallet.Wallet)
		nodeChan := make(chan *rpc.Server)
		d.CtX.WalletKill = make(chan struct{})
		d.CtX.NodeKill = make(chan struct{})
		d.CtX.Wallet = &atomic.Value{}
		d.CtX.Wallet.Store(false)
		d.CtX.Node = &atomic.Value{}
		d.CtX.Node.Store(false)
		var err error
		var wg sync.WaitGroup
		if !*d.CtX.Config.NodeOff {
			go func() {
				log.INFO("starting node")
				err = node.Main(d.CtX, shutdownChan, d.CtX.NodeKill, nodeChan, &wg)
				if err != nil {
		log.ERROR(err)
fmt.Println("error running node:", err)
					os.Exit(1)
				}
			}()
			log.DEBUG("waiting for nodeChan")
			d.CtX.RPCServer = <-nodeChan
			log.DEBUG("nodeChan sent")
			d.CtX.Node.Store(true)
		}
		if !*d.CtX.Config.WalletOff {
			go func() {
				log.INFO("starting wallet")
				err = walletmain.Main(d.CtX.Config, d.CtX.StateCfg,
					d.CtX.ActiveNet, walletChan, d.CtX.WalletKill, &wg)
				if err != nil {
		log.ERROR(err)
fmt.Println("error running wallet:", err)
					os.Exit(1)
				}
			}()
			log.DEBUG("waiting for walletChan")
			d.CtX.WalletServer = <-walletChan
			log.DEBUG("walletChan sent")
			d.CtX.Wallet.Store(true)
		}
		interrupt.AddHandler(func() {
			log.WARN("interrupt received, " +
				"shutting down shell modules")
			close(d.CtX.WalletKill)
			close(d.CtX.NodeKill)
		})
		gui(d)
		if !d.CtX.Node.Load().(bool) {
			close(d.CtX.WalletKill)
		}
		if !d.CtX.Wallet.Load().(bool) {
			close(d.CtX.NodeKill)
		}
		return err
	}
}
