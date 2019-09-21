package addresses

import (
	"github.com/parallelcointeam/parallelcoin/app/save"
	"github.com/parallelcointeam/parallelcoin/cmd/node/state"
	"github.com/parallelcointeam/parallelcoin/pkg/pod"
	"github.com/parallelcointeam/parallelcoin/pkg/wallet"
	wm "github.com/parallelcointeam/parallelcoin/pkg/wallet/addrmgr"
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
		WARN("refilling mining addresses")
		account, err := w.AccountNumber(wm.KeyScopeBIP0044,
			"default")
		if err != nil {
			ERROR("error getting account number ", err,
			)
		}
		for i := 0; i < toMake; i++ {
			addr, err := w.NewAddress(account, wm.KeyScopeBIP0044,
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
				ERROR("error adding new address ", err)
			}
		}
		if save.Pod(cfg) {
			WARN("saved config with new addresses")
		} else {
			ERROR("failed to save config")
		}
	}()
}
