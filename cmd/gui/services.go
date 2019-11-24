package gui

import (
	"fmt"
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/wallet"
	"os"
	"sync"
	"sync/atomic"
)

func Services(cx *conte.Xt) error {
	walletChan := make(chan *wallet.Wallet)
	cx.WalletKill = make(chan struct{})
	cx.Wallet = &atomic.Value{}
	cx.Wallet.Store(false)
	var err error
	var wg sync.WaitGroup
	if !*cx.Config.WalletOff {
		go func() {
			log.INFO("starting wallet")
			//utils.GetBiosMessage(view, "starting wallet")
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
	})
	return err
}
