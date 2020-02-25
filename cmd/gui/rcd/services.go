package rcd

import (
	"fmt"
	"github.com/p9c/pod/cmd/node"
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/wallet"
	"os"
	"sync"
	"sync/atomic"
)

func (r *RcVar) StartServices() (err error) {
	nodeChan := make(chan *rpc.Server)
	// Start Node
	err = r.DuoNodeService(nodeChan)
	if err != nil {
		log.ERROR(err)
	}
	log.DEBUG("waiting for nodeChan")
	r.cx.RPCServer = <-nodeChan
	log.DEBUG("nodeChan sent")
	r.cx.Node.Store(true)

	walletChan := make(chan *wallet.Wallet)
	// Start wallet
	err = r.DuoWalletService(walletChan)
	if err != nil {
		log.ERROR(err)
	}
	log.DEBUG("waiting for walletChan")
	r.cx.WalletServer = <-walletChan
	log.DEBUG("walletChan sent")
	r.cx.Wallet.Store(true)
	//r.Boot.IsBoot = false
	r.Ready <- struct{}{}
	return
}

func (r *RcVar) DuoWalletService(walletChan chan *wallet.Wallet) error {
	r.cx.WalletKill = make(chan struct{})
	r.cx.Wallet = &atomic.Value{}
	r.cx.Wallet.Store(false)
	var err error
	var wg sync.WaitGroup
	if !*r.cx.Config.WalletOff {
		go func() {
			log.INFO("starting wallet")
			//utils.GetBiosMessage(view, "starting wallet")
			err = walletmain.Main(r.cx.Config, r.cx.StateCfg,
				r.cx.ActiveNet, walletChan, r.cx.WalletKill, &wg)
			if err != nil {
				fmt.Println("error running wallet:", err)
				os.Exit(1)
			}
		}()
	}
	interrupt.AddHandler(func() {
		log.WARN("interrupt received, " +
			"shutting down shell modules")
		close(r.cx.WalletKill)
	})
	return err
}

func (r *RcVar) DuoNodeService(nodeChan chan *rpc.Server) error {
	r.cx.NodeKill = make(chan struct{})
	r.cx.Node = &atomic.Value{}
	r.cx.Node.Store(false)
	var err error
	var wg sync.WaitGroup
	if !*r.cx.Config.NodeOff {
		go func() {
			log.INFO(r.cx.Language.RenderText("goApp_STARTINGNODE"))
			//utils.GetBiosMessage(view, cx.Language.RenderText("goApp_STARTINGNODE"))
			err = node.Main(r.cx, nil, r.cx.NodeKill, nodeChan, &wg)
			if err != nil {
				log.INFO("error running node:", err)
				os.Exit(1)
			}
		}()

	}
	interrupt.AddHandler(func() {
		log.WARN("interrupt received, " +
			"shutting down node")
		close(r.cx.NodeKill)
	})
	return err
}
