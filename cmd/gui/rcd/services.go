package rcd

import (
	"fmt"
	"os"
	
	"github.com/p9c/pod/cmd/node"
	"github.com/p9c/pod/cmd/walletmain"
	log "github.com/p9c/logi"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func (r *RcVar) StartServices() (err error) {
	log.L.Debug("starting up services")
	// Start Node
	err = r.DuoNodeService()
	if err != nil {
		log.L.Error(err)
	}
	
	// Start wallet
	err = r.DuoWalletService()
	if err != nil {
		log.L.Error(err)
	}
	// r.Boot.IsBoot = false
	r.Ready <- struct{}{}
	return
}

func (r *RcVar) DuoWalletService() error {
	r.cx.WalletKill = make(chan struct{})
	r.cx.Wallet.Store(false)
	var err error
	if !*r.cx.Config.WalletOff {
		go func() {
			log.L.Info("starting wallet")
			// utils.GetBiosMessage(view, "starting wallet")
			err = walletmain.Main(r.cx)
			if err != nil {
				fmt.Println("error running wallet:", err)
				os.Exit(1)
			}
		}()
		r.cx.WalletServer = <-r.cx.WalletChan
		r.cx.Wallet.Store(true)
	}
	interrupt.AddHandler(func() {
		log.L.Warn("interrupt received, " +
			"shutting down shell modules")
		close(r.cx.WalletKill)
	})
	return err
}

func (r *RcVar) DuoNodeService() error {
	r.cx.NodeKill = make(chan struct{})
	r.cx.Node.Store(false)
	var err error
	if !*r.cx.Config.NodeOff {
		go func() {
			log.L.Info(r.cx.Language.RenderText("goApp_STARTINGNODE"))
			// utils.GetBiosMessage(view, cx.Language.RenderText("goApp_STARTINGNODE"))
			err = node.Main(r.cx, nil)
			if err != nil {
				log.L.Info("error running node:", err)
				os.Exit(1)
			}
		}()
		r.cx.RPCServer = <-r.cx.NodeChan
		r.cx.Node.Store(true)
	}
	interrupt.AddHandler(func() {
		log.L.Warn("interrupt received, " +
			"shutting down node")
		close(r.cx.NodeKill)
	})
	return err
}
