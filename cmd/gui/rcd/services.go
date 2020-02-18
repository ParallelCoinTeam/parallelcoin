package rcd

import (
	"fmt"
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/wallet"
	"os"
	"sync"
	"sync/atomic"
)

func (r *RcVar) Services(walletChan chan *wallet.Wallet) error {
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
