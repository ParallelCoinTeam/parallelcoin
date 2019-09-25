package addresses

import (
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/cmd/node/state"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/util/cl"
	"github.com/p9c/pod/pkg/wallet"
	waddrmgr "github.com/p9c/pod/pkg/wallet/addrmgr"
)

func RefillMiningAddresses(w *wallet.Wallet, cfg *pod.Config, stateCfg *state.Config) {
	go func() {
		// we make the list up to 1000 so the user does not have to attend to
		// this too often
		miningAddressLen := len(*cfg.MiningAddrs)
		toMake := 1000 - miningAddressLen
		if toMake < 1 {
			return
		}
		log <- cl.Warn{"refilling mining addresses", cl.Ine()}
		account, err := w.AccountNumber(waddrmgr.KeyScopeBIP0044,
			"default")
		if err != nil {
			log <- cl.Error{"error getting account number ", err,
				cl.Ine()}
		}
		for i := 0; i < toMake; i++ {
			addr, err := w.NewAddress(account, waddrmgr.KeyScopeBIP0044,
				true)
			if err == nil {
				// add them to the configuration to be saved
				*cfg.MiningAddrs = append(*cfg.MiningAddrs,
					addr.EncodeAddress())
				// add them to the active mining address list so they
				// are ready to use
				stateCfg.ActiveMiningAddrs = append(stateCfg.
					ActiveMiningAddrs, addr)
			} else {
				log <- cl.Error{"error adding new address ", err,
					cl.Ine()}
			}
		}
		if save.Pod(cfg) {
			log <- cl.Warn{"saved config with new addresses", cl.Ine()}
		} else {
			log <- cl.Error{"failed to save config", cl.Ine()}
		}
	}()
}
