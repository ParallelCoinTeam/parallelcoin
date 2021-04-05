package mining

import (
	"github.com/p9c/monorepo/pkg/btcaddr"
	"github.com/p9c/monorepo/pkg/podopts"
	
	"github.com/p9c/monorepo/cmd/node/state"
	wm "github.com/p9c/monorepo/pkg/waddrmgr"
	"github.com/p9c/monorepo/pkg/wallet"
)

// RefillMiningAddresses adds new addresses to the mining address pool for the miner
// todo: make this remove ones that have been used or received a payment or mined
func RefillMiningAddresses(w *wallet.Wallet, cfg *podopts.Config, stateCfg *state.Config) {
	if w == nil {
		D.Ln("trying to refill without a wallet")
		return
	}
	if cfg == nil {
		D.Ln("config is empty")
		return
	}
	miningAddressLen := len(cfg.MiningAddrs.S())
	toMake := 99 - miningAddressLen
	if miningAddressLen >= 99 {
		toMake = 0
	}
	if toMake < 1 {
		D.Ln("not making any new addresses")
		return
	}
	D.Ln("refilling mining addresses")
	account, e := w.AccountNumber(
		wm.KeyScopeBIP0044,
		"default",
	)
	if e != nil {
		E.Ln("error getting account number ", e)
	}
	for i := 0; i < toMake; i++ {
		var addr btcaddr.Address
		addr, e = w.NewAddress(
			account, wm.KeyScopeBIP0044,
			true,
		)
		if e == nil {
			// add them to the configuration to be saved
			cfg.MiningAddrs.Set(append(cfg.MiningAddrs.S(), addr.EncodeAddress()))
			// add them to the active mining address list so they
			// are ready to use
			stateCfg.ActiveMiningAddrs = append(stateCfg.ActiveMiningAddrs, addr)
		} else {
			E.Ln("error adding new address ", e)
		}
	}
	if podcfg.Save(cfg) {
		D.Ln("saved config with new addresses")
	} else {
		E.Ln("error adding new addresses", e)
	}
}
