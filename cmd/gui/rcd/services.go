package rcd

import (
	"fmt"
	"os"

	"github.com/p9c/pod/cmd/node"
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func (r *RcVar) StartServices() (err error) {
	Debug("starting up services")
	// Start Node
	err = r.DuoNodeService()
	if err != nil {
		Error(err)
	}
	r.cx.RPCServer = <-r.cx.NodeChan
	r.cx.Node.Store(true)

	// Start wallet
	err = r.DuoWalletService()
	if err != nil {
		Error(err)
	}
	r.cx.WalletServer = <-r.cx.WalletChan
	r.cx.Wallet.Store(true)
	r.cx.WalletServer.Rescan(nil, nil)
	r.Ready <- struct{}{}
	return
}

func (r *RcVar) DuoWalletService() error {
	r.cx.WalletKill = make(chan struct{})
	r.cx.Wallet.Store(false)
	var err error
	if !*r.cx.Config.WalletOff {
		go func() {
			Info("starting wallet")
			// utils.GetBiosMessage(view, "starting wallet")
			err = walletmain.Main(r.cx)
			if err != nil {
				fmt.Println("error running wallet:", err)
				os.Exit(1)
			}
		}()
	}
	interrupt.AddHandler(func() {
		Warn("interrupt received, " +
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
			Info(r.cx.Language.RenderText("goApp_STARTINGNODE"))
			// utils.GetBiosMessage(view, cx.Language.RenderText("goApp_STARTINGNODE"))
			err = node.Main(r.cx, nil)
			if err != nil {
				Info("error running node:", err)
				os.Exit(1)
			}
		}()
	}
	interrupt.AddHandler(func() {
		Warn("interrupt received, " +
			"shutting down node")
		close(r.cx.NodeKill)
	})
	return err
}
