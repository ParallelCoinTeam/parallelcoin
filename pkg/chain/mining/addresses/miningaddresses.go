package addresses

import (
	"github.com/stalker-loki/app/slog"
	"github.com/stalker-loki/pod/app/save"
	"github.com/stalker-loki/pod/cmd/node/state"
	"github.com/stalker-loki/pod/pkg/pod"
	"github.com/stalker-loki/pod/pkg/wallet"
	wm "github.com/stalker-loki/pod/pkg/wallet/addrmgr"
)

func RefillMiningAddresses(w *wallet.Wallet, cfg *pod.Config, stateCfg *state.Config) {
	// we make the list up to 1000 so the user does not have to attend to
	// this too often
	miningAddressLen := len(*cfg.MiningAddrs)
	toMake := 999 - miningAddressLen
	if miningAddressLen >= 999 {
		toMake = 18
	}
	if toMake < 3 {
		return
	}
	slog.Warn("refilling mining addresses")
	account, err := w.AccountNumber(wm.KeyScopeBIP0044,
		"default")
	if err != nil {
		slog.Error("error getting account number ", err)
	}
	for i := 0; i < toMake; i++ {
		addr, err := w.NewAddress(account, wm.KeyScopeBIP0044,
			true)
		if err == nil {
			// add them to the configuration to be saved
			*cfg.MiningAddrs = append(*cfg.MiningAddrs, addr.EncodeAddress())
			// add them to the active mining address list so they
			// are ready to use
			stateCfg.ActiveMiningAddrs = append(stateCfg.
				ActiveMiningAddrs, addr)
		} else {
			slog.Error("error adding new address ", err)
		}
	}
	if save.Pod(cfg) {
		slog.Warn("saved config with new addresses")
		// Info("you can now start up a node in the same config folder with fresh addresses ready to mine with")
		// os.Exit(0)
	} else {
		slog.Error("error adding new addresses", err)
	}
}
