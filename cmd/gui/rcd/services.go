package rcd

import (
	"fmt"
	"os"
	"sync/atomic"
	
	"github.com/p9c/pod/cmd/node"
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func (r *RcVar) StartServices() (err error) {
	log.DEBUG("starting up services")
	// Start Node
	err = r.DuoNodeService()
	if err != nil {
		log.ERROR(err)
	}
	log.DEBUG("waiting for nodeChan")
	r.cx.RPCServer = <-r.cx.NodeChan
	log.DEBUG("nodeChan sent")
	r.cx.Node.Store(true)
	// Start wallet
	err = r.DuoWalletService()
	if err != nil {
		log.ERROR(err)
	}
	log.DEBUG("waiting for walletChan")
	r.cx.WalletServer = <-r.cx.WalletChan
	log.DEBUG("walletChan sent")
	r.cx.Wallet.Store(true)
	// r.Boot.IsBoot = false
	r.Ready <- struct{}{}
	return
}

func (r *RcVar) DuoWalletService() error {
	r.cx.WalletKill = make(chan struct{})
	r.cx.Wallet = &atomic.Value{}
	r.cx.Wallet.Store(false)
	var err error
	if !*r.cx.Config.WalletOff {
		go func() {
			log.INFO("starting wallet")
			// utils.GetBiosMessage(view, "starting wallet")
			err = walletmain.Main(r.cx)
			if err != nil {
				fmt.Println("error running wallet:", err)
				os.Exit(1)
			}
		}()
		r.cx.WalletServer = <-r.cx.WalletChan
	}
	interrupt.AddHandler(func() {
		log.WARN("interrupt received, " +
			"shutting down shell modules")
		close(r.cx.WalletKill)
	})
	return err
}

func (r *RcVar) DuoNodeService() error {
	r.cx.NodeKill = make(chan struct{})
	r.cx.Node = &atomic.Value{}
	r.cx.Node.Store(false)
	var err error
	if !*r.cx.Config.NodeOff {
		go func() {
			log.INFO(r.cx.Language.RenderText("goApp_STARTINGNODE"))
			// utils.GetBiosMessage(view, cx.Language.RenderText("goApp_STARTINGNODE"))
			err = node.Main(r.cx, nil)
			if err != nil {
				log.INFO("error running node:", err)
				os.Exit(1)
			}
		}()
		r.cx.RPCServer = <-r.cx.NodeChan
	}
	interrupt.AddHandler(func() {
		log.WARN("interrupt received, " +
			"shutting down node")
		close(r.cx.NodeKill)
	})
	return err
}
