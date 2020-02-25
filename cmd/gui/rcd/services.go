package rcd

import (
	"fmt"
	"os"
	"sync/atomic"
	
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/wallet"
)

func (r *RcVar) Services(walletChan chan *wallet.Wallet) error {
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
	}
	interrupt.AddHandler(func() {
		log.WARN("interrupt received, " +
			"shutting down shell modules")
		close(r.cx.WalletKill)
	})
	return err
}
