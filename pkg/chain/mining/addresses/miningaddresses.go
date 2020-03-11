package addresses

import (
	"os"
	
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/cmd/node/state"
	log "github.com/p9c/logi"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/wallet"
	wm "github.com/p9c/pod/pkg/wallet/addrmgr"
)

func RefillMiningAddresses(w *wallet.Wallet, cfg *pod.Config, stateCfg *state.Config) {
	// we make the list up to 1000 so the user does not have to attend to
	// this too often
	miningAddressLen := len(*cfg.MiningAddrs)
	toMake := 999 - miningAddressLen
	if toMake < 3 {
		return
	}
	log.L.Warn("refilling mining addresses")
	account, err := w.AccountNumber(wm.KeyScopeBIP0044,
		"default")
	if err != nil {
		log.L.Error("error getting account number ", err)
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
			log.L.Error("error adding new address ", err)
		}
	}
	if save.Pod(cfg) {
		log.L.Warn("saved config with new addresses")
		log.L.Info("you can now start up a node in the same config folder with fresh addresses ready to mine with")
		os.Exit(0)
	} else {
		log.L.Error("error adding new addresses", err)
	}
}
